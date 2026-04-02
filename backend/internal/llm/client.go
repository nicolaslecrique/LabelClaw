package llm

import (
	"context"
	"encoding/json"
	"errors"
	"os"
)

var ErrProviderUnavailable = errors.New("llm provider unavailable")

type Client interface {
	GenerateLabelingPanel(ctx context.Context, request GenerateRequest) (GenerateResponse, error)
}

type GenerateRequest struct {
	SampleSchema json.RawMessage
	LabelSchema  json.RawMessage
	UIPrompt     string
}

type GenerateResponse struct {
	ComponentSource string          `json:"componentSource"`
	SampleData      json.RawMessage `json:"sampleData"`
}

type OpenAICompatibleConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

func LoadOpenAICompatibleConfigFromEnv() (OpenAICompatibleConfig, bool) {
	config := OpenAICompatibleConfig{
		BaseURL: os.Getenv("LLM_BASE_URL"),
		APIKey:  os.Getenv("LLM_API_KEY"),
		Model:   os.Getenv("LLM_MODEL"),
	}

	return config, config.BaseURL != "" && config.APIKey != "" && config.Model != ""
}

type UnavailableClient struct {
	reason string
}

func NewUnavailableClient(reason string) UnavailableClient {
	return UnavailableClient{reason: reason}
}

func (c UnavailableClient) GenerateLabelingPanel(_ context.Context, _ GenerateRequest) (GenerateResponse, error) {
	return GenerateResponse{}, errors.Join(ErrProviderUnavailable, errors.New(c.reason))
}

