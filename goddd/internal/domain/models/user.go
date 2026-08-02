package models

import (
	"time"

	"github.com/google/uuid"
)

// NewUser to create new users
type NewUser struct {
	Email    string
	Username string
}

// User is the core domain entity.
type User struct {
	ID        uuid.UUID
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
