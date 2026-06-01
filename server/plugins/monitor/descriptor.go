package monitor

import "graft/server/internal/plugin"

const (
	moduleID      = "monitor"
	moduleVersion = "0.1.0"
)

var moduleDependencies = []string{"user", "rbac"}

// NewModuleSpec exposes the monitor module's stable compile-time metadata and builder.
func NewModuleSpec() plugin.ModuleSpec {
	return plugin.ModuleSpec{
		ID:            moduleID,
		ModuleVersion: moduleVersion,
		Dependencies:  append([]string(nil), moduleDependencies...),
		Builder: plugin.BuilderFunc(func(plugin.BuildContext) (plugin.Plugin, error) {
			return NewPlugin(), nil
		}),
	}
}
