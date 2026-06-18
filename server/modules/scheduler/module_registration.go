// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"errors"
	"fmt"

	"graft/server/internal/i18n"
	"graft/server/internal/menu"
	"graft/server/internal/permission"
	schedulercontract "graft/server/modules/scheduler/contract"
)

const scheduledTaskMenuOrder = 104

func registerMessages(localizer *i18n.Service) error {
	if localizer == nil {
		return errors.New("i18n service is unavailable")
	}

	for _, locale := range []i18n.LocaleTag{i18n.LocaleZHCN, i18n.LocaleENUS} {
		for _, key := range schedulerMessageKeys() {
			matches := localizer.RegisteredMessageResources(locale, i18n.MessageKey(key.String()))
			if len(matches) == 0 {
				return fmt.Errorf("register scheduler module messages: locale resource %s missing key %s", locale, key)
			}
		}
	}

	return nil
}

func schedulerMessageKeys() []schedulercontract.MessageKey {
	return []schedulercontract.MessageKey{
		schedulercontract.ScheduledTaskNotFound,
		schedulercontract.ScheduledTaskAlreadyRunning,
		schedulercontract.ScheduledTaskInvalidRequest,
		schedulercontract.ScheduledTaskRunFailedNotificationTitle,
		schedulercontract.ScheduledTaskRunFailedNotificationMessage,
		schedulercontract.ScheduledTaskRunSucceededNotificationTitle,
		schedulercontract.ScheduledTaskRunSucceededNotificationMessage,
	}
}

func registerSchedulerPermissions(registry *permission.Registry, moduleName string) error {
	if registry == nil {
		return errors.New("permission registry is unavailable")
	}

	registry.Register(permission.Item{
		Code:           schedulercontract.ScheduledTaskReadPermission.String(),
		Name:           "Read Scheduled Tasks",
		DisplayKey:     "rbac.permissionCatalog.scheduledTaskRead.display",
		Description:    "Allows reading scheduled task runtime data and run history.",
		DescriptionKey: "rbac.permissionCatalog.scheduledTaskRead.description",
		Category:       "api",
		Module:         moduleName,
	})
	registry.Register(permission.Item{
		Code:           schedulercontract.ScheduledTaskCreatePermission.String(),
		Name:           "Create Scheduled Tasks",
		DisplayKey:     "rbac.permissionCatalog.scheduledTaskCreate.display",
		Description:    "Allows creating user-managed scheduled task instances.",
		DescriptionKey: "rbac.permissionCatalog.scheduledTaskCreate.description",
		Category:       "api",
		Module:         moduleName,
	})
	registry.Register(permission.Item{
		Code:           schedulercontract.ScheduledTaskUpdatePermission.String(),
		Name:           "Update Scheduled Tasks",
		DisplayKey:     "rbac.permissionCatalog.scheduledTaskUpdate.display",
		Description:    "Allows updating scheduled task definitions.",
		DescriptionKey: "rbac.permissionCatalog.scheduledTaskUpdate.description",
		Category:       "api",
		Module:         moduleName,
	})
	registry.Register(permission.Item{
		Code:           schedulercontract.ScheduledTaskDeletePermission.String(),
		Name:           "Delete Scheduled Tasks",
		DisplayKey:     "rbac.permissionCatalog.scheduledTaskDelete.display",
		Description:    "Allows deleting user-managed scheduled task instances.",
		DescriptionKey: "rbac.permissionCatalog.scheduledTaskDelete.description",
		Category:       "api",
		Module:         moduleName,
	})
	registry.Register(permission.Item{
		Code:           schedulercontract.ScheduledTaskRunPermission.String(),
		Name:           "Run Scheduled Tasks",
		DisplayKey:     "rbac.permissionCatalog.scheduledTaskRun.display",
		Description:    "Allows manually running scheduled task runtime jobs.",
		DescriptionKey: "rbac.permissionCatalog.scheduledTaskRun.description",
		Category:       "api",
		Module:         moduleName,
	})
	registry.Register(permission.Item{
		Code:           schedulercontract.ScheduledTaskEnablePermission.String(),
		Name:           "Enable Scheduled Tasks",
		DisplayKey:     "rbac.permissionCatalog.scheduledTaskEnable.display",
		Description:    "Allows enabling and disabling scheduled tasks.",
		DescriptionKey: "rbac.permissionCatalog.scheduledTaskEnable.description",
		Category:       "api",
		Module:         moduleName,
	})
	return nil
}

func registerSchedulerMenu(registry *menu.Registry, moduleName string) error {
	if registry == nil {
		return errors.New("menu registry is unavailable")
	}

	registry.Register(menu.Item{
		Code:       "scheduled-task.list",
		Title:      "定时任务",
		TitleKey:   schedulercontract.ScheduledTaskMenuTitle.String(),
		Path:       schedulercontract.ScheduledTaskMenuPath,
		Icon:       "time",
		Order:      scheduledTaskMenuOrder,
		Permission: schedulercontract.ScheduledTaskReadPermission.String(),
		Module:     moduleName,
	})
	return nil
}
