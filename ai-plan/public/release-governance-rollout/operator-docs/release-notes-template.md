# Release Notes Template Baseline

本模板是 `v0.1.0` 最小 release communication contract。它保证 operator 能直接消费 version、migration、config、
upgrade、rollback 和 support boundary 信息，但不代表这些字段已经被 workflow 自动生成。

## Header

- Release tag: `vX.Y.Z`
- Release date:
- Repository commit:
- Matching `server` artifact:
- Matching `web` artifact:

## Release Identity

- canonical release identity: repository Git tag
- Build identity visibility:
  - `CLI`
    - future `graft version` target fields: `version`、`git_commit`、`build_time_utc`、`git_tree_state`
  - `API`
    - operator-facing version endpoint support: `not promised in v0.1.0` or explicit status
  - `logs`
    - startup-log BuildInfo support: `not promised in v0.1.0` or explicit status

## Compatibility Summary

- release type: `patch` / `minor` / `major`
- breaking change status:
- supported upgrade path:
- unsupported mixed-version or mixed-tag cases:

## Migration Impact

- migration required: `yes/no`
- migration class: `additive` / `compatible` / `destructive`
- explicit migration entrypoint:
- backup/restore prerequisite:

## Config Impact

- config change classes involved:
- canonical keys affected:
- replacement or operator action:
- deprecated_in:
- removal_target:

## Upgrade Notes

- pre-upgrade checklist:
- runtime startup order:
- minimum post-upgrade verification:

## Rollback Notes

- rollback decision point:
- backup/config snapshot dependency:
- minimum post-rollback verification:

## Support Boundary

- officially supported deployment shape:
- unsupported deployment assumptions:
- experimental items explicitly labeled:
