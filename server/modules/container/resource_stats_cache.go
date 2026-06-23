package container

import (
	"context"
	"strings"
	"sync"
	"time"
)

const (
	containerResourceStatsCacheTTL         = 2 * time.Second
	containerResourceStatsCacheStaleWindow = 8 * time.Second
)

type resourceStatsLoader func(context.Context) ResourceSummary

type resourceStatsCache struct {
	mu          sync.Mutex
	ttl         time.Duration
	staleWindow time.Duration
	now         func() time.Time
	items       map[string]resourceStatsCacheEntry
	inflight    map[string]*resourceStatsLoad
}

type resourceStatsCacheEntry struct {
	summary    ResourceSummary
	freshUntil time.Time
	staleUntil time.Time
}

type resourceStatsLoad struct {
	done    chan struct{}
	summary ResourceSummary
}

func newResourceStatsCache(ttl time.Duration, staleWindow time.Duration) *resourceStatsCache {
	if ttl <= 0 {
		ttl = containerResourceStatsCacheTTL
	}
	if staleWindow <= 0 {
		staleWindow = containerResourceStatsCacheStaleWindow
	}
	return &resourceStatsCache{
		ttl:         ttl,
		staleWindow: staleWindow,
		now:         time.Now,
		items:       make(map[string]resourceStatsCacheEntry),
		inflight:    make(map[string]*resourceStatsLoad),
	}
}

func (c *resourceStatsCache) get(ctx context.Context, key string, loader resourceStatsLoader) ResourceSummary {
	key = strings.TrimSpace(key)
	if c == nil || loader == nil || key == "" {
		return unavailableResourceSummary(containerStatsIncompleteReason)
	}

	now := c.now()

	c.mu.Lock()
	entry, hasEntry := c.items[key]
	if summary, ok := resourceStatsFreshHit(entry, hasEntry, now); ok {
		c.mu.Unlock()
		return summary
	}
	if load, ok := c.inflight[key]; ok {
		return c.handleInflight(ctx, entry, hasEntry, now, load)
	}
	if summary, ok := c.serveStaleWhileRefresh(ctx, key, entry, hasEntry, now, loader); ok {
		return summary
	}

	load := &resourceStatsLoad{done: make(chan struct{})}
	c.inflight[key] = load
	c.mu.Unlock()

	summary := c.completeLoad(ctx, key, load, loader)
	return summary
}

func (c *resourceStatsCache) invalidate(keys ...string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		delete(c.items, key)
	}
}

func resourceStatsFreshHit(entry resourceStatsCacheEntry, ok bool, now time.Time) (ResourceSummary, bool) {
	if ok && now.Before(entry.freshUntil) {
		return entry.summary, true
	}
	return ResourceSummary{}, false
}

func (c *resourceStatsCache) handleInflight(
	ctx context.Context,
	entry resourceStatsCacheEntry,
	hasEntry bool,
	now time.Time,
	load *resourceStatsLoad,
) ResourceSummary {
	if hasEntry && now.Before(entry.staleUntil) {
		summary := entry.summary
		c.mu.Unlock()
		return summary
	}
	c.mu.Unlock()
	select {
	case <-load.done:
		return load.summary
	case <-ctx.Done():
		return unavailableResourceSummary(resourceStatsErrorReason(ctx.Err()))
	}
}

func (c *resourceStatsCache) serveStaleWhileRefresh(
	ctx context.Context,
	key string,
	entry resourceStatsCacheEntry,
	hasEntry bool,
	now time.Time,
	loader resourceStatsLoader,
) (ResourceSummary, bool) {
	if !hasEntry || !now.Before(entry.staleUntil) {
		return ResourceSummary{}, false
	}
	load := &resourceStatsLoad{done: make(chan struct{})}
	c.inflight[key] = load
	stale := entry.summary
	refreshCtx := context.WithoutCancel(ctx)
	c.mu.Unlock()
	go c.completeLoad(refreshCtx, key, load, loader)
	return stale, true
}

func (c *resourceStatsCache) completeLoad(ctx context.Context, key string, load *resourceStatsLoad, loader resourceStatsLoader) ResourceSummary {
	summary := loader(ctx)

	c.mu.Lock()
	defer c.mu.Unlock()
	recordedAt := c.now()
	load.summary = summary
	c.items[key] = resourceStatsCacheEntry{
		summary:    summary,
		freshUntil: recordedAt.Add(c.ttl),
		staleUntil: recordedAt.Add(c.ttl + c.staleWindow),
	}
	delete(c.inflight, key)
	close(load.done)
	return summary
}
