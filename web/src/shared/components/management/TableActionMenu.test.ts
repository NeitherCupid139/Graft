// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

import { mount } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent, h } from 'vue';

import TableActionMenu from './TableActionMenu.vue';

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => (key === 'components.commonTable.more' ? '更多' : key),
  }),
}));

const TButtonStub = defineComponent({
  name: 'TButtonStub',
  props: {
    disabled: Boolean,
  },
  emits: ['click'],
  setup(props, { emit, slots }) {
    return () =>
      h(
        'button',
        {
          disabled: props.disabled,
          onClick: () => emit('click'),
        },
        slots.default?.(),
      );
  },
});

const TDropdownStub = defineComponent({
  name: 'TDropdownStub',
  props: {
    options: {
      default: () => [],
      type: Array,
    },
  },
  setup(_props, { slots }) {
    return () => h('div', { 'data-testid': 'action-dropdown' }, slots.default?.());
  },
});

describe('TableActionMenu', () => {
  it('uses the shared localized fallback for the more action label', () => {
    const wrapper = mount(TableActionMenu, {
      global: {
        stubs: {
          't-button': TButtonStub,
          't-dropdown': TDropdownStub,
        },
      },
      props: {
        actions: [
          {
            label: 'components.commonTable.detail',
            value: 'detail',
          },
          {
            label: 'components.commonTable.detail',
            value: 'edit',
          },
        ],
      },
    });

    expect(wrapper.get('[data-testid="action-dropdown"]').text()).toContain('更多');
  });
});
