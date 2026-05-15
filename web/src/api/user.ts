import type { UserListResponse } from '@/api/model/userModel';
import { request } from '@/utils/request';

const Api = {
  Users: '/api/users',
} as const;

export function getUsers() {
  return request.get<UserListResponse>({
    url: Api.Users,
  });
}
