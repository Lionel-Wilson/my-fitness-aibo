package domain

import (
	"time"

	"github.com/google/uuid"
)

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
