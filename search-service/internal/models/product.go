package models

// Product represents a product in the search service.
type Product struct {
	ID          int64          `json:"id" bson:"id"`
	Name        string         `json:"name" bson:"name"`
	Category    string         `json:"category" bson:"category"`
	Price       float64        `json:"price" bson:"price"`
	Description string         `json:"description" bson:"description"`
	Images      []string       `json:"images" bson:"images"`
	Attributes  map[string]any `json:"attributes" bson:"attributes"`
}
