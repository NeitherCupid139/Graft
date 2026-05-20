export interface ServerStatusDependency {
  status: string;
}

export interface ServerStatusPlugin {
  name: string;
  status: string;
  version: string;
}

export interface ServerStatusServer {
  version: string;
  started_at: string;
  uptime_seconds: number;
  go_version: string;
  app_name: string;
  app_env: string;
}

export interface ServerStatusDependencies {
  database: ServerStatusDependency;
  redis: ServerStatusDependency;
}

export interface ServerStatusResponse {
  status: string;
  observed_at: string;
  server: ServerStatusServer;
  dependencies: ServerStatusDependencies;
  plugins: ServerStatusPlugin[];
}
