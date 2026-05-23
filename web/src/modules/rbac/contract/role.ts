import type { components } from '@/contracts/openapi/generated/schema';

export type RoleListItem = components['schemas']['RoleListItem'];
export type RoleListResponse = components['schemas']['RoleListResponse'];

// UserRoleBindingResponse matches the role binding contract used across user and rbac modules.
export interface UserRoleBindingResponse {
  role_ids: number[];
}
