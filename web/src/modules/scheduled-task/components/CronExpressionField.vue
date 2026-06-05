<template>
  <div class="cron-expression-field" :class="{ 'cron-expression-field--disabled': disabled }">
    <t-input-adornment class="cron-expression-field__control">
      <t-input
        v-model="rawExpression"
        data-testid="cron-expression-input"
        :disabled="disabled"
        :placeholder="t('scheduledTask.cronExpressionField.placeholder')"
        :status="inputStatus"
        clearable
        @change="handleRawInput"
        @blur="handleRawInput"
      />
      <template #append>
        <t-button
          data-testid="cron-config-button"
          theme="default"
          variant="text"
          :disabled="disabled"
          @click="dialogVisible = true"
        >
          {{ t('scheduledTask.cronExpressionField.configure') }}
        </t-button>
      </template>
    </t-input-adornment>

    <div class="cron-expression-field__meta" data-testid="cron-expression-meta">
      <t-tag v-if="displayValid" size="small" theme="success" variant="light">
        {{ t('scheduledTask.cronExpressionField.validStatus') }}
      </t-tag>
      <span :class="{ 'cron-expression-field__message--error': Boolean(invalidMessage) }">
        {{ invalidMessage || descriptionText }}
      </span>
    </div>

    <cron-schedule-dialog
      v-model:visible="dialogVisible"
      :model-value="rawExpression"
      :disabled="disabled"
      @confirm="handleDialogConfirm"
    />
  </div>
</template>
<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  type CronValidationResult,
  describeCronExpression,
  normalizeCronExpression,
  validateCronExpression,
} from '../utils/cron';
import { translateCronDescription, translateCronValidation } from '../utils/cron-i18n';
import CronScheduleDialog from './CronScheduleDialog.vue';

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

const { t } = useI18n();
const rawExpression = ref(props.modelValue);
const dialogVisible = ref(false);

const normalizedExpression = computed(() => normalizeCronExpression(rawExpression.value));
const validation = computed(() => validateCronExpression(rawExpression.value));
const description = computed(() => describeCronExpression(rawExpression.value));
const invalidMessage = computed(() => props.error || cronValidationMessageText(validation.value));
const displayValid = computed(() => validation.value.valid && !props.error);
const inputStatus = computed(() => (invalidMessage.value ? 'error' : 'default'));
const descriptionText = computed(() => translateCronDescription(description.value, t));

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

function handleRawInput(value: string | number) {
  applyExpression(String(value));
}

function handleDialogConfirm(expression: string) {
  applyExpression(expression);
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

function cronValidationMessageText(result: CronValidationResult) {
  return translateCronValidation(result, t);
}
</script>
<style scoped lang="less">
.cron-expression-field {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-8);
  min-width: 0;
}

.cron-expression-field__control {
  width: 100%;
}

.cron-expression-field__control :deep(.t-input-adornment__append) {
  background: var(--td-bg-color-container);
}

.cron-expression-field__meta {
  align-items: center;
  color: var(--td-text-color-secondary);
  display: flex;
  gap: var(--graft-density-gap-8);
  min-height: var(--td-comp-size-xs);
  min-width: 0;
}

.cron-expression-field__meta span {
  min-width: 0;
  overflow-wrap: anywhere;
}

.cron-expression-field__message--error {
  color: var(--td-error-color);
}

.cron-expression-field--disabled {
  opacity: 0.72;
}
</style>
