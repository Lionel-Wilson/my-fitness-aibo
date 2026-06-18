package domain

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID          uuid.UUID
	PlanID      uuid.UUID
	Name        string
	DayLabel    string
	OrderIndex  int
	DurationMin *int
	Notes       string
	CreatedAt   time.Time
}

type WorkoutInput struct {
	Name        string
	DayLabel    string
	OrderIndex  int
	DurationMin *int
	Notes       string
}
