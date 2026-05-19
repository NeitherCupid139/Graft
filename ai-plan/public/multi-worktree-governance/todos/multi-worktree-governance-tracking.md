# Multi Worktree Governance Tracking

## Topic

- Topic: `multi-worktree-governance`
- Branch: `refactor/server-module-boundaries`
- Worktree: repository root only; no dedicated long-lived worktree exists yet
- Scope: keep the shared recovery baseline short, freeze cross-worktree ownership truth, and prepare the first
  dedicated long-lived worktree/topic pair without reopening archived MVP recovery state.

## Goal

- Keep the active recovery entry focused on the current shared baseline, not on completed milestone history.
- Preserve the final `web` ownership model and the current `server` compile-time modular-monolith ownership model until
  dedicated worktrees are created.
- Keep historical detail available in topic-local archive snapshots instead of carrying it in the default recovery path.

## Repository Truth

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/design/AI任务追踪与恢复设计.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/插件与依赖注入设计.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/roadmap/MVP实施计划.md`
- `ai-plan/public/README.md`
- `ai-plan/public/multi-worktree-governance/roadmap/server-module-boundaries-plan.md`

## Current Recovery Point

- `mvp-extension-path` is complete and remains archived under `ai-plan/public/archive/mvp-extension-path/`; it is not
  part of the active recovery path.
- The repository root is still the only active worktree and is currently on branch
  `refactor/server-module-boundaries`.
- The frozen `web` ownership model is:
  - `shell-owned`: `web/src/app/**`, `web/src/layouts/**`, `web/src/router/**`, `web/src/config/**`,
    `web/src/utils/route/**`, `web/src/store/modules/{user,permission}.ts`, `web/src/locales/**`, and platform
    `web/src/contracts/**`
  - `module-owned`: `web/src/modules/<name>/**`
  - `shared-owned`: `web/src/shared/**`
- Root-level module-specific files under `web/src/api/**`, `web/src/api/model/**`, and
  `web/src/contracts/{user,rbac}/**` are not valid steady-state ownership surfaces and must not return.
- The frozen `server` ownership model is:
  - compile-time modular monolith only; no runtime plugin loading, discovery, hot-load lifecycle, or generalized
    service locator
  - plugin-first owned scope under `server/plugins/<name>/**`
  - shared stable backend boundary under `server/internal/pluginapi/**` and `server/internal/contract/**`
  - centralized generated hotspot limited to `server/internal/pluginregistry/generated.go`
  - `user_roles` final owner is `rbac`
  - new business logic must not flow back into `server/internal/store/**` or non-core-owned portions of
    `server/internal/ent/**`
- The latest backend ownership slices already landed:
  - `654c791` moved audit persistence into `server/plugins/audit/store/**` and `server/plugins/audit/storeent/**`
  - `5f45b31` removed the shared audit compatibility shim from `internal/store`
  - `799f1ff` removed the shared `user` store compatibility seam
  - `866582a` removed the shared `user/auth` seam, leaving `internal/store` as a doc-only placeholder rather than a
    business persistence owner
  - the remaining backend hotspot is now deeper ownership work around `server/internal/ent/**` and migration
    boundaries, not the already-removed shared store seams

## Shared Hotspots

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/**`
- `server/internal/pluginapi/**`
- `server/internal/contract/**`
- `server/internal/pluginregistry/generated.go`
- `server/cmd/graft/**`
- `server/internal/ent/**`
- `server/internal/ent/migrate/migrations/**`
- `server/plugins/*/ent/**`
- `server/plugins/*/migrations/**`
- `web/src/app/**`
- `web/src/shared/**`
- `web/src/router/index.ts`
- `web/src/layouts/**`
- `web/src/store/modules/user.ts`
- `web/src/store/modules/permission.ts`
- `web/src/locales/**`

## Active Risks

- If future backend slices reopen `server/internal/store/**` or shared `server/internal/ent/**` as business landing
  zones, the first real multi-worktree merge wave will recreate avoidable hotspot churn.
- If future frontend slices reintroduce module truth outside `web/src/modules/<name>/**`, the `web` ownership freeze
  becomes unenforceable.
- If the repository root branch changes again and `ai-plan/public/README.md` is not updated in the same slice, future
  startup recovery will land on stale branch/worktree assumptions.
- If the first dedicated worktree/topic pair is created without an explicit owned-scope definition, this governance
  topic will continue accumulating feature-specific history that belongs elsewhere.

## Latest Validation

- Latest backend validation carried by the active baseline before this compaction:
  - `cd server && go test ./plugins/rbac ./plugins/user`
  - `cd server && go test ./internal/store/...`
  - `cd server && env GOCACHE=/tmp/go-build go run ./cmd/graft validate backend --stage lint`
- This compaction slice rechecked the recovery-path shape with:
  - `git show --stat --oneline --decorate=short 654c791 5f45b31 799f1ff 866582a --`
  - `git diff -- ai-plan/public/multi-worktree-governance`
  - `wc -l ai-plan/public/multi-worktree-governance/todos/multi-worktree-governance-tracking.md ai-plan/public/multi-worktree-governance/traces/multi-worktree-governance-trace.md`

## Archive Pointers

- Pre-compaction tracking snapshot:
  `ai-plan/public/multi-worktree-governance/archive/todos/multi-worktree-governance-tracking-pre-compaction-2026-05-19.md`
- Pre-compaction trace snapshot:
  `ai-plan/public/multi-worktree-governance/archive/traces/multi-worktree-governance-trace-pre-compaction-2026-05-19.md`

## Immediate Next Step

- Keep `multi-worktree-governance` limited to shared baseline governance while the repository root remains the only
  active worktree.
- For the next backend slice, continue reducing deeper `internal/ent/**` and migration ownership hotspots without
  weakening the frozen `rbac` ownership over `user_roles`.
- Before creating the first dedicated long-lived worktree/topic pair, record its owned scope and shared-hotspot policy
  in `ai-plan/public/README.md` and give it its own tracking/trace files instead of extending this governance topic.
