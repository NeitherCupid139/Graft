---
name: graft-frontend-page-builder
description: Build or reshape Graft admin frontend pages using Vue 3, TypeScript, TDesign Vue Next, Pinia, Vue Router, Axios, UnoCSS, Graft page-type governance, and TDesign MCP. Use instead of generic Frontend App Builder, Frontend Design, Impeccable UI, Theme Factory, landing-page, or React/Tailwind/shadcn builder skills for Graft web work.
---

# Graft Frontend Page Builder

Use this skill for real Graft admin UI work, not standalone websites. The output must fit `web/src/modules/<name>` or shell-owned frontend boundaries and preserve existing route, menu, permission, i18n, and theme conventions.

## Workflow

1. Complete repository startup preflight and read `web/AGENTS.md`.
2. Classify the page type through existing Graft frontend governance, especially `$graft-web-vibe-coding` when page, shell, visual, copy, or prompt shaping is involved.
3. Use TDesign Vue Next components through `$graft-web-vibe-coding` and the TDesign MCP docs when component API, DOM, or changelog detail is needed.
4. Keep implementation inside the existing Graft module structure. Do not create a second app, framework baseline, global style system, or local design system.
5. Build dense admin surfaces: clear tables, forms, filters, detail panels, drawers, dialogs, status indicators, and action bars. Avoid marketing hero pages, decorative card-heavy layouts, and oversized display copy.
6. Keep visible copy i18n-safe and aligned with existing locale patterns.
7. Validate with the frontend entrypoint required by repository governance, normally `bun run check`, plus focused checks when appropriate.

## Replacement Map

- Frontend App Builder -> Graft module/page implementation under `web/src/modules/**`.
- Frontend Design / Impeccable -> Graft page-type workflow, TDesign components, repository tokens, responsive constraints.
- Theme Factory -> existing Graft theme tokens and TDesign theming only.
- Tailwind/shadcn/React starter -> reject unless repository docs are explicitly revised first.

## Constraints

- Do not introduce React, shadcn, Tailwind as a runtime baseline, or web/package dependency changes.
- Do not add visible in-app instructions about how to use the UI unless product copy already requires them.
- Do not use generic generated templates that bypass Graft module, API, contract, permission, or route ownership.
- Escalate to cross-boundary governance when the frontend symptom requires server descriptors, OpenAPI, typed contract, menu, route, or permission authority repair.
