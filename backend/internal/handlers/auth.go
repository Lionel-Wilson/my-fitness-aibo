package handlers

import (
	"net/http"
	"strings"

	"github.com/lionel/my-fitness-aibo/backend/internal/auth"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

// Signup creates a new account and returns an access token.
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var c credentials
	if !decode(w, r, &c) {
		return
	}
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))
	if !strings.Contains(c.Email, "@") {
		writeError(w, http.StatusBadRequest, "a valid email is required")
		return
	}
	if len(c.Password) < 8 {
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hash, err := auth.HashPassword(c.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user, err := h.Store.CreateUser(r.Context(), c.Email, hash)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	token, err := h.Tokens.Issue(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not issue token")
		return
	}
	writeJSON(w, http.StatusCreated, authResponse{Token: token, User: user})
}

// Login validates credentials and returns an access token.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var c credentials
	if !decode(w, r, &c) {
		return
	}
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))

	user, err := h.Store.GetUserByEmail(r.Context(), c.Email)
	if err != nil || !auth.CheckPassword(user.PasswordHash, c.Password) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := h.Tokens.Issue(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not issue token")
		return
	}
	writeJSON(w, http.StatusOK, authResponse{Token: token, User: user})
}

// Logout is a no-op for stateless JWTs; the client discards the token.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// Me returns the authenticated user.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, err := h.Store.GetUserByID(r.Context(), userID(r))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}
