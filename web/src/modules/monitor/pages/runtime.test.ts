import { flushPromises, mount } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent, h } from 'vue';

import RuntimePage from './runtime.vue';

const monitorApiMocks = vi.hoisted(() => ({
  getServerStatus: vi.fn(),
}));

const translations = vi.hoisted(
  (): Record<string, string> => ({
    'monitor.sectionTitle': 'Server Management',
    'monitor.shared.refresh': 'Refresh now',
    'monitor.shared.loadFailed': 'Failed to load server status',
    'monitor.shared.empty': 'No server-status snapshot is available yet',
    'monitor.shared.errorTitle': 'Snapshot request failed',
    'monitor.shared.notReported': 'Not reported',
    'monitor.runtimePage.title': 'Runtime',
    'monitor.runtimePage.subtitle':
      'Inspect the current Go process runtime snapshot and host environment. Host memory and Go Runtime memory stay explicitly separated here.',
    'monitor.runtimePage.snapshotReady': 'Snapshot ready',
    'monitor.runtimePage.snapshotPending': 'Waiting for snapshot',
    'monitor.runtimePage.memoryBoundaryTitle': 'Memory ownership boundary',
    'monitor.runtimePage.memoryBoundaryDescription':
      'Host memory reflects the full machine snapshot. Go Runtime memory only reflects allocations inside the current Go process.',
    'monitor.runtimePage.runtimeMemoryTitle': 'Go Runtime Memory',
    'monitor.runtimePage.processBuildTitle': 'Process and Build',
    'monitor.runtimePage.hostEnvironmentTitle': 'Host Environment',
    'monitor.runtimePage.snapshotContextTitle': 'Snapshot Context',
    'monitor.runtimePage.summary.uptime': 'Uptime',
    'monitor.runtimePage.summary.uptimeDescription': 'Process running duration',
    'monitor.runtimePage.summary.goroutines': 'Goroutines',
    'monitor.runtimePage.summary.goroutinesDescription': 'Current Go scheduler concurrency',
    'monitor.runtimePage.summary.goVersion': 'Go version',
    'monitor.runtimePage.summary.goVersionDescription': 'Version reported by the current process',
    'monitor.runtimePage.summary.gcCycles': 'GC cycles',
    'monitor.runtimePage.summary.gcCyclesDescription': 'Completed garbage collection rounds',
    'monitor.runtimePage.fields.runtimeAlloc': 'Runtime alloc',
    'monitor.runtimePage.fields.runtimeHeap': 'Heap in use',
    'monitor.runtimePage.fields.runtimeSys': 'Runtime sys',
    'monitor.runtimePage.fields.gcCycles': 'GC count',
    'monitor.runtimePage.fields.lastGc': 'Last GC time',
    'monitor.runtimePage.fields.buildVersion': 'Build version',
    'monitor.runtimePage.fields.gitCommit': 'Git commit',
    'monitor.runtimePage.fields.appName': 'Application',
    'monitor.runtimePage.fields.appEnv': 'Environment',
    'monitor.runtimePage.fields.startedAt': 'Started at',
    'monitor.runtimePage.fields.hostName': 'Host name',
    'monitor.runtimePage.fields.platform': 'Platform',
    'monitor.runtimePage.fields.cpuCores': 'CPU cores',
    'monitor.runtimePage.fields.hostMemory': 'Host memory',
    'monitor.runtimePage.fields.observedAt': 'Last observed',
    'monitor.runtimePage.fields.loadAverage': 'Load average',
    'monitor.runtimePage.fieldDescriptions.runtimeAlloc': 'Bytes currently allocated by the Go runtime.',
    'monitor.runtimePage.fieldDescriptions.runtimeHeap': 'Heap bytes actively used by the Go process.',
    'monitor.runtimePage.fieldDescriptions.runtimeSys': 'Total memory requested from the system by the Go runtime.',
    'monitor.runtimePage.fieldDescriptions.gcCycles': 'Current snapshot of completed GC cycles.',
    'monitor.runtimePage.fieldDescriptions.lastGc': 'Reserved until the backend exposes the most recent GC timestamp.',
    'monitor.runtimePage.fieldDescriptions.buildVersion': 'Uses the current backend version field when present.',
    'monitor.runtimePage.fieldDescriptions.gitCommit':
      'Reserved until the backend exposes an explicit commit identifier.',
    'monitor.runtimePage.fieldDescriptions.appName': 'Application identifier reported by the service.',
    'monitor.runtimePage.fieldDescriptions.appEnv': 'Current runtime environment label.',
    'monitor.runtimePage.fieldDescriptions.startedAt': 'First stable process startup time.',
    'monitor.runtimePage.fieldDescriptions.hostName': 'Host reported by the runtime snapshot.',
    'monitor.runtimePage.fieldDescriptions.platform': 'Operating system and architecture for the host.',
    'monitor.runtimePage.fieldDescriptions.cpuCores': 'Logical CPU core count visible to the process.',
    'monitor.runtimePage.fieldDescriptions.hostMemory': 'Whole-host memory snapshot, not Go Runtime memory.',
    'monitor.runtimePage.fieldDescriptions.observedAt': 'Timestamp of the latest aggregated snapshot.',
    'monitor.runtimePage.fieldDescriptions.loadAverage': '1m / 5m / 15m load averages from the current host.',
  }),
);

vi.mock('../api/server-status', () => ({
  getServerStatus: monitorApiMocks.getServerStatus,
}));

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => translations[key] ?? key,
  }),
}));

const passthroughStub = defineComponent({
  name: 'PassthroughStub',
  props: {
    title: {
      type: String,
      default: '',
    },
    theme: {
      type: String,
      default: '',
    },
  },
  setup(props, { slots }) {
    return () => h('div', { 'data-theme': props.theme || undefined }, [props.title, slots.default?.()]);
  },
});

const buttonStub = defineComponent({
  name: 'TButtonStub',
  setup(_props, { slots }) {
    return () => h('button', slots.default?.());
  },
});

function mountRuntimePage() {
  return mount(RuntimePage, {
    global: {
      stubs: {
        't-card': passthroughStub,
        't-tag': passthroughStub,
        't-button': buttonStub,
        't-empty': passthroughStub,
      },
    },
  });
}

function createResponse() {
  return {
    status: 'healthy',
    observed_at: '2026-05-21T10:30:00Z',
    server: {
      version: 'v0.3.2',
      started_at: '2026-05-21T08:30:00Z',
      uptime_seconds: 7200,
      go_version: 'go1.26.0',
      app_name: 'graft',
      app_env: 'prod',
    },
    runtime: {
      go_version: 'go1.26.0',
      host_name: 'node-a',
      operating_system: 'linux',
      architecture: 'amd64',
      cpu_cores: 8,
      load_average: {
        one_minute: 0.34,
        five_minutes: 0.28,
        fifteen_minutes: 0.22,
      },
      disk_usage: {
        path: '/',
        total_bytes: 0,
        used_bytes: 0,
        free_bytes: 0,
        used_percent: 0,
      },
      host_memory_total_bytes: 17179869184,
      host_memory_used_bytes: 8589934592,
      host_memory_free_bytes: 8589934592,
      host_memory_used_percent: 50,
      goroutines: 37,
      runtime_alloc_bytes: 41943040,
      runtime_heap_in_use_bytes: 31457280,
      runtime_sys_bytes: 83886080,
      runtime_gc_cycles: 12,
    },
    dependencies: {
      database: { status: 'healthy', detail: 'ok', latency_ms: 2.1 },
      redis: { status: 'disabled', detail: 'not configured', latency_ms: null },
    },
    summary: {
      total_dependencies: 2,
      healthy_dependencies: 1,
      degraded_dependencies: 0,
      unknown_dependencies: 0,
      disabled_dependencies: 1,
      total_plugins: 5,
      healthy_plugins: 4,
    },
    trend: {
      range: '10m',
      retention_seconds: 600,
      sample_interval_seconds: 5,
      points: [],
    },
    plugins: [],
  };
}

describe('monitor runtime page', () => {
  it('renders runtime and host sections with explicit memory separation and reserved fields', async () => {
    monitorApiMocks.getServerStatus.mockResolvedValue(createResponse());

    const wrapper = mountRuntimePage();
    await flushPromises();

    expect(wrapper.attributes('data-page-type')).toBe('overview-dashboard');
    expect(wrapper.text()).toContain('Runtime');
    expect(wrapper.text()).toContain('Go Runtime Memory');
    expect(wrapper.text()).toContain('Host Environment');
    expect(wrapper.text()).toContain('Host memory');
    expect(wrapper.text()).toContain('Git commit');
    expect(wrapper.text()).toContain('Not reported');
    expect(wrapper.text()).toContain('Build version');
    expect(wrapper.text()).toContain('v0.3.2');
  });
});
