import { flushPromises, mount } from '@vue/test-utils';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { defineComponent, h, KeepAlive } from 'vue';
import { createI18n } from 'vue-i18n';
import { createMemoryHistory, createRouter, RouterView } from 'vue-router';

import AuditLogsPage from './index.vue';

const auditApiMocks = vi.hoisted(() => ({
  getAuditLogs: vi.fn(async () => ({
    items: [
      {
        id: 1,
        actor_user_id: 1,
        actor_username: 'admin',
        actor_display_name: 'Admin',
        action: 'role.delete',
        resource_type: 'role',
        resource_id: '12',
        resource_name: 'Ops Admin',
        success: false,
        result: 'DENIED',
        risk_level: 'CRITICAL',
        target_type: 'ROLE',
        target_label: '角色',
        request_id: 'req-1',
        trace_id: 'trace-1',
        session_id: 'sess-1',
        ip: '127.0.0.1',
        user_agent: 'vitest',
        request_method: 'POST',
        request_path: '/api/roles/12/delete',
        status_code: 403,
        message: 'role removed',
        metadata: {
          trace_id: 'trace-1',
          session_id: 'sess-1',
        },
        created_at: '2026-05-27T08:00:00Z',
      },
    ],
    total: 1,
    page: 1,
    page_size: 10,
  })),
}));

vi.mock('../../api/audit', () => ({
  getAuditLogs: auditApiMocks.getAuditLogs,
}));

vi.mock('@/modules/shared/localized-api-error', () => ({
  resolveLocalizedErrorMessage: () => 'load failed',
}));

vi.mock('@/utils/logger', () => ({
  createLogger: () => ({
    debug: vi.fn(),
    error: vi.fn(),
  }),
}));

vi.mock('../../components/AuditFilters.vue', () => ({
  default: defineComponent({
    name: 'AuditFiltersStub',
    props: ['presets', 'activePreset', 'modelValue'],
    emits: ['search', 'reset', 'apply-preset', 'update:modelValue'],
    setup(props, { emit }) {
      return () =>
        h('div', [
          h('span', { 'data-testid': 'audit-filter-model' }, JSON.stringify(props.modelValue)),
          h('button', { 'data-testid': 'audit-search', onClick: () => emit('search') }, 'search'),
          h('button', { 'data-testid': 'audit-reset', onClick: () => emit('reset') }, 'reset'),
          h('button', { 'data-testid': 'audit-preset', onClick: () => emit('apply-preset', 'high-risk') }, 'preset'),
          h(
            'button',
            {
              'data-testid': 'audit-route-sync',
              onClick: () =>
                emit('update:modelValue', {
                  ...props.modelValue,
                  actor: 'route-admin',
                  success: 'all',
                  createdRange: ['2026-05-01 10:00:00', '2026-05-02 18:30:00'],
                  actionPrefixes: [],
                  actionKeywords: [],
                  requestPathPrefixes: [],
                  resourceTypes: [],
                  result: 'FAILED',
                  results: [],
                  sorters: [{ field: 'created_at', direction: 'asc' }],
                  riskLevels: [],
                }),
            },
            'sync-route',
          ),
        ]);
    },
  }),
}));

vi.mock('../../components/AuditTable.vue', () => ({
  default: defineComponent({
    name: 'AuditTableStub',
    props: ['rows', 'summary', 'footerSummary'],
    emits: ['detail', 'update:current', 'update:pageSize', 'page-change'],
    setup(props, { emit }) {
      return () =>
        h('div', [
          props.summary,
          props.footerSummary,
          h('span', JSON.stringify(props.rows)),
          h('button', { 'data-testid': 'audit-detail', onClick: () => emit('detail', props.rows?.[0]) }, 'detail'),
        ]);
    },
  }),
}));

vi.mock('../../components/AuditDetailDrawer.vue', () => ({
  default: defineComponent({
    name: 'AuditDetailDrawerStub',
    props: ['visible', 'record', 'monitorOrigin'],
    setup(props) {
      return () => h('div', [String(props.visible), props.record?.request_id, JSON.stringify(props.monitorOrigin)]);
    },
  }),
}));

const passthroughStub = defineComponent({
  name: 'PassthroughStub',
  props: ['title', 'description'],
  setup(props, { slots }) {
    return () => h('div', [props.title, props.description, slots.default?.(), slots.actions?.()]);
  },
});

const buttonStub = defineComponent({
  name: 'TButtonStub',
  emits: ['click'],
  setup(_, { emit, slots, attrs }) {
    return () => h('button', { ...attrs, onClick: () => emit('click') }, slots.default?.());
  },
});

const checkboxGroupStub = defineComponent({
  name: 'TCheckboxGroupStub',
  setup(_, { slots }) {
    return () => h('div', slots.default?.());
  },
});

const checkboxStub = defineComponent({
  name: 'TCheckboxStub',
  setup(_, { slots }) {
    return () => h('label', slots.default?.());
  },
});

const drawerStub = defineComponent({
  name: 'TDrawerStub',
  props: ['visible', 'header'],
  setup(_, { slots }) {
    return () => h('div', slots.default?.());
  },
});

const i18n = createI18n({
  legacy: false,
  locale: 'en-US',
  messages: {
    'en-US': {
      menu: {
        audit: {
          title: 'Security Audit',
          logs: {
            title: 'Audit Logs',
          },
        },
      },
      components: {
        commonTable: {
          operation: 'Operation',
        },
      },
      audit: {
        common: {
          unknownActor: 'Anonymous',
          unknownResource: 'Unknown resource',
          result: { SUCCESS: 'Success', FAILED: 'Business Failed', DENIED: 'Denied', ERROR: 'System Error' },
          risk: { LOW: 'Low', MEDIUM: 'Medium', HIGH: 'High', CRITICAL: 'Critical' },
        },
        logList: {
          title: 'Audit Logs',
          description:
            'Query security audit events and inspect request and permission context; health checks, monitor polling, and bootstrap requests are not part of the default audit dataset.',
          refresh: 'Refresh',
          columnSettings: 'Columns',
          retry: 'Retry',
          detail: 'View Details',
          more: 'More',
          footerTotal: '{count} audit events total',
          footerFiltered: '{count} records matched on this page',
          currentPageFiltered: 'Current page filter',
          loadFailed: 'Failed to load audit logs',
          errorTitle: 'Audit logs are temporarily unavailable',
          emptyTitle: 'No audit logs',
          emptyDescription: 'Adjust filters and try again.',
          detailTitle: 'Audit Detail',
          presets: {
            all: 'All',
            failedOperations: 'Failed Operations',
            todayAnomalies: "Today's Security Anomalies",
            rbacChanges: 'Permission Configuration Changes',
            permissionDenied: 'Permission Denied',
            sensitiveOps: 'Sensitive Operations',
            authFailed: 'Auth Failed',
            highRisk: 'High Risk',
          },
          actions: {
            search: 'Search',
            reset: 'Reset',
            backToMonitor: 'Back to monitor',
            addFilter: 'Add filter',
            showAdvanced: 'Advanced Filters',
            hideAdvanced: 'Hide Advanced',
          },
          filters: {
            keywordPlaceholder: 'Keyword: action, request ID, audit target, operated object',
            actorPlaceholder: 'Actor',
            actionPlaceholder: 'Action type',
            actionPrefixesPlaceholder: 'Select action groups',
            actionKeywordsPlaceholder: 'Type an action keyword and press Enter',
            successPlaceholder: 'Success state',
            datePlaceholder: 'Time range',
            sourcePlaceholder: 'Source',
            resourceTypePlaceholder: 'Audit target type',
            resourceTypesPlaceholder: 'Select target type set',
            resourceNamePlaceholder: 'Audit target / target name',
            resourceIdPlaceholder: 'Resource ID',
            resultPlaceholder: 'Result',
            resultsPlaceholder: 'Select result set',
            riskPlaceholder: 'Risk',
            riskLevelsPlaceholder: 'Select risk level set',
            sessionPlaceholder: 'Session ID',
            requestIdPlaceholder: 'Request ID',
            requestPathPrefixesPlaceholder: 'Type a request path prefix and press Enter',
          },
          builder: {
            title: 'Filter fields',
            hint: 'Choose a field and set its value. Active conditions appear as removable tags.',
            fields: {
              success: 'Success state',
              action: 'Action type',
              actionPrefixes: 'Action groups',
              actionKeywords: 'Action keywords',
              result: 'Result',
              results: 'Result set',
              riskLevel: 'Risk level',
              riskLevels: 'Risk level set',
              source: 'Event type',
              actor: 'Actor',
              resourceName: 'Audit target',
              resourceType: 'Target type',
              resourceTypes: 'Target type set',
              requestPathPrefixes: 'Request path prefixes',
              requestId: 'Request ID',
              session: 'Session ID',
              resourceId: 'Resource ID',
            },
          },
          filterOptions: {
            allActions: 'All actions',
            allSource: 'All source',
            allResourceTypes: 'All target types',
            auth: 'Authentication',
            authPrefix: 'Authentication actions',
            rbacPrefix: 'Permission configuration actions',
            role: 'Role',
            rolePrefix: 'Role actions',
            permission: 'Permission',
            permissionPrefix: 'Permission actions',
            session: 'Session',
            userResource: 'User',
            roleResource: 'Role',
            permissionResource: 'Permission',
            authResource: 'Authentication',
            allResults: 'All results',
            SUCCESS: 'Success',
            FAILED: 'Business Failed',
            DENIED: 'Denied',
            ERROR: 'System Error',
            allRisk: 'All risk',
            LOW: 'Low',
            MEDIUM: 'Medium',
            HIGH: 'High',
            CRITICAL: 'Critical',
          },
          columns: {
            action: 'Action',
            actor: 'Actor',
            resource: 'Audit Target',
            correlation: 'Request ID',
            result: 'Result',
            risk: 'Risk',
            createdAt: 'Time',
          },
          drawer: {
            messageFallback: 'No additional message',
            sections: {
              basic: 'Event Summary',
              request: 'Request Context',
              correlation: 'Related Context',
              risk: 'Risk',
              metadata: 'Metadata',
            },
            fields: {
              target: 'Audit Target',
              result: 'Result',
              requestId: 'Request ID',
              sessionId: 'Session ID',
              ip: 'IP',
              userAgent: 'User-Agent',
              method: 'Method',
              path: 'Path',
              status: 'Status',
            },
            related: {
              sameRequest: 'Same Request ID',
              sameActor: 'Recent Actions by Actor',
              sameResource: 'Recent Changes on Audit Target',
              empty: 'No more related records in the current list',
            },
            risk: {
              failedOperation: 'Failed operation',
              sensitiveOperation: 'Sensitive write',
              requestTrace: 'Request available',
            },
            actions: {
              viewRelatedRequest: 'View Related Request',
              copyMetadata: 'Copy JSON',
              copyMetadataSuccess: 'Metadata JSON copied',
              copyMetadataFail: 'Failed to copy metadata JSON',
              toggleMetadata: 'Show raw metadata',
            },
            metadataEmpty: 'No metadata is available for this event.',
          },
        },
      },
    },
  },
});

describe('AuditLogsPage', () => {
  beforeEach(() => {
    auditApiMocks.getAuditLogs.mockClear();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  async function mountPage(
    initialQuery: Record<string, string> = {
      created_from: '2026-05-30T07:21:04.000Z',
      created_to: '2026-05-31T07:21:04.000Z',
      results: 'DENIED',
    },
  ) {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/audit/logs', component: AuditLogsPage }],
    });

    await router.push({
      path: '/audit/logs',
      query: initialQuery,
    });
    await router.isReady();

    const replaceSpy = vi.spyOn(router, 'replace');
    const wrapper = mount(AuditLogsPage, {
      global: {
        plugins: [i18n, router],
        stubs: {
          'management-empty-state': passthroughStub,
          'management-page-content': passthroughStub,
          'management-page-header': passthroughStub,
          't-button': buttonStub,
          't-checkbox': checkboxStub,
          't-checkbox-group': checkboxGroupStub,
          't-drawer': drawerStub,
          't-space': passthroughStub,
        },
      },
    });

    await flushPromises();
    return { router, replaceSpy, wrapper };
  }

  async function mountKeepAliveHost(initialQuery: Record<string, string> = {}) {
    const OtherPage = defineComponent({
      name: 'OtherPageStub',
      setup: () => () => h('div', { 'data-testid': 'other-page' }, 'other'),
    });

    const RouterHost = defineComponent({
      name: 'RouterHost',
      setup() {
        return () =>
          h(RouterView, null, {
            default: ({ Component }: { Component: unknown }) => h(KeepAlive, null, () => [h(Component as never)]),
          });
      },
    });

    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/audit/logs', name: 'AuditLogList', component: AuditLogsPage },
        { path: '/users', name: 'UsersIndex', component: OtherPage },
      ],
    });

    await router.push({
      path: '/audit/logs',
      query: initialQuery,
    });
    await router.isReady();

    const replaceSpy = vi.spyOn(router, 'replace');
    const wrapper = mount(RouterHost, {
      global: {
        plugins: [i18n, router],
        stubs: {
          'management-empty-state': passthroughStub,
          'management-page-content': passthroughStub,
          'management-page-header': passthroughStub,
          't-button': buttonStub,
          't-checkbox': checkboxStub,
          't-checkbox-group': checkboxGroupStub,
          't-drawer': drawerStub,
          't-space': passthroughStub,
        },
      },
    });

    await flushPromises();
    return { router, replaceSpy, wrapper };
  }

  it('restores deep-link filters including created range and keeps backend request shape unchanged', async () => {
    const { wrapper } = await mountPage({
      actor: 'alice',
      created_from: '2026-05-01T10:00:00Z',
      created_to: '2026-05-02T18:30:00Z',
      result: 'FAILED',
    });

    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).toContain('"actor":"alice"');
    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).toContain(
      '"createdRange":["2026-05-01 18:00:00","2026-05-03 02:30:00"]',
    );
    expect(auditApiMocks.getAuditLogs).toHaveBeenLastCalledWith({
      page: 1,
      page_size: 10,
      actor: 'alice',
      result: 'FAILED',
      created_from: '2026-05-01T10:00:00.000Z',
      created_to: '2026-05-02T18:30:00.000Z',
      sort_by: 'created_at',
      sort_order: 'desc',
    });
  });

  it('loads explicit-range records and opens the detail drawer', async () => {
    const { wrapper } = await mountPage();

    expect(auditApiMocks.getAuditLogs).toHaveBeenCalledWith(
      expect.objectContaining({
        created_from: '2026-05-30T07:21:04.000Z',
        created_to: '2026-05-31T07:21:04.000Z',
        results: ['DENIED'],
        sort_by: 'created_at',
        sort_order: 'desc',
      }),
    );
    expect(wrapper.text()).not.toContain('security audit records shown');
    expect(wrapper.text()).not.toContain('Core fields only');
    expect(wrapper.text()).toContain('req-1');

    await wrapper.get('[data-testid="audit-detail"]').trigger('click');
    await flushPromises();
    expect(wrapper.text()).toContain('true');
    expect(wrapper.text()).toContain('req-1');
  });

  it('keeps monitor return context when syncing log filters', async () => {
    const { replaceSpy, router, wrapper } = await mountPage({
      created_from: '2026-05-30T07:21:04.000Z',
      created_to: '2026-05-31T07:21:04.000Z',
      results: 'DENIED',
      monitorView: 'overview',
      monitorTrendRange: '10m',
      monitorAnomalyKey: 'resource_cpu_pressure',
      monitorScopeRef: 'runtime:cpu',
    });

    auditApiMocks.getAuditLogs.mockClear();
    replaceSpy.mockClear();

    await wrapper.get('[data-testid="audit-route-sync"]').trigger('click');
    await wrapper.get('[data-testid="audit-search"]').trigger('click');
    await flushPromises();

    expect(replaceSpy).toHaveBeenCalledWith(
      expect.objectContaining({
        path: '/audit/logs',
        query: expect.objectContaining({
          monitorView: 'overview',
          monitorTrendRange: '10m',
          monitorAnomalyKey: 'resource_cpu_pressure',
          monitorScopeRef: 'runtime:cpu',
        }),
      }),
    );
    expect(router.currentRoute.value.query).toMatchObject({
      monitorView: 'overview',
      monitorTrendRange: '10m',
      monitorAnomalyKey: 'resource_cpu_pressure',
      monitorScopeRef: 'runtime:cpu',
    });
    expect(wrapper.text()).toContain('"view":"overview"');
  });

  it('applies quick preset from filters and refetches with unchanged query contract', async () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-05-31T07:21:04Z'));
    const { wrapper } = await mountPage();
    auditApiMocks.getAuditLogs.mockClear();

    await wrapper.get('[data-testid="audit-preset"]').trigger('click');
    await flushPromises();

    expect(auditApiMocks.getAuditLogs).toHaveBeenCalledWith({
      page: 1,
      page_size: 10,
      created_from: '2026-05-30T07:21:04.000Z',
      created_to: '2026-05-31T07:21:04.000Z',
      risk_levels: ['HIGH', 'CRITICAL'],
      sort_by: 'created_at',
      sort_order: 'desc',
    });
  });

  it('does not send an implicit preset when the route has no time range', async () => {
    const { wrapper } = await mountPage({});

    expect(auditApiMocks.getAuditLogs).toHaveBeenLastCalledWith({
      page: 1,
      page_size: 10,
      sort_by: 'created_at',
      sort_order: 'desc',
    });
    expect(wrapper.text()).toContain('req-1');
  });

  it('syncs interactive filters into route query for reload and sharing', async () => {
    const { replaceSpy, router, wrapper } = await mountPage();
    auditApiMocks.getAuditLogs.mockClear();
    replaceSpy.mockClear();

    await wrapper.get('[data-testid="audit-route-sync"]').trigger('click');
    await wrapper.get('[data-testid="audit-search"]').trigger('click');
    await flushPromises();

    expect(replaceSpy).toHaveBeenCalledWith(
      expect.objectContaining({
        path: '/audit/logs',
        query: expect.objectContaining({
          actor: 'route-admin',
          created_from: '2026-05-01T02:00:00.000Z',
          created_to: '2026-05-02T10:30:00.000Z',
          result: 'FAILED',
          sort_by: 'created_at',
          sort_order: 'asc',
        }),
      }),
    );
    expect(router.currentRoute.value.query).toMatchObject({
      actor: 'route-admin',
      created_from: '2026-05-01T02:00:00.000Z',
      created_to: '2026-05-02T10:30:00.000Z',
      result: 'FAILED',
      sort_by: 'created_at',
      sort_order: 'asc',
    });
    expect(auditApiMocks.getAuditLogs).toHaveBeenLastCalledWith(
      expect.objectContaining({
        result: 'FAILED',
        created_from: '2026-05-01T02:00:00.000Z',
        created_to: '2026-05-02T10:30:00.000Z',
        sort_by: 'created_at',
        sort_order: 'asc',
      }),
    );
  });

  it('preserves explicit created range over preset-derived display state', async () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-05-31T07:21:04Z'));

    const { wrapper } = await mountPage({
      created_from: '2026-05-01T10:00:00Z',
      created_to: '2026-05-02T18:30:00Z',
    });

    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).toContain(
      '"createdRange":["2026-05-01 18:00:00","2026-05-03 02:30:00"]',
    );
    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).not.toContain(
      '"createdRange":["2026-05-30 15:21:04","2026-05-31 15:21:04"]',
    );
  });

  it('ignores legacy route params and keeps only canonical visible filters', async () => {
    const { router, wrapper } = await mountPage({
      preset: 'last_24h',
      summary: 'failed-operations',
      risk_group: 'auth_failures',
      occurred_from: '2026-05-01T10:00:00Z',
      occurred_to: '2026-05-02T18:30:00Z',
      created_from: '2026-05-03T10:00:00Z',
      created_to: '2026-05-04T18:30:00Z',
      results: 'DENIED',
    });

    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).toContain(
      '"createdRange":["2026-05-03 18:00:00","2026-05-05 02:30:00"]',
    );
    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).not.toContain('"success":"false"');
    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).not.toContain(
      '"resourceTypes":["auth","session"]',
    );
    expect(router.currentRoute.value.query).toMatchObject({
      created_from: '2026-05-03T10:00:00.000Z',
      created_to: '2026-05-04T18:30:00.000Z',
      results: 'DENIED',
    });
    expect(router.currentRoute.value.query).not.toHaveProperty('preset');
    expect(router.currentRoute.value.query).not.toHaveProperty('summary');
    expect(router.currentRoute.value.query).not.toHaveProperty('risk_group');
    expect(router.currentRoute.value.query).not.toHaveProperty('occurred_from');
    expect(router.currentRoute.value.query).not.toHaveProperty('occurred_to');
    expect(auditApiMocks.getAuditLogs).toHaveBeenLastCalledWith({
      page: 1,
      page_size: 10,
      created_from: '2026-05-03T10:00:00.000Z',
      created_to: '2026-05-04T18:30:00.000Z',
      results: ['DENIED'],
      sort_by: 'created_at',
      sort_order: 'desc',
    });
  });

  it('writes back canonical query fields only after interactive changes', async () => {
    const { router, wrapper } = await mountPage({
      preset: 'last_24h',
      summary: 'failed-operations',
      risk_group: 'auth_failures',
      occurred_from: '2026-05-01T10:00:00Z',
      occurred_to: '2026-05-02T18:30:00Z',
      created_from: '2026-05-03T10:00:00Z',
      created_to: '2026-05-04T18:30:00Z',
    });

    await wrapper.get('[data-testid="audit-search"]').trigger('click');
    await flushPromises();

    expect(router.currentRoute.value.query).toMatchObject({
      created_from: '2026-05-03T10:00:00.000Z',
      created_to: '2026-05-04T18:30:00.000Z',
      sort_by: 'created_at',
      sort_order: 'desc',
    });
    expect(router.currentRoute.value.query).not.toHaveProperty('preset');
    expect(router.currentRoute.value.query).not.toHaveProperty('summary');
    expect(router.currentRoute.value.query).not.toHaveProperty('risk_group');
    expect(router.currentRoute.value.query).not.toHaveProperty('occurred_from');
    expect(router.currentRoute.value.query).not.toHaveProperty('occurred_to');
  });

  it('does not redirect back to audit logs after the kept-alive page is deactivated', async () => {
    const { replaceSpy, router, wrapper } = await mountKeepAliveHost({
      created_from: '2026-05-30T07:21:04.000Z',
      created_to: '2026-05-31T07:21:04.000Z',
      results: 'DENIED',
    });

    auditApiMocks.getAuditLogs.mockClear();
    replaceSpy.mockClear();

    await router.push({ path: '/users', query: { tab: 'active' } });
    await flushPromises();

    expect(router.currentRoute.value.path).toBe('/users');
    expect(router.currentRoute.value.query).toMatchObject({ tab: 'active' });
    expect(wrapper.get('[data-testid="other-page"]').text()).toBe('other');
    expect(replaceSpy).not.toHaveBeenCalledWith(
      expect.objectContaining({
        path: '/audit/logs',
      }),
    );
    expect(auditApiMocks.getAuditLogs).not.toHaveBeenCalled();
  });

  it('re-applies current route query when the kept-alive audit page is re-activated', async () => {
    const { router, wrapper } = await mountKeepAliveHost({
      created_from: '2026-05-30T07:21:04.000Z',
      created_to: '2026-05-31T07:21:04.000Z',
      results: 'DENIED',
    });

    await router.push({ path: '/users', query: { tab: 'active' } });
    await flushPromises();

    auditApiMocks.getAuditLogs.mockClear();

    await router.push({
      path: '/audit/logs',
      query: {
        resource_type: 'user',
        resource_name: 'Graft Admin',
        resource_id: '1',
      },
    });
    await flushPromises();

    expect(router.currentRoute.value.path).toBe('/audit/logs');
    expect(router.currentRoute.value.query).toMatchObject({
      resource_type: 'user',
      resource_name: 'Graft Admin',
      resource_id: '1',
    });
    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).toContain('"resourceType":"user"');
    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).toContain('"resourceName":"Graft Admin"');
    expect(wrapper.get('[data-testid="audit-filter-model"]').text()).toContain('"resourceId":"1"');
    expect(auditApiMocks.getAuditLogs).toHaveBeenLastCalledWith({
      page: 1,
      page_size: 10,
      resource_type: 'user',
      resource_name: 'Graft Admin',
      resource_id: '1',
      sort_by: 'created_at',
      sort_order: 'desc',
    });
  });
});
