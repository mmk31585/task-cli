package domain

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrInvalidStatus      = errors.New("invalid status")
	ErrEmptyDescription   = errors.New("the description is empty")
	ErrDescriptionTooLong = errors.New("the description is too long")
)
