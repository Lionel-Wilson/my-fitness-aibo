package context

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// WithUserID returns a context carrying the authenticated user ID.
func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserID returns the authenticated user ID from the context.
func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
