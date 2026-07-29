package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/formater"
	"github.com/mmk31585/task-cli/internal/service"
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

func (h *Handler) HandleAdd(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "add requires exactly 1 argument")
		osExit(1)
	}
	task, err := h.svc.AddTask(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		osExit(1)
	}
	tasktemp, err := formater.FormatTaskJSON(task)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		osExit(1)
	}
	fmt.Fprintln(os.Stdout, tasktemp)
}

func (h *Handler) HandleUpdate(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "update requires exactly 2 arguments: id description")
		osExit(1)
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid id")
		osExit(1)
	}
	task, err := h.svc.UpdateTask(id, args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
	}
	output, err := formater.FormatTaskJSON(task)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
	}
	fmt.Fprintln(os.Stdout, output)
}

func (h *Handler) HandleDelete(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "delete requires exactly 1 argument: id")
		osExit(1)
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid id")
		osExit(1)
	}
	if err := h.svc.DeleteTask(id); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
	}
	fmt.Fprintf(os.Stdout, `{"deleted": %d}`+"\n", id)
}

func (h *Handler) HandleList(args []string, table bool) {
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "list accepts 0 or 1 argument: [status]")
		osExit(1)
	}
	status := ""
	if len(args) == 1 {
		status = args[0]
	}
	tasks, err := h.svc.ListTasks(status)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
	}
	var output string
	if table {
		output, err = formater.FormatTasksTable(tasks)
	} else {
		output, err = formater.FormatTasksJSON(tasks)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
	}
	fmt.Fprintln(os.Stdout, output)
}

func (h *Handler) HandleMark(args []string, status domain.Status) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "mark requires exactly 1 argument: id")
		osExit(1)
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: invalid id")
		osExit(1)
	}
	task, err := h.svc.MarkTask(id, status)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
	}
	output, err := formater.FormatTaskJSON(task)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
	}
	fmt.Fprintln(os.Stdout, output)
}
