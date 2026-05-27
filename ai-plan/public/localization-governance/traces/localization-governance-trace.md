# Localization Governance Trace

## 2026-05-27 active topic initialized

- Created a new active public recovery topic `localization-governance`.
- Renamed the implementation workspace from the archived-topic-derived identity to:
  - worktree `/home/gewuyou/project/go/Graft-wt/feat/wt-localization-governance`
  - branch `feat/wt-localization-governance`
- Updated `ai-plan/public/README.md` so the recovery index no longer points to `None` for this active task.
- Recorded the frozen localization governance compatibility rules before any business-code implementation:
  - locale bundles must be future-compatible with plugin-provided sources
  - all locale keys require owner namespaces
  - menu, permission display, and error semantics use stable keys
  - frontend consumes `key + fallback`
  - backend registry remains registration/validation/fallback only
  - OpenAPI stays key-semantic only and does not become a multilingual copy store
