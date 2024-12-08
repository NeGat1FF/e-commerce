package models

type OrderStatus struct {
	ID   int    `json:"id" gorm:"type:int;primaryKey"`
	Name string `json:"name"`
}

const (
	OrderStatusPending  = 1
	OrderStatusAccepted = 2
	OrderStatusRejected = 3
)
