package repository

import (
	"context"

	"github.com/NeGat1FF/e-commerce/order-service/internal/models"
)

// OrderRepoInterface is an interface for order repository
type OrderRepoInterface interface {
	CreateOrder(ctx context.Context, order *models.Order, items []*models.OrderItem) (*models.OrderResponse, error)
	GetOrderByID(ctx context.Context, orderID string) (*models.OrderResponse, error)
	GetOrders(ctx context.Context, userID string) ([]*models.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, orderId string, status int) (*models.OrderResponse, error)
}
