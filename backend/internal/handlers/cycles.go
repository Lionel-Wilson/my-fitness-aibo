package handlers

import (
	"net/http"
	"time"
)

type createCycleRequest struct {
	Label string `json:"label"`
	Notes string `json:"notes"`
}

type updateCycleRequest struct {
	Label     string `json:"label"`
	Notes     string `json:"notes"`
	Completed *bool  `json:"completed"`
}

// ListCycles returns the cycles of a plan.
func (h *Handler) ListCycles(w http.ResponseWriter, r *http.Request) {
	planID, ok := pathUUID(w, r, "planID")
	if !ok {
		return
	}
	cycles, err := h.Store.ListCycles(r.Context(), userID(r), planID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cycles)
}

// CreateCycle starts the next cycle for a plan.
func (h *Handler) CreateCycle(w http.ResponseWriter, r *http.Request) {
	planID, ok := pathUUID(w, r, "planID")
	if !ok {
		return
	}
	var req createCycleRequest
	if !decode(w, r, &req) {
		return
	}
	cycle, err := h.Store.CreateCycle(r.Context(), userID(r), planID, req.Label, req.Notes)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, cycle)
}

// UpdateCycle updates a cycle's label, notes and completion state.
func (h *Handler) UpdateCycle(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "cycleID")
	if !ok {
		return
	}
	var req updateCycleRequest
	if !decode(w, r, &req) {
		return
	}
	var completedAt *time.Time
	if req.Completed != nil && *req.Completed {
		now := time.Now()
		completedAt = &now
	}
	cycle, err := h.Store.UpdateCycle(r.Context(), userID(r), id, req.Label, req.Notes, completedAt)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cycle)
}

// DeleteCycle deletes a cycle.
func (h *Handler) DeleteCycle(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "cycleID")
	if !ok {
		return
	}
	if err := h.Store.DeleteCycle(r.Context(), userID(r), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
