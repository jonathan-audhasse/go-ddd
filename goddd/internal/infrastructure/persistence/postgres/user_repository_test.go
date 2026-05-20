package postgres_test

import (
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"
	"goddd/internal/infrastructure/persistence/postgres"
	"goddd/internal/infrastructure/persistence/postgres/sqlcgen"
	"goddd/internal/infrastructure/persistence/testutils"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())

	repo := postgres.NewUserRepository(testDB)

	_, err := repo.FindByID(ctx, uuid.MustParse(gofakeit.UUID()))
	assert.ErrorIs(t, err, repository.ErrUserNotFound)
}

func TestUserRepository_FindByID(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, tx, rollback := testutils.WithTx(t, ctx, testDB)

	defer rollback()

	q := sqlcgen.New(tx)

	// insert user...
	user, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)
	res, err := repo.FindByID(ctx, user.ID.Bytes)
	require.NoError(t, err)

	// assert
	assert.Equal(t, uuid.UUID(user.ID.Bytes), res.ID)
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Username, res.Username)
}

func TestUserRepository_FindPaged(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, _, rollback := testutils.WithTx(t, ctx, testDB)

	defer rollback()

	repo := postgres.NewUserRepository(testDB)
	users := []models.NewUser{
		{Email: gofakeit.Email(), Username: gofakeit.Username()},
		{Email: gofakeit.Email(), Username: gofakeit.Username()},
		{Email: gofakeit.Email(), Username: gofakeit.Username()},
	}

	// add users
	_, err := repo.BulkCreates(ctx, users)
	require.NoError(t, err)

	// list a single item
	page := repository.Page{Limit: 1}
	res, err := repo.FindPaged(ctx, page)
	require.NoError(t, err)

	// assertion
	require.Equal(t, 1, len(res.Users))
	require.NotNil(t, res.NextCursor)
	cursor1 := res.Users[0].ID
	assert.Equal(t, cursor1, *res.NextCursor)

	// list the rest
	page = repository.Page{Limit: 99, Cursor: res.Users[0].ID}
	res, err = repo.FindPaged(ctx, page)
	require.NoError(t, err)

	// assertion
	assert.Equal(t, 2, len(res.Users))
	assert.Nil(t, res.NextCursor)
	assert.NotContains(t, lo.Map(res.Users, func(u models.User, _ int) uuid.UUID { return u.ID }), cursor1)
}

func TestUserRepository_Create(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, _, rollback := testutils.WithTx(t, ctx, testDB)

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
	assert.Equal(t, user.Email, res.Email)
	assert.Equal(t, user.Username, res.Username)
	assert.WithinDuration(t, res.CreatedAt, now, 10*time.Millisecond)
	assert.WithinDuration(t, res.UpdatedAt, now, 10*time.Millisecond)

	// try add the same user twice
	_, err = repo.Create(ctx, user)
	require.Error(t, err, repository.ErrUserEmailAlreadyExist)
}

func TestUserRepository_Update_Unknown_User(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	repo := postgres.NewUserRepository(testDB)

	user := models.User{
		ID:       uuid.MustParse(gofakeit.UUID()),
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	}

	// update the new user
	_, err := repo.Update(ctx, user)
	require.Error(t, err)
}

func TestUserRepository_Update(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, tx, rollback := testutils.WithTx(t, ctx, testDB)

	defer rollback()

	q := sqlcgen.New(tx)

	// insert user...
	newUser, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)

	user := models.User{
		ID:       newUser.ID.Bytes,
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	}

	// update the new user
	u, err := repo.Update(ctx, user)
	require.NoError(t, err)
	require.NotEmpty(t, u)

	// assert
	assert.Equal(t, user.Email, u.Email)
	assert.Equal(t, user.Username, u.Username)
}

func TestUserRepository_Update_Email_Already_Exists(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, tx, rollback := testutils.WithTx(t, ctx, testDB)

	defer rollback()

	q := sqlcgen.New(tx)

	// insert users ...
	newUser1, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)
	newUser2, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)

	user := models.User{
		ID:    newUser1.ID.Bytes,
		Email: newUser2.Email,
	}

	// try to update the user
	_, err = repo.Update(ctx, user)
	require.ErrorIs(t, err, repository.ErrUserEmailAlreadyExist)
}

func TestUserRepository_Delete_Unknown_User(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())

	repo := postgres.NewUserRepository(testDB)

	// delete
	err := repo.Delete(ctx, uuid.MustParse(gofakeit.UUID()))
	// assert
	assert.NoError(t, err)
}

func TestUserRepository_Delete(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())
	ctx, tx, rollback := testutils.WithTx(t, ctx, testDB)

	defer rollback()

	q := sqlcgen.New(tx)

	// insert user...
	user, err := q.CreateUser(ctx, &sqlcgen.CreateUserParams{
		Email:    gofakeit.Email(),
		Username: gofakeit.Username(),
	})
	require.NoError(t, err)

	repo := postgres.NewUserRepository(testDB)
	_, err = repo.FindByID(ctx, user.ID.Bytes)
	require.NoError(t, err)

	// delete
	err = repo.Delete(ctx, user.ID.Bytes)

	// assert user is not found
	_, err = repo.FindByID(ctx, user.ID.Bytes)
	assert.ErrorIs(t, err, repository.ErrUserNotFound)
}

func TestUserRepository_BulkCreates(t *testing.T) {
	ctx := log.Logger.WithContext(t.Context())

	// tx isolation
	ctx, tx, rollback := testutils.WithTx(t, ctx, testDB)
	defer rollback()

	q := sqlcgen.New(tx)

	// assert empty before added
	total, err := q.CountUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(0), total)

	repo := postgres.NewUserRepository(testDB)

	users := []models.NewUser{
		{Email: gofakeit.Email(), Username: gofakeit.Username()},
		{Email: gofakeit.Email(), Username: gofakeit.Username()},
		{Email: gofakeit.Email(), Username: gofakeit.Username()},
	}

	// add users
	count, err := repo.BulkCreates(ctx, users)
	require.NoError(t, err)
	require.Equal(t, int64(3), count)

	// assert users has been added
	total, err = q.CountUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
}
