package application

import (
	healthservice "goddd/internal/application/health"
	"goddd/internal/application/transaction"
	userservice "goddd/internal/application/user"
	"goddd/internal/domain/repository"
)

// Services holds the service object
type Services struct {
	*healthservice.HealthService
	*userservice.UserService
}

// NewServices instantiate services
func NewServices(repo *repository.Repository, tm transaction.TransactionManager) *Services {
	return &Services{
		HealthService: healthservice.NewHealthService(repo.HealthRepo),
		UserService:   userservice.NewUserService(repo.UserRepo, tm),
	}
}
