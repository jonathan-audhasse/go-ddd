package postgres_test

import (
	"goddd/internal/domain/models"
	"goddd/internal/infrastructure/persistence/postgres"
	"goddd/internal/infrastructure/persistence/postgres/sqlcgen"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())

	repo := postgres.NewUserRepository(testDB)

	_, err := repo.FindByID(ctx, uuid.MustParse(gofakeit.UUID()))
	require.ErrorIs(t, err, postgres.ErrFailedToFindUserById)
}

func TestUserRepository_FindByID(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, tx, rollback := withTx(t, ctx, testDB)

	defer rollback()

	q := sqlcgen.New(tx)

	// insert user...
	user, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)
	id := postgres.FromPgUUID(user.ID)
	res, err := repo.FindByID(ctx, id)
	require.NoError(t, err)

	// assert
	require.Equal(t, id, res.ID)
	require.Equal(t, user.Email, res.Email)
	require.Equal(t, user.Username, res.Username)
}

func TestUserRepository_Create_(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, _, rollback := withTx(t, ctx, testDB)

	defer rollback()

	repo := postgres.NewUserRepository(testDB)
	user := models.NewUser{Email: gofakeit.Email(), Username: gofakeit.Username()}

	now := time.Now()

	// insert user...
	res, err := repo.Create(ctx, user)
	require.NoError(t, err)

	// assert user has been added
	u, err := repo.FindByID(ctx, res.ID)
	require.NoError(t, err)
	require.NotEmpty(t, u)

	// assert
	require.Equal(t, user.Email, res.Email)
	require.Equal(t, user.Username, res.Username)
	require.WithinDuration(t, res.CreatedAt, now, 10*time.Millisecond)
	require.WithinDuration(t, res.UpdatedAt, now, 10*time.Millisecond)
}
