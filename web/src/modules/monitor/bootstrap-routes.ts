import type { BootstrapRouteRegistration } from '@/modules/types';

import { MONITOR_ROUTE_PATH } from './contract/paths';

export const monitorBootstrapRouteRegistrations: BootstrapRouteRegistration[] = [
  {
    menuPath: MONITOR_ROUTE_PATH.SERVER_STATUS,
    routeName: 'MonitorServerStatus',
    loadPage: () => import('./pages/index.vue'),
  },
];
