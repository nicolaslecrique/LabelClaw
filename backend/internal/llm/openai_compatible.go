package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAICompatibleClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewOpenAICompatibleClient(config OpenAICompatibleConfig) *OpenAICompatibleClient {
	return &OpenAICompatibleClient{
		baseURL: strings.TrimRight(config.BaseURL, "/"),
		apiKey:  config.APIKey,
		model:   config.Model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *OpenAICompatibleClient) GenerateLabelingPanel(ctx context.Context, request GenerateRequest) (GenerateResponse, error) {
	payload := openAIChatCompletionRequest{
		Model: c.model,
		Messages: []openAIMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role: "user",
				Content: buildUserPrompt(request),
			},
		},
		Temperature: 0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("encode llm request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("create llm request: %w", err)
	}

	httpRequest.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("call llm provider: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("read llm response: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return GenerateResponse{}, fmt.Errorf("llm provider returned %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var completion openAIChatCompletionResponse
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return GenerateResponse{}, fmt.Errorf("decode llm response: %w", err)
	}

	if len(completion.Choices) == 0 {
		return GenerateResponse{}, fmt.Errorf("llm response had no choices")
	}

	content, err := decodeOpenAIMessageContent(completion.Choices[0].Message.Content)
	if err != nil {
		return GenerateResponse{}, err
	}

	var generated GenerateResponse
	if err := json.Unmarshal([]byte(content), &generated); err != nil {
		return GenerateResponse{}, fmt.Errorf("decode generated payload: %w", err)
	}

	return generated, nil
}

type openAIChatCompletionRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func decodeOpenAIMessageContent(raw json.RawMessage) (string, error) {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text, nil
	}

	var chunks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(raw, &chunks); err != nil {
		return "", fmt.Errorf("decode llm message content: %w", err)
	}

	var builder strings.Builder
	for _, chunk := range chunks {
		if chunk.Type == "text" {
			builder.WriteString(chunk.Text)
		}
	}

	if builder.Len() == 0 {
		return "", fmt.Errorf("llm response content was empty")
	}

	return builder.String(), nil
}

func buildUserPrompt(request GenerateRequest) string {
	return fmt.Sprintf(`Create a React component for a dataset labeling interface.

Return only a JSON object with this exact shape:
{
  "componentSource": "string",
  "sampleData": {}
}

Rules for componentSource:
- Must be valid JSX, not Markdown.
- Must use exactly one named default export function: export default function LabelingPanel(...)
- Must not contain imports or require().
- Must use only React and built-in HTML elements.
- Must accept props: { sample, value, onChange }
- Must call onChange(nextValue) with output matching the label schema.
- Keep the UI concise and practical.

Input sample JSON Schema:
%s

Output label JSON Schema:
%s

UI prompt:
%s
`, request.SampleSchema, request.LabelSchema, request.UIPrompt)
}

const systemPrompt = `You generate safe React labeling panels for an existing host application.
Always return a single JSON object and nothing else.
The "componentSource" value must be a string of React component code.
The "sampleData" value must be valid JSON and must satisfy the input sample schema.`

