import type { BootstrapRouteRegistration } from '@/modules/types';

export const rbacBootstrapRouteRegistrations: BootstrapRouteRegistration[] = [
  {
    menuPath: '/roles',
    routeName: 'RoleList',
    loadPage: () => import('@/pages/role/index.vue'),
  },
];
