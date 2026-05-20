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
)
