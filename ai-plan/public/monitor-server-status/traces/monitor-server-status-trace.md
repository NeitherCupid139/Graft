# Monitor Server Status Trace

## 2026-05-20

- Completed the first minimal `monitor/server-status` slice across `server/plugins/monitor/**` and `web/src/modules/monitor/**`.
- Registered the backend plugin through the approved shared-hotspot exception in `server/internal/pluginregistry/generated.go`.
- Completed the runtime metadata snapshot follow-up and validated that slice with backend checks.
- Completed the richer dashboard follow-up inside the owned cross-boundary scope:
  - upgraded the `monitor` plugin response with runtime, summary, dependency-detail, plugin dependency, and in-memory trend data
  - finished the inherited dashboard page diff and aligned monitor module types, locales, and Vitest coverage
  - updated the monitor topic design/tracking docs and `web/AGENTS.md` to freeze theme-token and chart responsiveness rules for monitor-style dashboards
- Completed the IA-alignment follow-up for the same real page:
  - kept `/monitor/server-status` as the only real runtime page while registering a real backend `服务器管理` menu parent and assembling it into the shell route tree
  - upgraded the overview page with 5-second default auto refresh, visibility pause/resume, retry backoff, icon-assisted summary cards, grouped runtime sections, and a non-empty trend fallback
  - aligned locale catalogs and tests so breadcrumb/menu semantics render `服务器管理 / 服务器状态` without exposing an `index` crumb
  - updated monitor topic design/tracking docs to record future IA placeholders as design-only, not runtime contracts
- Full command- and file-level history for this stage stays in the session log; keep this trace as the concise recovery entrypoint.
