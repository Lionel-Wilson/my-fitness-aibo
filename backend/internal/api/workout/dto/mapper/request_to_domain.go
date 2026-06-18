package mapper

import (
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/api/workout/dto"
	workoutdomain "github.com/lionel/my-fitness-aibo/backend/internal/workout/domain"
)

func WorkoutRequestToDomain(req dto.WorkoutRequest) workoutdomain.WorkoutInput {
	return workoutdomain.WorkoutInput{
		Name:        strings.TrimSpace(req.Name),
		DayLabel:    req.DayLabel,
		OrderIndex:  req.OrderIndex,
		DurationMin: req.DurationMin,
		Notes:       req.Notes,
	}
}
