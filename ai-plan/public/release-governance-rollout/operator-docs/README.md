# Release Operator Docs Baseline

本目录是 `release-governance-rollout` 在 Phase 3 固定的最小 operator 文档集合。它服务于 `v0.1.0` 的
documentation-first release governance，不是产品内置帮助中心，也不是 docs-site 信息架构。

## Authority

- Phase 1 safety authority：
  - `ai-plan/public/release-governance-rollout/README.md`
  - `ai-plan/design/数据库表设计与迁移规范.md`
  - `ai-plan/design/服务端API边界与兼容治理规范.md`
- Phase 2 release identity authority：
  - `ai-plan/public/release-governance-rollout/README.md`
  - `README.md`
- operator docs consume those authorities; they do not replace them

## Canonical Mapping

| Topic | Canonical document | Purpose |
| --- | --- | --- |
| Installation | `operator-docs/install.md` | Supported install shape, prerequisites, artifact coordination |
| Configuration | `operator-docs/config-reference.md` | Config lifecycle, compatibility classes, operator config duties |
| Upgrade | `operator-docs/upgrade.md` | Supported upgrade path, migration classes, compatibility boundary |
| Rollback | `operator-docs/rollback.md` | Manual rollback decision points and minimum verification |
| Release notes | `operator-docs/release-notes-template.md` | Minimum release communication contract |
| Versioning | `operator-docs/README.md` | SemVer baseline, release identity, BuildInfo visibility |
| Support boundary | `operator-docs/README.md` | Supported, unsupported, and experimental scope |

## Release Identity And Versioning

- official release identity is the repository Git tag `vMAJOR.MINOR.PATCH`
- release-grade `server` and `web` artifacts plus release notes must come from the same release tag
- `SemVer` is the operator-facing compatibility language for repository releases
- breaking change means any release change that requires unsupported mixed-tag deployment, destructive schema action,
  incompatible stable config reinterpretation, or release-notes-defined operator intervention beyond ordinary patch
  expectations
- boundary by version type：
  - `patch`
    - bug fixes and safe additive adjustments only
    - must not silently rename, remove, or reinterpret stable config keys
    - must not require destructive migration steps
  - `minor`
    - may add features and bounded compatibility-managed changes
    - any `rename`、`semantic-change`、`removal`, or destructive operator action must be documented in release notes and
      upgrade notes
  - `major`
    - the first place where intentionally incompatible release governance changes may be planned as a normal release
      path
    - still requires explicit operator guidance and migration/rollback notes

## Build Identity Visibility

The minimum future BuildInfo baseline remains `version`、`git_commit`、`build_time_utc`、`git_tree_state`.

- `CLI`
  - future `graft version` is the canonical metadata readout once implemented
  - it must not require PostgreSQL, Redis, HTTP startup, or migration execution
- `API`
  - `v0.1.0` does not yet promise a dedicated operator-facing version API
- `logs`
  - `v0.1.0` does not yet promise startup-log BuildInfo emission as an authoritative support surface

Until those surfaces exist, the canonical operator-facing identity is:

- release tag
- published artifact names
- release notes

## Support Boundary

### Supported In `v0.1.0`

- one active repository release line at a time
- self-managed deployment of one `server` artifact and one `web` artifact from the same release tag
- explicit operator-run migration step through `graft migrate up` or `graft dev`
- documentation-first install, config, upgrade, rollback, and release note guidance

### Not Yet Promised In `v0.1.0`

- `Docker` / `Compose` / `Kubernetes` / hosted deployment support matrix
- automatic rollback tooling
- implicit startup migration or startup-time schema repair
- independent `server` / `web` official release trains
- operator-facing introspection UI beyond future minimal `graft version`

### Experimental Label Rule

- a capability is only `experimental` when release notes or operator docs explicitly label it that way
- internal code paths, draft scripts, or unpublished artifacts do not create an implied support promise
