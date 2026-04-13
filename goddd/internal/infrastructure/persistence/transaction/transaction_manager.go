package transaction

import (
	"context"
	"goddd/internal/application"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// txKey key to store the transation in the context
type txKey struct{}

// transactionManager receiver to holds transaction interface
type transactionManager struct {
    pool *pgxpool.Pool
}

// NewTransactionManager define a new instance of TransactionManager
func NewTransactionManager(pool *pgxpool.Pool) application.TransactionManager {
    return &transactionManager{pool: pool}
}

// GetTx retrieve the transaction from the context
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// Do run the function `fn` around an a transation
func (m *transactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
    // Begin the transaction
    tx, err := m.pool.Begin(ctx)
    if err != nil {
        return err
    }

    // defer rollback
    defer tx.Rollback(ctx)

    // Put tx into context (so repos can pick it up)
    ctxWithTx := context.WithValue(ctx, txKey{}, tx)

    if err := fn(ctxWithTx); err != nil {
        return err
    }

    return tx.Commit(ctx)
}