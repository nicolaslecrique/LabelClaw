package configuration

import (
	"bytes"
	"encoding/json"
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

type ValidationError struct {
	message string
}

func (e ValidationError) Error() string {
	return e.message
}

func IsValidationError(err error) bool {
	var validationError ValidationError

	return err != nil && AsValidationError(err, &validationError)
}

func AsValidationError(err error, target *ValidationError) bool {
	if err == nil {
		return false
	}

	validationError, ok := err.(ValidationError)
	if !ok {
		return false
	}

	if target != nil {
		*target = validationError
	}

	return true
}

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
		return ValidationError{message: fmt.Sprintf("%s is required", field)}
	}

	var decoded any
	if err := json.Unmarshal(value, &decoded); err != nil {
		return ValidationError{message: fmt.Sprintf("%s must be valid JSON: %v", field, err)}
	}

	if _, ok := decoded.(map[string]any); !ok {
		return ValidationError{message: fmt.Sprintf("%s must be a JSON object", field)}
	}

	return nil
}

func validateRequired(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return ValidationError{message: fmt.Sprintf("%s is required", field)}
	}

	return nil
}
