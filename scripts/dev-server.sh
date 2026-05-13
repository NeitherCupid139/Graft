#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
server_root="$repo_root/server"

if ! command -v atlas >/dev/null 2>&1; then
  echo "atlas CLI is required before starting the server." >&2
  echo "Run 'graft migrate up' after installing atlas, then start 'graft serve'." >&2
  exit 1
fi

cd "$server_root"

echo "==> Applying Atlas migrations"
go run ./cmd/graft migrate up

echo "==> Starting Graft server"
go run ./cmd/graft serve
