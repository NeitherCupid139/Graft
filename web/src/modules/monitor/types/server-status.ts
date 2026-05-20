import type { MonitorTrendRange } from '../contract/trend';

export interface ServerStatusDependency {
  status: string;
  detail: string;
  latency_ms: number | null;
}

export interface ServerStatusPlugin {
  name: string;
  status: string;
  version: string;
  depends_on: string[];
}

export interface ServerStatusServer {
  version: string;
  started_at: string;
  uptime_seconds: number;
  go_version: string;
  app_name: string;
  app_env: string;
}

export interface ServerStatusRuntime {
  go_version: string;
  host_name: string;
  operating_system: string;
  architecture: string;
  cpu_cores: number;
  goroutines: number;
  alloc_bytes: number;
  heap_in_use_bytes: number;
  system_memory_bytes: number;
  gc_cycles: number;
}

export interface ServerStatusDependencies {
  database: ServerStatusDependency;
  redis: ServerStatusDependency;
}

export interface ServerStatusSummary {
  total_dependencies: number;
  healthy_dependencies: number;
  degraded_dependencies: number;
  unknown_dependencies: number;
  disabled_dependencies: number;
  total_plugins: number;
  healthy_plugins: number;
}

export interface ServerStatusTrendPoint {
  observed_at: string;
  cpu_percent: number;
  goroutines: number;
  alloc_bytes: number;
  heap_in_use_bytes: number;
  system_memory_bytes: number;
}

export interface ServerStatusTrend {
  range: MonitorTrendRange;
  retention_seconds: number;
  sample_interval_seconds: number;
  points: ServerStatusTrendPoint[];
}

export interface ServerStatusResponse {
  status: string;
  observed_at: string;
  server: ServerStatusServer;
  runtime: ServerStatusRuntime;
  dependencies: ServerStatusDependencies;
  summary: ServerStatusSummary;
  trend: ServerStatusTrend;
  plugins: ServerStatusPlugin[];
}
