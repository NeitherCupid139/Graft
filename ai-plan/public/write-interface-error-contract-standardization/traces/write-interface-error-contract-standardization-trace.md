# Write-Interface Error Contract Standardization Trace

## 2026-05-23 `POST /api/users` sample closure in current worktree

- Bound this topic to the existing `feat/wt-openapi-contract-governance` worktree instead of creating a new worktree.
- Kept the accepted baseline unchanged:
  - `spec-first + TS-first + explicit server DTOs`
  - `web/src/utils/request.ts` remains the only frontend runtime transport truth
  - `server/internal/httpx` remains the backend envelope and localized-error owner
  - no `oapi-codegen`, `openapi-fetch`, TS runtime SDK, or generated Go server runtime wiring
- Standardized the `POST /api/users` sample decision that `data.field` uses the current request-contract field name.
- Narrowed the create-user password-policy error field from `new_password` to `password` at the user-plugin create route boundary only.
- Kept the mapping out of `httpx` and out of `request.ts`.
- Added the matching OpenAPI `400` examples for invalid argument and password-policy violation.
- Kept the follow-up conclusion unchanged: `oapi-codegen` is still deferred, and a future Go types-only spike is still blocked on broader write-interface hardening beyond this one sample.

## 2026-05-24 Phase 1 audit for covered write-route rollout

- Re-ran startup governance locally for this delegated round against root `AGENTS.md`, `.ai/environment/tools.ai.yaml`,
  `server/AGENTS.md`, `web/AGENTS.md`, `ai-plan/public/README.md`, and the active recovery traces.
- Audited the current covered rollout against the accepted `POST /api/users` sample and recorded the full matrix in
  `ai-plan/public/write-interface-error-contract-standardization/design/phase-1-audit-and-rollout-plan.md`.
- Confirmed the current Phase 1 status by route:
  - `POST /api/users`
    - aligned baseline
  - `POST /api/users/{id}/update`
    - partial: server field semantics exist, OpenAPI examples are generic, frontend update handling is still generic
  - `POST /api/users/{id}/status`
    - partial: server field semantics exist, OpenAPI examples are generic, frontend status handling is still generic
  - `POST /api/users/{id}/reset-password`
    - partial: server still emits `data.field=new_password`, which conflicts with the accepted request-field naming rule
  - `POST /api/roles`
    - partial: server field semantics exist, OpenAPI `400` is generic, frontend still uses generic error handling
  - `POST /api/roles/{id}/update`
    - partial: same gap pattern as role create, plus generic `404`
  - `POST /api/roles/{id}/permissions/assign`
    - partial: backend field semantics exist, OpenAPI write-error coverage is still incomplete, frontend still uses
      generic error handling
  - `POST /api/auth/login`
    - modeled but deferred from the rollout because it does not currently decide the write-route `data.field`
      convention
- Recorded the honest current readiness verdict:
  - `ready_for_oapi_codegen_types_only_spike: false`
- Locked the remaining loop plan as:
  - Phase 2: shared server-side field/data/error-contract alignment
  - Phase 3: OpenAPI error responses/examples alignment
  - Phase 4: web module-level structured field-error consumption alignment
  - Phase 5: focused validation closure and final readiness verdict docs

## 2026-05-24 Phase 2 server-side alignment evidence closure

- Re-checked the covered backend write routes against the actual request DTOs and route handlers before changing runtime
  code.
- Corrected the earlier Phase 1 audit note for `POST /api/users/{id}/reset-password`:
  - the current request-contract field name is already `new_password`
  - the existing backend `data.field=new_password` mapping is therefore aligned with the accepted rule that
    `data.field` must use the current request-contract field name
  - no `httpx` or plugin-runtime refactor was justified for this route in Phase 2
- Added focused user-plugin route tests proving the covered reset-password error behavior:
  - password-policy violation returns `AuthPasswordPolicyViolation` with `data.field=new_password`
  - password-reuse rejection returns `AuthPasswordReuseForbidden` with `data.field=new_password`
- Re-confirmed the covered RBAC write routes already have focused field-level error assertions for:
  - role create/update name conflicts -> `data.field=name`
  - role permission assignment invalid inputs -> `data.field=permission_ids`
- Phase 2 conclusion:
  - the covered server-side write-error contract surface is aligned without broadening ownership or introducing new
    abstractions
