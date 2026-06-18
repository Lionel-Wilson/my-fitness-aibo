package dto

import (
	"time"

	"github.com/google/uuid"

	cycledomain "github.com/lionel/my-fitness-aibo/backend/internal/cycle/domain"
)

type CycleResponse struct {
	ID          uuid.UUID  `json:"id"`
	PlanID      uuid.UUID  `json:"planId"`
	CycleNumber int        `json:"cycleNumber"`
	Label       string     `json:"label"`
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Notes       string     `json:"notes"`
}

func CycleFromDomain(c cycledomain.Cycle) CycleResponse {
	return CycleResponse{
		ID:          c.ID,
		PlanID:      c.PlanID,
		CycleNumber: c.CycleNumber,
		Label:       c.Label,
		StartedAt:   c.StartedAt,
		CompletedAt: c.CompletedAt,
		Notes:       c.Notes,
	}
}

func CyclesFromDomain(cycles []cycledomain.Cycle) []CycleResponse {
	out := make([]CycleResponse, len(cycles))
	for i, c := range cycles {
		out[i] = CycleFromDomain(c)
	}

	return out
}
