# Container Management MVP

## 当前状态摘要

- 当前主题目标是在 `Graft` 增加“运维管理 -> 容器管理”能力。
- 状态：`active`。
- 任务分类为 `cross-boundary`，涉及 OpenAPI、server module、permission/menu/audit/system-config、web module 和 i18n。
- Canonical design：`ai-plan/design/容器管理设计.md`。
- 当前仅完成 Phase 0 设计与恢复入口持久化；尚未修改 OpenAPI source、后端模块或前端页面。

## Recovery Receipt

- governance source：root `AGENTS.md`
- task class：`cross-boundary`
- recovery source：`parent topic`
- authority summary：`ai-plan/design/容器管理设计.md` + OpenAPI source + `server/modules/container` module contract/descriptor + `web/src/modules/container` bootstrap routes + permission/menu/audit/system-config/i18n governance docs

## Owned Scope

允许修改：

- `ai-plan/design/容器管理设计.md`
- `ai-plan/public/container-management-mvp/**`
- `ai-plan/public/README.md`
- `openapi/**`
- `server/modules/container/**`
- `server/internal/moduleregistry/generated.go`
- 必要的 backend generated OpenAPI contract 文件
- `web/src/modules/container/**`
- 必要的 `web` OpenAPI generated schema、module aggregation、route/menu/i18n 接入文件

禁止误触：

- 不得把产品能力命名为特定容器运行时管理。
- 不得使用 `ops.docker` 作为权限、菜单、OpenAPI tag 或前端业务前缀。
- 不得在 MVP 中实现 exec、终端、文件编辑、删除、prune、镜像构建/拉取/推送、容器创建、多节点、远程 Docker Host 或 Kubernetes。
- 不得新增单独操作日志表；操作记录复用 audit log。
- 不得把容器管理挂入既有“服务器管理”菜单。

## Phase Plan

- Phase 0：设计和 public topic 持久化。
- Phase 1：OpenAPI 与 contract source。
- Phase 2：后端模块骨架、菜单、权限、i18n、系统配置定义。
- Phase 3：后端 `DockerRuntime` adapter、API、权限、错误映射、审计。
- Phase 4：前端容器列表、详情 Drawer、日志 Drawer/Dialog、高危确认。
- Phase 5：测试、i18n、治理收尾、归档准备。

## Current Recovery Point

- 分支为 `feat/container-management-mvp`。
- 设计 authority 已落到 `ai-plan/design/容器管理设计.md`。
- active topic 已落到 `ai-plan/public/container-management-mvp/`。
- Phase 1 已完成：OpenAPI 与 generated contract artifacts。
- 当前下一批次：Phase 2 后端模块骨架、菜单、权限、i18n、系统配置定义。

## Validation Targets

```bash
git diff --check
cd server && go test ./internal/contract/openapi/...
node scripts/openapi-bundle.mjs
cd web && bun run openapi:types
cd web && bun run check
cd server && go run ./cmd/graft validate backend
```

按阶段选择最小正确验证；Phase 0 只要求文档结构和命名检查，后续实现阶段再运行 server/web 完成态验证。
