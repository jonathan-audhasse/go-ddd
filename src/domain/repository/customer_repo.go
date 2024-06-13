package repository

import (
	"errors"
	"goddd/src/domain/models"
)

var (
	// ErrNotFound is returned when a customer is not found.
	ErrNotFound = errors.New("the item was not found in the repository")
	// ErrFailedToAdd is returned when the could not be added to the repository.
	ErrFailedToAdd = errors.New("failed to add the customer to the repository")
	// ErrUpdateCustomer is returned when the customer could not be updated in the repository.
	ErrListCustomers = errors.New("failed to list customer from the repository")
)

// CustomerRepository is a interface that defines the rules around what a customer repository
type CustomerRepository interface {
	Get(string) (models.Customer, error)
	Add(models.Customer) error
	List() ([]models.Customer, error)
}
