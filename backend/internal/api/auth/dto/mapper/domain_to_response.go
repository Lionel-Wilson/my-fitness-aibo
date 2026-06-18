package mapper

import (
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/auth/dto"
	"github.com/lionel/my-fitness-aibo/backend/internal/auth/domain"
)

func CredentialsRequestToDomain(req dto.CredentialsRequest) domain.Credentials {
	return domain.Credentials{
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: req.Password,
	}
}

func AuthResultToResponse(result domain.AuthResult) dto.AuthResponse {
	return dto.AuthResponse{
		Token: result.Token,
		User:  dto.UserFromDomain(result.User),
	}
}

func MeResultToResponse(result domain.MeResult) dto.UserResponse {
	return dto.UserFromDomain(result.User)
}
