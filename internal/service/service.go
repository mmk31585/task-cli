package service

import (
	"errors"
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
	return t.repo.Add(desc)
}
func (t *TaskService) ListTasks(status string) ([]domain.Task, error) {
	tasks, err := t.repo.GetAll()
	if err != nil {
		return nil, err
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
		return domain.Task{}, err
	}
	task.Description = desc
	task.UpdatedAt = time.Now()
	return t.repo.Update(task)
}
func (t *TaskService) DeleteTask(id int64) error {
	if err := validator.ValidateID(id); err != nil {
		return err
	}
	_, err := t.repo.GetByID(id)
	if err != nil {
		return err
	}
	return t.repo.Delete(id)
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
		return domain.Task{}, err
	}
	task.Status = status
	task.UpdatedAt = time.Now()
	return t.repo.Update(task)
}
