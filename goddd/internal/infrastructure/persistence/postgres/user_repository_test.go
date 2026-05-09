package postgres_test

import (
	"goddd/internal/infrastructure/persistence/postgres"
	"goddd/internal/infrastructure/persistence/postgres/sqlcgen"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_FindByID(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, tx, rollback := withTx(t, ctx, testDB)

	defer rollback()

	q := sqlcgen.New(tx)

	// insert user...
	id := uuid.New()
	// now := time.Now().UTC()
	now := pgtype.Timestamptz{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	_, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		ID:           postgres.ToPgUUID(id), // or postgres.toPgUUID if not exported
		Email:        "john@example.com",
		Username:     "john",
		PasswordHash: "hash",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)

	user, err := repo.FindByID(ctx, id)
	require.NoError(t, err)

	// assert
	require.Equal(t, id, user.ID)
	require.Equal(t, "john@example.com", user.Email)
	require.Equal(t, "john", user.Username)
	require.Equal(t, "hash", user.PasswordHash)
}
