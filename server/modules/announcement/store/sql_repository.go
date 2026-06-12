// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"database/sql"
	"errors"
)

// SQLRepository persists Announcement Center state in module-owned SQL tables.
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository creates a SQL-backed announcement repository.
func NewSQLRepository(db *sql.DB) (*SQLRepository, error) {
	if db == nil {
		return nil, errors.New("announcement repository requires a non-nil sql db")
	}
	return &SQLRepository{db: db}, nil
}

// Ping verifies the repository can reach its SQL dependency.
func (r *SQLRepository) Ping(ctx context.Context) error {
	if r == nil || r.db == nil {
		return errors.New("announcement repository is unavailable")
	}
	return r.db.PingContext(ctx)
}

var _ Repository = (*SQLRepository)(nil)
