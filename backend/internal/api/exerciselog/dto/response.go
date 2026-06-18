package dto

import (
	"time"

	"github.com/google/uuid"

	exerciselogdomain "github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/domain"
)

type SetLogResponse struct {
	ID       uuid.UUID `json:"id"`
	SetIndex int       `json:"setIndex"`
	Side     string    `json:"side"`
	WeightKg *float64  `json:"weightKg"`
	Reps     *float64  `json:"reps"`
	Rpe      *float64  `json:"rpe"`
	IsDropSet bool     `json:"isDropSet"`
}

type ExerciseLogResponse struct {
	ID              uuid.UUID        `json:"id"`
	ExerciseID      uuid.UUID        `json:"exerciseId"`
	CycleID         uuid.UUID        `json:"cycleId"`
	Note            string           `json:"note"`
	WorkingWeightKg *float64         `json:"workingWeightKg"`
	CreatedAt       time.Time        `json:"createdAt"`
	Sets            []SetLogResponse `json:"sets"`
}

func SetLogFromDomain(s exerciselogdomain.SetLog) SetLogResponse {
	return SetLogResponse{
		ID:        s.ID,
		SetIndex:  s.SetIndex,
		Side:      s.Side,
		WeightKg:  s.WeightKg,
		Reps:      s.Reps,
		Rpe:       s.Rpe,
		IsDropSet: s.IsDropSet,
	}
}

func ExerciseLogFromDomain(l exerciselogdomain.ExerciseLog) ExerciseLogResponse {
	sets := make([]SetLogResponse, len(l.Sets))
	for i, s := range l.Sets {
		sets[i] = SetLogFromDomain(s)
	}

	return ExerciseLogResponse{
		ID:              l.ID,
		ExerciseID:      l.ExerciseID,
		CycleID:         l.CycleID,
		Note:            l.Note,
		WorkingWeightKg: l.WorkingWeightKg,
		CreatedAt:       l.CreatedAt,
		Sets:            sets,
	}
}

func ExerciseLogsFromDomain(logs []exerciselogdomain.ExerciseLog) []ExerciseLogResponse {
	out := make([]ExerciseLogResponse, len(logs))
	for i, l := range logs {
		out[i] = ExerciseLogFromDomain(l)
	}

	return out
}
