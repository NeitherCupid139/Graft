// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { shallowMount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import { defineComponent, h } from 'vue';
import { createI18n } from 'vue-i18n';

import type { AppLogItem } from '../types/app-log';
import AppLogDetailDrawer from './AppLogDetailDrawer.vue';

const LogJsonPanelStub = defineComponent({
  name: 'LogJsonPanelStub',
  props: [
    'title',
    'expandLabel',
    'collapseLabel',
    'copyLabel',
    'copySuccessLabel',
    'copyFailLabel',
    'emptyText',
    'value',
  ],
  setup(props) {
    return () => h('pre', { 'data-testid': `json-panel-${props.title}` }, JSON.stringify(props.value));
  },
});

const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  messages: {
    'zh-CN': {
      appLog: {
        actions: { copy: '复制', copyFail: '复制失败', copySuccess: '已复制' },
        columns: {
          component: '组件',
          occurredAt: '发生时间',
          operation: '事件 Key',
          severity: '级别',
        },
        detail: {
          basic: '基础信息',
          collapseContext: '收起完整应用日志 JSON',
          contextEmpty: '当前应用日志没有可展示的上下文。',
          copyContext: '复制 JSON',
          copyContextFail: '复制应用日志 JSON 失败',
          copyContextSuccess: '应用日志 JSON 已复制',
          correlation: '关联信息',
          error: '错误',
          expandContext: '展开完整应用日志 JSON',
          fields: '结构化字段',
          message: '消息',
          rawJson: '原始 JSON',
        },
        filters: {
          method: '方法',
          requestId: '请求 ID',
          route: '路由模板',
          traceId: 'Trace ID',
        },
        page: { detailTitle: '应用日志详情' },
        values: {
          emptyField: '-',
          noError: '无错误',
          noOperation: '未记录事件',
        },
      },
    },
  },
});

function appLogRecord(): AppLogItem {
  return {
    component: 'internal.dashboard',
    error: '',
    fields: { widget: 'summary' },
    id: 7,
    message: 'dashboard widget loaded',
    method: 'GET',
    occurred_at: '2026-06-13T08:00:00Z',
    operation: 'dashboard_widget_load',
    request_id: 'req-7',
    route: '/api/dashboard/summary',
    severity: 'info',
    trace_id: 'trace-7',
  } as AppLogItem;
}

describe('AppLogDetailDrawer', () => {
  it('renders structured fields and full raw JSON panels inside the drawer', () => {
    const record = appLogRecord();
    const wrapper = shallowMount(AppLogDetailDrawer, {
      props: {
        initialTab: 'raw',
        record,
        visible: true,
      },
      global: {
        plugins: [i18n],
        stubs: {
          LogJsonPanel: LogJsonPanelStub,
          TButton: true,
          TDescriptions: { template: '<section><slot /></section>' },
          TDescriptionsItem: { template: '<div><slot /></div>' },
          TDrawer: { template: '<aside><slot /></aside>' },
          TTag: { template: '<span><slot /></span>' },
          TTabPanel: { template: '<div><slot /></div>' },
          TTabs: { template: '<div><slot /></div>' },
        },
      },
    });

    expect(wrapper.text()).toContain('dashboard widget loaded');
    expect(wrapper.get('[data-testid="json-panel-结构化字段"]').text()).toContain('"widget":"summary"');
    expect(wrapper.get('[data-testid="json-panel-原始 JSON"]').text()).toContain('"id":7');
    expect(wrapper.get('[data-testid="json-panel-原始 JSON"]').text()).toContain('"request_id":"req-7"');
  });
});
