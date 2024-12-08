package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	UserID    uuid.UUID `json:"id" gorm:"column:user_id;type:uuid;primaryKey"`
	Token     string    `json:"token" gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time `json:"expires_at" gor:"column:expires_at;type:timestamp;not null"`
}
