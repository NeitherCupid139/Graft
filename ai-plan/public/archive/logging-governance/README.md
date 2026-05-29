# Logging Governance

## Status

- Topic: `logging-governance`
- Status: `archived`
- Loop mode: `topic-completion-loop`
- Worktree: `feat/wt-audit-plugin-mvp`
- Branch: `feat/logging-governance`
- Task class: `cross-boundary`
- Started: `2026-05-28`
- Closed: `2026-05-28`

## Goal

- Inventory how `Graft` currently uses logs, audit recording, request identifiers, and metrics-like events.
- Distinguish `Application / Runtime Log`, `Error Log`, `Access Log`, `Audit Log`, `Security Event`, and
  `Metrics Event` responsibilities.
- Identify concrete mixing or governance gaps without starting a runtime refactor in this topic.
- Produce an MVP-suitable logging governance design that keeps interfaces explicit and avoids a generic
  `Log(type, ...)` abstraction.

## Scope

- Recovery-doc scope:
  - `ai-plan/public/archive/logging-governance/**`
  - `ai-plan/public/README.md`
- Read-only exploration scope:
  - `server/**`
  - `web/**` only where needed for frontend logging/error reporting understanding
  - `temp/**` for the assessment output artifact requested in this round
- Forbidden scope for this round:
  - runtime code refactors
  - audit-plugin-mvp continuation or closeout work
  - OpenTelemetry rollout
  - request-log database productization
  - new frontend logging platform work

## Planned Batches

- Batch 0: topic setup and startup receipt
- Batch 1: `server` inventory for logger initialization, request/access logging, audit, security, and metrics signals
- Batch 2: `web` inventory for console/error reporting usage and integration risks
- Batch 3: architecture synthesis, migration staging, and assessment closeout docs

## Batch Closeout

- Completed:
  - Batch 0: startup receipt and bounded topic setup completed.
  - Batch 1: server inventory confirmed `zap` is already the main backend logger, with remaining gaps in stdlib log
    usage, Gin access logging, global request-id coverage, and metrics separation.
  - Batch 2: web inventory confirmed a structured local logger exists, but there is no global frontend error sink and
    request `traceId` is not yet promoted into frontend logger context.
  - Batch 3: produced the design assessment, recommended target architecture, phased migration plan, and acceptance
    criteria without widening into runtime code.
- Final assessment artifact:
  - `temp/logging-governance-assessment.md`

## Final Conclusions

- Backend recommendation preserves `zap` as the canonical application logging backend.
- Logging responsibilities should stay split across explicit surfaces instead of collapsing into a generic
  `Log(type, ...)` abstraction.
- Highest-priority follow-up is a future bounded backend implementation topic for global request correlation and
  structured access logging, then security/audit wiring cleanup, then frontend error/context enrichment, then metrics.
- This topic is closed as design evidence only. Any implementation follow-up must start as a new bounded topic.
