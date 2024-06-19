package repository

import (
	"goddd/src/domain/models"
)

// CustomerRepository is a interface that defines the rules around what a customer repository
type CustomerRepository interface {
	Get(string) (models.Customer, error)
	Add(models.Customer) error
	List() ([]models.Customer, error)
}
