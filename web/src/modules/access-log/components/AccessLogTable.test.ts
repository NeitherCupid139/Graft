// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import { defineComponent, h } from 'vue';
import { createI18n } from 'vue-i18n';

import type { AccessLogItem } from '../types/access-log';
import AccessLogTable from './AccessLogTable.vue';

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
  'accessLog.columns.durationMs': '耗时',
  'accessLog.columns.method': '方法',
  'accessLog.columns.occurredAt': '发生时间',
  'accessLog.columns.operation': '操作',
  'accessLog.columns.path': '路径',
  'accessLog.columns.requestId': '请求 ID',
  'accessLog.columns.startedAt': '开始时间',
  'accessLog.columns.statusCode': '状态码',
  'accessLog.columns.user': '用户',
  'accessLog.page.emptyTitle': '暂无访问日志',
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

function accessLogRow(): AccessLogItem {
  return {
    duration_ms: 1,
    id: 1,
    method: 'GET',
    occurred_at: '2026-06-13T08:00:00Z',
    path: '/api/access-log',
    request_id: 'req-1',
    route: '',
    started_at: '2026-06-13T08:00:00Z',
    status_code: 200,
    user_id: 1,
    username: 'graft',
  } as AccessLogItem;
}

describe('AccessLogTable', () => {
  it('forwards the table toolbar slot into the shared paged table header', () => {
    const wrapper = mount(AccessLogTable, {
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
        rows: [accessLogRow()],
        summary: '当前 1 条',
        total: 1,
        visibleColumnKeys: ['started_at', 'method', 'path'],
      },
      slots: {
        toolbar: '<button data-testid="table-refresh">刷新</button>',
      },
    });

    expect(wrapper.get('[data-testid="table-toolbar"]').text()).toContain('刷新');
    expect(wrapper.get('[data-testid="table-refresh"]').text()).toBe('刷新');
  });
});
