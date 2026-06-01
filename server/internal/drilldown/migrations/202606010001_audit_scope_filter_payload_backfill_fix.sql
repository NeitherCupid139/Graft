UPDATE "system_drilldown_scope"
SET
  "filter_payload" = NULL,
  "updated_at" = NOW()
WHERE "module" = 'audit'
  AND "scope" IN (
    'failed_operations',
    'high_risk_operations',
    'sensitive_operations',
    'auth_failures',
    'permission_denials',
    'rbac_changes',
    'critical_security'
  );
