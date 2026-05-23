import type { components } from '@/contracts/openapi/generated/schema';

import type { UserStatus } from '../contract/status';

export type UserListItem = components['schemas']['UserListItem'];
export type UserListResponse = components['schemas']['UserListResponse'];

export interface CreateUserPayload {
  username: string;
  display: string;
  password: string;
}

export interface UpdateUserPayload {
  username: string;
  display: string;
}

export interface UpdateUserStatusPayload {
  status: UserStatus;
}

export interface ResetUserPasswordPayload {
  new_password: string;
}
