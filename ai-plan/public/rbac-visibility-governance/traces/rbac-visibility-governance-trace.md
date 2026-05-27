# RBAC Visibility Governance Trace

## 2026-05-27 governance topic initialized

- Re-ran the current-turn startup preflight under root `AGENTS.md` for a `cross-boundary` slice.
- Read:
  - `AGENTS.md`
  - `server/AGENTS.md`
  - `web/AGENTS.md`
  - `ai-plan/public/README.md`
  - `ai-plan/public/rbac-further-development/traces/rbac-further-development-trace.md`
  - `ai-plan/public/rbac-further-development/todos/rbac-further-development-tracking.md`
  - `ai-plan/design/AI任务追踪与恢复设计.md`
- Confirmed the recovery index still listed no active topic even though the active implementation line had shifted to an RBAC visibility-governance direction on branch `feat/wt-rbac-further-development`.
- Opened `rbac-visibility-governance` as the new active topic for this worktree and branch pair.
- Recorded explicit guardrails for the topic:
  - no menu CRUD
  - no resource CRUD
  - no resource table
  - no migration of menu canonical truth from registry/bootstrap into database-owned CRUD
  - no reverse-parsed persisted resource model from permission codes
- Set the first planned delegated round to a read-only baseline audit of the current visibility chain.

## 2026-05-27 Batch 1 baseline audit mapped the current visibility chain

- Executed Batch 1 as a read-only delegated round under `graft-multi-agent-loop`.
- The delegated round stayed within owned scope and made no file edits.
- Confirmed the current closure path is already implemented end to end:
  - permission declaration and stable permission-code contracts on `server`
  - request-time API guard through `server/internal/httpx.RequirePermission`
  - permission-filtered bootstrap menus in `server/plugins/user/bootstrap.go`
  - bootstrap snapshot recovery and dynamic route mounting in `web/src/app/bootstrap/route-guards.ts`
  - bootstrap-menu-driven async route construction in `web/src/store/modules/permission.ts`
  - localized menu title resolution and layout navigation rendering in `web/src/utils/route/**` and `web/src/layouts/**`
  - button-level element visibility via existing `v-permission` infrastructure and page-local computed capability flags
- Confirmed the primary governance drift is now concentrated in frontend compatibility logic rather than backend API authorization:
  - `web` still normalizes legacy `/users`, `/roles`, `/permissions` paths into `/access-control/*`
  - `web` still synthesizes access-control hierarchy nodes that the backend does not declare explicitly
  - `web` still rewrites legacy `title_key` values into access-control-specific keys
  - critical RBAC/user button visibility is not yet standardized on `v-permission`
  - one frontend permission-name alias still maps two semantic names to the same backend permission code
- Accepted the delegated recommendation that Batch 2 should focus on canonical bootstrap menu and dynamic route alignment under Option A only.
