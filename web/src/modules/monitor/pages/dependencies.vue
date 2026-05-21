<template>
  <div class="monitor-detail-page" data-page-type="overview-dashboard">
    <header class="monitor-detail-page__header">
      <div class="monitor-detail-page__heading">
        <p class="monitor-detail-page__eyebrow">{{ t('monitor.sectionTitle') }}</p>
        <h1 class="monitor-detail-page__title">{{ t('monitor.dependenciesPage.title') }}</h1>
        <p class="monitor-detail-page__subtitle">{{ t('monitor.dependenciesPage.subtitle') }}</p>
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
        <h2 class="monitor-note__title">{{ t('monitor.dependenciesPage.noteTitle') }}</h2>
        <p class="monitor-note__description">{{ t('monitor.dependenciesPage.noteDescription') }}</p>
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
      <t-card
        v-for="service in serviceCards"
        :key="service.key"
        class="monitor-detail-page__status-card"
        :bordered="false"
      >
        <div class="monitor-service-card">
          <header class="monitor-service-card__header">
            <div class="monitor-service-card__heading">
              <h2 class="monitor-service-card__title">{{ service.title }}</h2>
              <p class="monitor-service-card__subtitle">{{ service.subtitle }}</p>
            </div>
            <t-tag :theme="service.theme" variant="light">{{ service.statusLabel }}</t-tag>
          </header>

          <div class="monitor-kv-grid">
            <div v-for="field in service.fields" :key="field.key" class="monitor-kv">
              <span class="monitor-kv__label">{{ field.label }}</span>
              <strong class="monitor-kv__value">{{ field.value }}</strong>
              <span v-if="field.description" class="monitor-kv__description">{{ field.description }}</span>
            </div>
          </div>

          <div v-if="service.futureEntry" class="monitor-service-card__future">
            {{ t('monitor.dependenciesPage.futureEntryDescription') }}
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
  dependencyStatusTheme,
  displayText,
  formatLatency,
  formatTimestamp,
  normalizeDependencyStatus,
  useServerStatusSnapshot,
} from './server-status-snapshot';

type DependencyCard = {
  key: string;
  title: string;
  subtitle: string;
  status: ReturnType<typeof normalizeDependencyStatus>;
  statusLabel: string;
  theme: 'success' | 'warning' | 'danger' | 'default';
  futureEntry: boolean;
  fields: Array<{
    key: string;
    label: string;
    value: string;
    description: string;
  }>;
};

const { t } = useI18n();
const { errorMessage, initialized, loading, observedAt, refreshSnapshot, serverStatus } = useServerStatusSnapshot();

const headerTheme = computed(() => dependencyStatusTheme(overallDependencyStatus.value));
const headerStatusLabel = computed(() => {
  switch (overallDependencyStatus.value) {
    case 'healthy':
      return t('monitor.dependenciesPage.statusHealthy');
    case 'abnormal':
      return t('monitor.dependenciesPage.statusAbnormal');
    case 'notConfigured':
      return t('monitor.dependenciesPage.statusNotConfigured');
    default:
      return t('monitor.dependenciesPage.statusUnknown');
  }
});

const summaryMetrics = computed(() => {
  const summary = serverStatus.value?.summary;

  return [
    {
      key: 'healthy',
      label: t('monitor.dependenciesPage.summary.healthy'),
      value: summary?.healthy_dependencies !== undefined ? String(summary.healthy_dependencies) : '--',
      description: t('monitor.dependenciesPage.summary.healthyDescription'),
    },
    {
      key: 'abnormal',
      label: t('monitor.dependenciesPage.summary.abnormal'),
      value: summary?.degraded_dependencies !== undefined ? String(summary.degraded_dependencies) : '--',
      description: t('monitor.dependenciesPage.summary.abnormalDescription'),
    },
    {
      key: 'notConfigured',
      label: t('monitor.dependenciesPage.summary.notConfigured'),
      value: summary?.disabled_dependencies !== undefined ? String(summary.disabled_dependencies) : '--',
      description: t('monitor.dependenciesPage.summary.notConfiguredDescription'),
    },
    {
      key: 'lastCheck',
      label: t('monitor.dependenciesPage.summary.lastCheck'),
      value: formatTimestamp(observedAt.value),
      description: t('monitor.dependenciesPage.summary.lastCheckDescription'),
    },
  ];
});

const serviceCards = computed(() => {
  const response = serverStatus.value;
  const observedLabel = formatTimestamp(response?.observed_at);
  const database = response?.dependencies.database;
  const redis = response?.dependencies.redis;

  return [
    buildServiceCard({
      key: 'postgresql',
      title: t('monitor.serverStatus.postgresqlLabel'),
      subtitle: t('monitor.dependenciesPage.postgresqlSubtitle'),
      status: normalizeDependencyStatus(database?.status),
      latency: database?.latency_ms,
      checkedAt: observedLabel,
      detail: database?.detail,
    }),
    buildServiceCard({
      key: 'redis',
      title: t('monitor.serverStatus.redisLabel'),
      subtitle: t('monitor.dependenciesPage.redisSubtitle'),
      status: normalizeDependencyStatus(redis?.status),
      latency: redis?.latency_ms,
      checkedAt: observedLabel,
      detail: redis?.detail,
    }),
    {
      key: 'future',
      title: t('monitor.dependenciesPage.futureEntryTitle'),
      subtitle: t('monitor.dependenciesPage.futureEntrySubtitle'),
      status: 'notConfigured',
      statusLabel: t('monitor.dependenciesPage.statusNotConfigured'),
      theme: dependencyStatusTheme('notConfigured'),
      futureEntry: true,
      fields: [
        {
          key: 'entry',
          label: t('monitor.dependenciesPage.fields.extensionEntry'),
          value: t('monitor.dependenciesPage.futureEntryLabel'),
          description: t('monitor.dependenciesPage.futureEntryHint'),
        },
      ],
    },
  ];
});

const overallDependencyStatus = computed(() => {
  const statuses = serviceCards.value
    .filter((service) => !service.futureEntry)
    .map((service) => service.status as ReturnType<typeof normalizeDependencyStatus> | undefined)
    .filter(Boolean);

  if (statuses.includes('abnormal')) {
    return 'abnormal';
  }

  if (statuses.includes('unknown')) {
    return 'unknown';
  }

  if (statuses.every((status) => status === 'notConfigured') && statuses.length > 0) {
    return 'notConfigured';
  }

  if (statuses.length > 0 && statuses.every((status) => status === 'healthy' || status === 'notConfigured')) {
    return 'healthy';
  }

  return 'unknown';
});

function buildServiceCard(options: {
  key: string;
  title: string;
  subtitle: string;
  status: ReturnType<typeof normalizeDependencyStatus>;
  latency?: number | null;
  checkedAt: string;
  detail?: string;
}): DependencyCard {
  return {
    key: options.key,
    title: options.title,
    subtitle: options.subtitle,
    status: options.status,
    statusLabel: dependencyStatusLabel(options.status),
    theme: dependencyStatusTheme(options.status),
    futureEntry: false,
    fields: [
      {
        key: 'latency',
        label: t('monitor.dependenciesPage.fields.latency'),
        value: formatLatency(options.latency),
        description: t('monitor.dependenciesPage.fieldDescriptions.latency'),
      },
      {
        key: 'checkedAt',
        label: t('monitor.dependenciesPage.fields.checkedAt'),
        value: options.checkedAt,
        description: t('monitor.dependenciesPage.fieldDescriptions.checkedAt'),
      },
      {
        key: 'errorInfo',
        label: t('monitor.dependenciesPage.fields.errorInfo'),
        value:
          options.status === 'abnormal' || options.status === 'unknown'
            ? displayText(options.detail)
            : t('monitor.dependenciesPage.noError'),
        description: t('monitor.dependenciesPage.fieldDescriptions.errorInfo'),
      },
      {
        key: 'detail',
        label: t('monitor.dependenciesPage.fields.detail'),
        value: displayText(options.detail),
        description: t('monitor.dependenciesPage.fieldDescriptions.detail'),
      },
    ],
  };
}

function dependencyStatusLabel(status: ReturnType<typeof normalizeDependencyStatus>) {
  switch (status) {
    case 'healthy':
      return t('monitor.dependenciesPage.statusHealthy');
    case 'abnormal':
      return t('monitor.dependenciesPage.statusAbnormal');
    case 'notConfigured':
      return t('monitor.dependenciesPage.statusNotConfigured');
    default:
      return t('monitor.dependenciesPage.statusUnknown');
  }
}
</script>
