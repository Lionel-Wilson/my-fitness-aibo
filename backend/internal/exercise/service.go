package exercise

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/exercise/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/exercise/mapper"
	exercisestorage "github.com/lionel/my-fitness-aibo/backend/internal/exercise/storage"
)

type Service interface {
	ListExercises(ctx context.Context, userID, workoutID uuid.UUID) ([]domain.Exercise, error)
	CreateExercise(ctx context.Context, userID, workoutID uuid.UUID, input domain.ExerciseInput) (domain.Exercise, error)
	GetExercise(ctx context.Context, userID, id uuid.UUID) (domain.Exercise, error)
	UpdateExercise(ctx context.Context, userID, id uuid.UUID, input domain.ExerciseInput) (domain.Exercise, error)
	DeleteExercise(ctx context.Context, userID, id uuid.UUID) error
}

type service struct {
	logger *log.Logger
	repo   exercisestorage.ExerciseRepository
}

func NewService(logger *log.Logger, repo exercisestorage.ExerciseRepository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) ListExercises(ctx context.Context, userID, workoutID uuid.UUID) ([]domain.Exercise, error) {
	exercises, err := s.repo.ListExercises(ctx, userID, workoutID)
	if err != nil {
		return nil, err
	}

	return mapper.EntitiesToDomain(exercises), nil
}

func (s *service) CreateExercise(ctx context.Context, userID, workoutID uuid.UUID, input domain.ExerciseInput) (domain.Exercise, error) {
	e, err := s.repo.CreateExercise(ctx, userID, workoutID, input)
	if err != nil {
		return domain.Exercise{}, err
	}

	return mapper.EntityToDomain(e), nil
}

func (s *service) GetExercise(ctx context.Context, userID, id uuid.UUID) (domain.Exercise, error) {
	e, err := s.repo.GetExercise(ctx, userID, id)
	if err != nil {
		return domain.Exercise{}, err
	}

	return mapper.EntityToDomain(e), nil
}

func (s *service) UpdateExercise(ctx context.Context, userID, id uuid.UUID, input domain.ExerciseInput) (domain.Exercise, error) {
	e, err := s.repo.UpdateExercise(ctx, userID, id, input)
	if err != nil {
		return domain.Exercise{}, err
	}

	return mapper.EntityToDomain(e), nil
}

func (s *service) DeleteExercise(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.DeleteExercise(ctx, userID, id)
}
