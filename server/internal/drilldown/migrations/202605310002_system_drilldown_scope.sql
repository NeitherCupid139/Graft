CREATE TABLE IF NOT EXISTS "system_drilldown_scope" (
  "id" BIGSERIAL PRIMARY KEY,
  "module" VARCHAR(64) NOT NULL,
  "scope" VARCHAR(128) NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT NOT NULL DEFAULT '',
  "target_type" VARCHAR(64) NOT NULL,
  "target_module" VARCHAR(64) NOT NULL,
  "target_page" VARCHAR(128) NOT NULL,
  "filter_payload" JSONB NULL,
  "enabled" BOOLEAN NOT NULL DEFAULT TRUE,
  "sort_order" INTEGER NOT NULL DEFAULT 0,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "uq_system_drilldown_scope_module_scope"
  ON "system_drilldown_scope" ("module", "scope");

CREATE INDEX IF NOT EXISTS "idx_system_drilldown_scope_module_enabled"
  ON "system_drilldown_scope" ("module", "enabled");

CREATE INDEX IF NOT EXISTS "idx_system_drilldown_scope_target"
  ON "system_drilldown_scope" ("target_module", "target_page");

INSERT INTO "system_drilldown_scope" (
  "module",
  "scope",
  "name",
  "description",
  "target_type",
  "target_module",
  "target_page",
  "filter_payload",
  "enabled",
  "sort_order"
) VALUES (
  'audit',
  'sensitive_operations',
  '敏感操作',
  '用于安全审计概览与日志页之间的业务钻取语义。',
  'log_query',
  'audit',
  'audit_logs',
  NULL,
  TRUE,
  100
)
ON CONFLICT ("module", "scope") DO UPDATE SET
  "name" = EXCLUDED."name",
  "description" = EXCLUDED."description",
  "target_type" = EXCLUDED."target_type",
  "target_module" = EXCLUDED."target_module",
  "target_page" = EXCLUDED."target_page",
  "enabled" = EXCLUDED."enabled",
  "sort_order" = EXCLUDED."sort_order",
  "updated_at" = NOW();

COMMENT ON TABLE "system_drilldown_scope" IS 'Platform-owned drilldown scope metadata registry. Query semantics stay in module registries.';
COMMENT ON COLUMN "system_drilldown_scope"."filter_payload" IS 'Reserved future extension payload. It is not query authority in v1.';
  