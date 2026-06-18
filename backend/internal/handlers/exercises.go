package handlers

import (
	"net/http"
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/store"
)

type exerciseRequest struct {
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

func (req exerciseRequest) toParams() store.ExerciseParams {
	return store.ExerciseParams{
		Name:         strings.TrimSpace(req.Name),
		OrderIndex:   req.OrderIndex,
		TargetSets:   req.TargetSets,
		RepLow:       req.RepLow,
		RepHigh:      req.RepHigh,
		RpeLow:       req.RpeLow,
		RpeHigh:      req.RpeHigh,
		RestSeconds:  req.RestSeconds,
		Instructions: req.Instructions,
		OrGroup:      req.OrGroup,
		IsOptional:   req.IsOptional,
		IsUnilateral: req.IsUnilateral,
	}
}

// ListExercises returns the exercises of a workout.
func (h *Handler) ListExercises(w http.ResponseWriter, r *http.Request) {
	workoutID, ok := pathUUID(w, r, "workoutID")
	if !ok {
		return
	}
	exercises, err := h.Store.ListExercises(r.Context(), userID(r), workoutID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, exercises)
}

// CreateExercise creates an exercise in a workout.
func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	workoutID, ok := pathUUID(w, r, "workoutID")
	if !ok {
		return
	}
	var req exerciseRequest
	if !decode(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	exercise, err := h.Store.CreateExercise(r.Context(), userID(r), workoutID, req.toParams())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, exercise)
}

// GetExercise returns a single exercise.
func (h *Handler) GetExercise(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "exerciseID")
	if !ok {
		return
	}
	exercise, err := h.Store.GetExercise(r.Context(), userID(r), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, exercise)
}

// UpdateExercise updates an exercise.
func (h *Handler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "exerciseID")
	if !ok {
		return
	}
	var req exerciseRequest
	if !decode(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	exercise, err := h.Store.UpdateExercise(r.Context(), userID(r), id, req.toParams())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, exercise)
}

// DeleteExercise deletes an exercise.
func (h *Handler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "exerciseID")
	if !ok {
		return
	}
	if err := h.Store.DeleteExercise(r.Context(), userID(r), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
