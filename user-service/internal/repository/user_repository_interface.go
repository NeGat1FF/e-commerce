package repository

import (
	"context"

	"github.com/NeGat1FF/e-commerce/user-service/internal/models"
)

type UserRepositoryInterface interface {
	// Create a new user
	CreateUser(ctx context.Context, user *models.User) error

	// Get a user by ID
	GetUserByID(ctx context.Context, id string) (*models.User, error)

	// Get a user by email
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// Update a user
	UpdateUser(ctx context.Context, user *models.User) error

	// Delete a user
	DeleteUser(ctx context.Context, id string) error

	// Create a verification token
	CreateVerificationToken(ctx context.Context, token *models.Token) error

	// Verify user email
	VerifyEmail(ctx context.Context, token string) error

	// Create reset password token
	CreateResetPasswordToken(ctx context.Context, token *models.Token) error

	// Add refresh token
	AddRefreshToken(ctx context.Context, tokenHash, userId string) error

	// Verify refresh token
	VerifyRefreshToken(ctx context.Context, tokenHash, userId string) error

	// Update refresh token
	UpdateRefreshToken(ctx context.Context, oldTokenHash, newTokenHash, userId string) error

	// Delete refresh token
	DeleteRefreshToken(ctx context.Context, tokenHash, userId string) error

	// Reset user password
	ResetPassword(ctx context.Context, token, password string) error
}
