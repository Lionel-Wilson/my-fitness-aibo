package dto

type CreateCycleRequest struct {
	Label string `json:"label"`
	Notes string `json:"notes"`
}

func (r CreateCycleRequest) Validate() error {
	return nil
}

type UpdateCycleRequest struct {
	Label     string `json:"label"`
	Notes     string `json:"notes"`
	Completed *bool  `json:"completed"`
}

func (r UpdateCycleRequest) Validate() error {
	return nil
}
