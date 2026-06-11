// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { describe, expect, it } from 'vitest';
import type { ComposerTranslation } from 'vue-i18n';

import type { NotificationItem } from '../types/notification';
import {
  resolveNotificationActionLabel,
  resolveNotificationCategory,
  resolveNotificationMessage,
  resolveNotificationResourceType,
  resolveNotificationTitle,
} from './presentation';

const messages: Record<string, string> = {
  'notification.category.task': '任务',
  'notification.resource.scheduledTaskRun': '定时任务运行记录',
  'notification.scheduler.taskSucceeded.action': '打开定时任务运行记录',
  'notification.scheduler.taskSucceeded.message': '{taskName} 已执行成功。',
  'notification.scheduler.taskSucceeded.title': '定时任务执行成功',
};

const t = ((key: string, context?: Record<string, unknown>) => {
  const template = messages[key];
  if (!template) return key;
  return template.replaceAll(/\{(\w+)\}/g, (_, name: string) => String(context?.[name] ?? ''));
}) as ComposerTranslation;

function notification(overrides: Partial<NotificationItem> = {}): NotificationItem {
  return {
    category: 'TASK',
    delivery_created_at: '2026-06-11T10:47:21Z',
    delivery_id: 1,
    event_id: 1,
    event_type: 'task_succeeded',
    message: 'Scheduled task Access log retention cleanup succeeded.',
    navigation: { kind: 'SCHEDULER_RUN', payload: {} },
    occurred_at: '2026-06-11T10:47:21Z',
    severity: 'info',
    source_module: 'scheduler',
    status: 'unread',
    target_ref: '1',
    target_type: 'USER',
    title: 'Scheduled task succeeded',
    ...overrides,
  };
}

describe('notification presentation resolver', () => {
  it('prefers key-first localized title, message, category, resource, and action labels', () => {
    const item = notification({
      action_label_key: 'notification.scheduler.taskSucceeded.action',
      category_key: 'notification.category.task',
      context: { taskName: '访问日志保留清理' },
      message_key: 'notification.scheduler.taskSucceeded.message',
      resource_type: 'scheduled_task_run',
      title_key: 'notification.scheduler.taskSucceeded.title',
    });

    expect(resolveNotificationTitle(item, t)).toBe('定时任务执行成功');
    expect(resolveNotificationMessage(item, t)).toBe('访问日志保留清理 已执行成功。');
    expect(resolveNotificationCategory(item, t)).toBe('任务');
    expect(resolveNotificationResourceType(item, t)).toBe('定时任务运行记录');
    expect(resolveNotificationActionLabel(item, t)).toBe('打开定时任务运行记录');
  });

  it('falls back to backend title and message when keys are missing', () => {
    const item = notification();

    expect(resolveNotificationTitle(item, t)).toBe('Scheduled task succeeded');
    expect(resolveNotificationMessage(item, t)).toBe('Scheduled task Access log retention cleanup succeeded.');
  });
});
