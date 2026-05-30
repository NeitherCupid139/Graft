import type { BootstrapRouteRegistration } from '@/modules/types';

import { ACCESS_LOG_BOOTSTRAP_ROUTE } from './contract/bootstrap';

export const accessLogBootstrapRouteRegistrations: BootstrapRouteRegistration[] = [
  {
    ...ACCESS_LOG_BOOTSTRAP_ROUTE.LIST,
    loadPage: () => import('./pages/list/index.vue'),
    meta: {
      tabGroup: 'access-log',
      pageKind: 'list',
      semanticTitle: {
        'zh-CN': '访问日志 - 访问日志查询',
        'en-US': 'Access Logs - Explorer',
      },
      breadcrumbTitle: {
        'zh-CN': '访问日志查询',
        'en-US': 'Access Log Explorer',
      },
      tabTitle: {
        'zh-CN': '访问日志 - 访问日志查询',
        'en-US': 'Access Logs - Explorer',
      },
    },
  },
];
