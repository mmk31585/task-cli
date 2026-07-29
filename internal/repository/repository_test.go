package repository

import (
	"os"
	"sync"
	"testing"

	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/storage"
)

func newTestRepo(t *testing.T) *JSONTaskRepository {
	t.Helper()
	store, err := storage.NewJSONStorage(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return NewJSONTaskRepository(store)
}

func TestConcurrentAdds(t *testing.T) {
	r := newTestRepo(t)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				_, err := r.Add("concurrent task")
				if err != nil {
					t.Logf("Add error: %v", err)
				}
			}
		}()
	}
	wg.Wait()

	tasks, err := r.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 100 {
		t.Fatalf("got %d tasks, want 100", len(tasks))
	}

	ids := make(map[int64]bool)
	for _, task := range tasks {
		if ids[task.ID] {
			t.Fatalf("duplicate ID: %d", task.ID)
		}
		ids[task.ID] = true
	}
}

func TestRepository(t *testing.T) {
	t.Run("Add then GetAll returns the added task", func(t *testing.T) {
		r := newTestRepo(t)
		task, err := r.Add("buy milk")
		if err != nil {
			t.Fatal(err)
		}
		if task.ID != 1 {
			t.Errorf("ID = %d, want 1", task.ID)
		}
		if task.Description != "buy milk" {
			t.Errorf("Description = %q, want %q", task.Description, "buy milk")
		}
		if task.Status != domain.StatusTodo {
			t.Errorf("Status = %q, want %q", task.Status, domain.StatusTodo)
		}
		tasks, err := r.GetAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 1 {
			t.Fatalf("GetAll returned %d tasks, want 1", len(tasks))
		}
		if tasks[0].ID != 1 {
			t.Errorf("ID = %d, want 1", tasks[0].ID)
		}
		if tasks[0].Description != "buy milk" {
			t.Errorf("Description = %q, want %q", tasks[0].Description, "buy milk")
		}
	})

	t.Run("Add multiple tasks returns correct count", func(t *testing.T) {
		r := newTestRepo(t)
		for i := 0; i < 5; i++ {
			_, err := r.Add("task")
			if err != nil {
				t.Fatal(err)
			}
		}
		tasks, err := r.GetAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 5 {
			t.Fatalf("GetAll returned %d tasks, want 5", len(tasks))
		}
		for i, task := range tasks {
			if task.ID != int64(i+1) {
				t.Errorf("tasks[%d].ID = %d, want %d", i, task.ID, i+1)
			}
		}
	})

	t.Run("GetByID returns correct task by ID", func(t *testing.T) {
		r := newTestRepo(t)
		first, _ := r.Add("first")
		second, _ := r.Add("second")
		third, _ := r.Add("third")

		got, err := r.GetByID(second.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.ID != second.ID {
			t.Errorf("ID = %d, want %d", got.ID, second.ID)
		}
		if got.Description != "second" {
			t.Errorf("Description = %q, want %q", got.Description, "second")
		}
		if got.ID == first.ID || got.ID == third.ID {
			t.Error("GetByID returned wrong task")
		}
	})

	t.Run("GetByID with unknown ID returns ErrTaskNotFound", func(t *testing.T) {
		r := newTestRepo(t)
		r.Add("task")
		_, err := r.GetByID(999)
		if err != domain.ErrTaskNotFound {
			t.Errorf("error = %v, want %v", err, domain.ErrTaskNotFound)
		}
	})

	t.Run("Update modifies description", func(t *testing.T) {
		r := newTestRepo(t)
		task, _ := r.Add("old desc")
		task.Description = "new desc"
		task.Status = domain.StatusDone
		updated, err := r.Update(task)
		if err != nil {
			t.Fatal(err)
		}
		if updated.Description != "new desc" {
			t.Errorf("Description = %q, want %q", updated.Description, "new desc")
		}
		if updated.Status != domain.StatusDone {
			t.Errorf("Status = %q, want %q", updated.Status, domain.StatusDone)
		}
		got, _ := r.GetByID(task.ID)
		if got.Description != "new desc" {
			t.Errorf("after GetAll: Description = %q, want %q", got.Description, "new desc")
		}
	})

	t.Run("Update with unknown ID returns ErrTaskNotFound", func(t *testing.T) {
		r := newTestRepo(t)
		r.Add("task")
		_, err := r.Update(domain.Task{ID: 999, Description: "x"})
		if err != domain.ErrTaskNotFound {
			t.Errorf("error = %v, want %v", err, domain.ErrTaskNotFound)
		}
	})

	t.Run("Delete removes task", func(t *testing.T) {
		r := newTestRepo(t)
		first, _ := r.Add("first")
		second, _ := r.Add("second")

		err := r.Delete(first.ID)
		if err != nil {
			t.Fatal(err)
		}
		tasks, _ := r.GetAll()
		if len(tasks) != 1 {
			t.Fatalf("GetAll returned %d tasks, want 1", len(tasks))
		}
		if tasks[0].ID != second.ID {
			t.Errorf("remaining task ID = %d, want %d", tasks[0].ID, second.ID)
		}
		_, err = r.GetByID(first.ID)
		if err != domain.ErrTaskNotFound {
			t.Error("deleted task should not be found")
		}
	})

	t.Run("Delete with unknown ID returns ErrTaskNotFound", func(t *testing.T) {
		r := newTestRepo(t)
		r.Add("task")
		err := r.Delete(999)
		if err != domain.ErrTaskNotFound {
			t.Errorf("error = %v, want %v", err, domain.ErrTaskNotFound)
		}
	})

	t.Run("Add returns error when storage read fails", func(t *testing.T) {
		dir := t.TempDir()
		os.MkdirAll(dir+"/tasks.json", 0755)
		store, _ := storage.NewJSONStorage(dir)
		r := NewJSONTaskRepository(store)
		_, err := r.Add("x")
		if err == nil {
			t.Error("Add() expected error, got nil")
		}
	})

	t.Run("GetAll, GetByID, Update, Delete return error when storage read fails", func(t *testing.T) {
		dir := t.TempDir()
		os.MkdirAll(dir+"/tasks.json", 0755)
		store, _ := storage.NewJSONStorage(dir)
		r := NewJSONTaskRepository(store)

		_, err := r.GetAll()
		if err == nil {
			t.Error("GetAll() expected error, got nil")
		}
		_, err = r.GetByID(1)
		if err == nil {
			t.Error("GetByID() expected error, got nil")
		}
		_, err = r.Update(domain.Task{ID: 1})
		if err == nil {
			t.Error("Update() expected error, got nil")
		}
		err = r.Delete(1)
		if err == nil {
			t.Error("Delete() expected error, got nil")
		}
	})

	t.Run("Add returns error when storage write fails", func(t *testing.T) {
		dir := t.TempDir()
		store, err := storage.NewJSONStorage(dir)
		if err != nil {
			t.Fatal(err)
		}
		r := NewJSONTaskRepository(store)
		if _, err := r.Add("first"); err != nil {
			t.Fatal(err)
		}
		os.Chmod(dir, 0444)
		t.Cleanup(func() { os.Chmod(dir, 0755) })
		_, err = r.Add("second")
		if err == nil {
			t.Error("Add() expected error on read-only dir, got nil")
		}
	})

	t.Run("all operations preserve other tasks", func(t *testing.T) {
		r := newTestRepo(t)
		_, _ = r.Add("task1")
		t2, _ := r.Add("task2")
		t3, _ := r.Add("task3")

		got, err := r.GetByID(t2.ID)
		if err != nil {
			t.Fatal(err)
		}
		got.Description = "updated task2"
		got.Status = domain.StatusInProgress
		r.Update(got)
		r.Delete(t3.ID)

		tasks, _ := r.GetAll()
		if len(tasks) != 2 {
			t.Fatalf("GetAll returned %d tasks, want 2", len(tasks))
		}
		if tasks[0].Description != "task1" && tasks[1].Description != "task1" {
			t.Error("task1 should still exist")
		}
		if tasks[0].Description != "updated task2" && tasks[1].Description != "updated task2" {
			t.Error("task2 should have updated description")
		}
	})
}
