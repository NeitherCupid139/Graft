# System Config Model Unification Tracking

## Current Goal

Convert the System Config exploration and planning discussion into durable repository design truth.

## Current Scope

- [x] Create repository-level System Config model and renderer design document.
- [x] Register this active topic in `ai-plan/public/README.md`.
- [x] Add topic recovery entrypoint, tracking file, and trace file.
- [x] Phase 1 implementation slice: UI consistency without backend model changes.
- [x] Phase 2 implementation slice: front-end renderer extraction.
- [x] Phase 3 backend registry/OpenAPI enhancement.
- [ ] Later implementation slice: Phase 4 typical config migration.
- [ ] Later implementation slice: Phase 5 validation and screenshot acceptance.

## Authority Discovery

- Design truth owner: `ai-plan/design/系统配置模型与渲染设计.md`.
- Runtime definition authority: module-owned `ConfigDefinition` and `config_schema`.
- Override persistence authority: `server/modules/system-config`.
- Wire contract authority: OpenAPI source under `openapi/**`.
- UI consumer authority: `web/src/modules/system-config` and shared renderer code under `web/src/shared/schema-form`.

## Current Risks

- Do not let derived `fields` views become a second schema truth.
- Do not migrate old flat keys without explicit compatibility expiry and cleanup conditions.
- Do not plan nested object or array visual builders before backend schema validation supports them.
- Do not let frontend key-specific mappings replace module-owned schema repair.

## Latest Validation

- Phase 3 backend registry/OpenAPI enhancement slice.
- Ran validation:
  - `cd server && go test ./internal/scheduler ./internal/configregistry ./modules/system-config`
  - `cd server && go run ./cmd/graft validate backend --stage lint`
- OpenAPI source and generated artifacts were not changed because Phase 3 did not add the optional derived `fields`
  response view.

## Next Step

Start Phase 4 typical config migration as the next separate slice after Phase 3 backend validation and
outer-orchestrator commit acceptance.
