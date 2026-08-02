package transaction

import (
	"context"
	"errors"
	"goddd/internal/application/transaction"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
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
func NewTransactionManager(pool *pgxpool.Pool) transaction.TransactionManager {
	return &transactionManager{pool: pool}
}

// GetTx retrieve the transaction from the context
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// WithTx Put tx into context
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
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
	ctxWithTx := WithTx(ctx, tx)

	if err := fn(ctxWithTx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Err(err).Msg("failed to commit transaction")
		return ErrFailedToCommit
	}

	return nil
}
