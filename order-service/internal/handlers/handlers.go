package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NeGat1FF/e-commerce/order-service/internal/models"
	"github.com/NeGat1FF/e-commerce/order-service/internal/service"
)

// OrderHandler is a struct that holds the order service
type OrderHandler struct {
	OrderService *service.OrderService
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

// CreateOrder creates a new order
func (oh *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq models.CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	jwt := strings.Split(auth, " ")[1]

	resp, err := oh.OrderService.CreateOrder(r.Context(), jwt, orderReq.Products)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetOrder gets an order
func (oh *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	resp, err := oh.OrderService.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (oh *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Authorization header is required", http.StatusBadRequest)
		return
	}

	jwt := strings.Split(auth, " ")[1]

	resp, err := oh.OrderService.GetOrders(r.Context(), jwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateOrder updates an order
func (oh *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var orderReq struct {
		Status int `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := oh.OrderService.UpdateOrderStatus(r.Context(), id, orderReq.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
