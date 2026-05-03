package transaction

import (
	"context"
	"errors"
	"goddd/internal/application"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedToBeginTransaction = errors.New("transaction: failed to begin transaction")
	ErrFailedToCommit           = errors.New("transaction: failed to commit")
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
	logger := log.Ctx(ctx)
	// Begin the transaction
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		logger.Err(err).Msg("failed to begin transaction")
		return ErrFailedToBeginTransaction
	}

	// defer rollback
	defer tx.Rollback(ctx)

	// Put tx into context (so repos can pick it up)
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctxWithTx); err != nil {
		logger.Err(err).Msg("failed to commit transaction")
		return ErrFailedToCommit
	}

	return tx.Commit(ctx)
}
