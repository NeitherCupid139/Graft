package systemconfig

import (
	"errors"

	"graft/server/internal/module"
)

const moduleID = "system-config"

// Module owns system configuration administrator overrides and HTTP management.
type Module struct {
	service *Service
}

// NewModule creates the system configuration module.
func NewModule(service *Service) (*Module, error) {
	if service == nil {
		return nil, errors.New("system config service is unavailable")
	}
	return &Module{service: service}, nil
}

// Register declares permissions, menu, messages, and management routes.
func (m *Module) Register(ctx *module.Context) error {
	if m == nil || m.service == nil {
		return errors.New("system config module service is unavailable")
	}
	if err := registerMessages(ctx.I18n); err != nil {
		return err
	}
	if err := registerSystemConfigPermissions(ctx.PermissionRegistry, moduleID); err != nil {
		return err
	}
	if err := registerSystemConfigMenu(ctx.MenuRegistry, moduleID); err != nil {
		return err
	}
	return registerSystemConfigRoutes(ctx, moduleID, m.service)
}

// Boot currently has no runtime work; definitions are registered by owner modules.
func (m *Module) Boot(_ *module.Context) error {
	return nil
}

// Shutdown currently has no resources to release.
func (m *Module) Shutdown(_ *module.Context) error {
	return nil
}
