package configuration

import "encoding/json"

type GenerateInput struct {
	SampleSchema json.RawMessage `json:"sampleSchema"`
	LabelSchema  json.RawMessage `json:"labelSchema"`
	UIPrompt     string          `json:"uiPrompt"`
}

type SavedConfiguration struct {
	SampleSchema    json.RawMessage `json:"sampleSchema"`
	LabelSchema     json.RawMessage `json:"labelSchema"`
	UIPrompt        string          `json:"uiPrompt"`
	SampleData      json.RawMessage `json:"sampleData"`
	ComponentSource string          `json:"componentSource"`
	UpdatedAt       string          `json:"updatedAt"`
}

