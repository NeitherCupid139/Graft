# Plugin Audit Correlation Governance Trace

## 2026-05-29 Batch 1 inventory and bounded fix

- Re-ran startup preflight from root `AGENTS.md`.
- Read:
  - root `AGENTS.md`
  - `server/AGENTS.md`
  - `.ai/environment/tools.ai.yaml`
  - `ai-plan/public/README.md`
  - archive-ready `logging-unification-rollout`
  - archive-ready evidence `request-correlation-access-logging`
- Confirmed task class remains `server`.
- Confirmed bounded authority:
  - `server/internal/httpx/**` owns canonical request correlation extraction
  - `server/internal/pluginapi/**` owns request-auth actor transfer
  - `server/plugins/audit/**` owns plugin-domain audit normalization before persistence

## Inventory

| Area | File | Call Site | Uses ctx | requestId | traceId | actorId | route | status | Action |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Request correlation authority | `server/internal/httpx/response.go` | `EnsureRequestID`, `RequestIDMiddleware` | partial | canonical | alias to requestId in response model | none | none | compliant but request-id-only | keep as request-id source; add request-audit context carrier |
| Auth/request middleware | `server/internal/httpx/authz.go` | `RequirePermission` | yes | canonical | canonical alias | from `RequestAuthContext` on security events | canonical route | missing plugin-domain transfer | inject canonical request-audit snapshot into `request.Context` |
| Security event bridge | `server/internal/httpx/authz.go` | `eventBusSecurityAuditPublisher.Publish` | yes | canonical | canonical alias | canonical | canonical route | compliant | preserve existing path |
| Unified audit metadata | `server/internal/audit/service.go` | `candidateMetadata` | no | canonical + legacy alias | canonical + legacy alias | canonical + legacy alias | canonical + legacy alias | compliant | preserve |
| Audit event contract | `server/internal/pluginapi/audit.go` | `AuditEvent` DTO | n/a | optional payload field | legacy implicit alias only | optional operator | optional request path | legacy/minimal | preserve DTO; do not add second model |
| Audit plugin request path | `server/plugins/audit/plugin.go` | `requestAuditCandidate` | yes | canonical | canonical alias | canonical from request auth | canonical route | compliant | preserve |
| Audit plugin domain-event path | `server/plugins/audit/plugin.go` | `recordEvent -> eventAuditCandidate` | before fix: no effective context enrichment | payload/manual only | payload/manual only | payload/manual only | payload/manual only | real risk | enrich from `context.Context` in unified path |
| User plugin domain audit publisher | `server/plugins/user/plugin.go` | `publishAudit` and event constructors | yes | legacy helper `currentRequestID(ctx)` | none | operator from request auth | none | partial/manual | leave helper unchanged; rely on unified enrichment |
| User route write adapters | `server/plugins/user/route_user_handlers.go` | `withAuditRequestID(...)` | yes | header copy only | none | inherited through request auth | none | partial/manual | preserve compatibility helper |
| RBAC plugin domain audit publisher | `server/plugins/rbac/write_service.go` | `publishAudit` and event constructors | yes | legacy helper `currentRBACRequestID(ctx)` | none | operator from request auth | none | partial/manual | leave helper unchanged; rely on unified enrichment |
| RBAC route write adapters | `server/plugins/rbac/route_write_handlers.go` | `withRBACAuditRequestID(...)` | yes | header copy only | none | inherited through request auth | none | partial/manual | preserve compatibility helper |
| Legacy alias point | `server/internal/audit/service.go` | metadata emit: `request_id`, `trace_id`, `request_method`, `request_path`, `actor_id` | n/a | legacy alias | legacy alias | legacy alias | legacy alias | compliant | preserve unchanged |

## Fix Summary

- Added canonical `httpx.RequestAuditContext` transfer on the `RequirePermission` middleware path.
- Updated audit-plugin domain-event normalization to merge:
  - explicit payload fields first
  - canonical request-audit context second
  - request-auth actor fallback when `payload.Operator` is absent
- Preserved current legacy request-id helper behavior in user/RBAC plugins.
- Preserved unified audit legacy aliases.

## Validation

- Passed:
  - `cd server && go test ./internal/httpx ./internal/audit ./plugins/user/... ./plugins/rbac/... ./plugins/audit/...`
