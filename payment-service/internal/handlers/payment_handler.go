package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/NeGat1FF/e-commerce/payment-service/internal/models"
	"github.com/NeGat1FF/e-commerce/payment-service/internal/service"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/webhook"
)

// PaymentHandler is a handler for the payment service
type PaymentHandler struct {
	service             *service.PaymentService
	stripeWebhookSecret string
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(service *service.PaymentService, stripeWebhookSecret string) *PaymentHandler {
	return &PaymentHandler{
		service:             service,
		stripeWebhookSecret: stripeWebhookSecret,
	}
}

// CreatePayment creates a new payment
func (ph *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	var req models.CreatePaymentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stripe.Key = ph.stripeWebhookSecret
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment := &models.Payment{
		OrderID: id,
		Amount:  req.Amount,
	}

	resp, err := ph.service.CreatePayment(r.Context(), payment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// GetPaymentByID retrieves a payment by its ID
func (ph *PaymentHandler) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	payment, err := ph.service.GetPaymentByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

// GetPaymentByOrderID retrieves a payment by its order ID
func (ph *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	payment, err := ph.service.GetPaymentByOrderID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

// UpdatePaymentStatus updates the status of a payment
func (ph *PaymentHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	// Handle a webhook from Stripe

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		ph.stripeWebhookSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		ph.service.UpdatePaymentStatus(r.Context(), event.ID, models.PaymentStatusSucceeded)
	case "payment_intent.payment_failed":
		ph.service.UpdatePaymentStatus(r.Context(), event.ID, models.PaymentStatusFailed)
	default:
		return
	}
}
