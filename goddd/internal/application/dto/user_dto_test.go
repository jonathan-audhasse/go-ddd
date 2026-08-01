package dto_test

import (
	"fmt"
	"goddd/internal/application/dto"
	"goddd/internal/domain/repository"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserRequest_Validate(t *testing.T) {

	testCases := []struct {
		name   string
		req    dto.CreateUserRequest
		expErr error
	}{
		{
			name:   "missing username",
			req:    dto.CreateUserRequest{},
			expErr: dto.ErrMissingUsername,
		},
		{
			name:   "invalid email",
			req:    dto.CreateUserRequest{Username: "john", Email: "invalid-mail"},
			expErr: fmt.Errorf("%w: email='invalid-mail'", dto.ErrInvalidEmail),
		},
		{
			name:   "succeed",
			req:    dto.CreateUserRequest{Username: "john", Email: "john@test.com"},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expErr, tc.req.Validate())
		})
	}
}

func TestCreateUsersRequest_Validate(t *testing.T) {

	testCases := []struct {
		name   string
		req    dto.CreateUsersRequest
		expErr error
	}{
		{
			name:   "missing username",
			req:    []dto.CreateUserRequest{{}},
			expErr: dto.ErrMissingUsername,
		},
		{
			name:   "invalid email",
			req:    []dto.CreateUserRequest{{Username: "john", Email: "invalid-mail"}},
			expErr: fmt.Errorf("%w: email='invalid-mail'", dto.ErrInvalidEmail),
		},
		{
			name:   "succeed",
			req:    []dto.CreateUserRequest{{Username: "john", Email: "john@test.com"}},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expErr, tc.req.Validate())
		})
	}
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
