package scheduler

import "graft/server/internal/plugin"

// NewModuleSpec exposes the scheduler module's stable compile-time metadata and builder.
func NewModuleSpec() plugin.ModuleSpec {
	instance := NewPlugin()

	return plugin.ModuleSpec{
		ID:            instance.Name(),
		ModuleVersion: instance.Version(),
		Dependencies:  append([]string(nil), instance.DependsOn()...),
		Builder:       plugin.BuilderFunc(func(plugin.BuildContext) (plugin.Plugin, error) { return NewPlugin(), nil }),
	}
}
