// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { Ref } from 'vue';
import type { Router } from 'vue-router';

import { formatLocaleDateTime, MEDIUM_DATE_TIME_FORMAT_OPTIONS } from '@/shared/observability';

export function openDashboardRoute(router: Router, location: string) {
  void router.push(location);
}

export function formatDashboardDateTime(value: string, locale?: string | Ref<string | undefined> | null) {
  return formatLocaleDateTime(value, locale, MEDIUM_DATE_TIME_FORMAT_OPTIONS);
}
