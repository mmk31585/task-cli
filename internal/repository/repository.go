package repository

import (
	"fmt"

	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/storage"
)

type TaskRepository interface {
	Add(description string) (domain.Task, error)
	GetAll() ([]domain.Task, error)
	GetByID(id int64) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id int64) error
}

type JSONTaskRepository struct {
	store *storage.JSONStorage
}

func NewJSONTaskRepository(store *storage.JSONStorage) *JSONTaskRepository {
	return &JSONTaskRepository{store: store}
}

func (r *JSONTaskRepository) Add(description string) (domain.Task, error) {
	tasks, err := r.store.Read()
	if err != nil {
		return domain.Task{}, fmt.Errorf("read tasks: %w", err)
	}
	nextID := int64(1)
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
	task := domain.Task{
		ID:          nextID,
		Description: description,
		Status:      domain.StatusTodo,
	}
	tasks = append(tasks, task)
	if err := r.store.Write(tasks); err != nil {
		return domain.Task{}, fmt.Errorf("write tasks: %w", err)
	}
	return task, nil
}

func (r *JSONTaskRepository) GetAll() ([]domain.Task, error) {
	tasks, err := r.store.Read()
	if err != nil {
		return nil, fmt.Errorf("read tasks: %w", err)
	}
	return tasks, nil
}

func (r *JSONTaskRepository) GetByID(id int64) (domain.Task, error) {
	tasks, err := r.store.Read()
	if err != nil {
		return domain.Task{}, fmt.Errorf("read tasks: %w", err)
	}
	for _, t := range tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return domain.Task{}, domain.ErrTaskNotFound
}

func (r *JSONTaskRepository) Update(task domain.Task) (domain.Task, error) {
	tasks, err := r.store.Read()
	if err != nil {
		return domain.Task{}, fmt.Errorf("read tasks: %w", err)
	}
	for i, t := range tasks {
		if t.ID == task.ID {
			tasks[i] = task
			if err := r.store.Write(tasks); err != nil {
				return domain.Task{}, fmt.Errorf("write tasks: %w", err)
			}
			return task, nil
		}
	}
	return domain.Task{}, domain.ErrTaskNotFound
}

func (r *JSONTaskRepository) Delete(id int64) error {
	tasks, err := r.store.Read()
	if err != nil {
		return fmt.Errorf("read tasks: %w", err)
	}
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return r.store.Write(tasks)
		}
	}
	return domain.ErrTaskNotFound
}
