package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/models"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/shopping-cart-service/proto"
)

type CartService struct {
	repo         repository.ShoppingCartRepoInterface
	priceService proto.PriceServiceClient
}

func NewCartService(repo repository.ShoppingCartRepoInterface, priceService proto.PriceServiceClient) *CartService {
	return &CartService{
		repo:         repo,
		priceService: priceService,
	}
}

// AddItem adds an item to the shopping cart
func (s *CartService) AddItem(ctx context.Context, userID string, item *models.Item) (models.GetCartResponse, error) {
	priceRes, err := s.priceService.GetPrice(ctx, &proto.PriceRequest{ProductId: fmt.Sprintf("%d", item.ItemID)})
	if err != nil {
		return models.GetCartResponse{}, err
	}

	price, err := strconv.ParseFloat(priceRes.Price, 64)
	if err != nil {
		return models.GetCartResponse{}, err
	}

	var cart models.Cart
	cart.UserID = userID
	cart.ItemID = item.ItemID
	cart.Quantity = item.Quantity
	cart.Price = price

	response, err := s.repo.AddItem(ctx, &cart)
	if err != nil {
		return models.GetCartResponse{}, err
	}

	return response, nil
}

// SetQuantity sets the quantity of an item in the shopping cart
func (s *CartService) SetQuantity(ctx context.Context, userID string, itemID int64, quantity int) (models.GetCartResponse, error) {
	response, err := s.repo.SetQuantity(ctx, userID, itemID, quantity)
	if err != nil {
		return models.GetCartResponse{}, err
	}

	return response, nil
}

// RemoveItem removes an item from the shopping cart
func (s *CartService) RemoveItem(ctx context.Context, userID string, itemID int64) (models.GetCartResponse, error) {
	response, err := s.repo.RemoveItem(ctx, userID, itemID)
	if err != nil {
		return models.GetCartResponse{}, err
	}

	return response, nil
}

// GetItems returns all items in the shopping cart
func (s *CartService) GetItems(ctx context.Context, userID string) (models.GetCartResponse, error) {
	return s.repo.GetItems(ctx, userID)
}

// DeleteCart deletes the shopping cart
func (s *CartService) DeleteCart(ctx context.Context, userID string) error {
	return s.repo.DeleteCart(ctx, userID)
}
