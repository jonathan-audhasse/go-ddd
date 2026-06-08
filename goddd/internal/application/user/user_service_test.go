package user_test

import (
	"context"
	"goddd/internal/application/apperror"
	"goddd/internal/application/dto"
	user "goddd/internal/application/user"
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"
	mockrepo "goddd/tests/mocks/repository"
	mocktm "goddd/tests/mocks/transaction"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"
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

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()
	mockReq := dto.ListUsersRequest{
		Limit:  gofakeit.Uint(),
		Cursor: lo.ToPtr(uuid.NewString()),
	}
	mockUsers := make([]models.User, 2)
	gofakeit.Struct(&mockUsers[0])
	gofakeit.Struct(&mockUsers[1])

	mockNextCursor := uuid.New()

	testCases := []struct {
		name      string
		req       dto.ListUsersRequest
		findPaged func(context.Context, repository.Page) (repository.PagedResult, error)
		expErr    error
		expRes    *dto.ListUsersResponse
	}{
		{
			name:   "failed to parse request",
			req:    dto.ListUsersRequest{Cursor: lo.ToPtr("invalid-id-format")},
			expErr: dto.ErrInvalidCursor,
		},
		{
			name: "failed to find page from repo",
			req:  mockReq,
			findPaged: func(context.Context, repository.Page) (repository.PagedResult, error) {
				return repository.PagedResult{}, gofakeit.Error()
			},
			expErr: apperror.ErrInternal,
		},
		{
			name: "succeed to retrieve from repo, no next cursor",
			req:  mockReq,
			findPaged: func(_ context.Context, page repository.Page) (repository.PagedResult, error) {
				require.Equal(t, mockReq.Limit, page.Limit)
				require.Equal(t, *mockReq.Cursor, page.Cursor.String())
				return repository.PagedResult{Users: mockUsers}, nil
			},
			expRes: &dto.ListUsersResponse{
				Users:      mockUsers,
				Limit:      mockReq.Limit,
				NextCursor: nil,
			},
		},
		{
			name: "succeed to retrieve from repo, next cursor provided",
			req:  mockReq,
			findPaged: func(_ context.Context, page repository.Page) (repository.PagedResult, error) {
				require.Equal(t, mockReq.Limit, page.Limit)
				require.Equal(t, *mockReq.Cursor, page.Cursor.String())
				return repository.PagedResult{Users: mockUsers, NextCursor: &mockNextCursor}, nil
			},
			expRes: &dto.ListUsersResponse{
				Users:      mockUsers,
				Limit:      mockReq.Limit,
				NextCursor: lo.ToPtr(mockNextCursor.String()),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// set mock behavior
			repo := mockrepo.NewMockUserRepository(t)
			if tc.findPaged != nil {
				repo.EXPECT().FindPaged(ctx, mock.AnythingOfType("repository.Page")).RunAndReturn(tc.findPaged)
			}
			srv := user.NewUserService(repo, nil)

			// call service
			res, err := srv.ListUsers(ctx, tc.req)
			require.Equal(t, tc.expErr, err)

			if tc.expErr != nil {
				return
			}

			assert.Equal(t, tc.expRes, res)
		})
	}
}
