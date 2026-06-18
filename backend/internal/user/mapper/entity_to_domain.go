package mapper

import (
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/user/domain"
)

func EntityToDomain(e entity.User) domain.User {
	return domain.User{
		ID:           e.ID,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		CreatedAt:    e.CreatedAt,
	}
}

func EntitiesToDomain(users []entity.User) []domain.User {
	out := make([]domain.User, len(users))
	for i, u := range users {
		out[i] = EntityToDomain(u)
	}

	return out
}
