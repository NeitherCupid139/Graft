// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { ComposerTranslation } from 'vue-i18n';

import type { NotificationItem } from '../types/notification';

export const NOTIFICATION_MVP_SOURCE_MODULES = [
  'notification',
  'scheduler',
  'audit',
  'system-config',
  'access-log',
] as const;

type NotificationContext = Record<string, unknown>;

const SOURCE_KEY_BY_MODULE: Record<string, string> = {
  'access-log': 'notification.source.accessControl',
  audit: 'notification.source.audit',
  notification: 'notification.source.system',
  scheduler: 'notification.source.scheduler',
  'system-config': 'notification.source.system',
};

const CATEGORY_KEY_BY_VALUE: Record<string, string> = {
  CONFIG: 'notification.category.system',
  OPERATIONS: 'notification.category.system',
  SECURITY: 'notification.category.security',
  SYSTEM: 'notification.category.system',
  TASK: 'notification.category.task',
};

const EVENT_TYPE_KEY_BY_VALUE: Record<string, string> = {
  task_succeeded: 'notification.event.taskSucceeded',
};

const RESOURCE_TYPE_KEY_BY_VALUE: Record<string, string> = {
  scheduled_task_run: 'notification.resource.scheduledTaskRun',
};

const DELIVERY_TYPE_KEY_BY_VALUE: Record<string, string> = {
  USER: 'notification.delivery.user',
};

function translateKey(t: ComposerTranslation, key: string | null | undefined, context?: NotificationContext) {
  const normalized = key?.trim();
  if (!normalized) return '';
  const translated = context ? t(normalized, context) : t(normalized);
  return translated === normalized ? '' : translated;
}

function notificationContext(item: NotificationItem): NotificationContext {
  return item.context && typeof item.context === 'object' ? item.context : {};
}

function fallbackValue(value: unknown, fallback: string) {
  return typeof value === 'string' && value.trim() ? value.trim() : fallback;
}

export function notificationSeverityTheme(severity: NotificationItem['severity']) {
  switch (severity) {
    case 'critical':
    case 'error':
      return 'danger';
    case 'warning':
      return 'warning';
    case 'info':
      return 'primary';
    default:
      return 'default';
  }
}

export function notificationStatusTheme(status: NotificationItem['status']) {
  return status === 'unread' ? 'primary' : 'default';
}

export function resolveNotificationTitle(item: NotificationItem, t: ComposerTranslation) {
  return translateKey(t, item.title_key, notificationContext(item)) || item.title;
}

export function resolveNotificationMessage(item: NotificationItem, t: ComposerTranslation) {
  return translateKey(t, item.message_key, notificationContext(item)) || item.message;
}

export function resolveNotificationCategory(item: NotificationItem, t: ComposerTranslation) {
  return resolveNotificationCategoryValue(item.category, t, item.category_key);
}

export function resolveNotificationCategoryValue(value: string, t: ComposerTranslation, key?: string | null) {
  return (
    translateKey(t, key) ||
    translateKey(t, CATEGORY_KEY_BY_VALUE[value]) ||
    fallbackValue(value, t('notification.values.emptyField'))
  );
}

export function resolveNotificationSource(item: NotificationItem, t: ComposerTranslation) {
  return resolveNotificationSourceValue(item.source_module, t, item.source_key);
}

export function resolveNotificationSourceValue(value: string, t: ComposerTranslation, key?: string | null) {
  return (
    translateKey(t, key) ||
    translateKey(t, SOURCE_KEY_BY_MODULE[value]) ||
    fallbackValue(value, t('notification.values.unknownSource'))
  );
}

export function resolveNotificationLevel(item: NotificationItem, t: ComposerTranslation) {
  return resolveNotificationLevelValue(item.severity, t, item.level_key);
}

export function resolveNotificationLevelValue(value: string, t: ComposerTranslation, key?: string | null) {
  return translateKey(t, key) || translateKey(t, `notification.level.${value}`) || value;
}

export function resolveNotificationStatus(item: NotificationItem, t: ComposerTranslation) {
  return translateKey(t, `notification.status.${item.status}`) || item.status;
}

export function resolveNotificationActionLabel(item: NotificationItem, t: ComposerTranslation) {
  return (
    translateKey(t, item.action_label_key, notificationContext(item)) ||
    fallbackValue(item.action_label, t('notification.action.openBusinessContext'))
  );
}

export function resolveNotificationEventType(item: NotificationItem, t: ComposerTranslation) {
  return (
    translateKey(t, item.event_type_key) ||
    translateKey(t, EVENT_TYPE_KEY_BY_VALUE[item.event_type]) ||
    fallbackValue(item.event_type, t('notification.values.emptyField'))
  );
}

export function resolveNotificationResourceType(item: NotificationItem, t: ComposerTranslation) {
  return (
    translateKey(t, RESOURCE_TYPE_KEY_BY_VALUE[item.resource_type ?? '']) ||
    fallbackValue(item.resource_type, t('notification.values.emptyField'))
  );
}

export function resolveNotificationDeliveryType(item: NotificationItem, t: ComposerTranslation) {
  return (
    translateKey(t, DELIVERY_TYPE_KEY_BY_VALUE[item.target_type]) ||
    fallbackValue(item.target_type, t('notification.values.emptyField'))
  );
}

export function resolveNotificationResultSummary(item: NotificationItem, t: ComposerTranslation) {
  return fallbackValue(notificationContext(item).resultSummary, t('notification.values.emptyField'));
}

export function formatNotificationDiagnosticValue(value: unknown, t: ComposerTranslation) {
  return fallbackValue(value, t('notification.values.emptyField'));
}
