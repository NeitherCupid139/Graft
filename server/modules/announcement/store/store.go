// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

// Package store defines Announcement Center module persistence contracts.
package store

import (
	"context"
	"time"
)

// Announcement stores the management announcement record.
type Announcement struct {
	ID        uint64
	Title     string
	Content   string
	Level     string
	Status    string
	Pinned    bool
	PublishAt *time.Time
	ExpireAt  *time.Time
	CreatedBy *uint64
	UpdatedBy *uint64
	DeletedBy *uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt int64
}

// AnnouncementRead stores one user's read fact for one announcement.
type AnnouncementRead struct {
	ID             uint64
	AnnouncementID uint64
	UserID         uint64
	ReadAt         time.Time
	CreatedAt      time.Time
}

// CreateInput describes one announcement draft insert.
type CreateInput struct {
	Title     string
	Content   string
	Level     string
	Status    string
	Pinned    bool
	PublishAt *time.Time
	ExpireAt  *time.Time
	ActorID   *uint64
}

// UpdateInput describes editable announcement fields.
type UpdateInput struct {
	Title     string
	Content   string
	Level     string
	Pinned    bool
	PublishAt *time.Time
	ExpireAt  *time.Time
	ActorID   *uint64
}

// Repository persists announcement records and per-user read facts.
type Repository interface {
	Ping(ctx context.Context) error
}
