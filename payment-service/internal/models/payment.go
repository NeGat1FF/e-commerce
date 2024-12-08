package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

type Payment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	StripeID  string    `json:"stripe_id"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid"`
	Status    int       `json:"status"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePaymentRequest struct {
	Amount float64 `json:"amount"`
}

type CreatePaymentResponse struct {
	ID            string                `json:"id"`
	PaymentIntent *stripe.PaymentIntent `json:"payment_intent"`
}

const (
	_ = iota
	PaymentStatusPending
	PaymentStatusSucceeded
	PaymentStatusFailed
)
