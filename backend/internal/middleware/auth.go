// Package middleware provides HTTP middleware, primarily JWT authentication.
package middleware

import (
	"net/http"
	"strings"

	commonauth "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/auth"
	commoncontext "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/context"
	commonMappers "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/mappers"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
)

// Authenticator validates the bearer token on each request and injects the
// authenticated user ID into the request context.
func Authenticator(tm *commonauth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				render.Json(w, http.StatusUnauthorized, commonMappers.ToSimpleErrorResponse("missing or malformed authorization header"))
				return
			}

			userID, err := tm.Verify(parts[1])
			if err != nil {
				render.Json(w, http.StatusUnauthorized, commonMappers.ToSimpleErrorResponse("invalid or expired token"))
				return
			}

			ctx := commoncontext.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
