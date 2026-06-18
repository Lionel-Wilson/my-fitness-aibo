package mapper

import (
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/progress/domain"
)

func EntityToDomain(e entity.ProgressPoint) domain.ProgressPoint {
	return domain.ProgressPoint{
		CycleID:     e.CycleID,
		CycleNumber: e.CycleNumber,
		Label:       e.Label,
		Side:        e.Side,
		TopWeightKg: e.TopWeightKg,
		VolumeKg:    e.VolumeKg,
		BestE1RM:    e.BestE1RM,
		TotalReps:   e.TotalReps,
		StartedAt:   e.StartedAt,
	}
}

func EntitiesToDomain(points []entity.ProgressPoint) []domain.ProgressPoint {
	out := make([]domain.ProgressPoint, len(points))
	for i, p := range points {
		out[i] = EntityToDomain(p)
	}

	return out
}
