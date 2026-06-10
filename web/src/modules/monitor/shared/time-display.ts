// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { Ref } from 'vue';

import { formatLocaleDateOnly, formatLocaleTimeOnly } from '@/shared/observability';

export function formatTimeOnly(value?: string | null, locale?: string | Ref<string | undefined> | null) {
  if (!value || Number.isNaN(new Date(value).getTime())) {
    return '--';
  }

  const formatted = formatLocaleTimeOnly(value, locale);
  return formatted === '-' ? '--' : formatted;
}

export function formatDateOnly(value?: string | null, locale?: string | Ref<string | undefined> | null) {
  if (!value || Number.isNaN(new Date(value).getTime())) {
    return '';
  }

  return formatLocaleDateOnly(value, locale);
}
