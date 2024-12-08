package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Product struct {
	ID       int64 `json:"product_id" gorm:"type:uuid;primaryKey"`
	Quantity int   `json:"quantity" gorm:"type:int"`
}

type Products []Product

// Value implements the driver.Valuer interface for saving into the database.
func (p Products) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface for reading from the database.
func (p *Products) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, p)
}
