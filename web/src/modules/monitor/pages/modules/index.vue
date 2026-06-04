<template>
  <server-status-page-shell
    :eyebrow="t('monitor.sectionTitle')"
    :title="t('monitor.moduleRuntime.title')"
    :description="t('monitor.moduleRuntime.subtitle')"
    compact-header
  >
    <template #toolbar>
      <div class="module-runtime-toolbar">
        <status-tag :label="headerStatusLabel" :status="headerStatus" />
        <t-button
          class="module-runtime-toolbar__button"
          theme="primary"
          size="small"
          :loading="loading"
          @click="refreshSnapshot"
        >
          <template #icon>
            <refresh-icon />
          </template>
          {{ t('monitor.moduleRuntime.actions.refresh') }}
        </t-button>
      </div>
    </template>

    <template #summary>
      <div v-for="metric in summaryMetrics" :key="metric.key" class="module-runtime-summary-card">
        <t-statistic :title="metric.label" :value="metric.value" :loading="loading" />
        <p>{{ metric.description }}</p>
      </div>
    </template>

    <template #feedback>
      <t-alert
        v-if="errorMessage"
        theme="error"
        :title="t('monitor.moduleRuntime.errorTitle')"
        :message="errorMessage"
      />
    </template>

    <section-card
      :title="t('monitor.moduleRuntime.table.title')"
      :description="t('monitor.moduleRuntime.table.description')"
      :min-height="420"
    >
      <template #actions>
        <t-alert class="module-runtime-table__note" theme="info" :message="t('monitor.moduleRuntime.table.note')" />
      </template>

      <t-table
        row-key="module_key"
        hover
        :data="items"
        :columns="columns"
        :loading="loading"
        :empty="emptyTableContent"
        table-layout="fixed"
        table-content-width="1180px"
      >
        <template #module_key="{ row }">
          <strong class="module-runtime-table__key">{{ row.module_key }}</strong>
        </template>

        <template #enabled="{ row }">
          <status-tag :label="booleanLabel(row.enabled)" :status="row.enabled ? 'healthy' : 'disabled'" />
        </template>

        <template #registered="{ row }">
          <status-tag :label="booleanLabel(row.registered)" :status="row.registered ? 'healthy' : 'unknown'" />
        </template>

        <template #health="{ row }">
          <status-tag :label="healthLabel(row.health)" :status="healthTone(row.health)" />
        </template>

        <template #dependencies="{ row }">
          <span class="module-runtime-table__muted">
            {{ dependencySummary(row.dependencies) }}
          </span>
        </template>

        <template #migration="{ row }">
          <div class="module-runtime-table__stack">
            <status-tag
              :label="migrationStatusLabel(row.migration_status.status)"
              :status="declaredTone(row.migration_status.status)"
            />
            <span>{{
              t('monitor.moduleRuntime.values.migrationDirCount', { count: row.migration_status.declared_dirs.length })
            }}</span>
          </div>
        </template>

        <template #schema="{ row }">
          <status-tag
            :label="schemaStatusLabel(row.schema_status.status)"
            :status="declaredTone(row.schema_status.status)"
          />
        </template>

        <template #config="{ row }">
          <status-tag
            :label="configStatusLabel(row.config_status.status)"
            :status="configTone(row.config_status.status)"
          />
        </template>

        <template #operation="{ row }">
          <t-button variant="text" theme="primary" size="small" @click="openDetail(row)">
            {{ t('monitor.moduleRuntime.actions.detail') }}
          </t-button>
        </template>
      </t-table>

      <t-empty
        v-if="initialized && !loading && !items.length && !errorMessage"
        class="module-runtime-empty"
        :description="t('monitor.moduleRuntime.empty')"
      />
    </section-card>

    <t-drawer
      v-model:visible="detailVisible"
      :header="detailHeader"
      :footer="false"
      size="520px"
      attach="body"
      destroy-on-close
    >
      <div v-if="selectedModule" class="module-runtime-detail">
        <t-descriptions :column="1" item-layout="vertical" size="small">
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.moduleKey')">
            {{ selectedModule.module_key }}
          </t-descriptions-item>
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.runtimeStatus')">
            <status-tag
              :label="runtimeStatusLabel(selectedModule.runtime_status)"
              :status="runtimeTone(selectedModule.runtime_status)"
            />
          </t-descriptions-item>
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.enablementSource')">
            {{ enablementSourceLabel(selectedModule.enablement_source) }}
          </t-descriptions-item>
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.dependencies')">
            <div v-if="selectedModule.dependencies.length" class="module-runtime-detail__list">
              <div
                v-for="dependency in selectedModule.dependencies"
                :key="dependency.module_key"
                class="module-runtime-detail__line"
              >
                <span>{{ dependency.module_key }}</span>
                <status-tag
                  :label="dependencyStatusLabel(dependency.status)"
                  :status="dependencyTone(dependency.status)"
                />
              </div>
            </div>
            <span v-else>{{ t('monitor.moduleRuntime.values.none') }}</span>
          </t-descriptions-item>
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.migrationDirs')">
            <div v-if="selectedModule.migration_status.declared_dirs.length" class="module-runtime-detail__chips">
              <t-tag
                v-for="directory in selectedModule.migration_status.declared_dirs"
                :key="directory"
                variant="light"
              >
                {{ directory }}
              </t-tag>
            </div>
            <span v-else>{{ t('monitor.moduleRuntime.values.none') }}</span>
          </t-descriptions-item>
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.schemaStatus')">
            <status-tag
              :label="schemaStatusLabel(selectedModule.schema_status.status)"
              :status="declaredTone(selectedModule.schema_status.status)"
            />
          </t-descriptions-item>
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.configStatus')">
            <status-tag
              :label="configStatusLabel(selectedModule.config_status.status)"
              :status="configTone(selectedModule.config_status.status)"
            />
          </t-descriptions-item>
          <t-descriptions-item :label="t('monitor.moduleRuntime.detail.diagnostics')">
            <div v-if="diagnosticEntries.length" class="module-runtime-detail__diagnostics">
              <div v-for="[key, value] in diagnosticEntries" :key="key" class="module-runtime-detail__diagnostic">
                <span>{{ key }}</span>
                <strong>{{ value }}</strong>
              </div>
            </div>
            <span v-else>{{ t('monitor.moduleRuntime.values.none') }}</span>
          </t-descriptions-item>
        </t-descriptions>
      </div>
    </t-drawer>
  </server-status-page-shell>
</template>
<script setup lang="ts">
import { RefreshIcon } from 'tdesign-icons-vue-next';
import type { TdBaseTableProps } from 'tdesign-vue-next';
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { getModuleRuntimeSnapshot } from '../../api/module-runtime';
import SectionCard from '../../components/SectionCard.vue';
import type { ServerStatusTone } from '../../components/server-status-ui';
import ServerStatusPageShell from '../../components/ServerStatusPageShell.vue';
import StatusTag from '../../components/StatusTag.vue';
import type {
  ModuleRuntimeConfigStatus,
  ModuleRuntimeDependency,
  ModuleRuntimeItem,
  ModuleRuntimeMigrationStatus,
  ModuleRuntimeSchemaStatus,
  ModuleRuntimeSnapshot,
} from '../../types/module-runtime';

const { t } = useI18n();

const snapshot = ref<ModuleRuntimeSnapshot | null>(null);
const loading = ref(false);
const initialized = ref(false);
const errorMessage = ref('');
const detailVisible = ref(false);
const selectedModule = ref<ModuleRuntimeItem | null>(null);

const items = computed(() => snapshot.value?.items ?? []);
const summary = computed(() => snapshot.value?.summary);

const summaryMetrics = computed(() => [
  {
    key: 'total',
    label: t('monitor.moduleRuntime.summary.total'),
    value: summary.value?.total_modules ?? 0,
    description: t('monitor.moduleRuntime.summary.totalDescription'),
  },
  {
    key: 'enabled',
    label: t('monitor.moduleRuntime.summary.enabled'),
    value: summary.value?.enabled_modules ?? 0,
    description: t('monitor.moduleRuntime.summary.enabledDescription'),
  },
  {
    key: 'healthy',
    label: t('monitor.moduleRuntime.summary.healthy'),
    value: summary.value?.healthy_modules ?? 0,
    description: t('monitor.moduleRuntime.summary.healthyDescription'),
  },
  {
    key: 'degradedUnknown',
    label: t('monitor.moduleRuntime.summary.degradedUnknown'),
    value: (summary.value?.degraded_modules ?? 0) + (summary.value?.unknown_modules ?? 0),
    description: t('monitor.moduleRuntime.summary.degradedUnknownDescription'),
  },
]);

const headerStatus = computed<ServerStatusTone>(() => {
  const currentSummary = summary.value;
  if (!currentSummary) {
    return 'unknown';
  }

  if (currentSummary.degraded_modules > 0) {
    return 'warning';
  }

  if (currentSummary.unknown_modules > 0) {
    return 'unknown';
  }

  return 'healthy';
});

const headerStatusLabel = computed(() => {
  switch (headerStatus.value) {
    case 'healthy':
      return t('monitor.moduleRuntime.status.ready');
    case 'warning':
      return t('monitor.moduleRuntime.status.attention');
    default:
      return t('monitor.moduleRuntime.status.unknown');
  }
});

const emptyTableContent = computed(() => (initialized.value && !loading.value ? t('monitor.moduleRuntime.empty') : ''));

const detailHeader = computed(() =>
  selectedModule.value
    ? t('monitor.moduleRuntime.detail.titleWithKey', { key: selectedModule.value.module_key })
    : t('monitor.moduleRuntime.detail.title'),
);

const diagnosticEntries = computed(() => Object.entries(selectedModule.value?.diagnostics ?? {}));

const columns = computed<TdBaseTableProps['columns']>(() => [
  {
    colKey: 'module_key',
    title: t('monitor.moduleRuntime.columns.moduleKey'),
    width: 180,
    fixed: 'left',
  },
  {
    colKey: 'enabled',
    title: t('monitor.moduleRuntime.columns.enabled'),
    width: 110,
  },
  {
    colKey: 'registered',
    title: t('monitor.moduleRuntime.columns.registered'),
    width: 120,
  },
  {
    colKey: 'health',
    title: t('monitor.moduleRuntime.columns.health'),
    width: 130,
  },
  {
    colKey: 'dependencies',
    title: t('monitor.moduleRuntime.columns.dependencies'),
    width: 150,
  },
  {
    colKey: 'migration',
    title: t('monitor.moduleRuntime.columns.migration'),
    width: 180,
  },
  {
    colKey: 'schema',
    title: t('monitor.moduleRuntime.columns.schema'),
    width: 140,
  },
  {
    colKey: 'config',
    title: t('monitor.moduleRuntime.columns.config'),
    width: 150,
  },
  {
    colKey: 'operation',
    title: t('monitor.moduleRuntime.columns.action'),
    width: 110,
    fixed: 'right',
  },
]);

onMounted(() => {
  void refreshSnapshot();
});

async function refreshSnapshot() {
  loading.value = true;
  errorMessage.value = '';

  try {
    snapshot.value = await getModuleRuntimeSnapshot();
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : t('monitor.moduleRuntime.errorFallback');
  } finally {
    loading.value = false;
    initialized.value = true;
  }
}

function openDetail(row: ModuleRuntimeItem) {
  selectedModule.value = row;
  detailVisible.value = true;
}

function booleanLabel(value: boolean) {
  return value ? t('monitor.moduleRuntime.values.yes') : t('monitor.moduleRuntime.values.no');
}

function healthTone(status: ModuleRuntimeItem['health']): ServerStatusTone {
  switch (status) {
    case 'healthy':
      return 'healthy';
    case 'degraded':
      return 'warning';
    case 'disabled':
      return 'disabled';
    default:
      return 'unknown';
  }
}

function healthLabel(status: ModuleRuntimeItem['health']) {
  return t(`monitor.moduleRuntime.health.${status}`);
}

function runtimeTone(status: ModuleRuntimeItem['runtime_status']): ServerStatusTone {
  switch (status) {
    case 'registered':
      return 'healthy';
    case 'degraded':
      return 'warning';
    case 'disabled':
      return 'disabled';
    default:
      return 'unknown';
  }
}

function runtimeStatusLabel(status: ModuleRuntimeItem['runtime_status']) {
  return t(`monitor.moduleRuntime.runtimeStatus.${status}`);
}

function declaredTone(
  status: ModuleRuntimeMigrationStatus['status'] | ModuleRuntimeSchemaStatus['status'],
): ServerStatusTone {
  return status === 'declared' ? 'healthy' : 'unknown';
}

function migrationStatusLabel(status: ModuleRuntimeMigrationStatus['status']) {
  return t(`monitor.moduleRuntime.migrationStatus.${status}`);
}

function schemaStatusLabel(status: ModuleRuntimeSchemaStatus['status']) {
  return t(`monitor.moduleRuntime.schemaStatus.${status}`);
}

function configTone(status: ModuleRuntimeConfigStatus['status']): ServerStatusTone {
  return status === 'not_required' ? 'disabled' : 'unknown';
}

function configStatusLabel(status: ModuleRuntimeConfigStatus['status']) {
  return t(`monitor.moduleRuntime.configStatus.${status}`);
}

function dependencyTone(status: ModuleRuntimeDependency['status']): ServerStatusTone {
  switch (status) {
    case 'satisfied':
      return 'healthy';
    case 'disabled':
      return 'disabled';
    default:
      return 'warning';
  }
}

function dependencyStatusLabel(status: ModuleRuntimeDependency['status']) {
  return t(`monitor.moduleRuntime.dependencyStatus.${status}`);
}

function dependencySummary(dependencies: ModuleRuntimeDependency[]) {
  if (!dependencies.length) {
    return t('monitor.moduleRuntime.values.none');
  }

  const satisfied = dependencies.filter((dependency) => dependency.status === 'satisfied').length;
  return t('monitor.moduleRuntime.values.dependencySummary', {
    satisfied,
    total: dependencies.length,
  });
}

function enablementSourceLabel(source: ModuleRuntimeItem['enablement_source']) {
  return t(`monitor.moduleRuntime.enablementSource.${source}`);
}
</script>
<style scoped lang="less">
.module-runtime-toolbar {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: flex-end;
}

.module-runtime-toolbar__button {
  white-space: nowrap;
}

.module-runtime-summary-card {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-border-level-1-color);
  border-radius: var(--td-radius-medium);
  min-height: 112px;
  padding: 16px;
}

.module-runtime-summary-card p {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
  margin: 8px 0 0;
}

.module-runtime-table__note {
  max-width: 520px;
}

.module-runtime-table__key {
  color: var(--td-text-color-primary);
  font: var(--td-font-body-medium);
  overflow-wrap: anywhere;
}

.module-runtime-table__muted {
  color: var(--td-text-color-secondary);
}

.module-runtime-table__stack {
  align-items: flex-start;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.module-runtime-table__stack span:last-child {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
}

.module-runtime-empty {
  margin-top: 20px;
}

.module-runtime-detail {
  padding-bottom: 12px;
}

.module-runtime-detail__list,
.module-runtime-detail__diagnostics {
  display: grid;
  gap: 10px;
}

.module-runtime-detail__line,
.module-runtime-detail__diagnostic {
  align-items: center;
  background: var(--td-bg-color-container-hover);
  border: 1px solid var(--td-border-level-1-color);
  border-radius: var(--td-radius-default);
  display: flex;
  gap: 12px;
  justify-content: space-between;
  min-width: 0;
  padding: 10px 12px;
}

.module-runtime-detail__line span,
.module-runtime-detail__diagnostic span,
.module-runtime-detail__diagnostic strong {
  min-width: 0;
  overflow-wrap: anywhere;
}

.module-runtime-detail__diagnostic span {
  color: var(--td-text-color-secondary);
}

.module-runtime-detail__diagnostic strong {
  color: var(--td-text-color-primary);
  font-weight: 500;
  text-align: right;
}

.module-runtime-detail__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media (width <= 767px) {
  .module-runtime-toolbar {
    justify-content: flex-start;
  }

  .module-runtime-summary-card {
    min-height: 96px;
  }
}
</style>
