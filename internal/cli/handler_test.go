package cli

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/service"
)

type mockRepo struct {
	tasks []domain.Task
}

func newMockRepo() *mockRepo {
	return &mockRepo{}
}

func nextID(tasks []domain.Task) int64 {
	var maxID int64
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

func (m *mockRepo) Add(desc string) (domain.Task, error) {
	task := domain.Task{
		ID:          nextID(m.tasks),
		Description: desc,
		Status:      domain.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.tasks = append(m.tasks, task)
	return task, nil
}

func (m *mockRepo) GetAll() ([]domain.Task, error) {
	return m.tasks, nil
}

func (m *mockRepo) GetByID(id int64) (domain.Task, error) {
	for _, t := range m.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return domain.Task{}, domain.ErrTaskNotFound
}

func (m *mockRepo) Update(task domain.Task) (domain.Task, error) {
	for i, t := range m.tasks {
		if t.ID == task.ID {
			m.tasks[i] = task
			return task, nil
		}
	}
	return domain.Task{}, domain.ErrTaskNotFound
}

func (m *mockRepo) Delete(id int64) error {
	for i, t := range m.tasks {
		if t.ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return domain.ErrTaskNotFound
}

func newTestHandler(t *testing.T) (*Handler, *mockRepo) {
	t.Helper()
	repo := newMockRepo()
	svc := service.NewTaskService(repo)
	return NewHandler(svc), repo
}

func runHandler(t *testing.T, fn func()) (stdout, stderr string, exitCode int) {
	t.Helper()

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	oldExit := osExit

	os.Stdout = wOut
	os.Stderr = wErr

	var recordedCode int
	osExit = func(code int) {
		recordedCode = code
		panic("os.Exit")
	}

	func() {
		defer func() {
			wOut.Close()
			wErr.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr
			osExit = oldExit
			recover()
		}()
		fn()
	}()

	out, _ := io.ReadAll(rOut)
	err, _ := io.ReadAll(rErr)

	return string(out), string(err), recordedCode
}

func TestHandleAdd(t *testing.T) {
	t.Run("correct args outputs JSON to stdout", func(t *testing.T) {
		h, _ := newTestHandler(t)
		stdout, stderr, code := runHandler(t, func() {
			h.HandleAdd([]string{"buy milk"})
		})
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if stderr != "" {
			t.Errorf("expected no stderr, got %q", stderr)
		}
		if !strings.Contains(stdout, `"description": "buy milk"`) {
			t.Errorf("stdout missing description: %s", stdout)
		}
		if !strings.Contains(stdout, `"status": "todo"`) {
			t.Errorf("stdout missing status: %s", stdout)
		}
	})

	t.Run("missing args writes error to stderr, exits 1", func(t *testing.T) {
		h, _ := newTestHandler(t)
		_, stderr, code := runHandler(t, func() {
			h.HandleAdd([]string{})
		})
		if code != 1 {
			t.Errorf("expected exit code 1, got %d", code)
		}
		if stderr == "" {
			t.Error("expected stderr output")
		}
	})
}

func TestHandleList(t *testing.T) {
	t.Run("outputs JSON array", func(t *testing.T) {
		h, repo := newTestHandler(t)
		repo.Add("task one")
		repo.Add("task two")

		stdout, stderr, code := runHandler(t, func() {
			h.HandleList([]string{}, false)
		})
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if stderr != "" {
			t.Errorf("expected no stderr, got %q", stderr)
		}
		if !strings.HasPrefix(strings.TrimSpace(stdout), "[") {
			t.Errorf("expected JSON array, got %s", stdout)
		}
	})

	t.Run("--table outputs formatted table", func(t *testing.T) {
		h, repo := newTestHandler(t)
		repo.Add("task one")

		stdout, stderr, code := runHandler(t, func() {
			h.HandleList([]string{}, true)
		})
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if stderr != "" {
			t.Errorf("expected no stderr, got %q", stderr)
		}
		if !strings.Contains(stdout, "ID") {
			t.Errorf("expected table header, got %s", stdout)
		}
		if !strings.Contains(stdout, "task one") {
			t.Errorf("expected task description, got %s", stdout)
		}
	})

	t.Run("filters by status", func(t *testing.T) {
		h, repo := newTestHandler(t)
		repo.Add("todo task")
		repo.Add("done task")
		h.svc.MarkTask(2, domain.StatusDone)

		stdout, stderr, code := runHandler(t, func() {
			h.HandleList([]string{"done"}, false)
		})
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if stderr != "" {
			t.Errorf("expected no stderr, got %q", stderr)
		}
		if strings.Contains(stdout, "todo task") {
			t.Errorf("expected only done tasks, got todo in output")
		}
		if !strings.Contains(stdout, "done task") {
			t.Errorf("expected done task in output")
		}
	})
}

func TestHandleUpdate(t *testing.T) {
	t.Run("correct args outputs updated task", func(t *testing.T) {
		h, repo := newTestHandler(t)
		repo.Add("old desc")

		stdout, stderr, code := runHandler(t, func() {
			h.HandleUpdate([]string{"1", "new desc"})
		})
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if stderr != "" {
			t.Errorf("expected no stderr, got %q", stderr)
		}
		if !strings.Contains(stdout, `"description": "new desc"`) {
			t.Errorf("expected updated description, got %s", stdout)
		}
	})

	t.Run("unknown ID returns error", func(t *testing.T) {
		h, _ := newTestHandler(t)
		_, stderr, code := runHandler(t, func() {
			h.HandleUpdate([]string{"999", "desc"})
		})
		if code != 1 {
			t.Errorf("expected exit code 1, got %d", code)
		}
		if stderr == "" {
			t.Error("expected stderr output")
		}
	})
}

func TestHandleDelete(t *testing.T) {
	t.Run("outputs success JSON", func(t *testing.T) {
		h, repo := newTestHandler(t)
		repo.Add("to delete")

		stdout, stderr, code := runHandler(t, func() {
			h.HandleDelete([]string{"1"})
		})
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if stderr != "" {
			t.Errorf("expected no stderr, got %q", stderr)
		}
		if !strings.Contains(stdout, `"deleted": 1`) {
			t.Errorf("expected delete JSON, got %s", stdout)
		}
	})

	t.Run("unknown ID returns error", func(t *testing.T) {
		h, _ := newTestHandler(t)
		_, stderr, code := runHandler(t, func() {
			h.HandleDelete([]string{"999"})
		})
		if code != 1 {
			t.Errorf("expected exit code 1, got %d", code)
		}
		if stderr == "" {
			t.Error("expected stderr output")
		}
	})
}

func TestHandleMark(t *testing.T) {
	t.Run("outputs updated task with new status", func(t *testing.T) {
		h, repo := newTestHandler(t)
		repo.Add("task")

		stdout, stderr, code := runHandler(t, func() {
			h.HandleMark([]string{"1"}, domain.StatusDone)
		})
		if code != 0 {
			t.Errorf("expected exit code 0, got %d", code)
		}
		if stderr != "" {
			t.Errorf("expected no stderr, got %q", stderr)
		}
		if !strings.Contains(stdout, `"status": "done"`) {
			t.Errorf("expected done status, got %s", stdout)
		}
	})

	t.Run("unknown ID returns error", func(t *testing.T) {
		h, _ := newTestHandler(t)
		_, stderr, code := runHandler(t, func() {
			h.HandleMark([]string{"999"}, domain.StatusDone)
		})
		if code != 1 {
			t.Errorf("expected exit code 1, got %d", code)
		}
		if stderr == "" {
			t.Error("expected stderr output")
		}
	})
}
