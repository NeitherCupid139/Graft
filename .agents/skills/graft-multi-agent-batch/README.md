# Graft Multi-Agent Batch

`graft-multi-agent-batch` is the repository workflow for one reviewable parallel wave of bounded subagent slices.

## Main-Agent Role

The main agent coordinates the wave:

- choose disjoint slices
- dispatch bounded subagents
- review returned results
- prepare validation and next safe slices

The main agent does not take back an active delegated implementation slice just because the worker is quiet, has no
visible diff yet, or answered a checkpoint.

## Worker Role

Each write-capable `worker` owns one bounded implementation slice until it returns:

- a final closeout
- an explicit blocked state
- an owned-scope conflict
- retry exhaustion

Checkpoint replies are health reports, not permission for the main agent to finish the slice locally.

## Retry Rule

If a worker fails to emit a usable final closeout:

1. retry the same bounded slice once with a fresh worker
2. pass the retry worker the prior failure reason and any relevant partial context
3. if retry still fails, mark the slice blocked or stop the batch wave explicitly

Do not silently recover the implementation in the main agent.

## Slow-Worker Example

If a worker has been quiet and still shows no diff after one wait window:

- do not assume it is stalled yet
- use at most a bounded checkpoint if the slice has gone quiet long enough
- if the checkpoint says `can_continue=true`, keep the same worker and resume waiting for final closeout
