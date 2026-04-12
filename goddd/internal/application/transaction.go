package application

import (
	"context"
	"goddd/internal/domain/repository"
)

type UnitOfWork interface {
    Users() repository.UserRepository
}

type TransactionManager interface {
    WithTx(ctx context.Context, fn func(uow UnitOfWork) error) error
}