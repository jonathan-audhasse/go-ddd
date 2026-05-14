package repository

import (
	"context"
	"errors"
)

var ErrFailToPing = errors.New("failed to ping repository")

// HealthRepository for health operation
type HealthRepository interface {
	// Ping tries to ping the repository
	Ping(context.Context) error
}
