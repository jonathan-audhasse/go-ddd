package models

import (
	"time"

	"github.com/google/uuid"
)

// User is the core domain entity. No infrastructure concerns here.
type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
 