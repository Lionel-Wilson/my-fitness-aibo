package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is an application account.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// CreateUserInput holds signup fields.
type CreateUserInput struct {
	Email        string
	PasswordHash string
}
