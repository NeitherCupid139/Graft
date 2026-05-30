<template>
  <div data-page-type="list-form-detail">
    <management-page-content>
      <management-page-header :title="t('accessLog.page.title')" :description="t('accessLog.page.description')">
        <template #eyebrow>{{ t('menu.accessLog.title') }}</template>
        <template #actions>
          <t-button theme="default" variant="outline" :loading="loading" @click="fetchAccessLogs">
            {{ t('accessLog.page.refresh') }}
          </t-button>
        </template>
      </management-page-header>

      <access-log-filters v-model="filters" :loading="loading" @reset="resetFilters" @search="handleSearch" />

      <management-empty-state
        v-if="listError && !loading"
        tone="error"
        :title="t('accessLog.page.errorTitle')"
        :description="listError"
      >
        <template #actions>
          <t-button theme="primary" variant="outline" @click="fetchAccessLogs">
            {{ t('accessLog.page.retry') }}
          </t-button>
        </template>
      </management-empty-state>

      <access-log-table
        v-else
        v-model:current="pagination.current"
        v-model:page-size="pagination.pageSize"
        :description="t('accessLog.page.description')"
        :footer-summary="footerSummary"
        :loading="loading"
        :rows="rows"
        :summary="tableSummary"
        :total="total"
        @detail="openDetail"
        @page-change="fetchAccessLogs"
      />
    </management-page-content>

    <access-log-detail-drawer v-model:visible="detailVisible" :record="detailRecord" />
  </div>
</template>
<script setup lang="ts">
import { MessagePlugin } from 'tdesign-vue-next';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { resolveLocalizedErrorMessage } from '@/modules/shared/localized-api-error';
import { ManagementEmptyState, ManagementPageContent, ManagementPageHeader } from '@/shared/components/management';

import { getAccessLogDetail, getAccessLogs } from '../../api/access-log';
import AccessLogDetailDrawer from '../../components/AccessLogDetailDrawer.vue';
import AccessLogFilters from '../../components/AccessLogFilters.vue';
import AccessLogTable from '../../components/AccessLogTable.vue';
import type { AccessLogFilterState, AccessLogItem, AccessLogQuery } from '../../types/access-log';

defineOptions({
  name: 'AccessLogListIndex',
});

const { t } = useI18n();

const loading = ref(false);
const listError = ref('');
const rows = ref<AccessLogItem[]>([]);
const total = ref(0);
const detailVisible = ref(false);
const detailRecord = ref<AccessLogItem | null>(null);
const pagination = ref({
  current: 1,
  pageSize: 20,
});
const filters = ref<AccessLogFilterState>(createDefaultFilters());

const tableSummary = computed(() => t('accessLog.page.summary', { count: rows.value.length }));
const footerSummary = computed(() => t('accessLog.page.footerTotal', { count: total.value }));

function createDefaultFilters(): AccessLogFilterState {
  return {
    requestId: '',
    userId: '',
    username: '',
    method: '',
    path: '',
    pathMatch: 'exact',
    route: '',
    statusCode: '',
    durationMinMs: '',
    durationMaxMs: '',
    occurredRange: [],
    sortBy: 'occurred_at',
    sortOrder: 'desc',
  };
}

function buildQuery(): AccessLogQuery {
  const query: AccessLogQuery = {
    page: pagination.value.current,
    page_size: pagination.value.pageSize,
    sort_by: filters.value.sortBy,
    sort_order: filters.value.sortOrder,
    path_match: filters.value.pathMatch,
  };

  if (filters.value.requestId) query.request_id = filters.value.requestId;
  if (filters.value.userId) query.user_id = Number(filters.value.userId);
  if (filters.value.username) query.username = filters.value.username;
  if (filters.value.method) query.method = filters.value.method;
  if (filters.value.path) query.path = filters.value.path;
  if (filters.value.route) query.route = filters.value.route;
  if (filters.value.statusCode) query.status_code = Number(filters.value.statusCode);
  if (filters.value.durationMinMs) query.duration_min_ms = Number(filters.value.durationMinMs);
  if (filters.value.durationMaxMs) query.duration_max_ms = Number(filters.value.durationMaxMs);
  if (filters.value.occurredRange[0]) query.occurred_from = toISOStringOrRaw(filters.value.occurredRange[0]);
  if (filters.value.occurredRange[1]) query.occurred_to = toISOStringOrRaw(filters.value.occurredRange[1]);

  return query;
}

async function fetchAccessLogs() {
  loading.value = true;
  listError.value = '';

  try {
    const response = await getAccessLogs(buildQuery());
    rows.value = response.items;
    total.value = response.total;
  } catch (error) {
    rows.value = [];
    total.value = 0;
    listError.value = resolveLocalizedErrorMessage(t, error, t('accessLog.page.loadFailed'));
    MessagePlugin.error(listError.value);
  } finally {
    loading.value = false;
  }
}

async function openDetail(row: AccessLogItem) {
  try {
    detailRecord.value = await getAccessLogDetail(Number(row.id));
    detailVisible.value = true;
  } catch (error) {
    MessagePlugin.error(resolveLocalizedErrorMessage(t, error, t('accessLog.page.loadFailed')));
  }
}

function resetFilters() {
  filters.value = createDefaultFilters();
  pagination.value.current = 1;
  void fetchAccessLogs();
}

function handleSearch() {
  pagination.value.current = 1;
  void fetchAccessLogs();
}

function toISOStringOrRaw(value: string) {
  const date = new Date(value.replace(' ', 'T'));
  return Number.isNaN(date.getTime()) ? value : date.toISOString();
}

void fetchAccessLogs();
</script>
