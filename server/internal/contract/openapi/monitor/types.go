package monitor

// ServerInterface is the minimal monitor-only generated handler contract used by this spike.
type ServerInterface interface {
	GetMonitorServerStatus(params GetMonitorServerStatusParams)
}
