package entity

import (
	"time"

	"github.com/google/uuid"
)

// Exercise is a movement within a workout.
type Exercise struct {
	ID           uuid.UUID
	WorkoutID    uuid.UUID
	Name         string
	OrderIndex   int
	TargetSets   *int
	RepLow       *int
	RepHigh      *int
	RpeLow       *float64
	RpeHigh      *float64
	RestSeconds  *int
	Instructions string
	OrGroup      string
	IsOptional   bool
	IsUnilateral bool
	CreatedAt    time.Time
}
