package handlers

import (
	"net/http"

	"github.com/lionel/my-fitness-aibo/backend/internal/store"
)

type upsertLogRequest struct {
	Note            string               `json:"note"`
	WorkingWeightKg *float64             `json:"workingWeightKg"`
	Sets            []store.SetLogParams `json:"sets"`
}

// ListExerciseLogs returns all cycle logs (with sets) for an exercise. Powers
// the cycle history grid.
func (h *Handler) ListExerciseLogs(w http.ResponseWriter, r *http.Request) {
	exerciseID, ok := pathUUID(w, r, "exerciseID")
	if !ok {
		return
	}
	// Confirm ownership so a non-owned exercise returns 404 rather than [].
	if _, err := h.Store.GetExercise(r.Context(), userID(r), exerciseID); err != nil {
		writeStoreError(w, err)
		return
	}
	logs, err := h.Store.ListExerciseLogs(r.Context(), userID(r), exerciseID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, logs)
}

// UpsertExerciseLog creates or replaces the log for an exercise in a cycle,
// including all its sets.
func (h *Handler) UpsertExerciseLog(w http.ResponseWriter, r *http.Request) {
	exerciseID, ok := pathUUID(w, r, "exerciseID")
	if !ok {
		return
	}
	cycleID, ok := pathUUID(w, r, "cycleID")
	if !ok {
		return
	}
	var req upsertLogRequest
	if !decode(w, r, &req) {
		return
	}
	log, err := h.Store.UpsertExerciseLog(r.Context(), userID(r), exerciseID, cycleID, req.Note, req.WorkingWeightKg, req.Sets)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, log)
}
