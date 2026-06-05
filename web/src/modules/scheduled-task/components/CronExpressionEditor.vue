<template>
  <div class="cron-expression-editor" :class="{ 'cron-expression-editor--disabled': disabled }">
    <t-card
      class="cron-expression-editor__card"
      :bordered="true"
      size="small"
      :title="t('scheduledTask.cronEditor.title')"
    >
      <t-space class="cron-expression-editor__stack" direction="vertical" size="small">
        <div class="cron-expression-editor__visual" data-testid="cron-visual-editor">
          <cron-light
            :model-value="visualExpression"
            :disabled="disabled"
            :locale="cronLocale"
            theme="ant"
            @update:model-value="handleVisualUpdate"
            @error="handleVisualError"
          />
        </div>

        <label class="cron-expression-editor__field">
          <span class="cron-expression-editor__label">{{ t('scheduledTask.cronEditor.expressionLabel') }}</span>
          <t-input
            v-model:value="rawExpression"
            class="cron-expression-editor__input"
            data-testid="cron-raw-input"
            :disabled="disabled"
            :placeholder="t('scheduledTask.cronEditor.expressionPlaceholder')"
            :status="inputStatus"
            :tips="inputTips"
            clearable
            @change="handleRawInput"
            @blur="handleRawInput"
          />
        </label>

        <div class="cron-expression-editor__meta">
          <t-space align="center" break-line size="small">
            <t-tag v-if="displayValid" data-testid="cron-valid-tag" theme="success" variant="light">
              {{ t('scheduledTask.cronEditor.validStatus') }}
            </t-tag>
            <span class="cron-expression-editor__description" data-testid="cron-description">
              {{ descriptionText }}
            </span>
          </t-space>
          <code class="cron-expression-editor__expression" data-testid="cron-normalized-expression">
            {{ normalizedExpression || t('scheduledTask.cronEditor.emptyExpression') }}
          </code>
        </div>

        <t-alert v-if="invalidMessage" data-testid="cron-invalid-alert" theme="error" :message="invalidMessage" />
      </t-space>
    </t-card>
  </div>
</template>
<script setup lang="ts">
import '@vue-js-cron/light/dist/light.css';

import { CronLight } from '@vue-js-cron/light';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  type CronDescriptionResult,
  type CronValidationResult,
  describeCronExpression,
  normalizeCronExpression,
  validateCronExpression,
} from '../utils/cron';

const props = withDefaults(
  defineProps<{
    modelValue: string;
    disabled?: boolean;
    error?: string;
  }>(),
  {
    disabled: false,
    error: '',
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: string];
  validate: [result: CronValidationResult & { normalizedExpression: string }];
}>();

const { locale, t } = useI18n();
const rawExpression = ref(props.modelValue);
const visualEditorError = ref('');

const normalizedExpression = computed(() => normalizeCronExpression(rawExpression.value));
const validation = computed(() => validateCronExpression(rawExpression.value));
const description = computed(() => describeCronExpression(rawExpression.value) as CronDescriptionResult);
const invalidMessage = computed(
  () => props.error || cronValidationMessageText(validation.value) || visualEditorError.value,
);
const displayValid = computed(() => validation.value.valid && !props.error && !visualEditorError.value);
const inputStatus = computed(() => (invalidMessage.value ? 'error' : 'default'));
const inputTips = computed(() => invalidMessage.value || t('scheduledTask.cronEditor.inputHint'));
const visualExpression = computed(() => {
  const fields = normalizedExpression.value.split(' ');
  return validation.value.valid && fields.length === 6 ? fields.slice(1).join(' ') : rawExpression.value;
});
const descriptionText = computed(() => translateDescription(description.value));
const cronLocale = computed(() => (String(locale.value).toLowerCase().startsWith('zh') ? 'zh' : 'en'));

watch(
  () => props.modelValue,
  (value) => {
    if (value !== rawExpression.value) {
      rawExpression.value = value;
    }
  },
);

watch(
  () =>
    [
      validation.value.valid,
      validation.value.messageKey,
      validation.value.messageParams,
      normalizedExpression.value,
    ] as const,
  () => {
    emitValidation();
  },
  { immediate: true },
);

function handleVisualUpdate(value: string) {
  visualEditorError.value = '';
  applyExpression(value);
}

function handleVisualError(value: string) {
  visualEditorError.value = value;
}

function handleRawInput(value: string | number) {
  visualEditorError.value = '';
  applyExpression(String(value));
}

function applyExpression(value: string) {
  const nextExpression = normalizeCronExpression(value);
  rawExpression.value = nextExpression;
  emit('update:modelValue', nextExpression);
  emitValidation();
}

function emitValidation() {
  emit('validate', {
    ...validation.value,
    normalizedExpression: normalizedExpression.value,
  });
}

function translateDescription(result: CronDescriptionResult) {
  switch (result.key) {
    case 'scheduledTask.cronDescription.everyMinute':
      return t('scheduledTask.cronDescription.everyMinute', result.params ?? {});
    case 'scheduledTask.cronDescription.everyNMinutes':
      return t('scheduledTask.cronDescription.everyNMinutes', result.params ?? {});
    case 'scheduledTask.cronDescription.hourly':
      return t('scheduledTask.cronDescription.hourly', result.params ?? {});
    case 'scheduledTask.cronDescription.daily':
      return t('scheduledTask.cronDescription.daily', result.params ?? {});
    case 'scheduledTask.cronDescription.weekly':
      return t('scheduledTask.cronDescription.weekly', result.params ?? {});
    case 'scheduledTask.cronDescription.monthly':
      return t('scheduledTask.cronDescription.monthly', result.params ?? {});
    case 'scheduledTask.cronDescription.custom':
      return t('scheduledTask.cronDescription.custom', result.params ?? {});
    case 'scheduledTask.cronDescription.invalid':
    default:
      return t('scheduledTask.cronDescription.invalid', result.params ?? {});
  }
}

function cronValidationMessageText(result: CronValidationResult) {
  switch (result.messageKey) {
    case 'scheduledTask.cronValidation.fieldCount':
      return t('scheduledTask.cronValidation.fieldCount', result.messageParams ?? {});
    case 'scheduledTask.cronValidation.stepRange':
      return t('scheduledTask.cronValidation.stepRange', result.messageParams ?? {});
    case 'scheduledTask.cronValidation.fieldRange':
      return t('scheduledTask.cronValidation.fieldRange', result.messageParams ?? {});
    default:
      return '';
  }
}
</script>
<style scoped lang="less">
.cron-expression-editor {
  color: var(--td-text-color-primary);
}

.cron-expression-editor__card {
  background: var(--td-bg-color-container);
}

.cron-expression-editor__stack {
  width: 100%;
}

.cron-expression-editor__visual {
  background: var(--td-bg-color-secondarycontainer);
  border: 1px solid var(--td-border-level-1-color);
  border-radius: var(--td-radius-default);
  overflow-x: auto;
  padding: var(--td-comp-paddingTB-s) var(--td-comp-paddingLR-s);
}

.cron-expression-editor__field {
  display: grid;
  gap: var(--td-comp-margin-xs);
}

.cron-expression-editor__label {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-medium);
}

.cron-expression-editor__input {
  width: 100%;
}

.cron-expression-editor__meta {
  display: grid;
  gap: var(--td-comp-margin-xs);
}

.cron-expression-editor__description {
  color: var(--td-text-color-secondary);
}

.cron-expression-editor__expression {
  background: var(--td-bg-color-secondarycontainer);
  border: 1px solid var(--td-border-level-1-color);
  border-radius: var(--td-radius-default);
  color: var(--td-text-color-primary);
  display: inline-flex;
  max-width: 100%;
  overflow-wrap: anywhere;
  padding: var(--td-comp-paddingTB-xs) var(--td-comp-paddingLR-s);
  width: fit-content;
}

.cron-expression-editor--disabled {
  opacity: 0.72;
}
</style>
