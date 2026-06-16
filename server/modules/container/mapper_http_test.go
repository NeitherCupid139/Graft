// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package container

import "testing"

func TestToResourceSummaryMapsDockerStatsFields(t *testing.T) {
	t.Parallel()

	resource := ResourceSummary{
		Available:                  true,
		StatsAvailable:             true,
		CPUPercent:                 float64Ptr(12.5),
		OnlineCPUs:                 int64Ptr(4),
		SystemCPUUsage:             int64Ptr(1000),
		TotalCPUUsage:              int64Ptr(200),
		CPUUsageInUsermode:         int64Ptr(70),
		CPUUsageInKernelmode:       int64Ptr(30),
		ThrottlingPeriods:          int64Ptr(11),
		ThrottlingThrottledPeriods: int64Ptr(3),
		ThrottlingThrottledTime:    int64Ptr(900),
		MemoryUsageBytes:           int64Ptr(256),
		MemoryLimitBytes:           int64Ptr(1024),
		MemoryPercent:              float64Ptr(25),
		MemoryCache:                int64Ptr(10),
		MemoryRSS:                  int64Ptr(20),
		MemoryActiveFile:           int64Ptr(30),
		MemoryInactiveFile:         int64Ptr(40),
		MemoryPgfault:              int64Ptr(50),
		MemoryPgmajfault:           int64Ptr(60),
		RxBytes:                    int64Ptr(107),
		TxBytes:                    int64Ptr(208),
		RxPackets:                  int64Ptr(12),
		TxPackets:                  int64Ptr(14),
		RxErrors:                   int64Ptr(12),
		TxErrors:                   int64Ptr(14),
		RxDropped:                  int64Ptr(18),
		TxDropped:                  int64Ptr(20),
		PIDsCurrent:                int64Ptr(5),
		PIDsLimit:                  int64Ptr(128),
	}

	mapped := toResourceSummary(resource)
	if mapped == nil {
		t.Fatalf("expected mapped resource summary")
	}
	assertFloatPtr(t, mapped.CpuPercent, 12.5, "mapped CPU percent")
	assertInt64Ptr(t, mapped.OnlineCpus, 4, "mapped online CPUs")
	assertInt64Ptr(t, mapped.SystemCpuUsage, 1000, "mapped system CPU usage")
	assertInt64Ptr(t, mapped.TotalCpuUsage, 200, "mapped total CPU usage")
	assertInt64Ptr(t, mapped.CpuUsageInUsermode, 70, "mapped CPU user mode usage")
	assertInt64Ptr(t, mapped.CpuUsageInKernelmode, 30, "mapped CPU kernel mode usage")
	assertInt64Ptr(t, mapped.ThrottlingPeriods, 11, "mapped CPU throttling periods")
	assertInt64Ptr(t, mapped.ThrottlingThrottledPeriods, 3, "mapped CPU throttled periods")
	assertInt64Ptr(t, mapped.ThrottlingThrottledTime, 900, "mapped CPU throttled time")
	assertInt64Ptr(t, mapped.MemoryUsageBytes, 256, "mapped memory usage bytes")
	assertInt64Ptr(t, mapped.MemoryLimitBytes, 1024, "mapped memory limit bytes")
	assertFloatPtr(t, mapped.MemoryPercent, 25, "mapped memory percent")
	assertInt64Ptr(t, mapped.MemoryCache, 10, "mapped memory cache")
	assertInt64Ptr(t, mapped.MemoryRss, 20, "mapped memory rss")
	assertInt64Ptr(t, mapped.MemoryActiveFile, 30, "mapped memory active file")
	assertInt64Ptr(t, mapped.MemoryInactiveFile, 40, "mapped memory inactive file")
	assertInt64Ptr(t, mapped.MemoryPgfault, 50, "mapped memory pgfault")
	assertInt64Ptr(t, mapped.MemoryPgmajfault, 60, "mapped memory pgmajfault")
	assertInt64Ptr(t, mapped.RxBytes, 107, "mapped rx bytes")
	assertInt64Ptr(t, mapped.TxBytes, 208, "mapped tx bytes")
	assertInt64Ptr(t, mapped.RxPackets, 12, "mapped rx packets")
	assertInt64Ptr(t, mapped.TxPackets, 14, "mapped tx packets")
	assertInt64Ptr(t, mapped.RxErrors, 12, "mapped rx errors")
	assertInt64Ptr(t, mapped.TxErrors, 14, "mapped tx errors")
	assertInt64Ptr(t, mapped.RxDropped, 18, "mapped rx dropped")
	assertInt64Ptr(t, mapped.TxDropped, 20, "mapped tx dropped")
	assertInt64Ptr(t, mapped.PidsCurrent, 5, "mapped pids current")
	assertInt64Ptr(t, mapped.PidsLimit, 128, "mapped pids limit")
}

func int64Ptr(value int64) *int64 {
	return &value
}

func float64Ptr(value float64) *float64 {
	return &value
}
