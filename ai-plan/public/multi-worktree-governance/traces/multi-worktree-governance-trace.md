# Multi Worktree Governance Trace

## 2026-05-19 long-lived worktree mapping rules made explicit

- Expanded the active public recovery docs so the root-only state is operationally clear:
  - the repository root remains the only active worktree
  - the root is acting as shared-baseline governance, not as a permanent feature-owned worktree
  - future long-lived worktrees must be mapped with their own active topic plus tracking/trace pair before feature
    recovery moves there
- Added an explicit shared-hotspot handling policy to the active recovery path:
  - shared hotspots remain limited and opt-in
  - plugin-owned worktrees should not take standing ownership of `internal/ent/**`,
    `internal/ent/migrate/migrations/**`, or `internal/pluginregistry/generated.go`
  - hotspot updates should be split into bounded governance slices or serialized instead of being co-owned by default

## 2026-05-19 active recovery compaction

- Archived the previous active tracking and trace files into topic-local snapshots because the default recovery path had
  grown past the point where it was useful as a startup entry:
  - `ai-plan/public/multi-worktree-governance/archive/todos/multi-worktree-governance-tracking-pre-compaction-2026-05-19.md`
  - `ai-plan/public/multi-worktree-governance/archive/traces/multi-worktree-governance-trace-pre-compaction-2026-05-19.md`
- Replaced the active tracking file with a short recovery entry that keeps only:
  - current branch/worktree truth
  - frozen ownership baselines
  - shared hotspots
  - active risks
  - latest validation
  - immediate next step

## 2026-05-19 backend boundary cleanup landed

- `654c791` moved audit persistence into plugin-local storage surfaces under `server/plugins/audit/**`, so audit no longer
  relies on the shared store path for its live repository wiring.
- `5f45b31` removed the shared audit compatibility shim from `internal/store`, closing the last shared audit transition path.
- `799f1ff` removed the shared user store compatibility seam, so user tests and reset helpers now exercise plugin-local
  store contracts instead of the old shared adapter.
- `866582a` removed the shared user/auth seam and left `internal/store` as documentation-only scaffolding, which means the
  active backend hotspot is no longer shared store cleanup but the deeper `internal/ent/**` and migration ownership boundary.

## Active Baseline After Compaction

- `mvp-extension-path` stays archived and is no longer part of the active recovery path.
- The repository root on `refactor/server-module-boundaries` remains the only active worktree.
- `web` baseline stays frozen on:
  - shell-owned `app/layouts/router/config/locales/platform-contract` surfaces
  - module-owned `web/src/modules/<name>/**`
  - shared-owned `web/src/shared/**`
- `server` baseline stays frozen on:
  - compile-time modular monolith
  - plugin-first ownership under `server/plugins/<name>/**`
  - shared backend boundary at `internal/pluginapi/**` and `internal/contract/**`
  - generated shared hotspot at `internal/pluginregistry/generated.go`
  - `user_roles -> rbac` ownership
- The latest landed backend milestone is now the full shared-store seam cleanup for audit, user, and user/auth:
  - live RBAC persistence is plugin-local
  - transitional shared audit/user compatibility paths are removed
  - `internal/store` is no longer a business persistence landing zone

## Historical Detail Pointer

- Full milestone history from `2026-05-17` through the pre-compaction `2026-05-19` slices now lives in:
  `ai-plan/public/multi-worktree-governance/archive/traces/multi-worktree-governance-trace-pre-compaction-2026-05-19.md`
- Use that snapshot only when a task explicitly needs older validation logs, intermediate migration notes, or the full
  chronology of the web/server/docs governance slices.

## Immediate Next Step

- Keep this topic focused on shared baseline governance until the first real dedicated worktree/topic pair exists.
- When the next slice becomes feature-owned instead of governance-owned, create a dedicated topic entry rather than
  appending another long phase ledger here.
