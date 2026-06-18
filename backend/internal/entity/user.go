package entity

import (
	"time"

	"github.com/google/uuid"
)

// User is an application account row.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
