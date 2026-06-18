package domain

import (
	"github.com/google/uuid"

	userdomain "github.com/lionel/my-fitness-aibo/backend/internal/user/domain"
)

type Credentials struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string
	User  userdomain.User
}

type MeResult struct {
	User userdomain.User
}

type UserID uuid.UUID
