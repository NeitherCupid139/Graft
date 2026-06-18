// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { BootstrapRouteRegistration, GlobalRouteRegistration } from '@/modules/types';
import { localizeRouteTitle, localizeRouteTitleKey } from '@/utils/route/title';

import { CONTAINER_BOOTSTRAP_ROUTE } from './contract/bootstrap';

const listRouteTitle = localizeRouteTitleKey('container.route.list.title');
const listBreadcrumbTitle = localizeRouteTitleKey('container.route.list.breadcrumb');
const operationsRouteTitle = localizeRouteTitle('Operations', 'menu.ops.title');
const detailRouteTitle = localizeRouteTitleKey('container.route.detail.title');
const detailBreadcrumbTitle = localizeRouteTitleKey('container.route.detail.breadcrumb');

export const containerBootstrapRouteRegistrations: BootstrapRouteRegistration[] = [
  {
    ...CONTAINER_BOOTSTRAP_ROUTE.LIST,
    loadPage: () => import('./pages/list/index.vue'),
    meta: {
      tabGroup: 'ops',
      pageKind: 'list',
      semanticTitle: operationsRouteTitle,
      breadcrumbTitle: listBreadcrumbTitle,
      tabTitle: listRouteTitle,
      title: operationsRouteTitle,
      titleKey: 'container.route.list.title',
    },
  },
];

export const containerGlobalRouteRegistrations: GlobalRouteRegistration[] = [
  {
    ...CONTAINER_BOOTSTRAP_ROUTE.DETAIL,
    loadPage: () => import('./pages/detail/index.vue'),
    meta: {
      hidden: false,
      hiddenMenu: true,
      keepAlive: false,
      pageKind: 'detail',
      pageSurface: 'form-detail',
      semanticTitle: detailRouteTitle,
      breadcrumbTitle: detailBreadcrumbTitle,
      tabGroup: 'ops',
      tabTitle: detailRouteTitle,
      title: detailRouteTitle,
      titleKey: 'container.route.detail.title',
    },
  },
];
