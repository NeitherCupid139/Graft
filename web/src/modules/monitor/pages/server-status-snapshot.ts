import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { getServerStatus } from '../api/server-status';
import { MONITOR_TREND_RANGE } from '../contract/trend';
import type { ServerStatusResponse } from '../types/server-status';

type TagTheme = 'success' | 'warning' | 'danger' | 'default';

export type DependencyDisplayStatus = 'healthy' | 'abnormal' | 'notConfigured' | 'unknown';

export function useServerStatusSnapshot() {
  const { t } = useI18n();

  const loading = ref(false);
  const initialized = ref(false);
  const errorMessage = ref('');
  const serverStatus = ref<ServerStatusResponse | null>(null);

  async function refreshSnapshot() {
    loading.value = true;
    errorMessage.value = '';

    try {
      serverStatus.value = await getServerStatus(MONITOR_TREND_RANGE.TEN_MINUTES);
    } catch (error) {
      errorMessage.value = resolveErrorMessage(error, t('monitor.shared.loadFailed'));
    } finally {
      loading.value = false;
      initialized.value = true;
    }
  }

  onMounted(() => {
    void refreshSnapshot();
  });

  return {
    loading,
    initialized,
    errorMessage,
    serverStatus,
    refreshSnapshot,
    observedAt: computed(() => serverStatus.value?.observed_at ?? ''),
  };
}

export function normalizeDependencyStatus(status?: string): DependencyDisplayStatus {
  switch ((status ?? '').trim().toLowerCase()) {
    case 'healthy':
      return 'healthy';
    case 'degraded':
      return 'abnormal';
    case 'disabled':
      return 'notConfigured';
    default:
      return 'unknown';
  }
}

export function dependencyStatusTheme(status: DependencyDisplayStatus): TagTheme {
  switch (status) {
    case 'healthy':
      return 'success';
    case 'abnormal':
      return 'danger';
    case 'notConfigured':
      return 'default';
    default:
      return 'warning';
  }
}

export function runtimeSnapshotTheme(status?: string): TagTheme {
  switch ((status ?? '').trim().toLowerCase()) {
    case 'healthy':
      return 'success';
    case 'degraded':
      return 'warning';
    case 'disabled':
      return 'default';
    default:
      return 'warning';
  }
}

export function formatBytes(bytes?: number | null) {
  if (!Number.isFinite(bytes) || bytes === null || bytes === undefined || bytes < 0) {
    return '--';
  }

  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let value = bytes;
  let unitIndex = 0;

  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex += 1;
  }

  const digits = value >= 100 || unitIndex === 0 ? 0 : value >= 10 ? 1 : 2;
  return `${value.toFixed(digits)} ${units[unitIndex]}`;
}

export function formatTimestamp(value?: string | null) {
  if (!value) {
    return '--';
  }

  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return '--';
  }

  return parsed.toLocaleString();
}

export function formatUptime(totalSeconds?: number | null) {
  if (!Number.isFinite(totalSeconds) || totalSeconds === null || totalSeconds === undefined || totalSeconds < 0) {
    return '--';
  }

  const remainingSeconds = Math.floor(totalSeconds);
  const days = Math.floor(remainingSeconds / 86400);
  const hours = Math.floor((remainingSeconds % 86400) / 3600);
  const minutes = Math.floor((remainingSeconds % 3600) / 60);
  const seconds = remainingSeconds % 60;
  const parts = [
    days > 0 ? `${days}d` : '',
    hours > 0 ? `${hours}h` : '',
    minutes > 0 ? `${minutes}m` : '',
    seconds > 0 || (days === 0 && hours === 0 && minutes === 0) ? `${seconds}s` : '',
  ].filter(Boolean);

  return parts.join(' ');
}

export function formatLatency(latencyMs?: number | null) {
  if (!Number.isFinite(latencyMs) || latencyMs === null || latencyMs === undefined) {
    return '--';
  }

  return `${latencyMs.toFixed(2)} ms`;
}

export function displayText(value?: string | null) {
  if (!value) {
    return '--';
  }

  const trimmed = value.trim();
  return trimmed.length > 0 ? trimmed : '--';
}

function resolveErrorMessage(error: unknown, fallbackMessage: string) {
  if (error instanceof Error) {
    const trimmed = error.message.trim();
    if (trimmed) {
      return trimmed;
    }
  }

  return fallbackMessage;
}
