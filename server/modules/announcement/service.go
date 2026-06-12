// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"context"
	"errors"
	"time"

	announcementstore "graft/server/modules/announcement/store"
)

var (
	errAnnouncementNotImplemented = errors.New("announcement service behavior is not implemented")
)

// AdminListQuery describes management-side announcement filters.
type AdminListQuery struct {
	Status   string
	Level    string
	Pinned   *bool
	Keyword  string
	Page     int
	PageSize int
	Sort     string
}

// UserListQuery describes current-user announcement filters.
type UserListQuery struct {
	UserID     uint64
	UnreadOnly bool
	Page       int
	PageSize   int
}

// Service owns announcement use cases.
type Service struct {
	repository announcementstore.Repository
}

// NewService creates the announcement service boundary.
func NewService(repository announcementstore.Repository) (*Service, error) {
	if repository == nil {
		return nil, errors.New("announcement repository is unavailable")
	}
	return &Service{repository: repository}, nil
}

// ListAdmin is intentionally unimplemented in Phase 1.
func (s *Service) ListAdmin(context.Context, AdminListQuery) error {
	return s.phaseOneNotImplemented()
}

// Create is intentionally unimplemented in Phase 1.
func (s *Service) Create(context.Context, announcementstore.CreateInput) error {
	return s.phaseOneNotImplemented()
}

// GetAdmin is intentionally unimplemented in Phase 1.
func (s *Service) GetAdmin(context.Context, uint64) error {
	return s.phaseOneNotImplemented()
}

// Update is intentionally unimplemented in Phase 1.
func (s *Service) Update(context.Context, uint64, announcementstore.UpdateInput) error {
	return s.phaseOneNotImplemented()
}

// Publish is intentionally unimplemented in Phase 1.
func (s *Service) Publish(context.Context, uint64, *time.Time) error {
	return s.phaseOneNotImplemented()
}

// Archive is intentionally unimplemented in Phase 1.
func (s *Service) Archive(context.Context, uint64) error {
	return s.phaseOneNotImplemented()
}

// Delete is intentionally unimplemented in Phase 1.
func (s *Service) Delete(context.Context, uint64, uint64) error {
	return s.phaseOneNotImplemented()
}

// ListCurrentUser is intentionally unimplemented in Phase 1.
func (s *Service) ListCurrentUser(context.Context, UserListQuery) error {
	return s.phaseOneNotImplemented()
}

// MarkRead is intentionally unimplemented in Phase 1.
func (s *Service) MarkRead(context.Context, uint64, uint64) error {
	return s.phaseOneNotImplemented()
}

// MarkAllRead is intentionally unimplemented in Phase 1.
func (s *Service) MarkAllRead(context.Context, uint64) error {
	return s.phaseOneNotImplemented()
}

// UnreadCount is intentionally unimplemented in Phase 1.
func (s *Service) UnreadCount(context.Context, uint64) error {
	return s.phaseOneNotImplemented()
}

func (s *Service) phaseOneNotImplemented() error {
	if s == nil || s.repository == nil {
		return errors.New("announcement service is unavailable")
	}
	return errAnnouncementNotImplemented
}
