// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import { defineComponent, h } from 'vue';
import { createI18n } from 'vue-i18n';

import type { AppLogItem } from '../types/app-log';
import AppLogTable from './AppLogTable.vue';

const AdvancedQueryPagedTableStub = defineComponent({
  name: 'AdvancedQueryPagedTableStub',
  props: ['columns'],
  setup(props, { slots }) {
    return () =>
      h('section', { 'data-testid': 'paged-table' }, [
        h('div', { 'data-testid': 'table-toolbar' }, slots.toolbar?.()),
        h(
          'div',
          { 'data-testid': 'table-columns' },
          (props.columns ?? []).map((column: { colKey: string }) => h('span', column.colKey)),
        ),
      ]);
  },
});

const translations: Record<string, string> = {
  'appLog.columns.actions': '操作',
  'appLog.columns.component': '组件',
  'appLog.columns.correlation': '关联字段',
  'appLog.columns.fields': '字段',
  'appLog.columns.message': '消息',
  'appLog.columns.occurredAt': '发生时间',
  'appLog.columns.operation': '操作',
  'appLog.columns.severity': '级别',
  'appLog.page.emptyTitle': '暂无应用日志',
};

const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  messages: {
    'zh-CN': translations,
  },
});

const TTagStub = defineComponent({
  name: 'TTagStub',
  setup(_props, { slots }) {
    return () => h('span', slots.default?.());
  },
});

function appLogRow(): AppLogItem {
  return {
    component: 'internal.dashboard',
    error: '',
    fields: {},
    id: 1,
    message: 'dashboard widget loaded',
    occurred_at: '2026-06-13T08:00:00Z',
    operation: 'dashboard_widget_load',
    request_id: 'req-1',
    severity: 'debug',
    trace_id: 'trace-1',
  } as AppLogItem;
}

describe('AppLogTable', () => {
  it('forwards the table toolbar slot into the shared paged table header', () => {
    const wrapper = mount(AppLogTable, {
      global: {
        plugins: [i18n],
        stubs: {
          AdvancedQueryPagedTable: AdvancedQueryPagedTableStub,
          't-tag': TTagStub,
        },
      },
      props: {
        current: 1,
        description: '列表说明',
        emptyDescription: '暂无数据',
        footerSummary: '共 1 条',
        pageSize: 20,
        rows: [appLogRow()],
        summary: '当前 1 条',
        total: 1,
        visibleColumnKeys: ['occurred_at', 'severity', 'component'],
      },
      slots: {
        toolbar: '<button data-testid="table-refresh">刷新</button>',
      },
    });

    expect(wrapper.get('[data-testid="table-toolbar"]').text()).toContain('刷新');
    expect(wrapper.get('[data-testid="table-refresh"]').text()).toBe('刷新');
  });
});
