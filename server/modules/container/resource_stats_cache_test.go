package container

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

//nolint:unused // The tested stale-refresh path is intentionally kept for cache-governance coverage.
func TestResourceStatsCacheReturnsStaleAndRefreshesInBackground(t *testing.T) {
	t.Parallel()

	start := time.Unix(1_700_000_000, 0)
	now := start
	cache := newResourceStatsCache(2*time.Second, 5*time.Second)
	cache.now = func() time.Time { return now }

	var calls atomic.Int64
	loader := func(context.Context) ResourceSummary {
		value := calls.Add(1)
		return fullResourceSummary(float64(value), float64(value*10))
	}

	first := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, first.CPUPercent, 1, "initial cpu percent")
	assertFloatPtr(t, first.MemoryPercent, 10, "initial memory percent")
	if calls.Load() != 1 {
		t.Fatalf("expected first load to call loader once, got %d", calls.Load())
	}

	now = now.Add(3 * time.Second)
	stale := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, stale.CPUPercent, 1, "stale cpu percent")
	assertFloatPtr(t, stale.MemoryPercent, 10, "stale memory percent")

	waitForCalls(t, &calls, 2)

	fresh := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, fresh.CPUPercent, 2, "refreshed cpu percent")
	assertFloatPtr(t, fresh.MemoryPercent, 20, "refreshed memory percent")
	if calls.Load() != 2 {
		t.Fatalf("expected refreshed cache hit to avoid extra load, got %d", calls.Load())
	}
}

//nolint:unused // The tested stale-refresh path is intentionally kept for cache-governance coverage.
func TestResourceStatsCacheRefreshFailurePreservesLastSuccessWithinStaleWindow(t *testing.T) {
	t.Parallel()

	start := time.Unix(1_700_000_000, 0)
	now := start
	cache := newResourceStatsCache(2*time.Second, 5*time.Second)
	cache.now = func() time.Time { return now }

	var calls atomic.Int64
	loader := func(context.Context) ResourceSummary {
		switch calls.Add(1) {
		case 1:
			return fullResourceSummary(0.6, 12)
		default:
			return unavailableResourceSummary(containerStatsTimeoutReason)
		}
	}

	first := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, first.CPUPercent, 0.6, "initial cpu percent")
	assertFloatPtr(t, first.MemoryPercent, 12, "initial memory percent")

	now = now.Add(3 * time.Second)
	stale := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, stale.CPUPercent, 0.6, "stale cpu percent")
	assertFloatPtr(t, stale.MemoryPercent, 12, "stale memory percent")

	waitForCalls(t, &calls, 2)

	afterFailure := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, afterFailure.CPUPercent, 0.6, "preserved cpu percent")
	assertFloatPtr(t, afterFailure.MemoryPercent, 12, "preserved memory percent")
	if !cache.items["container-1"].updatedAt.Equal(start) {
		t.Fatalf("expected last successful snapshot time to remain anchored at %s, got %s", start, cache.items["container-1"].updatedAt)
	}
}

//nolint:unused // The tested stale-refresh path is intentionally kept for cache-governance coverage.
func TestResourceStatsCachePartialRefreshPreservesLastFullSnapshot(t *testing.T) {
	t.Parallel()

	start := time.Unix(1_700_000_000, 0)
	now := start
	cache := newResourceStatsCache(2*time.Second, 5*time.Second)
	cache.now = func() time.Time { return now }

	var calls atomic.Int64
	loader := func(context.Context) ResourceSummary {
		switch calls.Add(1) {
		case 1:
			return fullResourceSummary(0.6, 12)
		default:
			return partialResourceSummary(cpuOnlySummary(0.8))
		}
	}

	first := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, first.CPUPercent, 0.6, "initial cpu percent")
	assertFloatPtr(t, first.MemoryPercent, 12, "initial memory percent")

	now = now.Add(3 * time.Second)
	stale := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, stale.CPUPercent, 0.6, "stale cpu percent")
	assertFloatPtr(t, stale.MemoryPercent, 12, "stale memory percent")

	waitForCalls(t, &calls, 2)

	afterPartial := cache.get(context.Background(), "container-1", loader)
	assertFloatPtr(t, afterPartial.CPUPercent, 0.6, "preserved cpu percent")
	assertFloatPtr(t, afterPartial.MemoryPercent, 12, "preserved memory percent")
}

//nolint:unused // The tested stale-refresh path is intentionally kept for cache-governance coverage.
func TestResourceStatsCacheWithoutPriorSnapshotReturnsUnavailableForPartialResult(t *testing.T) {
	t.Parallel()

	cache := newResourceStatsCache(2*time.Second, 5*time.Second)
	summary := cache.get(context.Background(), "container-1", func(context.Context) ResourceSummary {
		return partialResourceSummary(memoryOnlySummary(12))
	})

	if summary.Available || summary.StatsAvailable {
		t.Fatalf("expected partial result without prior snapshot to be unavailable, got %#v", summary)
	}
	if summary.CPUPercent != nil || summary.MemoryPercent != nil {
		t.Fatalf("expected unavailable result to avoid half snapshot, got %#v", summary)
	}
	if summary.UnavailableReason != containerStatsIncompleteReason {
		t.Fatalf("expected incomplete reason, got %q", summary.UnavailableReason)
	}
	if _, ok := cache.items["container-1"]; ok {
		t.Fatalf("expected partial result without prior snapshot to avoid caching")
	}
}

func fullResourceSummary(cpu float64, memoryPercent float64) ResourceSummary {
	return ResourceSummary{
		Available:      true,
		StatsAvailable: true,
		CPUPercent:     float64Ptr(cpu),
		MemoryPercent:  float64Ptr(memoryPercent),
	}
}

func cpuOnlySummary(cpu float64) ResourceSummary {
	return ResourceSummary{
		Available:      true,
		StatsAvailable: true,
		CPUPercent:     float64Ptr(cpu),
	}
}

func memoryOnlySummary(memoryPercent float64) ResourceSummary {
	return ResourceSummary{
		Available:      true,
		StatsAvailable: true,
		MemoryPercent:  float64Ptr(memoryPercent),
	}
}

func partialResourceSummary(summary ResourceSummary) ResourceSummary {
	summary.StatsErrorKey = "should-not-survive"
	summary.StatsErrorMessage = "should-not-survive"
	return summary
}

func waitForCalls(t *testing.T, calls *atomic.Int64, expected int64) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if calls.Load() >= expected {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("expected loader call count %d, got %d", expected, calls.Load())
}
