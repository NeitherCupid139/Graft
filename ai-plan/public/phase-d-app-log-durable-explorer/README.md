# Phase D App Log Durable Explorer

## Status

- Topic: `phase-d-app-log-durable-explorer`
- Status: `active`
- Task class: `cross-boundary`
- Recovery source: `parent topic`
  - `phase-d-app-log-retention-authz-and-storage-readiness`

## Startup Receipt

- governance source: root `AGENTS.md`
- task class: `cross-boundary`
- recovery source: `parent topic`
- authority summary:
  - `server/internal/logger/**` owns App Log semantics, durable storage, cleanup, and read API runtime once approved
  - `openapi/**` owns shared App Log wire contracts
  - `web/src/modules/app-log/**` is the downstream App Log Explorer consumer
  - `server/internal/httpx/**` remains Access Log authority and must not be reused as App Log storage
  - `server/internal/audit/**` and `server/modules/audit/**` remain Audit Log / Security Event authority

## Goal

Implement a bounded repository-owned App Log durable troubleshooting surface after approving the runtime authority:

- durable App Log persistence owned by `server/internal/logger/**`
- retention cleanup owned by the same logger boundary
- read-only list/detail API with an explicit permission
- OpenAPI contract and generated consumers
- `web/src/modules/app-log/**` Explorer page for operator troubleshooting

## Scope Guardrails

- Do not reuse `access_logs`, `audit_logs`, Redis, `access_log.read`, or `audit.read`.
- Do not add security-event standalone storage or metrics/tracing semantics.
- Do not persist access-owned fields such as path, status code, client IP, user agent, request size, or response size.
- Do not persist audit/security-owned actor, resource, action, decision, policy, permission, session, or credential fields.
- Keep App Log as runtime troubleshooting data, not audit evidence or access analytics.

## Loop State

- loop mode: `topic-completion-loop`
- max rounds: `3`
- max commits: `4`
- validation failure policy: stop on failure unless the next bounded batch is a direct fix
- completed batches:
  - Batch 1 backend approval and runtime foundation
- pending batches:
  1. Add OpenAPI contract and web App Log Explorer consumer.
  2. Final validation, governance alignment, and archive-readiness closeout.

## Batch 1 Backend Foundation

- Status: `completed`
- Owner: `server/internal/logger/**`
- Implemented:
  - PostgreSQL `app_logs` live migration under `server/internal/logger/migrations/**`
  - logger-owned `AppLogRepository` create/list/delete foundation
  - optional best-effort `AppLogger` repository sink that preserves zap output
  - `GRAFT_LOG_APP_LOG_PERSIST` and `GRAFT_LOG_APP_LOG_RETENTION` config
  - `logger.app-log-retention-cleanup` cleanup job
  - default migration-chain registration for `internal/logger/migrations`
- Retention defaults:
  - local/test/dev: 3 days
  - staging: 7 days
  - production: 14 days
- Deferred to Batch 2:
  - read permission/menu/API route contract
  - OpenAPI source and generated consumers
  - `web/src/modules/app-log/**` Explorer
