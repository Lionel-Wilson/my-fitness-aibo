package mapper

import (
	"strings"
	"time"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/plan/dto"
	plandomain "github.com/lionel/my-fitness-aibo/backend/internal/plan/domain"
)

func parseDate(s *string) (*time.Time, error) {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", strings.TrimSpace(*s))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func PlanRequestToDomain(req dto.PlanRequest) (plandomain.PlanInput, error) {
	start, err := parseDate(req.PeriodStart)
	if err != nil {
		return plandomain.PlanInput{}, err
	}

	end, err := parseDate(req.PeriodEnd)
	if err != nil {
		return plandomain.PlanInput{}, err
	}

	cycleLabel := strings.TrimSpace(req.CycleLabel)
	if cycleLabel == "" {
		cycleLabel = "Cycle"
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return plandomain.PlanInput{
		Name:        strings.TrimSpace(req.Name),
		Quality:     req.Quality,
		Description: req.Description,
		CycleLabel:  cycleLabel,
		PeriodStart: start,
		PeriodEnd:   end,
		IsActive:    isActive,
	}, nil
}
