<!--
  Copyright (c) 2025-2026 GeWuYou
  SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="container-page" data-page-type="list-form-detail">
    <management-page-header
      title-key="container.list.title"
      :title="t('container.list.title')"
      description-key="container.list.description"
      :description="t('container.list.description')"
      :source="{ labelKey: 'container.list.eyebrow', fallback: t('container.list.eyebrow') }"
    >
      <template #meta>
        <t-space break-line size="small">
          <t-tag :theme="runtimeStatusTheme" variant="light-outline">
            {{ t('container.list.runtimeLabel') }}: {{ runtimeSummary }}
          </t-tag>
          <t-tag theme="default" variant="light-outline">
            {{ t('container.list.runtimeContainers', runtimeCountParams) }}
          </t-tag>
          <t-tag theme="success" variant="light-outline">
            {{ t('container.list.runningCount', { count: runningCount }) }}
          </t-tag>
          <t-tag theme="warning" variant="light-outline">
            {{ t('container.list.stoppedCount', { count: stoppedCount }) }}
          </t-tag>
          <t-tag theme="danger" variant="light-outline">
            {{ t('container.list.errorCount', { count: errorCount }) }}
          </t-tag>
        </t-space>
      </template>
    </management-page-header>

    <management-toolbar>
      <template #filters>
        <t-input
          v-model="filters.keyword"
          class="management-list-search"
          clearable
          :placeholder="t('container.list.filters.searchPlaceholder')"
          @enter="applyFilters"
        >
          <template #prefix-icon><search-icon /></template>
        </t-input>
        <t-select
          v-model="filters.status"
          class="management-toolbar__select"
          :placeholder="t('container.list.filters.status')"
        >
          <t-option value="all" :label="t('container.list.filters.allStatuses')" />
          <t-option v-for="status in statusOptions" :key="status" :value="status" :label="stateLabel(status)" />
        </t-select>
        <t-button theme="primary" @click="applyFilters">
          {{ t('container.list.filters.query') }}
        </t-button>
        <t-button theme="default" variant="text" @click="resetFilters">
          {{ t('container.list.filters.reset') }}
        </t-button>
      </template>
    </management-toolbar>

    <management-table-card>
      <template #head>
        <div class="container-table-head">
          <p class="container-table-head__summary">
            {{ t('container.list.tableSummary', { count: filteredRows.length }) }}
          </p>
          <p>{{ t('container.list.tableHint') }}</p>
        </div>
      </template>
      <template #toolbar>
        <table-view-toolbar
          :column-settings-label="t('container.list.columnSettings')"
          :density-label="tableDensityLabel"
          :refresh-label="t('container.list.refresh')"
          :refresh-loading="loading"
          @column-settings="columnDrawerVisible = true"
          @density="toggleTableDensity"
          @refresh="refreshContainers"
        />
      </template>

      <t-alert v-if="listError.title" class="container-alert" theme="error" :title="listError.title">
        <p v-if="listError.hint" class="container-alert__hint">{{ listError.hint }}</p>
        <template #operation>
          <t-button theme="danger" variant="text" @click="refreshContainers">
            {{ t('container.list.retry') }}
          </t-button>
        </template>
      </t-alert>

      <div ref="tableHostRef" class="container-table-host" :data-table-mode="tableWidthPolicy.mode">
        <t-table
          row-key="id"
          :columns="visibleColumns"
          :data="paginatedRows"
          :loading="loading"
          :size="tableDensity"
          :table-content-width="tableWidthPolicy.tableContentWidth"
          cell-empty-content="-"
          table-layout="fixed"
          hover
        >
          <template #state="{ row }">
            <t-tag :theme="stateTheme(row.state)" variant="light-outline">
              {{ stateLabel(row.state) }}
            </t-tag>
          </template>

          <template #name="{ row }">
            <div class="container-identity">
              <span class="container-identity__name">{{ displayName(row) }}</span>
              <t-tooltip :content="row.id" placement="top-left">
                <span class="container-identity__id">{{ shortContainerId(row.id) }}</span>
              </t-tooltip>
            </div>
          </template>

          <template #image="{ row }">
            <div class="container-image">
              <span>{{ row.image }}</span>
              <span v-if="row.runtime" class="container-muted">{{ row.runtime }}</span>
            </div>
          </template>

          <template #ports="{ row }">
            <div v-if="visiblePortLabels(row).length" class="container-port-list">
              <t-tag v-for="port in visiblePortLabels(row)" :key="port" size="small" theme="default" variant="light">
                {{ port }}
              </t-tag>
              <t-tooltip
                v-if="hiddenPortLabels(row).length"
                :content="hiddenPortLabels(row).join(' / ')"
                placement="top"
              >
                <t-tag size="small" theme="primary" variant="light">
                  {{ t('container.list.morePorts', { count: hiddenPortLabels(row).length }) }}
                </t-tag>
              </t-tooltip>
            </div>
            <span v-else>-</span>
          </template>

          <template #runtime_status="{ row }">
            <div class="container-runtime-status">
              <span>{{ row.runtime || '-' }}</span>
              <span>{{ row.status || '-' }}</span>
            </div>
          </template>

          <template #image_id="{ row }">
            <span>{{ row.image_id || '-' }}</span>
          </template>

          <template #labels="{ row }">
            <span>{{ labelSummary(row) }}</span>
          </template>

          <template #created_at="{ row }">
            {{ formatTime(row.created_at) }}
          </template>

          <template #started_at="{ row }">
            {{ formatTime(row.started_at) }}
          </template>

          <template #restart_policy="{ row }">
            {{ row.restart_policy || '-' }}
          </template>

          <template #operation="{ row }">
            <t-space class="container-actions" size="small" align="center">
              <t-button
                v-permission="permissionCodes.DETAIL"
                data-testid="container-action-detail"
                theme="primary"
                variant="text"
                size="small"
                @click="openDetail(row)"
              >
                {{ t('container.list.actions.detail') }}
              </t-button>
              <table-action-menu
                :actions="rowActions(row)"
                :more-label="t('container.list.actions.more')"
                more-label-fallback="container.list.actions.more"
                @action="(action) => handleRowAction(action, row)"
              />
            </t-space>
          </template>

          <template #empty>
            <t-empty
              :title="t('container.list.emptyTitle')"
              :description="
                hasActiveFilters ? t('container.list.emptyFilteredDescription') : t('container.list.emptyDescription')
              "
            >
              <template v-if="hasActiveFilters" #action>
                <t-button theme="primary" variant="outline" @click="resetFilters">
                  {{ t('container.list.clearFilters') }}
                </t-button>
              </template>
            </t-empty>
          </template>
        </t-table>
      </div>

      <template #footer>
        <management-table-pagination :summary="footerSummary">
          <t-pagination
            v-model:current="pagination.current"
            v-model:page-size="pagination.pageSize"
            :page-size-options="[10, 20, 50, 100]"
            :show-page-number="true"
            :total="filteredRows.length"
            @change="handlePageChange"
          />
        </management-table-pagination>
      </template>
    </management-table-card>

    <advanced-query-column-drawer
      v-model:visible="columnDrawerVisible"
      v-model:selected-keys="visibleColumnKeys"
      :columns="columnSettingOptions"
      :default-selected-keys="DEFAULT_VISIBLE_COLUMNS"
      :disabled-keys="ALWAYS_VISIBLE_COLUMNS"
      :reset-label="t('container.list.resetColumns')"
      :title="t('container.list.columnSettings')"
    />

    <t-drawer
      v-model:visible="detailDrawerVisible"
      :header="t('container.list.detail.title')"
      :footer="false"
      size="560px"
    >
      <div class="container-drawer-panel">
        <t-alert v-if="detailError" theme="error" :title="detailError">
          <template #operation>
            <t-button v-if="selectedContainer" theme="danger" variant="text" @click="loadDetail(selectedContainer.id)">
              {{ t('container.list.retry') }}
            </t-button>
          </template>
        </t-alert>
        <t-loading :loading="detailLoading">
          <section v-if="activeDetail" class="container-detail-stack">
            <t-descriptions :title="t('container.list.detail.identity')" :column="1" bordered>
              <t-descriptions-item :label="t('container.list.fields.name')">{{
                displayName(activeDetail)
              }}</t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.id')">{{ activeDetail.id }}</t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.image')">{{
                activeDetail.image
              }}</t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.state')">
                {{ stateLabel(activeDetail.state) }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.status')">{{
                activeDetail.status
              }}</t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.restartPolicy')">
                {{ activeDetail.restart_policy || '-' }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.createdAt')">
                {{ formatTime(activeDetail.created_at) }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.startedAt')">
                {{ formatTime(activeDetail.started_at) }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.detail.command')">
                {{ joinList(activeDetail.command) }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.detail.entrypoint')">
                {{ joinList(activeDetail.entrypoint) }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.detail.workingDir')">
                {{ activeDetail.working_dir || '-' }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.detail.inspectUpdatedAt')">
                {{ formatTime(activeDetail.inspect_updated_at) }}
              </t-descriptions-item>
            </t-descriptions>

            <t-descriptions :title="t('container.list.detail.runtime')" :column="1" bordered>
              <t-descriptions-item :label="t('container.list.fields.runtime')">
                {{ activeDetail.runtime_info.runtime }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.endpoint')">
                {{ activeDetail.runtime_info.endpoint || '-' }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.apiVersion')">
                {{ activeDetail.runtime_info.api_version || '-' }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.serverVersion')">
                {{ activeDetail.runtime_info.server_version || '-' }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.operatingSystem')">
                {{ activeDetail.runtime_info.operating_system || '-' }}
              </t-descriptions-item>
              <t-descriptions-item :label="t('container.list.fields.architecture')">
                {{ activeDetail.runtime_info.architecture || '-' }}
              </t-descriptions-item>
            </t-descriptions>

            <section class="container-detail-section">
              <h3>{{ t('container.list.detail.mounts') }}</h3>
              <div v-if="activeDetail.mounts.length" class="container-detail-list">
                <div
                  v-for="mount in activeDetail.mounts"
                  :key="`${mount.type}:${mount.destination}`"
                  class="container-detail-item"
                >
                  <strong>{{ mount.destination }}</strong>
                  <span>{{ mount.type }} / {{ mount.mode || '-' }} / {{ mount.read_only ? 'ro' : 'rw' }}</span>
                  <span>{{ mount.source || mount.name || '-' }}</span>
                </div>
              </div>
              <t-empty v-else size="small" :description="t('container.list.detail.mountEmpty')" />
            </section>

            <section class="container-detail-section">
              <h3>{{ t('container.list.detail.networks') }}</h3>
              <div v-if="activeDetail.networks.length" class="container-detail-list">
                <div v-for="network in activeDetail.networks" :key="network.name" class="container-detail-item">
                  <strong>{{ network.name }}</strong>
                  <span>{{ network.ip_address || '-' }}</span>
                  <span>{{ network.gateway || network.mac_address || '-' }}</span>
                </div>
              </div>
              <t-empty v-else size="small" :description="t('container.list.detail.networkEmpty')" />
            </section>
          </section>
        </t-loading>
      </div>
    </t-drawer>

    <t-drawer v-model:visible="logsDrawerVisible" :header="logsDrawerTitle" :footer="false" size="720px">
      <div class="container-drawer-panel container-logs-panel">
        <t-form class="container-log-controls" layout="inline" label-align="top" :data="logQuery">
          <t-form-item :label="t('container.list.logs.tail')" name="tail">
            <t-input-number v-model="logQuery.tail" theme="normal" :min="1" :max="2000" :step="100" />
          </t-form-item>
          <t-form-item :label="t('container.list.logs.since')" name="since">
            <t-input v-model="logQuery.since" clearable :placeholder="t('container.list.logs.sincePlaceholder')" />
          </t-form-item>
          <t-form-item name="streams">
            <t-space break-line size="small">
              <t-checkbox v-model="logQuery.timestamps">{{ t('container.list.logs.timestamps') }}</t-checkbox>
              <t-checkbox v-model="logQuery.stdout">{{ t('container.list.logs.stdout') }}</t-checkbox>
              <t-checkbox v-model="logQuery.stderr">{{ t('container.list.logs.stderr') }}</t-checkbox>
            </t-space>
          </t-form-item>
          <t-form-item>
            <t-space size="small">
              <t-button theme="primary" :loading="logsLoading" @click="refreshLogs">
                {{ t('container.list.logs.refresh') }}
              </t-button>
              <t-button theme="default" variant="outline" :disabled="!activeLogs?.lines.length" @click="copyLogs">
                {{ t('container.list.logs.copy') }}
              </t-button>
            </t-space>
          </t-form-item>
        </t-form>

        <t-alert v-if="logsError" class="container-alert" theme="error" :title="logsError" />
        <t-alert
          v-if="activeLogs?.truncated"
          class="container-alert"
          theme="warning"
          :title="t('container.list.logs.truncated')"
        />

        <t-loading :loading="logsLoading">
          <pre v-if="activeLogs?.lines.length" class="container-log-output">{{ activeLogs.lines.join('\n') }}</pre>
          <t-empty v-else size="small" :description="t('container.list.logs.empty')" />
        </t-loading>
      </div>
    </t-drawer>
  </div>
</template>
<script setup lang="ts">
import { SearchIcon } from 'tdesign-icons-vue-next';
import type { TdBaseTableProps } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next';
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  buildVisibleColumns,
  ManagementPageHeader,
  ManagementTableCard,
  ManagementTablePagination,
  ManagementToolbar,
  resolveTableWidthPolicy,
  TableActionMenu,
  TableViewToolbar,
  useTableHostWidth,
} from '@/shared/components/management';
import { AdvancedQueryColumnDrawer } from '@/shared/components/query-list';
import { resolveLocalizedErrorMessage } from '@/shared/localized-api-error';
import { formatLocaleDateTime } from '@/shared/observability';
import { usePermissionStore } from '@/store';
import { createLogger } from '@/utils/logger';

import { getContainer, getContainerLogs, getContainers, runContainerAction } from '../../api/container';
import { CONTAINER_PERMISSION_CODE } from '../../contract/permissions';
import type {
  ContainerAction,
  ContainerDetail,
  ContainerFilters,
  ContainerLogQuery,
  ContainerLogResponse,
  ContainerPort,
  ContainerRuntimeInfo,
  ContainerState,
  ContainerSummary,
} from '../../types/container';

defineOptions({
  name: 'ContainerListIndex',
});

const { locale, t } = useI18n();
const permissionStore = usePermissionStore();
const logger = createLogger('container.list');
const permissionCodes = CONTAINER_PERMISSION_CODE;

const statusOptions: ContainerState[] = [
  'created',
  'running',
  'paused',
  'restarting',
  'removing',
  'exited',
  'dead',
  'unknown',
];
const DEFAULT_LOG_QUERY: Required<ContainerLogQuery> = {
  tail: 200,
  since: '',
  timestamps: false,
  stdout: true,
  stderr: true,
};
const CONTAINER_RUNTIME_DISABLED_MESSAGE_KEY = 'ops.container.error.runtimeDisabled';
const CONTAINER_COLUMN_STORAGE_KEY = 'graft.container.list.visibleColumns';
const DEFAULT_VISIBLE_COLUMNS = ['state', 'name', 'image', 'ports', 'runtime_status', 'created_at', 'operation'];
const ALWAYS_VISIBLE_COLUMNS = ['state', 'name', 'operation'];
const ALL_COLUMN_KEYS = [
  'state',
  'name',
  'image',
  'ports',
  'runtime_status',
  'created_at',
  'started_at',
  'restart_policy',
  'image_id',
  'labels',
  'operation',
];
const CONTAINER_PORT_VISIBLE_LIMIT = 2;
const CONTAINER_DEFAULT_PAGE_SIZE = 20;

type ListErrorState = {
  title: string;
  hint: string;
};
type RowAction = 'copy-id' | 'logs' | ContainerAction;

const loading = ref(false);
const listError = ref<ListErrorState>({ title: '', hint: '' });
const rows = ref<ContainerSummary[]>([]);
const runtime = ref<ContainerRuntimeInfo | null>(null);
const detailDrawerVisible = ref(false);
const detailLoading = ref(false);
const detailError = ref('');
const selectedContainer = ref<ContainerSummary | null>(null);
const activeDetail = ref<ContainerDetail | null>(null);
const logsDrawerVisible = ref(false);
const logsLoading = ref(false);
const logsError = ref('');
const activeLogs = ref<ContainerLogResponse | null>(null);
const actionLoadingKey = ref('');
const columnDrawerVisible = ref(false);
const visibleColumnKeys = ref<string[]>(loadVisibleColumnKeys());
const tableDensity = ref<'medium' | 'small'>('medium');
const filters = reactive<ContainerFilters>({
  keyword: '',
  status: 'all',
});
const logQuery = reactive<Required<ContainerLogQuery>>({ ...DEFAULT_LOG_QUERY });
const pagination = reactive({
  current: 1,
  pageSize: CONTAINER_DEFAULT_PAGE_SIZE,
});

const allColumns = computed<TdBaseTableProps['columns']>(() => [
  { title: t('container.list.columns.status'), colKey: 'state', width: 104, align: 'center', ellipsis: false },
  {
    title: t('container.list.columns.name'),
    colKey: 'name',
    minWidth: 260,
    ellipsis: { theme: 'default', placement: 'top-left' },
  },
  {
    title: t('container.list.columns.image'),
    colKey: 'image',
    minWidth: 280,
    ellipsis: { theme: 'default', placement: 'top-left' },
  },
  { title: t('container.list.columns.ports'), colKey: 'ports', width: 220, ellipsis: false },
  {
    title: t('container.list.columns.runtimeStatus'),
    colKey: 'runtime_status',
    minWidth: 220,
    ellipsis: { theme: 'default', placement: 'top-left' },
  },
  { title: t('container.list.columns.createdAt'), colKey: 'created_at', width: 168, align: 'center' },
  { title: t('container.list.columns.startedAt'), colKey: 'started_at', width: 168, align: 'center' },
  { title: t('container.list.columns.restartPolicy'), colKey: 'restart_policy', width: 140, align: 'center' },
  {
    title: t('container.list.columns.imageId'),
    colKey: 'image_id',
    width: 220,
    ellipsis: { theme: 'default', placement: 'top-left' },
  },
  {
    title: t('container.list.columns.labels'),
    colKey: 'labels',
    width: 180,
    ellipsis: { theme: 'default', placement: 'top-left' },
  },
  {
    title: t('container.list.columns.operation'),
    colKey: 'operation',
    width: 176,
    fixed: 'right',
    align: 'center',
    ellipsis: false,
  },
]);
const visibleColumns = computed<TdBaseTableProps['columns']>(() =>
  buildVisibleColumns(allColumns.value, visibleColumnKeys.value, ALWAYS_VISIBLE_COLUMNS),
);
const { tableHostRef, tableHostWidth } = useTableHostWidth(() => visibleColumns.value);
const tableWidthPolicy = computed(() => resolveTableWidthPolicy(visibleColumns.value, tableHostWidth.value));
const filteredRows = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();
  return rows.value.filter((row) => {
    const matchesStatus = filters.status === 'all' || row.state === filters.status;
    if (!matchesStatus) return false;
    if (!keyword) return true;
    return [
      row.id,
      shortContainerId(row.id),
      row.image,
      row.status,
      row.runtime,
      row.restart_policy,
      ...row.names,
      ...formatPorts(row.ports),
    ].some((value) => value?.toLowerCase().includes(keyword) ?? false);
  });
});
const paginatedRows = computed(() => {
  const start = (pagination.current - 1) * pagination.pageSize;
  return filteredRows.value.slice(start, start + pagination.pageSize);
});
const hasActiveFilters = computed(() => Boolean(filters.keyword.trim()) || filters.status !== 'all');
const runningCount = computed(() => rows.value.filter((row) => row.state === 'running').length);
const stoppedCount = computed(
  () => rows.value.filter((row) => row.state === 'exited' || row.state === 'created').length,
);
const errorCount = computed(() => rows.value.filter((row) => row.state === 'dead' || row.state === 'unknown').length);
const runtimeStatusTheme = computed(() => {
  if (runtime.value?.status === 'enabled') return 'success';
  if (runtime.value?.status === 'disabled') return 'warning';
  return 'danger';
});
const runtimeSummary = computed(() => {
  if (!runtime.value) return t('container.list.runtimeUnavailable');
  const version = runtime.value.server_version || runtime.value.api_version || '';
  return version ? `${runtime.value.runtime} / ${version}` : runtime.value.runtime;
});
const runtimeCountParams = computed(() => ({
  running: runtime.value?.containers_running ?? 0,
  total: runtime.value?.containers_total ?? rows.value.length,
}));
const tableDensityLabel = computed(() =>
  tableDensity.value === 'medium' ? t('container.list.compactDensity') : t('container.list.defaultDensity'),
);
const columnSettingOptions = computed(() => [
  { label: t('container.list.columns.status'), value: 'state' },
  { label: t('container.list.columns.name'), value: 'name' },
  { label: t('container.list.columns.image'), value: 'image' },
  { label: t('container.list.columns.ports'), value: 'ports' },
  { label: t('container.list.columns.runtimeStatus'), value: 'runtime_status' },
  { label: t('container.list.columns.createdAt'), value: 'created_at' },
  { label: t('container.list.columns.startedAt'), value: 'started_at' },
  { label: t('container.list.columns.restartPolicy'), value: 'restart_policy' },
  { label: t('container.list.columns.imageId'), value: 'image_id' },
  { label: t('container.list.columns.labels'), value: 'labels' },
  { label: t('container.list.columns.operation'), value: 'operation' },
]);
const footerSummary = computed(() => {
  if (!filteredRows.value.length) {
    return t('container.list.pagination.empty');
  }

  const start = (pagination.current - 1) * pagination.pageSize + 1;
  const end = Math.min(pagination.current * pagination.pageSize, filteredRows.value.length);
  return t('container.list.pagination.summary', {
    end,
    start,
    total: filteredRows.value.length,
  });
});
const logsDrawerTitle = computed(() => {
  const containerName = selectedContainer.value ? displayName(selectedContainer.value) : '';
  return containerName ? `${t('container.list.logs.title')} - ${containerName}` : t('container.list.logs.title');
});

onMounted(() => {
  void refreshContainers();
});

watch(
  visibleColumnKeys,
  (keys) => {
    const normalizedKeys = normalizeVisibleColumnKeys(keys);
    if (normalizedKeys.join('|') !== keys.join('|')) {
      visibleColumnKeys.value = normalizedKeys;
      return;
    }
    persistVisibleColumnKeys(normalizedKeys);
  },
  { deep: true },
);

watch(
  () => [filters.status, filteredRows.value.length, pagination.pageSize],
  () => {
    const lastPage = Math.max(1, Math.ceil(filteredRows.value.length / pagination.pageSize));
    if (pagination.current > lastPage) {
      pagination.current = lastPage;
    }
  },
);

async function refreshContainers() {
  loading.value = true;
  listError.value = { title: '', hint: '' };
  try {
    const payload = await getContainers();
    rows.value = payload.items;
    runtime.value = payload.runtime;
  } catch (error) {
    rows.value = [];
    runtime.value = null;
    listError.value = resolveListError(error);
    logger.error('failed to fetch containers', error);
  } finally {
    loading.value = false;
  }
}

function resolveListError(error: unknown): ListErrorState {
  if (isApiRequestErrorShape(error) && error.messageKey === CONTAINER_RUNTIME_DISABLED_MESSAGE_KEY) {
    return {
      title: t(CONTAINER_RUNTIME_DISABLED_MESSAGE_KEY),
      hint: t('container.list.runtimeDisabledHint'),
    };
  }

  return {
    title: resolveLocalizedErrorMessage(t, error, t('container.list.loadFailed')),
    hint: '',
  };
}

function isApiRequestErrorShape(error: unknown): error is { isApiRequestError: true; messageKey?: string } {
  return Boolean(error && typeof error === 'object' && (error as { isApiRequestError?: unknown }).isApiRequestError);
}

function applyFilters() {
  filters.keyword = filters.keyword.trim();
  pagination.current = 1;
}

function resetFilters() {
  filters.keyword = '';
  filters.status = 'all';
  pagination.current = 1;
}

async function openDetail(row: ContainerSummary) {
  selectedContainer.value = row;
  activeDetail.value = null;
  detailDrawerVisible.value = true;
  await loadDetail(row.id);
}

async function loadDetail(containerId: string) {
  detailLoading.value = true;
  detailError.value = '';
  try {
    activeDetail.value = await getContainer(containerId);
  } catch (error) {
    detailError.value = resolveLocalizedErrorMessage(t, error, t('container.list.detail.loadFailed'));
    logger.warn('failed to fetch container detail', error);
  } finally {
    detailLoading.value = false;
  }
}

async function openLogs(row: ContainerSummary) {
  selectedContainer.value = row;
  activeLogs.value = null;
  logsError.value = '';
  Object.assign(logQuery, DEFAULT_LOG_QUERY);
  logsDrawerVisible.value = true;
  await refreshLogs();
}

async function refreshLogs() {
  if (!selectedContainer.value) return;
  logsLoading.value = true;
  logsError.value = '';
  try {
    activeLogs.value = await getContainerLogs(selectedContainer.value.id, normalizeLogQuery());
  } catch (error) {
    logsError.value = resolveLocalizedErrorMessage(t, error, t('container.list.logs.loadFailed'));
    logger.warn('failed to fetch container logs', error);
  } finally {
    logsLoading.value = false;
  }
}

function normalizeLogQuery(): ContainerLogQuery {
  return {
    tail: logQuery.tail,
    since: logQuery.since.trim() || undefined,
    timestamps: logQuery.timestamps,
    stdout: logQuery.stdout,
    stderr: logQuery.stderr,
  };
}

async function copyLogs() {
  const text = activeLogs.value?.lines.join('\n') ?? '';
  if (!text) return;
  try {
    await navigator.clipboard.writeText(text);
    MessagePlugin.success(t('container.list.copySuccess'));
  } catch (error) {
    logger.warn('failed to copy container logs', error);
    MessagePlugin.error(t('container.list.copyError'));
  }
}

async function copyContainerId(row: ContainerSummary) {
  try {
    await navigator.clipboard.writeText(row.id);
    MessagePlugin.success(t('container.list.copyIdSuccess'));
  } catch (error) {
    logger.warn('failed to copy container id', error);
    MessagePlugin.error(t('container.list.copyIdError'));
  }
}

async function runAction(action: ContainerAction, row: ContainerSummary) {
  const key = actionKey(action, row);
  actionLoadingKey.value = key;
  try {
    const result = await runContainerAction(action, row.id);
    MessagePlugin.success(localizedActionMessage(result.message_key, result.message));
    await refreshContainers();
  } catch (error) {
    logger.warn('failed to run container action', error);
    MessagePlugin.error(resolveLocalizedErrorMessage(t, error, t('container.list.actionFailed')));
  } finally {
    if (actionLoadingKey.value === key) {
      actionLoadingKey.value = '';
    }
  }
}

function localizedActionMessage(messageKey?: string, fallback?: string) {
  if (messageKey) {
    const translated = t(messageKey);
    if (translated !== messageKey) {
      return translated;
    }
  }
  return fallback?.trim() || t('container.list.actionSuccess');
}

function actionKey(action: ContainerAction, row: ContainerSummary) {
  return `${action}:${row.id}`;
}

function rowActions(row: ContainerSummary) {
  const actions: Array<{
    disabled?: boolean;
    fallbackLabel: string;
    label: string;
    testId: string;
    value: RowAction;
  }> = [];

  if (permissionStore.hasPermission(permissionCodes.LOGS)) {
    actions.push({
      fallbackLabel: t('container.list.actions.logs'),
      label: 'container.list.actions.logs',
      testId: 'container-action-logs',
      value: 'logs',
    });
  }

  if (permissionStore.hasPermission(permissionCodes.START)) {
    actions.push({
      disabled: rowActionDisabled('start', row),
      fallbackLabel: t('container.list.actions.start'),
      label: 'container.list.actions.start',
      testId: 'container-action-start',
      value: 'start',
    });
  }

  if (permissionStore.hasPermission(permissionCodes.STOP)) {
    actions.push({
      disabled: rowActionDisabled('stop', row),
      fallbackLabel: t('container.list.actions.stop'),
      label: 'container.list.actions.stop',
      testId: 'container-action-stop',
      value: 'stop',
    });
  }

  if (permissionStore.hasPermission(permissionCodes.RESTART)) {
    actions.push({
      disabled: rowActionDisabled('restart', row),
      fallbackLabel: t('container.list.actions.restart'),
      label: 'container.list.actions.restart',
      testId: 'container-action-restart',
      value: 'restart',
    });
  }

  actions.push({
    fallbackLabel: t('container.list.actions.copyId'),
    label: 'container.list.actions.copyId',
    testId: 'container-action-copy-id',
    value: 'copy-id',
  });

  return actions;
}

function handleRowAction(action: string, row: ContainerSummary) {
  if (action === 'logs') {
    void openLogs(row);
    return;
  }

  if (action === 'copy-id') {
    void copyContainerId(row);
    return;
  }

  if (action === 'start' || action === 'stop' || action === 'restart') {
    const messageKey = `container.list.actions.confirm${action[0].toUpperCase()}${action.slice(1)}`;
    const confirmed = window.confirm(t(messageKey));
    if (confirmed) {
      void runAction(action as ContainerAction, row);
    }
  }
}

function rowActionDisabled(action: ContainerAction, row?: ContainerSummary) {
  if (!row) return false;
  if (action === 'start') return row.state === 'running' || row.state === 'removing';
  if (action === 'stop') return row.state !== 'running';
  return row.state === 'removing' || row.state === 'dead';
}

function handlePageChange(pageInfo: { current?: number; pageSize?: number }) {
  if (pageInfo.current) {
    pagination.current = pageInfo.current;
  }
  if (pageInfo.pageSize) {
    pagination.pageSize = pageInfo.pageSize;
  }
}

function displayName(row: ContainerSummary | ContainerDetail) {
  return row.names[0] || row.id;
}

function shortContainerId(id: string) {
  return id.length > 12 ? id.slice(0, 12) : id;
}

function formatPorts(ports: ContainerPort[]) {
  return ports.map((port) => {
    const target = `${port.private_port}/${port.type}`;
    if (port.public_port === undefined) {
      return target;
    }
    return `${port.ip ? `${port.ip}:` : ''}${port.public_port}->${target}`;
  });
}

function visiblePortLabels(row: ContainerSummary) {
  return formatPorts(row.ports).slice(0, CONTAINER_PORT_VISIBLE_LIMIT);
}

function hiddenPortLabels(row: ContainerSummary) {
  return formatPorts(row.ports).slice(CONTAINER_PORT_VISIBLE_LIMIT);
}

function labelSummary(row: ContainerSummary) {
  const count = Object.keys(row.labels ?? {}).length;
  return count ? t('container.list.labelCount', { count }) : '-';
}

function formatTime(value?: string | null) {
  return formatLocaleDateTime(value, locale);
}

function joinList(values?: string[]) {
  return values?.length ? values.join(' ') : '-';
}

function stateLabel(state: ContainerState) {
  return t(`container.list.states.${state}`);
}

function stateTheme(state: ContainerState) {
  if (state === 'running') return 'success';
  if (state === 'created' || state === 'paused' || state === 'restarting') return 'warning';
  if (state === 'dead') return 'danger';
  return 'default';
}

function toggleTableDensity() {
  tableDensity.value = tableDensity.value === 'medium' ? 'small' : 'medium';
}

function loadVisibleColumnKeys() {
  if (typeof window === 'undefined') {
    return [...DEFAULT_VISIBLE_COLUMNS];
  }

  try {
    const stored = window.localStorage.getItem(CONTAINER_COLUMN_STORAGE_KEY);
    if (!stored) {
      return [...DEFAULT_VISIBLE_COLUMNS];
    }
    const parsed = JSON.parse(stored);
    if (!Array.isArray(parsed)) {
      return [...DEFAULT_VISIBLE_COLUMNS];
    }

    const normalizedKeys = normalizeVisibleColumnKeys(parsed);
    persistVisibleColumnKeys(normalizedKeys);
    return normalizedKeys;
  } catch {
    return [...DEFAULT_VISIBLE_COLUMNS];
  }
}

function persistVisibleColumnKeys(keys: string[]) {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    window.localStorage.setItem(CONTAINER_COLUMN_STORAGE_KEY, JSON.stringify(keys));
  } catch {
    // Column settings are a convenience preference; list rendering must not depend on storage availability.
  }
}

function normalizeVisibleColumnKeys(keys: unknown[]) {
  const availableKeySet = new Set(ALL_COLUMN_KEYS);
  const nextKeys = new Set<string>();

  for (const key of keys) {
    if (typeof key === 'string' && availableKeySet.has(key)) {
      nextKeys.add(key);
    }
  }

  for (const key of ALWAYS_VISIBLE_COLUMNS) {
    nextKeys.add(key);
  }

  return ALL_COLUMN_KEYS.filter((key) => nextKeys.has(key));
}
</script>
<style scoped lang="less">
.container-page {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-16);
  min-width: 0;
}

.container-table-head,
.container-image,
.container-identity,
.container-detail-item {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-4);
  min-width: 0;
}

.container-table-head p,
.container-detail-section h3 {
  margin: 0;
}

.container-table-head__summary,
.container-identity__name,
.container-detail-item strong {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-small);
}

.container-table-head p:not(.container-table-head__summary),
.container-identity__id,
.container-muted,
.container-detail-item span {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
}

.container-identity__id {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-port-list,
.container-actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--graft-density-gap-6);
}

.container-runtime-status {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-4);
  min-width: 0;
}

.container-runtime-status span:first-child {
  color: var(--td-text-color-primary);
  font: var(--td-font-body-medium);
}

.container-runtime-status span:last-child {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-alert {
  margin-bottom: var(--graft-density-gap-12);
}

.container-table-host {
  max-width: 100%;
  min-width: 0;
  overflow-x: hidden;
  width: 100%;
}

.container-table-host[data-table-mode='scroll'] {
  overflow-x: auto;
}

.container-table-host :deep(.t-table__content) {
  min-width: 0;
}

.container-table-host :deep(.t-table__content table) {
  min-width: 100%;
}

.container-alert__hint {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
  margin: var(--graft-density-gap-4) 0 0;
}

.container-drawer-panel,
.container-detail-stack,
.container-logs-panel,
.container-detail-list {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-16);
}

.container-detail-section {
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-medium);
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-12);
  padding: var(--graft-density-gap-14);
}

.container-detail-list {
  gap: var(--graft-density-gap-10);
}

.container-detail-item {
  border-bottom: 1px solid var(--td-component-stroke);
  padding-bottom: var(--graft-density-gap-10);
}

.container-detail-item:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.container-log-controls {
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-medium);
  padding: var(--graft-density-gap-14);
}

.container-log-output {
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-medium);
  color: var(--td-text-color-primary);
  font-family: var(--td-font-family-monospace);
  line-height: var(--td-line-height-body-medium);
  margin: 0;
  max-height: min(60vh, 640px);
  overflow: auto;
  overflow-wrap: anywhere;
  padding: var(--graft-density-gap-14);
  white-space: pre-wrap;
}

@media (width <= 768px) {
  .container-actions {
    justify-content: flex-start;
  }
}
</style>
