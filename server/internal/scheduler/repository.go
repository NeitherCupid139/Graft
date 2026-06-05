package scheduler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"graft/server/internal/cronx"
)

const defaultRunListLimit = 20

// SQLRunRepository persists scheduler runtime history in scheduler_task_runs.
type SQLRunRepository struct {
	db *sql.DB
}

// NewSQLRunRepository builds a scheduler run-history repository from the shared SQL pool.
func NewSQLRunRepository(db *sql.DB) (*SQLRunRepository, error) {
	if db == nil {
		return nil, errors.New("scheduler run repository requires a non-nil sql db")
	}

	return &SQLRunRepository{db: db}, nil
}

// CreateRun inserts a running scheduler_task_runs row.
func (r *SQLRunRepository) CreateRun(ctx context.Context, run TaskRun) (TaskRun, error) {
	if r == nil || r.db == nil {
		return TaskRun{}, errors.New("scheduler run repository is unavailable")
	}
	if run.TaskKey == "" {
		return TaskRun{}, errors.New("scheduler run task key is required")
	}
	if run.StartedAt.IsZero() {
		return TaskRun{}, errors.New("scheduler run started_at is required")
	}
	if run.CreatedAt.IsZero() {
		run.CreatedAt = run.StartedAt
	}
	if run.Status == "" {
		run.Status = RunStatusRunning
	}

	row := r.db.QueryRowContext(ctx, `INSERT INTO scheduler_task_runs (
		task_key,
		task_name,
		owner,
		module,
		task_type,
		trigger_type,
		status,
		error,
		started_at,
		finished_at,
		duration_ms,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NULL, NULL, $10)
	RETURNING id`,
		run.TaskKey,
		run.TaskName,
		run.Owner,
		run.Module,
		string(run.TaskType),
		string(run.TriggerType),
		string(run.Status),
		run.Error,
		run.StartedAt.UTC(),
		run.CreatedAt.UTC(),
	)

	var id int64
	if err := row.Scan(&id); err != nil {
		return TaskRun{}, fmt.Errorf("create scheduler task run: %w", err)
	}
	run.ID = uint64(id)
	return run, nil
}

// FinishRun closes a running scheduler_task_runs row.
func (r *SQLRunRepository) FinishRun(
	ctx context.Context,
	id uint64,
	status RunStatus,
	finishedAt time.Time,
	errorMessage string,
) (TaskRun, error) {
	if r == nil || r.db == nil {
		return TaskRun{}, errors.New("scheduler run repository is unavailable")
	}
	if id == 0 {
		return TaskRun{}, errors.New("scheduler run id is required")
	}
	if id > uint64(math.MaxInt64) {
		return TaskRun{}, errors.New("scheduler run id is too large")
	}
	if status == "" {
		return TaskRun{}, errors.New("scheduler run status is required")
	}
	if finishedAt.IsZero() {
		return TaskRun{}, errors.New("scheduler run finished_at is required")
	}

	var startedAt time.Time
	sqlID := int64(id)
	if err := r.db.QueryRowContext(ctx, `SELECT started_at FROM scheduler_task_runs WHERE id = $1`, sqlID).Scan(&startedAt); err != nil {
		return TaskRun{}, fmt.Errorf("read scheduler task run start: %w", err)
	}
	durationMS := finishedAt.UTC().Sub(startedAt.UTC()).Milliseconds()
	if durationMS < 0 {
		durationMS = 0
	}

	result, err := r.db.ExecContext(ctx, `UPDATE scheduler_task_runs
	SET status = $1,
		error = $2,
		finished_at = $3,
		duration_ms = $4
	WHERE id = $5`,
		string(status),
		errorMessage,
		finishedAt.UTC(),
		durationMS,
		sqlID,
	)
	if err != nil {
		return TaskRun{}, fmt.Errorf("update scheduler task run: %w", err)
	}
	_ = result

	row := r.db.QueryRowContext(ctx, `SELECT id, task_key, task_name, owner, module, task_type, trigger_type, status, error, started_at, finished_at, duration_ms, created_at
	FROM scheduler_task_runs
	WHERE id = $1`, sqlID)
	run, err := scanTaskRun(row)
	if err != nil {
		return TaskRun{}, fmt.Errorf("finish scheduler task run: %w", err)
	}

	return run, nil
}

// ListRuns returns one stable page of task run history.
func (r *SQLRunRepository) ListRuns(ctx context.Context, query RunListQuery) (RunListResult, error) {
	if r == nil || r.db == nil {
		return RunListResult{}, errors.New("scheduler run repository is unavailable")
	}
	if query.TaskKey == "" {
		return RunListResult{}, errors.New("scheduler run task key is required")
	}
	if query.Limit <= 0 {
		query.Limit = defaultRunListLimit
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM scheduler_task_runs WHERE task_key = $1`, query.TaskKey).Scan(&total); err != nil {
		return RunListResult{}, fmt.Errorf("count scheduler task runs: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, task_key, task_name, owner, module, task_type, trigger_type, status, error, started_at, finished_at, duration_ms, created_at
	FROM scheduler_task_runs
	WHERE task_key = $1
	ORDER BY started_at DESC, id DESC
	LIMIT $2 OFFSET $3`, query.TaskKey, query.Limit, query.Offset)
	if err != nil {
		return RunListResult{}, fmt.Errorf("list scheduler task runs: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	items := make([]TaskRun, 0, query.Limit)
	for rows.Next() {
		run, scanErr := scanTaskRun(rows)
		if scanErr != nil {
			return RunListResult{}, scanErr
		}
		items = append(items, run)
	}
	if err := rows.Err(); err != nil {
		return RunListResult{}, fmt.Errorf("iterate scheduler task runs: %w", err)
	}

	return RunListResult{Items: items, Total: total}, nil
}

// LatestRunByTask returns the latest persisted run for one task key.
func (r *SQLRunRepository) LatestRunByTask(ctx context.Context, taskKey string) (TaskRun, bool, error) {
	if r == nil || r.db == nil {
		return TaskRun{}, false, errors.New("scheduler run repository is unavailable")
	}
	if taskKey == "" {
		return TaskRun{}, false, errors.New("scheduler run task key is required")
	}

	row := r.db.QueryRowContext(ctx, `SELECT id, task_key, task_name, owner, module, task_type, trigger_type, status, error, started_at, finished_at, duration_ms, created_at
	FROM scheduler_task_runs
	WHERE task_key = $1
	ORDER BY started_at DESC, id DESC
	LIMIT 1`, taskKey)

	run, err := scanTaskRun(row)
	if errors.Is(err, sql.ErrNoRows) {
		return TaskRun{}, false, nil
	}
	if err != nil {
		return TaskRun{}, false, err
	}

	return run, true, nil
}

type taskRunScanner interface {
	Scan(dest ...any) error
}

func scanTaskRun(scanner taskRunScanner) (TaskRun, error) {
	var run TaskRun
	var id int64
	var taskType string
	var triggerType string
	var status string
	var finishedAt sql.NullTime
	var durationMS sql.NullInt64

	if err := scanner.Scan(
		&id,
		&run.TaskKey,
		&run.TaskName,
		&run.Owner,
		&run.Module,
		&taskType,
		&triggerType,
		&status,
		&run.Error,
		&run.StartedAt,
		&finishedAt,
		&durationMS,
		&run.CreatedAt,
	); err != nil {
		return TaskRun{}, err
	}

	run.ID = uint64(id)
	run.TaskType = cronTaskType(taskType)
	run.TriggerType = TriggerType(triggerType)
	run.Status = RunStatus(status)
	if finishedAt.Valid {
		finished := finishedAt.Time
		run.FinishedAt = &finished
	}
	if durationMS.Valid {
		duration := durationMS.Int64
		run.DurationMS = &duration
	}

	return run, nil
}

func cronTaskType(value string) cronx.TaskType {
	if value == "" {
		return cronx.TaskTypeCron
	}
	return cronx.TaskType(value)
}
