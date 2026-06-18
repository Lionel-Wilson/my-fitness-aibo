// Package middleware provides HTTP middleware, primarily JWT authentication.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/lionel/my-fitness-aibo/backend/internal/auth"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// Authenticator validates the bearer token on each request and injects the
// authenticated user ID into the request context.
func Authenticator(tm *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"missing or malformed authorization header"}`, http.StatusUnauthorized)
				return
			}

			userID, err := tm.Verify(parts[1])
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserID returns the authenticated user ID from the request context. The bool is
// false when the request was not authenticated.
func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
