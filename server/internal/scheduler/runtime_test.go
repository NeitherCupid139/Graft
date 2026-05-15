package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"

	"graft/server/internal/cronx"
)

// TestRegisterJobRejectsInvalidDeclarations 验证调度器会拒绝缺失执行入口或非法表达式的任务声明。
func TestRegisterJobRejectsInvalidDeclarations(t *testing.T) {
	runtime := New(zap.NewNop())

	if err := runtime.RegisterJob(cronx.Job{Name: "", Schedule: "* * * * * *", Run: func(context.Context) error { return nil }}); err == nil {
		t.Fatal("expected empty job name to fail")
	}
	if err := runtime.RegisterJob(cronx.Job{Name: "cleanup", Schedule: "", Run: func(context.Context) error { return nil }}); err == nil {
		t.Fatal("expected empty schedule to fail")
	}
	if err := runtime.RegisterJob(cronx.Job{Name: "cleanup", Schedule: "* * * * * *"}); err == nil {
		t.Fatal("expected missing run function to fail")
	}
}

// TestRegisterJobRejectsDuplicateName 验证重复任务名会在注册阶段显式失败。
func TestRegisterJobRejectsDuplicateName(t *testing.T) {
	runtime := New(zap.NewNop())
	job := cronx.Job{
		Name:     "cleanup",
		Schedule: "*/1 * * * * *",
		Run:      func(context.Context) error { return nil },
	}

	if err := runtime.RegisterJob(job); err != nil {
		t.Fatalf("register first job: %v", err)
	}
	if err := runtime.RegisterJob(job); err == nil {
		t.Fatal("expected duplicate registration to fail")
	}
}

// TestStartAndStopRunsRegisteredJob 验证最小调度器可以启动、执行一次任务并正常停止。
func TestStartAndStopRunsRegisteredJob(t *testing.T) {
	runtime := New(zap.NewNop())
	triggered := make(chan struct{}, 1)

	if err := runtime.RegisterJob(cronx.Job{
		Name:     "heartbeat",
		Schedule: "*/1 * * * * *",
		Run: func(ctx context.Context) error {
			select {
			case triggered <- struct{}{}:
			default:
			}
			return nil
		},
	}); err != nil {
		t.Fatalf("register job: %v", err)
	}

	if err := runtime.Start(); err != nil {
		t.Fatalf("start runtime: %v", err)
	}
	defer func() {
		_ = runtime.Stop(context.Background())
	}()

	select {
	case <-triggered:
	case <-time.After(1500 * time.Millisecond):
		t.Fatal("expected scheduled job to run")
	}

	if err := runtime.Stop(context.Background()); err != nil {
		t.Fatalf("stop runtime: %v", err)
	}
}

// TestRemoveJobPreventsFutureExecution 验证移除任务后后续调度不会再次触发该任务。
func TestRemoveJobPreventsFutureExecution(t *testing.T) {
	runtime := New(zap.NewNop())
	triggered := make(chan struct{}, 2)

	if err := runtime.RegisterJob(cronx.Job{
		Name:     "cleanup",
		Schedule: "*/1 * * * * *",
		Run: func(ctx context.Context) error {
			triggered <- struct{}{}
			return nil
		},
	}); err != nil {
		t.Fatalf("register job: %v", err)
	}
	if err := runtime.Start(); err != nil {
		t.Fatalf("start runtime: %v", err)
	}

	select {
	case <-triggered:
	case <-time.After(1500 * time.Millisecond):
		t.Fatal("expected first scheduled execution")
	}

	if err := runtime.RemoveJob("cleanup"); err != nil {
		t.Fatalf("remove job: %v", err)
	}

	select {
	case <-triggered:
		t.Fatal("expected removed job not to run again")
	case <-time.After(1200 * time.Millisecond):
	}

	if err := runtime.Stop(context.Background()); err != nil {
		t.Fatalf("stop runtime: %v", err)
	}
}

// TestStopHonorsContextCancellation 验证 Stop 会把外部取消信号作为稳定错误返回。
func TestStopHonorsContextCancellation(t *testing.T) {
	runtime := New(zap.NewNop())
	if err := runtime.Start(); err != nil {
		t.Fatalf("start runtime: %v", err)
	}

	stopCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := runtime.Stop(stopCtx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}
