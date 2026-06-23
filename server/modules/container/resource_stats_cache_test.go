package container

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestResourceStatsCacheReturnsStaleAndRefreshesInBackground(t *testing.T) {
	t.Parallel()

	start := time.Unix(1_700_000_000, 0)
	now := start
	cache := newResourceStatsCache(2*time.Second, 5*time.Second)
	cache.now = func() time.Time { return now }

	var calls atomic.Int64
	loader := func(context.Context) ResourceSummary {
		value := calls.Add(1)
		memory := value * 100
		return ResourceSummary{
			Available:        true,
			StatsAvailable:   true,
			MemoryUsageBytes: &memory,
		}
	}

	first := cache.get(context.Background(), "container-1", loader)
	assertInt64Ptr(t, first.MemoryUsageBytes, 100, "initial memory usage")
	if calls.Load() != 1 {
		t.Fatalf("expected first load to call loader once, got %d", calls.Load())
	}

	now = now.Add(3 * time.Second)
	stale := cache.get(context.Background(), "container-1", loader)
	assertInt64Ptr(t, stale.MemoryUsageBytes, 100, "stale memory usage")

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if calls.Load() >= 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected stale read to trigger one background refresh, got %d", calls.Load())
	}

	fresh := cache.get(context.Background(), "container-1", loader)
	assertInt64Ptr(t, fresh.MemoryUsageBytes, 200, "refreshed memory usage")
	if calls.Load() != 2 {
		t.Fatalf("expected refreshed cache hit to avoid extra load, got %d", calls.Load())
	}
}
