// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	announcementstore "graft/server/modules/announcement/store"
)

type memoryAnnouncementRepository struct {
	mu     sync.Mutex
	nextID uint64
	items  map[uint64]announcementstore.Announcement
}

func newMemoryAnnouncementRepository() *memoryAnnouncementRepository {
	return &memoryAnnouncementRepository{
		nextID: 1,
		items:  make(map[uint64]announcementstore.Announcement),
	}
}

func (r *memoryAnnouncementRepository) Ping(context.Context) error {
	return nil
}

func (r *memoryAnnouncementRepository) ListAdmin(_ context.Context, query announcementstore.ListQuery) (announcementstore.ListResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	items := make([]announcementstore.Announcement, 0, len(r.items))
	for _, item := range r.items {
		if item.DeletedAt != 0 {
			continue
		}
		if query.Status != "" && item.Status != query.Status {
			continue
		}
		if query.Level != "" && item.Level != query.Level {
			continue
		}
		if query.Pinned != nil && item.Pinned != *query.Pinned {
			continue
		}
		if query.Keyword != "" && !strings.Contains(strings.ToLower(item.Title+" "+item.Content), strings.ToLower(query.Keyword)) {
			continue
		}
		items = append(items, item)
	}
	sort.SliceStable(items, func(i int, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt) || items[i].ID > items[j].ID
	})
	total := len(items)
	start := query.Offset
	if start > total {
		start = total
	}
	end := start + query.Limit
	if end > total {
		end = total
	}
	return announcementstore.ListResult{Items: append([]announcementstore.Announcement(nil), items[start:end]...), Total: total}, nil
}

func (r *memoryAnnouncementRepository) Create(_ context.Context, input announcementstore.CreateInput) (announcementstore.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	item := announcementstore.Announcement{
		ID:        r.nextID,
		Title:     strings.TrimSpace(input.Title),
		Content:   strings.TrimSpace(input.Content),
		Level:     strings.TrimSpace(input.Level),
		Status:    strings.TrimSpace(input.Status),
		Pinned:    input.Pinned,
		PublishAt: cloneTime(input.PublishAt),
		ExpireAt:  cloneTime(input.ExpireAt),
		CreatedBy: cloneUint64(input.ActorID),
		UpdatedBy: cloneUint64(input.ActorID),
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.nextID++
	r.items[item.ID] = item
	return item, nil
}

func (r *memoryAnnouncementRepository) GetAdmin(_ context.Context, id uint64) (announcementstore.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok || item.DeletedAt != 0 {
		return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
	}
	return item, nil
}

func (r *memoryAnnouncementRepository) Update(_ context.Context, id uint64, input announcementstore.UpdateInput) (announcementstore.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok || item.DeletedAt != 0 {
		return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
	}
	item.Title = strings.TrimSpace(input.Title)
	item.Content = strings.TrimSpace(input.Content)
	item.Level = strings.TrimSpace(input.Level)
	item.Pinned = input.Pinned
	item.PublishAt = cloneTime(input.PublishAt)
	item.ExpireAt = cloneTime(input.ExpireAt)
	item.UpdatedBy = cloneUint64(input.ActorID)
	item.UpdatedAt = time.Now().UTC()
	r.items[id] = item
	return item, nil
}

func (r *memoryAnnouncementRepository) Publish(_ context.Context, id uint64, publishAt time.Time, actorID *uint64) (announcementstore.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok || item.DeletedAt != 0 {
		return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
	}
	item.Status = "published"
	if item.PublishAt == nil {
		publishAt = publishAt.UTC()
		item.PublishAt = &publishAt
	}
	item.UpdatedBy = cloneUint64(actorID)
	item.UpdatedAt = time.Now().UTC()
	r.items[id] = item
	return item, nil
}

func (r *memoryAnnouncementRepository) Archive(_ context.Context, id uint64, actorID *uint64) (announcementstore.Announcement, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok || item.DeletedAt != 0 {
		return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
	}
	item.Status = "archived"
	item.UpdatedBy = cloneUint64(actorID)
	item.UpdatedAt = time.Now().UTC()
	r.items[id] = item
	return item, nil
}

func (r *memoryAnnouncementRepository) Delete(_ context.Context, id uint64, actorID uint64, deletedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok || item.DeletedAt != 0 {
		return announcementstore.ErrAnnouncementNotFound
	}
	item.DeletedBy = cloneUint64(&actorID)
	item.DeletedAt = deletedAt.UTC().Unix()
	item.UpdatedBy = cloneUint64(&actorID)
	item.UpdatedAt = deletedAt.UTC()
	r.items[id] = item
	return nil
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := value.UTC()
	return &cloned
}

func cloneUint64(value *uint64) *uint64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

var _ announcementstore.Repository = (*memoryAnnouncementRepository)(nil)
