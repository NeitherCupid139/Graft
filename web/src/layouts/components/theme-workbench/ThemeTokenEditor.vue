<template>
  <div class="theme-token-editor">
    <div class="editor-toolbar">
      <t-button size="small" variant="text" theme="danger" class="clear-button" @click="clearCurrentGroup">
        {{ t('layout.setting.workbench.actions.clearGroup') }}
      </t-button>
    </div>

    <div v-if="tokenDefinitions.length" class="token-grid">
      <div v-for="token in tokenDefinitions" :key="token.key" class="token-item">
        <div class="token-meta">
          <div class="token-label">{{ token.label }}</div>
          <div class="token-key">{{ token.key }}</div>
        </div>
        <div class="token-preview-rail">
          <div
            v-if="showPreviewSwatch(token.key)"
            class="token-preview"
            :style="{ background: getResolvedTokenValue(token.key) }"
          />
          <div v-else class="token-preview token-preview--text">
            <span>Aa</span>
          </div>
          <div class="token-preview-sample">
            <span class="token-preview-sample__line" />
            <span class="token-preview-sample__line token-preview-sample__line--short" />
          </div>
        </div>
        <div class="token-inputs">
          <label v-if="showColorInput(token.key)" class="color-input">
            <input
              type="color"
              :value="toHex(getInputValue(token.key))"
              @input="updateToken(token.key, ($event.target as HTMLInputElement).value)"
            />
          </label>
          <t-input
            class="token-input"
            :model-value="getInputValue(token.key)"
            @update:model-value="(value) => updateDraftValue(token.key, String(value ?? ''))"
            @change="(value) => commitToken(token.key, String(value ?? ''))"
            @blur="() => commitToken(token.key)"
          />
          <t-button
            size="small"
            variant="text"
            class="reset-button"
            :disabled="!hasTokenOverride(token.key)"
            @click="resetToken(token.key)"
          >
            {{ t('layout.setting.workbench.token.reset') }}
          </t-button>
        </div>
      </div>
    </div>

    <t-empty v-else :description="t('layout.setting.workbench.token.empty')" />
  </div>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue';

import { t } from '@/locales';
import { useSettingStore } from '@/store';
import type { ThemeTokenDefinition, ThemeTokenGroupKey } from '@/types/theme';
import type { ModeType } from '@/utils/types';

const props = defineProps<{
  tokenDefinitions: ThemeTokenDefinition[];
  groupKey: ThemeTokenGroupKey;
  mode: ModeType;
}>();

const settingStore = useSettingStore();
const draftValues = ref<Record<string, string>>({});

watch(
  () => [props.groupKey, props.mode],
  () => {
    draftValues.value = {};
  },
);

const getResolvedTokenValue = (tokenKey: string) => {
  const modeTokens = settingStore.themeResolvedTokens[props.mode];
  return modeTokens[tokenKey] ?? '';
};

const getInputValue = (tokenKey: string) => {
  return draftValues.value[tokenKey] ?? getResolvedTokenValue(tokenKey);
};

const updateDraftValue = (tokenKey: string, tokenValue: string) => {
  draftValues.value = {
    ...draftValues.value,
    [tokenKey]: tokenValue,
  };
};

const hasTokenOverride = (tokenKey: string) => {
  return Object.prototype.hasOwnProperty.call(settingStore.themeTokenOverrides[props.mode], tokenKey);
};

const resetToken = (tokenKey: string) => {
  const nextDraftValues = { ...draftValues.value };
  delete nextDraftValues[tokenKey];
  draftValues.value = nextDraftValues;
  settingStore.clearThemeTokenGroup(props.mode, [tokenKey]);
};

const commitToken = (tokenKey: string, tokenValue?: string) => {
  const resolvedValue = (tokenValue ?? getInputValue(tokenKey)).trim();

  if (!resolvedValue) {
    resetToken(tokenKey);
    return;
  }

  settingStore.updateThemeToken(props.mode, tokenKey, resolvedValue);
  const nextDraftValues = { ...draftValues.value };
  delete nextDraftValues[tokenKey];
  draftValues.value = nextDraftValues;
};

const updateToken = (tokenKey: string, tokenValue: string) => {
  updateDraftValue(tokenKey, tokenValue);
  commitToken(tokenKey, tokenValue);
};

const clearCurrentGroup = () => {
  settingStore.clearThemeTokenGroup(
    props.mode,
    props.tokenDefinitions.map((token) => token.key),
  );
  draftValues.value = {};
};

const showPreviewSwatch = (tokenKey: string) => /color|background|border/i.test(tokenKey);
const showColorInput = (tokenKey: string) => showPreviewSwatch(tokenKey) && !/shadow/i.test(tokenKey);

const toHex = (value: string) => {
  if (!value) return '#0052d9';

  const canvas = document.createElement('canvas');
  const context = canvas.getContext('2d');
  if (!context) return '#0052d9';

  context.fillStyle = value;
  const resolved = context.fillStyle;
  if (resolved.startsWith('#')) {
    return resolved;
  }

  const matched = resolved.match(/\d+/g);
  if (!matched?.length) {
    return '#0052d9';
  }

  const [red, green, blue] = matched.slice(0, 3).map((item) => Number(item).toString(16).padStart(2, '0'));
  return `#${red}${green}${blue}`;
};
</script>
<style lang="less" scoped>
@import './theme-surface.less';

.theme-token-editor {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.editor-toolbar {
  align-items: center;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.clear-button {
  opacity: 0.72;
}

.token-grid {
  display: grid;
  gap: 12px;
}

.token-item {
  .theme-workbench-surface();

  align-items: flex-start;
  gap: 12px;
  grid-template-columns: minmax(0, 1.3fr) 112px minmax(220px, 0.9fr);
  padding: 14px 16px;
}

.token-meta {
  display: grid;
  gap: 4px;
  min-width: 0;
}

.token-label {
  color: var(--td-text-color-primary);
  font: var(--td-font-body-medium);
  font-weight: 600;
}

.token-key {
  color: var(--td-text-color-placeholder);
  font-size: 12px;
  word-break: break-all;
}

.token-preview-rail {
  align-items: center;
  background: color-mix(in srgb, var(--td-bg-color-page) 84%, var(--td-bg-color-container));
  border: 1px solid color-mix(in srgb, var(--td-component-stroke) 88%, transparent);
  border-radius: 12px;
  display: grid;
  gap: 10px;
  grid-template-columns: 32px minmax(0, 1fr);
  min-height: 48px;
  padding: 8px 10px;
}

.token-inputs {
  align-items: center;
  display: grid;
  gap: 10px;
  grid-template-columns: auto minmax(160px, 188px) auto;
  justify-content: end;
  min-width: 0;
}

.token-input {
  min-width: 0;
}

.color-input {
  align-items: center;
  display: inline-flex;
}

.color-input input {
  appearance: none;
  background: transparent;
  border: 0;
  cursor: pointer;
  height: 32px;
  padding: 0;
  width: 32px;
}

.token-preview {
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-default);
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 14%),
    0 4px 14px rgb(15 23 42 / 10%);
  height: 32px;
  width: 32px;
}

.token-preview--text {
  align-items: center;
  background: var(--td-bg-color-container);
  color: var(--td-text-color-primary);
  display: inline-flex;
  font-size: 14px;
  font-weight: 700;
  justify-content: center;
}

.token-preview-sample {
  display: grid;
  gap: 6px;
}

.token-preview-sample__line {
  background: color-mix(in srgb, var(--td-brand-color) 12%, var(--td-text-color-placeholder));
  border-radius: 999px;
  display: block;
  height: 7px;
  width: 100%;
}

.token-preview-sample__line--short {
  width: 72%;
}

.reset-button {
  opacity: 0.76;
}

@media (width <= 768px) {
  .editor-toolbar {
    justify-content: flex-start;
  }

  .token-item {
    grid-template-columns: 1fr;
  }

  .token-inputs {
    grid-template-columns: auto 1fr auto;
  }
}
</style>
