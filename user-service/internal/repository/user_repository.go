package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/NeGat1FF/e-commerce/user-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserWithEmailExists  = errors.New("user with email already exists")
	ErrUserWithPhoneExists  = errors.New("user with phone already exists")
	ErrResetTokenInvalid    = errors.New("invalid or expired reset token")
	ErrNewPasswordSameAsOld = errors.New("new password cannot be the same as the old password")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create a new user
func (u *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	tx := u.db.WithContext(ctx).Create(user)
	if tx.Error != nil {
		errMsg := tx.Error.Error()
		if strings.Contains(errMsg, "violates unique constraint") {
			if strings.Contains(errMsg, "email") {
				return ErrUserWithEmailExists
			}
			if strings.Contains(errMsg, "phone") {
				return ErrUserWithPhoneExists
			}
		}
		return tx.Error
	}

	return nil
}

func (u *UserRepository) CreateVerificationToken(ctx context.Context, token *models.Token) error {
	tx := u.db.WithContext(ctx).Table("email_verifications").Save(token)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u *UserRepository) VerifyEmail(ctx context.Context, token string) error {
	u.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx).Exec(`
	UPDATE users
	SET email_verified = TRUE
	FROM email_verifications
	WHERE users.id = email_verifications.user_id
		AND email_verifications.token = ?
		AND email_verifications.expires_at > NOW();
	`, token)
		if tx.Error != nil {
			return tx.Error
		}

		if tx.RowsAffected == 0 {
			return errors.New("invalid or expired token")
		}

		tx = tx.WithContext(ctx).Exec(`
	DELETE FROM email_verifications
	WHERE token = ?;
	`, token)
		if tx.Error != nil {
			return tx.Error
		}

		return nil

	},
		nil,
	)

	return nil
}

func (u *UserRepository) AddRefreshToken(ctx context.Context, tokenHash, userId string) error {
	tx := u.db.WithContext(ctx).Exec(`
	UPDATE users
	SET refresh_tokens = array_append(refresh_tokens, ?)
	WHERE id = ?;
	`, tokenHash, userId)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u *UserRepository) VerifyRefreshToken(ctx context.Context, tokenHash, userId string) error {
	tx := u.db.WithContext(ctx).Raw(`
	SELECT COUNT(*) 
	FROM users 
	WHERE ? = any(refresh_tokens) 
		AND id = ?;`, tokenHash, userId)
	if tx.Error != nil {
		return tx.Error
	}

	var count int64

	tx.Scan(&count)

	if count != 1 {
		return errors.New("invalid refresh token")
	}

	return nil
}

func (u *UserRepository) UpdateRefreshToken(ctx context.Context, oldTokenHash, newTokenHash, userId string) error {
	u.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx).Exec(`
	UPDATE users
	SET refresh_tokens = array_remove(refresh_tokens, ?)
	WHERE ? = any(refresh_tokens)
		AND id = ?;
	`, oldTokenHash, oldTokenHash, userId)
		if tx.Error != nil {
			return tx.Error
		}

		if tx.RowsAffected == 0 {
			return errors.New("invalid refresh token")
		}

		tx = tx.WithContext(ctx).Exec(`
	UPDATE users
	SET refresh_tokens = array_append(refresh_tokens, ?)
	WHERE id = ?;
	`, newTokenHash, userId)
		if tx.Error != nil {
			return tx.Error
		}

		return nil
	},
		nil,
	)

	return nil
}

func (u *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash, userId string) error {
	tx := u.db.WithContext(ctx).Exec(`
	UPDATE users
	SET refresh_tokens = array_remove(refresh_tokens, ?)
	WHERE id = ?;
	`, tokenHash, userId)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("invalid refresh token")
	}

	return nil
}

func (u *UserRepository) CreateResetPasswordToken(ctx context.Context, token *models.Token) error {
	tx := u.db.WithContext(ctx).Table("password_resets").Save(token)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u *UserRepository) ResetPassword(ctx context.Context, token, password string) error {
	var result string

	err := u.db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		tx = tx.WithContext(ctx).Raw(`
		SELECT reset_password(?, ?)
		`, token, password)
		if tx.Error != nil {
			// return any error will rollback
			return tx.Error
		}

		tx.Scan(&result)

		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {
		return err
	}

	switch result {
	case "Invalid or expired reset token.":
		return ErrResetTokenInvalid
	case "New password cannot be the same as the old password.":
		return ErrNewPasswordSameAsOld
	case "Password updated successfully.":
		return nil
	default:
		return errors.New(result)
	}
}

// Get a user by ID
func (u *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	tx := u.db.WithContext(ctx).First(user, "id = ?", id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, tx.Error
	}
	return user, nil
}

// Get a user by email
func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	tx := u.db.WithContext(ctx).First(user, "email = ?", email)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, tx.Error
	}
	return user, nil
}

// Update a user
func (u *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	tx := u.db.WithContext(ctx).Updates(user)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete a user
func (u *UserRepository) DeleteUser(ctx context.Context, id string) error {
	tx := u.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
