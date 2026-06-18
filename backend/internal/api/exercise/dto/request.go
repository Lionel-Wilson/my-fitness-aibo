package dto

type ExerciseRequest struct {
	Name         string   `json:"name"`
	OrderIndex   int      `json:"orderIndex"`
	TargetSets   *int     `json:"targetSets"`
	RepLow       *int     `json:"repLow"`
	RepHigh      *int     `json:"repHigh"`
	RpeLow       *float64 `json:"rpeLow"`
	RpeHigh      *float64 `json:"rpeHigh"`
	RestSeconds  *int     `json:"restSeconds"`
	Instructions string   `json:"instructions"`
	OrGroup      string   `json:"orGroup"`
	IsOptional   bool     `json:"isOptional"`
	IsUnilateral bool     `json:"isUnilateral"`
}

func (r ExerciseRequest) Validate() error {
	return nil
}
