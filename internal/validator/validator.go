package validator

import (
	"errors"
	"strings"
)

var (
	ErrEmptyDescription   = errors.New("description cannot be empty")
	ErrDescriptionTooLong = errors.New("description must not exceed 500 characters")
	ErrIDIsInvalid        = errors.New("id is invalid")
)

func ValidateDescription(desc string) error {
	if strings.TrimSpace(desc) == "" {
		return ErrEmptyDescription
	}
	if len(desc) > 500 {
		return ErrDescriptionTooLong
	}
	return nil
}

func ValidateID(id int64) error {
	if id < 0 {
		return ErrIDIsInvalid
	}
	return nil
}
