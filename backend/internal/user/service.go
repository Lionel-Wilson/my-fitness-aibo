package user

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/user/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/user/mapper"
	userstorage "github.com/lionel/my-fitness-aibo/backend/internal/user/storage"
)

type Service interface {
	CreateUser(ctx context.Context, input domain.CreateUserInput) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type service struct {
	logger *log.Logger
	repo   userstorage.UserRepository
}

func NewService(logger *log.Logger, repo userstorage.UserRepository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) CreateUser(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	u, err := s.repo.CreateUser(ctx, input.Email, input.PasswordHash)
	if err != nil {
		return domain.User{}, err
	}

	return mapper.EntityToDomain(u), nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return mapper.EntityToDomain(u), nil
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return mapper.EntityToDomain(u), nil
}
