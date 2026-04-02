package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAICompatibleClientGenerateLabelingPanel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("unexpected auth header: %s", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"choices": [
				{
					"message": {
						"content": "{\"componentSource\":\"export default function LabelingPanel() { return <div />; }\",\"sampleData\":{\"title\":\"Hello\"}}"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewOpenAICompatibleClient(OpenAICompatibleConfig{
		BaseURL: server.URL,
		APIKey:  "secret",
		Model:   "test-model",
	})

	response, err := client.GenerateLabelingPanel(context.Background(), GenerateRequest{
		SampleSchema: json.RawMessage(`{"type":"object"}`),
		LabelSchema:  json.RawMessage(`{"type":"string"}`),
		UIPrompt:     "Render a compact UI",
	})
	if err != nil {
		t.Fatalf("generate panel: %v", err)
	}

	if response.ComponentSource == "" {
		t.Fatalf("expected component source")
	}

	if string(response.SampleData) != `{"title":"Hello"}` {
		t.Fatalf("unexpected sampleData: %s", response.SampleData)
	}
}

