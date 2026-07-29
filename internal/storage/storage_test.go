package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmk31585/task-cli/internal/domain"
)

func TestStorage(t *testing.T) {
	t.Run("non-existent file returns empty slice", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		tasks, err := s.Read()
		if err != nil {
			t.Errorf("Read() error = %v, want nil", err)
		}
		if tasks == nil {
			t.Error("Read() returned nil slice, want non-nil empty slice")
		}
		if len(tasks) != 0 {
			t.Errorf("Read() returned %d tasks, want 0", len(tasks))
		}
	})

	t.Run("valid file returns parsed tasks", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		raw := `[{"id":1,"description":"test","status":"done","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}]`
		if err := os.WriteFile(s.filePath, []byte(raw), 0644); err != nil {
			t.Fatal(err)
		}
		tasks, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(tasks) != 1 {
			t.Fatalf("Read() returned %d tasks, want 1", len(tasks))
		}
		if tasks[0].ID != 1 {
			t.Errorf("ID = %d, want 1", tasks[0].ID)
		}
		if tasks[0].Description != "test" {
			t.Errorf("Description = %q, want %q", tasks[0].Description, "test")
		}
		if tasks[0].Status != domain.StatusDone {
			t.Errorf("Status = %q, want %q", tasks[0].Status, domain.StatusDone)
		}
	})

	t.Run("corrupted JSON returns ErrCorruptedFile", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(s.filePath, []byte("{broken json"), 0644); err != nil {
			t.Fatal(err)
		}
		_, err = s.Read()
		if !errors.Is(err, ErrCorruptedFile) {
			t.Errorf("Read() error = %v, want %v", err, ErrCorruptedFile)
		}
	})

	t.Run("empty file returns empty slice", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(s.filePath, []byte("[]"), 0644); err != nil {
			t.Fatal(err)
		}
		tasks, err := s.Read()
		if err != nil {
			t.Errorf("Read() error = %v, want nil", err)
		}
		if len(tasks) != 0 {
			t.Errorf("Read() returned %d tasks, want 0", len(tasks))
		}
	})

	t.Run("Write creates the file", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(s.filePath); !os.IsNotExist(err) {
			t.Fatal("file should not exist before Write")
		}
		if err := s.Write([]domain.Task{}); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(s.filePath); err != nil {
			t.Errorf("file should exist after Write: %v", err)
		}
	})

	t.Run("Write produces valid JSON", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		tasks := []domain.Task{
			{
				ID:          1,
				Description: "test task",
				Status:      domain.StatusTodo,
			},
		}
		if err := s.Write(tasks); err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(s.filePath)
		if err != nil {
			t.Fatal(err)
		}
		if !json.Valid(raw) {
			t.Errorf("Write produced invalid JSON: %s", string(raw))
		}
		readTasks, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		if len(readTasks) != 1 {
			t.Fatalf("Read() returned %d tasks, want 1", len(readTasks))
		}
		if readTasks[0].ID != tasks[0].ID {
			t.Errorf("ID = %d, want %d", readTasks[0].ID, tasks[0].ID)
		}
		if readTasks[0].Description != tasks[0].Description {
			t.Errorf("Description = %q, want %q", readTasks[0].Description, tasks[0].Description)
		}
		if readTasks[0].Status != tasks[0].Status {
			t.Errorf("Status = %q, want %q", readTasks[0].Status, tasks[0].Status)
		}
	})

	t.Run("atomic write cleans up temp file", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		tmpPath := s.filePath + ".tmp"

		if err := s.Write([]domain.Task{}); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
			t.Errorf("temp file %q should be removed after Write", tmpPath)
		}
		if _, err := os.Stat(s.filePath); err != nil {
			t.Errorf("target file should exist: %v", err)
		}

		if err := s.Write([]domain.Task{{ID: 1, Description: "updated"}}); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
			t.Errorf("temp file %q should be removed after second Write", tmpPath)
		}
	})

	t.Run("Write creates data directory if missing", func(t *testing.T) {
		baseDir := t.TempDir()
		subDir := filepath.Join(baseDir, "sub", "nested")
		s, err := NewJSONStorage(subDir)
		if err != nil {
			t.Fatal(err)
		}
		if err := s.Write([]domain.Task{{ID: 1, Description: "test"}}); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(s.filePath); err != nil {
			t.Errorf("file should exist after Write: %v", err)
		}
	})

	t.Run("TASK_CLI_DATA_DIR overrides default path", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv("TASK_CLI_DATA_DIR", tmpDir)

		s, err := NewJSONStorage("")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(s.filePath, tmpDir) {
			t.Errorf("filePath = %q, want prefix %q", s.filePath, tmpDir)
		}
		if !strings.HasSuffix(s.filePath, "/tasks.json") {
			t.Errorf("filePath = %q, want suffix %q", s.filePath, "/tasks.json")
		}
	})

	t.Run("Read fails when file is not readable", func(t *testing.T) {
		s, err := NewJSONStorage(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(s.filePath, []byte("[]"), 0644); err != nil {
			t.Fatal(err)
		}
		os.Chmod(s.filePath, 0000)
		t.Cleanup(func() { os.Chmod(s.filePath, 0644) })

		_, err = s.Read()
		if err == nil {
			t.Error("Read() expected error, got nil")
		}
	})

	t.Run("Write fails when directory is not writable", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "readonly")
		s, err := NewJSONStorage(dir)
		if err != nil {
			t.Fatal(err)
		}
		parent := filepath.Dir(s.filePath)
		if err := os.MkdirAll(parent, 0755); err != nil {
			t.Fatal(err)
		}
		os.Chmod(parent, 0444)
		t.Cleanup(func() { os.Chmod(parent, 0755) })

		err = s.Write([]domain.Task{{ID: 1, Description: "x"}})
		if err == nil {
			t.Error("Write() expected error, got nil")
		}
	})

}
