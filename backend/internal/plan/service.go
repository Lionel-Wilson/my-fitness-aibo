package plan

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/plan/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/plan/mapper"
	planstorage "github.com/lionel/my-fitness-aibo/backend/internal/plan/storage"
)

type Service interface {
	ListPlans(ctx context.Context, userID uuid.UUID) ([]domain.Plan, error)
	CreatePlan(ctx context.Context, userID uuid.UUID, input domain.PlanInput) (domain.Plan, error)
	GetPlan(ctx context.Context, userID, id uuid.UUID) (domain.Plan, error)
	UpdatePlan(ctx context.Context, userID, id uuid.UUID, input domain.PlanInput) (domain.Plan, error)
	DeletePlan(ctx context.Context, userID, id uuid.UUID) error
}

type service struct {
	logger *log.Logger
	repo   planstorage.PlanRepository
}

func NewService(logger *log.Logger, repo planstorage.PlanRepository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) ListPlans(ctx context.Context, userID uuid.UUID) ([]domain.Plan, error) {
	plans, err := s.repo.ListPlans(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapper.EntitiesToDomain(plans), nil
}

func (s *service) CreatePlan(ctx context.Context, userID uuid.UUID, input domain.PlanInput) (domain.Plan, error) {
	p, err := s.repo.CreatePlan(ctx, userID, input)
	if err != nil {
		return domain.Plan{}, err
	}

	return mapper.EntityToDomain(p), nil
}

func (s *service) GetPlan(ctx context.Context, userID, id uuid.UUID) (domain.Plan, error) {
	p, err := s.repo.GetPlan(ctx, userID, id)
	if err != nil {
		return domain.Plan{}, err
	}

	return mapper.EntityToDomain(p), nil
}

func (s *service) UpdatePlan(ctx context.Context, userID, id uuid.UUID, input domain.PlanInput) (domain.Plan, error) {
	p, err := s.repo.UpdatePlan(ctx, userID, id, input)
	if err != nil {
		return domain.Plan{}, err
	}

	return mapper.EntityToDomain(p), nil
}

func (s *service) DeletePlan(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.DeletePlan(ctx, userID, id)
}
