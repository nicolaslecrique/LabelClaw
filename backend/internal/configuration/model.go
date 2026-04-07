package configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ActiveConfiguration struct {
	SampleJSONSchema json.RawMessage `json:"sampleJsonSchema"`
	LabelJSONSchema  json.RawMessage `json:"labelJsonSchema"`
	UIPrompt         string          `json:"uiPrompt"`
	ComponentCode    string          `json:"componentCode"`
}

type GenerateRequest struct {
	SampleJSONSchema json.RawMessage `json:"sampleJsonSchema"`
	LabelJSONSchema  json.RawMessage `json:"labelJsonSchema"`
	UIPrompt         string          `json:"uiPrompt"`
}

var ErrGenerateNotImplemented = errors.New("configuration generation is not implemented yet")

func (c ActiveConfiguration) Validate() error {
	if err := validateSchema("sampleJsonSchema", c.SampleJSONSchema); err != nil {
		return err
	}

	if err := validateSchema("labelJsonSchema", c.LabelJSONSchema); err != nil {
		return err
	}

	if err := validateRequired("uiPrompt", c.UIPrompt); err != nil {
		return err
	}

	if err := validateRequired("componentCode", c.ComponentCode); err != nil {
		return err
	}

	return nil
}

func (r GenerateRequest) Validate() error {
	if err := validateSchema("sampleJsonSchema", r.SampleJSONSchema); err != nil {
		return err
	}

	if err := validateSchema("labelJsonSchema", r.LabelJSONSchema); err != nil {
		return err
	}

	if err := validateRequired("uiPrompt", r.UIPrompt); err != nil {
		return err
	}

	return nil
}

func validateSchema(field string, value json.RawMessage) error {
	if len(bytes.TrimSpace(value)) == 0 {
		return fmt.Errorf("%s is required", field)
	}

	var decoded any
	if err := json.Unmarshal(value, &decoded); err != nil {
		return fmt.Errorf("%s must be valid JSON: %w", field, err)
	}

	if _, ok := decoded.(map[string]any); !ok {
		return fmt.Errorf("%s must be a JSON object", field)
	}

	return nil
}

func validateRequired(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}

	return nil
}
