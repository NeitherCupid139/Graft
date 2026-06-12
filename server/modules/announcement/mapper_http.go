// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"math"

	announcementstore "graft/server/modules/announcement/store"
)

func toAnnouncementListResponse(result AdminListResult) map[string]any {
	items := make([]map[string]any, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toAnnouncementItem(item))
	}
	return map[string]any{
		"items":     items,
		"total":     result.Total,
		"page":      result.Page,
		"page_size": result.PageSize,
	}
}

func toAnnouncementItem(item announcementstore.Announcement) map[string]any {
	return map[string]any{
		"id":         safeInt64(item.ID),
		"title":      item.Title,
		"content":    item.Content,
		"level":      item.Level,
		"status":     item.Status,
		"pinned":     item.Pinned,
		"publish_at": item.PublishAt,
		"expire_at":  item.ExpireAt,
		"created_by": safeOptionalInt64(item.CreatedBy),
		"updated_by": safeOptionalInt64(item.UpdatedBy),
		"deleted_by": safeOptionalInt64(item.DeletedBy),
		"created_at": item.CreatedAt,
		"updated_at": item.UpdatedAt,
	}
}

func safeInt64(value uint64) int64 {
	if value > uint64(math.MaxInt64) {
		return math.MaxInt64
	}
	return int64(value)
}

func safeOptionalInt64(value *uint64) *int64 {
	if value == nil {
		return nil
	}
	converted := safeInt64(*value)
	return &converted
}
