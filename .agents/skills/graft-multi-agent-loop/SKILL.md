---
name: graft-multi-agent-loop
description: Repository-specific loop orchestrator for Graft multi-agent tasks. Use when one bounded task should run through repeated main-agent-managed rounds of `graft-multi-agent-task` until closeout stops emitting a next-session startup prompt or an execution budget stops the loop.
---

# Graft Multi-Agent Loop

Use this skill when a `Graft` task should run as a sequence of bounded delegated rounds under one main-agent session.

Treat root `AGENTS.md` as the only governance source. This skill is only an outer automation wrapper around
`graft-multi-agent-task`; it does not define a second startup path, a second validation contract, or a second commit
workflow.

## When To Use

Use this skill when all of the following are true:

* the task should be executed through `graft-multi-agent-task`
* the task may require multiple future-session handoffs before it is actually complete
* you want the main agent to keep coordinating delegated rounds until closeout says to stop or a budget is exhausted

Typical triggers:

* `run this as a looped multi-agent task`
* `continue this multi-agent task automatically until it finishes`
* `use graft-multi-agent-loop for this bounded slice`

## Workflow

1. Ensure the current turn already has the startup receipt required by root `AGENTS.md`.
2. Confirm the owned scope and explicit budget before starting the loop:
   - `max_rounds`
   - `max_files_changed`
   - `max_commits`
   - `max_runtime_minutes`
   - `allowed_scopes`
   - validation failure policy
3. Keep the critical path in the main agent and use delegated rounds only for bounded execution slices:
   - use `graft-multi-agent-task` to define each round
   - use `graft-multi-agent-batch` only when the round itself benefits from parallel subagent work
   - keep budget tracking, stop conditions, and final acceptance in the main agent
4. For each round, require the delegated execution path to:
   - run the slice through `$graft-multi-agent-task`
   - keep closeout human-readable
   - emit the required machine-readable JSON closeout result
5. Let the main agent decide whether to continue based on:
   - closeout JSON
   - the presence or absence of `Next-session startup prompt:`
   - repeated prompts
   - scope expansion
   - risk level
   - remaining budget
6. If a delegated round stalls, omits closeout, or returns contradictory closeout:
   - do not silently continue as an ordinary task
   - either mark the loop `blocked`, or recover the coherent owned slice locally and then close out explicitly through
     `graft-task-closeout`
   - keep the stop reason explicit in the final closeout
7. Stop when:
   - no further next-session startup prompt is emitted
   - the closeout JSON says `continue: false`
   - a budget limit is exhausted
   - validation fails under a stop-on-failure policy
   - the delegated round expands scope or reports high risk

## Output Contract

Every delegated round run through this loop must end with:

1. a concise human-readable closeout
2. `Next-session startup prompt: <prompt>` when a future round is required
3. a fenced JSON block containing:
   - `closeout_status`
   - `continue`
   - `next_prompt`
   - `stop_reason`
   - `validation`
   - `commit`
   - `consumed_budget`
   - `remaining_budget`
   - `scope_expanded`
   - `risk_level`

The main agent treats the JSON block as the primary control surface. The keyword line is a human-readable mirror, not a
replacement control plane.

## Boundaries

* do not use this skill as a substitute for `graft-boot`
* do not bypass `graft-multi-agent-task`; this skill only orchestrates repeated delegated rounds of it
* do not let the loop broaden ownership beyond the declared `allowed_scopes`
* do not treat the loop as permission to skip closeout, validation, or scoped commit rules
* do not let a stalled or malformed delegated round silently downgrade into untracked main-agent execution
* do not assume a delegated round can inherit unstated governance; the round prompt must restate the inherited context
