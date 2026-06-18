// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { RouteRecordRaw } from 'vue-router';

import { getDefaultLocale, type SupportedLocale } from '@/contracts/i18n/locales';
import { renderLocalizedTitle, resolveRouteLocalizedTitle } from '@/utils/route/meta';
import type { AppRouteMeta } from '@/utils/types';

import type { DashboardQuickActionLink } from './quick-action-links';
import { normalizeDashboardRoutePath, normalizeJoinedDashboardRoutePath } from './route-paths';

type QuickActionSource = Pick<RouteRecordRaw, 'path' | 'children' | 'name' | 'meta'>;

export function buildDashboardQuickActionLinks(routes: RouteRecordRaw[], locale: SupportedLocale = getDefaultLocale()) {
  return collectLeafLinks(routes, locale).sort(compareQuickActions);
}

function collectLeafLinks(
  routes: QuickActionSource[],
  locale: SupportedLocale,
  parentPath = '',
): DashboardQuickActionLink[] {
  return routes.flatMap((route) => {
    const routeMeta = toRouteMeta(route.meta);
    const fullPath = normalizeJoinedDashboardRoutePath(parentPath, String(route.path ?? ''));
    if (!fullPath || routeMeta?.hidden || routeMeta?.hiddenMenu) {
      return [];
    }

    const visibleChildren = collectLeafLinks(route.children ?? [], locale, fullPath);
    if (visibleChildren.length > 0) {
      return visibleChildren;
    }

    if (!isQuickActionLeaf(route, fullPath)) {
      return [];
    }

    const title =
      renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'tab'), locale) ||
      renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'page'), locale) ||
      renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'breadcrumb'), locale) ||
      fullPath;

    return [
      {
        id: String(route.name ?? fullPath),
        module_key: inferModuleKey(fullPath),
        icon: typeof routeMeta?.icon === 'string' ? routeMeta.icon : undefined,
        order: routeMeta?.orderNo ?? 0,
        route_location: fullPath,
        title,
        title_key: routeMeta?.titleKey,
      },
    ];
  });
}

function isQuickActionLeaf(route: QuickActionSource, fullPath: string) {
  const routeMeta = toRouteMeta(route.meta);
  if (!route.name || !fullPath) {
    return false;
  }

  if (routeMeta?.single) {
    return !routeMeta.hidden && !routeMeta.hiddenMenu;
  }

  if ((route.children?.length ?? 0) > 0) {
    return false;
  }

  return !routeMeta?.hidden && !routeMeta?.hiddenMenu;
}

function toRouteMeta(meta: unknown) {
  return (meta ?? undefined) as AppRouteMeta | undefined;
}

function inferModuleKey(path: string) {
  const segments = normalizeDashboardRoutePath(path).split('/').filter(Boolean);
  if (segments.length === 0) {
    return 'dashboard';
  }

  const [first, second] = segments;
  if (first === 'access-control') {
    return second === 'overview' ? 'access-control' : 'rbac';
  }

  if (first === 'logs') {
    return second === 'access' ? 'access-log' : second === 'app' ? 'app-log' : first;
  }

  return first;
}

function compareQuickActions(left: DashboardQuickActionLink, right: DashboardQuickActionLink) {
  if (left.order !== right.order) {
    return left.order - right.order;
  }

  return left.id.localeCompare(right.id);
}
