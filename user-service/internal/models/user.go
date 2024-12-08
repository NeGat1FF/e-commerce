package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Name          string         `json:"first_name" gorm:"type:varchar(255);not null;default:null"`
	Surname       string         `json:"last_name" gorm:"type:varchar(255);not null;default:null"`
	Email         string         `json:"email" gorm:"type:varchar(255);unique;not null;default:null"`
	EmailVerified bool           `json:"email_verified" gorm:"column:email_verified;type:boolean;default:false"`
	Phone         string         `json:"phone" gorm:"type:varchar(15);unique;default:null"`
	Password      string         `json:"password" gorm:"type:varchar(255);not null;default:null"`
	RefreshTokens pq.StringArray `json:"refresh_tokens" gorm:"column:refresh_tokens;type:varchar(32)[];default:null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
