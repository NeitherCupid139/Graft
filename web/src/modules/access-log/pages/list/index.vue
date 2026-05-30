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
        :empty-description="emptyDescription"
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
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import { resolveLocalizedErrorMessage as resolveAccessLogErrorMessage } from '@/modules/shared/localized-api-error';
import { ManagementEmptyState, ManagementPageContent, ManagementPageHeader } from '@/shared/components/management';
import { createLogger as createModuleLogger } from '@/utils/logger';

import { getAccessLogDetail, getAccessLogs } from '../../api/access-log';
import AccessLogDetailDrawer from '../../components/AccessLogDetailDrawer.vue';
import AccessLogFilters from '../../components/AccessLogFilters.vue';
import AccessLogTable from '../../components/AccessLogTable.vue';
import {
  buildAccessLogLocation,
  buildAccessLogRequestLocation,
  buildAccessLogTraceLocation,
  parseAccessLogRouteQuery,
} from '../../contract/deep-link';
import type { AccessLogFilterState, AccessLogItem, AccessLogQuery } from '../../types/access-log';

defineOptions({
  name: 'AccessLogListIndex',
});

const { t } = useI18n();
const logger = createModuleLogger('access-log.list');
const route = useRoute();
const router = useRouter();

const loading = ref(false);
const listError = ref('');
const rows = ref<AccessLogItem[]>([]);
const total = ref(0);
const detailVisible = ref(false);
const detailRecord = ref<AccessLogItem | null>(null);
const applyingRoute = ref(false);
const pagination = ref({
  current: 1,
  pageSize: 20,
});
const filters = ref<AccessLogFilterState>(createDefaultFilters());
const deepLinkCorrelation = ref<'requestId' | 'traceId' | null>(null);

const tableSummary = computed(() => t('accessLog.page.summary', { count: rows.value.length }));
const footerSummary = computed(() => t('accessLog.page.footerTotal', { count: total.value }));
const emptyDescription = computed(() => {
  if (deepLinkCorrelation.value === 'requestId') {
    return t('accessLog.page.emptyRequestDescription');
  }
  if (deepLinkCorrelation.value === 'traceId') {
    return t('accessLog.page.emptyTraceDescription');
  }
  return t('accessLog.page.emptyDescription');
});
const routeLocation = computed(() => {
  if (filters.value.requestId) {
    return buildAccessLogRequestLocation(filters.value.requestId);
  }

  if (filters.value.traceId) {
    return buildAccessLogTraceLocation(filters.value.traceId);
  }

  return buildAccessLogLocation({});
});

function createDefaultFilters(): AccessLogFilterState {
  return {
    requestId: '',
    traceId: '',
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
  if (filters.value.traceId) query.trace_id = filters.value.traceId;
  if (filters.value.userId) query.user_id = Number(filters.value.userId);
  if (filters.value.username) query.username = filters.value.username;
  if (filters.value.method) query.method = filters.value.method;
  if (filters.value.path) query.path = filters.value.path;
  if (filters.value.route) query.route = filters.value.route;
  if (filters.value.statusCode) query.status_code = Number(filters.value.statusCode);
  if (filters.value.durationMinMs) query.duration_min_ms = Number(filters.value.durationMinMs);
  if (filters.value.durationMaxMs) query.duration_max_ms = Number(filters.value.durationMaxMs);
  if (filters.value.occurredRange[0]) query.occurred_from = normalizeOccurredAt(filters.value.occurredRange[0]);
  if (filters.value.occurredRange[1]) query.occurred_to = normalizeOccurredAt(filters.value.occurredRange[1]);

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
    logger.error('failed to fetch access logs', error);
    listError.value = resolveAccessLogErrorMessage(t, error, t('accessLog.page.loadFailed'));
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
    MessagePlugin.error(resolveAccessLogErrorMessage(t, error, t('accessLog.page.loadFailed')));
  }
}

function resetFilters() {
  filters.value = createDefaultFilters();
  pagination.value.current = 1;
  void updateRouteQuery();
}

function handleSearch() {
  pagination.value.current = 1;
  void updateRouteQuery();
}

function normalizeOccurredAt(value: string) {
  const date = new Date(value.replace(' ', 'T'));
  return Number.isFinite(date.getTime()) ? date.toISOString() : value;
}

function applyRouteFilters() {
  const { request_id: requestId = '', trace_id: traceId = '' } = parseAccessLogRouteQuery(route.query);
  filters.value = withCorrelationFilters(filters.value, requestId, traceId);
  deepLinkCorrelation.value = detectCorrelationMode(requestId, traceId);
}

function isCurrentRouteQuery(targetQuery: Record<string, string>) {
  const currentQuery = buildAccessLogLocation(route.query).query as Record<string, string>;
  const targetEntries = Object.entries(targetQuery);

  return (
    Object.keys(currentQuery).length === targetEntries.length &&
    targetEntries.every(([key, value]) => currentQuery[key] === value)
  );
}

function withCorrelationFilters(
  baseFilters: AccessLogFilterState,
  requestId: string,
  traceId: string,
): AccessLogFilterState {
  return {
    ...baseFilters,
    requestId,
    traceId,
  };
}

function detectCorrelationMode(requestId: string, traceId: string): 'requestId' | 'traceId' | null {
  if (requestId) {
    return 'requestId';
  }

  return traceId ? 'traceId' : null;
}

async function updateRouteQuery() {
  if (applyingRoute.value) {
    return;
  }

  if (isCurrentRouteQuery(routeLocation.value.query as Record<string, string>)) {
    await fetchAccessLogs();
    return;
  }

  await router.replace(routeLocation.value);
}

watch(
  () => [route.query.request_id, route.query.trace_id],
  ([requestId, traceId]) => {
    applyingRoute.value = true;
    try {
      const hasQuery = Boolean(requestId) || Boolean(traceId);
      if (!hasQuery && !filters.value.requestId && !filters.value.traceId) {
        deepLinkCorrelation.value = null;
      } else {
        applyRouteFilters();
      }
    } finally {
      applyingRoute.value = false;
    }
    pagination.value.current = 1;
    void fetchAccessLogs();
  },
  { immediate: true },
);
</script>
