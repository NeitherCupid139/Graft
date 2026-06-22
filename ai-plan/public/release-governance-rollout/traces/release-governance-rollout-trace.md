# Release Governance Rollout Trace

## 2026-06-22 Phase 0 topic bootstrap

- 已确认 `ai-plan/public/release-readiness-governance-audit/**` 的 docs-only 审计主题已经 `archive-ready`，不应继续作为新的 loop 承载体。
- 已将 `release-readiness-governance-audit` 从 active topic 索引迁出，并迁入 `ai-plan/public/archive/` 作为上游审计证据。
- 新建 `release-governance-rollout` 作为后续 `v0.1.0 P0` 治理落地的 active topic。
- 固定后续 loop 为 `topic-completion-loop`，并采用串行单 worker round：
  - `phase-1-release-safety-governance`
  - `phase-2-release-identity-and-policy`
  - `phase-3-release-operator-docs-baseline`
- 已固定全局 guardrail：
  - 当前 topic 不直接实施 release workflow、Docker/Compose、Kubernetes、托管平台支持
  - 当前 topic 先收口 authority 和文档口径，再决定是否开实现型 topic

## Suggested Implementation Order

- 第一优先级：Release Safety Governance
  - 先固定 migration/config/upgrade safety governance，避免后续版本与文档治理建立在不稳定的 operator path 上
- 第二优先级：Release Identity And Policy
  - 再固定 `BuildInfo`、`graft version`、release policy 和 support boundary，避免文档引用不存在的 version contract
- 第三优先级：Release Operator Docs Baseline
  - 最后把前两批 authority 真正沉淀为最小 operator 文档集合

## 2026-06-22 Phase 1 release safety governance

- 已使用 `server/internal/cli/{migrate,serve,validate}.go` 与 `server/internal/config/config.go` 作为 authority
  evidence，确认当前仓库的真实运行边界是：
  - `graft migrate up` / `graft dev` 承担显式迁移入口
  - `graft serve` 保持纯运行时启动，不隐式迁移
- 已把 `v0.1.0` migration safety governance 固定为：
  - forward-only live migration governance
  - 升级前必须先验证 backup/restore readiness
  - rollback 仅承诺文档化的 operator decision points、data risk 与最小验证步骤
  - 不承诺自动数据库回滚、自动配置回滚或 helper tooling
- 已把 `v0.1.0` config compatibility governance 固定为：
  - `additive` / `default-change` / `rename` / `semantic-change` / `removal`
  - patch release 不允许静默 rename/removal/semantic change
  - minor release 的高风险 config 变更必须同时记录 replacement、removal target、operator action、
    release notes 与 upgrade notes
  - startup warning、alias bridge、config rewrite helper 保持 deferred
- 已把 operator upgrade path 的最小治理口径固定为：
  - 版本目标确认
  - backup/restore readiness
  - config diff 与 operator action 检查
  - 显式 migration step
  - post-upgrade verification
- 已修复 topic authority drift：
  - 补齐 `ai-plan/public/release-governance-rollout/startup-prompt.md`
  - 使 topic README 的 recovery 引用重新可用

## 2026-06-22 Phase 2 release identity and policy

- 已使用 `server/internal/cli/root.go`、`.github/workflows/release.yml`、`.github/workflows/publish.yml`、
  `web/package.json` 与 `README.md` 作为 authority evidence，确认当前仓库的真实 release identity 状态是：
  - 发布 tag 与 artifact 名称已经存在 authority
  - `server` 端尚无统一 BuildInfo 注入模型
  - `graft` 根命令当前还没有 `version` subcommand
- 已把 `v0.1.0` release identity baseline 固定为：
  - 官方 release identity 以 Git tag `vMAJOR.MINOR.PATCH` 为准
  - future `BuildInfo` 最小字段集固定为 `version`、`git_commit`、`build_time_utc`、`git_tree_state`
  - `BuildInfo.version` 使用 bare semver；release tag 保持 `v` 前缀
  - 在 BuildInfo / `graft version` 真实实现落地前，tag、artifact filename 与 release notes 共同构成当前
    operator-facing canonical identity
- 已把 future `graft version` 的最小 contract 固定为：
  - 纯 metadata readout
  - 不依赖数据库、Redis、HTTP 启动或 migration 执行
  - release build 至少暴露 `version`、`git_commit`、`build_time_utc`、`git_tree_state`
  - 本批次不把该命令误写成“已存在支持”
- 已把 `v0.1.0` release policy / support boundary 固定为：
  - 一个 active repository release line
  - 不承诺 LTS、多 minor 并行维护或独立 `server` / `web` 官方 release cadence
  - 官方 operator 支持包是同一 tag 下的 `server` artifact、`web` artifact 与 release notes
- 已把 `server` / `web` / migration version coordination 固定为：
  - `server` / `web` / release notes 必须来自同一 release tag
  - 混用不同 tag 的 `server` / `web` artifact 在 `v0.1.0` 下不受支持
  - migration version 仅是内部排序号，不是 product version 或 compatibility label

## 2026-06-22 Phase 3 release operator docs baseline

- 已先执行 coverage check，确认 Phase 1/2 已覆盖基础 authority，但 operator-facing 落点仍缺：
  - supported vs unsupported upgrade path
  - migration class as operator guidance
  - config lifecycle as consumable notes
  - BuildInfo visibility contract across `CLI` / `API` / `logs`
  - versioning/support boundary mapping
- 已新增最小 operator 文档集合：
  - `operator-docs/README.md`
  - `operator-docs/install.md`
  - `operator-docs/config-reference.md`
  - `operator-docs/upgrade.md`
  - `operator-docs/rollback.md`
  - `operator-docs/release-notes-template.md`
- 已把 `Installation` / `Configuration` / `Upgrade` / `Rollback` / `Release notes` / `Versioning` /
  `Support boundary` 固定到单一 canonical mapping，避免散落在 topic 摘要里不可直接消费
- 已把 `Upgrade Safety Boundary` 明确为：
  - 只支持 official repository tag 到 official repository tag 的升级路径
  - mixed-tag 部署、隐式 startup migration、跳过 release notes/config review 属于 unsupported or discouraged path
  - operator 必须承担 backup/restore readiness、config snapshot、explicit migration 与 post-upgrade verification
- 已把 `Migration Governance Details` 明确为：
  - `additive`
  - `compatible`
  - `destructive`
  - 并固定其对 `patch` / `minor` / `major` 的最小 release boundary
- 已把 `Configuration Lifecycle` 明确为：
  - 新配置默认值原则
  - deprecation record baseline
  - rename/removal 的 operator 手工动作边界
  - 不把 alias bridge、startup warning 或 rewrite helper 伪装成现成能力
- 已把 `Build Identity Visibility` 明确为：
  - future `graft version` 是 canonical `CLI` surface
  - `API` / `logs` 目前都不是 `v0.1.0` promised canonical identity surface
  - 在实现前，release tag、artifact name、release notes 仍是 operator-facing canonical identity
- 已把 `experimental` 固定为 explicit label rule：
  - 只有 release notes 或 operator docs 明确标记 experimental，operator 才能按 experimental 理解
