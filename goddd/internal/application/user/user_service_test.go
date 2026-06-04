package user_test

import (
	"context"
	"goddd/internal/application/apperror"
	"goddd/internal/application/dto"
	user "goddd/internal/application/user"
	"goddd/internal/domain/models"
	mockrepo "goddd/tests/mocks/repository"
	mocktm "goddd/tests/mocks/transaction"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateNewUser(t *testing.T) {
	ctx := context.Background()

	tm := mocktm.NewMockTransactionManager(t)

	mockErr := gofakeit.Error()
	expUser := models.User{}
	req := dto.CreateUserRequest{}

	gofakeit.Struct(&expUser)
	gofakeit.Struct(&req)

	tm.EXPECT().
		Do(ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})

	testCases := []struct {
		name   string
		create func(context.Context, models.NewUser) (models.User, error)
		expErr error
	}{
		{
			name: "failed to insert into repository",
			create: func(context.Context, models.NewUser) (models.User, error) {
				return models.User{}, mockErr
			},
			expErr: apperror.ErrInternal,
		},
		{
			name: "succeed to create new user",
			create: func(_ context.Context, u models.NewUser) (models.User, error) {
				require.Equal(t, models.NewUser{Email: req.Email, Username: req.Username}, u)
				return expUser, nil
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mockrepo.NewMockUserRepository(t)
			repo.EXPECT().Create(ctx, mock.AnythingOfType("models.NewUser")).RunAndReturn(tc.create)
			srv := user.NewUserService(repo, tm)

			res, err := srv.CreateNewUser(ctx, req)
			require.Equal(t, tc.expErr, err)

			if tc.expErr != nil {
				return
			}

			assert.Equal(t, expUser, res)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()

	var expUser models.User
	gofakeit.Struct(&expUser)

	id := uuid.New()

	testCases := []struct {
		name    string
		mockErr error
		expErr  error
	}{
		{
			name:    "failed to insert into repository",
			mockErr: gofakeit.Error(),
			expErr:  apperror.ErrInternal,
		},
		{
			name:    "succeed to create new user",
			mockErr: nil,
			expErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := mockrepo.NewMockUserRepository(t)
			repo.EXPECT().FindByID(mock.Anything, id).Return(&expUser, tc.mockErr)
			srv := user.NewUserService(repo, nil)

			res, err := srv.GetUser(ctx, id)
			require.Equal(t, tc.expErr, err)

			if tc.expErr != nil {
				return
			}

			require.Equal(t, &expUser, res)
		})
	}
}
