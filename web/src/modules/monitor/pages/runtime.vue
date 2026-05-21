<template>
  <div class="monitor-detail-page monitor-detail-page--runtime" data-page-type="overview-dashboard">
    <header class="monitor-detail-page__header">
      <div class="monitor-detail-page__heading">
        <p class="monitor-detail-page__eyebrow">{{ t('monitor.sectionTitle') }}</p>
        <h1 class="monitor-detail-page__title">{{ t('monitor.runtimePage.title') }}</h1>
        <p class="monitor-detail-page__subtitle">{{ t('monitor.runtimePage.subtitle') }}</p>
        <div class="monitor-runtime-meta-bar">
          <div v-for="item in snapshotMetaItems" :key="item.key" class="monitor-runtime-meta-bar__item">
            <span class="monitor-runtime-header-summary__label">{{ item.label }}</span>
            <strong class="monitor-runtime-header-summary__value">{{ item.value }}</strong>
          </div>
        </div>
      </div>
      <div class="monitor-detail-page__actions">
        <t-tag :theme="headerTheme" variant="light">{{ headerStatusLabel }}</t-tag>
        <t-button theme="primary" variant="outline" :loading="loading" @click="refreshSnapshot">
          {{ t('monitor.shared.refresh') }}
        </t-button>
      </div>
    </header>

    <div class="monitor-runtime-tip" role="note">
      <span class="monitor-runtime-tip__content">{{ t('monitor.runtimePage.memoryBoundaryNotice') }}</span>
    </div>

    <t-card v-if="errorMessage" class="monitor-detail-page__note is-warning" :bordered="false">
      <div class="monitor-note">
        <h2 class="monitor-note__title">{{ t('monitor.shared.errorTitle') }}</h2>
        <p class="monitor-note__description">{{ errorMessage }}</p>
      </div>
    </t-card>

    <section class="monitor-detail-page__grid monitor-detail-page__grid--summary monitor-runtime-summary-grid">
      <t-card
        v-for="metric in summaryMetrics"
        :key="metric.key"
        class="monitor-detail-page__card monitor-runtime-summary-card"
        :bordered="false"
      >
        <div class="monitor-summary-metric">
          <span class="monitor-summary-metric__label">{{ metric.label }}</span>
          <strong class="monitor-summary-metric__value">{{ metric.value }}</strong>
          <span class="monitor-summary-metric__description">{{ metric.description }}</span>
        </div>
      </t-card>
    </section>

    <section class="monitor-runtime-layout">
      <t-card
        class="monitor-detail-page__card monitor-runtime-memory-card monitor-runtime-layout__main"
        :bordered="false"
      >
        <div class="monitor-runtime-card">
          <header class="monitor-runtime-card__heading">
            <div class="monitor-runtime-card__copy">
              <h2 class="monitor-runtime-card__title">{{ t('monitor.runtimePage.runtimeMemoryTitle') }}</h2>
              <p class="monitor-runtime-card__description">{{ t('monitor.runtimePage.runtimeMemoryDescription') }}</p>
            </div>
          </header>

          <div class="monitor-runtime-memory">
            <div class="monitor-runtime-memory__hero">
              <div v-for="field in runtimePrimaryFields" :key="field.key" class="monitor-runtime-memory__primary">
                <span class="monitor-runtime-memory__label">{{ field.label }}</span>
                <strong class="monitor-runtime-memory__value">{{ field.value }}</strong>
                <span class="monitor-runtime-memory__description">{{ field.description }}</span>
              </div>
            </div>

            <div class="monitor-runtime-memory__meta">
              <div v-for="field in runtimeSecondaryFields" :key="field.key" class="monitor-runtime-memory__secondary">
                <span class="monitor-runtime-memory__label">{{ field.label }}</span>
                <strong class="monitor-runtime-memory__value is-secondary">{{ field.value }}</strong>
                <span class="monitor-runtime-memory__description">{{ field.description }}</span>
              </div>
            </div>
          </div>
        </div>
      </t-card>

      <div class="monitor-runtime-layout__sidebar">
        <t-card class="monitor-detail-page__card monitor-runtime-info-card" :bordered="false">
          <div class="monitor-runtime-card">
            <header class="monitor-runtime-card__heading">
              <div class="monitor-runtime-card__copy">
                <h2 class="monitor-runtime-card__title">{{ t('monitor.runtimePage.processBuildTitle') }}</h2>
              </div>
            </header>

            <dl class="monitor-description-list">
              <div v-for="field in processAndBuildFields" :key="field.key" class="monitor-description-list__item">
                <dt class="monitor-description-list__term">{{ field.label }}</dt>
                <dd class="monitor-description-list__value">{{ field.value }}</dd>
              </div>
            </dl>
          </div>
        </t-card>

        <t-card class="monitor-detail-page__card monitor-runtime-info-card" :bordered="false">
          <div class="monitor-runtime-card">
            <header class="monitor-runtime-card__heading">
              <div class="monitor-runtime-card__copy">
                <h2 class="monitor-runtime-card__title">{{ t('monitor.runtimePage.hostEnvironmentTitle') }}</h2>
                <p class="monitor-runtime-card__description">
                  {{ t('monitor.runtimePage.serverEnvironmentDescription') }}
                </p>
              </div>
            </header>

            <dl class="monitor-description-list">
              <div v-for="field in hostEnvironmentFields" :key="field.key" class="monitor-description-list__item">
                <dt class="monitor-description-list__term">{{ field.label }}</dt>
                <dd class="monitor-description-list__value">{{ field.value }}</dd>
              </div>
            </dl>
          </div>
        </t-card>
      </div>
    </section>

    <t-empty v-if="initialized && !serverStatus && !loading" :description="t('monitor.shared.empty')" />
  </div>
</template>
<script setup lang="ts">
import './detail-page.less';

import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  displayText,
  formatBytes,
  formatTimestamp,
  formatUptime,
  runtimeSnapshotTheme,
  useServerStatusSnapshot,
} from './server-status-snapshot';

const { t } = useI18n();
const { errorMessage, initialized, loading, observedAt, refreshSnapshot, serverStatus } = useServerStatusSnapshot();

const headerTheme = computed(() => runtimeSnapshotTheme(serverStatus.value?.status));
const headerStatusLabel = computed(() =>
  serverStatus.value ? t('monitor.runtimePage.snapshotReady') : t('monitor.runtimePage.snapshotPending'),
);
const notReportedLabel = computed(() => t('monitor.shared.notReported'));

const snapshotMetaItems = computed(() => [
  {
    key: 'observedAt',
    label: t('monitor.runtimePage.fields.observedAt'),
    value: formatSnapshotTimestamp(observedAt.value),
  },
  {
    key: 'refreshFrequency',
    label: t('monitor.runtimePage.fields.refreshFrequency'),
    value: formatRefreshFrequency(serverStatus.value?.trend.sample_interval_seconds),
  },
  {
    key: 'loadAverage',
    label: t('monitor.runtimePage.fields.loadAverage'),
    value: formatLoadAverage(serverStatus.value?.runtime.load_average),
  },
]);

const summaryMetrics = computed(() => {
  const response = serverStatus.value;

  return [
    {
      key: 'uptime',
      label: t('monitor.runtimePage.summary.uptime'),
      value: formatSnapshotUptime(response?.server.uptime_seconds),
      description: t('monitor.runtimePage.summary.uptimeDescription'),
    },
    {
      key: 'goroutines',
      label: t('monitor.runtimePage.summary.goroutines'),
      value: formatCount(response?.runtime.goroutines),
      description: t('monitor.runtimePage.summary.goroutinesDescription'),
    },
    {
      key: 'runtimeAlloc',
      label: t('monitor.runtimePage.summary.runtimeAlloc'),
      value: formatSnapshotBytes(response?.runtime.runtime_alloc_bytes),
      description: t('monitor.runtimePage.summary.runtimeAllocDescription'),
    },
    {
      key: 'gcCycles',
      label: t('monitor.runtimePage.summary.gcCycles'),
      value: formatCount(response?.runtime.runtime_gc_cycles),
      description: t('monitor.runtimePage.summary.gcCyclesDescription'),
    },
  ];
});

const runtimePrimaryFields = computed(() => {
  const runtime = serverStatus.value?.runtime;

  return [
    {
      key: 'alloc',
      label: t('monitor.runtimePage.fields.runtimeAlloc'),
      value: formatSnapshotBytes(runtime?.runtime_alloc_bytes),
      description: t('monitor.runtimePage.fieldDescriptions.runtimeAlloc'),
    },
    {
      key: 'heap',
      label: t('monitor.runtimePage.fields.runtimeHeap'),
      value: formatSnapshotBytes(runtime?.runtime_heap_in_use_bytes),
      description: t('monitor.runtimePage.fieldDescriptions.runtimeHeap'),
    },
  ];
});

const runtimeSecondaryFields = computed(() => {
  const runtime = serverStatus.value?.runtime;

  return [
    {
      key: 'sys',
      label: t('monitor.runtimePage.fields.runtimeSys'),
      value: formatSnapshotBytes(runtime?.runtime_sys_bytes),
      description: t('monitor.runtimePage.fieldDescriptions.runtimeSys'),
    },
    {
      key: 'gcCycles',
      label: t('monitor.runtimePage.fields.gcCycles'),
      value: formatCount(runtime?.runtime_gc_cycles),
      description: t('monitor.runtimePage.fieldDescriptions.gcCycles'),
    },
    {
      key: 'lastGc',
      label: t('monitor.runtimePage.fields.lastGc'),
      value: notReportedLabel.value,
      description: t('monitor.runtimePage.fieldDescriptions.lastGc'),
    },
  ];
});

const processAndBuildFields = computed(() => {
  const server = serverStatus.value?.server;
  const runtime = serverStatus.value?.runtime;

  return [
    {
      key: 'appName',
      label: t('monitor.runtimePage.fields.appName'),
      value: displaySnapshotText(server?.app_name),
    },
    {
      key: 'appEnv',
      label: t('monitor.runtimePage.fields.appEnv'),
      value: displaySnapshotText(server?.app_env),
    },
    {
      key: 'buildVersion',
      label: t('monitor.runtimePage.fields.buildVersion'),
      value: displaySnapshotText(server?.version),
    },
    {
      key: 'gitCommit',
      label: t('monitor.runtimePage.fields.gitCommit'),
      value: notReportedLabel.value,
    },
    {
      key: 'startedAt',
      label: t('monitor.runtimePage.fields.startedAt'),
      value: formatSnapshotTimestamp(server?.started_at),
    },
    {
      key: 'goVersion',
      label: t('monitor.runtimePage.fields.goVersion'),
      value: displaySnapshotText(runtime?.go_version ?? server?.go_version),
    },
  ];
});

const hostEnvironmentFields = computed(() => {
  const runtime = serverStatus.value?.runtime;

  return [
    {
      key: 'hostName',
      label: t('monitor.runtimePage.fields.hostName'),
      value: displaySnapshotText(runtime?.host_name),
    },
    {
      key: 'platform',
      label: t('monitor.runtimePage.fields.platform'),
      value:
        runtime?.operating_system && runtime?.architecture
          ? `${runtime.operating_system} / ${runtime.architecture}`
          : notReportedLabel.value,
    },
    {
      key: 'cpuCores',
      label: t('monitor.runtimePage.fields.cpuCores'),
      value: formatCount(runtime?.cpu_cores),
    },
    {
      key: 'hostMemory',
      label: t('monitor.runtimePage.fields.hostMemory'),
      value: formatHostMemory(runtime),
    },
  ];
});

function displaySnapshotText(value?: string | null) {
  const formatted = displayText(value);
  return formatted === '--' ? notReportedLabel.value : formatted;
}

function formatSnapshotBytes(value?: number | null) {
  const formatted = formatBytes(value);
  return formatted === '--' ? notReportedLabel.value : formatted;
}

function formatSnapshotTimestamp(value?: string | null) {
  const formatted = formatTimestamp(value);
  return formatted === '--' ? notReportedLabel.value : formatted;
}

function formatSnapshotUptime(value?: number | null) {
  const formatted = formatUptime(value);
  return formatted === '--' ? notReportedLabel.value : formatted;
}

function formatCount(value?: number | null) {
  return Number.isFinite(value) ? String(value) : notReportedLabel.value;
}

function formatRefreshFrequency(value?: number | null) {
  return Number.isFinite(value) ? `${value}${t('monitor.runtimePage.refreshFrequencyUnit')}` : notReportedLabel.value;
}

function formatLoadAverage(
  value?: {
    one_minute?: number;
    five_minutes?: number;
    fifteen_minutes?: number;
  } | null,
) {
  const oneMinute = value?.one_minute;
  const fiveMinutes = value?.five_minutes;
  const fifteenMinutes = value?.fifteen_minutes;

  if (!Number.isFinite(oneMinute) || !Number.isFinite(fiveMinutes) || !Number.isFinite(fifteenMinutes)) {
    return notReportedLabel.value;
  }

  const safeOneMinute = Number(oneMinute);
  const safeFiveMinutes = Number(fiveMinutes);
  const safeFifteenMinutes = Number(fifteenMinutes);

  return `${safeOneMinute.toFixed(2)} / ${safeFiveMinutes.toFixed(2)} / ${safeFifteenMinutes.toFixed(2)}`;
}

function formatHostMemory(
  value?: {
    host_memory_total_bytes?: number;
    host_memory_used_bytes?: number;
    host_memory_used_percent?: number;
  } | null,
) {
  if (!value) {
    return notReportedLabel.value;
  }

  const used = formatSnapshotBytes(value.host_memory_used_bytes);
  const total = formatSnapshotBytes(value.host_memory_total_bytes);

  if (used === notReportedLabel.value || total === notReportedLabel.value) {
    return notReportedLabel.value;
  }

  const hostMemoryUsedPercent = value.host_memory_used_percent;
  const percent =
    typeof hostMemoryUsedPercent === 'number' && Number.isFinite(hostMemoryUsedPercent)
      ? ` (${hostMemoryUsedPercent.toFixed(0)}%)`
      : '';
  return `${used} / ${total}${percent}`;
}
</script>
