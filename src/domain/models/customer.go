package models

import (
	"fmt"
	"goddd/src/pkg/errors"

	"github.com/google/uuid"
)

// Customer domain model
type Customer struct {
	// ID
	ID string `json:"id"`
	// name
	Name string `json:"name"`
	// email
	Email string `json:"email"`
}

func NewCustomer(name, email string) (Customer, error) {
	if name == "" {
		return Customer{}, errors.InvalidFormat.New("customer name should not be empty")
	}
	return Customer{ID: fmt.Sprintf("cust_%s", uuid.NewString()[:8]), Name: name, Email: email}, nil
}
