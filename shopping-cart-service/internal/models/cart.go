package models

type Item struct {
	ItemID   int64   `json:"id" gorm:"column:item_id"`
	Quantity int     `json:"quantity" gorm:"column:quantity"`
	Price    float64 `json:"price" gorm:"column:price"`
}

type GetCartResponse struct {
	UserID     string  `json:"user_id"`
	Items      []Item  `json:"items"`
	TotalPrice float64 `json:"total_price"`
}

// Cart represents a shopping cart
type Cart struct {
	UserID   string  `gorm:"column:user_id;primaryKey"`
	ItemID   int64   `gorm:"column:item_id;primaryKey"`
	Quantity int     `gorm:"column:quantity;not null;default:1"`
	Price    float64 `gorm:"column:price;not null;default:0"`
}
