package validator

import (
	"strings"
	"testing"
)

func TestValidateDescription(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		err := ValidateDescription("")
		if err != ErrEmptyDescription {
			t.Errorf("got %v, want %v", err, ErrEmptyDescription)
		}
	})

	t.Run("whitespace-only returns error", func(t *testing.T) {
		err := ValidateDescription("   ")
		if err != ErrEmptyDescription {
			t.Errorf("got %v, want %v", err, ErrEmptyDescription)
		}
	})

	t.Run("too long returns error", func(t *testing.T) {
		long := strings.Repeat("a", 501)
		err := ValidateDescription(long)
		if err != ErrDescriptionTooLong {
			t.Errorf("got %v, want %v", err, ErrDescriptionTooLong)
		}
	})

	t.Run("valid description returns nil", func(t *testing.T) {
		err := ValidateDescription("Buy milk")
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("max length (500) returns nil", func(t *testing.T) {
		max := strings.Repeat("a", 500)
		err := ValidateDescription(max)
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})
}

func TestValidateID(t *testing.T) {
	t.Run("negative ID returns error", func(t *testing.T) {
		err := ValidateID(-5)
		if err != ErrIDIsInvalid {
			t.Errorf("got %v, want %v", err, ErrIDIsInvalid)
		}
	})

	t.Run("positive ID returns nil", func(t *testing.T) {
		err := ValidateID(3)
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("zero ID returns nil", func(t *testing.T) {
		err := ValidateID(0)
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})
}
