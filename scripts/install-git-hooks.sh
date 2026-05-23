#!/usr/bin/env sh
set -eu

repo_root="$(git rev-parse --show-toplevel)"

cd "$repo_root"

git config core.hooksPath .husky

chmod 755 .husky/commit-msg .husky/pre-commit .husky/pre-push

if [ -f "scripts/run_python.sh" ]; then
  chmod 755 "scripts/run_python.sh"
fi

printf 'git hooks installed at %s\n' "$(git config --get core.hooksPath)"
