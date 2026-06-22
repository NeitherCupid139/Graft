# Install Guide Baseline

本文件只定义 `v0.1.0` operator install baseline。它不承诺 `Docker`、`Compose`、`Kubernetes`、托管平台或一键安装器。

## Authority

- release identity authority：`../README.md` 的 Phase 2 accepted authority
- migration and config safety authority：`../README.md` 的 Phase 1 accepted authority

## Supported Install Shape

- install one official `server` artifact and one official `web` artifact from the same release tag
- prepare PostgreSQL and Redis before runtime startup
- prepare runtime config from the canonical server config keys before first startup
- use `graft migrate up` or `graft dev` for explicit schema application before normal `graft serve` startup when the
  target release includes schema changes

## Unsupported Install Assumptions

- mixing `server` and `web` artifacts from different official release tags
- expecting `graft serve` to auto-apply migrations
- relying on unpublished deployment manifests as official install guidance

## Minimum Prerequisites

- target release tag identified
- matching `server` and `web` release artifacts identified
- PostgreSQL connectivity verified
- Redis connectivity verified
- canonical config values prepared

## Initial Operator Checklist

1. Confirm the target release tag and matching artifacts.
2. Prepare the config file or environment values from canonical config keys only.
3. Run explicit migration before normal runtime startup when the release notes call for schema change.
4. Start the runtime with `graft serve` only after prerequisites and migrations are complete.
