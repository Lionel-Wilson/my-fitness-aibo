package mapper

import (
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/plan/domain"
)

func EntityToDomain(e entity.Plan) domain.Plan {
	return domain.Plan{
		ID:          e.ID,
		UserID:      e.UserID,
		Name:        e.Name,
		Quality:     e.Quality,
		Description: e.Description,
		CycleLabel:  e.CycleLabel,
		PeriodStart: e.PeriodStart,
		PeriodEnd:   e.PeriodEnd,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
	}
}

func EntitiesToDomain(plans []entity.Plan) []domain.Plan {
	out := make([]domain.Plan, len(plans))
	for i, p := range plans {
		out[i] = EntityToDomain(p)
	}

	return out
}
