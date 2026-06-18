package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/auth/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/user"
	userdomain "github.com/lionel/my-fitness-aibo/backend/internal/user/domain"
	commonauth "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/auth"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service interface {
	Signup(ctx context.Context, creds domain.Credentials) (domain.AuthResult, error)
	Login(ctx context.Context, creds domain.Credentials) (domain.AuthResult, error)
	Me(ctx context.Context, userID uuid.UUID) (domain.MeResult, error)
}

type service struct {
	logger       *log.Logger
	userService  user.Service
	tokenManager *commonauth.TokenManager
}

func NewService(logger *log.Logger, userService user.Service, tokenManager *commonauth.TokenManager) Service {
	return &service{
		logger:       logger,
		userService:  userService,
		tokenManager: tokenManager,
	}
}

func (s *service) Signup(ctx context.Context, creds domain.Credentials) (domain.AuthResult, error) {
	hash, err := commonauth.HashPassword(creds.Password)
	if err != nil {
		return domain.AuthResult{}, fmt.Errorf("hash password: %w", err)
	}

	u, err := s.userService.CreateUser(ctx, userdomain.CreateUserInput{
		Email:        creds.Email,
		PasswordHash: hash,
	})
	if err != nil {
		return domain.AuthResult{}, err
	}

	token, err := s.tokenManager.Issue(u.ID)
	if err != nil {
		return domain.AuthResult{}, fmt.Errorf("issue token: %w", err)
	}

	return domain.AuthResult{Token: token, User: u}, nil
}

func (s *service) Login(ctx context.Context, creds domain.Credentials) (domain.AuthResult, error) {
	u, err := s.userService.GetUserByEmail(ctx, creds.Email)
	if err != nil {
		if errors.Is(err, commonErrors.ErrNotFound) {
			return domain.AuthResult{}, ErrInvalidCredentials
		}

		return domain.AuthResult{}, err
	}

	if !commonauth.CheckPassword(u.PasswordHash, creds.Password) {
		return domain.AuthResult{}, ErrInvalidCredentials
	}

	token, err := s.tokenManager.Issue(u.ID)
	if err != nil {
		return domain.AuthResult{}, fmt.Errorf("issue token: %w", err)
	}

	return domain.AuthResult{Token: token, User: u}, nil
}

func (s *service) Me(ctx context.Context, userID uuid.UUID) (domain.MeResult, error) {
	u, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return domain.MeResult{}, err
	}

	return domain.MeResult{User: u}, nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
