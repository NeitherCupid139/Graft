// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package contract

// ConfigKey identifies a stable container management system configuration key.
type ConfigKey string

// String returns the canonical config key.
func (k ConfigKey) String() string {
	return string(k)
}

const (
	// ContainerEnabledConfig enables the container management module.
	ContainerEnabledConfig ConfigKey = "ops.container.enabled"
	// ContainerRuntimeConfig selects the runtime adapter.
	ContainerRuntimeConfig ConfigKey = "ops.container.runtime"
	// ContainerDockerEndpointConfig stores the first local runtime endpoint.
	ContainerDockerEndpointConfig ConfigKey = "ops.container.docker.endpoint"
	// ContainerLogsDefaultTailConfig stores the default log tail size.
	ContainerLogsDefaultTailConfig ConfigKey = "ops.container.logs.default_tail"
	// ContainerLogsMaxTailConfig stores the maximum log tail size.
	ContainerLogsMaxTailConfig ConfigKey = "ops.container.logs.max_tail"
	// ContainerDangerousActionsEnabledConfig enables start, stop, and restart actions.
	ContainerDangerousActionsEnabledConfig ConfigKey = "ops.container.dangerous_actions_enabled"
)
