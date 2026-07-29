package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/formater"
	"github.com/mmk31585/task-cli/internal/service"
	"github.com/mmk31585/task-cli/internal/validator"
)

var osExit = os.Exit

type Handler struct {
	svc *service.TaskService
}

func NewHandler(service *service.TaskService) *Handler {
	return &Handler{
		svc: service,
	}
}

func printError(err error) {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		fmt.Fprintln(os.Stderr, "Error: task not found")
	case errors.Is(err, validator.ErrEmptyDescription):
		fmt.Fprintln(os.Stderr, "Error: description cannot be empty")
	case errors.Is(err, validator.ErrDescriptionTooLong):
		fmt.Fprintln(os.Stderr, "Error: description must not exceed 500 characters")
	case errors.Is(err, service.ErrStatusInvalid):
		fmt.Fprintln(os.Stderr, "Error: invalid status")
	case errors.Is(err, validator.ErrIDIsInvalid):
		fmt.Fprintln(os.Stderr, "Error: invalid id")
	default:
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func (h *Handler) HandleAdd(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: add requires exactly 1 argument (description)")
		osExit(1)
	}
	task, err := h.svc.AddTask(args[0])
	if err != nil {
		printError(err)
		osExit(1)
	}
	tasktemp, err := formater.FormatTaskJSON(task)
	if err != nil {
		printError(err)
		osExit(1)
	}
	fmt.Fprintln(os.Stdout, tasktemp)
}

func (h *Handler) HandleUpdate(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: update requires exactly 2 arguments (id, description)")
		osExit(1)
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printError(validator.ErrIDIsInvalid)
		osExit(1)
	}
	task, err := h.svc.UpdateTask(id, args[1])
	if err != nil {
		printError(err)
		osExit(1)
	}
	output, err := formater.FormatTaskJSON(task)
	if err != nil {
		printError(err)
		osExit(1)
	}
	fmt.Fprintln(os.Stdout, output)
}

func (h *Handler) HandleDelete(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: delete requires exactly 1 argument (id)")
		osExit(1)
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printError(validator.ErrIDIsInvalid)
		osExit(1)
	}
	if err := h.svc.DeleteTask(id); err != nil {
		printError(err)
		osExit(1)
	}
	fmt.Fprintf(os.Stdout, `{"deleted": %d}`+"\n", id)
}

func (h *Handler) HandleList(args []string, table bool) {
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: list accepts 0 or 1 argument [status]")
		osExit(1)
	}
	status := ""
	if len(args) == 1 {
		status = args[0]
	}
	tasks, err := h.svc.ListTasks(status)
	if err != nil {
		printError(err)
		osExit(1)
	}
	var output string
	if table {
		output, err = formater.FormatTasksTable(tasks)
	} else {
		output, err = formater.FormatTasksJSON(tasks)
	}
	if err != nil {
		printError(err)
		osExit(1)
	}
	fmt.Fprint(os.Stdout, output)
}

func (h *Handler) HandleMark(args []string, status domain.Status) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: mark requires exactly 1 argument (id)")
		osExit(1)
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printError(validator.ErrIDIsInvalid)
		osExit(1)
	}
	task, err := h.svc.MarkTask(id, status)
	if err != nil {
		printError(err)
		osExit(1)
	}
	output, err := formater.FormatTaskJSON(task)
	if err != nil {
		printError(err)
		osExit(1)
	}
	fmt.Fprintln(os.Stdout, output)
}
