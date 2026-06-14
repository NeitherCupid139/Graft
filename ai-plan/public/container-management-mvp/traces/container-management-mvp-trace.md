# Container Management MVP Trace

## 2026-06-14

- 启动任务并确认 task class 为 `cross-boundary`。
- 当前分支从 `feat/frontend-vibe-toolchain-ui-fixes` 重命名为 `feat/container-management-mvp`。
- 建立容器管理设计 authority：`ai-plan/design/容器管理设计.md`。
- 建立 public recovery topic：`ai-plan/public/container-management-mvp/README.md`。
- 明确菜单 IA 为 `运维管理 / 容器管理`，不复用既有 `服务器管理`。
- 明确业务命名使用 `container` / `ops.container`，特定运行时只作为第一版 adapter 名称 `DockerRuntime` 或配置键
  `ops.container.docker.endpoint`。
- 明确 MVP 不做 exec、终端、文件编辑、删除、prune、镜像构建/拉取/推送、容器创建、远程 Docker Host、
  Kubernetes、长周期资源采集或单独操作日志表。
- Phase 0 仅持久化设计与恢复入口，不修改 OpenAPI source、后端模块或前端页面。
- Phase 1 新增 OpenAPI `container` tag（展示语义“容器管理”）和六个固定 API path：
  `GET /api/ops/containers`、`GET /api/ops/containers/{id}`、`GET /api/ops/containers/{id}/logs`、
  `POST /api/ops/containers/{id}/start`、`POST /api/ops/containers/{id}/stop`、
  `POST /api/ops/containers/{id}/restart`。
- Phase 1 新增 `ContainerSummary`、`ContainerDetail`、`ContainerPort`、`ContainerMount`、`ContainerNetwork`、
  `ContainerLogResponse`、`ContainerActionResponse`、`ContainerRuntimeInfo` 等 OpenAPI schema，并声明详情响应不暴露敏感环境变量。
- Phase 1 同步 `openapi/dist/openapi.bundle.json`、`server/internal/contract/openapi/generated/types.gen.go`、
  `server/internal/contract/openapi/container/zz_generated.container.go` 和
  `web/src/contracts/openapi/generated/schema.ts`。
- Phase 1 未实现 runtime handlers、DockerRuntime adapter 业务逻辑或前端页面。
- Phase 1 验证通过：`cd server && go test ./internal/contract/openapi/...`、`node scripts/openapi-bundle.mjs`、
  `cd web && bun run openapi:types`、`cd web && bun run openapi:types:check`、`git diff --check`。

## Loop Batch State

```json
{
  "loop_mode": "topic-completion-loop",
  "completed_batches": [
    "phase-0-design-topic-persistence",
    "phase-1-openapi-contract-source"
  ],
  "pending_batches": [
    "phase-2-server-module-foundation",
    "phase-3-server-runtime-api-audit",
    "phase-4-web-container-management-ui",
    "phase-5-validation-governance-closeout"
  ],
  "current_batch": "phase-1-openapi-contract-source",
  "next_batch": "phase-2-server-module-foundation",
  "closeout_status": "active"
}
```
