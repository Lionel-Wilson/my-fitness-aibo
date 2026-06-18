package domain

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Quality     string
	Description string
	CycleLabel  string
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	IsActive    bool
	CreatedAt   time.Time
}

type PlanInput struct {
	Name        string
	Quality     string
	Description string
	CycleLabel  string
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	IsActive    bool
}
