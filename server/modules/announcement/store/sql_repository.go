// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	defaultListLimit          = 20
	maxListLimit              = 100
	placeholderGrowthEstimate = 8
	adminFilterCapacity       = 4
	sortUpdatedDesc           = "updated_desc"
	sortPublishDesc           = "publish_desc"
	sortPinnedPublishDesc     = "pinned_publish_desc"
	statusPublished           = "published"
	statusArchived            = "archived"
)

// SQLRepository persists Announcement Center state in module-owned SQL tables.
type SQLRepository struct {
	db          *sql.DB
	placeholder placeholderStyle
}

// NewSQLRepository creates a SQL-backed announcement repository.
func NewSQLRepository(db *sql.DB) (*SQLRepository, error) {
	if db == nil {
		return nil, errors.New("announcement repository requires a non-nil sql db")
	}
	return &SQLRepository{db: db, placeholder: detectPlaceholderStyle(db)}, nil
}

// Ping verifies the repository can reach its SQL dependency.
func (r *SQLRepository) Ping(ctx context.Context) error {
	if err := r.ensureReady(); err != nil {
		return err
	}
	return r.db.PingContext(ctx)
}

// ListAdmin returns non-deleted announcements for the management list.
func (r *SQLRepository) ListAdmin(ctx context.Context, query ListQuery) (ListResult, error) {
	if err := r.ensureReady(); err != nil {
		return ListResult{}, err
	}
	query = normalizeListQuery(query)
	where, args, err := buildAdminWhere(query)
	if err != nil {
		return ListResult{}, err
	}

	//nolint:gosec // Predicates and ordering come from fixed fragments; values stay parameterized.
	countSQL := r.placeholder.rebind(fmt.Sprintf(`SELECT COUNT(*) FROM announcements WHERE %s`, strings.Join(where, " AND ")))
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return ListResult{}, fmt.Errorf("count announcements: %w", err)
	}

	args = append(args, query.Limit, query.Offset)
	//nolint:gosec // Predicates and ordering come from fixed fragments; values stay parameterized.
	rows, err := r.db.QueryContext(ctx, r.placeholder.rebind(fmt.Sprintf(`SELECT %s
		FROM announcements
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?`, announcementColumns(), strings.Join(where, " AND "), adminOrderBy(query.Sort))), args...)
	if err != nil {
		return ListResult{}, fmt.Errorf("list announcements: %w", err)
	}
	defer closeRows(rows)

	items, err := scanAnnouncements(rows)
	if err != nil {
		return ListResult{}, err
	}
	return ListResult{Items: items, Total: total}, nil
}

// Create inserts one management announcement.
func (r *SQLRepository) Create(ctx context.Context, input CreateInput) (Announcement, error) {
	if err := r.ensureReady(); err != nil {
		return Announcement{}, err
	}
	input = normalizeCreateInput(input)
	if input.Title == "" || input.Content == "" || input.Level == "" || input.Status == "" {
		return Announcement{}, ErrInvalidInput
	}
	if input.ExpireAt != nil && input.PublishAt != nil && !input.ExpireAt.After(*input.PublishAt) {
		return Announcement{}, ErrInvalidInput
	}
	now := time.Now().UTC()
	return scanAnnouncement(r.db.QueryRowContext(ctx, r.placeholder.rebind(`INSERT INTO announcements (
			title, content, level, status, pinned, publish_at, expire_at, created_by, updated_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING `+announcementColumns()),
		input.Title,
		input.Content,
		input.Level,
		input.Status,
		input.Pinned,
		input.PublishAt,
		input.ExpireAt,
		input.ActorID,
		input.ActorID,
		now,
		now,
	))
}

// GetAdmin returns one non-deleted announcement by id.
func (r *SQLRepository) GetAdmin(ctx context.Context, id uint64) (Announcement, error) {
	if err := r.ensureReady(); err != nil {
		return Announcement{}, err
	}
	targetID, err := toDBID(id)
	if err != nil {
		return Announcement{}, err
	}
	item, err := scanAnnouncement(r.db.QueryRowContext(ctx, r.placeholder.rebind(`SELECT `+announcementColumns()+`
		FROM announcements
		WHERE id = ? AND deleted_at = 0`), targetID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Announcement{}, ErrAnnouncementNotFound
		}
		return Announcement{}, fmt.Errorf("get announcement: %w", err)
	}
	return item, nil
}

// Update replaces editable management fields for one non-deleted announcement.
func (r *SQLRepository) Update(ctx context.Context, id uint64, input UpdateInput) (Announcement, error) {
	if err := r.ensureReady(); err != nil {
		return Announcement{}, err
	}
	targetID, err := toDBID(id)
	if err != nil {
		return Announcement{}, err
	}
	input = normalizeUpdateInput(input)
	if input.Title == "" || input.Content == "" || input.Level == "" {
		return Announcement{}, ErrInvalidInput
	}
	if input.ExpireAt != nil && input.PublishAt != nil && !input.ExpireAt.After(*input.PublishAt) {
		return Announcement{}, ErrInvalidInput
	}
	item, err := r.updateAnnouncement(ctx, targetID, input)
	if err != nil {
		return Announcement{}, err
	}
	return item, nil
}

func (r *SQLRepository) updateAnnouncement(ctx context.Context, targetID int64, input UpdateInput) (Announcement, error) {
	item, err := scanAnnouncement(r.db.QueryRowContext(ctx, r.placeholder.rebind(`UPDATE announcements
		SET title = ?, content = ?, level = ?, pinned = ?, publish_at = ?, expire_at = ?, updated_by = ?, updated_at = ?
		WHERE id = ? AND deleted_at = 0
		RETURNING `+announcementColumns()),
		input.Title,
		input.Content,
		input.Level,
		input.Pinned,
		input.PublishAt,
		input.ExpireAt,
		input.ActorID,
		time.Now().UTC(),
		targetID,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Announcement{}, ErrAnnouncementNotFound
		}
		return Announcement{}, fmt.Errorf("update announcement: %w", err)
	}
	return item, nil
}

// Publish marks one announcement published and records an effective publish time.
func (r *SQLRepository) Publish(ctx context.Context, id uint64, publishAt time.Time, actorID *uint64) (Announcement, error) {
	if err := r.ensureReady(); err != nil {
		return Announcement{}, err
	}
	targetID, err := toDBID(id)
	if err != nil {
		return Announcement{}, err
	}
	if publishAt.IsZero() {
		return Announcement{}, ErrInvalidInput
	}
	item, err := scanAnnouncement(r.db.QueryRowContext(ctx, r.placeholder.rebind(`UPDATE announcements
		SET status = ?, publish_at = COALESCE(publish_at, ?), updated_by = ?, updated_at = ?
		WHERE id = ? AND deleted_at = 0
		RETURNING `+announcementColumns()),
		statusPublished,
		publishAt.UTC(),
		actorID,
		time.Now().UTC(),
		targetID,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Announcement{}, ErrAnnouncementNotFound
		}
		return Announcement{}, fmt.Errorf("publish announcement: %w", err)
	}
	return item, nil
}

// Archive marks one non-deleted announcement archived.
func (r *SQLRepository) Archive(ctx context.Context, id uint64, actorID *uint64) (Announcement, error) {
	if err := r.ensureReady(); err != nil {
		return Announcement{}, err
	}
	targetID, err := toDBID(id)
	if err != nil {
		return Announcement{}, err
	}
	item, err := scanAnnouncement(r.db.QueryRowContext(ctx, r.placeholder.rebind(`UPDATE announcements
		SET status = ?, updated_by = ?, updated_at = ?
		WHERE id = ? AND deleted_at = 0
		RETURNING `+announcementColumns()),
		statusArchived,
		actorID,
		time.Now().UTC(),
		targetID,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Announcement{}, ErrAnnouncementNotFound
		}
		return Announcement{}, fmt.Errorf("archive announcement: %w", err)
	}
	return item, nil
}

// Delete soft-deletes one non-deleted announcement.
func (r *SQLRepository) Delete(ctx context.Context, id uint64, actorID uint64, deletedAt time.Time) error {
	if err := r.ensureReady(); err != nil {
		return err
	}
	targetID, err := toDBID(id)
	if err != nil {
		return err
	}
	if deletedAt.IsZero() {
		return ErrInvalidInput
	}
	result, err := r.db.ExecContext(ctx, r.placeholder.rebind(`UPDATE announcements
		SET deleted_by = ?, deleted_at = ?, updated_by = ?, updated_at = ?
		WHERE id = ? AND deleted_at = 0`), nullableUint64(actorID), deletedAt.UTC().Unix(), nullableUint64(actorID), deletedAt.UTC(), targetID)
	if err != nil {
		return fmt.Errorf("delete announcement: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read delete announcement rows affected: %w", err)
	}
	if affected == 0 {
		return ErrAnnouncementNotFound
	}
	return nil
}

func (r *SQLRepository) ensureReady() error {
	if r == nil || r.db == nil {
		return errors.New("announcement repository is unavailable")
	}
	return nil
}

func normalizeCreateInput(input CreateInput) CreateInput {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	input.Level = strings.TrimSpace(input.Level)
	input.Status = strings.TrimSpace(input.Status)
	if input.PublishAt != nil {
		publishAt := input.PublishAt.UTC()
		input.PublishAt = &publishAt
	}
	if input.ExpireAt != nil {
		expireAt := input.ExpireAt.UTC()
		input.ExpireAt = &expireAt
	}
	return input
}

func normalizeUpdateInput(input UpdateInput) UpdateInput {
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	input.Level = strings.TrimSpace(input.Level)
	if input.PublishAt != nil {
		publishAt := input.PublishAt.UTC()
		input.PublishAt = &publishAt
	}
	if input.ExpireAt != nil {
		expireAt := input.ExpireAt.UTC()
		input.ExpireAt = &expireAt
	}
	return input
}

func normalizeListQuery(query ListQuery) ListQuery {
	query.Status = strings.TrimSpace(query.Status)
	query.Level = strings.TrimSpace(query.Level)
	query.Keyword = strings.TrimSpace(query.Keyword)
	query.Sort = strings.TrimSpace(query.Sort)
	if query.Limit <= 0 {
		query.Limit = defaultListLimit
	}
	if query.Limit > maxListLimit {
		query.Limit = maxListLimit
	}
	if query.Offset < 0 {
		query.Offset = 0
	}
	return query
}

func buildAdminWhere(query ListQuery) ([]string, []any, error) {
	where := []string{"deleted_at = 0"}
	args := make([]any, 0, adminFilterCapacity)
	if query.Status != "" {
		args = append(args, query.Status)
		where = append(where, "status = ?")
	}
	if query.Level != "" {
		args = append(args, query.Level)
		where = append(where, "level = ?")
	}
	if query.Pinned != nil {
		args = append(args, *query.Pinned)
		where = append(where, "pinned = ?")
	}
	if query.Keyword != "" {
		keyword := "%" + strings.ToLower(query.Keyword) + "%"
		args = append(args, keyword, keyword)
		where = append(where, "(LOWER(title) LIKE ? OR LOWER(content) LIKE ?)")
	}
	return where, args, nil
}

func adminOrderBy(sort string) string {
	switch sort {
	case sortPinnedPublishDesc:
		return "pinned DESC, publish_at DESC NULLS LAST, id DESC"
	case sortPublishDesc:
		return "publish_at DESC NULLS LAST, id DESC"
	case "", sortUpdatedDesc:
		return "updated_at DESC, id DESC"
	default:
		return "updated_at DESC, id DESC"
	}
}

func announcementColumns() string {
	return `id, title, content, level, status, pinned, publish_at, expire_at,
		created_by, updated_by, deleted_by, created_at, updated_at, deleted_at`
}

func scanAnnouncement(scanner interface{ Scan(dest ...any) error }) (Announcement, error) {
	var item Announcement
	var publishAt sql.NullTime
	var expireAt sql.NullTime
	var createdBy sql.NullInt64
	var updatedBy sql.NullInt64
	var deletedBy sql.NullInt64
	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&item.Content,
		&item.Level,
		&item.Status,
		&item.Pinned,
		&publishAt,
		&expireAt,
		&createdBy,
		&updatedBy,
		&deletedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	); err != nil {
		return Announcement{}, err
	}
	if publishAt.Valid {
		item.PublishAt = &publishAt.Time
	}
	if expireAt.Valid {
		item.ExpireAt = &expireAt.Time
	}
	if createdBy.Valid {
		value, err := uint64FromDBID(createdBy.Int64)
		if err != nil {
			return Announcement{}, err
		}
		item.CreatedBy = &value
	}
	if updatedBy.Valid {
		value, err := uint64FromDBID(updatedBy.Int64)
		if err != nil {
			return Announcement{}, err
		}
		item.UpdatedBy = &value
	}
	if deletedBy.Valid {
		value, err := uint64FromDBID(deletedBy.Int64)
		if err != nil {
			return Announcement{}, err
		}
		item.DeletedBy = &value
	}
	return item, nil
}

func scanAnnouncements(rows *sql.Rows) ([]Announcement, error) {
	items := make([]Announcement, 0)
	for rows.Next() {
		item, err := scanAnnouncement(rows)
		if err != nil {
			return nil, fmt.Errorf("scan announcement row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate announcement rows: %w", err)
	}
	return items, nil
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}

func nullableUint64(value uint64) any {
	if value == 0 {
		return nil
	}
	return value
}

type placeholderStyle int

const (
	placeholderDollar placeholderStyle = iota
	placeholderQuestion
)

func detectPlaceholderStyle(db *sql.DB) placeholderStyle {
	if db == nil || db.Driver() == nil {
		return placeholderDollar
	}
	driverType := strings.ToLower(reflect.TypeOf(db.Driver()).String())
	if strings.Contains(driverType, "sqlite") {
		return placeholderQuestion
	}
	return placeholderDollar
}

func (s placeholderStyle) rebind(query string) string {
	if s == placeholderQuestion {
		return query
	}
	var builder strings.Builder
	builder.Grow(len(query) + placeholderGrowthEstimate)
	index := 1
	for _, current := range query {
		if current != '?' {
			builder.WriteRune(current)
			continue
		}
		builder.WriteByte('$')
		builder.WriteString(strconv.Itoa(index))
		index++
	}
	return builder.String()
}

func toDBID(value uint64) (int64, error) {
	if value == 0 || value > uint64(^uint64(0)>>1) {
		return 0, ErrInvalidInput
	}
	return int64(value), nil
}

func uint64FromDBID(value int64) (uint64, error) {
	if value < 0 {
		return 0, ErrInvalidInput
	}
	return uint64(value), nil
}

var _ Repository = (*SQLRepository)(nil)
