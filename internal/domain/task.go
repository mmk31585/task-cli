package domain

import "time"

type Status string

type Task struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

func IsValidStatus(status Status) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false

	}
}
