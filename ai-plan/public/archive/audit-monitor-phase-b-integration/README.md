# Audit Monitor Phase B Integration

## Status

- Topic: `audit-monitor-phase-b-integration`
- Status: `archived`
- Loop mode: `topic-completion-loop`
- Task class: `cross-boundary`
- Scope: design-only, archive-ready
- Parent evidence:
  - `ai-plan/public/archive/observability-development-governance/README.md`
  - `ai-plan/public/archive/metrics-governance/README.md`

## Goal

Define the Phase B operator workflow that connects `audit` and `monitor` without widening into a logging, metrics, or tracing platform.

Required outcome:

- an authority-first `monitor -> audit` evidence navigation model
- an authority-first `audit -> monitor` incident-context model
- a focused security timeline drilldown design
- a contract/backend/frontend gap analysis
- a phased implementation roadmap

## Maturity Assessment

- Phase A: completed
- Phase B: partially implemented
- Phase C: mostly implemented
- Phase D: not started
- Phase E: not started

Current repository state:

- `audit` already owns canonical audit overview, audit logs, grouped risk analytics, trend, and security timeline
- `monitor` already owns canonical server-status, runtime, dependencies, plugin health, and short-retention trend data
- both modules still operate as separate read surfaces
- there is no canonical anomaly model, no evidence-link contract, no monitor-to-audit navigation, and no security-incident drilldown page

## Authority Discovery

### Audit Authority

- Audit overview semantic owner:
  - `server/plugins/audit/store/audit.go`
- Audit overview shared HTTP owner:
  - `openapi/paths/audit.overview.yaml`
  - `openapi/components/schemas/audit-overview-response.yaml`
- Audit logs semantic owner:
  - `server/plugins/audit/store/audit.go`
- Audit logs shared HTTP owner:
  - `openapi/paths/audit.logs.yaml`
  - `openapi/components/schemas/audit-log-list-item.yaml`
  - `openapi/components/schemas/audit-log-list-response.yaml`
- Security timeline owner:
  - backend semantics stay in `server/plugins/audit/store/audit.go`
  - shared wire shape is only the nested `security_timeline` fragment inside `openapi/components/schemas/audit-overview-response.yaml`
- Current drilldown owner:
  - backend only supports audit-log filtering
  - frontend-only route/deep-link logic lives in `web/src/modules/audit/contract/deep-link.ts`

Audit authority gaps:

- no audit-owned incident/drilldown read model
- no standalone security-timeline drilldown contract
- current deep-link contract is only partially backed by OpenAPI
- `traceId`, `session`, `preset`, `keyword`, and `actor` remain frontend-only route state instead of shared drilldown authority

### Monitor Authority

- Monitor route/menu/permission owner:
  - `server/plugins/monitor/contract/route.go`
  - `server/plugins/monitor/contract/permission.go`
  - `server/plugins/monitor/contract/message.go`
- Monitor runtime implementation owner:
  - `server/plugins/monitor/plugin.go`
- Monitor shared HTTP owner:
  - `openapi/paths/monitor.server-status.yaml`
  - `openapi/components/schemas/server-status-*.yaml`
- Monitor trend owner:
  - `server/plugins/monitor/contract/trend.go`
  - `openapi/components/parameters/trend-range-query.yaml`
  - `openapi/components/schemas/server-status-trend.yaml`
  - `openapi/components/schemas/server-status-trend-point.yaml`

Monitor authority gaps:

- no canonical anomaly DTO or anomaly-definition contract
- no evidence-link field on the monitor payload
- no drilldown route contract from monitor surfaces into audit evidence
- parts of anomaly semantics currently live in frontend-only threshold logic in `web/src/modules/monitor/pages/overview/index.vue`

### Cross-Boundary Authority

- Anomaly definitions should be owned by `server/plugins/monitor/**`, with the shared wire model defined in `openapi/**`
- Audit evidence context should be owned by `audit`, because audit query semantics and incident grouping belong to the audit read model
- Shared evidence-link wire shape should live in `openapi/**`
- Monitor may emit audit evidence links, but it must emit them using an audit-owned context shape rather than inventing monitor-local audit query rules
- Security incident drilldown should be owned by `audit`, because the incident seed and related-event expansion come from audit authority first
- Related monitor context should remain a monitor-owned attached panel or linked read model, explicitly bounded by current short-retention trend authority

### Required Authority Repairs

1. Promote audit deep-link semantics from frontend-only route helpers into an audit-owned shared context contract.
2. Move monitor anomaly semantics from frontend threshold code into backend/OpenAPI authority.
3. Introduce a dedicated audit incident/drilldown read model instead of treating `security_timeline -> request_id filter` as the final design.

## Observability Integration Design

### Canonical Link Model

Canonical direction:

`Monitor Resource -> Evidence Link -> Audit Context`

Definitions:

- Monitor resource:
  - a backend-owned monitor subject such as dependency, plugin runtime status, runtime metric family, or resource-usage card
- Evidence link:
  - a backend-produced cross-boundary link object attached to a monitor anomaly
- Audit context:
  - an audit-owned filter or incident seed that the audit module can resolve without frontend compatibility mapping

Proposed shared wire model:

- `EvidenceLink`
  - `target_kind`
    - `audit-context`
    - `audit-incident`
  - `link_state`
    - `available`
    - `empty`
    - `unsupported`
    - `unavailable`
  - `title`
  - `reason`
  - `time_window`
  - `audit_context`
    - audit-owned filter DTO using only backend-supported semantics
  - `incident_seed`
    - optional audit event id when the link points to an incident drilldown

Authority rule:

- `audit_context` is not a frontend query-string invention
- `audit_context` is the canonical reusable cross-module evidence payload
- web route building stays downstream and deterministic

### Monitor Anomaly Model

Phase B should add a bounded monitor anomaly read model, not a platform alerting system.

Canonical owner:

- `server/plugins/monitor/**`
- shared wire surface in `openapi/**`

Bounded anomaly classes:

- `dependency_status_degraded`
- `dependency_status_unknown`
- `plugin_dependency_missing`
- `resource_cpu_pressure`
- `resource_memory_pressure`
- `resource_disk_pressure`
- `runtime_goroutine_pressure`
- `runtime_heap_pressure`
- `system_load_pressure`

Each anomaly should expose:

- `anomaly_key`
- `scope_kind`
  - `dependency`
  - `plugin`
  - `runtime`
  - `resource`
- `scope_ref`
- `severity`
- `status`
- `observed_at`
- `summary`
- `evidence_links[]`

This stays plugin-local to `monitor` and does not reopen platform-wide metrics governance.

### Operator Journey

#### Dependency Failure

1. Operator lands on monitor dependencies or overview.
2. A degraded dependency card exposes one or more backend-owned anomalies.
3. Clicking the anomaly opens an anomaly panel or route-owned detail state.
4. The anomaly panel renders `evidence_links[]`.
5. If an audit evidence link is `available`, the operator jumps into audit logs or an audit incident using the canonical audit context.
6. If link state is `empty` or `unsupported`, the UI says that no audit evidence exists under current authority instead of fabricating a correlation.

Recommended evidence targets:

- permission-denied spike around the same window
- failed-auth spike around the same window
- admin or configuration changes involving the affected dependency resource when such audit records exist

#### Runtime Anomaly

1. Operator sees a runtime or resource anomaly on monitor overview/runtime.
2. The anomaly detail shows recent trend context and evidence links.
3. Audit navigation points to security or administrative activity in the same bounded window, not to a fake “system error log”.
4. If no audit evidence exists, monitor remains authoritative and the UI states that the anomaly has no audit correlation.

#### Resource Anomaly

1. Operator selects a CPU, memory, disk, or goroutine anomaly.
2. Monitor opens the anomaly detail with the backend-owned explanation and threshold status.
3. The operator can pivot to:
  - audit evidence
  - runtime page
  - dependencies page
4. Audit evidence uses the same `EvidenceLink` contract, keeping navigation deterministic across anomaly types.

### Evidence Journey

#### Failed Authentication Spike

1. Operator starts from audit overview, audit logs, or the future incident page.
2. The audit module resolves an incident context from a seed event or grouped evidence filter.
3. The incident page shows:
  - seed event
  - related events
  - related actors
  - related resources
  - related requests
  - related monitor context
4. The related monitor context links to current runtime/dependency state plus recent trend context when retention still covers the window.

#### Permission Denial Spike

1. Operator starts from audit evidence.
2. Audit expands related denied events and affected actors/resources.
3. The incident page offers monitor context links focused on current plugin/dependency/runtime health.
4. Monitor context is marked `current` or `retention-limited`; it is never presented as historical truth beyond current monitor authority.

#### Security Incident

1. Operator clicks a security timeline event.
2. Audit opens a dedicated incident page keyed by the seed event id.
3. The page builds a bounded related-event graph using audit authority first.
4. A monitor panel shows:
  - current dependency and runtime status
  - recent trend context if still retained
  - explicit “historical monitor state unavailable” messaging when the event is outside retention

## Security Timeline Drilldown Design

### Current Limitation

Current flow:

- timeline event
- request-id filter on `/audit/logs`

This is insufficient because it loses incident grouping, actors, resources, and related runtime context.

### Target Workflow

- `Timeline Event`
- `Incident Context`
- `Related Audit Events`
- `Related Actors`
- `Related Resources`
- `Related Requests`
- `Related Monitor State`

### Proposed Contract Shape

New audit-owned endpoint:

- `GET /api/audit/incidents/{event_id}`

Audit incident response should contain:

- `seed_event`
- `incident`
  - `incident_key`
  - `title`
  - `summary`
  - `risk_level`
  - `started_at`
  - `ended_at`
  - `correlation_reason`
- `related_events`
  - bounded list
  - optional pagination token for later expansion
- `related_actors`
- `related_resources`
- `related_requests`
- `monitor_context`
  - `state`
    - `available`
    - `partial`
    - `unavailable`
  - `reason`
  - `current_status`
  - `trend_window`
  - `monitor_links[]`

Optional audit evidence payload reuse:

- `AuditContext`
  - stable filter fields supported by backend/OpenAPI
- `AuditIncidentSeed`
  - `event_id`

### Correlation Rules

Correlation remains bounded and explicit:

- primary seed:
  - event id from `security_timeline`
- first-order joins:
  - `request_id`
  - `actor_user_id`
  - `resource_type`
  - `resource_id`
  - `session_id` when audit authority is extended to support it
- bounded time window:
  - derived by audit backend, not the frontend

### Required Routes

- New audit page route:
  - `/audit/incidents/:eventId`
- New frontend route contract owner:
  - `web/src/modules/audit/**` as consumer of an audit-owned shared contract
- Optional monitor deep-link query additions:
  - focused anomaly key
  - focused trend range
  - source incident id

### Required UI States

- Incident loading
- Incident unavailable
- No related events
- Monitor context available
- Monitor context partial due to retention
- Monitor context unavailable
- Audit evidence link empty or unsupported

## Gap Analysis

### Contract Gaps

- missing monitor anomaly contract
- missing shared evidence-link contract
- missing audit-owned evidence-context contract
- missing audit incident/drilldown endpoint contract
- missing monitor-context contract for the incident page
- missing backend-backed `session_id` and `trace_id` drilldown semantics if those remain required

### Backend Gaps

- monitor plugin lacks anomaly classification output
- monitor plugin lacks evidence-link generation
- audit plugin lacks incident read-model aggregation
- audit repository lacks bounded correlation queries for incident expansion
- audit plugin cannot currently return a dedicated incident payload
- monitor plugin only retains short-window trend points and does not retain historical dependency snapshots

### Frontend Gaps

- no monitor -> audit navigation on overview/runtime/dependencies pages
- no audit -> monitor navigation beyond generic page switching
- no incident drilldown page
- audit deep-link logic is only partially backed by backend authority
- monitor anomaly thresholds are partly frontend-owned today

### UX Gaps

- current security timeline click is too shallow
- current monitor views surface health but not explainable evidence pivots
- there is no operator workspace that preserves anomaly, evidence, and runtime context in one workflow
- current UX cannot honestly distinguish “no audit evidence exists” from “link not implemented”

## Roadmap

### Phase B1: Minimum Valuable Integration

- Authority owner:
  - monitor anomaly semantics: `server/plugins/monitor/**`
  - audit evidence context: `server/plugins/audit/**`
  - shared wire shape: `openapi/**`
- Contract changes:
  - add monitor anomaly fields to `server-status` response
  - add shared `EvidenceLink`
  - add shared audit-owned `AuditContext`
- Backend changes:
  - move current monitor threshold logic into backend-owned anomaly classification
  - generate bounded audit evidence links for supported monitor anomalies
- Frontend changes:
  - render anomaly cards/panels on monitor pages
  - add direct monitor -> audit navigation from evidence links
- Validation:
  - backend: `graft validate backend --stage lint`, targeted `go test` for `plugins/monitor`, `plugins/audit`, related internal packages
  - web: `bun run check`

### Phase B2: Security Timeline Drilldown

- Authority owner:
  - incident read model: `server/plugins/audit/**`
  - shared wire shape: `openapi/**`
- Contract changes:
  - add `GET /api/audit/incidents/{event_id}`
  - add incident DTOs and related monitor-context DTOs
  - extend audit context only where backend semantics truly exist
- Backend changes:
  - implement incident aggregation and bounded correlation queries
  - attach best-effort related monitor context with explicit retention-state semantics
- Frontend changes:
  - add `/audit/incidents/:eventId`
  - change security timeline click target from audit-list filter to incident page
  - add related actors/resources/requests panels and monitor context panel
- Validation:
  - backend targeted tests for incident query logic and HTTP mapping
  - web `bun run check`
  - direct route/component tests for incident page drilldown states

### Phase B3: Observability Workspace Completion

- Authority owner:
  - audit module and monitor module remain separate authorities
  - workspace behavior is integration-only, not a new platform authority
- Contract changes:
  - optional refinement of evidence-link states and monitor focus-link params
  - no new infra contracts
- Backend changes:
  - refine supported evidence-link coverage
  - fill bounded gaps such as session-backed incident joins if justified
- Frontend changes:
  - preserve operator context between monitor anomaly detail, audit incident, and monitor pages
  - add clearer empty/unsupported/unavailable evidence states
- Validation:
  - backend and web full changed-scope validation
  - route/navigation integration coverage

## Recommended First Implementation Slice

Start with Phase B1, not the incident page.

Reason:

- the biggest current authority drift is monitor anomaly semantics living partly in the frontend
- monitor-to-audit navigation cannot be made canonical until anomaly and evidence-link ownership move to backend/OpenAPI
- once `EvidenceLink` and `AuditContext` exist, the later incident drilldown page can reuse them instead of inventing a second correlation model

First slice:

1. Define monitor anomaly DTOs and evidence-link DTOs in `openapi/**`.
2. Move the current monitor overview threshold logic into `server/plugins/monitor/**`.
3. Define an audit-owned `AuditContext` contract limited to backend-supported filter semantics.
4. Render monitor anomaly actions in `web/src/modules/monitor/**` using only backend-provided links.

## Non-Goals Preserved

- no logging platform
- no metrics platform
- no OpenTelemetry
- no Prometheus
- no Grafana
- no tracing platform
- no frontend compatibility layer that hides incorrect backend authority

## Closeout

- Result: `archive-ready design`
- Final judgment:
  - the `audit <-> monitor` integration model is fully defined
  - the security timeline drilldown workflow is fully defined
  - required contracts, backend work, frontend work, UX gaps, and roadmap phases are identified
  - authority repairs are explicit and upstream-first
- Next-session prompt:
  - `Re-run startup preflight from root AGENTS.md. Treat audit-monitor-phase-b-integration as archived design evidence. Start with Phase B1 authority repair in server/plugins/monitor/**, server/plugins/audit/**, openapi/**, and web/src/modules/monitor/** only if the next slice is implementation-focused.`
