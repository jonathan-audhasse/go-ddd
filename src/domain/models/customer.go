package models

import (
	"fmt"
	"goddd/src/pkg/errors"

	"github.com/google/uuid"
)

// Customer domain model
type Customer struct {
	// ID
	Id string `json:"id" db:"id"`
	// name
	Name string `json:"name" db:"name"`
	// email
	Email string `json:"email" db:"email"`
}

func NewCustomer(name, email string) (Customer, error) {
	if name == "" {
		return Customer{}, errors.InvalidFormat.New("customer name should not be empty")
	}
	return Customer{Id: fmt.Sprintf("cust_%s", uuid.NewString()[:8]), Name: name, Email: email}, nil
}
