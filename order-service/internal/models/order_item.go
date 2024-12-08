package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;references:orders(id);primaryKey"`
	ProductID int64     `json:"product_id" gorm:"primaryKey"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
