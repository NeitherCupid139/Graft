# Rollback Guide Baseline

本文件定义 `v0.1.0` 的 rollback/recovery baseline。它是 documentation-first、operator-controlled 的最小口径，
不是自动数据库回滚、自动配置回滚或一键恢复承诺。

## Authority

- rollback boundary authority：`../README.md`
- release identity authority：`../README.md`

## Rollback Support Boundary

- rollback is a manual operator decision
- the repository does not promise down migration support
- the repository does not promise automatic config rollback helpers
- rollback guidance must be read together with the target release notes and upgrade notes

## Required Prerequisites

- known source release tag
- known target release tag
- verified database backup and restore capability
- preserved pre-change config snapshot
- identified rollback decision point

## Decision Points

- whether schema changes have already been applied
- whether data written after upgrade would be invalidated by restore
- whether config changes can be manually reversed using the preserved snapshot
- whether the incident can be mitigated without changing release version

## Minimum Operator Actions

1. Stop and assess whether rollback is safer than forward repair.
2. Use the preserved release tag, backup state, and config snapshot as recovery anchors.
3. Restore data and config only through the documented operator-controlled process.
4. Re-run minimum verification after recovery.

## Minimum Post-Rollback Verification

- runtime starts cleanly
- expected release artifact set is restored
- required dependent services reconnect
- critical login or health checks succeed
