# Plugin Audit Correlation Governance

## Topic Status

- Status: `archive-ready`
- Task class: `server`
- Worktree: `feat/wt-audit-plugin-mvp`
- Batch: `batch-1-inventory-fix-and-closeout`

## Goal

Close the documented non-goal left by `logging-unification-rollout`: plugin-owned domain audit events must inherit the
canonical request correlation and actor semantics that already exist in `server/internal/httpx/**` and the unified
audit path.

## Authority Summary

- `server/internal/httpx/**` owns canonical HTTP request correlation extraction.
- `server/internal/pluginapi/**` owns stable request-auth actor transfer.
- `server/plugins/audit/**` owns plugin-domain audit event normalization before persistence.

## Bounded Outcome

- plugin-domain audit events now inherit `requestId`, `traceId`, `actorId`, `route`, `method`, `clientIp`, and
  `userAgent` from `context.Context` in the unified audit path when publishers do not provide them explicitly
- explicit event payload fields still win over inferred context values
- legacy aliases such as `request_id` and `trace_id` remain unchanged
- RBAC and user plugin legacy request-id helpers remain in place for compatibility, but they are no longer the only
  path for request correlation propagation
