package repo

import (
	"context"

	"github.com/NeGat1FF/e-commerce/payment-service/internal/models"
)

// PaymentRepoInterface is an interface for the payment repository
type PaymentRepoInterface interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, id string) (*models.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id string, status int) (*models.Payment, error)
}
