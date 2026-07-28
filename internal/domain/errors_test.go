package domain

import (
	"errors"
	"fmt"
	"testing"
)

var errs = []error{ErrInvalidStatus, ErrTaskNotFound, ErrDescriptionTooLong, ErrEmptyDescription}

func TestErrorsValue(t *testing.T) {
	for _, err := range errs {
		if err == nil {
			t.Errorf("expected non-nil error, got nil")
		}
	}
}

func TestErrorsComparable(t *testing.T) {
	for _, sentinel := range errs {

		if !errors.Is(sentinel, sentinel) {
			t.Errorf("errors.Is(%v, %v) = false, want true", sentinel, sentinel)
		}

		wrapped := fmt.Errorf("wrapped: %w", sentinel)
		if !errors.Is(wrapped, sentinel) {
			t.Errorf("errors.Is(wrapped, %v) = false, want true", sentinel)
		}

		other := fmt.Errorf("some other error")
		if errors.Is(other, sentinel) {
			t.Errorf("errors.Is(other, %v) = true, want false", sentinel)
		}
	}
}
