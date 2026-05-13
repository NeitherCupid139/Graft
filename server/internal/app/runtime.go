// Package app 组装 Graft 的显式运行时外壳。
package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"graft/server/internal/config"
	"graft/server/internal/container"
	"graft/server/internal/cronx"
	"graft/server/internal/database"
	"graft/server/internal/httpx"
	"graft/server/internal/menu"
	"graft/server/internal/permission"
	"graft/server/internal/plugin"
	"graft/server/internal/redisx"
	"graft/server/internal/store"
	"graft/server/internal/store/entstore"
)

// Runtime 持有 MVP 运行时的核心资源与插件生命周期执行入口。
//
// Runtime 把配置、数据库、Redis、HTTP 服务、注册中心和插件管理器集中
// 到一个显式对象中，方便在失败路径和正常关闭路径统一回收资源。
type Runtime struct {
	config             *config.Config
	database           *database.Resources
	redis              *redis.Client
	server             *httpx.Server
	services           *container.Container
	stores             store.Factory
	menuRegistry       *menu.Registry
	permissionRegistry *permission.Registry
	cronRegistry       *cronx.Registry
	pluginManager      *plugin.Manager
}

// NewRuntime 使用给定插件构造显式的 MVP 运行时外壳。
func NewRuntime(plugins ...plugin.Plugin) (*Runtime, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	databaseResources, err := database.Open(cfg.Database)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisx.Open(context.Background(), cfg.Redis)
	if err != nil {
		_ = database.Close(databaseResources)
		return nil, err
	}

	server := httpx.NewServer()
	services := container.New()
	stores := entstore.NewFactory(databaseResources.Client)
	menuRegistry := menu.NewRegistry()
	permissionRegistry := permission.NewRegistry()
	cronRegistry := cronx.NewRegistry()
	pluginManager := plugin.NewManager()

	runtime := &Runtime{
		config:             cfg,
		database:           databaseResources,
		redis:              redisClient,
		server:             server,
		services:           services,
		stores:             stores,
		menuRegistry:       menuRegistry,
		permissionRegistry: permissionRegistry,
		cronRegistry:       cronRegistry,
		pluginManager:      pluginManager,
	}

	if err := runtime.registerCoreServices(); err != nil {
		_ = runtime.closeCoreResources()
		return nil, err
	}

	runtime.registerCoreRoutes(server.Engine())

	for _, current := range plugins {
		if err := runtime.pluginManager.RegisterPlugin(current); err != nil {
			_ = runtime.closeCoreResources()
			return nil, err
		}
	}

	return runtime, nil
}

// Run 先执行插件注册与启动，再启动 HTTP 服务。
//
// 如果任一阶段失败，Run 会按已启动的实际范围反向释放插件与核心资源，
// 避免把半初始化状态泄漏到调用方。
func (r *Runtime) Run(runCtx context.Context) error {
	pluginCtx := &plugin.Context{
		Config:             r.config,
		Redis:              r.redis,
		Router:             r.server.Engine().Group("/api"),
		Services:           r.services,
		Stores:             r.stores,
		MenuRegistry:       r.menuRegistry,
		PermissionRegistry: r.permissionRegistry,
		CronRegistry:       r.cronRegistry,
	}

	ordered, err := r.pluginManager.Ordered()
	if err != nil {
		return err
	}

	booted := make([]plugin.Plugin, 0, len(ordered))
	for _, p := range ordered {
		if err := p.Register(pluginCtx); err != nil {
			return r.cleanupAfterFailure(pluginCtx, booted, fmt.Errorf("register plugin %s: %w", p.Name(), err))
		}
	}

	for _, p := range ordered {
		if err := p.Boot(pluginCtx); err != nil {
			return r.cleanupAfterFailure(pluginCtx, booted, fmt.Errorf("boot plugin %s: %w", p.Name(), err))
		}
		booted = append(booted, p)
	}

	if err := r.server.Run(runCtx, r.config.HTTP.Addr); err != nil {
		return r.cleanupAfterFailure(pluginCtx, booted, err)
	}

	if err := shutdownPlugins(pluginCtx, booted); err != nil {
		return r.cleanupAfterFailure(pluginCtx, nil, err)
	}

	if err := r.closeCoreResources(); err != nil {
		return err
	}

	return nil
}

func (r *Runtime) registerCoreRoutes(engine *gin.Engine) {
	engine.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"menus":       len(r.menuRegistry.Items()),
			"permissions": len(r.permissionRegistry.Items()),
			"jobs":        len(r.cronRegistry.Items()),
		})
	})
}

func (r *Runtime) registerCoreServices() error {
	if err := r.services.RegisterSingleton((*config.Config)(nil), func(resolver container.Resolver) (any, error) {
		return r.config, nil
	}); err != nil {
		return err
	}

	if err := r.services.RegisterSingleton((*store.Factory)(nil), func(resolver container.Resolver) (any, error) {
		return r.stores, nil
	}); err != nil {
		return err
	}

	return r.services.RegisterSingleton((*redis.Client)(nil), func(resolver container.Resolver) (any, error) {
		return r.redis, nil
	})
}

func shutdownPlugins(ctx *plugin.Context, ordered []plugin.Plugin) error {
	var shutdownErr error
	for i := len(ordered) - 1; i >= 0; i-- {
		// 关闭顺序必须与启动顺序相反，避免后启动的依赖还未释放时，上游
		// 插件先被销毁，导致清理逻辑访问失效资源。
		if err := ordered[i].Shutdown(ctx); err != nil {
			shutdownErr = errors.Join(shutdownErr, fmt.Errorf("shutdown plugin %s: %w", ordered[i].Name(), err))
		}
	}

	return shutdownErr
}

func (r *Runtime) closeCoreResources() error {
	var closeErr error
	if r.redis != nil {
		if err := r.redis.Close(); err != nil {
			closeErr = errors.Join(closeErr, fmt.Errorf("close redis: %w", err))
		}
		r.redis = nil
	}

	if r.database != nil {
		if err := database.Close(r.database); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
		r.database = nil
	}

	return closeErr
}

func (r *Runtime) cleanupAfterFailure(ctx *plugin.Context, booted []plugin.Plugin, cause error) error {
	var err error = cause
	if shutdownErr := shutdownPlugins(ctx, booted); shutdownErr != nil {
		err = errors.Join(err, shutdownErr)
	}
	if closeErr := r.closeCoreResources(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}
	return err
}
