package errors

import "errors"

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("not found")
)
