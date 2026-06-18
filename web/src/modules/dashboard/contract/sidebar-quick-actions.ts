// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { RouteRecordRaw } from 'vue-router';

import { getDefaultLocale, type SupportedLocale } from '@/contracts/i18n/locales';
import { renderLocalizedTitle, resolveRouteLocalizedTitle } from '@/utils/route/meta';
import type { AppRouteMeta } from '@/utils/types';

import type { DashboardQuickActionLink } from './quick-action-links';
import { normalizeDashboardRoutePath, normalizeJoinedDashboardRoutePath } from './route-paths';

type QuickActionSource = Pick<RouteRecordRaw, 'path' | 'children' | 'name' | 'meta'>;
type QuickActionParent = {
  groupKey?: string;
  groupLabel?: string;
};

/**
 * Builds a sorted list of dashboard quick action links from Vue Router routes.
 *
 * @param routes - The array of Vue Router route records to extract quick action links from.
 * @returns An array of DashboardQuickActionLink objects sorted by order (ascending), then by id (lexicographically).
 */
export function buildDashboardQuickActionLinks(routes: RouteRecordRaw[], locale: SupportedLocale = getDefaultLocale()) {
  return collectLeafLinks(routes, locale).sort(compareQuickActions);
}

/**
 * Recursively collects visible leaf routes and transforms them into dashboard quick action links.
 *
 * @returns An array of dashboard quick action links for visible leaf routes
 */
function collectLeafLinks(
  routes: QuickActionSource[],
  locale: SupportedLocale,
  parentPath = '',
  parent?: QuickActionParent,
): DashboardQuickActionLink[] {
  return routes.flatMap((route) => {
    const routeMeta = toRouteMeta(route.meta);
    const fullPath = normalizeJoinedDashboardRoutePath(parentPath, String(route.path ?? ''));
    if (!fullPath || routeMeta?.hidden || routeMeta?.hiddenMenu) {
      return [];
    }

    const nextParent = resolveParent(routeMeta, locale, parent);
    const visibleChildren = collectLeafLinks(route.children ?? [], locale, fullPath, nextParent);
    if (visibleChildren.length > 0) {
      return visibleChildren;
    }

    if (!isQuickActionLeaf(route, fullPath)) {
      return [];
    }

    const title =
      renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'breadcrumb'), locale) ||
      renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'page'), locale) ||
      fullPath;
    const fullLabel = renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'tab'), locale) || title;
    const routeGroup =
      renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'page'), locale) || parent?.groupLabel;
    const isSingleLeaf = Boolean(routeMeta?.single);

    return [
      {
        id: String(route.name ?? fullPath),
        full_label: fullLabel,
        group: isSingleLeaf ? routeGroup : parent?.groupLabel,
        group_key: isSingleLeaf ? routeMeta?.titleKey : parent?.groupKey,
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

/**
 * Determines if a route is eligible as a quick-action leaf node.
 *
 * @returns `true` if the route has a name and valid path, is not hidden, and either is marked as single or has no children, `false` otherwise.
 */
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

/**
 * Types route metadata as AppRouteMeta.
 *
 * @param meta - The raw metadata value
 * @returns The metadata as `AppRouteMeta`, or `undefined` if the input was nullish
 */
function toRouteMeta(meta: unknown) {
  return (meta ?? undefined) as AppRouteMeta | undefined;
}

/**
 * Resolves grouping context for a route based on its metadata and locale.
 *
 * @returns A grouping context with `groupLabel` and `groupKey`, or the parent context if no new values are available
 */
function resolveParent(routeMeta: AppRouteMeta | undefined, locale: SupportedLocale, parent?: QuickActionParent) {
  const groupLabel =
    renderLocalizedTitle(resolveRouteLocalizedTitle(routeMeta, 'page'), locale) ||
    renderLocalizedTitle(routeMeta?.title, locale) ||
    parent?.groupLabel;
  const groupKey = routeMeta?.titleKey || parent?.groupKey;

  if (!groupLabel && !groupKey) {
    return parent;
  }

  return {
    groupKey,
    groupLabel,
  };
}

/**
 * Derives a module key from a dashboard route path.
 *
 * Maps route paths to module keys based on path structure:
 * - 'access-control/overview' → 'access-control'
 * - 'access-control/*' → 'rbac'
 * - 'logs/access' → 'access-log'
 * - 'logs/app' → 'app-log'
 * - 'logs/*' → 'logs'
 * - Other paths → first path component
 * - Empty path → 'dashboard'
 *
 * @param path - The dashboard route path
 * @returns The module key derived from the path
 */
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

/**
 * Orders two dashboard quick action links by their order field, then by their ID.
 *
 * @returns A negative number if `left` should come before `right`, zero if equal, positive otherwise.
 */
function compareQuickActions(left: DashboardQuickActionLink, right: DashboardQuickActionLink) {
  if (left.order !== right.order) {
    return left.order - right.order;
  }

  return left.id.localeCompare(right.id);
}
