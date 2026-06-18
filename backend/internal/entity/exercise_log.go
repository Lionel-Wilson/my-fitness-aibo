package entity

import (
	"time"

	"github.com/google/uuid"
)

// SetLog is one logged set within an exercise log.
type SetLog struct {
	ID            uuid.UUID
	ExerciseLogID uuid.UUID
	SetIndex      int
	Side          string
	WeightKg      *float64
	Reps          *float64
	Rpe           *float64
	IsDropSet     bool
}

// ExerciseLog is an exercise's record for one cycle.
type ExerciseLog struct {
	ID              uuid.UUID
	ExerciseID      uuid.UUID
	CycleID         uuid.UUID
	Note            string
	WorkingWeightKg *float64
	CreatedAt       time.Time
	Sets            []SetLog
}
