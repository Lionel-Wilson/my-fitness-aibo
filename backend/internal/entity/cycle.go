package entity

import (
	"time"

	"github.com/google/uuid"
)

// Cycle is the plan-level progression unit.
type Cycle struct {
	ID          uuid.UUID
	PlanID      uuid.UUID
	CycleNumber int
	Label       string
	StartedAt   time.Time
	CompletedAt *time.Time
	Notes       string
}
