package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicolaslecrique/LabelClaw/backend/internal/configuration"
)

var ErrNotFound = errors.New("active configuration not found")

type FileStore struct {
	path string
	mu   sync.RWMutex
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Load(ctx context.Context) (configuration.ActiveConfiguration, error) {
	select {
	case <-ctx.Done():
		return configuration.ActiveConfiguration{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return configuration.ActiveConfiguration{}, ErrNotFound
		}

		return configuration.ActiveConfiguration{}, err
	}

	var current configuration.ActiveConfiguration
	if err := json.Unmarshal(data, &current); err != nil {
		return configuration.ActiveConfiguration{}, err
	}

	return current, nil
}

func (s *FileStore) Save(ctx context.Context, current configuration.ActiveConfiguration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(s.path), "active-config-*.json")
	if err != nil {
		return err
	}

	tempPath := tempFile.Name()

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		_ = os.Remove(tempPath)
		return err
	}

	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	if err := os.Rename(tempPath, s.path); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	return nil
}
