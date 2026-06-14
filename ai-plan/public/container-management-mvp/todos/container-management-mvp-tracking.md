# Container Management MVP Tracking

## Topic

Container Management MVP

## Scope

实现“运维管理 -> 容器管理”模块，第一版通过本机 Docker socket 支持容器列表、详情、日志读取和
`start` / `stop` / `restart`，并接入权限、审计、系统配置与前端管理页面。

## Repository Truth

- `AGENTS.md`
- `server/AGENTS.md`
- `web/AGENTS.md`
- `ai-plan/design/容器管理设计.md`
- `ai-plan/design/项目设计.md`
- `ai-plan/design/模块与依赖注入设计.md`
- `ai-plan/design/前端架构设计.md`
- `ai-plan/design/契约治理与魔法值治理规范.md`
- `ai-plan/design/系统配置模型与渲染设计.md`
- `ai-plan/design/TDesign-MCP-辅助开发规范.md`
- `ai-plan/design/AI任务追踪与恢复设计.md`

## Authority Discovery

- `server/modules/container/**` 是后端业务模块 authority。
- `openapi/**` 是 wire contract authority。
- `web/src/modules/container/**` 是前端业务模块 authority。
- 权限、菜单、错误 key、route fragment、message key 属于模块 contract。
- `DockerRuntime` 只是第一版 runtime adapter，不是业务命名前缀。
- 第一版不建数据库表；容器状态来自 Docker API，操作记录复用 audit log。

## Current Recovery Point

- Phase 0 已完成：设计文档与 public topic 持久化。
- 当前分支：`feat/container-management-mvp`。
- 下一批次：Phase 1 OpenAPI 与 contract source。
- 下一批次不得实现完整 DockerRuntime 或前端页面；只建立 OpenAPI source、generated artifacts 和最小 contract 常量。

## Task Checklist

- [x] Phase 0：设计和 topic 持久化
- [ ] Phase 1：OpenAPI 与 contract source
- [ ] Phase 2：后端模块骨架、菜单、权限、i18n、系统配置定义
- [ ] Phase 3：后端 Docker runtime、API、权限、错误映射、审计
- [ ] Phase 4：前端容器列表、详情 Drawer、日志 Drawer/Dialog、高危确认
- [ ] Phase 5：测试、i18n、治理收尾、归档准备

## Risks

- Docker socket 权限依赖部署用户和宿主机配置。
- `dangerous_actions_enabled` 默认关闭，写操作必须被显式启用。
- container id/name path 参数必须严格 encode/unescape/validate。
- 日志读取必须限制 `tail <= 2000` 并设置超时。
- OpenAPI source、generated Go、generated TypeScript 和前端消费层容易 drift，必须按 authority-first 顺序推进。
- 不得泄露敏感 env、secret、token 或 Docker inspect 原始敏感字段。

## Last Validation

- Phase 0 只修改 `ai-plan/**` 文档。
- 直接验证目标：
  - `git diff --check`
  - `rg "Docker 管理|ops\\.docker" ai-plan`，并确认命中只出现在禁止项或验证说明中。
  - `rg "容器管理|ops\\.container|/api/ops/containers" ai-plan`
