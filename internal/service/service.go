package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/repository"
	"github.com/mmk31585/task-cli/internal/validator"
)

var (
	ErrStatusInvalid = errors.New("status is not valid")
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repository,
	}
}
func (t *TaskService) AddTask(desc string) (domain.Task, error) {
	if err := validator.ValidateDescription(desc); err != nil {
		return domain.Task{}, err
	}
	task, err := t.repo.Add(desc)
	if err != nil {
		return domain.Task{}, fmt.Errorf("add task: %w", err)
	}
	return task, nil
}
func (t *TaskService) ListTasks(status string) ([]domain.Task, error) {
	tasks, err := t.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	if status == "" {
		return tasks, nil
	}
	if !domain.IsValidStatus(domain.Status(status)) {
		return nil, ErrStatusInvalid
	}
	return filterTask(tasks, domain.Status(status)), nil
}
func filterTask(tasks []domain.Task, status domain.Status) []domain.Task {
	var out []domain.Task
	for _, t := range tasks {
		if t.Status == status {
			out = append(out, t)
		}
	}
	return out
}

func (t *TaskService) UpdateTask(id int64, desc string) (domain.Task, error) {
	if err := validator.ValidateID(id); err != nil {
		return domain.Task{}, err
	}
	if err := validator.ValidateDescription(desc); err != nil {
		return domain.Task{}, err
	}
	task, err := t.repo.GetByID(id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}
	task.Description = desc
	task.UpdatedAt = time.Now()
	task, err = t.repo.Update(task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}
	return task, nil
}
func (t *TaskService) DeleteTask(id int64) error {
	if err := validator.ValidateID(id); err != nil {
		return err
	}
	_, err := t.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}
	if err := t.repo.Delete(id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}
func (t *TaskService) MarkTask(id int64, status domain.Status) (domain.Task, error) {
	if err := validator.ValidateID(id); err != nil {
		return domain.Task{}, err
	}
	if !domain.IsValidStatus(status) {
		return domain.Task{}, ErrStatusInvalid
	}
	task, err := t.repo.GetByID(id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}
	task.Status = status
	task.UpdatedAt = time.Now()
	task, err = t.repo.Update(task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}
	return task, nil
}
