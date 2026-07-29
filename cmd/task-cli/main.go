package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mmk31585/task-cli/internal/cli"
	"github.com/mmk31585/task-cli/internal/domain"
	"github.com/mmk31585/task-cli/internal/repository"
	"github.com/mmk31585/task-cli/internal/service"
	"github.com/mmk31585/task-cli/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	store, err := storage.NewJSONStorage("")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	repo := repository.NewJSONTaskRepository(store)
	svc := service.NewTaskService(repo)
	handler := cli.NewHandler(svc)

	switch os.Args[1] {
	case "add":
		handler.HandleAdd(os.Args[2:])
	case "update":
		handler.HandleUpdate(os.Args[2:])
	case "delete":
		handler.HandleDelete(os.Args[2:])
	case "list":
		fs := flag.NewFlagSet("list", flag.ExitOnError)
		table := fs.Bool("table", false, "display output in table format")
		fs.Parse(os.Args[2:])
		handler.HandleList(fs.Args(), *table)
	case "mark-in-progress":
		handler.HandleMark(os.Args[2:], domain.StatusInProgress)
	case "mark-done":
		handler.HandleMark(os.Args[2:], domain.StatusDone)
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Fprint(os.Stderr, `Task Tracker CLI

Usage:
  task-cli <command> [arguments]

Commands:
  add              Add a new task
  update           Update a task description
  delete           Delete a task
  mark-in-progress Mark a task as in progress
  mark-done        Mark a task as done
  list             List tasks (optionally filtered by status)

Flags:
  --table          Display output in table format (default: JSON)
  --help, -h       Show help
`)
}
