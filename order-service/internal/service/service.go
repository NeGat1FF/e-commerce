package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NeGat1FF/e-commerce/order-service/internal/models"
	"github.com/NeGat1FF/e-commerce/order-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/order-service/internal/utils"
	"github.com/NeGat1FF/e-commerce/order-service/proto"
	"github.com/google/uuid"
)

type PriceServiceInterface interface {
	GetPrice(ctx context.Context, productID int64) (float64, error)
}

type OrderService struct {
	repo         repository.OrderRepoInterface
	priceService proto.PriceServiceClient
	secret       string
}

func NewOrderService(repo repository.OrderRepoInterface, priceService proto.PriceServiceClient, secret string) *OrderService {
	return &OrderService{
		repo:         repo,
		priceService: priceService,
		secret:       secret,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, jwt string, items models.Products) (*models.OrderResponse, error) {
	claims, err := utils.ValidateToken(jwt, s.secret)
	if err != nil {
		return nil, err
	}

	var OrderItems []*models.OrderItem

	var total float64

	userUid, err := uuid.Parse(claims["uid"].(string))
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		ID:     uuid.New(),
		UserID: userUid,
		Status: models.OrderStatusPending,
	}

	for _, item := range items {
		priceRes, err := s.priceService.GetPrice(ctx, &proto.PriceRequest{ProductId: fmt.Sprintf("%d", item.ID)})
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(priceRes.Price, 64)
		if err != nil {
			return nil, err
		}
		OrderItems = append(OrderItems, &models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ID,
			Quantity:  item.Quantity,
			Price:     price,
		})

		total += price * float64(item.Quantity)
	}

	order.Total = total

	res, err := s.repo.CreateOrder(ctx, order, OrderItems)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetOrder returns an order by its ID
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	res, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetOrders returns all orders
func (s *OrderService) GetOrders(ctx context.Context, jwt string) ([]*models.OrderResponse, error) {
	claims, err := utils.ValidateToken(jwt, s.secret)
	if err != nil {
		return nil, err
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		return nil, fmt.Errorf("uid not found in claims")
	}

	res, err := s.repo.GetOrders(ctx, uid)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, status int) (*models.OrderResponse, error) {
	return s.repo.UpdateOrderStatus(ctx, orderID, status)
}
