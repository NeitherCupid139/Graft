import type { ApiRequestError } from '@/types/axios';

export type RoleFormField = 'name' | 'display' | 'description';
export type RolePermissionField = 'permission_ids';

const roleFormFieldMap: Record<string, RoleFormField> = {
  name: 'name',
  display: 'display',
  description: 'description',
};

const rolePermissionFieldMap: Record<string, RolePermissionField> = {
  permission_ids: 'permission_ids',
};

export function resolveRoleFormFieldError(error: ApiRequestError): RoleFormField | null {
  const field = readField(error.responseData);
  if (!field) {
    return null;
  }

  return roleFormFieldMap[field] ?? null;
}

export function resolveRolePermissionFieldError(error: ApiRequestError): RolePermissionField | null {
  const field = readField(error.responseData);
  if (!field) {
    return null;
  }

  return rolePermissionFieldMap[field] ?? null;
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
