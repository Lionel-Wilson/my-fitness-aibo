package errors

import "errors"

// ErrNotFound is returned when a row does not exist or is not owned by the user.
var ErrNotFound = errors.New("not found")

// ErrConflict is returned on a unique-constraint violation (e.g. duplicate email).
var ErrConflict = errors.New("conflict")
