#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
server_root="$repo_root/server"

cd "$server_root"
exec go run ./cmd/graft dev "$@"
