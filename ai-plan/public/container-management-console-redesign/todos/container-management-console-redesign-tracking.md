# Container Management Console Redesign Tracking

## Topic

Container Management Console Redesign

## Scope

Optimize the existing container management page from an MVP CRUD-style list into a practical operations-console list
experience, then close backend/OpenAPI fields, controlled operations, validation, and governance.

## Repository Truth

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/design/容器管理设计.md`
- `ai-plan/design/分页列表页统一规范与收敛计划.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/design/系统配置模型与渲染设计.md`
- `ai-plan/design/TDesign-MCP-辅助开发规范.md`
- `web/docs/frontend-i18n-guidelines.md`
- `web/docs/frontend-log-page-guidelines.md`

## Current Recovery Point

- Phase 1 wide-screen list convergence is implemented and validated on the frontend.
- Container list now keeps refresh in the TableCard toolbar, uses local pagination, column settings, density toggle, and
  a stable detail/logs/more action pattern.
- Default visible columns are status, container, image, ports, runtime/status, created time, and operation; `started_at`
  and `restart_policy` are optional column settings.
- Backend/OpenAPI field and server-pagination authority remains intentionally deferred to Phase 3.
- Next batch: `phase-2-detail-logs-drawers`.

## Task Checklist

- [x] Phase 0: planning and topic persistence.
- [x] Phase 1: wide-screen list convergence.
  - [x] PageHeader overview chips.
  - [x] Move refresh exclusively to TableCard toolbar.
  - [x] Tight FilterBar query layout.
  - [x] Default column redesign.
  - [x] Move `started_at` and `restart_policy` to optional columns.
  - [x] Add column settings.
  - [x] Add density toggle.
  - [x] Add `ManagementTablePagination`.
  - [x] Convert row actions to detail/logs/more.
  - [x] Preserve fill/internal-scroll policy and prevent page-level horizontal scroll.
- [ ] Phase 2: detail and logs Drawer improvements.
  - [ ] Detail Drawer width and sections.
  - [ ] Raw JSON collapsed area.
  - [ ] Copy ID.
  - [ ] Logs Drawer auto-refresh, copy, timestamps, and error states.
- [ ] Phase 3: backend/OpenAPI field and pagination enhancement.
  - [ ] List pagination and filters.
  - [ ] List summary.
  - [ ] Health field.
  - [ ] Stats/resource fields with graceful unavailability.
  - [ ] Ports/network/Compose summaries.
  - [ ] Nullable rules and generated artifacts.
- [ ] Phase 4: controlled operations closure.
  - [ ] Operation availability flags.
  - [ ] Optional remove endpoint and permission if implemented.
  - [ ] Dangerous action gate.
  - [ ] Audit for all write operations.
  - [ ] Confirm dialogs and batch operations.
  - [ ] Error codes and i18n keys.
- [ ] Phase 5: experience polish and governance closeout.
  - [ ] Auto-refresh polish.
  - [ ] Stats refresh behavior.
  - [ ] Column preference persistence.
  - [ ] Runtime extension copy and field names.
  - [ ] Final validation.
  - [ ] Browser evidence under `.ai/artifacts/browser/container-page-width-check` or equivalent.
  - [ ] Update `ai-plan/design/容器管理设计.md`.
  - [ ] Archive this topic.
  - [ ] Delete `ai-plan/dolist/container-management-console-redesign-plan.md`.

## Batch State

```json
{
  "loop_mode": "topic-completion-loop",
  "completed_batches": [
    "phase-0-planning-topic-persistence",
    "phase-1-wide-screen-list-convergence"
  ],
  "pending_batches": [
    "phase-2-detail-logs-drawers",
    "phase-3-backend-openapi-fields-pagination",
    "phase-4-controlled-operations-closure",
    "phase-5-polish-validation-governance-closeout"
  ],
  "current_batch": null,
  "next_batch": "phase-2-detail-logs-drawers",
  "closeout_status": "active"
}
```
