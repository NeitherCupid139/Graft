package scheduler

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"graft/server/internal/cronx"
)

func TestSQLRunRepositoryPersistsRunLifecycle(t *testing.T) {
	db := newSchedulerRepositoryTestDB(t)
	repo, err := NewSQLRunRepository(db)
	if err != nil {
		t.Fatalf("new repository: %v", err)
	}

	startedAt := time.Date(2026, 6, 5, 8, 0, 0, 0, time.UTC)
	run, err := repo.CreateRun(context.Background(), TaskRun{
		TaskKey:     "audit.audit-log-retention-cleanup",
		TaskName:    "audit.audit-log-retention-cleanup",
		Owner:       "audit",
		Module:      "audit",
		TaskType:    cronx.TaskTypeCron,
		TriggerType: TriggerTypeManual,
		Status:      RunStatusRunning,
		StartedAt:   startedAt,
		CreatedAt:   startedAt,
	})
	if err != nil {
		t.Fatalf("create run: %v", err)
	}

	finishedAt := startedAt.Add(1500 * time.Millisecond)
	finished, err := repo.FinishRun(context.Background(), run.ID, RunStatusSuccess, finishedAt, "")
	if err != nil {
		t.Fatalf("finish run: %v", err)
	}
	if finished.Status != RunStatusSuccess || finished.DurationMS == nil || *finished.DurationMS != 1500 {
		t.Fatalf("unexpected finished run: %#v", finished)
	}

	result, err := repo.ListRuns(context.Background(), RunListQuery{TaskKey: "audit.audit-log-retention-cleanup"})
	if err != nil {
		t.Fatalf("list runs: %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 {
		t.Fatalf("expected one run, got %#v", result)
	}

	latest, ok, err := repo.LatestRunByTask(context.Background(), "audit.audit-log-retention-cleanup")
	if err != nil {
		t.Fatalf("latest run: %v", err)
	}
	if !ok || latest.ID != run.ID {
		t.Fatalf("expected latest run %d, got ok=%v run=%#v", run.ID, ok, latest)
	}
}

func newSchedulerRepositoryTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() {
		_ = db.Close()
	})

	_, err = db.Exec(`CREATE TABLE scheduler_task_runs (
		id integer PRIMARY KEY AUTOINCREMENT,
		task_key text NOT NULL,
		task_name text NOT NULL DEFAULT '',
		owner text NOT NULL DEFAULT '',
		module text NOT NULL DEFAULT '',
		task_type text NOT NULL DEFAULT 'cron',
		trigger_type text NOT NULL,
		status text NOT NULL,
		error text NOT NULL DEFAULT '',
		started_at datetime NOT NULL,
		finished_at datetime NULL,
		duration_ms integer NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}

	return db
}
