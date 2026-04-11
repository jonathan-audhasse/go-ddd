package services

import (
	"goddd/domain/repository"
)

type Services struct {
	*HealthService
	*UserService
	*CustomerService
}

// Factory
func NewServices(repo *repository.Repository) *Services {
	return &Services{
		HealthService:   NewHealthService(repo.HealthRepo),
		UserService:     NewUserService(repo.UserRepo),
		CustomerService: NewCustomerService(repo.CustomerRepo),
	}
}
