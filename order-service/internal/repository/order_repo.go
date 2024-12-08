package repository

import (
	"context"
	"fmt"

	"github.com/NeGat1FF/e-commerce/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// CreateOrder creates a new order
func (r *OrderRepo) CreateOrder(ctx context.Context, order *models.Order, items []*models.OrderItem) (*models.OrderResponse, error) {
	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil

	})

	return r.GetOrderByID(ctx, order.ID.String())
}

// GetOrderByID gets an order by its ID
func (r *OrderRepo) GetOrderByID(ctx context.Context, orderID string) (*models.OrderResponse, error) {
	order := &models.OrderResponse{}

	tx := r.db.Raw(`
WITH product_details AS (
    SELECT
        oi.order_id,
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'product_id', oi.product_id,
                'quantity', oi.quantity
            )
        ) AS products
    FROM
        order_items oi
    GROUP BY
        oi.order_id
)
SELECT
    o.id AS id,
    o.user_id AS user_id,
    pd.products AS products,
    o.total AS total,
	st.name as status,
    o.created_at AS created_at,
    o.updated_at AS updated_at
FROM
    orders o
JOIN
    product_details pd
ON
    o.id = pd.order_id
JOIN
    order_statuses st
ON
    o.status = st.id
WHERE o.id = ?;
	`, orderID).Scan(order)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return order, nil
}

// GetOrders gets all orders of a user
func (r *OrderRepo) GetOrders(ctx context.Context, userID string) ([]*models.OrderResponse, error) {
	orders := []*models.OrderResponse{}

	tx := r.db.Raw(`
WITH product_details AS (
    SELECT
        oi.order_id,
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'product_id', oi.product_id,
                'quantity', oi.quantity
            )
        ) AS products
    FROM
        order_items oi
    GROUP BY
        oi.order_id
)
SELECT
    o.id AS id,
    o.user_id AS user_id,
    pd.products AS products,
    o.total AS total,
	st.name as status,
    o.created_at AS created_at,
    o.updated_at AS updated_at
FROM
    orders o
JOIN
    product_details pd
ON
    o.id = pd.order_id
JOIN
    order_statuses st
ON
    o.status = st.id
WHERE o.user_id = ?;
	`, userID).Scan(&orders)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return orders, nil
}

// UpdateOrder updates an order
func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, orderId string, status int) (*models.OrderResponse, error) {

	fmt.Println(orderId, status)

	tx := r.db.Exec(`
	UPDATE orders
	SET status = ?
	WHERE id = ?;`, status, orderId)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return r.GetOrderByID(ctx, orderId)
}
