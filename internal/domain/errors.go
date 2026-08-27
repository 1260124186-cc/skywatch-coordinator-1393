package domain

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrConflict       = errors.New("state conflict")
	ErrReleaseBlocked = errors.New("campaign cannot be released yet")
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e FieldError) Error() string { return e.Field + ": " + e.Message }
