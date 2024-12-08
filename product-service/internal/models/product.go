package models

// Product represents the internal model of a product that includes quantity
type Product struct {
	ID          int64          `json:"id" bson:"id"`
	Name        string         `json:"name" bson:"name"`
	Category    string         `json:"category" bson:"category"`
	Price       float64        `json:"price" bson:"price"`
	Description string         `json:"description" bson:"description"`
	Quantity    int64          `json:"quantity" bson:"quantity"`
	Images      []string       `json:"images" bson:"images"`
	Attributes  map[string]any `json:"attributes" bson:"attributes"`
}

// UserProduct represents the model of a product that is exposed to the user
type UserProduct struct {
	ID          int64          `json:"id" bson:"id"`
	Name        string         `json:"name" bson:"name"`
	Category    string         `json:"category" bson:"category"`
	Price       float64        `json:"price" bson:"price"`
	Description string         `json:"description" bson:"description"`
	Images      []string       `json:"images" bson:"images"`
	Attributes  map[string]any `json:"attributes" bson:"attributes"`
}
