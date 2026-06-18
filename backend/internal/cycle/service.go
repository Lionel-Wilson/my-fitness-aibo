package cycle

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/cycle/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/cycle/mapper"
	cyclestorage "github.com/lionel/my-fitness-aibo/backend/internal/cycle/storage"
)

type Service interface {
	ListCycles(ctx context.Context, userID, planID uuid.UUID) ([]domain.Cycle, error)
	CreateCycle(ctx context.Context, userID, planID uuid.UUID, input domain.CreateCycleInput) (domain.Cycle, error)
	UpdateCycle(ctx context.Context, userID, id uuid.UUID, input domain.UpdateCycleInput) (domain.Cycle, error)
	DeleteCycle(ctx context.Context, userID, id uuid.UUID) error
}

type service struct {
	logger *log.Logger
	repo   cyclestorage.CycleRepository
}

func NewService(logger *log.Logger, repo cyclestorage.CycleRepository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) ListCycles(ctx context.Context, userID, planID uuid.UUID) ([]domain.Cycle, error) {
	cycles, err := s.repo.ListCycles(ctx, userID, planID)
	if err != nil {
		return nil, err
	}

	return mapper.EntitiesToDomain(cycles), nil
}

func (s *service) CreateCycle(ctx context.Context, userID, planID uuid.UUID, input domain.CreateCycleInput) (domain.Cycle, error) {
	c, err := s.repo.CreateCycle(ctx, userID, planID, input)
	if err != nil {
		return domain.Cycle{}, err
	}

	return mapper.EntityToDomain(c), nil
}

func (s *service) UpdateCycle(ctx context.Context, userID, id uuid.UUID, input domain.UpdateCycleInput) (domain.Cycle, error) {
	c, err := s.repo.UpdateCycle(ctx, userID, id, input)
	if err != nil {
		return domain.Cycle{}, err
	}

	return mapper.EntityToDomain(c), nil
}

func (s *service) DeleteCycle(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.DeleteCycle(ctx, userID, id)
}
