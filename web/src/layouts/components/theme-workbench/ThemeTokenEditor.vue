<template>
  <div class="theme-token-editor">
    <div class="editor-toolbar">
      <div class="mode-switch">
        <span class="toolbar-label">{{ t('layout.setting.workbench.token.targetMode') }}</span>
        <t-radio-group v-model="activeMode" variant="default-filled">
          <t-radio-button value="light">{{ t('layout.setting.workbench.token.light') }}</t-radio-button>
          <t-radio-button value="dark">{{ t('layout.setting.workbench.token.dark') }}</t-radio-button>
        </t-radio-group>
      </div>
      <t-button size="small" variant="outline" @click="clearCurrentGroup">
        {{ t('layout.setting.workbench.actions.clearGroup') }}
      </t-button>
    </div>

    <div v-if="tokenDefinitions.length" class="token-grid">
      <div v-for="token in tokenDefinitions" :key="token.key" class="token-item">
        <div class="token-meta">
          <div class="token-label">{{ token.label }}</div>
          <div class="token-key">{{ token.key }}</div>
        </div>
        <div class="token-inputs">
          <label v-if="showColorInput(token.key)" class="color-input">
            <input
              type="color"
              :value="toHex(getTokenValue(token.key))"
              @input="updateToken(token.key, ($event.target as HTMLInputElement).value)"
            />
          </label>
          <div
            v-if="showPreviewSwatch(token.key)"
            class="token-preview"
            :style="{ background: getTokenValue(token.key) }"
          />
          <t-input
            :model-value="getTokenValue(token.key)"
            @update:model-value="(value) => updateToken(token.key, value)"
          />
        </div>
      </div>
    </div>

    <t-empty v-else :description="t('layout.setting.workbench.token.empty')" />
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue';

import { t } from '@/locales';
import { useSettingStore } from '@/store';
import type { ThemeTokenDefinition, ThemeTokenGroupKey } from '@/types/theme';
import type { ModeType } from '@/utils/types';

const { tokenDefinitions } = defineProps<{
  tokenDefinitions: ThemeTokenDefinition[];
  groupKey: ThemeTokenGroupKey;
}>();

const settingStore = useSettingStore();
const activeMode = ref<ModeType>('light');

const getTokenValue = (tokenKey: string) => {
  const modeTokens = settingStore.themeResolvedTokens[activeMode.value];
  return modeTokens[tokenKey] ?? '';
};

const updateToken = (tokenKey: string, tokenValue: string) => {
  if (!tokenValue) return;
  settingStore.updateThemeToken(activeMode.value, tokenKey, tokenValue);
};

const clearCurrentGroup = () => {
  settingStore.clearThemeTokenGroup(
    activeMode.value,
    tokenDefinitions.map((token) => token.key),
  );
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
.theme-token-editor {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.editor-toolbar {
  align-items: center;
  display: flex;
  gap: 12px;
  justify-content: space-between;
}

.mode-switch {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.toolbar-label {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-medium);
}

.token-grid {
  display: grid;
  gap: 12px;
}

.token-item {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-large);
  display: grid;
  gap: 12px;
  padding: 14px 16px;
}

.token-meta {
  display: grid;
  gap: 4px;
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

.token-inputs {
  align-items: center;
  display: grid;
  gap: 10px;
  grid-template-columns: auto auto 1fr;
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
  height: 28px;
  width: 28px;
}
</style>
