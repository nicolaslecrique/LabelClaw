package storage

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/nicolaslecrique/LabelClaw/backend/internal/configuration"
)

func TestFileStoreLoadReturnsNotFoundWhenFileDoesNotExist(t *testing.T) {
	t.Parallel()

	store := NewFileStore(filepath.Join(t.TempDir(), "active-config.json"))

	_, err := store.Load()
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFileStoreSaveAndLoadRoundTrip(t *testing.T) {
	t.Parallel()

	store := NewFileStore(filepath.Join(t.TempDir(), "active-config.json"))
	expected := configuration.ActiveConfiguration{
		SampleJSONSchema: json.RawMessage(`{"type":"object","properties":{"input":{"type":"string"}}}`),
		LabelJSONSchema:  json.RawMessage(`{"type":"object","properties":{"label":{"type":"string"}}}`),
		UIPrompt:         "Render a simple text labeling component.",
		ComponentCode:    "export function Panel() { return null; }",
	}

	if err := store.Save(expected); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	actual, err := store.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	assertJSONEqual(t, actual.SampleJSONSchema, expected.SampleJSONSchema)
	assertJSONEqual(t, actual.LabelJSONSchema, expected.LabelJSONSchema)

	if actual.UIPrompt != expected.UIPrompt {
		t.Fatalf("unexpected ui prompt: got %q want %q", actual.UIPrompt, expected.UIPrompt)
	}

	if actual.ComponentCode != expected.ComponentCode {
		t.Fatalf("unexpected component code: got %q want %q", actual.ComponentCode, expected.ComponentCode)
	}
}

func assertJSONEqual(t *testing.T, actual json.RawMessage, expected json.RawMessage) {
	t.Helper()

	var actualValue any
	if err := json.Unmarshal(actual, &actualValue); err != nil {
		t.Fatalf("failed to unmarshal actual JSON: %v", err)
	}

	var expectedValue any
	if err := json.Unmarshal(expected, &expectedValue); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %v", err)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Fatalf("unexpected JSON value: got %s want %s", actual, expected)
	}
}
