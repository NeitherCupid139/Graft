# Upgrade Guide Baseline

本文件定义 `v0.1.0` 的 operator upgrade baseline。它把 Phase 1/2 的治理口径转成可执行的操作边界，但不承诺
自动化 preflight、自动迁移、自动回滚或环境编排支持。

## Authority

- migration and config safety authority：`../README.md`
- release identity and versioning authority：`../README.md`

## Supported Upgrade Path

- upgrade from one official repository release tag to another official repository release tag
- keep `server` artifact, `web` artifact, and release notes aligned to the same target tag
- apply explicit migration before normal runtime startup when the target release includes schema changes
- keep database backup/restore readiness and a pre-change config snapshot before live migration

## Unsupported Or Discouraged Upgrade Path

- mixed-tag `server` / `web` deployment
- relying on `graft serve` for migration side effects
- skipping release notes, upgrade notes, or config review for releases with governed changes
- treating migration ordering numbers as product compatibility promises

## Upgrade Compatibility Principles

- forward-only migration governance is the only supported live schema evolution baseline in `v0.1.0`
- release tag is the only official product compatibility label
- migration file versions are internal ordering identifiers, not release versions
- config compatibility is governed by explicit change class, not by implicit legacy fallback

## Migration Governance Details

- `additive migration`
  - adds schema elements without invalidating the current release data model
  - acceptable in `patch` or `minor` when release notes stay honest and no destructive operator step is required
- `compatible migration`
  - changes schema shape but keeps the current operator path compatible when paired with explicit release notes and
    config/operator actions
  - acceptable in `minor`; use caution in `patch` and avoid when it creates rollout ambiguity
- `destructive migration`
  - removes or irreversibly rewrites schema/state in a way that narrows rollback choices or operator tolerance
  - not a normal `patch` change
  - only acceptable in governed `minor` or later release planning when release notes, upgrade notes, rollback decision
    points, and maintenance-window expectations are explicit
  - any incompatible destructive evolution should be treated as `major` planning by default unless repository authority
    later narrows the rule

## Operator Responsibility Boundary

1. Confirm the source release and target release tag.
2. Read release notes and upgrade notes before touching runtime state.
3. Verify backup and restore readiness.
4. Preserve a pre-change config snapshot.
5. Run explicit migration before normal startup when required.
6. Perform minimum post-upgrade verification and record rollback decision points.

## Minimum Upgrade Checklist

- source release identified
- target release identified
- matching artifacts verified
- migration requirement verified
- config changes reviewed
- backup/restore readiness verified
- post-upgrade verification plan prepared
