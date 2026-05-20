package contract

// MenuMessageKey identifies a stable monitor-plugin menu title message key.
type MenuMessageKey string

// String returns the canonical menu message key value.
func (k MenuMessageKey) String() string {
	return string(k)
}

const (
	// ServerStatusMenuTitle identifies the localized title for the server-status menu.
	ServerStatusMenuTitle MenuMessageKey = "menu.monitor.server_status.title"
)
