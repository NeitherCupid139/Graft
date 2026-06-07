// Package storeent provides SQL persistence for system-config overrides.
package storeent

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	systemconfigstore "graft/server/modules/system-config/store"
)

type repository struct {
	db *sql.DB
}

// NewRepository builds a SQL-backed system config override repository.
func NewRepository(db *sql.DB) (systemconfigstore.Repository, error) {
	if db == nil {
		return nil, errors.New("system config repository requires a non-nil sql db")
	}
	return &repository{db: db}, nil
}

func (r *repository) GetOverride(ctx context.Context, key string) (systemconfigstore.Override, error) {
	if r == nil || r.db == nil {
		return systemconfigstore.Override{}, errors.New("system config repository is unavailable")
	}

	row := r.db.QueryRowContext(
		ctx,
		`SELECT key, override_value, updated_at FROM system_config_values WHERE key = $1`,
		strings.TrimSpace(key),
	)
	override, err := scanOverride(row)
	if errors.Is(err, sql.ErrNoRows) {
		return systemconfigstore.Override{}, systemconfigstore.ErrOverrideNotFound
	}
	if err != nil {
		return systemconfigstore.Override{}, fmt.Errorf("get system config override: %w", err)
	}
	return override, nil
}

func (r *repository) SetOverride(ctx context.Context, key string, value json.RawMessage) (systemconfigstore.Override, error) {
	if r == nil || r.db == nil {
		return systemconfigstore.Override{}, errors.New("system config repository is unavailable")
	}

	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO system_config_values (key, override_value, updated_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (key)
		 DO UPDATE SET override_value = EXCLUDED.override_value, updated_at = NOW()
		 RETURNING key, override_value, updated_at`,
		strings.TrimSpace(key),
		value,
	)
	override, err := scanOverride(row)
	if err != nil {
		return systemconfigstore.Override{}, fmt.Errorf("set system config override: %w", err)
	}
	return override, nil
}

func (r *repository) DeleteOverride(ctx context.Context, key string) error {
	if r == nil || r.db == nil {
		return errors.New("system config repository is unavailable")
	}

	if _, err := r.db.ExecContext(ctx, `DELETE FROM system_config_values WHERE key = $1`, strings.TrimSpace(key)); err != nil {
		return fmt.Errorf("delete system config override: %w", err)
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanOverride(row rowScanner) (systemconfigstore.Override, error) {
	var override systemconfigstore.Override
	var updatedAt time.Time
	if err := row.Scan(&override.Key, &override.Value, &updatedAt); err != nil {
		return systemconfigstore.Override{}, err
	}
	override.UpdatedAt = updatedAt.UTC()
	return override, nil
}
