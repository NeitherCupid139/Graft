import { API_CODE, type ApiCode, type ApiResponseCode } from '@/contracts/api/codes';
import type { LocalizedTitle } from '@/contracts/i18n/locales';

export { API_CODE };
export type { ApiCode, ApiResponseCode };
export type {
  BootstrapLocale,
  BootstrapMenu,
  BootstrapResponse,
  CompleteRequiredPasswordChangePayload,
  LoginPayload,
  LoginResponse,
  LoginUser,
} from '@/modules/auth/types/auth';

export interface ApiSuccessEnvelope<T> {
  success: true;
  code: ApiCode;
  message: string;
  traceId: string;
  data: T;
  messageKey?: string;
  locale?: string;
}

export interface ApiErrorEnvelope {
  success: false;
  code: ApiResponseCode;
  message: string;
  traceId: string;
  data?: Record<string, unknown> | null;
  messageKey?: string;
  locale?: string;
}

export type ApiEnvelope<T> = ApiSuccessEnvelope<T> | ApiErrorEnvelope;

export interface AppBootstrapRouteMeta {
  title: LocalizedTitle;
  titleKey?: string;
  icon?: string;
  permission?: string;
}
