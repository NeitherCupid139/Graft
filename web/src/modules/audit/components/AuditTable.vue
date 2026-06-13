<!--
  Copyright (c) 2025-2026 GeWuYou
  SPDX-License-Identifier: Apache-2.0
-->

<template>
  <management-table-card>
    <template #head>
      <div class="table-head">
        <div>
          <p class="table-head__summary">{{ summary }}</p>
          <p class="table-head__description">{{ description }}</p>
        </div>
        <t-tag v-if="localFilterActive" theme="default" variant="light-outline" size="small">
          {{ t('audit.logList.currentPageFiltered') }}
        </t-tag>
      </div>
    </template>
    <template v-if="$slots.toolbar" #toolbar>
      <slot name="toolbar" />
    </template>

    <div ref="tableHostRef" class="audit-log-table__host" :data-table-mode="tableWidthPolicy.mode">
      <t-table
        row-key="id"
        :columns="columns"
        :data="rows"
        :loading="loading"
        table-layout="fixed"
        :table-content-width="tableWidthPolicy.tableContentWidth"
        cell-empty-content="-"
        hover
        @row-click="handleRowClick"
      >
        <template #action="{ row }">
          <div class="stack-cell">
            <strong>{{ actionTitle(row, t) }}</strong>
            <span class="stack-cell__secondary">{{ actionCategoryLabel(row, t) }}</span>
          </div>
        </template>

        <template #actor="{ row }">
          <div class="stack-cell">
            <strong>{{ actorLabel(row, t) }}</strong>
            <span class="stack-cell__secondary">{{ actorSecondaryLabel(row) }}</span>
          </div>
        </template>

        <template #resource="{ row }">
          <div class="stack-cell">
            <strong>{{ resourceLabel(row, t) }}</strong>
            <span class="stack-cell__secondary">{{ reasonForRecord(row, t) }}</span>
          </div>
        </template>

        <template #correlation="{ row }">
          <log-id-text
            :display-value="requestIdForRecord(row)"
            :tooltip="requestIdForRecord(row)"
            v-bind="technicalCopyLabels"
          />
        </template>

        <template #trace_id="{ row }">
          <log-id-text
            :display-value="traceIdForRecord(row)"
            :tooltip="traceIdForRecord(row)"
            v-bind="technicalCopyLabels"
          />
        </template>

        <template #session_id="{ row }">
          <log-id-text
            :display-value="row.session_id || '-'"
            :tooltip="row.session_id || '-'"
            v-bind="technicalCopyLabels"
          />
        </template>

        <template #ip="{ row }">
          <log-id-text :display-value="row.ip || '-'" :tooltip="row.ip || '-'" v-bind="technicalCopyLabels" />
        </template>

        <template #result="{ row }">
          <t-tag :theme="resultTone(row)" variant="light-outline" size="small">
            {{ resultLabel(row, t) }}
          </t-tag>
        </template>

        <template #risk="{ row }">
          <t-tag :theme="riskTone(row)" variant="light-outline" size="small">
            {{ riskLabel(row, t) }}
          </t-tag>
        </template>

        <template #created_at="{ row }">
          <span>{{ formatAuditTimestamp(row.created_at, locale) }}</span>
        </template>

        <template #empty>
          <div class="table-empty-state">
            <t-empty :title="t('audit.logList.emptyTitle')" :description="t('audit.logList.emptyDescription')" />
          </div>
        </template>
      </t-table>
    </div>

    <template #footer>
      <management-table-pagination :summary="footerSummary">
        <t-pagination
          v-model:current="current"
          v-model:page-size="pageSize"
          :total="total"
          :page-size-options="[10, 20, 50, 100]"
          @change="emitPageChange"
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
  createIdentifierColumn,
  createMainTextColumn,
  createStatusColumn,
  createTechnicalColumn,
  createTimeColumn,
  ManagementTableCard,
  ManagementTablePagination,
  resolveManagedColumns,
  resolveTableWidthPolicy,
  useTableHostWidth,
} from '@/shared/components/management';
import { LogIdText } from '@/shared/observability';

import {
  actionCategoryLabel,
  actionTitle,
  actorLabel,
  actorSecondaryLabel,
  formatAuditTimestamp,
  reasonForRecord,
  requestIdForRecord,
  resourceLabel,
  resultLabel,
  resultTone,
  riskLabel,
  riskTone,
  traceIdForRecord,
} from '../shared/presentation';
import type { AuditLogListItem } from '../types/audit';

const props = defineProps<{
  description?: string;
  footerSummary: string;
  loading?: boolean;
  localFilterActive?: boolean;
  rows: AuditLogListItem[];
  summary?: string;
  total: number;
  visibleColumnKeys?: string[];
}>();

const emit = defineEmits<{
  (e: 'detail', row: AuditLogListItem): void;
  (e: 'update:current', value: number): void;
  (e: 'update:pageSize', value: number): void;
  (e: 'page-change'): void;
}>();

const { t, locale } = useI18n();
const technicalCopyLabels = computed(() => ({
  copyable: true,
  copyLabel: t('audit.logList.drawer.actions.copyRequestId'),
  copySuccessLabel: t('audit.logList.drawer.actions.copyRequestIdSuccess'),
  copyFailLabel: t('audit.logList.drawer.actions.copyRequestIdFail'),
}));

const current = defineModel<number>('current', { required: true });
const pageSize = defineModel<number>('pageSize', { required: true });

const columns = computed<TdBaseTableProps['columns']>(() => {
  void locale.value;

  const allColumns: TdBaseTableProps['columns'] = [
    createMainTextColumn(t('audit.logList.columns.action'), 'action', 260),
    createIdentifierColumn(t('audit.logList.columns.actor'), 'actor', 168),
    createIdentifierColumn(t('audit.logList.columns.resource'), 'resource', 208),
    createTechnicalColumn(t('audit.logList.columns.correlation'), 'correlation', 248),
    createTechnicalColumn(t('audit.logList.columns.traceId'), 'trace_id', 260),
    createTechnicalColumn(t('audit.logList.columns.sessionId'), 'session_id', 220),
    createIdentifierColumn(t('audit.logList.columns.ip'), 'ip', 160),
    createStatusColumn(t('audit.logList.columns.result'), 'result', 132),
    createStatusColumn(t('audit.logList.columns.risk'), 'risk', 120),
    createTimeColumn(t('audit.logList.columns.createdAt'), 'created_at', 200),
  ];

  return resolveManagedColumns(allColumns, props.visibleColumnKeys);
});

const { tableHostRef, tableHostWidth } = useTableHostWidth(() => columns.value);

const tableWidthPolicy = computed(() => resolveTableWidthPolicy(columns.value, tableHostWidth.value));

function emitPageChange() {
  emit('page-change');
}

function handleRowClick(context: { row: unknown }) {
  emit('detail', context.row as AuditLogListItem);
}
</script>
<style scoped lang="less">
@import '@/shared/observability/log-table-cells.less';

.table-head {
  align-items: center;
  display: flex;
  gap: var(--graft-density-gap-12);
  justify-content: space-between;
}

.table-head__summary,
.table-head__description {
  color: var(--td-text-color-secondary);
  margin: 0;
}

.log-table-stack-cells();

.table-empty-state {
  padding: var(--graft-density-gap-24) 0 var(--graft-density-gap-8);
}

.audit-log-table__host {
  max-width: 100%;
  min-width: 0;
  overflow-x: hidden;
}

.audit-log-table__host[data-table-mode='scroll'] {
  overflow-x: auto;
}

@media (width <= 768px) {
  .table-head {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
