package dto_test

import (
	"goddd/internal/application/dto"
	"goddd/internal/domain/models"
	"goddd/internal/domain/repository"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserRequest_ToNewUser(t *testing.T) {
	req := dto.CreateUserRequest{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
	}
	expRes := models.NewUser{
		Username: req.Username,
		Email:    req.Email,
	}
	assert.Equal(t, expRes, req.ToNewUser())
}

func TestListUsersRequest_ToNewUser(t *testing.T) {
	randLimit := gofakeit.Uint()
	invalidCursor := lo.ToPtr("invalid-id-format")
	validId := uuid.New()

	testCases := []struct {
		name   string
		req    dto.ListUsersRequest
		expRes repository.Page
		expErr error
	}{
		{
			name:   "no cursor provided",
			req:    dto.ListUsersRequest{Limit: randLimit},
			expRes: repository.Page{Limit: randLimit, Cursor: uuid.Nil},
		},
		{
			name:   "failed to parse cursor to uuid",
			req:    dto.ListUsersRequest{Limit: randLimit, Cursor: invalidCursor},
			expErr: dto.ErrInvalidCursor,
		},
		{
			name:   "valid cursor id format",
			req:    dto.ListUsersRequest{Limit: randLimit, Cursor: lo.ToPtr(validId.String())},
			expRes: repository.Page{Limit: randLimit, Cursor: validId},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.req.ToPage()
			require.Equal(t, tc.expErr, err)

			if tc.expErr == nil {
				assert.Equal(t, tc.expRes, res)
			}
		})
	}

}
