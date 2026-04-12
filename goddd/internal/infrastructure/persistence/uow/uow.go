package uow

import (
	"goddd/internal/infrastructure/persistence/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)


type UnitOfWork struct {
    userRepo  *postgres.UserRepository
}

func New(pool *pgxpool.Pool) *UnitOfWork {
    uRepo := postgres.NewUserRepository(pool)
    return &UnitOfWork{
        userRepo:  uRepo,
    }
}

func (u *UnitOfWork) Users() *postgres.UserRepository {
    return u.userRepo
}
