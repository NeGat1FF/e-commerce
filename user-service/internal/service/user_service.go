package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/NeGat1FF/e-commerce/user-service/internal/models"
	"github.com/NeGat1FF/e-commerce/user-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/user-service/internal/utils"
	"github.com/NeGat1FF/e-commerce/user-service/pkg/logger"
	"github.com/NeGat1FF/e-commerce/user-service/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrInvalidJWT              = errors.New("invalid JWT")
	ErrInvalidEmailOrPassword  = errors.New("invalid email or password")
	ErrInternalServer          = errors.New("internal server error")
	ErrInvalidEmail            = errors.New("invalid email")
	ErrInvalidPasswordLength   = errors.New("password length must be between 8 and 64 characters")
	ErrInvalidFaildToSendEmail = errors.New("failed to send email")
	ErrEmailAlreadyVerified    = errors.New("email is already verified")
)

// UserService describes the service.
type UserService struct {
	repo        repository.UserRepositoryInterface
	mailService proto.MailServiceClient
	jwtSecret   string
	addr        string
	port        string
}

// NewUserService creates a new user service.
func NewUserService(repo repository.UserRepositoryInterface, mailService proto.MailServiceClient, jwtSecret string, addr string, port string) *UserService {
	return &UserService{
		repo:        repo,
		mailService: mailService,
		jwtSecret:   jwtSecret,
		addr:        addr,
		port:        port,
	}
}

// GenerateAccesssAndRefreshTokens generates access and refresh tokens.
func (s *UserService) GenerateAccesssAndRefreshTokens(userID, role string, accessExp, refreshExp time.Duration) (string, string, error) {
	accessClaims := map[string]interface{}{
		"uid":  userID,
		"type": "access",
		"role": role,
		"exp":  time.Now().Add(accessExp).Unix(),
	}

	refreshClaims := map[string]interface{}{
		"uid":  userID,
		"type": "refresh",
		"role": role,
		"exp":  time.Now().Add(refreshExp).Unix(),
	}

	accsessToken, err := utils.GenerateJWT(accessClaims, s.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to generate access token", zap.Error(err))
		return "", "", ErrInternalServer
	}

	refreshToken, err := utils.GenerateJWT(refreshClaims, s.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", ErrInternalServer
	}

	return accsessToken, refreshToken, nil
}

// CreateUser creates a new user.
func (s *UserService) RegisterUser(ctx context.Context, user *models.User) (string, string, error) {
	logger.Logger.Info("Creating new user")
	user.ID = uuid.New()

	// Check if email is valid
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return "", "", ErrInvalidEmail
	}

	// Check if password length is between 8 and 64 characters
	if len(user.Password) < 8 || len(user.Password) > 64 {
		return "", "", ErrInvalidPasswordLength
	}

	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.Logger.Error("Failed to hash password", zap.Error(err))
		return "", "", err
	}

	user.Password = hashedPass

	access, refresh, err := s.GenerateAccesssAndRefreshTokens(user.ID.String(), "user", time.Hour*1, time.Hour*24*7)
	if err != nil {
		logger.Logger.Error("Failed to generate access and refresh tokens", zap.Error(err))
		return "", "", ErrInternalServer
	}

	user.RefreshTokens = append(user.RefreshTokens, utils.HashToken(refresh))

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		logger.Logger.Error("Faield to create new user", zap.Error(err))
		return "", "", err
	}
	logger.Logger.Info("User created successfully")

	token, err := utils.GenerateToken(128)
	if err != nil {
		logger.Logger.Error("Failed to generate token", zap.Error(err))
		return "", "", err
	}

	// Create verification token
	verToken := &models.Token{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	err = s.repo.CreateVerificationToken(ctx, verToken)
	if err != nil {
		logger.Logger.Error("Failed to create verification token", zap.Error(err))
		return "", "", err
	}

	// Send email with verification token
	res, err := s.mailService.SendMail(context.Background(), &proto.MailRequest{
		To:   []string{user.Email},
		Type: proto.NotificationType_EMAIL_CONFIRMATION,
		Data: map[string]string{
			"Link": fmt.Sprintf("http://%s:%s/api/v1/user/verify_email?token=%s", s.addr, s.port, token),
		},
	})

	if err != nil || !res.Success {
		logger.Logger.Error("Failed to send email", zap.Error(err))
		return "", "", ErrInvalidFaildToSendEmail
	}

	return access, refresh, nil
}

func (s *UserService) ResendVerificationEmail(ctx context.Context, jwt string) error {
	token, err := utils.GenerateToken(128)
	if err != nil {
		logger.Logger.Error("Failed to generate token", zap.Error(err))
		return err
	}

	claims, err := utils.ValidateJWT(jwt, s.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to validate JWT", zap.Error(err))
		return ErrInvalidJWT
	}

	if utils.ValidateClaims(claims, "access") {
		logger.Logger.Error("Invalid JWT claims")
		return ErrInvalidJWT
	}

	uid := claims["uid"].(string)

	// Get user by ID and check if email is already verified
	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		logger.Logger.Error("Failed to get user by ID", zap.Error(err))
		return err
	}

	if user.EmailVerified {
		return ErrEmailAlreadyVerified
	}

	// Create verification token
	verToken := &models.Token{
		UserID:    uuid.MustParse(uid),
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	err = s.repo.CreateVerificationToken(ctx, verToken)
	if err != nil {
		logger.Logger.Error("Failed to create verification token", zap.Error(err))
		return err
	}

	// Send email with verification token
	res, err := s.mailService.SendMail(context.Background(), &proto.MailRequest{
		To:   []string{user.Email},
		Type: proto.NotificationType_EMAIL_CONFIRMATION,
		Data: map[string]string{
			"Link": fmt.Sprintf("http://%s:%s/api/v1/user/verify_email?token=%s", s.addr, s.port, token),
		},
	})

	if err != nil || !res.Success {
		logger.Logger.Error("Failed to send email", zap.Error(err))
		return ErrInvalidFaildToSendEmail
	}
	return nil
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {
	err := s.repo.VerifyEmail(ctx, token)
	if err != nil {
		logger.Logger.Error("Failed to verify email", zap.Error(err))
		return err
	}
	return nil
}

func (s *UserService) RequestResetPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Logger.Error("Failed to get user by email", zap.Error(err))
		return err
	}

	token, err := utils.GenerateToken(128)
	if err != nil {
		logger.Logger.Error("Failed to generate token", zap.Error(err))
		return err
	}

	// Create reset password token
	resetToken := &models.Token{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	err = s.repo.CreateResetPasswordToken(ctx, resetToken)
	if err != nil {
		logger.Logger.Error("Failed to create reset password token", zap.Error(err))
		return err
	}

	// Send email with reset password token
	res, err := s.mailService.SendMail(context.Background(), &proto.MailRequest{
		To:   []string{email},
		Type: proto.NotificationType_PASSWORD_RESET,
		Data: map[string]string{
			"Link":       fmt.Sprintf("http://%s:%s/api/v1/user/reset_password?token=%s", s.addr, s.port, token),
			"Exparation": fmt.Sprintf("%v hour/s", time.Now().Add(time.Hour*24)),
		},
	})

	if err != nil || !res.Success {
		logger.Logger.Error("Failed to send email", zap.Error(err))
		return ErrInvalidFaildToSendEmail
	}
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, token, password string) error {
	if len(password) < 8 || len(password) > 64 {
		return ErrInvalidPasswordLength
	}

	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		logger.Logger.Error("Failed to hash password", zap.Error(err))
		return err
	}

	err = s.repo.ResetPassword(ctx, token, hashedPass)
	if err != nil {
		logger.Logger.Error("Failed to reset password", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, user *models.User) (string, string, error) {
	localUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", "", ErrInvalidEmailOrPassword
	}

	if !utils.ComparePasswords(localUser.Password, user.Password) {
		return "", "", ErrInvalidEmailOrPassword
	}

	access, refresh, err := s.GenerateAccesssAndRefreshTokens(localUser.ID.String(), "user", time.Hour*1, time.Hour*24*7)
	if err != nil {
		logger.Logger.Error("Failed to generate access and refresh tokens", zap.Error(err))
		return "", "", err
	}

	err = s.repo.AddRefreshToken(ctx, utils.HashToken(refresh), localUser.ID.String())
	if err != nil {
		logger.Logger.Error("Failed to add refresh token", zap.Error(err))
		return "", "", err
	}

	return access, refresh, nil
}

func (s *UserService) RefreshTokens(ctx context.Context, jwt string) (string, string, error) {
	claims, err := utils.ValidateJWT(jwt, s.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to validate JWT", zap.Error(err))
		return "", "", ErrInvalidJWT
	}

	if utils.ValidateClaims(claims, "refresh") {
		logger.Logger.Error("Invalid JWT claims")

		// Check if refresh token is expired
		if claims.VerifyExpiresAt(time.Now().Unix(), true) {
			logger.Logger.Error("Refresh token is expired")
			s.repo.DeleteRefreshToken(ctx, utils.HashToken(jwt), claims["uid"].(string))
		}

		return "", "", ErrInvalidJWT
	}

	user, err := s.repo.GetUserByID(ctx, claims["uid"].(string))
	if err != nil {
		logger.Logger.Error("Failed to get user by ID", zap.Error(err))
		return "", "", err
	}

	err = s.repo.VerifyRefreshToken(ctx, utils.HashToken(jwt), user.ID.String())
	if err != nil {
		logger.Logger.Error("Failed to verify refresh token", zap.Error(err))
		return "", "", err
	}

	access, refresh, err := s.GenerateAccesssAndRefreshTokens(claims["uid"].(string), "user", time.Hour*1, time.Hour*24*7)
	if err != nil {
		logger.Logger.Error("Failed to generate access and refresh tokens", zap.Error(err))
		return "", "", ErrInternalServer
	}

	err = s.repo.UpdateRefreshToken(ctx, utils.HashToken(jwt), utils.HashToken(refresh), user.ID.String())
	if err != nil {
		logger.Logger.Error("Failed to update refresh token", zap.Error(err))
		return "", "", err
	}

	return access, refresh, nil
}

// UpdateUser updates a user.
func (s *UserService) UpdateUser(ctx context.Context, jwt string, user *models.User) error {
	claims, err := utils.ValidateJWT(jwt, s.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to validate JWT", zap.Error(err))
		return ErrInvalidJWT
	}

	if utils.ValidateClaims(claims, "access") {
		logger.Logger.Error("Invalid JWT claims")
		return ErrInvalidJWT
	}

	uid := claims["uid"].(string)

	userID, err := uuid.Parse(uid)
	if err != nil {
		logger.Logger.Error("Failed to parse user ID", zap.String("id", uid), zap.Error(err))
		return err
	}

	user.ID = userID

	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		if err == repository.ErrUserNotFound {
			logger.Logger.Error("User not found by ID", zap.String("id", uid))
			return err
		}
		logger.Logger.Error("Failed to update user", zap.String("id", user.ID.String()), zap.Error(err))
		return err
	}
	return nil
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(ctx context.Context, jwt string) error {
	logger.Logger.Info("Deleting user by ID")

	claims, err := utils.ValidateJWT(jwt, s.jwtSecret)
	if err != nil {
		logger.Logger.Error("Failed to validate JWT", zap.Error(err))
		return ErrInvalidJWT
	}

	if utils.ValidateClaims(claims, "access") {
		logger.Logger.Error("Invalid JWT claims")
		return ErrInvalidJWT
	}

	id := claims["uid"].(string)

	err = s.repo.DeleteUser(ctx, id)
	if err != nil {
		if err == repository.ErrUserNotFound {
			logger.Logger.Error("User not found by ID", zap.String("id", id))
			return err
		}
		logger.Logger.Error("Failed to delete user by ID", zap.String("id", id), zap.Error(err))
		return err
	}
	logger.Logger.Info("User deleted successfully")
	return nil
}
