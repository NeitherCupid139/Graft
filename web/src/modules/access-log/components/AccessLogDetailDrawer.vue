<template>
  <t-drawer
    :visible="visible"
    :header="t('accessLog.page.detailTitle')"
    :footer="false"
    destroy-on-close
    placement="right"
    size="640px"
    @update:visible="$emit('update:visible', $event)"
  >
    <div v-if="record" class="access-log-detail">
      <section class="access-log-detail__section">
        <h4>{{ t('accessLog.detail.basic') }}</h4>
        <div class="access-log-detail__grid">
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.columns.occurredAt') }}</span
            ><strong>{{ record.occurred_at }}</strong>
          </div>
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.columns.method') }}</span
            ><strong>{{ record.method }}</strong>
          </div>
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.columns.statusCode') }}</span
            ><strong>{{ record.status_code }}</strong>
          </div>
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.columns.durationMs') }}</span
            ><strong>{{ record.duration_ms }} ms</strong>
          </div>
          <div class="access-log-detail__item access-log-detail__item--full">
            <span>{{ t('accessLog.columns.path') }}</span
            ><strong>{{ record.path }}</strong>
          </div>
          <div class="access-log-detail__item access-log-detail__item--full">
            <span>{{ t('accessLog.detail.route') }}</span
            ><strong>{{ record.route || '-' }}</strong>
          </div>
          <div class="access-log-detail__item access-log-detail__item--full">
            <span>{{ t('accessLog.detail.requestId') }}</span>
            <div class="access-log-detail__copy-line">
              <strong class="access-log-detail__mono">{{ record.request_id }}</strong>
              <t-button size="small" theme="default" variant="text" @click="copyValue(record.request_id)">{{
                t('accessLog.actions.copy')
              }}</t-button>
            </div>
          </div>
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.detail.userId') }}</span
            ><strong>{{ record.user_id ?? '-' }}</strong>
          </div>
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.detail.username') }}</span
            ><strong>{{ record.username || '-' }}</strong>
          </div>
        </div>
      </section>

      <section class="access-log-detail__section">
        <h4>{{ t('accessLog.detail.network') }}</h4>
        <div class="access-log-detail__grid">
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.detail.clientIp') }}</span
            ><strong>{{ record.client_ip || '-' }}</strong>
          </div>
          <div class="access-log-detail__item access-log-detail__item--full">
            <span>{{ t('accessLog.detail.userAgent') }}</span
            ><strong>{{ record.user_agent || '-' }}</strong>
          </div>
        </div>
      </section>

      <section class="access-log-detail__section">
        <h4>{{ t('accessLog.detail.size') }}</h4>
        <div class="access-log-detail__grid">
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.detail.requestSize') }}</span
            ><strong>{{ record.request_size ?? '-' }}</strong>
          </div>
          <div class="access-log-detail__item">
            <span>{{ t('accessLog.detail.responseSize') }}</span
            ><strong>{{ record.response_size ?? '-' }}</strong>
          </div>
        </div>
      </section>
    </div>
  </t-drawer>
</template>
<script setup lang="ts">
import { MessagePlugin } from 'tdesign-vue-next';
import { useI18n } from 'vue-i18n';

import type { AccessLogItem } from '../types/access-log';

defineProps<{
  record: AccessLogItem | null;
  visible: boolean;
}>();

defineEmits<{
  (e: 'update:visible', value: boolean): void;
}>();

const { t } = useI18n();

async function copyValue(value: string) {
  try {
    await navigator.clipboard.writeText(value);
    MessagePlugin.success(t('accessLog.actions.copySuccess'));
  } catch {
    MessagePlugin.error(t('accessLog.actions.copyFail'));
  }
}
</script>
<style scoped lang="less">
.access-log-detail {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.access-log-detail__section h4 {
  margin: 0 0 12px;
}

.access-log-detail__grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.access-log-detail__item {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-border);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px;
}

.access-log-detail__item--full {
  grid-column: 1 / -1;
}

.access-log-detail__copy-line {
  align-items: center;
  display: flex;
  gap: 8px;
  justify-content: space-between;
}

.access-log-detail__mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

@media (width <= 768px) {
  .access-log-detail__grid {
    grid-template-columns: 1fr;
  }
}
</style>
