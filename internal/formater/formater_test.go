package formater

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mmk31585/task-cli/internal/domain"
)

func TestFormatTaskJSON(t *testing.T) {
	task := domain.Task{
		ID:          1,
		Description: "buy milk",
		Status:      domain.StatusTodo,
		CreatedAt:   time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC),
	}

	got, err := FormatTaskJSON(task)
	if err != nil {
		t.Fatal(err)
	}

	if !json.Valid([]byte(got)) {
		t.Error("output is not valid JSON")
	}

	if !strings.Contains(got, `"id": 1`) {
		t.Error("missing id")
	}
	if !strings.Contains(got, `"buy milk"`) {
		t.Error("missing description")
	}
	if !strings.Contains(got, `"todo"`) {
		t.Error("missing status")
	}
	if !strings.Contains(got, `"created_at"`) {
		t.Error("missing created_at")
	}
	if !strings.Contains(got, `"updated_at"`) {
		t.Error("missing updated_at")
	}
}

func TestFormatTasksJSON(t *testing.T) {
	tasks := []domain.Task{
		{ID: 1, Description: "first", Status: domain.StatusTodo, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Description: "second", Status: domain.StatusDone, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	got, err := FormatTasksJSON(tasks)
	if err != nil {
		t.Fatal(err)
	}

	if !json.Valid([]byte(got)) {
		t.Error("output is not valid JSON")
	}

	if !strings.HasPrefix(strings.TrimSpace(got), "[") {
		t.Error("expected JSON array to start with [")
	}
	if !strings.HasSuffix(strings.TrimSpace(got), "]") {
		t.Error("expected JSON array to end with ]")
	}

	if !strings.Contains(got, "first") {
		t.Error("missing first task")
	}
	if !strings.Contains(got, "second") {
		t.Error("missing second task")
	}
}

func TestFormatTasksTable(t *testing.T) {
	t.Run("has headers and data", func(t *testing.T) {
		tasks := []domain.Task{
			{ID: 1, Description: "buy milk", Status: domain.StatusTodo, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, Description: "write docs", Status: domain.StatusInProgress, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		got, err := FormatTasksTable(tasks)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(got, "ID") {
			t.Error("missing ID header")
		}
		if !strings.Contains(got, "Description") {
			t.Error("missing Description header")
		}
		if !strings.Contains(got, "Status") {
			t.Error("missing Status header")
		}
		if !strings.Contains(got, "Created At") {
			t.Error("missing Created At header")
		}
		if !strings.Contains(got, "Updated At") {
			t.Error("missing Updated At header")
		}

		if !strings.Contains(got, "buy milk") {
			t.Error("missing first task description")
		}
		if !strings.Contains(got, "write docs") {
			t.Error("missing second task description")
		}
	})

	t.Run("empty returns error and message", func(t *testing.T) {
		got, err := FormatTasksTable([]domain.Task{})
		if err != ErrNotFoundTask {
			t.Errorf("expected ErrNotFoundTask, got %v", err)
		}
		if got != "No tasks found\n" {
			t.Errorf("expected 'No tasks found\\n', got %q", got)
		}
	})

	t.Run("single task", func(t *testing.T) {
		tasks := []domain.Task{
			{ID: 1, Description: "only one", Status: domain.StatusDone, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		got, err := FormatTasksTable(tasks)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(got, "only one") {
			t.Error("missing task description")
		}
	})

	t.Run("many tasks", func(t *testing.T) {
		var tasks []domain.Task
		for i := int64(1); i <= 100; i++ {
			tasks = append(tasks, domain.Task{
				ID: i, Description: "task", Status: domain.StatusTodo,
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			})
		}

		got, err := FormatTasksTable(tasks)
		if err != nil {
			t.Fatal(err)
		}

		for i := int64(1); i <= 100; i++ {
			if !strings.Contains(got, fmt.Sprintf("%d", i)) {
				t.Errorf("missing ID %d in output", i)
			}
		}
	})
}
