package dto

type PlanRequest struct {
	Name        string  `json:"name"`
	Quality     string  `json:"quality"`
	Description string  `json:"description"`
	CycleLabel  string  `json:"cycleLabel"`
	PeriodStart *string `json:"periodStart"`
	PeriodEnd   *string `json:"periodEnd"`
	IsActive    *bool   `json:"isActive"`
}

func (r PlanRequest) Validate() error {
	return nil
}
