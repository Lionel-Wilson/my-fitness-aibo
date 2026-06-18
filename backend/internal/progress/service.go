package progress

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/progress/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/progress/mapper"
	progressstorage "github.com/lionel/my-fitness-aibo/backend/internal/progress/storage"
)

type Service interface {
	ExerciseProgress(ctx context.Context, userID, exerciseID uuid.UUID) ([]domain.ProgressPoint, error)
}

type service struct {
	logger *log.Logger
	repo   progressstorage.ProgressRepository
}

func NewService(logger *log.Logger, repo progressstorage.ProgressRepository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) ExerciseProgress(ctx context.Context, userID, exerciseID uuid.UUID) ([]domain.ProgressPoint, error) {
	points, err := s.repo.ExerciseProgress(ctx, userID, exerciseID)
	if err != nil {
		return nil, err
	}

	return mapper.EntitiesToDomain(points), nil
}
