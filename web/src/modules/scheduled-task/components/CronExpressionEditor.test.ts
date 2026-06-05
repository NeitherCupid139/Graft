import { mount } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';
import { defineComponent, h } from 'vue';
import { createI18n } from 'vue-i18n';

import CronExpressionEditor from './CronExpressionEditor.vue';

vi.mock('@vue-js-cron/light/dist/light.css', () => ({}));

vi.mock('@vue-js-cron/light', () => ({
  CronLight: defineComponent({
    name: 'CronLight',
    props: {
      disabled: Boolean,
      locale: {
        type: String,
        default: '',
      },
      modelValue: {
        type: String,
        default: '',
      },
      theme: {
        type: String,
        default: '',
      },
    },
    emits: ['update:model-value', 'error'],
    setup(props, { emit }) {
      return () =>
        h('div', { 'data-testid': 'cron-light', 'data-disabled': String(props.disabled) }, [
          h('span', { 'data-testid': 'cron-light-value' }, props.modelValue),
          h('button', { 'data-testid': 'cron-light-emit', onClick: () => emit('update:model-value', '*/5 * * * *') }),
          h('button', { 'data-testid': 'cron-light-error', onClick: () => emit('error', 'Visual editor error') }),
        ]);
    },
  }),
}));

const tDesignStubs = {
  TAlert: defineComponent({
    name: 'TAlert',
    props: {
      message: {
        type: String,
        default: '',
      },
      theme: {
        type: String,
        default: '',
      },
    },
    setup(props) {
      return () => h('div', { 'data-testid': 'cron-invalid-alert', 'data-theme': props.theme }, props.message);
    },
  }),
  TCard: defineComponent({
    name: 'TCard',
    props: {
      title: {
        type: String,
        default: '',
      },
    },
    setup(props, { slots }) {
      return () => h('section', [h('h3', props.title), slots.default?.()]);
    },
  }),
  TInput: defineComponent({
    name: 'TInput',
    props: {
      disabled: Boolean,
      modelValue: {
        type: String,
        default: '',
      },
      placeholder: {
        type: String,
        default: '',
      },
      status: {
        type: String,
        default: '',
      },
      tips: {
        type: String,
        default: '',
      },
      value: {
        type: String,
        default: '',
      },
    },
    emits: ['update:value', 'change', 'blur'],
    setup(props, { emit }) {
      return () =>
        h('input', {
          'data-testid': 'cron-raw-input',
          disabled: props.disabled,
          placeholder: props.placeholder,
          value: props.value ?? props.modelValue,
          'data-status': props.status,
          'data-tips': props.tips,
          onInput: (event: Event) => {
            const value = (event.target as HTMLInputElement).value;
            emit('update:value', value);
            emit('change', value);
          },
          onBlur: (event: Event) => emit('blur', (event.target as HTMLInputElement).value),
        });
    },
  }),
  TSpace: defineComponent({
    name: 'TSpace',
    setup(_props, { slots }) {
      return () => h('div', slots.default?.());
    },
  }),
  TTag: defineComponent({
    name: 'TTag',
    props: {
      theme: {
        type: String,
        default: '',
      },
      variant: {
        type: String,
        default: '',
      },
    },
    setup(props, { slots }) {
      return () => h('span', { 'data-testid': 'cron-valid-tag', 'data-theme': props.theme }, slots.default?.());
    },
  }),
};

const i18n = createI18n({
  legacy: false,
  locale: 'en-US',
  messages: {
    'en-US': {
      scheduledTask: {
        cronDescription: {
          custom: 'Current expression: {expression}',
          daily: 'Runs every day at {hour}:00.',
          everyMinute: 'Runs once every minute.',
          everyNMinutes: 'Runs once every {interval} minutes.',
          hourly: 'Runs hourly.',
          invalid: 'Invalid expression.',
          monthly: 'Runs monthly.',
          weekly: 'Runs weekly.',
        },
        cronValidation: {
          fieldCount:
            'Cron expression must contain either {unixFields} Unix fields or {secondsFields} seconds-based fields.',
          fieldRange: 'Cron {field} field must be * or a number between {min} and {max}.',
          stepRange: 'Cron {field} step must be between {min} and {max}.',
        },
        cronEditor: {
          emptyExpression: 'No expression',
          expressionLabel: 'Cron Expression',
          expressionPlaceholder: 'For example */5 * * * *',
          inputHint: 'Use 5-field Unix or 6-field seconds Cron.',
          title: 'Cron Schedule',
          validStatus: 'Valid',
        },
      },
    },
  },
});

function mountEditor(props: { modelValue?: string; disabled?: boolean; error?: string } = {}) {
  return mount(CronExpressionEditor, {
    props: {
      modelValue: props.modelValue ?? '* * * * *',
      disabled: props.disabled,
      error: props.error,
    },
    global: {
      components: tDesignStubs,
      plugins: [i18n],
    },
  });
}

describe('CronExpressionEditor', () => {
  it('emits normalized 6-field cron when the raw input receives 5 fields', async () => {
    const wrapper = mountEditor();

    await wrapper.get('[data-testid="cron-raw-input"]').setValue('*/5 * * * *');

    expect(wrapper.emitted('update:modelValue')?.at(-1)).toEqual(['0 */5 * * * *']);
    expect(wrapper.emitted('validate')?.at(-1)).toEqual([
      {
        normalizedExpression: '0 */5 * * * *',
        valid: true,
      },
    ]);
  });

  it('normalizes values emitted by the third-party cron editor wrapper', async () => {
    const wrapper = mountEditor();

    await wrapper.get('[data-testid="cron-light-emit"]').trigger('click');

    expect(wrapper.emitted('update:modelValue')?.at(-1)).toEqual(['0 */5 * * * *']);
  });

  it('indicates invalid expression state and emits validation metadata', async () => {
    const wrapper = mountEditor({ modelValue: '0 0 24 * * *' });

    expect(wrapper.find('[data-testid="cron-valid-tag"]').exists()).toBe(false);
    expect(wrapper.get('[data-testid="cron-invalid-alert"]').text()).toContain(
      'Cron hours field must be * or a number between 0 and 23.',
    );
    expect(wrapper.emitted('validate')?.at(-1)).toEqual([
      {
        messageKey: 'scheduledTask.cronValidation.fieldRange',
        messageParams: { field: 'hours', min: 0, max: 23 },
        normalizedExpression: '0 0 24 * * *',
        valid: false,
      },
    ]);
  });

  it('propagates disabled state to direct editable controls and the visual editor', () => {
    const wrapper = mountEditor({ disabled: true });

    expect(wrapper.get<HTMLInputElement>('[data-testid="cron-raw-input"]').element.disabled).toBe(true);
    expect(wrapper.get('[data-testid="cron-light"]').attributes('data-disabled')).toBe('true');
  });

  it('mounts without importing the scheduled task list page', () => {
    const wrapper = mountEditor();

    expect(wrapper.findComponent({ name: 'CronExpressionEditor' }).exists()).toBe(true);
    expect(wrapper.findComponent({ name: 'ScheduledTaskListPage' }).exists()).toBe(false);
  });
});
