package mapper

import (
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/domain"
)

func EntityToDomain(e entity.ExerciseLog) domain.ExerciseLog {
	sets := make([]domain.SetLog, len(e.Sets))
	for i, s := range e.Sets {
		sets[i] = SetLogEntityToDomain(s)
	}

	return domain.ExerciseLog{
		ID:              e.ID,
		ExerciseID:      e.ExerciseID,
		CycleID:         e.CycleID,
		Note:            e.Note,
		WorkingWeightKg: e.WorkingWeightKg,
		CreatedAt:       e.CreatedAt,
		Sets:            sets,
	}
}

func SetLogEntityToDomain(s entity.SetLog) domain.SetLog {
	return domain.SetLog{
		ID:            s.ID,
		ExerciseLogID: s.ExerciseLogID,
		SetIndex:      s.SetIndex,
		Side:          s.Side,
		WeightKg:      s.WeightKg,
		Reps:          s.Reps,
		Rpe:           s.Rpe,
		IsDropSet:     s.IsDropSet,
	}
}

func EntitiesToDomain(logs []entity.ExerciseLog) []domain.ExerciseLog {
	out := make([]domain.ExerciseLog, len(logs))
	for i, l := range logs {
		out[i] = EntityToDomain(l)
	}

	return out
}
