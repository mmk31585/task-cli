package domain

import (
	"reflect"
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{name: "todo returns true", status: StatusTodo, want: true},
		{name: "in-progress returns true", status: StatusInProgress, want: true},
		{name: "done returns true", status: StatusDone, want: true},
		{name: "invalid returns false", status: "invalid", want: false},
		{name: "empty returns false", status: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestTaskJsonStruct(t *testing.T) {
	if reflect.TypeOf(Task{}).Field(0).Tag.Get("json") != "id" {
		t.Errorf("we expect for first tag id but got %v", reflect.TypeOf(Task{}).Field(0).Tag.Get("json"))
	}
	if reflect.TypeOf(Task{}).Field(1).Tag.Get("json") != "description" {
		t.Errorf("we expect for first tag description but got %v", reflect.TypeOf(Task{}).Field(0).Tag.Get("json"))
	}
	if reflect.TypeOf(Task{}).Field(2).Tag.Get("json") != "status" {
		t.Errorf("we expect for first tag status but got %v", reflect.TypeOf(Task{}).Field(0).Tag.Get("json"))
	}
	if reflect.TypeOf(Task{}).Field(3).Tag.Get("json") != "created_at" {
		t.Errorf("we expect for first tag created_at but got %v", reflect.TypeOf(Task{}).Field(0).Tag.Get("json"))
	}
	if reflect.TypeOf(Task{}).Field(4).Tag.Get("json") != "updated_at" {
		t.Errorf("we expect for first tag updated_at but got %v", reflect.TypeOf(Task{}).Field(0).Tag.Get("json"))
	}
}
