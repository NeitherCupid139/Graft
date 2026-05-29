# Request Correlation Access Logging

## Status

- Topic: `request-correlation-access-logging`
- Status: `archived`
- Branch: `feat/request-correlation-access-logging`
- Task class: `server`
- Loop mode: `topic-completion-loop`
- Started: `2026-05-29`

## Goal

Implement Phase 1 of the archived `logging-governance` plan:

- make request correlation global for all HTTP routes, including root routes such as `/healthz`
- replace Gin default access logging with a structured `zap`-backed access logging path
- keep `zap` as the canonical backend application logging baseline
- keep the slice bounded to backend runtime changes required for request correlation and access logging

## Recovery Inputs

- `ai-plan/public/README.md`
- `ai-plan/public/archive/logging-governance/README.md`
- `ai-plan/public/archive/logging-governance/todos/logging-governance-tracking.md`
- `ai-plan/public/archive/logging-governance/traces/logging-governance-trace.md`
- `temp/logging-governance-assessment.md`

## Scope

- Owned scope:
  - `ai-plan/public/archive/request-correlation-access-logging/**`
  - `ai-plan/public/README.md`
  - `server/internal/httpx/**`
  - `server/internal/logger/**`
  - `server/internal/audit/**` only if required for bounded correlation-field alignment
  - `server/cmd/graft/**` only if directly required by validation or runtime wiring for this slice
  - bounded backend tests in directly affected packages
- Forbidden scope:
  - `web/**`
  - plugin-wide audit redesign beyond bounded field alignment
  - OpenTelemetry, Prometheus, or remote telemetry platform work
  - generic `Log(type, ...)` abstraction
  - unrelated CLI log backend cleanup unless required for this slice

## Acceptance Targets

- all HTTP routes receive one stable request correlation field path
- access logs become structured runtime output and no longer rely on Gin default logger as the long-term solution
- request correlation fields are available to access logs and response envelope/header paths consistently
- backend validation for owned scope runs through the repository backend entrypoint or a justified narrower slice

## Planned Batches

1. Batch 1: implement global request correlation middleware wiring and structured access logger middleware in
   `server/internal/httpx`.
2. Batch 2: align bounded tests and validation evidence, then verify Gin default logger no longer owns access-log output.
3. Batch 3: archive-readiness check, recovery-doc updates, scoped commit evaluation, and handoff prompt if needed.

## Closeout

- Status: `archived`
- Archive reason:
  - the bounded Phase 1 backend slice completed its planned batches
  - validation evidence for the owned `server` scope is recorded
  - later logging work moved to the broader archived `logging-unification-rollout` topic
- Final validation:
  - `cd server && go test ./internal/httpx ./internal/app`
  - `cd server && go test -cover ./internal/httpx`
