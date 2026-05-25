import type { paths } from '@/contracts/openapi/generated/schema';
import { AUTH_API_PATH } from '@/modules/auth/contract/paths';
import type {
  BootstrapResponse,
  CompleteRequiredPasswordChangePayload,
  LoginPayload,
  LoginResponse,
} from '@/modules/auth/contract/types';
import { request } from '@/utils/request';

type LoginPath = (typeof AUTH_API_PATH)['LOGIN'];
type BootstrapPath = (typeof AUTH_API_PATH)['BOOTSTRAP'];
type RefreshPath = (typeof AUTH_API_PATH)['REFRESH'];
type LogoutPath = (typeof AUTH_API_PATH)['LOGOUT'];
type PostAuthLoginOperation = paths[LoginPath]['post'];
type GetAuthBootstrapOperation = paths[BootstrapPath]['get'];
type PostAuthRefreshOperation = paths[RefreshPath]['post'];
type PostAuthLogoutOperation = paths[LogoutPath]['post'];
type PostAuthLoginResponse = PostAuthLoginOperation['responses']['200']['content']['application/json'];
type GetAuthBootstrapResponse = GetAuthBootstrapOperation['responses']['200']['content']['application/json'];
type PostAuthRefreshResponse = PostAuthRefreshOperation['responses']['200']['content']['application/json'];
type PostAuthLogoutResponse = PostAuthLogoutOperation['responses']['200']['content']['application/json'];
type PostAuthLoginResponseData = NonNullable<PostAuthLoginResponse['data']>;
type GetAuthBootstrapResponseData = NonNullable<GetAuthBootstrapResponse['data']>;
type PostAuthRefreshResponseData = NonNullable<PostAuthRefreshResponse['data']>;
type PostAuthLogoutResponseData = PostAuthLogoutResponse['data'];

// Keep generated request/response typing at the module API boundary; callers still own form-local state.
export function login(payload: LoginPayload) {
  return request.post<PostAuthLoginResponseData>({
    url: AUTH_API_PATH.LOGIN,
    data: payload,
  });
}

export function refresh() {
  return request.post<LoginResponse & PostAuthRefreshResponseData>({
    url: AUTH_API_PATH.REFRESH,
  });
}

export async function logout(): Promise<void> {
  await request.post<PostAuthLogoutResponseData>({
    url: AUTH_API_PATH.LOGOUT,
  });
}

export function completeRequiredPasswordChange(payload: CompleteRequiredPasswordChangePayload) {
  return request.post<void>({
    url: AUTH_API_PATH.COMPLETE_REQUIRED_PASSWORD_CHANGE,
    data: payload,
  });
}

export function getBootstrap() {
  return request.get<BootstrapResponse & GetAuthBootstrapResponseData>({
    url: AUTH_API_PATH.BOOTSTRAP,
  });
}
