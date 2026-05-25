# OAPI Generated Server/Client Governance Spike

## Summary

This topic owns the `monitor/server-status` pilot for generated server/client governance constraints.

- Topic: `oapi-generated-server-client-governance-spike`
- Task class: `cross-boundary`
- Branch: `feat/oapi-generated-server-client-governance-spike`
- Current recommendation: `implement_monitor_server_and_client_spike`

## Scope

- Writable scope:
  - `ai-plan/public/oapi-generated-server-client-governance-spike/**`
  - `server/internal/contract/openapi/monitor/**`
  - `server/plugins/monitor/**`
  - `web/src/modules/monitor/**`
- Read-only context:
  - `openapi/**`
  - `server/internal/httpx/**`
  - `web/src/utils/request.ts`

## Pilot Rules

- Keep `monitor/server-status` as the only operation in the pilot.
- Keep `httpx` as the backend envelope and localized error owner.
- Keep `request.ts` as the frontend transport/runtime owner.
- Do not broaden the pilot to `auth`, `user`, or `rbac`.
- Do not introduce a second global router or transport truth.

## Current Implementation Shape

- Backend:
  - generated monitor-only server bindings live in `server/internal/contract/openapi/monitor/generated/**`
  - monitor plugin still owns explicit route registration
  - generated layer constrains parameter binding and handler interface only
- Frontend:
  - monitor API uses operation-bound generated typing
  - pages still consume module API only
  - `request.ts` remains the only runtime transport adapter

## Validation Expectation

- `git diff --check`
- `cd web && bun run openapi:types:check`
- `cd web && bun run check`
- `cd server && go run ./cmd/graft validate backend`
- `scripts/run_python.sh scripts/magic_value/check_magic_values.py --mode ci`
