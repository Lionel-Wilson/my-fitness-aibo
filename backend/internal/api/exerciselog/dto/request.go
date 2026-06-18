package dto

type SetLogRequest struct {
	SetIndex  int      `json:"setIndex"`
	Side      string   `json:"side"`
	WeightKg  *float64 `json:"weightKg"`
	Reps      *float64 `json:"reps"`
	Rpe       *float64 `json:"rpe"`
	IsDropSet bool     `json:"isDropSet"`
}

type UpsertLogRequest struct {
	Note            string            `json:"note"`
	WorkingWeightKg *float64          `json:"workingWeightKg"`
	Sets            []SetLogRequest   `json:"sets"`
}

func (r UpsertLogRequest) Validate() error {
	return nil
}
