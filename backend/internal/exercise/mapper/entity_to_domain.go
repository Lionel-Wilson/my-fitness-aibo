package mapper

import (
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/exercise/domain"
)

func EntityToDomain(e entity.Exercise) domain.Exercise {
	return domain.Exercise{
		ID:           e.ID,
		WorkoutID:    e.WorkoutID,
		Name:         e.Name,
		OrderIndex:   e.OrderIndex,
		TargetSets:   e.TargetSets,
		RepLow:       e.RepLow,
		RepHigh:      e.RepHigh,
		RpeLow:       e.RpeLow,
		RpeHigh:      e.RpeHigh,
		RestSeconds:  e.RestSeconds,
		Instructions: e.Instructions,
		OrGroup:      e.OrGroup,
		IsOptional:   e.IsOptional,
		IsUnilateral: e.IsUnilateral,
		CreatedAt:    e.CreatedAt,
	}
}

func EntitiesToDomain(exercises []entity.Exercise) []domain.Exercise {
	out := make([]domain.Exercise, len(exercises))
	for i, e := range exercises {
		out[i] = EntityToDomain(e)
	}

	return out
}
