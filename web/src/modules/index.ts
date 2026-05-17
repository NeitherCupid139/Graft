import { rbacBootstrapRouteRegistrations } from './rbac/bootstrap-routes';
import type { BootstrapRouteRegistration } from './types';
import { userBootstrapRouteRegistrations } from './user/bootstrap-routes';

const allBootstrapRouteRegistrations: BootstrapRouteRegistration[] = [
  ...rbacBootstrapRouteRegistrations,
  ...userBootstrapRouteRegistrations,
];

const bootstrapRouteRegistrationMap = new Map<string, BootstrapRouteRegistration>();

for (const registration of allBootstrapRouteRegistrations) {
  if (bootstrapRouteRegistrationMap.has(registration.menuPath)) {
    throw new Error(`duplicate bootstrap route registration: ${registration.menuPath}`);
  }

  bootstrapRouteRegistrationMap.set(registration.menuPath, registration);
}

export function getBootstrapRouteRegistration(menuPath: string) {
  return bootstrapRouteRegistrationMap.get(menuPath);
}
