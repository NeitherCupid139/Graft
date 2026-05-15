package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"graft/server/internal/config"
)

const (
	defaultSmokeHealthPath = "/healthz"
	defaultSmokeTimeout    = 10 * time.Second
	defaultSmokeProbeDelay = 200 * time.Millisecond
)

// smokeValidateOptions 封装最小运行时 smoke 验证的显式输入。
type smokeValidateOptions struct {
	migrationDir string
	healthPath   string
	timeout      time.Duration
}

var smokeMigrateRunner = func(cmd *cobra.Command, migrationDir string) error {
	return runMigrateUp(cmd, migrateUpOptions{migrationDir: migrationDir})
}

var smokeServeRunner = runServe
var smokeLoadConfig = config.Load
var smokeHealthChecker = waitForSmokeHealth

// newValidateCommand 创建后端显式验证命令树。
//
// 这里的命令只编排仓库内已经存在的迁移与运行时入口，不负责隐式拉起
// disposable 基础设施，避免把环境准备魔法塞进 core 或 CLI 黑盒里。
func newValidateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "validate",
		Short: "Run explicit backend validation commands",
	}

	command.AddCommand(newValidateSmokeCommand())
	return command
}

// newValidateSmokeCommand 创建最小 smoke 验证子命令。
//
// 该命令显式执行一次迁移、启动运行时、探测 `/healthz`，然后主动停止服务，
// 用来把现有 disposable PostgreSQL/Redis 验证流压缩成一个仓库内入口。
func newValidateSmokeCommand() *cobra.Command {
	var opts smokeValidateOptions

	command := &cobra.Command{
		Use:   "smoke",
		Short: "Run migrations and probe the runtime health endpoint once",
		Long: "graft validate smoke is an explicit backend smoke validation command. " +
			"It runs migrations first, starts the runtime against already-prepared infrastructure, " +
			"waits for the health endpoint to become ready, and then shuts the runtime down.",
		Example:      "  graft validate smoke\n  graft validate smoke --timeout 15s",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidateSmoke(cmd, args, opts)
		},
	}

	command.Flags().StringVar(&opts.migrationDir, "dir", defaultMigrationDir, "migration directory")
	command.Flags().StringVar(&opts.healthPath, "health-path", defaultSmokeHealthPath, "health probe path")
	command.Flags().DurationVar(&opts.timeout, "timeout", defaultSmokeTimeout, "maximum time to wait for the smoke health probe")
	return command
}

// runValidateSmoke 执行最小运行时 smoke 验证闭环。
//
// 顺序语义：
//   - 先执行显式迁移，保持 schema 变更入口仍然可见。
//   - 再启动运行时并轮询健康检查，避免把成功判断退化为“进程未立刻退出”。
//   - 健康检查成功后主动取消运行时上下文，验证服务可以完成一次最小启动与关闭。
func runValidateSmoke(cmd *cobra.Command, args []string, opts smokeValidateOptions) error {
	if err := smokeMigrateRunner(cmd, opts.migrationDir); err != nil {
		return fmt.Errorf("run smoke migrations: %w", err)
	}

	cfg, err := smokeLoadConfig()
	if err != nil {
		return fmt.Errorf("load smoke config: %w", err)
	}

	probeURL, err := buildSmokeProbeURL(cfg.HTTP.Addr, opts.healthPath)
	if err != nil {
		return fmt.Errorf("build smoke probe url: %w", err)
	}

	parentCtx := cmd.Context()
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	runCtx, cancelRun := context.WithCancel(parentCtx)
	defer cancelRun()

	serveCommand := &cobra.Command{}
	serveCommand.SetContext(runCtx)
	serveCommand.SetOut(cmd.OutOrStdout())
	serveCommand.SetErr(cmd.ErrOrStderr())

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- smokeServeRunner(serveCommand, nil)
	}()

	probeCtx, cancelProbe := context.WithTimeout(runCtx, opts.timeout)
	defer cancelProbe()

	healthErrCh := make(chan error, 1)
	go func() {
		healthErrCh <- smokeHealthChecker(probeCtx, probeURL)
	}()

	select {
	case serveErr := <-serveErrCh:
		cancelProbe()
		if serveErr == nil {
			return errors.New("smoke runtime exited before health probe completed")
		}
		return fmt.Errorf("run smoke server: %w", serveErr)
	case err := <-healthErrCh:
		if err != nil {
			cancelRun()
			serveErr := <-serveErrCh
			if serveErr != nil {
				return errors.Join(
					fmt.Errorf("wait for smoke health check: %w", err),
					fmt.Errorf("run smoke server: %w", serveErr),
				)
			}
			return fmt.Errorf("wait for smoke health check: %w", err)
		}
	}

	cancelRun()
	if err := <-serveErrCh; err != nil {
		return fmt.Errorf("shutdown smoke server: %w", err)
	}

	return nil
}

// buildSmokeProbeURL 把监听地址转换为本地健康检查 URL。
//
// 为了兼容 `:8080`、`0.0.0.0:8080` 与 `127.0.0.1:8080` 这类常见监听写法，
// 这里会把 wildcard 地址归一化为 loopback 探测目标。
func buildSmokeProbeURL(listenAddr string, healthPath string) (string, error) {
	normalizedPath := strings.TrimSpace(healthPath)
	if normalizedPath == "" {
		return "", errors.New("health path is required")
	}
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}

	host, port, err := net.SplitHostPort(strings.TrimSpace(listenAddr))
	if err != nil {
		return "", fmt.Errorf("parse listen address %q: %w", listenAddr, err)
	}

	switch strings.TrimSpace(host) {
	case "", "0.0.0.0", "::":
		host = "127.0.0.1"
	}

	return "http://" + net.JoinHostPort(host, port) + normalizedPath, nil
}

// waitForSmokeHealth 轮询健康检查接口直到成功或超时。
//
// 这里把“服务真的完成最小启动”的判定收敛为一次 `200 OK` 响应，避免把
// smoke 验证简化成只看进程是否存活。
func waitForSmokeHealth(ctx context.Context, probeURL string) error {
	client := &http.Client{
		Timeout: time.Second,
	}

	ticker := time.NewTicker(defaultSmokeProbeDelay)
	defer ticker.Stop()

	var lastErr error
	for {
		if err := probeSmokeHealthOnce(ctx, client, probeURL); err == nil {
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			if lastErr != nil {
				return fmt.Errorf("probe %s: %w", probeURL, lastErr)
			}
			return fmt.Errorf("probe %s: %w", probeURL, ctx.Err())
		case <-ticker.C:
		}
	}
}

func probeSmokeHealthOnce(ctx context.Context, client *http.Client, probeURL string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, probeURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(io.LimitReader(response.Body, 256))
		if readErr != nil {
			return fmt.Errorf("health probe returned %s and read body failed: %w", response.Status, readErr)
		}

		trimmedBody := strings.TrimSpace(string(body))
		if trimmedBody == "" {
			return fmt.Errorf("health probe returned %s", response.Status)
		}

		return fmt.Errorf("health probe returned %s: %s", response.Status, trimmedBody)
	}

	return nil
}
