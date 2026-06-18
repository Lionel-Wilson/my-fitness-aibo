// Package handlers contains the HTTP handlers and router for the API.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/auth"
	"github.com/lionel/my-fitness-aibo/backend/internal/middleware"
	"github.com/lionel/my-fitness-aibo/backend/internal/store"
)

// Handler bundles the dependencies shared by all HTTP handlers.
type Handler struct {
	Store *store.Store
	Tokens *auth.TokenManager
}

// New returns a Handler.
func New(s *store.Store, tm *auth.TokenManager) *Handler {
	return &Handler{Store: s, Tokens: tm}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// writeStoreError maps store sentinel errors to HTTP status codes.
func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, store.ErrConflict):
		writeError(w, http.StatusConflict, "already exists")
	default:
		log.Printf("store error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

// decode parses the JSON request body into dst, returning false (and writing a
// 400) on failure.
func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

// userID returns the authenticated user ID; it always succeeds on routes behind
// the Authenticator middleware.
func userID(r *http.Request) uuid.UUID {
	id, _ := middleware.UserID(r.Context())
	return id
}

// pathUUID parses a UUID URL parameter, writing a 400 and returning false on failure.
func pathUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid "+key)
		return uuid.Nil, false
	}
	return id, true
}
