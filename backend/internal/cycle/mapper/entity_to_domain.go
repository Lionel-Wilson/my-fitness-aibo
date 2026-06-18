package mapper

import (
	"github.com/lionel/my-fitness-aibo/backend/internal/cycle/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
)

func EntityToDomain(e entity.Cycle) domain.Cycle {
	return domain.Cycle{
		ID:          e.ID,
		PlanID:      e.PlanID,
		CycleNumber: e.CycleNumber,
		Label:       e.Label,
		StartedAt:   e.StartedAt,
		CompletedAt: e.CompletedAt,
		Notes:       e.Notes,
	}
}

func EntitiesToDomain(cycles []entity.Cycle) []domain.Cycle {
	out := make([]domain.Cycle, len(cycles))
	for i, c := range cycles {
		out[i] = EntityToDomain(c)
	}

	return out
}
