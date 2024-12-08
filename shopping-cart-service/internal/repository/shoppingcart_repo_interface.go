package repository

import (
	"context"

	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/models"
)

type ShoppingCartRepoInterface interface {
	// AddItem adds an item to the shopping cart
	AddItem(ctx context.Context, item *models.Cart) (models.GetCartResponse, error)

	// SetQuantity sets the quantity of an item in the shopping cart
	SetQuantity(ctx context.Context, userID string, itemID int64, quantity int) (models.GetCartResponse, error)

	// RemoveItem removes an item from the shopping cart
	RemoveItem(ctx context.Context, userID string, itemID int64) (models.GetCartResponse, error)

	// GetItems returns all items in the shopping cart
	GetItems(ctx context.Context, userID string) (models.GetCartResponse, error)

	// DeleteCart deletes the shopping cart
	DeleteCart(ctx context.Context, userID string) error
}
