-- Copyright (c) 2025-2026 GeWuYou
-- SPDX-License-Identifier: Apache-2.0

ALTER TABLE "notification_events"
  ADD COLUMN IF NOT EXISTS "category_key" character varying NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS "source_key" character varying NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS "level_key" character varying NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS "event_type_key" character varying NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS "action_label_key" character varying NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS "action_label" text NOT NULL DEFAULT '';

COMMENT ON COLUMN "notification_events"."category_key" IS '通知分类展示国际化键';
COMMENT ON COLUMN "notification_events"."source_key" IS '通知来源展示国际化键';
COMMENT ON COLUMN "notification_events"."level_key" IS '通知级别展示国际化键';
COMMENT ON COLUMN "notification_events"."event_type_key" IS '通知事件类型展示国际化键';
COMMENT ON COLUMN "notification_events"."action_label_key" IS '通知业务动作文案国际化键';
COMMENT ON COLUMN "notification_events"."action_label" IS '通知业务动作文案回退快照';
