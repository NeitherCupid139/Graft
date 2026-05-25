# OAPI Generated Server/Client Governance Spike Trace

## 2026-05-25 implementation spike

- Renamed the worktree topic branch from `feat/wt-oapi-codegen-types-only-spike` to `feat/oapi-generated-server-client-governance-spike`.
- Added a monitor-only generated OpenAPI server contract package under `server/internal/contract/openapi/monitor/**`.
- Kept the backend router ownership explicit:
  - the monitor plugin still registers `GET /api/monitor/server-status` itself
  - the generated layer only owns parameter binding and compile-time handler interface conformance
- Rejected `strict-server` for this pilot implementation because it would force the response envelope away from `httpx`.
- Updated the monitor frontend module API so the server-status call is now operation-bound to generated OpenAPI typings while still running through `request.ts`.
- Kept page/module boundaries unchanged:
  - pages still call module API helpers
  - no page directly consumes generated client/runtime code

## Validation Notes

- The generated Go server binding emits the expected OpenAPI `3.1.x` warning from `oapi-codegen`.
- That warning does not block the pilot, but it remains a real governance risk for future broader rollout.

## Next-Session Startup Prompt

```text
使用 $graft-multi-agent-loop。

governance source: root AGENTS.md
task class: cross-boundary
recovery source:
  - current repository state
  - ai-plan/public/oapi-generated-server-client-governance-spike/README.md
  - ai-plan/public/oapi-generated-server-client-governance-spike/traces/oapi-generated-server-client-governance-spike-trace.md
  - ai-plan/public/README.md
branch / worktree:
  - feat/oapi-generated-server-client-governance-spike
owned scope:
  - ai-plan/public/oapi-generated-server-client-governance-spike/**
  - server/internal/contract/openapi/monitor/**
  - server/plugins/monitor/**
  - web/src/modules/monitor/**
forbidden scope:
  - 不迁移 auth/user/rbac
  - 不修改 OpenAPI spec 语义
  - 不改 request.ts 运行时责任
  - 不改 httpx envelope/error ownership
objective:
  - 完成 monitor/server-status generated server/client 试点的后续验证、测试补强和 go/no-go closeout
validation:
  - git diff --check
  - cd web && bun run openapi:types:check
  - cd web && bun run check
  - cd server && go run ./cmd/graft validate backend
  - scripts/run_python.sh scripts/magic_value/check_magic_values.py --mode ci
```
