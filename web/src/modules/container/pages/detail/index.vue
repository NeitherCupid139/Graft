<!--
  Copyright (c) 2025-2026 GeWuYou
  SPDX-License-Identifier: Apache-2.0
-->

<template>
  <div class="container-detail-page" data-page-type="operations-detail">
    <management-page-header
      :breadcrumb="detailBreadcrumb"
      :title="pageTitle"
      :description="detail ? detail.image : t('container.detail.description')"
      :source="{ labelKey: 'container.list.eyebrow', fallback: t('container.list.eyebrow') }"
    >
      <template #meta>
        <t-space class="container-detail-header-meta" break-line size="small">
          <span v-if="detail" class="container-detail-header-id">{{ shortContainerId(detail) }}</span>
          <t-tag v-if="detail" :theme="stateTheme(detail.state)" variant="light-outline">
            {{ stateLabel(detail.state) }}
          </t-tag>
          <t-tag v-if="detail" :theme="healthTheme(detail.health)" variant="light-outline">
            {{ healthLabel(detail.health) }}
          </t-tag>
          <t-tag v-if="detail?.runtime" theme="default" variant="light-outline">
            {{ detail.runtime }}
          </t-tag>
          <t-tag v-if="detail?.inspect_updated_at" theme="default" variant="light-outline">
            {{ t('container.detail.inspectUpdatedAt') }}: {{ formatTime(detail.inspect_updated_at) }}
          </t-tag>
        </t-space>
      </template>
      <template #actions>
        <t-space break-line size="small">
          <t-button theme="default" variant="outline" @click="goBack">
            {{ t('container.detail.back') }}
          </t-button>
          <t-button theme="primary" :loading="loading" @click="loadDetail">
            {{ t('container.detail.refresh') }}
          </t-button>
        </t-space>
      </template>
    </management-page-header>

    <t-alert v-if="error" theme="error" :title="error">
      <template #operation>
        <t-button theme="danger" variant="text" @click="loadDetail">
          {{ t('container.list.retry') }}
        </t-button>
      </template>
    </t-alert>

    <t-loading :loading="loading">
      <template v-if="detail">
        <section class="container-detail-summary">
          <t-card
            class="container-detail-summary-card container-detail-summary-card--identity"
            size="small"
            :bordered="false"
            :title="t('container.detail.summary.identity')"
          >
            <div class="container-detail-summary-list">
              <div class="container-detail-kv">
                <span>{{ t('container.list.fields.name') }}</span>
                <strong>{{ displayName(detail) }}</strong>
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.list.fields.image') }}</span>
                <copyable-detail-value
                  :copy-label="t('container.detail.copy')"
                  :value="detail.image"
                  :display-value="detail.image"
                  @copy="copyDetailText"
                />
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.list.fields.id') }}</span>
                <copyable-detail-value
                  :value="detail.id"
                  :display-value="shortContainerId(detail)"
                  :copy-label="t('container.detail.copy')"
                  code
                  data-testid="summary-container-id-copy"
                  @copy="copyDetailText"
                />
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.list.fields.runtime') }}</span>
                <strong>{{ runtimeLabel(detail) }}</strong>
              </div>
            </div>
          </t-card>
          <t-card
            class="container-detail-summary-card container-detail-summary-card--runtime"
            size="small"
            :bordered="false"
            :title="t('container.detail.summary.runtime')"
          >
            <div class="container-detail-summary-list">
              <div class="container-detail-kv container-detail-kv--inline">
                <span>{{ t('container.list.fields.status') }}</span>
                <t-tag :theme="stateTheme(detail.state)" variant="light-outline">
                  {{ stateLabel(detail.state) }}
                </t-tag>
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.list.fields.state') }}</span>
                <code>{{ detail.state || '-' }}</code>
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.list.fields.startedAt') }}</span>
                <strong>{{ formatTime(detail.started_at) }}</strong>
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.detail.health.status') }}</span>
                <strong>{{ healthLabel(detail.health) }}</strong>
              </div>
            </div>
          </t-card>
          <t-card
            class="container-detail-summary-card container-detail-summary-card--resources"
            size="small"
            :bordered="false"
            :title="t('container.detail.summary.resources')"
          >
            <div class="container-detail-summary__resource">
              <div class="container-detail-resource-meter container-detail-resource-meter--cpu">
                <div class="container-detail-resource-meter__content">
                  <span>{{ t('container.detail.resources.cpu') }}</span>
                  <strong>{{ formatPercent(detail.resource?.cpu_percent) }}</strong>
                </div>
                <t-progress
                  theme="circle"
                  size="small"
                  :label="formatPercent(detail.resource?.cpu_percent)"
                  :percentage="toProgressPercent(detail.resource?.cpu_percent)"
                />
              </div>
              <div class="container-detail-resource-meter container-detail-resource-meter--memory">
                <div class="container-detail-resource-meter__content">
                  <span>{{ t('container.detail.resources.memory') }}</span>
                  <strong>{{ memorySummary(detail) }}</strong>
                  <em>{{ formatPercent(detail.resource?.memory_percent) }}</em>
                </div>
                <t-progress
                  theme="line"
                  size="small"
                  :label="false"
                  :percentage="toProgressPercent(detail.resource?.memory_percent)"
                />
              </div>
            </div>
          </t-card>
          <t-card
            class="container-detail-summary-card container-detail-summary-card--network"
            size="small"
            :bordered="false"
            :title="t('container.detail.summary.network')"
          >
            <div class="container-detail-summary-list">
              <div class="container-detail-kv">
                <span>{{ t('container.detail.network.primaryIp') }}</span>
                <strong>{{ detail.primary_ip || '-' }}</strong>
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.detail.network.summary') }}</span>
                <strong>{{ networkSummary(detail) }}</strong>
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.detail.network.name') }}</span>
                <strong>{{ primaryNetworkName(detail) }}</strong>
              </div>
              <div class="container-detail-kv">
                <span>{{ t('container.detail.network.ports') }}</span>
                <div v-if="detail.ports.length" class="container-detail-port-list">
                  <t-tag
                    v-for="port in detail.ports"
                    :key="portLabel(port)"
                    class="container-detail-port-chip"
                    theme="default"
                    variant="light-outline"
                  >
                    {{ portLabel(port) }}
                  </t-tag>
                </div>
                <strong v-else>{{ t('container.detail.network.noPublicPorts') }}</strong>
              </div>
            </div>
          </t-card>
        </section>

        <t-card class="container-detail-tabs-card" :bordered="true">
          <t-tabs v-model:value="activeTab" theme="card" @change="handleTabChange">
            <t-tab-panel value="overview" :label="t('container.detail.tabs.overview')" :destroy-on-hide="false">
              <section class="container-detail-section container-detail-section--overview">
                <container-overview-panel
                  :copy-label="t('container.detail.copy')"
                  :sections="overviewSections"
                  @copy="copyDetailText"
                />
              </section>
            </t-tab-panel>

            <t-tab-panel value="resources" :label="t('container.detail.tabs.resources')" :destroy-on-hide="false">
              <section class="container-detail-section">
                <div class="container-detail-resource-grid">
                  <metric-card
                    :title="t('container.detail.resources.cpu')"
                    :value="resourceMetrics.cpu.value"
                    :description="resourceMetrics.cpu.description"
                    :progress="resourceMetrics.cpu.progress"
                    :progress-label="resourceMetrics.cpu.progressLabel"
                  />
                  <metric-card
                    :title="t('container.detail.resources.memory')"
                    :value="resourceMetrics.memory.value"
                    :description="resourceMetrics.memory.description"
                    :progress="resourceMetrics.memory.progress"
                    :progress-label="resourceMetrics.memory.progressLabel"
                  />
                  <article class="container-detail-resource-status-card">
                    <div class="container-detail-resource-status-card__content">
                      <span class="container-detail-resource-status-card__title">
                        {{ t('container.detail.resources.status') }}
                      </span>
                      <strong>{{ resourceMetrics.status.value }}</strong>
                      <span>{{ resourceMetrics.status.description }}</span>
                    </div>
                    <t-tag :theme="resourceMetrics.status.theme" variant="light-outline">
                      {{ resourceMetrics.status.value }}
                    </t-tag>
                  </article>
                </div>
                <section class="container-resource-detail-section">
                  <div class="container-resource-detail-section__title">
                    {{ t('container.detail.resources.detail') }}
                  </div>
                  <div class="container-resource-detail-section__body">
                    <div v-for="row in resourceDetailRows" :key="row.key" class="container-resource-detail-row">
                      <span class="container-resource-detail-row__label">{{ row.label }}</span>
                      <span class="container-resource-detail-row__value">
                        <t-tag v-if="row.type === 'tag'" :theme="row.theme" variant="light-outline">
                          {{ row.value }}
                        </t-tag>
                        <span v-else>{{ row.value }}</span>
                      </span>
                    </div>
                  </div>
                </section>
              </section>
            </t-tab-panel>

            <t-tab-panel value="logs" :label="t('container.detail.tabs.logs')" :destroy-on-hide="false">
              <section class="container-detail-section">
                <log-viewer
                  v-model:line-limit="logLineLimit"
                  :lines="logs?.lines ?? []"
                  :loading="logsLoading"
                  :error="logsError"
                  :truncated="logs?.truncated"
                  :refresh-label="t('container.detail.logs.refresh')"
                  :copy-label="t('container.detail.copy')"
                  :search-placeholder="t('container.detail.logs.searchPlaceholder')"
                  :wrap-label="t('container.detail.logs.wrap')"
                  :follow-tail-label="t('container.detail.logs.followTail')"
                  :empty-label="t('container.detail.logs.empty')"
                  :truncated-label="t('container.detail.logs.truncated')"
                  :copy-success-label="t('container.detail.copySuccess')"
                  :copy-error-label="t('container.detail.copyError')"
                  @refresh="loadLogs"
                />
              </section>
            </t-tab-panel>

            <t-tab-panel value="health" :label="t('container.detail.tabs.health')" :destroy-on-hide="false">
              <section class="container-detail-section">
                <t-descriptions :column="2" item-layout="vertical" bordered table-layout="fixed">
                  <t-descriptions-item :label="t('container.detail.health.status')">
                    {{ healthLabel(detail.health) }}
                  </t-descriptions-item>
                  <t-descriptions-item :label="t('container.detail.health.restartCount')">
                    {{ detail.restart_count ?? '-' }}
                  </t-descriptions-item>
                  <t-descriptions-item :label="t('container.list.fields.restartPolicy')">
                    {{ detail.restart_policy || '-' }}
                  </t-descriptions-item>
                  <t-descriptions-item :label="t('container.list.detail.inspectUpdatedAt')">
                    {{ formatTime(detail.inspect_updated_at) }}
                  </t-descriptions-item>
                </t-descriptions>
              </section>
            </t-tab-panel>

            <t-tab-panel value="config" :label="t('container.detail.tabs.config')" :destroy-on-hide="false">
              <section class="container-detail-section">
                <t-descriptions :column="2" item-layout="vertical" bordered table-layout="fixed">
                  <t-descriptions-item :label="t('container.list.detail.command')">
                    {{ joinList(detail.command) }}
                  </t-descriptions-item>
                  <t-descriptions-item :label="t('container.list.detail.entrypoint')">
                    {{ joinList(detail.entrypoint) }}
                  </t-descriptions-item>
                  <t-descriptions-item :label="t('container.list.detail.workingDir')">
                    {{ detail.working_dir || '-' }}
                  </t-descriptions-item>
                </t-descriptions>
                <div class="container-detail-subsection">
                  <h3>{{ t('container.detail.config.environment') }}</h3>
                  <t-table
                    v-if="environmentRows.length"
                    row-key="name"
                    size="small"
                    :columns="environmentColumns"
                    :data="environmentRows"
                    :pagination="undefined"
                    table-layout="fixed"
                    cell-empty-content="-"
                  >
                    <template #value="{ row }">
                      <span>{{ row.value || '-' }}</span>
                    </template>
                    <template #policy="{ row }">
                      <t-tag :theme="policyTheme(row.policy)" variant="light-outline">
                        {{ policyLabel(row.policy) }}
                      </t-tag>
                    </template>
                    <template #operation="{ row }">
                      <t-button
                        v-if="row.copyable"
                        data-testid="env-copy"
                        size="small"
                        theme="default"
                        variant="text"
                        @click="copyEnvironmentValue(row)"
                      >
                        {{ t('container.detail.copy') }}
                      </t-button>
                    </template>
                  </t-table>
                  <t-empty v-else size="small" :description="t('container.detail.config.environmentUnavailable')" />
                </div>
              </section>
            </t-tab-panel>

            <t-tab-panel value="network" :label="t('container.detail.tabs.network')" :destroy-on-hide="false">
              <section class="container-detail-section">
                <t-table
                  v-if="detail.networks.length"
                  row-key="name"
                  size="small"
                  :columns="networkColumns"
                  :data="detail.networks"
                  :pagination="undefined"
                  table-layout="fixed"
                  cell-empty-content="-"
                />
                <t-empty v-else size="small" :description="t('container.list.detail.networkEmpty')" />
              </section>
            </t-tab-panel>

            <t-tab-panel value="storage" :label="t('container.detail.tabs.storage')" :destroy-on-hide="false">
              <section class="container-detail-section">
                <t-table
                  v-if="detail.mounts.length"
                  row-key="destination"
                  size="small"
                  :columns="mountColumns"
                  :data="detail.mounts"
                  :pagination="undefined"
                  table-layout="fixed"
                  cell-empty-content="-"
                >
                  <template #read_only="{ row }">
                    {{ row.read_only ? 'ro' : 'rw' }}
                  </template>
                </t-table>
                <t-empty v-else size="small" :description="t('container.list.detail.mountEmpty')" />
              </section>
            </t-tab-panel>

            <t-tab-panel value="raw" :label="t('container.detail.tabs.raw')" :destroy-on-hide="false">
              <section class="container-detail-section">
                <json-viewer
                  :value="detail"
                  :title="t('container.detail.raw.title')"
                  :description="t('container.detail.raw.description')"
                  :root-label="t('container.detail.raw.root')"
                  :source-label="t('container.detail.raw.source')"
                  :tree-label="t('container.detail.raw.tree')"
                  :copy-label="t('container.detail.copy')"
                  :copy-success-label="t('container.detail.copySuccess')"
                  :copy-error-label="t('container.detail.copyError')"
                  :empty-label="t('container.detail.raw.empty')"
                  :error-label="t('container.detail.raw.error')"
                />
              </section>
            </t-tab-panel>
          </t-tabs>
        </t-card>
      </template>

      <t-empty v-else-if="!error" size="small" :description="t('container.detail.empty')" />
    </t-loading>
  </div>
</template>
<script setup lang="ts">
import type { TableProps } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next/es/message';
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import { ManagementPageHeader } from '@/shared/components/management';
import { MetricCard } from '@/shared/components/metrics';
import { resolveLocalizedErrorMessage } from '@/shared/localized-api-error';
import {
  copyText as copyTextToClipboard,
  formatBytes,
  formatLocaleDateTime,
  formatPercent,
  JsonViewer,
  LogViewer,
  toProgressPercent,
} from '@/shared/observability';
import { createLogger } from '@/utils/logger';

import { getContainer, getContainerLogs } from '../../api/container';
import { CONTAINER_BOOTSTRAP_ROUTE } from '../../contract/bootstrap';
import type { ContainerDetail, ContainerHealth, ContainerLogResponse, ContainerState } from '../../types/container';
import ContainerOverviewPanel from './components/ContainerOverviewPanel.vue';
import CopyableDetailValue from './components/CopyableDetailValue.vue';
import type { ContainerOverviewInfoSection } from './components/overview';

defineOptions({
  name: 'ContainerDetailIndex',
});

type DetailTab = 'overview' | 'resources' | 'logs' | 'health' | 'config' | 'network' | 'storage' | 'raw';
type EnvironmentPolicy = 'plain' | 'masked' | 'hidden' | 'unknown';
type EnvironmentRow = {
  copyable: boolean;
  name: string;
  policy: EnvironmentPolicy;
  rawValue: string;
  value: string;
};
type ResourceStatusTheme = 'success' | 'warning' | 'default';
type ResourceDetailRow =
  | {
      key: string;
      label: string;
      type: 'text';
      value: string;
    }
  | {
      key: string;
      label: string;
      theme: ResourceStatusTheme;
      type: 'tag';
      value: string;
    };

const DETAIL_TABS: DetailTab[] = ['overview', 'resources', 'logs', 'health', 'config', 'network', 'storage', 'raw'];
const DEFAULT_LOG_QUERY = {
  tail: 200,
  since: undefined,
  timestamps: false,
  stdout: true,
  stderr: true,
};

const { locale, t } = useI18n();
const route = useRoute();
const router = useRouter();
const logger = createLogger('container.detail');

const detail = ref<ContainerDetail | null>(null);
const loading = ref(false);
const error = ref('');
const logs = ref<ContainerLogResponse | null>(null);
const logsLoading = ref(false);
const logsError = ref('');
const logLineLimit = ref(DEFAULT_LOG_QUERY.tail);
const activeTab = ref<DetailTab>(normalizeTab(route.query.tab));

const containerId = computed(() => String(route.params.id ?? '').trim());
const detailBreadcrumb = computed(() => [
  { labelKey: 'container.list.eyebrow', fallback: t('container.list.eyebrow') },
  { labelKey: 'container.detail.title', fallback: t('container.detail.title') },
]);
const pageTitle = computed(() => {
  if (detail.value) {
    return displayName(detail.value);
  }
  return containerId.value || t('container.detail.title');
});
const environmentRows = computed(() => normalizeEnvironmentRows(detail.value));
const resourceMetrics = computed(() => {
  const current = detail.value;
  const resource = current?.resource;
  const cpuValue = formatPercent(resource?.cpu_percent);
  const memoryPercent = formatPercent(resource?.memory_percent);
  const status = current ? resourceStatus(current) : emptyResourceStatus();

  return {
    cpu: {
      description:
        cpuValue === '-' ? t('container.detail.resources.noData') : t('container.detail.resources.currentSnapshot'),
      progress: toProgressPercent(resource?.cpu_percent),
      progressLabel: cpuValue,
      value: cpuValue,
    },
    memory: {
      description: memoryPercent,
      progress: toProgressPercent(resource?.memory_percent),
      progressLabel: memoryPercent,
      value: current ? memorySummary(current) : '-',
    },
    status,
  };
});
const resourceDetailRows = computed<ResourceDetailRow[]>(() => {
  const current = detail.value;
  if (!current) {
    return [];
  }
  const status = resourceStatus(current);
  const resource = current.resource;

  return [
    {
      key: 'cpu',
      label: t('container.detail.resources.cpu'),
      type: 'text',
      value: formatPercent(resource?.cpu_percent),
    },
    {
      key: 'memory-usage',
      label: t('container.detail.resources.memoryUsage'),
      type: 'text',
      value: formatBytes(resource?.memory_usage_bytes),
    },
    {
      key: 'memory-limit',
      label: t('container.detail.resources.memoryLimit'),
      type: 'text',
      value: formatBytes(resource?.memory_limit_bytes),
    },
    {
      key: 'memory-percent',
      label: t('container.detail.resources.memoryPercent'),
      type: 'text',
      value: formatPercent(resource?.memory_percent),
    },
    {
      key: 'status',
      label: t('container.detail.resources.status'),
      theme: status.theme,
      type: 'tag',
      value: status.value,
    },
    {
      key: 'collected-at',
      label: t('container.detail.resources.collectedAt'),
      type: 'text',
      value: status.collectedAt,
    },
  ];
});
const overviewSections = computed<ContainerOverviewInfoSection[]>(() => {
  const current = detail.value;
  if (!current) {
    return [];
  }

  const imageId = readableImageId(current.image_id);

  return [
    {
      key: 'basic',
      title: t('container.detail.overview.basicInfo'),
      rows: [
        {
          displayValue: displayName(current),
          key: 'name',
          label: t('container.detail.overview.fields.name'),
          type: 'text',
        },
        {
          code: true,
          copyValue: current.id,
          displayValue: shortContainerId(current),
          key: 'container-id',
          label: t('container.detail.overview.fields.containerId'),
          testId: 'container-id-copy',
          type: 'copy',
        },
        {
          copyValue: current.image,
          displayValue: current.image || '-',
          key: 'image',
          label: t('container.detail.overview.fields.image'),
          type: 'copy',
        },
        {
          code: true,
          copyValue: imageId,
          displayValue: shortIdentifier(imageId),
          key: 'image-id',
          label: t('container.detail.overview.fields.imageId'),
          testId: 'image-id-copy',
          type: 'copy',
        },
        {
          displayValue: runtimeLabel(current),
          key: 'runtime',
          label: t('container.detail.overview.fields.runtime'),
          type: 'text',
        },
      ],
    },
    {
      key: 'runtime',
      title: t('container.detail.overview.runtimeInfo'),
      rows: [
        {
          key: 'status',
          label: t('container.detail.overview.fields.status'),
          tagLabel: stateLabel(current.state),
          tagTheme: stateTheme(current.state),
          type: 'tag',
        },
        {
          displayValue: current.state || '-',
          key: 'state',
          label: t('container.detail.overview.fields.state'),
          type: 'text',
        },
        {
          key: 'health',
          label: t('container.detail.overview.fields.health'),
          tagLabel: healthLabel(current.health),
          tagTheme: healthTheme(current.health),
          type: 'tag',
        },
        {
          displayValue: formatTime(current.created_at),
          key: 'created-at',
          label: t('container.detail.overview.fields.createdAt'),
          type: 'text',
        },
        {
          displayValue: formatTime(current.started_at),
          key: 'started-at',
          label: t('container.detail.overview.fields.startedAt'),
          type: 'text',
        },
        {
          displayValue: formatTime(current.inspect_updated_at),
          key: 'updated-at',
          label: t('container.detail.overview.fields.updatedAt'),
          type: 'text',
        },
      ],
    },
    {
      key: 'resource-network',
      title: t('container.detail.overview.resourceNetwork'),
      rows: [
        {
          displayValue: formatPercent(current.resource?.cpu_percent),
          key: 'cpu',
          label: t('container.detail.resources.cpu'),
          type: 'text',
        },
        {
          displayValue: memorySummary(current),
          key: 'memory',
          label: t('container.detail.resources.memory'),
          type: 'text',
        },
        {
          displayValue: current.primary_ip || '-',
          key: 'primary-ip',
          label: t('container.detail.network.primaryIp'),
          type: 'text',
        },
        {
          displayValue: networkSummary(current),
          key: 'network-mode',
          label: t('container.detail.overview.fields.networkMode'),
          type: 'text',
        },
        {
          displayValue: primaryNetworkName(current),
          key: 'network-name',
          label: t('container.detail.overview.fields.networkName'),
          type: 'text',
        },
        {
          emptyLabel: t('container.detail.network.noPublicPorts'),
          key: 'ports',
          label: t('container.detail.network.ports'),
          ports: current.ports.map((port) => portLabel(port)),
          type: 'ports',
        },
      ],
    },
  ];
});
const environmentColumns = computed<TableProps['columns']>(() => [
  { colKey: 'name', title: t('container.detail.config.envName'), minWidth: 220, ellipsis: true },
  { colKey: 'value', title: t('container.detail.config.envValue'), minWidth: 260, ellipsis: true },
  { colKey: 'policy', title: t('container.detail.config.envPolicy'), width: 160, align: 'center' },
  { colKey: 'operation', title: t('container.detail.operation'), width: 112, align: 'center' },
]);
const networkColumns = computed<TableProps['columns']>(() => [
  { colKey: 'name', title: t('container.detail.network.name'), minWidth: 180, ellipsis: true },
  { colKey: 'ip_address', title: t('container.detail.network.ipAddress'), minWidth: 160, ellipsis: true },
  { colKey: 'gateway', title: t('container.detail.network.gateway'), minWidth: 160, ellipsis: true },
  { colKey: 'mac_address', title: t('container.detail.network.macAddress'), minWidth: 180, ellipsis: true },
]);
const mountColumns = computed<TableProps['columns']>(() => [
  { colKey: 'destination', title: t('container.detail.storage.destination'), minWidth: 240, ellipsis: true },
  { colKey: 'source', title: t('container.detail.storage.source'), minWidth: 260, ellipsis: true },
  { colKey: 'type', title: t('container.detail.storage.type'), width: 120, align: 'center' },
  { colKey: 'mode', title: t('container.detail.storage.mode'), width: 120, align: 'center' },
  { colKey: 'read_only', title: t('container.detail.storage.access'), width: 120, align: 'center' },
]);

onMounted(() => {
  void loadDetail();
  if (activeTab.value === 'logs') {
    void loadLogs();
  }
});

watch(
  () => route.params.id,
  () => {
    resetDetailState();
    void loadDetail();
    if (activeTab.value === 'logs') {
      void loadLogs();
    }
  },
);

watch(
  () => route.query.tab,
  (tab) => {
    const normalized = normalizeTab(tab);
    activeTab.value = normalized;
    if (normalized === 'logs' && !logs.value) {
      void loadLogs();
    }
  },
);

watch(logLineLimit, () => {
  if (activeTab.value === 'logs') {
    void loadLogs();
  }
});

async function loadDetail() {
  if (!containerId.value) {
    detail.value = null;
    logs.value = null;
    error.value = t('container.detail.missingId');
    return;
  }

  loading.value = true;
  error.value = '';
  try {
    detail.value = await getContainer(containerId.value);
  } catch (loadError) {
    detail.value = null;
    error.value = resolveLocalizedErrorMessage(t, loadError, t('container.list.detail.loadFailed'));
    logger.warn('failed to fetch container detail', loadError);
  } finally {
    loading.value = false;
  }
}

async function loadLogs() {
  if (!containerId.value) {
    logs.value = null;
    return;
  }
  logsLoading.value = true;
  logsError.value = '';
  try {
    logs.value = await getContainerLogs(containerId.value, {
      ...DEFAULT_LOG_QUERY,
      tail: logLineLimit.value,
    });
  } catch (loadError) {
    logsError.value = resolveLocalizedErrorMessage(t, loadError, t('container.list.logs.loadFailed'));
    logger.warn('failed to fetch container logs', loadError);
  } finally {
    logsLoading.value = false;
  }
}

function resetDetailState() {
  detail.value = null;
  error.value = '';
  logs.value = null;
  logsError.value = '';
}

function handleTabChange(value: string | number) {
  const tab = normalizeTab(value);
  activeTab.value = tab;
  void router.replace({
    params: route.params,
    query: {
      ...route.query,
      tab,
    },
  });
  if (tab === 'logs' && !logs.value) {
    void loadLogs();
  }
}

function goBack() {
  if (window.history.length > 1) {
    router.back();
    return;
  }
  void router.push({ name: CONTAINER_BOOTSTRAP_ROUTE.LIST.routeName });
}

async function copyEnvironmentValue(row: EnvironmentRow) {
  await copyDetailText(row.rawValue);
}

async function copyDetailText(text: string) {
  if (!text) return;
  const copied = await copyTextToClipboard(text);
  if (copied) {
    MessagePlugin.success(t('container.detail.copySuccess'));
    return;
  }
  MessagePlugin.error(t('container.detail.copyError'));
}

function normalizeTab(value: unknown): DetailTab {
  const raw = Array.isArray(value) ? value[0] : value;
  return typeof raw === 'string' && DETAIL_TABS.includes(raw as DetailTab) ? (raw as DetailTab) : 'overview';
}

function normalizeEnvironmentRows(nextDetail: ContainerDetail | null): EnvironmentRow[] {
  const detailRecord = readUnknownRecord(nextDetail);
  const source = detailRecord?.environment;
  if (!Array.isArray(source)) {
    return [];
  }

  return source.flatMap((item) => {
    const record = readUnknownRecord(item);
    const name = readString(record?.name ?? record?.key);
    if (!name) {
      return [];
    }

    const rawPolicy = readString(record?.policy ?? record?.visibility ?? record?.state);
    const masked = record?.masked === true;
    const rawValue = readRawString(record?.value);
    const policy = normalizeEnvironmentPolicy(
      rawPolicy,
      readString(detailRecord?.environment_policy),
      masked,
      rawValue,
    );
    const value = policy === 'hidden' ? '' : rawValue;

    return [
      {
        copyable: Boolean(rawValue) && policy !== 'hidden',
        name,
        policy,
        rawValue,
        value: value || environmentValueFallback(policy),
      },
    ];
  });
}

function readUnknownRecord(value: unknown): Record<string, unknown> | null {
  return value && typeof value === 'object' && !Array.isArray(value) ? (value as Record<string, unknown>) : null;
}

function readString(value: unknown) {
  return typeof value === 'string' ? value.trim() : '';
}

function readRawString(value: unknown) {
  return typeof value === 'string' ? value : '';
}

function normalizeEnvironmentPolicy(
  value: string,
  detailPolicy = '',
  masked = false,
  rawValue = '',
): EnvironmentPolicy {
  if (value === 'plain' || value === 'masked' || value === 'hidden') {
    return value;
  }
  if (rawValue && !masked) {
    return 'plain';
  }
  if (masked) {
    return 'masked';
  }
  if (detailPolicy === 'plain' || detailPolicy === 'masked' || detailPolicy === 'hidden') {
    return detailPolicy;
  }
  return 'unknown';
}

function environmentValueFallback(policy: EnvironmentPolicy) {
  if (policy === 'masked') return t('container.detail.config.maskedValue');
  if (policy === 'hidden') return t('container.detail.config.hiddenValue');
  return '-';
}

function policyLabel(policy: EnvironmentPolicy) {
  return t(`container.detail.config.policy.${policy}`);
}

function policyTheme(policy: EnvironmentPolicy) {
  if (policy === 'plain') return 'success';
  if (policy === 'masked') return 'warning';
  if (policy === 'hidden') return 'danger';
  return 'default';
}

function displayName(row: ContainerDetail) {
  return row.name || row.names[0] || shortContainerId(row);
}

function stateLabel(state: ContainerState) {
  return t(`container.list.states.${state}`);
}

function healthLabel(health?: ContainerHealth | null) {
  return t(`container.list.health.${health || 'none'}`);
}

function healthTheme(health?: ContainerHealth | null) {
  if (health === 'healthy') return 'success';
  if (health === 'unhealthy') return 'danger';
  if (health === 'starting') return 'warning';
  return 'default';
}

function stateTheme(state: ContainerState) {
  if (state === 'running') return 'success';
  if (state === 'created' || state === 'paused' || state === 'restarting') return 'warning';
  if (state === 'dead') return 'danger';
  return 'default';
}

function formatTime(value?: string | null) {
  return formatLocaleDateTime(value, locale);
}

function joinList(values?: string[]) {
  return values?.length ? values.join(' ') : '-';
}

function runtimeLabel(nextDetail: ContainerDetail) {
  return nextDetail.runtime || nextDetail.runtime_info?.runtime || '-';
}

function shortContainerId(nextDetail: ContainerDetail) {
  return shortIdentifier(nextDetail.id, nextDetail.short_id, 12);
}

function networkSummary(nextDetail: ContainerDetail) {
  return nextDetail.network_summary || primaryNetworkName(nextDetail);
}

function primaryNetworkName(nextDetail: ContainerDetail) {
  return nextDetail.networks[0]?.name || '-';
}

function resourceStatus(nextDetail: ContainerDetail) {
  const resource = nextDetail.resource;
  if (resource?.stats_available || resource?.available) {
    return {
      collectedAt: formatTime(nextDetail.inspect_updated_at),
      description: formatTime(nextDetail.inspect_updated_at),
      theme: 'success' as const,
      value: t('container.detail.resources.available'),
    };
  }
  if (resource?.stats_error_message || resource?.stats_error_key || resource?.unavailable_reason) {
    return {
      collectedAt: '-',
      description: resource.stats_error_message || resource.stats_error_key || resource.unavailable_reason || '-',
      theme: 'warning' as const,
      value: t('container.detail.resources.unavailable'),
    };
  }
  return emptyResourceStatus();
}

function emptyResourceStatus() {
  return {
    collectedAt: '-',
    description: '-',
    theme: 'default' as const,
    value: t('container.detail.resources.noData'),
  };
}

function memorySummary(nextDetail: ContainerDetail) {
  const resource = nextDetail.resource;
  const usage = formatBytes(resource?.memory_usage_bytes);
  const limit = formatBytes(resource?.memory_limit_bytes);
  if (usage === '-' && limit === '-') {
    return '-';
  }
  return `${usage} / ${limit}`;
}

function shortIdentifier(value?: string | null, preferred?: string | null, maxLength = 28) {
  const normalized = value?.trim();
  if (!normalized) return '-';
  if (preferred?.trim()) return preferred.trim();
  if (normalized.length <= maxLength) return normalized;
  if (maxLength <= 12) return normalized.slice(0, maxLength);
  return `${normalized.slice(0, 18)}...${normalized.slice(-10)}`;
}

function readableImageId(value?: string | null) {
  const normalized = value?.trim();
  if (!normalized) return '-';
  return normalized.startsWith('sha256:') ? normalized.slice('sha256:'.length) : normalized;
}

function portLabel(port: ContainerDetail['ports'][number]) {
  const privatePort = port.private_port ? `${port.private_port}` : '-';
  const publicPort = port.public_port ? `${port.public_port}` : '';
  return publicPort ? `${publicPort}:${privatePort}/${port.type}` : `${privatePort}/${port.type}`;
}
</script>
<style scoped lang="less">
.container-detail-page {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-16);
  min-width: 0;
}

.container-detail-summary {
  display: grid;
  gap: var(--graft-density-gap-12);
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.container-detail-summary-card {
  background: color-mix(in srgb, var(--td-bg-color-container) 92%, transparent);
  border: 1px solid color-mix(in srgb, var(--td-component-stroke) 64%, transparent);
  height: 100%;
  min-width: 0;
}

.container-detail-summary-card :deep(.t-card__body) {
  height: calc(100% - var(--td-comp-size-xxxl));
}

.container-detail-metric,
.container-detail-section,
.container-detail-subsection {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-8);
  min-width: 0;
}

.container-detail-summary-list {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-10);
  height: 100%;
  min-width: 0;
}

.container-detail-kv {
  display: grid;
  gap: var(--graft-density-gap-4);
  min-width: 0;
}

.container-detail-kv--inline {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
}

.container-detail-kv > span {
  color: var(--td-text-color-placeholder);
  font: var(--td-font-body-small);
}

.container-detail-kv strong,
.container-detail-kv code {
  color: var(--td-text-color-primary);
  font: var(--td-font-body-medium);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.container-detail-kv code,
.container-detail-header-id {
  font-family: var(
    --td-font-family-mono,
    ui-monospace,
    SFMono-Regular,
    Menlo,
    Monaco,
    Consolas,
    'Liberation Mono',
    monospace
  );
}

.container-detail-header-id {
  align-items: center;
  color: var(--td-text-color-secondary);
  display: inline-flex;
  font: var(--td-font-body-small);
  min-height: var(--td-comp-size-xs);
}

.container-detail-header-meta {
  max-width: min(100%, 680px);
  min-width: 0;
}

.container-detail-header-meta :deep(.t-space-item) {
  min-width: 0;
}

.container-detail-header-meta :deep(.t-tag) {
  max-width: 100%;
}

.container-detail-summary__resource,
.container-detail-resource-grid {
  display: grid;
  gap: var(--graft-density-gap-10);
  grid-template-columns: repeat(2, minmax(0, 1fr));
  min-width: 0;
}

.container-detail-resource-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.container-detail-summary__resource {
  grid-template-columns: 1fr;
}

.container-detail-metric strong,
.container-detail-subsection h3,
.container-detail-overview-group h3 {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-small);
  margin: 0;
}

.container-detail-metric span {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
  overflow-wrap: anywhere;
}

.container-detail-tag-row,
.container-detail-port-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--graft-density-gap-6);
  min-width: 0;
}

.container-detail-resource-meter {
  border-radius: var(--td-radius-medium);
  display: flex;
  gap: var(--graft-density-gap-10);
  min-width: 0;
}

.container-detail-resource-meter--cpu {
  align-items: center;
  justify-content: space-between;
}

.container-detail-resource-meter--memory {
  flex-direction: column;
  justify-content: center;
}

.container-detail-resource-meter__content {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: var(--graft-density-gap-4);
  min-width: 0;
}

.container-detail-resource-meter__content span {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
}

.container-detail-resource-meter__content strong {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-small);
  min-width: 0;
  overflow-wrap: anywhere;
}

.container-detail-resource-meter__content em {
  color: var(--td-text-color-placeholder);
  font: var(--td-font-body-small);
  font-style: normal;
}

.container-detail-resource-meter--memory :deep(.t-progress) {
  width: 100%;
}

.container-detail-port-chip {
  font-family: var(
    --td-font-family-mono,
    ui-monospace,
    SFMono-Regular,
    Menlo,
    Monaco,
    Consolas,
    'Liberation Mono',
    monospace
  );
  max-width: 100%;
}

.container-detail-tabs-card {
  min-width: 0;
}

.container-detail-tabs-card :deep(.t-card__body) {
  padding: 0;
}

.container-detail-tabs-card :deep(.t-tabs__content) {
  padding-top: var(--graft-density-gap-12);
}

/*
 * Short detail tabs use the page scrollbar. Long-form tabs such as logs and raw JSON own
 * their internal scrolling so the page does not fight a second nested scrollbar.
 */
.container-detail-section {
  padding: 0;
}

.container-detail-section--overview {
  min-height: 0;
  padding: 0 var(--graft-density-gap-16) var(--graft-density-gap-16);
}

.container-detail-resource-status-card {
  align-items: center;
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-medium);
  display: flex;
  gap: var(--graft-density-gap-12);
  justify-content: space-between;
  min-width: 0;
  padding: var(--graft-density-gap-14);
}

.container-detail-resource-status-card__content {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: var(--graft-density-gap-6);
  min-width: 0;
}

.container-detail-resource-status-card__content strong {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-medium);
  min-width: 0;
  overflow-wrap: anywhere;
}

.container-detail-resource-status-card__content span {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
  min-width: 0;
  overflow-wrap: anywhere;
}

.container-detail-resource-status-card__title {
  color: var(--td-text-color-secondary);
}

.container-resource-detail-section {
  background: var(--td-bg-color-container);
  border: 1px solid color-mix(in srgb, var(--td-component-stroke) 72%, transparent);
  border-radius: var(--td-radius-medium);
  display: flex;
  flex-direction: column;
  margin-top: var(--graft-density-gap-12);
  min-width: 0;
  overflow: hidden;
  width: 100%;
}

.container-resource-detail-section__title {
  background: color-mix(in srgb, var(--td-bg-color-container) 86%, var(--td-bg-color-page));
  border-bottom: 1px solid color-mix(in srgb, var(--td-component-stroke) 72%, transparent);
  color: var(--td-text-color-primary);
  font: var(--td-font-title-small);
  font-weight: 600;
  line-height: 22px;
  padding: var(--graft-density-gap-12) var(--graft-density-gap-16);
}

.container-resource-detail-section__body {
  display: flex;
  flex-direction: column;
  min-width: 0;
  padding: var(--graft-density-gap-8) var(--graft-density-gap-16);
}

.container-resource-detail-row {
  align-items: center;
  column-gap: var(--graft-density-gap-16);
  display: grid;
  grid-template-columns: 112px minmax(0, 1fr);
  min-height: 36px;
  min-width: 0;
  width: 100%;
}

.container-resource-detail-row + .container-resource-detail-row {
  border-top: 1px solid color-mix(in srgb, var(--td-component-stroke) 30%, transparent);
}

.container-resource-detail-row__label {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
  line-height: 20px;
  min-width: 0;
}

.container-resource-detail-row__value {
  align-items: center;
  color: var(--td-text-color-primary);
  display: inline-flex;
  font: var(--td-font-body-small);
  font-weight: 500;
  gap: var(--graft-density-gap-6);
  line-height: 22px;
  min-width: 0;
  overflow: hidden;
}

.container-resource-detail-row__value > span:not(.t-tag) {
  display: inline-block;
  max-width: 100%;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (width <= 1360px) {
  .container-detail-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (width <= 960px) {
  .container-detail-resource-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (width <= 720px) {
  .container-detail-summary,
  .container-detail-summary__resource {
    grid-template-columns: 1fr;
  }

  .container-detail-resource-grid {
    grid-template-columns: 1fr;
  }

  .container-resource-detail-row {
    align-items: flex-start;
    gap: var(--graft-density-gap-4);
    grid-template-columns: 1fr;
    padding: var(--graft-density-gap-8) 0;
  }

  .container-resource-detail-row__value {
    width: 100%;
  }

  .container-detail-header-meta {
    max-width: 100%;
  }

  .container-detail-header-meta :deep(.t-space) {
    width: 100%;
  }

  .container-detail-header-meta :deep(.t-space-item) {
    max-width: 100%;
  }
}
</style>
