package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorsValue(t *testing.T) {
	if ErrTaskNotFound == nil {
		t.Error("expected non-nil error, got nil")
	}
}

func TestErrorsComparable(t *testing.T) {
	if !errors.Is(ErrTaskNotFound, ErrTaskNotFound) {
		t.Errorf("errors.Is(ErrTaskNotFound, ErrTaskNotFound) = false, want true")
	}

	wrapped := fmt.Errorf("wrapped: %w", ErrTaskNotFound)
	if !errors.Is(wrapped, ErrTaskNotFound) {
		t.Errorf("errors.Is(wrapped, ErrTaskNotFound) = false, want true")
	}

	other := fmt.Errorf("some other error")
	if errors.Is(other, ErrTaskNotFound) {
		t.Errorf("errors.Is(other, ErrTaskNotFound) = true, want false")
	}
}
