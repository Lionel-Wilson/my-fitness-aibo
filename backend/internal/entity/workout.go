package entity

import (
	"time"

	"github.com/google/uuid"
)

// Workout is a single training day within a plan.
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
