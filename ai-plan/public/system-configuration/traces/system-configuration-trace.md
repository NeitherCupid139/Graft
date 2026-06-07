# System Configuration Trace

## 2026-06-07 topic start

- Re-ran startup preflight for a `cross-boundary` task:
  - read root `AGENTS.md`
  - read `.ai/environment/tools.ai.yaml`
  - read `server/AGENTS.md`
  - read `web/AGENTS.md`
  - read `$graft-multi-agent-loop` and `$graft-multi-agent-task`
- Confirmed topic authority:
  - `server/internal/configregistry` should own the platform declaration registry
  - `server/modules/system-config` should own persistence, API, menu, permissions, and service behavior
  - OpenAPI source should own wire contracts
  - web should consume backend menu/permission/OpenAPI contracts
- Renamed branch from `feat/scheduled-task-mvp` to `feat/system-configuration`.
- User constraints incorporated:
  - plan sensitive/masked fields in OpenAPI from the start
  - do not place the registry in `server/internal/config`
  - store administrator overrides only in `system_config_values`
  - keep `ConfigDefinition` as module-registered authority, not persisted truth

## 2026-06-07 Batch 1 retry worker

- Implemented backend authority baseline:
  - added `server/internal/configregistry` for module-registered `ConfigDefinition` authority, validation, value typing, and masking metadata
  - exposed the registry through `module.Context` and the runtime service container
  - added `server/modules/system-config` with typed route/permission/message contracts, menu placement at `/server/system-config`, route handlers, service/store boundary, and override-only SQL migration
  - wired `system-config` into compile-time module registration with module-owned migration path
- Implemented OpenAPI source baseline:
  - added list/detail/update/reset paths
  - added schemas for system config items, list responses, update request, and envelopes
  - modeled `sensitive`, `masked`, nullable value fields, and `masked_placeholder` so sensitive effective/default/current values are not returned as plaintext
- Generation note:
  - added the backend systemconfig generated-type package needed by the new routes
  - current OpenAPI freshness automation does not yet include a `backend-system-config` target; extending that checker is left for a later validation/automation slice

## 2026-06-07 Batch 1 owner acceptance

- Corrected the API/menu boundary before commit:
  - API source and generated bundle use `/api/system-configs`
  - the Service Management menu path remains `/server/system-config`
  - menu order is `105`, after Scheduled Tasks (`104`)
- Re-ran focused validation:
  - `cd server && go test ./internal/configregistry ./modules/system-config ./internal/moduleregistry ./internal/app`
  - `cd server && go build ./cmd/graft`
  - `cd server && go test ./internal/contract/openapi/...`
  - `git diff --check`
