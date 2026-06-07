// Package store defines system-config module persistence contracts.
package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// ErrOverrideNotFound indicates that no administrator override exists for the key.
var ErrOverrideNotFound = errors.New("system config override not found")

// Override stores administrator-provided JSON for one config key.
type Override struct {
	Key       string
	Value     json.RawMessage
	UpdatedAt time.Time
}

// Repository persists administrator overrides only.
type Repository interface {
	GetOverride(ctx context.Context, key string) (Override, error)
	SetOverride(ctx context.Context, key string, value json.RawMessage) (Override, error)
	DeleteOverride(ctx context.Context, key string) error
}
