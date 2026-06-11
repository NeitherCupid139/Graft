# AI Environment Inventory

`.ai/environment/` stores generated environment truth for `Graft`.

## Files

- `tools.raw.yaml`
  - Raw, repository-relevant environment facts collected from the current machine.
- `tools.ai.yaml`
  - AI-facing summary derived from `tools.raw.yaml`.
  - Prefer reading this file first during startup and task planning.

## Refresh Commands

```bash
bash scripts/collect-dev-environment.sh --check
bash scripts/collect-dev-environment.sh --write
python3 scripts/generate-ai-environment.py
```

## Project-Local Browser Agent

`graft-web-browser-agent` uses a project-local Python environment for AI-assisted frontend screenshots and simple
browser interactions:

```bash
.agents/skills/graft-web-browser-agent/scripts/bootstrap.sh
```

Generated runtime files are intentionally ignored:

- `.ai/venv/`
- `.ai/ms-playwright/`
- `.ai/artifacts/browser/`
- `.headroom/`

Playwright MCP may be configured in the user-level Codex MCP list as an exploration aid, but the project-local browser
agent remains the auditable artifact path for screenshots, text snapshots, login summaries, and cleanup.

## MCP Inventory

The environment inventory records whether these user-level Codex MCP servers are configured:

- `codegraph`
- `tdesign`
- `context7`
- `github`
- `playwright`
- `headroom`

These entries are generated local capability facts. They do not make MCP a repository runtime dependency, CI gate, hook,
or required contributor setup.

## Headroom

Headroom is an optional local AI context compression tool. The inventory records both the `headroom` CLI and the
user-level `headroom` MCP server when present.

Recommended privacy-first wrapper entry:

```bash
HEADROOM_TELEMETRY=off .ai/venv/bin/headroom wrap codex
```

Recommended MCP entry:

```bash
.ai/venv/bin/python -m pip install "headroom-ai[proxy]"
codex mcp add headroom -- .ai/venv/bin/headroom mcp serve
```

Use Headroom for local compression, retrieval, and stats only by default. Do not enable `--memory`, `--learn`, or any
automatic write to `AGENTS.md`, Codex `instructions.md`, `ai-plan/lessons/**`, or recovery documents without a separate
human-confirmed governance slice.

Use the skill cleanup script at task closeout when the user chooses to remove browser artifacts:

```bash
.agents/skills/graft-web-browser-agent/scripts/cleanup.sh --session <session>
```

## Rules

- Do not hand-maintain `tools.raw.yaml` or `tools.ai.yaml`.
- Refresh both files when repository toolchain expectations or environment guidance change.
- Keep secrets, machine-specific credentials, and private URLs out of the inventory.
- Read `tools.ai.yaml` first during repository startup; use `tools.raw.yaml` only when the AI-facing summary is missing
  or insufficient.
- Keep the generated inventory aligned with the repository's current local toolchain so docs and automation can reference one fact source instead of restating divergent command matrices.
- The inventory is environment truth, not startup or validation governance: root `AGENTS.md` remains the only startup
  governance source, and repository entrypoints such as `graft validate backend` / `bun run check` remain validation
  truth.
