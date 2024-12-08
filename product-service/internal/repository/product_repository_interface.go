package repository

import (
	"context"

	"github.com/NeGat1FF/e-commerce/product-service/internal/models"
)

// ProductRepository defines the methods that any
// data storage provider needs to implement to get products.
type ProductRepository interface {
	// GetProductsByCategory retrieves products by category with pagination.
	GetProductsByCategory(ctx context.Context, category string, page, limit int) ([]models.UserProduct, error)

	// GetProductByID retrieves a product by its ID.
	GetProductByID(ctx context.Context, id int64) (models.UserProduct, error)

	// CreateProduct adds a new product to the catalog.
	CreateProduct(ctx context.Context, product models.Product) error

	// UpdateProduct updates an existing product.
	UpdateProduct(ctx context.Context, id int64, updateFields map[string]any) error

	// DeleteProduct removes a product from the catalog.
	DeleteProduct(ctx context.Context, id int64) error

	// GetStock retrieves the stock quantity of a product.
	GetStock(ctx context.Context, id int64) (int64, error)

	// AddStock increases the stock quantity of a product.
	AddStock(ctx context.Context, id int64, quantity int64) error

	// ReduceStock decreases the stock quantity of a product.
	ReduceStock(ctx context.Context, id int64, quantity int64) error
}
