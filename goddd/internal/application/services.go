package application

import (
	"goddd/internal/application/health"
	"goddd/internal/application/transaction"
	"goddd/internal/application/user"
	"goddd/internal/domain/repository"
)

type Services struct {
	Health health.HealthService
	User   user.UserService
}

// NewServices instantiates new services
func NewServices(
	repo *repository.Repositories,
	tm transaction.TransactionManager,
	pinger health.Pinger,
) *Services {
	return &Services{
		Health: health.NewHealthService(pinger),
		User:   user.NewUserService(repo.UserRepo, tm),
	}
}
