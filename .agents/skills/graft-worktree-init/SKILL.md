---
name: graft-worktree-init
description: Repository-specific worktree creation workflow for Graft. Use when creating or rebuilding a long-lived or temporary local git worktree and the setup should follow the repository's shared local-resource rules without hard-coded root paths.
---

# Graft Worktree Init

Use this skill when a `Graft` task needs a new local git worktree.

Treat root `AGENTS.md` as the governance source of truth. This skill provides the repository-standard worktree creation
path; it does not create a second boot chain or recovery policy.

## Inputs

- required: `branch_name`
- optional: `base_branch`
  - defaults to `main`
- optional: `rebuild`
  - defaults to `false`
- optional: `repo_dir`
  - only when auto-detection fails or the user explicitly wants an override
- optional: `worktree_root`
  - only when the default sibling `<repo-name>-wt` root is not correct

## Workflow

1. Run startup preflight from `AGENTS.md` before substantive repository conclusions.
2. Inspect `.worktree-shared.json` in the repository root.
3. Use `python3 .agents/skills/graft-worktree-init/scripts/graft_worktree_init.py ...` as the only helper entrypoint.
4. Let the helper auto-detect:
   - canonical `repo_dir` from the shared git common dir
   - default `worktree_root` as the sibling `<repo-name>-wt`
5. If the helper reports that `worktree_root` cannot be inferred safely, confirm the path with the user and rerun with
   `--worktree-root`.
6. Review the printed execution plan before mutating anything.
7. After success, tell the user:
   - created worktree path
   - branch/base branch used
   - any shared local resources linked
   - any optional local resources skipped because the source file was missing

## Shared Local Resource Rules

- `.local` is legacy and must not be created or depended on.
- Shared local resources are defined only by `.worktree-shared.json`.
- The helper creates relative symlinks from the new worktree back to the canonical repository root.
- `web/.env.development` is shared through that symlink model when it exists in the canonical repo root.
- Optional local files missing in the canonical repo root should produce a warning, not a hard failure.

## Guardrails

- do not hard-code machine-specific `ROOT_DIR`, `REPO_DIR`, or `WORKTREE_ROOT`
- do not copy per-worktree local env files when a shared symlink is sufficient
- do not treat `.local` as an active repository convention
- do not write a second manifest or second shared-resource truth outside `.worktree-shared.json`
