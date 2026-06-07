// Package configregistry owns module-registered system configuration definitions.
package configregistry

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const maskedPlaceholder = "******"

var keyPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[._-][a-z0-9]+)*$`)

// ValueType identifies the JSON value shape accepted by one config definition.
type ValueType string

const (
	ValueTypeString  ValueType = "string"
	ValueTypeNumber  ValueType = "number"
	ValueTypeInteger ValueType = "integer"
	ValueTypeBoolean ValueType = "boolean"
	ValueTypeObject  ValueType = "object"
	ValueTypeArray   ValueType = "array"
)

// Definition declares one module-owned system configuration key.
//
// Definitions are registered by modules during Register. They are canonical
// metadata and must not be copied into system_config_values as database truth.
type Definition struct {
	Key             string
	Module          string
	Group           string
	Title           string
	TitleKey        string
	Description     string
	DescriptionKey  string
	Type            ValueType
	Schema          json.RawMessage
	DefaultValue    json.RawMessage
	Sensitive       bool
	Required        bool
	RestartRequired bool
	Permission      string
	Order           int
}

// Snapshot returns an immutable copy safe for callers to retain.
func (d Definition) Snapshot() Definition {
	cloned := d
	cloned.Schema = cloneRawMessage(d.Schema)
	cloned.DefaultValue = cloneRawMessage(d.DefaultValue)
	return cloned
}

// MaskedPlaceholder returns the canonical masked value sentinel for sensitive values.
func MaskedPlaceholder() string {
	return maskedPlaceholder
}

func validateDefinition(definition Definition) error {
	key := strings.TrimSpace(definition.Key)
	if key == "" {
		return errors.New("config definition key is required")
	}
	if !keyPattern.MatchString(key) {
		return fmt.Errorf("config definition key %q is invalid", definition.Key)
	}
	if strings.TrimSpace(definition.Module) == "" {
		return fmt.Errorf("config definition %s module is required", key)
	}
	if strings.TrimSpace(definition.Group) == "" {
		return fmt.Errorf("config definition %s group is required", key)
	}
	if strings.TrimSpace(definition.Title) == "" && strings.TrimSpace(definition.TitleKey) == "" {
		return fmt.Errorf("config definition %s title or title key is required", key)
	}
	if !slices.Contains(validValueTypes(), definition.Type) {
		return fmt.Errorf("config definition %s type %q is invalid", key, definition.Type)
	}
	if err := validateJSONObject(definition.Schema, "schema", key); err != nil {
		return err
	}
	if err := validateDefaultValue(definition.DefaultValue, definition.Type, key); err != nil {
		return err
	}
	return nil
}

func validValueTypes() []ValueType {
	return []ValueType{
		ValueTypeString,
		ValueTypeNumber,
		ValueTypeInteger,
		ValueTypeBoolean,
		ValueTypeObject,
		ValueTypeArray,
	}
}

func validateJSONObject(raw json.RawMessage, label string, key string) error {
	if len(raw) == 0 {
		return nil
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return fmt.Errorf("config definition %s %s is invalid JSON: %w", key, label, err)
	}
	if _, ok := decoded.(map[string]any); !ok {
		return fmt.Errorf("config definition %s %s must be a JSON object", key, label)
	}
	return nil
}

func validateDefaultValue(raw json.RawMessage, valueType ValueType, key string) error {
	if len(raw) == 0 {
		return fmt.Errorf("config definition %s default value is required", key)
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return fmt.Errorf("config definition %s default value is invalid JSON: %w", key, err)
	}

	switch valueType {
	case ValueTypeString:
		if _, ok := decoded.(string); !ok {
			return fmt.Errorf("config definition %s default value must be string", key)
		}
	case ValueTypeNumber:
		if _, ok := decoded.(float64); !ok {
			return fmt.Errorf("config definition %s default value must be number", key)
		}
	case ValueTypeInteger:
		number, ok := decoded.(float64)
		if !ok || number != float64(int64(number)) {
			return fmt.Errorf("config definition %s default value must be integer", key)
		}
	case ValueTypeBoolean:
		if _, ok := decoded.(bool); !ok {
			return fmt.Errorf("config definition %s default value must be boolean", key)
		}
	case ValueTypeObject:
		if _, ok := decoded.(map[string]any); !ok {
			return fmt.Errorf("config definition %s default value must be object", key)
		}
	case ValueTypeArray:
		if _, ok := decoded.([]any); !ok {
			return fmt.Errorf("config definition %s default value must be array", key)
		}
	}
	return nil
}

func cloneRawMessage(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	cloned := make(json.RawMessage, len(raw))
	copy(cloned, raw)
	return cloned
}
