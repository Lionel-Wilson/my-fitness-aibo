// Package metrics holds strength-training calculations used across the API.
package metrics

// EstimatedOneRepMax returns the estimated one-rep max using the Epley formula:
//
//	1RM = weight * (1 + reps/30)
//
// A single rep returns the weight itself. Non-positive weight or reps yields 0.
func EstimatedOneRepMax(weight float64, reps int) float64 {
	if weight <= 0 || reps <= 0 {
		return 0
	}
	return weight * (1 + float64(reps)/30.0)
}
