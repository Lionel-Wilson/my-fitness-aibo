package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cycle struct {
	ID          uuid.UUID
	PlanID      uuid.UUID
	CycleNumber int
	Label       string
	StartedAt   time.Time
	CompletedAt *time.Time
	Notes       string
}

type CreateCycleInput struct {
	Label string
	Notes string
}

type UpdateCycleInput struct {
	Label       string
	Notes       string
	CompletedAt *time.Time
}
