# Phase 1 Audit And Rollout Plan

## Objective

Audit the currently modeled write-route surface against the accepted `POST /api/users` sample and define the smallest
remaining phase plan needed before the repository can honestly claim it is ready for an isolated `oapi-codegen` Go
types-only spike.

## Accepted Sample Baseline

The accepted sample remains `POST /api/users`:

- backend envelope semantics stay `success`, `code`, `message`, `messageKey`, `locale`, `traceId`, `data`
- `server/internal/httpx` remains the only backend envelope owner
- `web/src/utils/request.ts` remains the only frontend transport/runtime owner
- `data.field` means the current request-contract field name for the covered write route
- field mapping stays at route or module level, not in `httpx` and not in `request.ts`

Current concrete sample evidence:

- `openapi/paths/users.yaml` includes explicit `400` examples for `username` and `password`
- `server/plugins/user/route_errors.go` maps create-user errors into stable `data.field` values
- `web/src/modules/user/error-adapter.ts` consumes structured create-user field errors locally

## Covered Route Audit

### Legend

- `aligned`
  - already meets the sample pattern closely enough for this topic
- `partial`
  - some layers are aligned, but at least one required server/openapi/web layer is still generic or missing
- `defer`
  - already modeled, but not required to settle the current write-route field/error convention

### Audit Matrix

| Route | Current state | Gap vs sample | Phase owner |
| --- | --- | --- | --- |
| `POST /api/users` | `aligned` | none for Phase 1 baseline | done |
| `POST /api/users/{id}/update` | `partial` | server already emits `data.field` for `username` and `body`, but OpenAPI `400`/`404` only use generic `error-response`, and the user page only binds structured field errors during create mode | Phase 2, 3, 4 |
| `POST /api/users/{id}/status` | `partial` | server already emits `data.field=status` or `data.field=id`, but OpenAPI has no concrete error examples and the user page still uses only generic toast handling | Phase 2, 3, 4 |
| `POST /api/users/{id}/reset-password` | `partial` | server already returns password-policy/reuse errors as `data.field=new_password`, which matches the current request contract; OpenAPI/web are still generic | Phase 3, 4 |
| `POST /api/roles` | `partial` | RBAC server already emits `data.field=name`, but OpenAPI `400` is generic and the RBAC page has no structured field-error consumption path | Phase 2, 3, 4 |
| `POST /api/roles/{id}/update` | `partial` | same as role create, plus `404` is still modeled generically | Phase 2, 3, 4 |
| `POST /api/roles/{id}/permissions/assign` | `partial` | backend already uses `data.field=permission_ids` for invalid inputs, but OpenAPI currently lacks explicit write-error coverage for this route and the RBAC permission dialog still falls back to generic error handling | Phase 2, 3, 4 |
| `POST /api/auth/login` | `defer` | already modeled, but this route currently settles credential semantics rather than the covered `data.field` convention; no additional auth work is required unless later readiness review finds an unresolved shared envelope mismatch | hold unless Phase 5 proves otherwise |

## Concrete Evidence By Layer

### Server

- `server/plugins/user/route_errors.go`
  - already centralizes user write-route error mapping
  - update/status/create use request-field-oriented `data.field`
  - reset-password already uses the request-contract field `new_password`; only OpenAPI/web still needed alignment
- `server/plugins/rbac/route_errors.go`
  - already centralizes RBAC write-route error mapping
  - current stable field outputs are route-local and explicit: `name`, `permission_ids`, or the provided `invalidField`

### OpenAPI

- `openapi/paths/users.yaml`
  - create-user already includes concrete `400` examples matching runtime semantics
- `openapi/paths/users.update.yaml`
  - `400` and `404` are still schema-only, with no concrete examples
- `openapi/paths/users.status.yaml`
  - `400` and `404` are still schema-only, with no concrete examples
- `openapi/paths/users.reset-password.yaml`
  - `400` and `404` are still schema-only, with no concrete examples
- `openapi/paths/roles.list.yaml`
  - role create `400` is still schema-only
- `openapi/paths/roles.update.yaml`
  - `400` and `404` are still schema-only
- `openapi/paths/roles.assign-permissions.yaml`
  - no explicit `400`/`404` write-error coverage is modeled yet
- `openapi/paths/auth.login.yaml`
  - modeled, but not a deciding blocker for the current request-field convention

### Web

- `web/src/modules/user/error-adapter.ts`
  - only create-user field binding is implemented today
- `web/src/modules/user/pages/index.vue`
  - create-user consumes structured field errors locally
  - update, status, and reset-password still fall back to generic error toasts
- `web/src/modules/rbac/pages/index.vue`
  - role submit and permission assignment still fall back to generic error messages
- no RBAC field-error adapter equivalent exists yet

## Readiness Verdict After Phase 1

`ready_for_oapi_codegen_types_only_spike: false`

Reason:

- the repo has one accepted write-route sample, but the broader currently modeled write surface is still only partially
  aligned
- at least one server route (`POST /api/users/{id}/reset-password`) still violates the accepted request-field naming
- OpenAPI examples/responses are not yet consistent across the covered rollout
- frontend structured field-error consumption is still limited to create-user only

## Remaining Phase Plan

### Phase 2: Shared Server-Side Alignment

Goal:

- standardize covered backend write-route error outputs without moving ownership out of plugin-local handlers or
  `httpx`

Required outcomes:

- keep `httpx` as the only envelope writer
- keep plugin-local route registration and handler ownership explicit
- verify all covered user/RBAC write routes emit stable `data.field` values only where the route can actually act on a
  request field
- add or tighten focused tests for the covered user and RBAC write routes

Expected validation:

- `git diff --check`
- `git status --short`
- `rg` scans for covered field names and message keys
- `cd server && go test ./plugins/user/...`
- `cd server && go test ./plugins/rbac/...`

### Phase 3: OpenAPI Error Responses And Examples

Goal:

- make the covered OpenAPI write-route error responses match real backend behavior

Required outcomes:

- add explicit `400` and `404` examples where the backend already has stable semantics
- add missing response entries for covered write routes such as role-permission assignment where the route already has
  stable backend behavior
- keep shared envelope schema semantics unchanged

Expected validation:

- `git diff --check`
- `git status --short`
- `rg` scans across `openapi/**`, `server/plugins/user/**`, and `server/plugins/rbac/**`
- `cd server && go run ./cmd/graft validate backend --stage openapi`
- rerun focused plugin tests if server assertions change together with spec examples

### Phase 4: Web Structured Field-Error Consumption

Goal:

- consume structured field errors locally in the covered `user` and `rbac` modules without pushing field semantics into
  `request.ts`

Required outcomes:

- keep `request.ts` as transport truth only
- extend `user` module handling beyond create-user where the route now exposes structured field errors
- add bounded RBAC module adapters or equivalent local handling for `name` and `permission_ids`
- leave auth untouched unless Phase 5 shows it is required for the final readiness verdict

Expected validation:

- `git diff --check`
- `git status --short`
- `rg` scans across touched web/module files
- `cd web && bun run test:run src/modules/user/pages/index.test.ts`
- `cd web && bun run test:run <focused rbac tests>`
- `cd web && bun run openapi:types:check`
- `cd web && bun run check` when the slice is broad enough

### Phase 5: Validation Closure And Readiness Verdict

Goal:

- prove the covered rollout is aligned and record the final go/no-go result for a future isolated Go types-only spike

Required outcomes:

- rerun the strongest honest validations for all touched server/openapi/web scope
- document final aligned route coverage
- set `ready_for_oapi_codegen_types_only_spike` to `true` only if all covered write routes match the accepted pattern
- if still `false`, list only narrow, evidence-backed blockers

Expected validation:

- `git diff --check`
- `git status --short`
- `rg` consistency scans across docs, `openapi`, `server`, and `web`
- `cd server && go test ./plugins/user/...`
- `cd server && go test ./plugins/rbac/...`
- `cd server && go run ./cmd/graft validate backend --stage openapi`
- `cd web && bun run openapi:types:check`
- `cd web && bun run check`

## Phase Sequencing Rule

Do not start the next phase from an uncommitted validated phase boundary. Each completed phase must:

1. finish the bounded slice
2. run the strongest honest validation for that slice
3. commit the validated owned scope with `$graft-commit`
4. emit the next-session startup prompt for the following phase
