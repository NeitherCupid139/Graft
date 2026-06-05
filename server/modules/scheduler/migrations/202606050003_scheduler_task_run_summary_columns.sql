ALTER TABLE "scheduler_task_runs"
  ADD COLUMN IF NOT EXISTS "result_summary" text NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS "error_message" text NOT NULL DEFAULT '';

COMMENT ON COLUMN "scheduler_task_runs"."result_summary" IS '运行结果摘要，不保存完整响应体';
COMMENT ON COLUMN "scheduler_task_runs"."error_message" IS '失败错误摘要，为空表示无错误';
