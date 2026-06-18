package dto

import (
	"time"

	"github.com/google/uuid"

	exercisedomain "github.com/lionel/my-fitness-aibo/backend/internal/exercise/domain"
)

type ExerciseResponse struct {
	ID           uuid.UUID `json:"id"`
	WorkoutID    uuid.UUID `json:"workoutId"`
	Name         string    `json:"name"`
	OrderIndex   int       `json:"orderIndex"`
	TargetSets   *int      `json:"targetSets"`
	RepLow       *int      `json:"repLow"`
	RepHigh      *int      `json:"repHigh"`
	RpeLow       *float64  `json:"rpeLow"`
	RpeHigh      *float64  `json:"rpeHigh"`
	RestSeconds  *int      `json:"restSeconds"`
	Instructions string    `json:"instructions"`
	OrGroup      string    `json:"orGroup"`
	IsOptional   bool      `json:"isOptional"`
	IsUnilateral bool      `json:"isUnilateral"`
	CreatedAt    time.Time `json:"createdAt"`
}

func ExerciseFromDomain(e exercisedomain.Exercise) ExerciseResponse {
	return ExerciseResponse{
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

func ExercisesFromDomain(exercises []exercisedomain.Exercise) []ExerciseResponse {
	out := make([]ExerciseResponse, len(exercises))
	for i, e := range exercises {
		out[i] = ExerciseFromDomain(e)
	}

	return out
}
