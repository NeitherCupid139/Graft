# Config Reference And Lifecycle Baseline

本文件定义 `v0.1.0` operator-facing config compatibility baseline。它固定可长期维护的说明结构，但不伪装成
machine-readable config inventory 或自动兼容桥接实现。

## Authority

- Phase 1 config safety authority：`../README.md`
- repository runtime baseline：`/README.md`

## Stable Config Change Classes

- `additive`
  - 新增稳定配置项，不移除既有 key
- `default-change`
  - 既有 key 保持不变，但默认值或推荐值发生调整
- `rename`
  - canonical key 发生替换，旧 key 进入兼容治理
- `semantic-change`
  - key 名称不变，但语义、单位、取值解释或安全后果发生变化
- `removal`
  - 稳定 key 被计划移除

## Default Value Principle

- 新配置项应优先提供安全、最小惊讶、可文档化的默认值
- 如果没有合理默认值，operator 文档必须明确标记为 required input，而不是假定隐式环境值
- `default-change` 必须记录行为影响与 operator 是否需要显式覆盖

## Lifecycle Rules

- `patch`
  - may introduce `additive`
  - may introduce low-risk `default-change` only when behavior impact is explicitly documented
  - must not silently introduce `rename`、`semantic-change`、`removal`
- `minor`
  - may introduce `rename`、`semantic-change`、`removal` only with release notes and upgrade notes
  - must record replacement, removal target, and operator action
- `major`
  - may carry intentionally incompatible config lifecycle changes as a planned boundary
  - still requires explicit replacement and operator guidance

## Deprecation Record Baseline

For every governed `rename`、`semantic-change`、`removal`, keep at least:

- config key or config group
- canonical owner
- change class
- deprecated_in
- removal_target
- replacement
- operator action required
- release-notes required
- upgrade-notes required

## Rename And Removal Compatibility

- `v0.1.0` does not assume startup deprecation warnings, alias bridges, or config rewrite helpers
- when a rename happens in a supported release, operator docs must treat the new key as canonical and must describe the
  manual migration action
- removal must not be described as safe by default; operator docs must state when old keys stop being accepted

## Operator Responsibility

- keep a pre-change config snapshot before upgrade
- compare the target release notes against current config usage
- update keys and values using canonical documentation rather than inferred legacy behavior
