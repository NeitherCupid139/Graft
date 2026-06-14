// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"errors"

	"graft/server/internal/module"
)

// Module declares the container management module foundation.
type Module struct{}

// NewModule creates a container management module instance.
func NewModule() *Module {
	return &Module{}
}

// Register declares container menu, permissions, messages, and system config definitions.
func (m *Module) Register(ctx *module.Context) error {
	if m == nil {
		return errors.New("container module is unavailable")
	}
	if err := registerMessages(ctx.I18n); err != nil {
		return err
	}
	if err := registerPermissions(ctx.PermissionRegistry, moduleID); err != nil {
		return err
	}
	if err := registerMenu(ctx.MenuRegistry, moduleID); err != nil {
		return err
	}
	return registerConfig(ctx.I18n, ctx.ConfigRegistry)
}

// Boot currently has no runtime work; DockerRuntime is implemented in a later phase.
func (m *Module) Boot(_ *module.Context) error {
	return nil
}

// Shutdown currently has no runtime resources to release.
func (m *Module) Shutdown(_ *module.Context) error {
	return nil
}
