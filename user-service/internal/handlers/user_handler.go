package handlers

import (
	"net/http"
	"strings"

	"github.com/NeGat1FF/e-commerce/user-service/internal/models"
	"github.com/NeGat1FF/e-commerce/user-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Register registers a new user.
func (h *UserHandler) Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ac, rf, err := h.service.RegisterUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "access_token": ac, "refresh_token": rf})
}

// Login logs in a user.
func (h *UserHandler) Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

// VerifyEmail verifies the email.
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Request.URL.Query().Get("token")

	if err := h.service.VerifyEmail(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// ResendVerificationEmail resends the verification email.
func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	jwt := c.GetHeader("Authorization")
	if jwt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	jwt = strings.Split(jwt, " ")[1]

	if err := h.service.ResendVerificationEmail(c.Request.Context(), jwt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

// ForgotPassword sends a reset password email.
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	email := c.Request.URL.Query().Get("email")

	if err := h.service.RequestResetPassword(c.Request.Context(), email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset password email sent successfully"})
}

// ResetPassword resets the password.
func (h *UserHandler) ResetPassword(c *gin.Context) {
	token := c.Request.URL.Query().Get("token")

	var body map[string]string

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), token, body["password"]); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// RefreshToken refreshes the access and refresh tokens.
func (h *UserHandler) RefreshToken(c *gin.Context) {

	var token struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.RefreshTokens(c.Request.Context(), token.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

// UpdateUser updates a user.
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt := c.GetHeader("Authorization")
	if jwt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	jwt = strings.Split(jwt, " ")[1]

	if err := h.service.UpdateUser(c.Request.Context(), jwt, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user.
func (h *UserHandler) DeleteUser(c *gin.Context) {
	jwt := c.GetHeader("Authorization")

	if jwt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	jwt = strings.Split(jwt, " ")[1]

	if err := h.service.DeleteUser(c.Request.Context(), jwt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
