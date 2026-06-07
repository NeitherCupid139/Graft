package contract

// PermissionCode identifies a stable system configuration permission contract.
type PermissionCode string

// String returns the wire-format permission code.
func (c PermissionCode) String() string {
	return string(c)
}

const (
	SystemConfigReadPermission  PermissionCode = "system-config.read"
	SystemConfigWritePermission PermissionCode = "system-config.write"
)
