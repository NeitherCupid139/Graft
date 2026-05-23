import { API_CODE } from '@/contracts/api/codes';
import type { ApiRequestError } from '@/types/axios';

export type UserFormField = 'username' | 'display' | 'password';
export type ResetPasswordField = 'password';

const createUserFieldMap: Record<string, UserFormField> = {
  username: 'username',
  display: 'display',
  password: 'password',
  new_password: 'password',
};

const updateUserFieldMap: Record<string, UserFormField> = {
  username: 'username',
  display: 'display',
};

const resetPasswordFieldMap: Record<string, ResetPasswordField> = {
  password: 'password',
  new_password: 'password',
};

export function resolveUserFormFieldError(error: ApiRequestError, mode: 'create' | 'edit'): UserFormField | null {
  const field = readField(error.responseData);
  if (!field) {
    return mode === 'create' && error.code === API_CODE.AUTH_PASSWORD_POLICY_VIOLATION ? 'password' : null;
  }

  return (mode === 'create' ? createUserFieldMap : updateUserFieldMap)[field] ?? null;
}

export function resolveResetPasswordFieldError(error: ApiRequestError): ResetPasswordField | null {
  const field = readField(error.responseData);
  if (!field) {
    return error.code === API_CODE.AUTH_PASSWORD_POLICY_VIOLATION ||
      error.code === API_CODE.AUTH_PASSWORD_REUSE_FORBIDDEN
      ? 'password'
      : null;
  }

  return resetPasswordFieldMap[field] ?? null;
}

function readField(payload: unknown): string | null {
  if (!payload || typeof payload !== 'object' || !('data' in payload)) {
    return null;
  }

  const data = (payload as { data?: unknown }).data;
  if (!data || typeof data !== 'object' || !('field' in data)) {
    return null;
  }

  const field = (data as { field?: unknown }).field;
  return typeof field === 'string' && field.trim() !== '' ? field : null;
}
