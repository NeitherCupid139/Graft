import { AUDIT_ROUTE_PATH } from './paths';

export type AuditLogsRouteQuery = Partial<{
  preset: string;
  keyword: string;
  actor: string;
  action: string;
  source: string;
  resourceType: string;
  resourceName: string;
  resourceId: string;
  result: string;
  riskLevel: string;
  session: string;
  requestId: string;
  traceId: string;
}>;

function trimQueryValue(value: unknown) {
  return typeof value === 'string' ? value.trim() : '';
}

function normalizeAuditLogsRouteQuery(query: AuditLogsRouteQuery) {
  return Object.fromEntries(
    Object.entries(query)
      .map(([key, value]) => [key, trimQueryValue(value)])
      .filter(([, value]) => value !== ''),
  ) as Record<string, string>;
}

export function buildAuditLogsLocation(query: AuditLogsRouteQuery) {
  return {
    path: AUDIT_ROUTE_PATH.LOGS,
    query: normalizeAuditLogsRouteQuery(query),
  };
}

export function buildAuditResourceLocation(resourceType: string, resourceId: string, resourceName?: string) {
  return buildAuditLogsLocation({
    resourceType,
    resourceName,
    resourceId,
  });
}

export function buildAuditRequestLocation(requestId: string) {
  return buildAuditLogsLocation({
    requestId,
  });
}
