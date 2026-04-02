package configuration

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
)

func TestFileStoreLoadMissing(t *testing.T) {
	store := NewFileStore(filepath.Join(t.TempDir(), "missing.json"))

	_, err := store.Load()
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFileStoreSaveAndLoad(t *testing.T) {
	store := NewFileStore(filepath.Join(t.TempDir(), "active-config.json"))
	expected := SavedConfiguration{
		SampleSchema:    json.RawMessage(`{"type":"object"}`),
		LabelSchema:     json.RawMessage(`{"type":"string"}`),
		UIPrompt:        "Render a textarea",
		SampleData:      json.RawMessage(`{"text":"hello"}`),
		ComponentSource: "export default function Panel() { return <div /> }",
		UpdatedAt:       "2026-04-02T10:30:00Z",
	}

	if err := store.Save(expected); err != nil {
		t.Fatalf("save: %v", err)
	}

	actual, err := store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if string(actual.SampleSchema) != string(expected.SampleSchema) {
		t.Fatalf("unexpected sample schema: %s", actual.SampleSchema)
	}
	if string(actual.LabelSchema) != string(expected.LabelSchema) {
		t.Fatalf("unexpected label schema: %s", actual.LabelSchema)
	}
	if string(actual.SampleData) != string(expected.SampleData) {
		t.Fatalf("unexpected sample data: %s", actual.SampleData)
	}
	if actual.ComponentSource != expected.ComponentSource {
		t.Fatalf("unexpected component source: %s", actual.ComponentSource)
	}
	if actual.UpdatedAt != expected.UpdatedAt {
		t.Fatalf("unexpected updatedAt: %s", actual.UpdatedAt)
	}
}

