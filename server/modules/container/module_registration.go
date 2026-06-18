// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"errors"
	"fmt"

	"graft/server/internal/i18n"
	"graft/server/internal/menu"
	"graft/server/internal/permission"
	containercontract "graft/server/modules/container/contract"
)

const (
	operationsMenuOrderRoot = 50
	containerMenuOrderList  = 51
)

func registerMessages(localizer *i18n.Service) error {
	if localizer == nil {
		return errors.New("i18n service is unavailable")
	}

	for _, locale := range []i18n.LocaleTag{i18n.LocaleZHCN, i18n.LocaleENUS} {
		for _, key := range containerLocaleBackedMessageKeys() {
			matches := localizer.RegisteredMessageResources(locale, i18n.MessageKey(key))
			if len(matches) == 0 {
				return fmt.Errorf("register container module messages: locale resource %s missing key %s", locale, key)
			}
		}
	}
	return nil
}

func containerLocaleBackedMessageKeys() []string {
	keys := make([]string, 0, len(containerMessageKeys))
	for _, key := range containerMessageKeys {
		if key == containercontract.OperationsMenuTitle.String() || key == containercontract.ContainerMenuTitle.String() {
			continue
		}
		keys = append(keys, key)
	}
	return keys
}

var containerMessageKeys = []string{
	containercontract.OperationsMenuTitle.String(),
	containercontract.ContainerMenuTitle.String(),
	containercontract.ContainerRuntimeDisabled.String(),
	containercontract.ContainerRuntimeSocketMissing.String(),
	containercontract.ContainerRuntimePermissionDenied.String(),
	containercontract.ContainerRuntimeUnavailable.String(),
	containercontract.ContainerNotFound.String(),
	containercontract.ContainerMountNotFound.String(),
	containercontract.ContainerInvalidRef.String(),
	containercontract.ContainerInvalidListQuery.String(),
	containercontract.ContainerInvalidBatchAction.String(),
	containercontract.ContainerInvalidState.String(),
	containercontract.ContainerLogsTooLarge.String(),
	containercontract.ContainerInvalidLogQuery.String(),
	containercontract.ContainerTimeout.String(),
	containercontract.ContainerMountUsageUnsupported.String(),
	containercontract.ContainerDangerousActionsDisabled.String(),
	containercontract.ContainerActionStartCompleted.String(),
	containercontract.ContainerActionStopCompleted.String(),
	containercontract.ContainerActionRestartCompleted.String(),
	containercontract.ContainerActionRemoveCompleted.String(),
	containercontract.ContainerBatchActionCompleted.String(),
	containercontract.ContainerBatchActionPartial.String(),
	containercontract.ContainerBatchActionFailed.String(),
}

func registerPermissions(registry *permission.Registry, moduleName string) error {
	if registry == nil {
		return errors.New("permission registry is unavailable")
	}

	for _, item := range permissionItems(moduleName) {
		registry.Register(item)
	}
	return nil
}

func permissionItems(moduleName string) []permission.Item {
	return []permission.Item{
		{
			Code:           containercontract.ContainerViewPermission.String(),
			Name:           "View Containers",
			DisplayKey:     "rbac.permissionCatalog.containerView.display",
			Description:    "Allows reading the container list.",
			DescriptionKey: "rbac.permissionCatalog.containerView.description",
			Category:       "api",
			Module:         moduleName,
		},
		{
			Code:           containercontract.ContainerDetailPermission.String(),
			Name:           "View Container Details",
			DisplayKey:     "rbac.permissionCatalog.containerDetail.display",
			Description:    "Allows reading container details.",
			DescriptionKey: "rbac.permissionCatalog.containerDetail.description",
			Category:       "api",
			Module:         moduleName,
		},
		{
			Code:           containercontract.ContainerEnvironmentPermission.String(),
			Name:           "Read Container Environment Values",
			DisplayKey:     "rbac.permissionCatalog.containerEnvironment.display",
			Description:    "Allows reading container environment variable values when the display policy permits them.",
			DescriptionKey: "rbac.permissionCatalog.containerEnvironment.description",
			Category:       "api",
			Module:         moduleName,
		},
		{
			Code:           containercontract.ContainerLogsPermission.String(),
			Name:           "Read Container Logs",
			DisplayKey:     "rbac.permissionCatalog.containerLogs.display",
			Description:    "Allows reading bounded container logs.",
			DescriptionKey: "rbac.permissionCatalog.containerLogs.description",
			Category:       "api",
			Module:         moduleName,
		},
		{
			Code:           containercontract.ContainerStartPermission.String(),
			Name:           "Start Containers",
			DisplayKey:     "rbac.permissionCatalog.containerStart.display",
			Description:    "Allows starting containers when dangerous actions are enabled.",
			DescriptionKey: "rbac.permissionCatalog.containerStart.description",
			Category:       "api",
			Module:         moduleName,
		},
		{
			Code:           containercontract.ContainerStopPermission.String(),
			Name:           "Stop Containers",
			DisplayKey:     "rbac.permissionCatalog.containerStop.display",
			Description:    "Allows stopping containers when dangerous actions are enabled.",
			DescriptionKey: "rbac.permissionCatalog.containerStop.description",
			Category:       "api",
			Module:         moduleName,
		},
		{
			Code:           containercontract.ContainerRestartPermission.String(),
			Name:           "Restart Containers",
			DisplayKey:     "rbac.permissionCatalog.containerRestart.display",
			Description:    "Allows restarting containers when dangerous actions are enabled.",
			DescriptionKey: "rbac.permissionCatalog.containerRestart.description",
			Category:       "api",
			Module:         moduleName,
		},
		{
			Code:           containercontract.ContainerRemovePermission.String(),
			Name:           "Remove Containers",
			DisplayKey:     "rbac.permissionCatalog.containerRemove.display",
			Description:    "Allows removing containers when dangerous actions are enabled.",
			DescriptionKey: "rbac.permissionCatalog.containerRemove.description",
			Category:       "api",
			Module:         moduleName,
		},
	}
}

func registerMenu(registry *menu.Registry, moduleName string) error {
	if registry == nil {
		return errors.New("menu registry is unavailable")
	}

	registry.Register(menu.Item{
		Code:       "ops.root",
		Title:      "",
		TitleKey:   containercontract.OperationsMenuTitle.String(),
		Path:       containercontract.ContainerMenuRootPath,
		Icon:       "tools",
		Order:      operationsMenuOrderRoot,
		Permission: "",
		Module:     moduleName,
	})
	registry.Register(menu.Item{
		Code:                     "container.list",
		Title:                    "",
		TitleKey:                 containercontract.ContainerMenuTitle.String(),
		Path:                     containercontract.ContainerMenuPath,
		Icon:                     "server",
		Order:                    containerMenuOrderList,
		Permission:               containercontract.ContainerViewPermission.String(),
		VisibleWhenConfigEnabled: containercontract.ContainerRuntimeEnabledConfig.String(),
		Module:                   moduleName,
	})
	return nil
}
