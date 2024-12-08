package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"

	"github.com/NeGat1FF/e-commerce/payment-service/internal/models"
	"github.com/NeGat1FF/e-commerce/payment-service/internal/repo"
)

type PaymentService struct {
	repo         repo.PaymentRepoInterface
	stripeSecret string
}

func NewPaymentService(repo repo.PaymentRepoInterface, stripeSecret string) *PaymentService {
	return &PaymentService{
		repo:         repo,
		stripeSecret: stripeSecret,
	}
}

// CreatePayment creates a new payment
func (ps *PaymentService) CreatePayment(ctx context.Context, payment *models.Payment) (*models.CreatePaymentResponse, error) {
	intent, err := ps.createPaymentIntent(int64(payment.Amount))
	if err != nil {
		return nil, err
	}

	payment.ID = uuid.New()
	payment.StripeID = intent.ID
	payment.Status = models.PaymentStatusPending

	err = ps.repo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	resp := &models.CreatePaymentResponse{
		ID:            payment.ID.String(),
		PaymentIntent: intent,
	}
	return resp, nil
}

// GetPaymentByID retrieves a payment by its ID
func (ps *PaymentService) GetPaymentByID(ctx context.Context, id string) (*models.Payment, error) {
	return ps.repo.GetPaymentByID(ctx, id)
}

// GetPaymentByOrderID retrieves a payment by its order ID
func (ps *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	return ps.repo.GetPaymentByOrderID(ctx, orderID)
}

// UpdatePaymentStatus updates the status of a payment
func (ps *PaymentService) UpdatePaymentStatus(ctx context.Context, id string, status int) (*models.Payment, error) {
	return ps.repo.UpdatePaymentStatus(ctx, id, status)
}

func (ps *PaymentService) createPaymentIntent(amount int64) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount), // in smallest unit (cents for USD)
		Currency: stripe.String("usd"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card", // Payment method
		}),
	}
	stripe.Key = ps.stripeSecret
	return paymentintent.New(params)
}
