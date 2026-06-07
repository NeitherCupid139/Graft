package systemconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"graft/server/internal/configregistry"
	systemconfigstore "graft/server/modules/system-config/store"
)

var (
	errDefinitionNotFound = errors.New("system config definition not found")
	errInvalidConfigValue = errors.New("invalid system config value")
)

// ValueSnapshot is the service read model for one effective config value.
type ValueSnapshot struct {
	Definition     configregistry.Definition
	EffectiveValue json.RawMessage
	DefaultValue   json.RawMessage
	OverrideValue  json.RawMessage
	HasOverride    bool
	Masked         bool
}

// Service merges module-registered definitions with administrator overrides.
type Service struct {
	registry *configregistry.Registry
	store    systemconfigstore.Repository
}

// NewService creates the system configuration service boundary.
func NewService(registry *configregistry.Registry, store systemconfigstore.Repository) (*Service, error) {
	if registry == nil {
		return nil, errors.New("config registry is unavailable")
	}
	if store == nil {
		return nil, errors.New("system config store is unavailable")
	}
	return &Service{registry: registry, store: store}, nil
}

func (s *Service) List(ctx context.Context) ([]ValueSnapshot, error) {
	definitions := s.registry.Items()
	items := make([]ValueSnapshot, 0, len(definitions))
	for _, definition := range definitions {
		item, err := s.snapshot(ctx, definition)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Service) Get(ctx context.Context, key string) (ValueSnapshot, error) {
	definition, ok := s.registry.Get(key)
	if !ok {
		return ValueSnapshot{}, errDefinitionNotFound
	}
	return s.snapshot(ctx, definition)
}

func (s *Service) Update(ctx context.Context, key string, value json.RawMessage) (ValueSnapshot, error) {
	definition, ok := s.registry.Get(key)
	if !ok {
		return ValueSnapshot{}, errDefinitionNotFound
	}
	if err := validateValueForDefinition(definition, value); err != nil {
		return ValueSnapshot{}, err
	}
	if _, err := s.store.SetOverride(ctx, definition.Key, value); err != nil {
		return ValueSnapshot{}, err
	}
	return s.snapshot(ctx, definition)
}

func (s *Service) Reset(ctx context.Context, key string) (ValueSnapshot, error) {
	definition, ok := s.registry.Get(key)
	if !ok {
		return ValueSnapshot{}, errDefinitionNotFound
	}
	if err := s.store.DeleteOverride(ctx, definition.Key); err != nil {
		return ValueSnapshot{}, err
	}
	return s.snapshot(ctx, definition)
}

func (s *Service) snapshot(ctx context.Context, definition configregistry.Definition) (ValueSnapshot, error) {
	override, err := s.store.GetOverride(ctx, definition.Key)
	hasOverride := true
	if errors.Is(err, systemconfigstore.ErrOverrideNotFound) {
		hasOverride = false
	} else if err != nil {
		return ValueSnapshot{}, err
	}

	effective := definition.DefaultValue
	var overrideValue json.RawMessage
	if hasOverride {
		effective = override.Value
		overrideValue = override.Value
	}

	item := ValueSnapshot{
		Definition:     definition.Snapshot(),
		EffectiveValue: cloneRawMessage(effective),
		DefaultValue:   cloneRawMessage(definition.DefaultValue),
		OverrideValue:  cloneRawMessage(overrideValue),
		HasOverride:    hasOverride,
		Masked:         definition.Sensitive,
	}
	if definition.Sensitive {
		item.EffectiveValue = nil
		item.DefaultValue = nil
		item.OverrideValue = nil
	}
	return item, nil
}

func validateValueForDefinition(definition configregistry.Definition, value json.RawMessage) error {
	if len(value) == 0 {
		return fmt.Errorf("%w: value is required", errInvalidConfigValue)
	}
	var decoded any
	if err := json.Unmarshal(value, &decoded); err != nil {
		return fmt.Errorf("%w: %v", errInvalidConfigValue, err)
	}
	switch definition.Type {
	case configregistry.ValueTypeString:
		if _, ok := decoded.(string); !ok {
			return fmt.Errorf("%w: %s must be string", errInvalidConfigValue, definition.Key)
		}
	case configregistry.ValueTypeNumber:
		if _, ok := decoded.(float64); !ok {
			return fmt.Errorf("%w: %s must be number", errInvalidConfigValue, definition.Key)
		}
	case configregistry.ValueTypeInteger:
		number, ok := decoded.(float64)
		if !ok || number != float64(int64(number)) {
			return fmt.Errorf("%w: %s must be integer", errInvalidConfigValue, definition.Key)
		}
	case configregistry.ValueTypeBoolean:
		if _, ok := decoded.(bool); !ok {
			return fmt.Errorf("%w: %s must be boolean", errInvalidConfigValue, definition.Key)
		}
	case configregistry.ValueTypeObject:
		if _, ok := decoded.(map[string]any); !ok {
			return fmt.Errorf("%w: %s must be object", errInvalidConfigValue, definition.Key)
		}
	case configregistry.ValueTypeArray:
		if _, ok := decoded.([]any); !ok {
			return fmt.Errorf("%w: %s must be array", errInvalidConfigValue, definition.Key)
		}
	default:
		return fmt.Errorf("%w: %s has unsupported type %s", errInvalidConfigValue, definition.Key, definition.Type)
	}
	return nil
}

func cleanKey(key string) string {
	return strings.TrimSpace(key)
}

func cloneRawMessage(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	cloned := make(json.RawMessage, len(raw))
	copy(cloned, raw)
	return cloned
}
