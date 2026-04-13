package application

import (
	"context"
)

// TransactionManager is a transaction runner that handles commit/rollback
// This is for Unit of Work: it defines the transactional boundary.
type TransactionManager interface {
	// Do run the function fn wrap in a trasaction.
    Do(ctx context.Context, fn func(ctx context.Context) error) error
}