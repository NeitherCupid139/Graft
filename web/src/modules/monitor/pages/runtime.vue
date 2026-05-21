<template>
  <div class="monitor-detail-page" data-page-type="overview-dashboard">
    <header class="monitor-detail-page__header">
      <div class="monitor-detail-page__heading">
        <p class="monitor-detail-page__eyebrow">{{ t('monitor.sectionTitle') }}</p>
        <h1 class="monitor-detail-page__title">{{ t('monitor.runtimePage.title') }}</h1>
        <p class="monitor-detail-page__subtitle">{{ t('monitor.runtimePage.subtitle') }}</p>
      </div>
      <div class="monitor-detail-page__actions">
        <t-tag :theme="headerTheme" variant="light">{{ headerStatusLabel }}</t-tag>
        <t-button theme="primary" variant="outline" :loading="loading" @click="refreshSnapshot">
          {{ t('monitor.shared.refresh') }}
        </t-button>
      </div>
    </header>

    <t-card class="monitor-detail-page__note" :bordered="false">
      <div class="monitor-note">
        <h2 class="monitor-note__title">{{ t('monitor.runtimePage.memoryBoundaryTitle') }}</h2>
        <p class="monitor-note__description">{{ t('monitor.runtimePage.memoryBoundaryDescription') }}</p>
      </div>
    </t-card>

    <t-card v-if="errorMessage" class="monitor-detail-page__note is-warning" :bordered="false">
      <div class="monitor-note">
        <h2 class="monitor-note__title">{{ t('monitor.shared.errorTitle') }}</h2>
        <p class="monitor-note__description">{{ errorMessage }}</p>
      </div>
    </t-card>

    <section class="monitor-detail-page__grid monitor-detail-page__grid--summary">
      <t-card v-for="metric in summaryMetrics" :key="metric.key" class="monitor-detail-page__card" :bordered="false">
        <div class="monitor-summary-metric">
          <span class="monitor-summary-metric__label">{{ metric.label }}</span>
          <strong class="monitor-summary-metric__value">{{ metric.value }}</strong>
          <span class="monitor-summary-metric__description">{{ metric.description }}</span>
        </div>
      </t-card>
    </section>

    <section class="monitor-detail-page__grid monitor-detail-page__grid--detail">
      <t-card class="monitor-detail-page__card" :bordered="false" :title="t('monitor.runtimePage.runtimeMemoryTitle')">
        <div class="monitor-kv-grid">
          <div v-for="field in runtimeMemoryFields" :key="field.key" class="monitor-kv">
            <span class="monitor-kv__label">{{ field.label }}</span>
            <strong class="monitor-kv__value">{{ field.value }}</strong>
            <span v-if="field.description" class="monitor-kv__description">{{ field.description }}</span>
          </div>
        </div>
      </t-card>

      <t-card class="monitor-detail-page__card" :bordered="false" :title="t('monitor.runtimePage.processBuildTitle')">
        <div class="monitor-kv-grid">
          <div v-for="field in processAndBuildFields" :key="field.key" class="monitor-kv">
            <span class="monitor-kv__label">{{ field.label }}</span>
            <strong class="monitor-kv__value">{{ field.value }}</strong>
            <span v-if="field.description" class="monitor-kv__description">{{ field.description }}</span>
          </div>
        </div>
      </t-card>

      <t-card
        class="monitor-detail-page__card"
        :bordered="false"
        :title="t('monitor.runtimePage.hostEnvironmentTitle')"
      >
        <div class="monitor-kv-grid">
          <div v-for="field in hostEnvironmentFields" :key="field.key" class="monitor-kv">
            <span class="monitor-kv__label">{{ field.label }}</span>
            <strong class="monitor-kv__value">{{ field.value }}</strong>
            <span v-if="field.description" class="monitor-kv__description">{{ field.description }}</span>
          </div>
        </div>
      </t-card>

      <t-card
        class="monitor-detail-page__card"
        :bordered="false"
        :title="t('monitor.runtimePage.snapshotContextTitle')"
      >
        <div class="monitor-kv-grid">
          <div v-for="field in snapshotContextFields" :key="field.key" class="monitor-kv">
            <span class="monitor-kv__label">{{ field.label }}</span>
            <strong class="monitor-kv__value">{{ field.value }}</strong>
            <span v-if="field.description" class="monitor-kv__description">{{ field.description }}</span>
          </div>
        </div>
      </t-card>
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

const summaryMetrics = computed(() => {
  const response = serverStatus.value;

  return [
    {
      key: 'uptime',
      label: t('monitor.runtimePage.summary.uptime'),
      value: formatUptime(response?.server.uptime_seconds),
      description: t('monitor.runtimePage.summary.uptimeDescription'),
    },
    {
      key: 'goroutines',
      label: t('monitor.runtimePage.summary.goroutines'),
      value:
        response?.runtime.goroutines !== undefined
          ? String(response.runtime.goroutines)
          : t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.summary.goroutinesDescription'),
    },
    {
      key: 'goVersion',
      label: t('monitor.runtimePage.summary.goVersion'),
      value: displayText(response?.runtime.go_version ?? response?.server.go_version),
      description: t('monitor.runtimePage.summary.goVersionDescription'),
    },
    {
      key: 'gcCycles',
      label: t('monitor.runtimePage.summary.gcCycles'),
      value:
        response?.runtime.runtime_gc_cycles !== undefined
          ? String(response.runtime.runtime_gc_cycles)
          : t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.summary.gcCyclesDescription'),
    },
  ];
});

const runtimeMemoryFields = computed(() => {
  const runtime = serverStatus.value?.runtime;

  return [
    {
      key: 'alloc',
      label: t('monitor.runtimePage.fields.runtimeAlloc'),
      value: formatBytes(runtime?.runtime_alloc_bytes),
      description: t('monitor.runtimePage.fieldDescriptions.runtimeAlloc'),
    },
    {
      key: 'heap',
      label: t('monitor.runtimePage.fields.runtimeHeap'),
      value: formatBytes(runtime?.runtime_heap_in_use_bytes),
      description: t('monitor.runtimePage.fieldDescriptions.runtimeHeap'),
    },
    {
      key: 'sys',
      label: t('monitor.runtimePage.fields.runtimeSys'),
      value: formatBytes(runtime?.runtime_sys_bytes),
      description: t('monitor.runtimePage.fieldDescriptions.runtimeSys'),
    },
    {
      key: 'gcCycles',
      label: t('monitor.runtimePage.fields.gcCycles'),
      value:
        runtime?.runtime_gc_cycles !== undefined ? String(runtime.runtime_gc_cycles) : t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.fieldDescriptions.gcCycles'),
    },
    {
      key: 'lastGc',
      label: t('monitor.runtimePage.fields.lastGc'),
      value: t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.fieldDescriptions.lastGc'),
    },
  ];
});

const processAndBuildFields = computed(() => {
  const server = serverStatus.value?.server;

  return [
    {
      key: 'buildVersion',
      label: t('monitor.runtimePage.fields.buildVersion'),
      value: displayText(server?.version),
      description: t('monitor.runtimePage.fieldDescriptions.buildVersion'),
    },
    {
      key: 'gitCommit',
      label: t('monitor.runtimePage.fields.gitCommit'),
      value: t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.fieldDescriptions.gitCommit'),
    },
    {
      key: 'appName',
      label: t('monitor.runtimePage.fields.appName'),
      value: displayText(server?.app_name),
      description: t('monitor.runtimePage.fieldDescriptions.appName'),
    },
    {
      key: 'appEnv',
      label: t('monitor.runtimePage.fields.appEnv'),
      value: displayText(server?.app_env),
      description: t('monitor.runtimePage.fieldDescriptions.appEnv'),
    },
    {
      key: 'startedAt',
      label: t('monitor.runtimePage.fields.startedAt'),
      value: formatTimestamp(server?.started_at),
      description: t('monitor.runtimePage.fieldDescriptions.startedAt'),
    },
  ];
});

const hostEnvironmentFields = computed(() => {
  const runtime = serverStatus.value?.runtime;

  return [
    {
      key: 'hostName',
      label: t('monitor.runtimePage.fields.hostName'),
      value: displayText(runtime?.host_name),
      description: t('monitor.runtimePage.fieldDescriptions.hostName'),
    },
    {
      key: 'platform',
      label: t('monitor.runtimePage.fields.platform'),
      value:
        runtime?.operating_system && runtime?.architecture
          ? `${runtime.operating_system} / ${runtime.architecture}`
          : t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.fieldDescriptions.platform'),
    },
    {
      key: 'cpuCores',
      label: t('monitor.runtimePage.fields.cpuCores'),
      value: runtime?.cpu_cores !== undefined ? String(runtime.cpu_cores) : t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.fieldDescriptions.cpuCores'),
    },
    {
      key: 'hostMemory',
      label: t('monitor.runtimePage.fields.hostMemory'),
      value: runtime
        ? `${formatBytes(runtime.host_memory_used_bytes)} / ${formatBytes(runtime.host_memory_total_bytes)}`
        : t('monitor.shared.notReported'),
      description: t('monitor.runtimePage.fieldDescriptions.hostMemory'),
    },
  ];
});

const snapshotContextFields = computed(() => [
  {
    key: 'observedAt',
    label: t('monitor.runtimePage.fields.observedAt'),
    value: formatTimestamp(observedAt.value),
    description: t('monitor.runtimePage.fieldDescriptions.observedAt'),
  },
  {
    key: 'loadAverage',
    label: t('monitor.runtimePage.fields.loadAverage'),
    value: serverStatus.value
      ? [
          serverStatus.value.runtime.load_average.one_minute.toFixed(2),
          serverStatus.value.runtime.load_average.five_minutes.toFixed(2),
          serverStatus.value.runtime.load_average.fifteen_minutes.toFixed(2),
        ].join(' / ')
      : t('monitor.shared.notReported'),
    description: t('monitor.runtimePage.fieldDescriptions.loadAverage'),
  },
]);
</script>
