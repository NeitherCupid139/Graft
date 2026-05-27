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

## 2026-05-27 Batch 2 aligned bootstrap menus and dynamic routes to canonical access-control paths

- Executed Batch 2 as a delegated worker round under `graft-multi-agent-loop`.
- Accepted the worker-owned breaking migration decision already authorized for this topic:
  - keep only `/access-control/users`, `/access-control/roles`, `/access-control/permissions`
  - remove frontend compatibility for `/users`, `/roles`, `/permissions`
  - remove frontend historical access-control `title_key` rewrite compatibility
- Tightened backend bootstrap truth instead of preserving frontend adapter magic:
  - RBAC menu registration now declares explicit `/access-control` root with localized `menu.access_control.title`
  - RBAC bootstrap-related icon metadata is now declared canonically in the backend menu registry path
  - bootstrap menu responses are now emitted in stable access-control-first order
- Simplified frontend route transformation:
  - removed legacy access-control path normalization
  - removed historical title-key rewrite compatibility
  - removed legacy path constants from the access-control bootstrap contract
- Revalidated the owned-scope implementation directly:
  - `cd web && bun run check`
  - `cd server && go test ./plugins/rbac ./plugins/user`
  - `git diff --check`
- All above validation passed.
- The next governance drift is now button-level permission visibility standardization rather than route/menu path truth.

## 2026-05-27 Batch 3 tightened critical button-level permission visibility on owned web surfaces

- Executed Batch 3 as a delegated worker round under `graft-multi-agent-loop`.
- Kept the slice within the owned `web` scope and Option A only; no menu CRUD, resource CRUD, or DTO truth expansion.
- Standardized critical visibility on the existing `v-permission` directive for the owned RBAC and access-control surfaces:
  - role create
  - role edit
  - role permission assignment
  - access-control overview entry actions
  - access-control overview quick links for users, roles, and permissions
- Matched the frontend visibility gates to the backend guard semantics already declared on:
  - `role.create`
  - `role.update`
  - `role.read`
  - `role.permission.assign`
  - `permission.read`
  - `user.read`
  - `user.create`
- Tightened access-control overview loading behavior so the page no longer calls guarded user/role/permission APIs when the
  current session lacks the corresponding read permission; this prevents overview load failures caused by hidden-but-still-
  requested protected endpoints.
- Added focused tests for:
  - directive-driven hiding of role create and assign-permissions actions
  - directive-driven hiding of access-control permission entry points
  - permission-aware overview fetch suppression when read permissions are missing
- Explicitly observed remaining out-of-scope drift:
  - `web/src/modules/user/pages/index.vue` still relies on page-local permission booleans and dropdown disable states
    rather than the same explicit `v-permission` visibility pattern for all dangerous actions
  - that drift was not edited because the delegated round write scope did not include `web/src/modules/user/**`
- Revalidated the owned-scope implementation directly:
  - `cd web && bun run check`
  - `git diff --check`
- All above validation passed.

## 2026-05-27 Batch 3 standardized critical visibility on owned RBAC/access-control surfaces

- Executed Batch 3 as a delegated worker round under `graft-multi-agent-loop`.
- Accepted the worker result because it stayed inside the declared owned scope and matched the current Option A target.
- Standardized key owned-scope visibility on the canonical `v-permission` path:
  - role create
  - role edit
  - role permission assignment
  - access-control overview entry actions
  - access-control overview quick links for users, roles, and permissions
- Confirmed the worker also tightened behavior, not just presentation:
  - the access-control overview now suppresses guarded fetches when the current session lacks the matching read
    permission, so the page no longer performs known-protected reads merely to fail behind backend guards
- Revalidated the owned-scope implementation directly:
  - `cd web && bun run check`
  - `git diff --check`
- All above validation passed.
- Remaining known drift after this round:
  - `web/src/modules/user/pages/index.vue` still needs the same visibility-governance tightening on dangerous actions
  - backend guard consistency across RBAC and adjacent management routes still needs its dedicated Batch 4 audit
