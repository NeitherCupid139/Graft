// Package plugin 定义运行时插件契约与生命周期管理能力。
package plugin

import (
	"errors"
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"graft/server/internal/config"
	"graft/server/internal/container"
	"graft/server/internal/cronx"
	"graft/server/internal/menu"
	"graft/server/internal/permission"
	"graft/server/internal/store"
)

// Plugin 定义所有后端插件都必须实现的稳定生命周期契约。
type Plugin interface {
	// Name 返回插件的稳定标识，用于依赖声明和运行时元数据。
	Name() string
	// Version 返回当前插件版本。
	Version() string
	// DependsOn 返回当前插件依赖的插件名称列表。
	DependsOn() []string
	// Register 负责声明路由、权限、菜单、任务和公开服务。
	Register(ctx *Context) error
	// Boot 在所有插件完成注册后启动运行时行为。
	Boot(ctx *Context) error
	// Shutdown 在停止阶段释放插件资源，调用顺序与启动顺序相反。
	Shutdown(ctx *Context) error
}

// Context 向插件暴露允许使用的显式运行时句柄。
//
// 这里聚合的是插件生命周期真正需要的核心能力，目的是让插件通过稳定
// 边界接入平台，而不是直接触碰 core 内部实现细节。
type Context struct {
	Config             *config.Config
	Redis              *redis.Client
	Router             gin.IRouter
	Services           *container.Container
	Stores             store.Factory
	MenuRegistry       *menu.Registry
	PermissionRegistry *permission.Registry
	CronRegistry       *cronx.Registry
}

// Manager 负责维护插件集合并按依赖关系排序。
type Manager struct {
	plugins []Plugin
}

// NewManager 创建一个空的插件管理器。
func NewManager() *Manager {
	return &Manager{plugins: make([]Plugin, 0)}
}

// RegisterPlugin 在运行时启动前向管理器注册一个插件。
func (m *Manager) RegisterPlugin(p Plugin) error {
	if p == nil {
		return errors.New("plugin is required")
	}

	for _, existing := range m.plugins {
		if existing.Name() == p.Name() {
			return fmt.Errorf("plugin already registered: %s", p.Name())
		}
	}

	m.plugins = append(m.plugins, p)
	return nil
}

// Ordered 按声明的依赖关系返回插件启动顺序。
//
// 这里使用显式拓扑排序而不是隐式注册顺序，避免插件接入规模增加后因为
// 注册位置变化而打破稳定的启动语义。
func (m *Manager) Ordered() ([]Plugin, error) {
	total := len(m.plugins)
	if total == 0 {
		return nil, nil
	}

	index := make(map[string]Plugin, total)
	inDegree := make(map[string]int, total)
	edges := make(map[string][]string, total)

	for _, p := range m.plugins {
		name := p.Name()
		index[name] = p
		inDegree[name] = 0
	}

	for _, p := range m.plugins {
		for _, dependency := range p.DependsOn() {
			if _, ok := index[dependency]; !ok {
				return nil, fmt.Errorf("plugin %s depends on missing plugin %s", p.Name(), dependency)
			}

			edges[dependency] = append(edges[dependency], p.Name())
			inDegree[p.Name()]++
		}
	}

	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue)

	ordered := make([]Plugin, 0, total)
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		ordered = append(ordered, index[name])

		for _, next := range edges[name] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
				sort.Strings(queue)
			}
		}
	}

	if len(ordered) != total {
		return nil, errors.New("plugin dependency cycle detected")
	}

	return ordered, nil
}
