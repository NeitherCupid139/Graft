// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import type { BootstrapRouteRegistration } from '@/modules/types';
import { localizeRouteTitleKey } from '@/utils/route/title';

import { ANNOUNCEMENT_BOOTSTRAP_ROUTE } from './contract/bootstrap';

const managementRouteTitle = localizeRouteTitleKey('announcement.route.management.title');
const managementBreadcrumbTitle = localizeRouteTitleKey('announcement.route.management.breadcrumb');

export const announcementBootstrapRouteRegistrations: BootstrapRouteRegistration[] = [
  {
    ...ANNOUNCEMENT_BOOTSTRAP_ROUTE.MANAGEMENT,
    loadPage: () => import('./pages/management/index.vue'),
    meta: {
      pageKind: 'list',
      semanticTitle: managementRouteTitle,
      breadcrumbTitle: managementBreadcrumbTitle,
      tabGroup: 'server',
      tabTitle: managementRouteTitle,
    },
  },
];
