package mapper

import (
	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/exerciselog/dto"
	exerciselogdomain "github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/domain"
)

func UpsertLogRequestToDomain(exerciseID, cycleID uuid.UUID, req dto.UpsertLogRequest) exerciselogdomain.UpsertExerciseLogInput {
	sets := make([]exerciselogdomain.SetLogInput, len(req.Sets))
	for i, s := range req.Sets {
		sets[i] = exerciselogdomain.SetLogInput{
			SetIndex:  s.SetIndex,
			Side:      s.Side,
			WeightKg:  s.WeightKg,
			Reps:      s.Reps,
			Rpe:       s.Rpe,
			IsDropSet: s.IsDropSet,
		}
	}

	return exerciselogdomain.UpsertExerciseLogInput{
		ExerciseID:      exerciseID,
		CycleID:         cycleID,
		Note:            req.Note,
		WorkingWeightKg: req.WorkingWeightKg,
		Sets:            sets,
	}
}
