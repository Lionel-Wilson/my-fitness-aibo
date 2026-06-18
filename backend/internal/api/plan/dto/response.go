package dto

import (
	"time"

	"github.com/google/uuid"

	plandomain "github.com/lionel/my-fitness-aibo/backend/internal/plan/domain"
)

type PlanResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Quality     string     `json:"quality"`
	Description string     `json:"description"`
	CycleLabel  string     `json:"cycleLabel"`
	PeriodStart *time.Time `json:"periodStart"`
	PeriodEnd   *time.Time `json:"periodEnd"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
}

func PlanFromDomain(p plandomain.Plan) PlanResponse {
	return PlanResponse{
		ID:          p.ID,
		Name:        p.Name,
		Quality:     p.Quality,
		Description: p.Description,
		CycleLabel:  p.CycleLabel,
		PeriodStart: p.PeriodStart,
		PeriodEnd:   p.PeriodEnd,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
	}
}

func PlansFromDomain(plans []plandomain.Plan) []PlanResponse {
	out := make([]PlanResponse, len(plans))
	for i, p := range plans {
		out[i] = PlanFromDomain(p)
	}

	return out
}
