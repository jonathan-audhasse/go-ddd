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

func NewServices(repo *repository.Repository, tm transaction.TransactionManager) *Services {
	return &Services{
		Health: health.NewHealthService(repo.HealthRepo),
		User:   user.NewUserService(repo.UserRepo, tm),
	}
}
