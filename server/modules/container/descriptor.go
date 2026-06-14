// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package container

import "graft/server/internal/module"

const moduleID = "container"

// NewModuleSpec exposes the container module's stable compile-time metadata and builder.
func NewModuleSpec() module.Spec {
	return module.Spec{
		ID:           moduleID,
		Dependencies: []string{"user", "auth", "rbac", "system-config"},
		Builder: module.BuilderFunc(func(module.BuildContext) (module.Module, error) {
			return NewModule(), nil
		}),
	}
}
