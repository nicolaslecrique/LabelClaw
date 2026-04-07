package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/nicolaslecrique/LabelClaw/backend/internal/configuration"
	"github.com/nicolaslecrique/LabelClaw/backend/internal/storage"
)

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newTestHandler(t))
	defer server.Close()

	response, err := http.Get(server.URL + "/api/health")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", response.StatusCode, http.StatusOK)
	}
}

func TestGetCurrentConfigurationReturnsNotFoundWhenMissing(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newTestHandler(t))
	defer server.Close()

	response, err := http.Get(server.URL + "/api/configuration/current")
	if err != nil {
		t.Fatalf("get current request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d want %d", response.StatusCode, http.StatusNotFound)
	}
}

func TestPutCurrentConfigurationPersistsConfiguration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newTestHandler(t))
	defer server.Close()

	current := configuration.ActiveConfiguration{
		SampleJSONSchema: json.RawMessage(`{"type":"object","properties":{"input":{"type":"string"}}}`),
		LabelJSONSchema:  json.RawMessage(`{"type":"object","properties":{"label":{"type":"string"}}}`),
		UIPrompt:         "Render a clean classification form.",
		ComponentCode:    "export function Panel() { return null; }",
	}

	payload, err := json.Marshal(current)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	request, err := http.NewRequest(http.MethodPut, server.URL+"/api/configuration/current", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("build request failed: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("put request failed: %v", err)
	}
	response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected put status: got %d want %d", response.StatusCode, http.StatusNoContent)
	}

	getResponse, err := http.Get(server.URL + "/api/configuration/current")
	if err != nil {
		t.Fatalf("follow-up get request failed: %v", err)
	}
	defer getResponse.Body.Close()

	if getResponse.StatusCode != http.StatusOK {
		t.Fatalf("unexpected get status: got %d want %d", getResponse.StatusCode, http.StatusOK)
	}

	var saved configuration.ActiveConfiguration
	if err := json.NewDecoder(getResponse.Body).Decode(&saved); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if string(saved.SampleJSONSchema) != string(current.SampleJSONSchema) {
		t.Fatalf("unexpected sample schema: got %s want %s", saved.SampleJSONSchema, current.SampleJSONSchema)
	}

	if string(saved.LabelJSONSchema) != string(current.LabelJSONSchema) {
		t.Fatalf("unexpected label schema: got %s want %s", saved.LabelJSONSchema, current.LabelJSONSchema)
	}

	if saved.UIPrompt != current.UIPrompt {
		t.Fatalf("unexpected ui prompt: got %q want %q", saved.UIPrompt, current.UIPrompt)
	}

	if saved.ComponentCode != current.ComponentCode {
		t.Fatalf("unexpected component code: got %q want %q", saved.ComponentCode, current.ComponentCode)
	}
}

func TestGenerateConfigurationReturnsNotImplemented(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(newTestHandler(t))
	defer server.Close()

	requestBody := configuration.GenerateRequest{
		SampleJSONSchema: json.RawMessage(`{"type":"object","properties":{"input":{"type":"string"}}}`),
		LabelJSONSchema:  json.RawMessage(`{"type":"object","properties":{"label":{"type":"string"}}}`),
		UIPrompt:         "Generate a panel with radio buttons.",
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	response, err := http.Post(server.URL+"/api/configuration/generate", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("generate request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotImplemented {
		t.Fatalf("unexpected status: got %d want %d", response.StatusCode, http.StatusNotImplemented)
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	store := storage.NewFileStore(filepath.Join(t.TempDir(), "active-config.json"))
	service := configuration.NewService(store)

	return NewHandler(service)
}
