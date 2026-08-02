package transaction_test

import (
	"context"
	"goddd/internal/infrastructure/persistence/transaction"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestTransactionManager_Do_Commits(t *testing.T) {
	ctx := context.Background()

	tm := transaction.NewTransactionManager(testDB)

	email := gofakeit.Email()
	username := gofakeit.Username()

	err := tm.Do(ctx, func(ctx context.Context) error {
		// tx must exist in ctx
		tx, ok := transaction.GetTx(ctx)
		require.True(t, ok)
		require.NotNil(t, tx)

		_, err := tx.Exec(ctx,
			`INSERT INTO users (email, username) VALUES ($1, $2)`,
			email, username,
		)
		return err
	})

	require.NoError(t, err)

	// outside transaction, row should exist (commit happened)
	var count int
	err = testDB.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE email = $1`,
		email,
	).Scan(&count)

	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestTransactionManager_Do_Rollbacks(t *testing.T) {
	ctx := context.Background()

	tm := transaction.NewTransactionManager(testDB)

	email := gofakeit.Email()
	username := gofakeit.Username()
	mockErr := gofakeit.Error()

	err := tm.Do(ctx, func(ctx context.Context) error {
		tx, ok := transaction.GetTx(ctx)
		require.True(t, ok)

		_, err := tx.Exec(ctx,
			`INSERT INTO users (email, username) VALUES ($1, $2)`,
			email, username,
		)
		require.NoError(t, err)

		// force failure
		return mockErr
	})

	// expect ErrFailedToCommit when failed
	require.ErrorIs(t, err, mockErr)

	// row should NOT exist because rollback happened
	var count int
	err = testDB.QueryRow(ctx,
		`SELECT COUNT(*) FROM users WHERE email = $1`,
		email,
	).Scan(&count)

	require.NoError(t, err)
	require.Equal(t, 0, count)
}
