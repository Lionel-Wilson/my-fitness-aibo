package entity

import (
	"time"

	"github.com/google/uuid"
)

// ProgressPoint is one cycle's aggregated performance for an exercise.
type ProgressPoint struct {
	CycleID     uuid.UUID
	CycleNumber int
	Label       string
	Side        string
	TopWeightKg float64
	VolumeKg    float64
	BestE1RM    float64
	TotalReps   float64
	StartedAt   time.Time
}
