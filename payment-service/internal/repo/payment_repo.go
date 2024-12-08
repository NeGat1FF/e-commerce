package repo

import (
	"context"

	"github.com/NeGat1FF/e-commerce/payment-service/internal/models"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

// CreatePayment creates a new payment
func (pr *PaymentRepo) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return pr.db.Create(payment).Error
}

// GetPaymentByID retrieves a payment by its ID
func (pr *PaymentRepo) GetPaymentByID(ctx context.Context, id string) (*models.Payment, error) {
	payment := &models.Payment{}
	err := pr.db.Model(&models.Payment{}).Where("id = ?", id).First(payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// GetPaymentByOrderID retrieves a payment by its order ID
func (pr *PaymentRepo) GetPaymentByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	payment := &models.Payment{}
	err := pr.db.Model(&models.Payment{}).Where("order_id = ?", orderID).First(payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (pr *PaymentRepo) UpdatePaymentStatus(ctx context.Context, id string, status int) (*models.Payment, error) {
	tx := pr.db.Model(&models.Payment{}).Where("id = ?", id).Update("status", status).Error
	if tx != nil {
		return nil, tx
	}

	return pr.GetPaymentByID(ctx, id)
}
