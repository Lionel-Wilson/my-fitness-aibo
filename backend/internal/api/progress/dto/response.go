package dto

import (
	"time"

	"github.com/google/uuid"

	progressdomain "github.com/lionel/my-fitness-aibo/backend/internal/progress/domain"
)

type ProgressPointResponse struct {
	CycleID     uuid.UUID `json:"cycleId"`
	CycleNumber int       `json:"cycleNumber"`
	Label       string    `json:"label"`
	Side        string    `json:"side"`
	TopWeightKg float64   `json:"topWeightKg"`
	VolumeKg    float64   `json:"volumeKg"`
	BestE1RM    float64   `json:"bestE1rm"`
	TotalReps   float64   `json:"totalReps"`
	StartedAt   time.Time `json:"startedAt"`
}

func ProgressPointFromDomain(p progressdomain.ProgressPoint) ProgressPointResponse {
	return ProgressPointResponse{
		CycleID:     p.CycleID,
		CycleNumber: p.CycleNumber,
		Label:       p.Label,
		Side:        p.Side,
		TopWeightKg: p.TopWeightKg,
		VolumeKg:    p.VolumeKg,
		BestE1RM:    p.BestE1RM,
		TotalReps:   p.TotalReps,
		StartedAt:   p.StartedAt,
	}
}

func ProgressPointsFromDomain(points []progressdomain.ProgressPoint) []ProgressPointResponse {
	out := make([]ProgressPointResponse, len(points))
	for i, p := range points {
		out[i] = ProgressPointFromDomain(p)
	}

	return out
}
