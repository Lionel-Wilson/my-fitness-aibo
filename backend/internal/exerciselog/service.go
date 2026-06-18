package exerciselog

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/mapper"
	exerciselogstorage "github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/storage"
)

type Service interface {
	ListExerciseLogs(ctx context.Context, userID, exerciseID uuid.UUID) ([]domain.ExerciseLog, error)
	UpsertExerciseLog(ctx context.Context, userID uuid.UUID, input domain.UpsertExerciseLogInput) (domain.ExerciseLog, error)
}

type service struct {
	logger *log.Logger
	repo   exerciselogstorage.ExerciseLogRepository
}

func NewService(logger *log.Logger, repo exerciselogstorage.ExerciseLogRepository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) ListExerciseLogs(ctx context.Context, userID, exerciseID uuid.UUID) ([]domain.ExerciseLog, error) {
	logs, err := s.repo.ListExerciseLogs(ctx, userID, exerciseID)
	if err != nil {
		return nil, err
	}

	return mapper.EntitiesToDomain(logs), nil
}

func (s *service) UpsertExerciseLog(ctx context.Context, userID uuid.UUID, input domain.UpsertExerciseLogInput) (domain.ExerciseLog, error) {
	log, err := s.repo.UpsertExerciseLog(ctx, userID, input)
	if err != nil {
		return domain.ExerciseLog{}, err
	}

	return mapper.EntityToDomain(log), nil
}
