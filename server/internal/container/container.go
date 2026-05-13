// Package container 提供核心与插件共用的显式单例注册能力。
package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Provider 定义单例服务的构造函数。
type Provider func(resolver Resolver) (any, error)

// Registry 定义显式单例注册能力。
type Registry interface {
	// RegisterSingleton 使用给定 key 注册单例构造函数。
	RegisterSingleton(key any, provider Provider) error
}

// Resolver 定义显式单例解析能力。
type Resolver interface {
	// Resolve 根据给定 key 返回单例实例。
	Resolve(key any) (any, error)
}

// Container 是运行时使用的显式单例容器。
//
// Container 负责保存单例构造函数、缓存已构建实例，并在并发解析时
// 复用同一次构建过程，避免重复创建共享资源。Container 支持并发访问。
type Container struct {
	// mu 串行化 provider、实例缓存和构建中状态的读写。
	mu sync.RWMutex
	// providers 保存服务 key 到单例构造函数的映射。
	providers map[string]Provider
	// instances 缓存已经成功构建的共享单例。
	instances map[string]any
	// inflight 记录正在构建的单例，确保并发调用共享同一次构建结果。
	inflight map[string]*inflightCall
}

type inflightCall struct {
	done chan struct{}
	val  any
	err  error
}

// New 创建一个空的单例容器。
func New() *Container {
	return &Container{
		providers: make(map[string]Provider),
		instances: make(map[string]any),
		inflight:  make(map[string]*inflightCall),
	}
}

// RegisterSingleton 为一个服务 key 注册唯一的单例构造函数。
func (c *Container) RegisterSingleton(key any, provider Provider) error {
	if provider == nil {
		return errors.New("provider is required")
	}

	name := keyName(key)

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.providers[name]; exists {
		return fmt.Errorf("service already registered: %s", name)
	}

	c.providers[name] = provider
	return nil
}

// Resolve 解析指定 key 对应的单例实例。
//
// 如果实例尚未创建，Resolve 会触发一次构建并缓存结果；如果同一实例
// 正在被其它 goroutine 构建，当前调用会等待并复用那次构建结果。
func (c *Container) Resolve(key any) (any, error) {
	name := keyName(key)

	c.mu.Lock()
	if instance, ok := c.instances[name]; ok {
		c.mu.Unlock()
		return instance, nil
	}

	// 第一个进入的调用负责实际构建，后续并发调用等待 inflight 结果，
	// 这样可以避免高成本 provider 被重复执行。
	if call, ok := c.inflight[name]; ok {
		c.mu.Unlock()
		<-call.done
		if call.err != nil {
			return nil, fmt.Errorf("build service %s: %w", name, call.err)
		}
		return call.val, nil
	}

	provider, ok := c.providers[name]
	if !ok {
		c.mu.Unlock()
		return nil, fmt.Errorf("service not registered: %s", name)
	}

	call := &inflightCall{done: make(chan struct{})}
	c.inflight[name] = call
	c.mu.Unlock()

	instance, err := provider(c)

	c.mu.Lock()
	if err == nil {
		c.instances[name] = instance
	}
	call.val = instance
	call.err = err
	close(call.done)
	delete(c.inflight, name)
	c.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("build service %s: %w", name, err)
	}
	return instance, nil
}

func keyName(key any) string {
	if key == nil {
		return "<nil>"
	}

	return reflect.TypeOf(key).String()
}
