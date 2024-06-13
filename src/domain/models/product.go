package models

import "github.com/google/uuid"

// Product domain model
type Product struct {
	// ID
	ID uuid.UUID `json:"id"`
	// name
	Name string `json:"name"`
	// Price
	Price int `json:"price"`
}
