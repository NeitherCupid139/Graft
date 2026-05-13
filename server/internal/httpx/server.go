// Package httpx 提供 `server` 运行时使用的 HTTP 服务封装。
package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Server 封装运行时使用的 Gin 引擎与 HTTP 服务实例。
//
// Server 负责把 `Run` / `Shutdown` 的生命周期归属集中到一个显式对象中，
// 避免并发启动或停止时出现状态竞争。Server 支持并发调用生命周期方法。
type Server struct {
	engine *gin.Engine
	mu     sync.Mutex
	server *http.Server
}

// NewServer 创建 MVP 运行时使用的最小 Gin 服务外壳。
func NewServer() *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	return &Server{engine: engine}
}

// Engine 返回供 core 和插件注册路由使用的根路由。
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

// Run 启动 HTTP 服务，并把服务生命周期绑定到给定上下文。
func (s *Server) Run(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	if err := s.bindRunningServer(srv); err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err, ok := <-errCh:
		s.clearRunningServer(srv)
		if !ok {
			return nil
		}
		return fmt.Errorf("listen and serve: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr := s.Shutdown(shutdownCtx)
		<-errCh
		return shutdownErr
	}
}

// Shutdown 在服务运行时执行优雅关闭。
func (s *Server) Shutdown(ctx context.Context) error {
	server := s.detachRunningServer()
	if server == nil {
		return nil
	}

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return nil
}

func (s *Server) bindRunningServer(server *http.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 这里显式串行化运行中服务指针的所有权，避免并发 Run / Shutdown
	// 在半完成状态上竞争，导致重复启动或错误清理。
	if s.server != nil {
		return errors.New("http server already running")
	}

	s.server = server
	return nil
}

func (s *Server) detachRunningServer() *http.Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	server := s.server
	s.server = nil
	return server
}

func (s *Server) clearRunningServer(server *http.Server) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.server == server {
		s.server = nil
	}
}
