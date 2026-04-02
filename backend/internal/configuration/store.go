package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrNotFound = errors.New("configuration not found")

type Store interface {
	Load() (SavedConfiguration, error)
	Save(SavedConfiguration) error
}

type FileStore struct {
	path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Load() (SavedConfiguration, error) {
	content, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return SavedConfiguration{}, ErrNotFound
	}
	if err != nil {
		return SavedConfiguration{}, fmt.Errorf("read config: %w", err)
	}

	var config SavedConfiguration
	if err := json.Unmarshal(content, &config); err != nil {
		return SavedConfiguration{}, fmt.Errorf("decode config: %w", err)
	}

	return config, nil
}

func (s *FileStore) Save(config SavedConfiguration) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	content, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	tempFile, err := os.CreateTemp(filepath.Dir(s.path), "active-config-*.json")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	tempPath := tempFile.Name()
	if _, err := tempFile.Write(content); err != nil {
		tempFile.Close()
		_ = os.Remove(tempPath)
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("close temp file: %w", err)
	}

	if err := os.Rename(tempPath, s.path); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("replace config: %w", err)
	}

	return nil
}

