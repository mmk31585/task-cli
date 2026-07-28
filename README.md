# Task Tracker CLI

A minimalist, production-grade CLI task tracker written in **Go** following **Clean Architecture**, **Domain-Driven Design**, and **SOLID** principles. Persists tasks to a JSON file using only the Go standard library.

---

## Features

- **Create** tasks with a description
- **List** tasks filtered by status (`todo`, `in-progress`, `done`)
- **Update** task descriptions
- **Delete** tasks
- **Mark** tasks as `in-progress` / `done`
- **JSON** output by default with optional human-readable `--table` format
- **UUIDv7** identifiers (time-ordered, unique)
- **No external dependencies** — 100% Go standard library

---

## Goals

1. Demonstrate professional Go project structure with layered architecture.
2. Serve as a portfolio piece that showcases Clean Architecture, DDD, and SOLID in a small CLI project.
3. Provide a real, usable tool for personal task tracking.

---

## Requirements

- Go 1.24+

---

## Installation

```bash
# Clone the repository
git clone https://github.com/mmk31585/task-cli.git
cd task-cli

# Build
go build -o task-cli ./cmd/task-cli

# (Optional) Install to $GOPATH/bin
go install ./cmd/task-cli
```

---

## Usage

```
Usage: task-cli <command> [arguments]

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
```

---

## Example Commands

```bash
# Add a new task
task-cli add "Buy groceries"

# Update a task
task-cli update 1 "Buy groceries and cook dinner"

# Delete a task
task-cli delete 1

# Mark a task as in progress
task-cli mark-in-progress 1

# Mark a task as done
task-cli mark-done 1

# List all tasks (JSON default)
task-cli list

# List tasks by status
task-cli list done
task-cli list todo
task-cli list in-progress

# Display as a formatted table
task-cli list --table
task-cli list done --table
```

---

## Project Structure

```
task-cli/
├── cmd/
│   └── task-cli/
│       └── main.go              # Application entry point
├── internal/
│   ├── cli/
│   │   └── handler.go           # CLI command handlers (thin layer)
│   ├── domain/
│   │   ├── task.go              # Task entity, status type
│   │   └── errors.go            # Domain-specific error definitions
│   ├── formater/
│   │   └── formater.go          # JSON and table output formatting
│   ├── repository/
│   │   └── repository.go        # TaskRepository interface + JSON implementation
│   ├── service/
│   │   └── service.go           # Business logic / use-case orchestration
│   ├── storage/
│   │   └── storage.go           # Low-level JSON file read/write + locking
│   └── validator/
│       └── validator.go         # Input validation rules
├── pkg/
│   └── uuid/
│       └── uuid.go              # UUIDv7 generator (public, reusable)
├── docs/
│   └── software-design-document.md  # Complete architecture document
├── go.mod
├── go.sum
└── README.md
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────┐
│                  cmd/task-cli                │
│           (Entry Point / Main)              │
├─────────────────────────────────────────────┤
│               internal/cli                   │
│          (CLI Handlers — thin layer)         │
├─────────────────────────────────────────────┤
│              internal/service               │
│        (Business Logic / Use Cases)          │
├─────────────────────────────────────────────┤
│          internal/repository                 │
│   ┌─────────────────────────────────────┐   │
│   │  TaskRepository (interface)         │   │
│   │  JSONTaskRepository (implementation)│   │
│   └─────────────────────────────────────┘   │
├─────────────────────────────────────────────┤
│           internal/domain                   │
│      (Entities, Value Objects, Errors)      │
├─────────────────────────────────────────────┤
│  internal/storage  │  internal/validator    │
│  (JSON File I/O)   │  (Input Validation)   │
├─────────────────────────────────────────────┤
│  internal/formater │  pkg/uuid             │
│  (Output Display)  │  (UUIDv7 Generator)   │
└─────────────────────────────────────────────┘
```

**Dependency flow:** `cmd → cli → service → repository (impl) → storage | domain`

Higher layers depend on abstractions (interfaces), never on concrete implementations. The domain layer has zero external dependencies.

---

## Development Workflow

```bash
# Run directly
go run ./cmd/task-cli list

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code
go vet ./...

# Build
go build -o task-cli ./cmd/task-cli
```

---

## Testing

| Layer       | Scope                          | Approach                     |
|-------------|--------------------------------|------------------------------|
| Domain      | Entity validation, state transitions | Pure unit tests, no mocks   |
| Service     | Business logic orchestration    | Unit tests with mocked repository |
| Repository  | JSON file operations            | Unit tests with temp files   |
| CLI         | Command parsing + output        | Table-driven integration tests |
| Formater    | JSON/table rendering            | Golden file comparison       |

See [docs/software-design-document.md](./docs/software-design-document.md#13-testing-strategy) for details.

---

## Future Improvements

- SQLite repository implementation (swap via interface)
- REST API server
- Web UI
- Terminal UI (TUI)
- Tags / categories
- Priority levels
- Due dates / deadlines
- Recurring tasks
- Data import/export (CSV, JSON)
- Cloud sync

---

## Contributing

Contributions are welcome. Please open an issue first to discuss changes.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit with [conventional commits](https://www.conventionalcommits.org/)
4. Open a Pull Request

---

## License

MIT

---

## Implementation Roadmap

### Progress Tracking

| Phase | Area | Est. Effort | Status |
|---|---|---|---|
| 0 | Definition of Ready | — | Pending |
| 1 | Bootstrap | 15 min | Pending |
| 2 | Domain | 30 min | Pending |
| 3 | ID Generation (`pkg/uuid`) | 20 min | Pending |
| 4 | JSON Storage | 45 min | Pending |
| 5 | Repository | 30 min | Pending |
| 6 | Validator | 15 min | Pending |
| 7 | Formatter | 20 min | Pending |
| 8 | Service (Business Logic) | 45 min | Pending |
| 9 | CLI Handler + Commands | 1 hr | Pending |
| 10 | Error Handling | 30 min | Pending |
| 11 | Testing | 2 hr | Pending |
| 12 | Documentation | 1 hr | Pending |
| 13 | Refactoring & Code Quality | 1 hr | Pending |
| 14 | Release v1.0.0 | 30 min | Pending |

---

### Definition of Ready

Before any implementation begins:

- [ ] All ADRs from the Software Design Document are reviewed and accepted
- [ ] The package directory structure exists (`cmd/`, `internal/*`, `pkg/`)
- [ ] Go 1.24+ is installed (`go version`)
- [ ] `go.mod` is initialized with the correct module path
- [ ] The project builds with zero dependencies (`go build ./...`)
- [ ] `go vet ./...` passes on an empty project
- [ ] `gofmt -s .` produces no diffs
- [ ] All stakeholders agree on the command interface (names, flags, arguments)
- [ ] The JSON data format is finalized (field names, types, timestamps)
- [ ] Git branch `main` is ready with an initial commit

---

### Phase 1 — Bootstrap

**Objective:** Initialize the Go module and project skeleton. Ensure the build pipeline works.

**Deliverables:**
- `go.mod` / `go.sum`
- Empty Go files in each package to reserve the directory structure
- Working `go build`, `go vet`, `go test`

**TODO Checklist:**

- ✅ Run `go mod init github.com/mmk31585/task-cli`
- ✅ Create `cmd/task-cli/main.go` with a minimal `func main()` that prints a help message
- ✅ Create placeholder files in each `internal/` package (one `.go` file per package with `package <name>`)
- ✅ Create `pkg/uuid/uuid.go` placeholder
- ✅ Verify `go build ./...` succeeds
- ✅ Verify `go vet ./...` succeeds
- ✅ Verify `go test ./...` succeeds (zero tests is OK at this stage)
- ✅ Verify `gofmt -s .` produces no diffs
- ✅ Commit: `chore: bootstrap Go module and project skeleton`

**Testing Checklist:**

- ✅ No tests required for this phase
- ✅ CI passes `go build` and `go vet`

**Definition of Done:**

- ✅ `go build -o /dev/null ./cmd/task-cli` produces a binary
- ✅ `go vet ./...` is clean
- ✅ Every package directory has at least one `.go` file
- ✅ `go.mod` exists with module path `github.com/mmk31585/task-cli`

---

### Phase 2 — Domain Layer

**Objective:** Define the core domain types — `Task`, `Status`, and domain errors.

**Deliverables:**
- `internal/domain/task.go` — `Task` struct, `Status` type, constants
- `internal/domain/errors.go` — Sentinel errors

**TODO Checklist:**

- ✅ Define `Status` as a typed string with constants `StatusTodo`, `StatusInProgress`, `StatusDone`
- ✅ Define `IsValidStatus(s Status) bool` function
- ✅ Define `Task` struct with fields: `ID`, `Description`, `Status`, `CreatedAt`, `UpdatedAt`
- ✅ Add JSON struct tags (`json:"id"`, `json:"description"`, etc.)
- ✅ Use `time.Time` for timestamp fields
- ✅ Define sentinel errors: `ErrTaskNotFound`, `ErrInvalidStatus`, `ErrEmptyDescription`, `ErrDescriptionTooLong`
- ✅ Verify package compiles: `go build ./internal/domain/...`
- ✅ Commit: `feat(domain): add Task entity, Status type, and domain errors`

**Testing Checklist:**

- ✅ Test `IsValidStatus("todo")` returns `true`
- ✅ Test `IsValidStatus("invalid")` returns `false`
- ✅ Test sentinel errors are non-nil and comparable with `errors.Is`
- ✅ Test `Task` struct fields have correct JSON tags

**Definition of Done:**

- ✅ All domain types compile and pass tests
- ✅ Package imports only `time` and `encoding/json` from stdlib
- ✅ Package has zero imports from other project packages



### Phase 4 — JSON Storage

**Objective:** Implement low-level JSON file read/write with atomic writes.

**Deliverables:**
- `internal/storage/storage.go` — `JSONStorage` struct with `Read` and `Write`

**TODO Checklist:**

- [ ] Define `JSONStorage` struct with `filePath string` field
- [ ] Write `NewJSONStorage(dataDir string) *JSONStorage` constructor
- [ ] Resolve default data directory as `~/.task-cli` using `os.UserHomeDir()`
- [ ] Support `TASK_CLI_DATA_DIR` environment variable override
- [ ] Implement `Read() ([]domain.Task, error)`:
  - [ ] Return empty slice if file does not exist (first run)
  - [ ] Use `os.ReadFile`
  - [ ] Unmarshal with `json.Unmarshal`
  - [ ] Return `ErrCorruptedFile` on JSON syntax errors
- [ ] Implement `Write(tasks []domain.Task) error`:
  - [ ] Marshal with `json.MarshalIndent` (2-space indent)
  - [ ] Write to a temporary file in the same directory
  - [ ] Use `os.Rename` for atomic replacement
  - [ ] Create data directory with `os.MkdirAll` if needed
- [ ] Ensure all exported functions have Go doc comments
- [ ] Commit: `feat(storage): add JSON file storage with atomic writes`

**Testing Checklist:**

- [ ] Test `Read` on non-existent file returns empty slice (not error)
- [ ] Test `Read` on valid file returns parsed tasks
- [ ] Test `Read` on corrupted JSON returns `ErrCorruptedFile`
- [ ] Test `Write` creates the file
- [ ] Test `Write` produces valid JSON (read back and compare)
- [ ] Test atomic write: write, then verify original file is intact
- [ ] Test `Write` creates the data directory if missing
- [ ] Test `TASK_CLI_DATA_DIR` env var overrides the default path
- [ ] Test `Read` on empty file returns empty slice

**Definition of Done:**

- [ ] All CRUD operations on the file layer work correctly
- [ ] Atomic write prevents file corruption on failure
- [ ] Coverage > 90%

---

### Phase 5 — Repository

**Objective:** Define the `TaskRepository` interface and implement the JSON-backed repository.

**Deliverables:**
- `internal/repository/repository.go` — Interface + `JSONTaskRepository`

**TODO Checklist:**

- [ ] Define `TaskRepository` interface with:
  - [ ] `Add(task domain.Task) error`
  - [ ] `GetAll() ([]domain.Task, error)`
  - [ ] `GetByID(id string) (domain.Task, error)`
  - [ ] `Update(task domain.Task) error`
  - [ ] `Delete(id string) error`
- [ ] Define `JSONTaskRepository` struct with `store *storage.JSONStorage` field
- [ ] Write `NewJSONTaskRepository(store *storage.JSONStorage) *JSONTaskRepository`
- [ ] Implement `Add`: read all, append, write all
- [ ] Implement `GetAll`: read all, return
- [ ] Implement `GetByID`: read all, linear scan, return `ErrTaskNotFound` if missing
- [ ] Implement `Update`: read all, find by ID, replace, write all
- [ ] Implement `Delete`: read all, filter out by ID, write all
- [ ] Commit: `feat(repository): add TaskRepository interface and JSON implementation`

**Testing Checklist:**

- [ ] Test `Add` then `GetAll` returns the added task
- [ ] Test `Add` multiple tasks returns correct count
- [ ] Test `GetByID` returns correct task by ID
- [ ] Test `GetByID` with unknown ID returns `ErrTaskNotFound`
- [ ] Test `Update` modifies description
- [ ] Test `Update` with unknown ID returns `ErrTaskNotFound`
- [ ] Test `Delete` removes task
- [ ] Test `Delete` with unknown ID returns `ErrTaskNotFound`
- [ ] Test all operations preserve other tasks (no data loss)

**Definition of Done:**

- [ ] `JSONTaskRepository` satisfies `TaskRepository` interface (checked by Go compiler)
- [ ] All operations work correctly with a temp file backing store
- [ ] Coverage > 90%

---

### Phase 6 — Validator

**Objective:** Implement pure input validation functions.

**Deliverables:**
- `internal/validator/validator.go` — Validation functions

**TODO Checklist:**

- [ ] Implement `ValidateDescription(desc string) error`:
  - [ ] Return `ErrEmptyDescription` if `desc` is empty or whitespace-only
  - [ ] Return `ErrDescriptionTooLong` if `len(desc) > 500`
  - [ ] Return `nil` otherwise
- [ ] Implement `ValidateID(id string) error`:
  - [ ] Return error if `id` is empty
- [ ] Ensure all functions are pure (no I/O, no state)
- [ ] Commit: `feat(validator): add input validation functions`

**Testing Checklist:**

- [ ] Test `ValidateDescription("")` returns error
- [ ] Test `ValidateDescription("  ")` returns error
- [ ] Test `ValidateDescription(strings.Repeat("a", 501))` returns error
- [ ] Test `ValidateDescription("Buy milk")` returns `nil`
- [ ] Test `ValidateDescription(strings.Repeat("a", 500))` returns `nil`
- [ ] Test `ValidateID("")` returns error
- [ ] Test `ValidateID("abc")` returns `nil`

**Definition of Done:**

- [ ] All validation functions are pure (no side effects)
- [ ] Boundary conditions tested (0, 1, 500, 501 characters)
- [ ] Coverage 100%

---

### Phase 7 — Formatter

**Objective:** Implement output formatting — JSON and table.

**Deliverables:**
- `internal/formater/formater.go` — Formatting functions

**TODO Checklist:**

- [ ] Implement `FormatTaskJSON(task domain.Task) (string, error)`:
  - [ ] Use `json.MarshalIndent` with 2-space indent
  - [ ] Return formatted JSON string
- [ ] Implement `FormatTasksJSON(tasks []domain.Task) (string, error)`:
  - [ ] Marshal the entire slice as a JSON array
- [ ] Implement `FormatTasksTable(tasks []domain.Task) (string, error)`:
  - [ ] Use `text/tabwriter`
  - [ ] Columns: ID, Description, Status, Created At, Updated At
  - [ ] Header row with column names
  - [ ] Align columns properly
  - [ ] Handle empty slice (print "No tasks found")
- [ ] Commit: `feat(formater): add JSON and table output formatters`

**Testing Checklist:**

- [ ] Test `FormatTaskJSON` produces valid JSON with correct fields
- [ ] Test `FormatTasksJSON` produces a JSON array
- [ ] Test `FormatTasksTable` produces a formatted table with headers
- [ ] Test `FormatTasksTable` with empty slice prints "No tasks found"
- [ ] Test table format aligns columns correctly (fixed-width test)
- [ ] Test edge case: single task, many tasks

**Definition of Done:**

- [ ] Both output formats work correctly
- [ ] Table output uses `text/tabwriter` for alignment
- [ ] Coverage > 95%

---

### Phase 8 — Service (Business Logic)

**Objective:** Implement the core business logic / use case orchestration layer.

**Deliverables:**
- `internal/service/service.go` — `TaskService` struct and methods

**TODO Checklist:**

- [ ] Define `TaskService` struct with `repo repository.TaskRepository` field
- [ ] Write `NewTaskService(repo repository.TaskRepository) *TaskService`
- [ ] Implement `AddTask(description string) (domain.Task, error)`:
  - [ ] Call `validator.ValidateDescription`
  - [ ] Generate UUIDv7 via `uuid.NewV7()`
  - [ ] Create `domain.Task` with status `StatusTodo`, timestamps set to `time.Now()`
  - [ ] Call `repo.Add`
  - [ ] Return the created task
- [ ] Implement `ListTasks(status string) ([]domain.Task, error)`:
  - [ ] Call `repo.GetAll()`
  - [ ] If status is empty, return all tasks
  - [ ] If status is set, filter by status (validate status first)
- [ ] Implement `UpdateTask(id, description string) (domain.Task, error)`:
  - [ ] Validate description
  - [ ] Get task by ID
  - [ ] Update description and `UpdatedAt`
  - [ ] Persist via `repo.Update`
- [ ] Implement `DeleteTask(id string) error`:
  - [ ] Check task exists (call `GetByID` — let it return `ErrTaskNotFound`)
  - [ ] Delete via `repo.Delete`
- [ ] Implement `MarkTask(id string, status domain.Status) (domain.Task, error)`:
  - [ ] Validate status with `domain.IsValidStatus`
  - [ ] Get task by ID
  - [ ] Set status and `UpdatedAt`
  - [ ] Persist via `repo.Update`
- [ ] Commit: `feat(service): add business logic layer with use cases`

**Testing Checklist:**

- [ ] Test `AddTask` with valid description returns task with `StatusTodo`
- [ ] Test `AddTask` with empty description returns error (no I/O)
- [ ] Test `AddTask` with long description returns error (no I/O)
- [ ] Test `ListTasks("")` returns all tasks
- [ ] Test `ListTasks("done")` returns only done tasks
- [ ] Test `ListTasks("invalid")` returns error
- [ ] Test `UpdateTask` with valid input updates description and `UpdatedAt`
- [ ] Test `UpdateTask` with unknown ID returns `ErrTaskNotFound`
- [ ] Test `DeleteTask` removes task
- [ ] Test `DeleteTask` with unknown ID returns `ErrTaskNotFound`
- [ ] Test `MarkTask` transitions status correctly
- [ ] Test `MarkTask` with unknown ID returns `ErrTaskNotFound`
- [ ] Test `MarkTask` with invalid status returns error

**Definition of Done:**

- [ ] All use cases operate correctly with a mock repository
- [ ] No business logic leaks into the CLI or repository layers
- [ ] Coverage > 90%

---

### Phase 9 — CLI Handler and Commands

**Objective:** Implement the CLI entry point, command dispatch, and all individual commands.

**Deliverables:**
- `internal/cli/handler.go` — CLI handler struct
- `cmd/task-cli/main.go` — Full entry point with DI wiring

**TODO Checklist:**

- [ ] Define `Handler` struct with `svc *service.TaskService` field
- [ ] Write `NewHandler(svc *service.TaskService) *Handler`
- [ ] Implement `HandleAdd(args []string)`:
  - [ ] Require exactly 1 argument (description)
  - [ ] Call `svc.AddTask`
  - [ ] Format output with `formater.FormatTaskJSON`
  - [ ] Handle error (print to stderr, exit 1)
- [ ] Implement `HandleUpdate(args []string)`:
  - [ ] Require exactly 2 arguments (id, description)
  - [ ] Call `svc.UpdateTask`
  - [ ] Output the updated task as JSON
- [ ] Implement `HandleDelete(args []string)`:
  - [ ] Require exactly 1 argument (id)
  - [ ] Call `svc.DeleteTask`
  - [ ] Output success message `{"deleted": "<id>"}`
- [ ] Implement `HandleList(args []string, table bool)`:
  - [ ] Accept 0 or 1 argument (optional status filter)
  - [ ] Call `svc.ListTasks`
  - [ ] If `table` flag is set, call `formater.FormatTasksTable`
  - [ ] Otherwise call `formater.FormatTasksJSON`
- [ ] Implement `HandleMark(args []string, status domain.Status)`:
  - [ ] Require exactly 1 argument (id)
  - [ ] Call `svc.MarkTask`
  - [ ] Output the updated task as JSON
- [ ] In `cmd/task-cli/main.go`:
  - [ ] Parse `os.Args[1]` for the command name
  - [ ] Wire dependencies: `storage → repository → service → handler`
  - [ ] Use `flag.FlagSet` per command for `--table`, `--help`
  - [ ] Implement `help` command / `--help` flag
  - [ ] Dispatch to the correct handler method
  - [ ] Set `os.Exit(0)` on success, `os.Exit(1)` on error
- [ ] Test all commands manually:
  - [ ] `task-cli add "Test task"` → success
  - [ ] `task-cli list` → shows tasks
  - [ ] `task-cli list --table` → table format
  - [ ] `task-cli list done` → filtered
  - [ ] `task-cli update <id> "New desc"` → updated
  - [ ] `task-cli delete <id>` → deleted
  - [ ] `task-cli mark-in-progress <id>` → status changed
  - [ ] `task-cli mark-done <id>` → status changed
  - [ ] `task-cli add ""` → error message
  - [ ] `task-cli` (no args) → help text
  - [ ] `task-cli --help` → help text
  - [ ] `task-cli unknown` → "unknown command" error
- [ ] Commit: `feat(cli): add CLI handler, command dispatch, and all commands`

**Testing Checklist:**

- [ ] Test `HandleAdd` with correct args outputs JSON to stdout
- [ ] Test `HandleAdd` with missing args writes error to stderr, exits 1
- [ ] Test `HandleList` outputs JSON array
- [ ] Test `HandleList --table` outputs formatted table
- [ ] Test `HandleList done` outputs only done tasks
- [ ] Test `HandleUpdate` with correct args outputs updated task
- [ ] Test `HandleDelete` outputs success JSON
- [ ] Test `HandleMark` outputs updated task with new status
- [ ] Test unknown command returns error
- [ ] Test no command prints help
- [ ] Test `--help` flag prints help
- [ ] Test all error paths produce non-zero exit code

**Definition of Done:**

- [ ] All 7 commands work correctly from the terminal
- [ ] Error messages go to stderr, output goes to stdout
- [ ] Exit codes are correct (0 success, 1 error)
- [ ] `--help` and `--table` work as documented

---

### Phase 10 — Error Handling

**Objective:** Polish error handling — ensure consistent wrapping, user-friendly messages, and proper propagation.

**Deliverables:**
- Updates to `internal/service/service.go`, `internal/cli/handler.go`, `cmd/task-cli/main.go`

**TODO Checklist:**

- [ ] Ensure every error returned from storage is wrapped with context in the repository layer
- [ ] Ensure every error returned from repository is wrapped with context in the service layer
- [ ] Ensure sentinel errors (`ErrTaskNotFound`, `ErrInvalidStatus`, etc.) are preserved with `%w` wrapping
- [ ] In CLI handler, use `errors.Is` to detect sentinel errors and print user-friendly messages
- [ ] Handle unexpected panics: add `defer/recover` in `main.go` to prevent stack traces leaking to the user
- [ ] Ensure all `fmt.Fprintln(os.Stderr, ...)` calls use consistent formatting: `"Error: <message>"`
- [ ] Test error messages read naturally (e.g., "task with ID abc-123 not found" not "error: task not found")
- [ ] Commit: `refactor(errors): polish error wrapping and user-facing messages`

**Testing Checklist:**

- [ ] Test that `ErrTaskNotFound` propagated through CLI shows a user-friendly message
- [ ] Test that `ErrCorruptedFile` shows a message suggesting how to fix it
- [ ] Test that an unexpected panic in the service layer is caught and does not crash with stack trace
- [ ] Test all error paths in every command

**Definition of Done:**

- [ ] No raw errors are printed to the user
- [ ] All errors use `%w` wrapping to preserve the error chain
- [ ] `errors.Is` works across all layers
- [ ] Panics are caught and converted to user-friendly errors

---

### Phase 11 — Testing

**Objective:** Write the complete test suite across all layers.

**Deliverables:**
- Test files for every package
- Coverage report

**TODO Checklist:**

- [ ] Write domain tests: `internal/domain/task_test.go`
  - [ ] Test `IsValidStatus` for all valid and invalid values
  - [ ] Test sentinel error comparisons
- [ ] Write UUID tests: `pkg/uuid/uuid_test.go`
  - [ ] Test format, uniqueness, time-ordering
- [ ] Write storage tests: `internal/storage/storage_test.go`
  - [ ] Use `os.CreateTemp` for isolated test directories
  - [ ] Test read/write/corruption/atomicity
- [ ] Write repository tests: `internal/repository/repository_test.go`
  - [ ] Use real `JSONStorage` with temp files (not mocks)
  - [ ] Test all CRUD operations
- [ ] Write validator tests: `internal/validator/validator_test.go`
  - [ ] Test boundary conditions
- [ ] Write formatter tests: `internal/formater/formater_test.go`
  - [ ] Test JSON and table output
  - [ ] Golden file comparison for table format
- [ ] Write service tests: `internal/service/service_test.go`
  - [ ] Use a hand-written in-memory mock repository
  - [ ] Test all use cases, error paths, edge cases
- [ ] Write CLI handler tests: `internal/cli/handler_test.go`
  - [ ] Capture stdout/stderr with `bytes.Buffer`
  - [ ] Use a mock service
  - [ ] Test each command
- [ ] Run `go test -race ./...` and fix any races
- [ ] Run `go test -coverprofile=coverage.out ./...` and verify coverage
- [ ] Review coverage report: `go tool cover -html=coverage.out`
- [ ] Commit: `test: add complete test suite for all layers`

**Testing Checklist:**

- [ ] All tests pass: `go test ./...`
- [ ] Race detector clean: `go test -race ./...`
- [ ] Coverage > 80%: `go test -cover ./...`
- [ ] Every exported function has at least one test
- [ ] Edge cases covered: empty list, invalid input, corrupted file, missing file

**Definition of Done:**

- [ ] `go test -race -cover ./...` passes with > 80% coverage
- [ ] Mock repository is available for service tests
- [ ] All test files follow the same pattern
- [ ] Tests are order-independent (can run in any order)

---

### Phase 12 — Documentation

**Objective:** Finalize all project documentation.

**Deliverables:**
- Complete `README.md`
- Complete `docs/software-design-document.md`
- Go doc comments on all exported symbols

**TODO Checklist:**

- [ ] Review `README.md` for completeness
- [ ] Review `docs/software-design-document.md` for accuracy
- [ ] Ensure all exported types, functions, and constants have Go doc comments
- [ ] Add package-level doc comments to each package
- [ ] Verify `go doc ./...` output is readable
- [ ] Ensure every ADR in the SDD is up to date with actual implementation decisions
- [ ] Commit: `docs: finalize project documentation`

**Testing Checklist:**

- [ ] `go doc ./...` produces no errors
- [ ] All doc comments are complete sentences
- [ ] No `TODO` comments remain in the code

**Definition of Done:**

- [ ] README is complete and accurate
- [ ] SDD matches the implemented code
- [ ] All exported symbols are documented

---

### Phase 13 — Refactoring and Code Quality

**Objective:** Review the entire codebase for quality, consistency, and adherence to design principles.

**Deliverables:**
- Clean, consistent codebase passing all quality checks

**Code Quality Checklist:**

- [ ] `gofmt -s .` produces no diffs
- [ ] `go vet ./...` is clean
- [ ] `go mod tidy` has been run
- [ ] No duplicate code (DRY principle)
- [ ] Clean Architecture dependency rule is respected (no inward violations)
- [ ] SOLID principles are followed throughout
- [ ] KISS is respected (no unnecessary abstractions)
- [ ] YAGNI is respected (no unused code or speculative features)
- [ ] No `init()` functions
- [ ] No global variables (except sentinel errors)
- [ ] No panics (except recovered in main)
- [ ] No unused imports or variables (`go vet` would catch these)
- [ ] All function signatures are consistent
- [ ] No magic numbers or strings (use constants)
- [ ] No commented-out code
- [ ] No `_` test files without matching source files
- [ ] All files have a consistent license header (optional)
- [ ] `go mod verify` passes
- [ ] `go build -o /dev/null ./cmd/task-cli` succeeds
- [ ] Commit: `refactor: code quality review and cleanup`

**Definition of Done:**

- [ ] Entire checklist above is satisfied
- [ ] No warnings from `go vet` or `gofmt`
- [ ] Codebase follows the project's coding standards

---

### Phase 14 — Release v1.0.0

**Objective:** Tag and publish the first stable release.

**Deliverables:**
- Git tag `v1.0.0`
- Release notes

**Release Checklist:**

- [ ] All previous phases are complete and committed
- [ ] `main` branch is up to date
- [ ] `go build -o task-cli ./cmd/task-cli` produces a working binary
- [ ] Final manual smoke test of all commands:
  - [ ] `task-cli add "Smoke test"` → creates task
  - [ ] `task-cli list` → shows task
  - [ ] `task-cli list --table` → table format OK
  - [ ] `task-cli update <id> "Updated"` → updates
  - [ ] `task-cli mark-in-progress <id>` → status changes
  - [ ] `task-cli mark-done <id>` → status changes
  - [ ] `task-cli list done` → filter works
  - [ ] `task-cli list todo` → filter works
  - [ ] `task-cli delete <id>` → deletes
  - [ ] `task-cli list` → empty list
  - [ ] `task-cli add ""` → error shown
  - [ ] `task-cli --help` → help shown
- [ ] `go test -race -cover ./...` passes
- [ ] Create and push tag:
  ```bash
  git tag -a v1.0.0 -m "v1.0.0: Initial stable release"
  git push origin v1.0.0
  ```
- [ ] Write release notes summarizing features

**Definition of Done (Project):**

- [ ] All functional requirements (FR-01 through FR-12) are implemented
- [ ] All non-functional requirements (NFR-01 through NFR-07) are satisfied
- [ ] Zero external dependencies
- [ ] Binary compiles as a single static executable
- [ ] Test coverage exceeds 80%
- [ ] All code quality checks pass
- [ ] Documentation is complete and accurate
- [ ] Git tag `v1.0.0` is published
- [ ] Project is ready for public consumption

---

### Testing Milestones

| Milestone | Phase | Verification |
|---|---|---|
| Domain tests pass | Phase 2 | `go test ./internal/domain/...` |
| UUID generation works | Phase 3 | `go test ./pkg/uuid/...` |
| File I/O works (temp files) | Phase 4 | `go test ./internal/storage/...` |
| CRUD operations correct | Phase 5 | `go test ./internal/repository/...` |
| Validation rules correct | Phase 6 | `go test ./internal/validator/...` |
| Output formatting correct | Phase 7 | `go test ./internal/formater/...` |
| Business logic correct | Phase 8 | `go test ./internal/service/...` |
| CLI integration works | Phase 9 | Manual testing + `go test ./internal/cli/...` |
| Error messages are user-friendly | Phase 10 | Manual error path testing |
| Full test suite passes | Phase 11 | `go test -race -cover ./...` |
| Documentation is accurate | Phase 12 | Visual review |
| Code quality gates pass | Phase 13 | `gofmt`, `go vet`, lint checks |
| Release is ready | Phase 14 | Full smoke test + tag |

---

### Suggested Git Branch Strategy

```
main              ●──────●──────●──────●───────────────● v1.0.0
                   \    / \    / \    /                 /
phase/1-bootstrap  ●──●   ●──●   ●──●                 /
phase/2-domain           ●──●                          /
phase/3-uuid                 ●──●                      /
phase/4-storage                 ●──●                  /
...                                                     /
                                                       /
release/v1.0.0  ──────────────────────────────────────●
```

**Rules:**

- `main` is always in a working state (builds, tests pass).
- Each phase gets a branch: `phase/<number>-<name>` (e.g., `phase/1-bootstrap`).
- Branches are created from `main`, merged back via PR.
- Use squash merge to keep history clean.
- After all phases, create a `release/v1.0.0` branch for final testing and tag.
- Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/).

**Branch naming patterns:**

| Branch | Purpose |
|---|---|
| `phase/1-bootstrap` | Project skeleton |
| `phase/2-domain` | Domain layer |
| `phase/3-uuid` | UUID generation |
| `phase/4-storage` | JSON storage |
| `phase/5-repository` | Repository |
| `phase/6-validator` | Validation |
| `phase/7-formatter` | Output formatting |
| `phase/8-service` | Business logic |
| `phase/9-cli` | CLI handlers |
| `phase/10-error-handling` | Error handling polish |
| `phase/11-testing` | Test suite |
| `phase/12-documentation` | Documentation |
| `phase/13-refactoring` | Code quality |
| `release/v1.0.0` | Release preparation |
