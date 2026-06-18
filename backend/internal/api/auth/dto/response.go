package dto

import (
	"time"

	"github.com/google/uuid"

	userdomain "github.com/lionel/my-fitness-aibo/backend/internal/user/domain"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func UserFromDomain(u userdomain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
