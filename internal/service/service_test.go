package service

import (
	"errors"
	"testing"
	"time"

	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/validator"
)

type mockRepo struct {
	tasks  map[int64]domain.Task
	nextID int64
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		tasks:  make(map[int64]domain.Task),
		nextID: 1,
	}
}

func (m *mockRepo) Add(desc string) (domain.Task, error) {
	task := domain.Task{
		ID:          m.nextID,
		Description: desc,
		Status:      domain.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.tasks[m.nextID] = task
	m.nextID++
	return task, nil
}

func (m *mockRepo) GetAll() ([]domain.Task, error) {
	out := make([]domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		out = append(out, t)
	}
	return out, nil
}

func (m *mockRepo) GetByID(id int64) (domain.Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrTaskNotFound
	}
	return t, nil
}

func (m *mockRepo) Update(task domain.Task) (domain.Task, error) {
	if _, ok := m.tasks[task.ID]; !ok {
		return domain.Task{}, domain.ErrTaskNotFound
	}
	m.tasks[task.ID] = task
	return task, nil
}

func (m *mockRepo) Delete(id int64) error {
	if _, ok := m.tasks[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(m.tasks, id)
	return nil
}

func TestAddTask(t *testing.T) {
	t.Run("valid description returns task with StatusTodo", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		task, err := svc.AddTask("buy milk")
		if err != nil {
			t.Fatal(err)
		}
		if task.Description != "buy milk" {
			t.Errorf("expected 'buy milk', got %q", task.Description)
		}
		if task.Status != domain.StatusTodo {
			t.Errorf("expected StatusTodo, got %v", task.Status)
		}
		if task.ID != 1 {
			t.Errorf("expected ID 1, got %d", task.ID)
		}
		if task.CreatedAt.IsZero() {
			t.Error("expected non-zero CreatedAt")
		}
		if task.UpdatedAt.IsZero() {
			t.Error("expected non-zero UpdatedAt")
		}
	})

	t.Run("empty description returns error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		_, err := svc.AddTask("")
		if !errors.Is(err, validator.ErrEmptyDescription) {
			t.Errorf("expected ErrEmptyDescription, got %v", err)
		}
	})

	t.Run("long description returns error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		desc := string(make([]byte, 501))
		_, err := svc.AddTask(desc)
		if !errors.Is(err, validator.ErrDescriptionTooLong) {
			t.Errorf("expected ErrDescriptionTooLong, got %v", err)
		}
	})
}

func TestListTasks(t *testing.T) {
	t.Run("empty status returns all tasks", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		if _, err := svc.AddTask("one"); err != nil {
			t.Fatal(err)
		}
		if _, err := svc.AddTask("two"); err != nil {
			t.Fatal(err)
		}

		tasks, err := svc.ListTasks("")
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 2 {
			t.Errorf("expected 2 tasks, got %d", len(tasks))
		}
	})

	t.Run("status filter returns matching tasks", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		if _, err := svc.AddTask("todo task"); err != nil {
			t.Fatal(err)
		}
		doneTask, err := svc.AddTask("done task")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.MarkTask(doneTask.ID, domain.StatusDone); err != nil {
			t.Fatal(err)
		}

		tasks, err := svc.ListTasks("done")
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 1 {
			t.Errorf("expected 1 done task, got %d", len(tasks))
		}
		if tasks[0].Description != "done task" {
			t.Errorf("expected 'done task', got %q", tasks[0].Description)
		}
	})

	t.Run("invalid status returns error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		_, err := svc.ListTasks("invalid")
		if !errors.Is(err, ErrStatusInvalid) {
			t.Errorf("expected ErrStatusInvalid, got %v", err)
		}
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("valid input updates description and UpdatedAt", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		original, _ := svc.AddTask("old desc")
		time.Sleep(time.Millisecond)

		updated, err := svc.UpdateTask(original.ID, "new desc")
		if err != nil {
			t.Fatal(err)
		}
		if updated.Description != "new desc" {
			t.Errorf("expected 'new desc', got %q", updated.Description)
		}
		if updated.UpdatedAt.Equal(original.UpdatedAt) {
			t.Error("expected UpdatedAt to change")
		}
		if updated.CreatedAt != original.CreatedAt {
			t.Error("expected CreatedAt to stay the same")
		}
	})

	t.Run("unknown ID returns ErrTaskNotFound", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		_, err := svc.UpdateTask(999, "desc")
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("invalid ID returns validation error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		_, err := svc.UpdateTask(-1, "desc")
		if !errors.Is(err, validator.ErrIDIsInvalid) {
			t.Errorf("expected ErrIDIsInvalid, got %v", err)
		}
	})

	t.Run("empty description returns validation error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		task, _ := svc.AddTask("desc")
		_, err := svc.UpdateTask(task.ID, "")
		if !errors.Is(err, validator.ErrEmptyDescription) {
			t.Errorf("expected ErrEmptyDescription, got %v", err)
		}
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("removes existing task", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		task, _ := svc.AddTask("to delete")

		err := svc.DeleteTask(task.ID)
		if err != nil {
			t.Fatal(err)
		}

		_, err = svc.repo.GetByID(task.ID)
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Error("expected task to be deleted")
		}
	})

	t.Run("unknown ID returns ErrTaskNotFound", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		err := svc.DeleteTask(999)
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("invalid ID returns validation error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		err := svc.DeleteTask(-1)
		if !errors.Is(err, validator.ErrIDIsInvalid) {
			t.Errorf("expected ErrIDIsInvalid, got %v", err)
		}
	})
}

func TestMarkTask(t *testing.T) {
	t.Run("transitions status correctly", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		task, _ := svc.AddTask("task")

		task, err := svc.MarkTask(task.ID, domain.StatusInProgress)
		if err != nil {
			t.Fatal(err)
		}
		if task.Status != domain.StatusInProgress {
			t.Errorf("expected StatusInProgress, got %v", task.Status)
		}

		task, err = svc.MarkTask(task.ID, domain.StatusDone)
		if err != nil {
			t.Fatal(err)
		}
		if task.Status != domain.StatusDone {
			t.Errorf("expected StatusDone, got %v", task.Status)
		}
	})

	t.Run("unknown ID returns ErrTaskNotFound", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		_, err := svc.MarkTask(999, domain.StatusDone)
		if !errors.Is(err, domain.ErrTaskNotFound) {
			t.Errorf("expected ErrTaskNotFound, got %v", err)
		}
	})

	t.Run("invalid status returns error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		task, _ := svc.AddTask("task")
		_, err := svc.MarkTask(task.ID, domain.Status("bogus"))
		if !errors.Is(err, ErrStatusInvalid) {
			t.Errorf("expected ErrStatusInvalid, got %v", err)
		}
	})

	t.Run("invalid ID returns validation error", func(t *testing.T) {
		svc := NewTaskService(newMockRepo())
		_, err := svc.MarkTask(-1, domain.StatusDone)
		if !errors.Is(err, validator.ErrIDIsInvalid) {
			t.Errorf("expected ErrIDIsInvalid, got %v", err)
		}
	})
}
