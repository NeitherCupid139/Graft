package contract

// MessageKey identifies a stable system configuration message key.
type MessageKey string

// String returns the canonical message key value.
func (k MessageKey) String() string {
	return string(k)
}

const (
	SystemConfigMenuTitle      MessageKey = "menu.server.system_config.title"
	SystemConfigNotFound       MessageKey = "system_config.not_found"
	SystemConfigInvalidRequest MessageKey = "system_config.invalid_request"
)
