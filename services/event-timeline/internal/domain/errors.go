package domain

import "errors"

var (
	ErrValidation       = errors.New("validation error")
	ErrNotFound         = errors.New("not found")
	ErrSequenceConflict = errors.New("sequence conflict")
	ErrIdempotencyConflict = errors.New("idempotency conflict")
	ErrUnauthorized     = errors.New("unauthorized")
)
