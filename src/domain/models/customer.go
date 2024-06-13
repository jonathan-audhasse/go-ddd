package models

import (
	"errors"
	"fmt"

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

var (
	// ErrInvalidCustomer is returned when customer's name is empty
	ErrInvalidCustomer = errors.New("a customer has to have an valid person")
)

func NewCustomer(name, email string) (Customer, error) {
	if name == "" {
		return Customer{}, ErrInvalidCustomer
	}
	return Customer{ID: fmt.Sprintf("cust_%s", uuid.NewString()), Name: name, Email: email}, nil
}
