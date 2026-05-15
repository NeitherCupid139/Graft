package cli

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"graft/server/internal/config"
)

// TestRunValidateSmokeRunsMigrateBeforeServe 验证 smoke 验证会先执行迁移，
// 再等待健康检查成功，最后主动停止运行时。
func TestRunValidateSmokeRunsMigrateBeforeServe(t *testing.T) {
	originalMigrateRunner := smokeMigrateRunner
	originalServeRunner := smokeServeRunner
	originalLoadConfig := smokeLoadConfig
	originalHealthChecker := smokeHealthChecker
	defer func() {
		smokeMigrateRunner = originalMigrateRunner
		smokeServeRunner = originalServeRunner
		smokeLoadConfig = originalLoadConfig
		smokeHealthChecker = originalHealthChecker
	}()

	var steps []string
	serveStarted := make(chan struct{})

	smokeMigrateRunner = func(cmd *cobra.Command, migrationDir string) error {
		steps = append(steps, "migrate:"+migrationDir)
		return nil
	}
	smokeLoadConfig = func() (*config.Config, error) {
		return &config.Config{
			HTTP: config.HTTPConfig{Addr: ":18080"},
		}, nil
	}
	smokeServeRunner = func(cmd *cobra.Command, args []string) error {
		steps = append(steps, "serve-start")
		close(serveStarted)
		<-cmd.Context().Done()
		steps = append(steps, "serve-stop")
		return nil
	}
	smokeHealthChecker = func(ctx context.Context, probeURL string) error {
		<-serveStarted
		steps = append(steps, "health:"+probeURL)
		return nil
	}

	err := runValidateSmoke(&cobra.Command{}, nil, smokeValidateOptions{
		migrationDir: defaultMigrationDir,
		healthPath:   defaultSmokeHealthPath,
		timeout:      time.Second,
	})
	if err != nil {
		t.Fatalf("run validate smoke: %v", err)
	}

	expected := []string{
		"migrate:" + defaultMigrationDir,
		"serve-start",
		"health:http://127.0.0.1:18080/healthz",
		"serve-stop",
	}
	if !reflect.DeepEqual(steps, expected) {
		t.Fatalf("expected %v, got %v", expected, steps)
	}
}

// TestRunValidateSmokeStopsAfterMigrationFailure 验证迁移失败时不会继续启动运行时。
func TestRunValidateSmokeStopsAfterMigrationFailure(t *testing.T) {
	originalMigrateRunner := smokeMigrateRunner
	originalServeRunner := smokeServeRunner
	defer func() {
		smokeMigrateRunner = originalMigrateRunner
		smokeServeRunner = originalServeRunner
	}()

	smokeMigrateRunner = func(cmd *cobra.Command, migrationDir string) error {
		return errors.New("migrate failed")
	}
	smokeServeRunner = func(cmd *cobra.Command, args []string) error {
		t.Fatal("serve runner should not be called")
		return nil
	}

	err := runValidateSmoke(&cobra.Command{}, nil, smokeValidateOptions{
		migrationDir: defaultMigrationDir,
		healthPath:   defaultSmokeHealthPath,
		timeout:      time.Second,
	})
	if err == nil {
		t.Fatal("expected smoke validation error")
	}
	if !strings.Contains(err.Error(), "run smoke migrations") {
		t.Fatalf("expected migration context, got %v", err)
	}
}

// TestRunValidateSmokeReturnsServeFailure 验证运行时在健康检查前退出时会立刻返回服务错误。
func TestRunValidateSmokeReturnsServeFailure(t *testing.T) {
	originalMigrateRunner := smokeMigrateRunner
	originalServeRunner := smokeServeRunner
	originalLoadConfig := smokeLoadConfig
	originalHealthChecker := smokeHealthChecker
	defer func() {
		smokeMigrateRunner = originalMigrateRunner
		smokeServeRunner = originalServeRunner
		smokeLoadConfig = originalLoadConfig
		smokeHealthChecker = originalHealthChecker
	}()

	smokeMigrateRunner = func(cmd *cobra.Command, migrationDir string) error {
		return nil
	}
	smokeLoadConfig = func() (*config.Config, error) {
		return &config.Config{
			HTTP: config.HTTPConfig{Addr: ":18080"},
		}, nil
	}
	smokeServeRunner = func(cmd *cobra.Command, args []string) error {
		return errors.New("listen failed")
	}
	smokeHealthChecker = func(ctx context.Context, probeURL string) error {
		<-ctx.Done()
		return ctx.Err()
	}

	err := runValidateSmoke(&cobra.Command{}, nil, smokeValidateOptions{
		migrationDir: defaultMigrationDir,
		healthPath:   defaultSmokeHealthPath,
		timeout:      time.Second,
	})
	if err == nil {
		t.Fatal("expected smoke validation error")
	}
	if !strings.Contains(err.Error(), "run smoke server") {
		t.Fatalf("expected serve context, got %v", err)
	}
}

// TestRunValidateSmokeReturnsHealthFailure 验证健康检查失败时会停止运行时并返回探测错误。
func TestRunValidateSmokeReturnsHealthFailure(t *testing.T) {
	originalMigrateRunner := smokeMigrateRunner
	originalServeRunner := smokeServeRunner
	originalLoadConfig := smokeLoadConfig
	originalHealthChecker := smokeHealthChecker
	defer func() {
		smokeMigrateRunner = originalMigrateRunner
		smokeServeRunner = originalServeRunner
		smokeLoadConfig = originalLoadConfig
		smokeHealthChecker = originalHealthChecker
	}()

	smokeMigrateRunner = func(cmd *cobra.Command, migrationDir string) error {
		return nil
	}
	smokeLoadConfig = func() (*config.Config, error) {
		return &config.Config{
			HTTP: config.HTTPConfig{Addr: ":18080"},
		}, nil
	}
	smokeServeRunner = func(cmd *cobra.Command, args []string) error {
		<-cmd.Context().Done()
		return nil
	}
	smokeHealthChecker = func(ctx context.Context, probeURL string) error {
		return errors.New("health failed")
	}

	err := runValidateSmoke(&cobra.Command{}, nil, smokeValidateOptions{
		migrationDir: defaultMigrationDir,
		healthPath:   defaultSmokeHealthPath,
		timeout:      time.Second,
	})
	if err == nil {
		t.Fatal("expected smoke validation error")
	}
	if !strings.Contains(err.Error(), "wait for smoke health check") {
		t.Fatalf("expected health-check context, got %v", err)
	}
}

// TestBuildSmokeProbeURLUsesLoopbackForWildcard 验证 wildcard 监听地址会转换为本地可探测的 loopback URL。
func TestBuildSmokeProbeURLUsesLoopbackForWildcard(t *testing.T) {
	testCases := []struct {
		name     string
		addr     string
		path     string
		expected string
	}{
		{
			name:     "empty host",
			addr:     ":8080",
			path:     "/healthz",
			expected: "http://127.0.0.1:8080/healthz",
		},
		{
			name:     "ipv4 wildcard",
			addr:     "0.0.0.0:8080",
			path:     "healthz",
			expected: "http://127.0.0.1:8080/healthz",
		},
		{
			name:     "localhost",
			addr:     "127.0.0.1:8080",
			path:     "/healthz",
			expected: "http://127.0.0.1:8080/healthz",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := buildSmokeProbeURL(testCase.addr, testCase.path)
			if err != nil {
				t.Fatalf("build smoke probe url: %v", err)
			}
			if actual != testCase.expected {
				t.Fatalf("expected %s, got %s", testCase.expected, actual)
			}
		})
	}
}

// TestNewRootCommandRegistersValidateSmoke 验证根命令始终注册 `validate smoke` 子命令。
func TestNewRootCommandRegistersValidateSmoke(t *testing.T) {
	command := NewRootCommand()

	found, _, err := command.Find([]string{"validate", "smoke"})
	if err != nil {
		t.Fatalf("find validate smoke command: %v", err)
	}
	if found == nil || found.Name() != "smoke" {
		t.Fatalf("expected smoke command, got %#v", found)
	}
}
