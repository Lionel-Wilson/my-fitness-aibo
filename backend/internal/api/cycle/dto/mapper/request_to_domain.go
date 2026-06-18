package mapper

import (
	"time"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/cycle/dto"
	cycledomain "github.com/lionel/my-fitness-aibo/backend/internal/cycle/domain"
)

func CreateCycleRequestToDomain(req dto.CreateCycleRequest) cycledomain.CreateCycleInput {
	return cycledomain.CreateCycleInput{
		Label: req.Label,
		Notes: req.Notes,
	}
}

func UpdateCycleRequestToDomain(req dto.UpdateCycleRequest) cycledomain.UpdateCycleInput {
	var completedAt *time.Time
	if req.Completed != nil && *req.Completed {
		now := time.Now()
		completedAt = &now
	}

	return cycledomain.UpdateCycleInput{
		Label:       req.Label,
		Notes:       req.Notes,
		CompletedAt: completedAt,
	}
}
