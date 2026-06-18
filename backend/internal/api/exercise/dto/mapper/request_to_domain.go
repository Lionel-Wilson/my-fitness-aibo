package mapper

import (
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/exercise/dto"
	exercisedomain "github.com/lionel/my-fitness-aibo/backend/internal/exercise/domain"
)

func ExerciseRequestToDomain(req dto.ExerciseRequest) exercisedomain.ExerciseInput {
	return exercisedomain.ExerciseInput{
		Name:         strings.TrimSpace(req.Name),
		OrderIndex:   req.OrderIndex,
		TargetSets:   req.TargetSets,
		RepLow:       req.RepLow,
		RepHigh:      req.RepHigh,
		RpeLow:       req.RpeLow,
		RpeHigh:      req.RpeHigh,
		RestSeconds:  req.RestSeconds,
		Instructions: req.Instructions,
		OrGroup:      req.OrGroup,
		IsOptional:   req.IsOptional,
		IsUnilateral: req.IsUnilateral,
	}
}
