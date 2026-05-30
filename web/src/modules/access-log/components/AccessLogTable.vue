<template>
  <management-table-card>
    <template #head>
      <section class="table-head" aria-label="access-log-table-head">
        <p class="table-head__description">{{ description }}</p>
        <p class="table-head__summary">{{ summary }}</p>
      </section>
    </template>
    <t-table
      row-key="id"
      :columns="columns"
      :data="rows"
      :loading="loading"
      table-layout="fixed"
      :table-content-width="tableContentWidth"
      cell-empty-content="-"
      hover
    >
      <template #method="{ row }">
        <t-tag theme="primary" variant="light-outline" size="small">{{ row.method }}</t-tag>
      </template>
      <template #path="{ row }">
        <div class="stack-cell">
          <strong>{{ row.path }}</strong>
          <span class="stack-cell__secondary">{{ row.route || '-' }}</span>
        </div>
      </template>
      <template #status_code="{ row }">
        <t-tag
          :theme="row.status_code >= 500 ? 'danger' : row.status_code >= 400 ? 'warning' : 'success'"
          variant="light-outline"
          size="small"
        >
          {{ row.status_code }}
        </t-tag>
      </template>
      <template #duration_ms="{ row }">
        <span>{{ row.duration_ms }} ms</span>
      </template>
      <template #user="{ row }">
        <div class="stack-cell">
          <strong>{{ row.username || '-' }}</strong>
          <span class="stack-cell__secondary">{{ row.user_id ?? '-' }}</span>
        </div>
      </template>
      <template #request_id="{ row }">
        <strong class="table-mono">{{ row.request_id }}</strong>
      </template>
      <template #occurred_at="{ row }">
        <span>{{ formatCompactDateTime(row.occurred_at) }}</span>
      </template>
      <template #operation="{ row }">
        <table-action-menu
          :actions="[{ label: t('accessLog.actions.detail'), testId: 'access-log-detail', value: 'detail' }]"
          :more-label="t('accessLog.actions.detail')"
          @action="() => $emit('detail', row)"
        />
      </template>
      <template #empty>
        <div class="table-empty-state">
          <t-empty :title="t('accessLog.page.emptyTitle')" :description="t('accessLog.page.emptyDescription')" />
        </div>
      </template>
    </t-table>
    <template #footer>
      <management-table-pagination :summary="footerSummary">
        <t-pagination
          v-model:current="current"
          v-model:page-size="pageSize"
          :total="total"
          :page-size-options="[10, 20, 50, 100]"
          @change="$emit('page-change')"
        />
      </management-table-pagination>
    </template>
  </management-table-card>
</template>
<script setup lang="ts">
import type { TdBaseTableProps } from 'tdesign-vue-next';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  calculateTableContentWidth,
  createActionColumn,
  createTextColumn,
  createTimeColumn,
  formatCompactDateTime,
  ManagementTableCard,
  ManagementTablePagination,
  TableActionMenu,
} from '@/shared/components/management';

import type { AccessLogItem } from '../types/access-log';

defineProps<{
  description: string;
  footerSummary: string;
  loading?: boolean;
  rows: AccessLogItem[];
  summary: string;
  total: number;
}>();

defineEmits<{
  (e: 'detail', row: AccessLogItem): void;
  (e: 'page-change'): void;
}>();

const current = defineModel<number>('current', { required: true });
const pageSize = defineModel<number>('pageSize', { required: true });

const { t, locale } = useI18n();

const columns = computed<TdBaseTableProps['columns']>(() => {
  void locale.value;
  return [
    createTimeColumn(t('accessLog.columns.occurredAt'), 'occurred_at', 176),
    createTextColumn(t('accessLog.columns.method'), 'method', { width: 110, fixed: 'left' }),
    createTextColumn(t('accessLog.columns.path'), 'path', { minWidth: 280 }),
    createTextColumn(t('accessLog.columns.statusCode'), 'status_code', { width: 110 }),
    createTextColumn(t('accessLog.columns.durationMs'), 'duration_ms', { width: 120 }),
    createTextColumn(t('accessLog.columns.user'), 'user', { width: 180 }),
    createTextColumn(t('accessLog.columns.requestId'), 'request_id', { width: 220 }),
    createActionColumn(t('components.commonTable.operation'), 96),
  ];
});

const tableContentWidth = computed(() => calculateTableContentWidth(columns.value));
</script>
<style scoped lang="less">
.table-head__summary,
.table-head__description,
.stack-cell__secondary {
  color: var(--td-text-color-secondary);
  margin: 0;
}

.stack-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.table-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
}

.table-empty-state {
  padding: 24px 0 8px;
}
</style>
