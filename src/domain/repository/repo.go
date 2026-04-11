package repository

import "goddd/domain/models"

type Repository struct {
	HealthRepo   HealthRepository
	UserRepo     UserRepository
	CustomerRepo CustomerRepository
}

// Healt repository interface
type HealthRepository interface {
	// retrieve the health check value
	HealthValue() (string, error)
}

// User repository interface
type UserRepository interface {
	// find user by id
	GetById(id string) (models.User, error)
	// find user by username
	GetByUsername(username string) (models.User, error)
}

// CustomerRepository is a interface that defines the rules around what a customer repository
type CustomerRepository interface {
	Get(string) (models.Customer, error)
	// add customer to repository
	Add(models.Customer) error
	// list customer for a given user
	ListByUser(userid string) ([]models.Customer, error)
	Delete(string) error
}
