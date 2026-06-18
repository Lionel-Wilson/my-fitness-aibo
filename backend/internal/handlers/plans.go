package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/lionel/my-fitness-aibo/backend/internal/store"
)

type planRequest struct {
	Name        string  `json:"name"`
	Quality     string  `json:"quality"`
	Description string  `json:"description"`
	CycleLabel  string  `json:"cycleLabel"`
	PeriodStart *string `json:"periodStart"`
	PeriodEnd   *string `json:"periodEnd"`
	IsActive    *bool   `json:"isActive"`
}

func parseDate(s *string, w http.ResponseWriter) (*time.Time, bool) {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil, true
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(*s))
	if err != nil {
		writeError(w, http.StatusBadRequest, "dates must be in YYYY-MM-DD format")
		return nil, false
	}
	return &t, true
}

func (req planRequest) toParams(w http.ResponseWriter) (store.PlanParams, bool) {
	start, ok := parseDate(req.PeriodStart, w)
	if !ok {
		return store.PlanParams{}, false
	}
	end, ok := parseDate(req.PeriodEnd, w)
	if !ok {
		return store.PlanParams{}, false
	}
	cycleLabel := strings.TrimSpace(req.CycleLabel)
	if cycleLabel == "" {
		cycleLabel = "Cycle"
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	return store.PlanParams{
		Name:        strings.TrimSpace(req.Name),
		Quality:     req.Quality,
		Description: req.Description,
		CycleLabel:  cycleLabel,
		PeriodStart: start,
		PeriodEnd:   end,
		IsActive:    isActive,
	}, true
}

// ListPlans returns the user's plans.
func (h *Handler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.Store.ListPlans(r.Context(), userID(r))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, plans)
}

// CreatePlan creates a plan.
func (h *Handler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var req planRequest
	if !decode(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	params, ok := req.toParams(w)
	if !ok {
		return
	}
	plan, err := h.Store.CreatePlan(r.Context(), userID(r), params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, plan)
}

// GetPlan returns a single plan.
func (h *Handler) GetPlan(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "planID")
	if !ok {
		return
	}
	plan, err := h.Store.GetPlan(r.Context(), userID(r), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}

// UpdatePlan updates a plan.
func (h *Handler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "planID")
	if !ok {
		return
	}
	var req planRequest
	if !decode(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	params, ok := req.toParams(w)
	if !ok {
		return
	}
	plan, err := h.Store.UpdatePlan(r.Context(), userID(r), id, params)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}

// DeletePlan deletes a plan.
func (h *Handler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	id, ok := pathUUID(w, r, "planID")
	if !ok {
		return
	}
	if err := h.Store.DeletePlan(r.Context(), userID(r), id); err != nil {
		writeStoreError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
