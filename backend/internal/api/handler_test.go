package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicolas/labelclaw/backend/internal/configuration"
	"github.com/nicolas/labelclaw/backend/internal/llm"
)

type stubGenerator struct {
	response llm.GenerateResponse
	err      error
}

func (s stubGenerator) GenerateLabelingPanel(_ context.Context, _ llm.GenerateRequest) (llm.GenerateResponse, error) {
	return s.response, s.err
}

func TestGenerateConfigurationRejectsInvalidSchema(t *testing.T) {
	store := configuration.NewFileStore(t.TempDir() + "/active.json")
	handler := NewServer(Dependencies{
		Store:         store,
		Generator:     stubGenerator{},
		StaticHandler: http.NotFoundHandler(),
		AllowedOrigin: "*",
	})

	request := httptest.NewRequest(http.MethodPost, "/api/configuration/generate", bytes.NewBufferString(`{
		"sampleSchema": {"type":"object"},
		"labelSchema": {"type": },
		"uiPrompt": "Render a form"
	}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
}

func TestGenerateConfigurationRejectsInvalidProviderOutput(t *testing.T) {
	store := configuration.NewFileStore(t.TempDir() + "/active.json")
	handler := NewServer(Dependencies{
		Store: store,
		Generator: stubGenerator{
			response: llm.GenerateResponse{
				ComponentSource: "export default function LabelingPanel() { return <div />; }",
				SampleData:      json.RawMessage(`{"title":123}`),
			},
		},
		StaticHandler: http.NotFoundHandler(),
		AllowedOrigin: "*",
	})

	request := httptest.NewRequest(http.MethodPost, "/api/configuration/generate", bytes.NewBufferString(`{
		"sampleSchema": {"type":"object","properties":{"title":{"type":"string"}},"required":["title"]},
		"labelSchema": {"type":"string"},
		"uiPrompt": "Render a title"
	}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", response.Code)
	}
}

func TestSaveAndLoadConfiguration(t *testing.T) {
	store := configuration.NewFileStore(t.TempDir() + "/active.json")
	handler := NewServer(Dependencies{
		Store:         store,
		Generator:     stubGenerator{},
		StaticHandler: http.NotFoundHandler(),
		AllowedOrigin: "*",
	})

	savePayload := `{
		"sampleSchema": {"type":"object","properties":{"title":{"type":"string"}},"required":["title"]},
		"labelSchema": {"type":"string"},
		"uiPrompt": "Render a textarea",
		"sampleData": {"title":"Hello"},
		"componentSource": "export default function LabelingPanel({ sample, value, onChange }) { return <textarea aria-label='Label output' value={value ?? ''} onChange={(event) => onChange(event.target.value)} />; }"
	}`
	saveRequest := httptest.NewRequest(http.MethodPut, "/api/configuration/current", bytes.NewBufferString(savePayload))
	saveResponse := httptest.NewRecorder()

	handler.ServeHTTP(saveResponse, saveRequest)

	if saveResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", saveResponse.Code)
	}

	loadRequest := httptest.NewRequest(http.MethodGet, "/api/configuration/current", nil)
	loadResponse := httptest.NewRecorder()

	handler.ServeHTTP(loadResponse, loadRequest)

	if loadResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", loadResponse.Code)
	}

	var config configuration.SavedConfiguration
	if err := json.Unmarshal(loadResponse.Body.Bytes(), &config); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if config.UpdatedAt == "" {
		t.Fatalf("expected updatedAt to be set")
	}
}

func TestGenerateConfigurationReturnsServiceUnavailableWhenProviderIsMissing(t *testing.T) {
	store := configuration.NewFileStore(t.TempDir() + "/active.json")
	handler := NewServer(Dependencies{
		Store:         store,
		Generator:     stubGenerator{err: errors.Join(llm.ErrProviderUnavailable, errors.New("missing env"))},
		StaticHandler: http.NotFoundHandler(),
		AllowedOrigin: "*",
	})

	request := httptest.NewRequest(http.MethodPost, "/api/configuration/generate", bytes.NewBufferString(`{
		"sampleSchema": {"type":"object"},
		"labelSchema": {"type":"string"},
		"uiPrompt": "Render a form"
	}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.Code)
	}
}
