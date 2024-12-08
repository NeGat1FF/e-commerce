package main

import (
	"net/http"

	"github.com/NeGat1FF/e-commerce/payment-service/internal/config"
	"github.com/NeGat1FF/e-commerce/payment-service/internal/db"
	"github.com/NeGat1FF/e-commerce/payment-service/internal/handlers"
	"github.com/NeGat1FF/e-commerce/payment-service/internal/repo"
	"github.com/NeGat1FF/e-commerce/payment-service/internal/service"
)

func main() {

	cfg := config.LoadConfig()

	db, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	repo := repo.NewPaymentRepo(db)
	service := service.NewPaymentService(repo, cfg.StripeSecret)
	handler := handlers.NewPaymentHandler(service, cfg.StripeWebhookSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/order/{id}/payment", handler.CreatePayment)
	mux.HandleFunc("GET /api/v1/order/{id}/payment", handler.GetPaymentByOrderID)
	mux.HandleFunc("GET /api/v1/payments/{id}", handler.GetPaymentByID)
	mux.HandleFunc("POST /api/v1/webhook", handler.UpdatePaymentStatus)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
