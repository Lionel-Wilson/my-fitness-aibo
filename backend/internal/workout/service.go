package workout

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/workout/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/workout/mapper"
	workoutstorage "github.com/lionel/my-fitness-aibo/backend/internal/workout/storage"
)

type Service interface {
	ListWorkouts(ctx context.Context, userID, planID uuid.UUID) ([]domain.Workout, error)
	CreateWorkout(ctx context.Context, userID, planID uuid.UUID, input domain.WorkoutInput) (domain.Workout, error)
	GetWorkout(ctx context.Context, userID, id uuid.UUID) (domain.Workout, error)
	UpdateWorkout(ctx context.Context, userID, id uuid.UUID, input domain.WorkoutInput) (domain.Workout, error)
	DeleteWorkout(ctx context.Context, userID, id uuid.UUID) error
}

type service struct {
	logger *log.Logger
	repo   workoutstorage.WorkoutRepository
}

func NewService(logger *log.Logger, repo workoutstorage.WorkoutRepository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) ListWorkouts(ctx context.Context, userID, planID uuid.UUID) ([]domain.Workout, error) {
	workouts, err := s.repo.ListWorkouts(ctx, userID, planID)
	if err != nil {
		return nil, err
	}

	return mapper.EntitiesToDomain(workouts), nil
}

func (s *service) CreateWorkout(ctx context.Context, userID, planID uuid.UUID, input domain.WorkoutInput) (domain.Workout, error) {
	w, err := s.repo.CreateWorkout(ctx, userID, planID, input)
	if err != nil {
		return domain.Workout{}, err
	}

	return mapper.EntityToDomain(w), nil
}

func (s *service) GetWorkout(ctx context.Context, userID, id uuid.UUID) (domain.Workout, error) {
	w, err := s.repo.GetWorkout(ctx, userID, id)
	if err != nil {
		return domain.Workout{}, err
	}

	return mapper.EntityToDomain(w), nil
}

func (s *service) UpdateWorkout(ctx context.Context, userID, id uuid.UUID, input domain.WorkoutInput) (domain.Workout, error) {
	w, err := s.repo.UpdateWorkout(ctx, userID, id, input)
	if err != nil {
		return domain.Workout{}, err
	}

	return mapper.EntityToDomain(w), nil
}

func (s *service) DeleteWorkout(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.DeleteWorkout(ctx, userID, id)
}
