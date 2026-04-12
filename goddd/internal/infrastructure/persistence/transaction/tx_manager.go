package transaction

import (
	"context"
	"goddd/internal/infrastructure/persistence/transaction"
	"goddd/internal/infrastructure/persistence/uow"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
    pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool) *Manager {
    return &Manager{pool: pool}
}

func (m *Manager) WithTx(ctx context.Context, fn func(uow *transaction.UnitOfWork) error) error {
    tx, err := m.pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    unit := uow.New(tx)

    if err := fn(unit); err != nil {
        return err
    }

    return tx.Commit(ctx)
}