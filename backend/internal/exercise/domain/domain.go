package domain

import (
	"time"

	"github.com/google/uuid"
)

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

type ExerciseInput struct {
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
}
