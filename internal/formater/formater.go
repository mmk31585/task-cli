package formater

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/mmk31585/task-cli/internal/domain"
)

var (
	ErrNotFoundTask = errors.New("no tasks found")
)

func FormatTaskJSON(task domain.Task) (string, error) {
	by, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return "", err
	}
	return string(by), nil
}
func FormatTasksJSON(tasks []domain.Task) (string, error) {
	by, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return "", err
	}
	return string(by), nil
}
func FormatTasksTable(tasks []domain.Task) (string, error) {
	if len(tasks) == 0 {
		return "No tasks found\n", ErrNotFoundTask
	}
	var buf bytes.Buffer
	writer := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tDescription\tStatus\tCreated At\tUpdated At")

	for _, t := range tasks {
		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\n",
			t.ID, t.Description, t.Status,
			t.CreatedAt.Format(time.RFC3339),
			t.UpdatedAt.Format(time.RFC3339))
	}
	writer.Flush()
	return buf.String(), nil
}
