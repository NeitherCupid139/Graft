package systemconfig

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"graft/server/internal/configregistry"
	systemconfigstore "graft/server/modules/system-config/store"
)

func TestServiceMasksSensitiveValues(t *testing.T) {
	service := newTestService(t, configregistry.Definition{
		Key:          "auth.jwt_secret",
		Module:       "auth",
		Group:        "security",
		Title:        "JWT Secret",
		Type:         configregistry.ValueTypeString,
		DefaultValue: json.RawMessage(`"secret"`),
		Sensitive:    true,
	})

	item, err := service.Get(context.Background(), "auth.jwt_secret")
	if err != nil {
		t.Fatalf("get sensitive config: %v", err)
	}
	if !item.Masked || item.EffectiveValue != nil || item.DefaultValue != nil {
		t.Fatalf("sensitive values must be masked, got %#v", item)
	}
}

func TestServiceStoresOverridesOnly(t *testing.T) {
	repo := newMemoryRepo()
	service := newTestServiceWithRepo(t, repo, configregistry.Definition{
		Key:          "scheduler.timeout",
		Module:       "scheduler",
		Group:        "runtime",
		Title:        "Timeout",
		Type:         configregistry.ValueTypeString,
		DefaultValue: json.RawMessage(`"30s"`),
	})

	item, err := service.Get(context.Background(), "scheduler.timeout")
	if err != nil {
		t.Fatalf("get default config: %v", err)
	}
	if item.HasOverride {
		t.Fatal("default-only config must not imply persisted override")
	}

	item, err = service.Update(context.Background(), "scheduler.timeout", json.RawMessage(`"60s"`))
	if err != nil {
		t.Fatalf("update override: %v", err)
	}
	if !item.HasOverride || string(item.EffectiveValue) != `"60s"` {
		t.Fatalf("expected effective override, got %#v", item)
	}
	if len(repo.values) != 1 {
		t.Fatalf("expected only one override row, got %d", len(repo.values))
	}
}

func TestServiceRejectsMismatchedValueType(t *testing.T) {
	service := newTestService(t, configregistry.Definition{
		Key:          "audit.retention_days",
		Module:       "audit",
		Group:        "retention",
		Title:        "Retention Days",
		Type:         configregistry.ValueTypeInteger,
		DefaultValue: json.RawMessage(`30`),
	})

	if _, err := service.Update(context.Background(), "audit.retention_days", json.RawMessage(`"30"`)); err == nil {
		t.Fatal("expected value type error")
	}
}

func newTestService(t *testing.T, definition configregistry.Definition) *Service {
	t.Helper()
	return newTestServiceWithRepo(t, newMemoryRepo(), definition)
}

func newTestServiceWithRepo(t *testing.T, repo *memoryRepo, definition configregistry.Definition) *Service {
	t.Helper()
	registry := configregistry.NewRegistry()
	if err := registry.Register(definition); err != nil {
		t.Fatalf("register definition: %v", err)
	}
	service, err := NewService(registry, repo)
	if err != nil {
		t.Fatalf("create service: %v", err)
	}
	return service
}

type memoryRepo struct {
	values map[string]json.RawMessage
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{values: make(map[string]json.RawMessage)}
}

func (r *memoryRepo) GetOverride(_ context.Context, key string) (systemconfigstore.Override, error) {
	value, ok := r.values[key]
	if !ok {
		return systemconfigstore.Override{}, systemconfigstore.ErrOverrideNotFound
	}
	return systemconfigstore.Override{Key: key, Value: cloneRawMessage(value)}, nil
}

func (r *memoryRepo) SetOverride(_ context.Context, key string, value json.RawMessage) (systemconfigstore.Override, error) {
	if len(value) == 0 {
		return systemconfigstore.Override{}, errors.New("value is required")
	}
	r.values[key] = cloneRawMessage(value)
	return systemconfigstore.Override{Key: key, Value: cloneRawMessage(value)}, nil
}

func (r *memoryRepo) DeleteOverride(_ context.Context, key string) error {
	delete(r.values, key)
	return nil
}
