package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/mmk31585/task-cli/internal/domain"
)

const fileName = "tasks.json"

var ErrCorruptedFile = errors.New("failed to unmarshal the JSON file")

type JSONStorage struct {
	filePath string
	mu       sync.Mutex
	lockFile *os.File
}

func NewJSONStorage(dataDir string) (*JSONStorage, error) {
	dir := dataDir
	if dir == "" {
		dir = os.Getenv("TASK_CLI_DATA_DIR")
	}
	if dir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home directory: %w", err)
		}
		dir = homeDir + "/.task-cli"
	}
	return &JSONStorage{
		filePath: dir + "/" + fileName,
	}, nil
}

func (s *JSONStorage) Read() ([]domain.Task, error) {
	raw, err := os.ReadFile(s.filePath)
	if errors.Is(err, fs.ErrNotExist) {
		return []domain.Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var tasks []domain.Task
	if err := json.Unmarshal(raw, &tasks); err != nil {
		return nil, ErrCorruptedFile
	}
	return tasks, nil
}

func (s *JSONStorage) Write(tasks []domain.Task) error {
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal tasks: %w", err)
	}
	tmpPath := s.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmpPath, s.filePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func (s *JSONStorage) Lock() error {
	f, err := os.OpenFile(s.filePath+".lock", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open lock file: %w", err)
	}
	if err := platformLock(f); err != nil {
		f.Close()
		return fmt.Errorf("acquire lock: %w", err)
	}
	s.mu.Lock()
	s.lockFile = f
	s.mu.Unlock()
	return nil
}

func (s *JSONStorage) Unlock() error {
	s.mu.Lock()
	lf := s.lockFile
	s.lockFile = nil
	s.mu.Unlock()
	if lf == nil {
		return nil
	}
	if err := platformUnlock(lf); err != nil {
		return fmt.Errorf("release lock: %w", err)
	}
	if err := lf.Close(); err != nil {
		return fmt.Errorf("close lock file: %w", err)
	}
	return nil
}
