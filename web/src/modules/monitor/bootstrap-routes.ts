import type { BootstrapRouteRegistration } from '@/modules/types';

import { MONITOR_ROUTE_PATH } from './contract/paths';

export const monitorBootstrapRouteRegistrations: BootstrapRouteRegistration[] = [
  {
    menuPath: MONITOR_ROUTE_PATH.SERVER_STATUS_OVERVIEW,
    routeName: 'MonitorServerStatusOverview',
    loadPage: () => import('./pages/index.vue'),
  },
];
