# AI Plan Public Recovery Index

`ai-plan/public/README.md` is the shared recovery index used after `AGENTS.md` startup preflight. It should stay short,
list only active topics, and help the current branch or worktree land on the right recovery documents without scanning
every public artifact.

## Recovery Rules

1. Read this file only after startup preflight from the root `AGENTS.md`.
2. If the current branch or worktree appears in the map below, read the listed topics in priority order.
3. Read the parent topic tracking and trace files before reading any subtopic files.
4. If the parent topic defines subtopics, continue into the relevant `server` / `web` / other bounded subtopic based on
   the current task shape.
5. If there is no match, fall back to scanning active topic directories.
6. Ignore `ai-plan/public/archive/**` by default unless historical context is explicitly requested.

## Mapping Conventions

- Prefer recording both `Worktree` and `Branch` when a long-lived local worktree already exists.
- If a long-lived worktree has not been created yet, branch-only mapping is acceptable as a temporary state.
- Short-lived branches for hotfixes or narrow fixes should not be added as default active-topic mappings.
- When a future topic is planned but its worktree has not been created yet, keep the current active-topic mapping stable
  and record the intended split in the relevant tracking documents first.

## Active Topics

- `multi-worktree-governance`
  - Purpose: govern the repository on local `main` so shared contracts, recovery rules, and ownership boundaries are
    ready before spawning long-lived worktrees from local branches.
  - Tracking: `ai-plan/public/multi-worktree-governance/todos/multi-worktree-governance-tracking.md`
  - Trace: `ai-plan/public/multi-worktree-governance/traces/multi-worktree-governance-trace.md`
  - Recovery note: use this topic only for shared baseline governance on `main`; once a long-lived worktree/topic pair
    is actually created, map that worktree to its own active topic instead of keeping the work on `main`.

## Worktree To Active Topic Map

- Worktree: `primary-main`
- Branch: `main`
  - Priority 1: `multi-worktree-governance`
  - Note: `primary-main` is the logical label for the repository root worktree on local `main`; it exists only to
    track shared baseline governance before dedicated long-lived worktrees are created.
