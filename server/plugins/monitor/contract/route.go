package contract

// JoinRoute combines a route group path with a route fragment.
func JoinRoute(group, fragment string) string {
	return group + fragment
}

const (
	// MonitorGroup identifies the monitor route group.
	MonitorGroup = "/monitor"

	// ServerStatusRoute identifies the server-status route fragment.
	ServerStatusRoute = "/server-status"

	// OverviewRoute identifies the overview route fragment under server-status.
	OverviewRoute = "/overview"

	// ServerStatusMenuPath identifies the second-level server-status menu path.
	ServerStatusMenuPath = MonitorGroup + ServerStatusRoute

	// ServerStatusOverviewMenuPath identifies the third-level overview menu path.
	ServerStatusOverviewMenuPath = ServerStatusMenuPath + OverviewRoute
)
