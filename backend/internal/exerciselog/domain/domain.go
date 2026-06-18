package domain

import (
	"time"

	"github.com/google/uuid"
)

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

type SetLogInput struct {
	SetIndex  int
	Side      string
	WeightKg  *float64
	Reps      *float64
	Rpe       *float64
	IsDropSet bool
}

type ExerciseLog struct {
	ID              uuid.UUID
	ExerciseID      uuid.UUID
	CycleID         uuid.UUID
	Note            string
	WorkingWeightKg *float64
	CreatedAt       time.Time
	Sets            []SetLog
}

type UpsertExerciseLogInput struct {
	ExerciseID      uuid.UUID
	CycleID         uuid.UUID
	Note            string
	WorkingWeightKg *float64
	Sets            []SetLogInput
}
