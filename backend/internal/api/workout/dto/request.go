package dto

type WorkoutRequest struct {
	Name        string `json:"name"`
	DayLabel    string `json:"dayLabel"`
	OrderIndex  int    `json:"orderIndex"`
	DurationMin *int   `json:"durationMin"`
	Notes       string `json:"notes"`
}

func (r WorkoutRequest) Validate() error {
	return nil
}
