package configuration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	componentSourcePattern = regexp.MustCompile(`(?m)^\s*export\s+default\s+function\s+[A-Za-z_$][\w$]*\s*\(`)
	importPattern          = regexp.MustCompile(`(?m)^\s*import\b`)
)

func ValidateSchemaJSON(raw json.RawMessage) error {
	if len(bytes.TrimSpace(raw)) == 0 {
		return errors.New("schema must not be empty")
	}

	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return fmt.Errorf("schema must be valid JSON: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", bytes.NewReader(raw)); err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	if _, err := compiler.Compile("schema.json"); err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	return nil
}

func ValidateDataAgainstSchema(schemaRaw json.RawMessage, dataRaw json.RawMessage) error {
	if len(bytes.TrimSpace(dataRaw)) == 0 {
		return errors.New("data must not be empty")
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", bytes.NewReader(schemaRaw)); err != nil {
		return fmt.Errorf("load schema: %w", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	var data any
	if err := json.Unmarshal(dataRaw, &data); err != nil {
		return fmt.Errorf("data must be valid JSON: %w", err)
	}

	if err := schema.Validate(data); err != nil {
		return fmt.Errorf("validate data against schema: %w", err)
	}

	return nil
}

func ValidateGeneratedComponentSource(source string) error {
	trimmed := strings.TrimSpace(source)
	if trimmed == "" {
		return errors.New("component source must not be empty")
	}

	if importPattern.MatchString(trimmed) {
		return errors.New("generated component must not contain imports")
	}

	if strings.Contains(trimmed, "require(") {
		return errors.New("generated component must not call require()")
	}

	if !componentSourcePattern.MatchString(trimmed) {
		return errors.New("generated component must use `export default function ComponentName(...)`")
	}

	return nil
}

