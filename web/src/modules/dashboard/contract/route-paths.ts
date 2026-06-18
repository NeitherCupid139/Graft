// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

export function normalizeDashboardRoutePath(path: string) {
  const trimmed = path.trim();
  if (!trimmed) {
    return '';
  }

  const pathOnly = trimmed.split(/[?#]/, 1)[0] ?? '';
  const withLeadingSlash = pathOnly.startsWith('/') ? pathOnly : `/${pathOnly}`;
  return withLeadingSlash === '/' ? withLeadingSlash : withLeadingSlash.replace(/\/+$/, '');
}

export function normalizeJoinedDashboardRoutePath(parentPath: string, routePath: string) {
  if (routePath.startsWith('/')) {
    return normalizeDashboardRoutePath(routePath);
  }

  if (!routePath) {
    return normalizeDashboardRoutePath(parentPath);
  }

  return normalizeDashboardRoutePath(`${parentPath}/${routePath}`);
}
