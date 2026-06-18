package handlers

import (
	"net/http"
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/store"
)

type workoutRequest struct {
	Name        string `json:"name"`
	DayLabel    string `json:"dayLabel"`
	OrderIndex  int    `json:"orderIndex"`
	DurationMin *int   `json:"durationMin"`
	Notes       string `json:"notes"`
}

func (req workoutRequest) toParams() store.WorkoutParams {
	return store.WorkoutParams{
		Name:        strings.TrimSpace(req.Name),
		DayLabel:    req.DayLabel,
		OrderIndex:  req.OrderIndex,
		DurationMin: req.DurationMin,
		Notes:       req.Notes,
	}
}

// ListWorkouts returns the workouts of a plan.
func (h *Handler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
	planID, ok := pathUUID(w, r, "planID")
	if !ok {
		return
	}
	workouts, err := h.Store.ListWorkouts(r.Context(), userID(r), planID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, workouts)
}

// CreateWorkout creates a workout in a plan.
func (h *Handler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	planID, ok := pathUUID(w, r, "planID")
	if !ok {
		return
	}
	var req workoutRequest
	if !decode(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	workout, err := h.Store.CreateWorkout(r.Context(), userID(r), planID, req.toParams())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, workout)
}

// GetWorkout returns a single workout.
func (h *Handler) GetWorkout(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "workoutID")
	if !ok {
		return
	}
	workout, err := h.Store.GetWorkout(r.Context(), userID(r), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, workout)
}

// UpdateWorkout updates a workout.
func (h *Handler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "workoutID")
	if !ok {
		return
	}
	var req workoutRequest
	if !decode(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	workout, err := h.Store.UpdateWorkout(r.Context(), userID(r), id, req.toParams())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, workout)
}

// DeleteWorkout deletes a workout.
func (h *Handler) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "workoutID")
	if !ok {
		return
	}
	if err := h.Store.DeleteWorkout(r.Context(), userID(r), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
