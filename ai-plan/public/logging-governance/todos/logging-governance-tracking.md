# Logging Governance Tracking

## Topic

- Topic: `logging-governance`
- Status: `archived`
- Goal:
  - complete a bounded logging inventory and governance design for MVP-stage `Graft`
  - keep the round read-only with runtime implementation deferred
- Recovery source:
  - `ai-plan/public/README.md`
  - `ai-plan/public/logging-governance/README.md`
  - `ai-plan/public/logging-governance/todos/logging-governance-tracking.md`
  - `ai-plan/public/logging-governance/traces/logging-governance-trace.md`
- Worktree: `feat/wt-audit-plugin-mvp`
- Branch: `feat/logging-governance`
- Task class: `cross-boundary`
- Loop mode: `topic-completion-loop`

## Startup Receipt

- Governance source: `root AGENTS.md`
- Task class: `cross-boundary`
- Recovery source: `parent topic`
  - `ai-plan/public/README.md`
  - `ai-plan/public/logging-governance/README.md`
  - `ai-plan/public/logging-governance/todos/logging-governance-tracking.md`
  - `ai-plan/public/logging-governance/traces/logging-governance-trace.md`
- Authority summary:
  - current authority for the evaluation target was the existing runtime logging behavior in `server/**`
  - `web/**` remained a secondary read-only input for frontend console/error reporting and trace-context handling

## Scope

- Owned doc scope:
  - `ai-plan/public/logging-governance/**`
  - `ai-plan/public/README.md`
  - `temp/logging-governance-assessment.md`
- Read-only exploration scope:
  - `server/**`
  - `web/**` only if needed
- Out-of-scope this round:
  - runtime code changes
  - database schema changes
  - audit UI extension
  - observability platform rollout

## Batch State

- Current batch: `Batch 3 - Architecture synthesis and closeout docs`
- Completed batches:
  - `Batch 0 - Topic setup and startup receipt`
  - `Batch 1 - Server logging inventory`
  - `Batch 2 - Web logging/error inventory`
  - `Batch 3 - Architecture synthesis and closeout docs`
- Pending batches:
  - none

## Batch Findings Snapshot

- Batch 1:
  - `zap` is already the primary backend logger via `server/internal/logger/logger.go` and runtime/plugin injection.
  - remaining drift includes stdlib `log` in CLI entrypoints, Ent debug logging in `user` and `rbac`, Gin default
    access logging, and request-id coverage that does not fully cover root routes.
- Batch 2:
  - frontend already uses a structured logger wrapper around browser console via `consola`.
  - no direct `console.*` drift was reported in `web/src`, but there is still no global frontend error sink and
    frontend log metadata does not automatically carry backend `traceId`.
- Batch 3:
  - recommended architecture keeps `AppLogger`, `AccessLogger`, `AuditRecorder`, security events, and
    `MetricsEmitter` as separate explicit responsibilities.
  - highest-priority implementation follow-up is backend request correlation plus structured access logging, not new
    telemetry products.

## Current Recovery Point

- The design-only loop is complete and archived.
- Final assessment lives at `temp/logging-governance-assessment.md`.
- Any runtime implementation should start as a new bounded topic and preserve `zap` as the backend baseline unless new
  evidence justifies a replacement.

## Validation Run

- Executed:
  - `git diff --check`
- Not executed by design:
  - `cd server && go run ./cmd/graft validate backend`
  - `cd web && bun run check`

## Archive Risks To Watch

- request correlation remains incomplete until access logging and request-id middleware become global
- audit evidence quality still depends on partly manual trace/request context propagation
- frontend runtime errors can still escape centralized capture until a future global sink is introduced
- metrics remain a distinct future capability and should not be backfilled by parsing logs
