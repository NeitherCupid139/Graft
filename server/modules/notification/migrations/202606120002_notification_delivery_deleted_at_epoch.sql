-- Copyright (c) 2025-2026 GeWuYou
-- SPDX-License-Identifier: Apache-2.0

DROP INDEX IF EXISTS "notification_deliveries_recipient_created";
DROP INDEX IF EXISTS "notification_deliveries_unread";

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = current_schema()
      AND table_name = 'notification_deliveries'
      AND column_name = 'deleted_at'
      AND data_type <> 'bigint'
  ) THEN
    ALTER TABLE "notification_deliveries"
      ALTER COLUMN "deleted_at" DROP DEFAULT,
      ALTER COLUMN "deleted_at" TYPE bigint
        USING CASE
          WHEN "deleted_at" IS NULL THEN 0
          ELSE EXTRACT(EPOCH FROM "deleted_at")::bigint
        END,
      ALTER COLUMN "deleted_at" SET DEFAULT 0,
      ALTER COLUMN "deleted_at" SET NOT NULL;
  ELSE
    UPDATE "notification_deliveries"
    SET "deleted_at" = 0
    WHERE "deleted_at" IS NULL;

    ALTER TABLE "notification_deliveries"
      ALTER COLUMN "deleted_at" SET DEFAULT 0,
      ALTER COLUMN "deleted_at" SET NOT NULL;
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS "notification_deliveries_recipient_created"
  ON "notification_deliveries" ("recipient_user_id", "created_at" DESC, "id" DESC)
  WHERE "deleted_at" = 0;

CREATE INDEX IF NOT EXISTS "notification_deliveries_unread"
  ON "notification_deliveries" ("recipient_user_id", "id" DESC)
  WHERE "read_at" IS NULL AND "deleted_at" = 0;

COMMENT ON TABLE "notification_deliveries" IS '通知中心用户投递状态表';
COMMENT ON COLUMN "notification_deliveries"."deleted_at" IS '软删除时间戳，0 表示未删除，非 0 表示删除发生时的 Unix 秒';
COMMENT ON INDEX "notification_deliveries_recipient_created" IS '当前用户可见通知列表排序索引，仅覆盖未删除投递';
COMMENT ON INDEX "notification_deliveries_unread" IS '当前用户未读通知计数索引，仅覆盖未删除投递';
