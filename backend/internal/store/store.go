// Package store is the data-access layer: typed models and hand-written pgx
// queries. Every read/write is scoped to the owning user.
package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a row does not exist or is not owned by the user.
var ErrNotFound = errors.New("not found")

// ErrConflict is returned on a unique-constraint violation (e.g. duplicate email).
var ErrConflict = errors.New("conflict")

// Store wraps the connection pool and exposes query methods.
type Store struct {
	pool *pgxpool.Pool
}

// New returns a Store backed by the given pool.
func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}
