package handlers

import "net/http"

// ExerciseProgress returns per-cycle aggregates for charting an exercise's
// strength and volume over time.
func (h *Handler) ExerciseProgress(w http.ResponseWriter, r *http.Request) {
	exerciseID, ok := pathUUID(w, r, "exerciseID")
	if !ok {
		return
	}
	if _, err := h.Store.GetExercise(r.Context(), userID(r), exerciseID); err != nil {
		writeStoreError(w, err)
		return
	}
	points, err := h.Store.ExerciseProgress(r.Context(), userID(r), exerciseID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, points)
}
