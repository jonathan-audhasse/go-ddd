package apperror_test

import (
	"goddd/internal/application/apperror"
	"goddd/internal/domain/repository"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestToAppError(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		res  error
	}{
		{name: "repository.ErrUserNotFound", err: repository.ErrUserNotFound, res: apperror.ErrNotFound},
		{name: "repository.ErrUserEmailAlreadyExist", err: repository.ErrUserEmailAlreadyExist, res: apperror.ErrUnprocessable},
		{name: "repository.ErrFailToPing", err: repository.ErrFailToPing, res: apperror.ErrUnavailable},
		{name: "repository.ErrFailedToFindUserById", err: repository.ErrFailedToFindUserById, res: apperror.ErrInternal},
		{name: "repository.ErrFailedToCreateUser", err: repository.ErrFailedToCreateUser, res: apperror.ErrInternal},
		{name: "repository.ErrFailedToUpdateUser", err: repository.ErrFailedToUpdateUser, res: apperror.ErrInternal},
		{name: "repository.ErrFailedToDeleteUser", err: repository.ErrFailedToDeleteUser, res: apperror.ErrInternal},
		{name: "repository.ErrFailedToListUsers", err: repository.ErrFailedToListUsers, res: apperror.ErrInternal},
		{name: "random error", err: gofakeit.Error(), res: apperror.ErrInternal},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, apperror.ToAppError(tc.err), tc.res)
		})
	}
}
