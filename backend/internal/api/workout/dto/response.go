package dto

import (
	"time"

	"github.com/google/uuid"

	workoutdomain "github.com/lionel/my-fitness-aibo/backend/internal/workout/domain"
)

type WorkoutResponse struct {
	ID          uuid.UUID `json:"id"`
	PlanID      uuid.UUID `json:"planId"`
	Name        string    `json:"name"`
	DayLabel    string    `json:"dayLabel"`
	OrderIndex  int       `json:"orderIndex"`
	DurationMin *int      `json:"durationMin"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"createdAt"`
}

func WorkoutFromDomain(w workoutdomain.Workout) WorkoutResponse {
	return WorkoutResponse{
		ID:          w.ID,
		PlanID:      w.PlanID,
		Name:        w.Name,
		DayLabel:    w.DayLabel,
		OrderIndex:  w.OrderIndex,
		DurationMin: w.DurationMin,
		Notes:       w.Notes,
		CreatedAt:   w.CreatedAt,
	}
}

func WorkoutsFromDomain(workouts []workoutdomain.Workout) []WorkoutResponse {
	out := make([]WorkoutResponse, len(workouts))
	for i, w := range workouts {
		out[i] = WorkoutFromDomain(w)
	}

	return out
}
