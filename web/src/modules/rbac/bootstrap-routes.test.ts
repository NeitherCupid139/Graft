import { describe, expect, it } from 'vitest';

import { rbacBootstrapRouteRegistrations } from './bootstrap-routes';
import { RBAC_BOOTSTRAP_ROUTE } from './contract/bootstrap';

describe('rbac bootstrap route registrations', () => {
  it('uses the canonical RBAC bootstrap identity contract values', () => {
    expect(rbacBootstrapRouteRegistrations).toHaveLength(2);
    expect(rbacBootstrapRouteRegistrations[0]).toMatchObject(RBAC_BOOTSTRAP_ROUTE.ROLE_LIST);
    expect(rbacBootstrapRouteRegistrations[1]).toMatchObject(RBAC_BOOTSTRAP_ROUTE.PERMISSION_LIST);
  });
});
