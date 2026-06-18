package mapper

import (
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/workout/domain"
)

func EntityToDomain(e entity.Workout) domain.Workout {
	return domain.Workout{
		ID:          e.ID,
		PlanID:      e.PlanID,
		Name:        e.Name,
		DayLabel:    e.DayLabel,
		OrderIndex:  e.OrderIndex,
		DurationMin: e.DurationMin,
		Notes:       e.Notes,
		CreatedAt:   e.CreatedAt,
	}
}

func EntitiesToDomain(workouts []entity.Workout) []domain.Workout {
	out := make([]domain.Workout, len(workouts))
	for i, w := range workouts {
		out[i] = EntityToDomain(w)
	}

	return out
}
