package repository

import (
	"context"

	"github.com/NeGat1FF/e-commerce/shopping-cart-service/internal/models"
	"gorm.io/gorm"
)

type ShoppingCartRepo struct {
	db *gorm.DB
}

func NewShoppingCartRepo(db *gorm.DB) *ShoppingCartRepo {
	return &ShoppingCartRepo{
		db: db,
	}
}

// AddItem adds an item to the shopping cart
func (repo *ShoppingCartRepo) AddItem(ctx context.Context, item *models.Cart) (models.GetCartResponse, error) {
	tx := repo.db.WithContext(ctx).Table("cart").Save(item)
	if tx.Error != nil {
		return models.GetCartResponse{}, tx.Error
	}

	return repo.GetItems(ctx, item.UserID)
}

// SetQuantity sets quantity of an item in the shopping cart
func (repo *ShoppingCartRepo) SetQuantity(ctx context.Context, userID string, itemID int64, quantity int) (models.GetCartResponse, error) {
	tx := repo.db.WithContext(ctx).Exec(`
	UPDATE cart
	SET quantity = ?
	WHERE user_id = ? AND item_id = ?
	`, quantity, userID, itemID)
	if tx.Error != nil {
		return models.GetCartResponse{}, tx.Error
	}

	return repo.GetItems(ctx, userID)
}

// RemoveItem removes an item from the shopping cart
func (repo *ShoppingCartRepo) RemoveItem(ctx context.Context, userID string, itemID int64) (models.GetCartResponse, error) {
	tx := repo.db.WithContext(ctx).Exec(`
	DELETE FROM cart
	WHERE user_id = ? AND item_id = ?
	`, userID, itemID)
	if tx.Error != nil {
		return models.GetCartResponse{}, tx.Error
	}

	return repo.GetItems(ctx, userID)
}

// GetItems returns all items in the shopping cart
func (repo *ShoppingCartRepo) GetItems(ctx context.Context, userID string) (models.GetCartResponse, error) {
	tx := repo.db.WithContext(ctx).Raw(`
	SELECT 
		user_id,
		item_id,
		quantity,
		price,
    	SUM(quantity * price) OVER (PARTITION BY user_id) AS total_price
	FROM 
		cart
	WHERE 
		user_id = ?
	`, userID)

	if tx.Error != nil {
		return models.GetCartResponse{}, tx.Error
	}

	rows, err := tx.Rows()
	if err != nil {
		return models.GetCartResponse{}, err
	}
	defer rows.Close()

	var response models.GetCartResponse
	var items []models.Item
	var totalPrice float64

	for rows.Next() {
		var item models.Item
		var currentUserID string

		if err := rows.Scan(&currentUserID, &item.ItemID, &item.Quantity, &item.Price, &totalPrice); err != nil {
			return response, err
		}

		// Add item to the list
		items = append(items, item)
		response.UserID = currentUserID
	}

	response.Items = items
	response.TotalPrice = totalPrice

	return response, nil

}

func (repo *ShoppingCartRepo) DeleteCart(ctx context.Context, userID string) error {
	tx := repo.db.WithContext(ctx).Exec(`
	DELETE FROM cart
	WHERE user_id = ?
	`, userID)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
