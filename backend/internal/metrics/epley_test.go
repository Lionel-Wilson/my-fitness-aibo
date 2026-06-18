package metrics

import (
	"math"
	"testing"
)

func TestEstimatedOneRepMax(t *testing.T) {
	tests := []struct {
		name   string
		weight float64
		reps   int
		want   float64
	}{
		{"single rep returns weight", 100, 1, 100 * (1 + 1.0/30)},
		{"ten reps", 100, 10, 100 * (1 + 10.0/30)},
		{"zero reps", 100, 0, 0},
		{"zero weight", 0, 5, 0},
		{"negative weight", -50, 5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EstimatedOneRepMax(tt.weight, tt.reps)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("EstimatedOneRepMax(%v, %d) = %v, want %v", tt.weight, tt.reps, got, tt.want)
			}
		})
	}
}
