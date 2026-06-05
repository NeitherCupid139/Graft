<template>
  <div class="cron-expression-field" :class="{ 'cron-expression-field--disabled': disabled }">
    <div class="cron-expression-field__control">
      <div class="cron-expression-field__input">
        <t-input
          v-model="rawExpression"
          data-testid="cron-expression-input"
          :disabled="disabled"
          :placeholder="t('scheduledTask.cronExpressionField.placeholder')"
          :status="inputStatus"
          @change="handleRawInput"
          @blur="handleInputBlur"
          @focus="inputFocused = true"
        >
          <template #suffix>
            <span class="cron-expression-field__clear-space">
              <button
                v-show="showClearControl"
                class="cron-expression-field__clear-button"
                type="button"
                tabindex="-1"
                :aria-label="t('scheduledTask.cronExpressionField.clear')"
                @click.stop="handleClearExpression"
                @mousedown.prevent
              >
                <close-circle-filled-icon class="cron-expression-field__clear-icon" aria-hidden="true" />
              </button>
            </span>
          </template>
        </t-input>
      </div>
      <t-button
        class="cron-expression-field__configure"
        data-testid="cron-config-button"
        theme="default"
        variant="text"
        :disabled="disabled"
        @click="dialogVisible = true"
      >
        {{ t('scheduledTask.cronExpressionField.configure') }}
      </t-button>
    </div>

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
import { CloseCircleFilledIcon } from 'tdesign-icons-vue-next';
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
const inputFocused = ref(false);

const normalizedExpression = computed(() => normalizeCronExpression(rawExpression.value));
const validation = computed(() => validateCronExpression(rawExpression.value));
const description = computed(() => describeCronExpression(rawExpression.value));
const invalidMessage = computed(() => props.error || cronValidationMessageText(validation.value));
const displayValid = computed(() => validation.value.valid && !props.error);
const inputStatus = computed(() => (invalidMessage.value ? 'error' : 'default'));
const descriptionText = computed(() => translateCronDescription(description.value, t));
const showClearControl = computed(() => Boolean(rawExpression.value && !props.disabled && inputFocused.value));

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

function handleInputBlur(value: string | number) {
  inputFocused.value = false;
  handleRawInput(value);
}

function handleClearExpression() {
  applyExpression('');
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
  align-items: stretch;
  display: flex;
  gap: var(--graft-density-gap-8);
  min-width: 0;
  width: 100%;
}

.cron-expression-field__input {
  flex: 1 1 0;
  min-width: 0;
}

.cron-expression-field__clear-space {
  align-items: center;
  block-size: 1em;
  display: inline-flex;
  inline-size: var(--td-font-size-body-large);
  justify-content: center;
  position: relative;
}

.cron-expression-field__clear-button {
  align-items: center;
  appearance: none;
  background: transparent;
  border: 0;
  color: var(--td-text-color-placeholder);
  cursor: pointer;
  display: inline-flex;
  inset: 0;
  justify-content: center;
  padding: 0;
  position: absolute;
  transition: color 0.2s linear;
}

.cron-expression-field__clear-button:hover {
  color: var(--td-text-color-secondary);
}

.cron-expression-field__clear-icon {
  font-size: var(--td-font-size-body-large);
}

.cron-expression-field__configure {
  flex: 0 0 auto;
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
