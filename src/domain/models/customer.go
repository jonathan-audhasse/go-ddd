package models

import (
	"fmt"
	"goddd/pkg/errors"

	"github.com/google/uuid"
)

// Customer domain model
type Customer struct {
	// ID
	Id string `json:"id" db:"id"`
	// user id
	UserId string `json:"-" db:"user_id"`
	// name
	Name string `json:"name" db:"name"`
	// email
	Email string `json:"email" db:"email"`
}

func NewCustomer(userid, name, email string) (Customer, error) {
	// check input
	if userid == "" {
		return Customer{}, errors.InternalError.New("user must be specified")
	}
	if name == "" {
		return Customer{}, errors.InvalidFormat.New("customer name should not be empty")
	}
	return Customer{
		Id:     fmt.Sprintf("cust_%s", uuid.NewString()[:8]),
		UserId: userid,
		Name:   name,
		Email:  email,
	}, nil
}
