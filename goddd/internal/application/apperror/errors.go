package apperror

import (
	"errors"
	"fmt"
	"goddd/internal/domain/repository"
)

var (
	ErrInvalid       = errors.New("invalid input")
	ErrInternal      = errors.New("internal error")
	ErrNotFound      = errors.New("not found")
	ErrUnavailable   = errors.New("unavailable item")
	ErrUnprocessable = errors.New("unprocessable item")
)

// ToAppError map error to application error
func ToAppError(err error) error {
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case errors.Is(err, repository.ErrUserEmailAlreadyExist):
		return fmt.Errorf("%w: %s", ErrUnprocessable, err.Error())
	case errors.Is(err, repository.ErrFailToPing):
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	}

	return ErrInternal
}
