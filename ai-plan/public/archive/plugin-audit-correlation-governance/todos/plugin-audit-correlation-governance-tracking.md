# Plugin Audit Correlation Governance Tracking

## Status

- Topic status: `archived`
- Current batch: `batch-1-inventory-fix-and-closeout`
- Continue required: `false`

## Completed

- inventoried audit record entry points, plugin-owned domain audit publishers, correlation sources, compliant call
  sites, missing-field call sites, and legacy alias points
- implemented canonical request-audit context propagation in `server/internal/httpx`
- updated `server/plugins/audit` to enrich plugin-domain audit candidates from `context.Context`
- added bounded tests for requestId, traceId, actorId, nil-HTTP-context safety, and legacy alias preservation
- updated public recovery status for the new topic

## Validation

- passed:
  - `cd server && go test ./internal/httpx ./internal/audit ./plugins/user/... ./plugins/rbac/... ./plugins/audit/...`
