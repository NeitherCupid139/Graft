// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

// Package store defines Announcement Center module persistence contracts.
package store

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrInvalidInput indicates an announcement store input violates the module persistence contract.
	ErrInvalidInput = errors.New("announcement invalid input")
	// ErrAnnouncementNotFound indicates no non-deleted announcement exists for the requested id.
	ErrAnnouncementNotFound = errors.New("announcement not found")
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

// ListQuery describes management-side announcement filters.
type ListQuery struct {
	Status  string
	Level   string
	Pinned  *bool
	Keyword string
	Sort    string
	Limit   int
	Offset  int
}

// ListResult returns a paginated management announcement page.
type ListResult struct {
	Items []Announcement
	Total int
}

// Repository persists announcement records and per-user read facts.
type Repository interface {
	Ping(ctx context.Context) error
	ListAdmin(ctx context.Context, query ListQuery) (ListResult, error)
	Create(ctx context.Context, input CreateInput) (Announcement, error)
	GetAdmin(ctx context.Context, id uint64) (Announcement, error)
	Update(ctx context.Context, id uint64, input UpdateInput) (Announcement, error)
	Publish(ctx context.Context, id uint64, publishAt time.Time, actorID *uint64) (Announcement, error)
	Archive(ctx context.Context, id uint64, actorID *uint64) (Announcement, error)
	Delete(ctx context.Context, id uint64, actorID uint64, deletedAt time.Time) error
}
