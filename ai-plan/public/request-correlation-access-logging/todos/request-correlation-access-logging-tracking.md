# Request Correlation Access Logging Tracking

## Topic

- Topic: `request-correlation-access-logging`
- Status: `active`
- Goal:
  - implement global request correlation for all HTTP routes
  - introduce structured `zap`-backed access logging
  - keep the implementation bounded to Phase 1 backend logging governance follow-up
- Recovery source:
  - `ai-plan/public/README.md`
  - archived `logging-governance`
  - `temp/logging-governance-assessment.md`
- Branch: `feat/request-correlation-access-logging`
- Task class: `server`
- Loop mode: `topic-completion-loop`

## Startup Receipt

- Governance source: `root AGENTS.md`
- Task class: `server`
- Recovery source: `archived topic evidence`
  - `ai-plan/public/logging-governance/**`
  - `temp/logging-governance-assessment.md`
- Authority summary:
  - `server/internal/httpx/**` owns request middleware and HTTP access logging behavior
  - `server/internal/logger/**` remains the canonical backend logger baseline

## Scope

- Owned scope:
  - `ai-plan/public/request-correlation-access-logging/**`
  - `ai-plan/public/README.md`
  - `server/internal/httpx/**`
  - `server/internal/logger/**`
  - bounded backend tests in directly affected packages
- Forbidden scope:
  - `web/**`
  - `server/plugins/**` unless a direct bounded validation fix becomes necessary
  - OpenTelemetry or metrics platform work
  - unrelated log backend consolidation

## Batch State

- Completed batches:
  - `batch-1-global-correlation-and-access-logger`
  - `batch-2-tests-and-validation`
- Pending batches:
  - `batch-3-closeout-and-archive-check`
- Current batch:
  - `batch-3-closeout-and-archive-check`
- Next batch:
  - `batch-3-closeout-and-archive-check`

## Budget

- Max rounds: `3`
- Max files changed: `12`
- Max commits: `1`
- Max runtime minutes: `90`
- Validation failure policy: `stop-on-failure`
- Checkpoint budget per round: `1`

## Current Recovery Point

- Topic opened from archived `logging-governance` design evidence on `2026-05-29`.
- Batch 1 accepted:
  - `httpx.NewServer` now mounts global request-correlation middleware plus zap-backed structured access logging.
  - root routes and plugin routes share the same request-id entry path instead of relying on plugin-local middleware
    only.
  - direct runtime wiring now passes the runtime logger into `httpx.NewServer`.
- Validation already recorded for Batch 1:
  - `cd server && go test ./internal/httpx ./internal/app`
- Batch 2 accepted:
  - added the smallest extra assertion coverage for access-log severity routing
  - confirmed no further owned-scope runtime changes are justified before closeout
- Validation recorded for Batch 2:
  - `cd server && go test ./internal/httpx ./internal/app`
  - `cd server && go test -cover ./internal/httpx`
- Next target is Batch 3:
  - run archive-readiness and commit-eligibility evaluation for the bounded server slice.
