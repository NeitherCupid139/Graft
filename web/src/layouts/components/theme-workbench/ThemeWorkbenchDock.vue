<template>
  <div class="theme-workbench-dock">
    <t-button
      class="theme-workbench-dock__main"
      :class="{ 'theme-workbench-dock__main--active': settingStore.showThemeWorkbench }"
      :title="dockMainTitle"
      variant="outline"
      @click="toggleOverview"
    >
      <template #icon>
        <t-icon name="palette" size="20px" />
      </template>
      <span class="theme-workbench-dock__action-label">{{ t('layout.setting.workbench.dock.title') }}</span>
    </t-button>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue';

import { t } from '@/locales';
import { useSettingStore } from '@/store';

const settingStore = useSettingStore();

const dockMainTitle = computed(() => {
  if (settingStore.showThemeWorkbench && settingStore.activeThemeWorkbenchGroup === 'overview') {
    return undefined;
  }

  return t('layout.setting.workbench.dock.title');
});

const toggleOverview = () => {
  if (settingStore.showThemeWorkbench) {
    settingStore.cancelThemeDraft();
    return;
  }

  settingStore.openThemeWorkbench('overview');
};
</script>
<style lang="less" scoped>
.theme-workbench-dock {
  align-items: center;
  backdrop-filter: blur(22px) saturate(155%);
  background:
    linear-gradient(135deg, rgb(255 255 255 / 58%), rgb(255 255 255 / 20%)),
    color-mix(in srgb, var(--td-bg-color-container) 84%, transparent);
  border: 1px solid color-mix(in srgb, var(--td-component-stroke) 52%, rgb(255 255 255 / 46%));
  border-radius: 28px;
  bottom: calc(24px + env(safe-area-inset-bottom, 0px));
  box-shadow:
    0 14px 34px rgb(15 23 42 / 12%),
    inset 0 1px 0 rgb(255 255 255 / 42%);
  box-sizing: border-box;
  display: inline-flex;
  flex-wrap: nowrap;
  justify-content: center;
  left: 50%;
  max-width: calc(100vw - 24px);
  padding: 6px;
  position: fixed;
  transform: translateX(-50%);
  width: max-content;
  z-index: 1090;
}

.theme-workbench-dock__main {
  flex: 0 0 auto;
  min-width: 48px;
}

.theme-workbench-dock__main--active {
  background: color-mix(in srgb, var(--td-brand-color) 10%, var(--td-bg-color-container));
  border-color: color-mix(in srgb, var(--td-brand-color) 24%, transparent);
  box-shadow:
    0 8px 18px color-mix(in srgb, var(--td-brand-color) 12%, transparent),
    inset 0 1px 0 rgb(255 255 255 / 28%);
  color: var(--td-brand-color);
}

:deep(.t-button--variant-outline) {
  backdrop-filter: blur(14px);
  background: color-mix(in srgb, var(--td-bg-color-container) 78%, transparent);
  border-color: color-mix(in srgb, var(--td-component-stroke) 44%, transparent);
  box-shadow:
    0 6px 16px rgb(15 23 42 / 7%),
    inset 0 1px 0 rgb(255 255 255 / 28%);
  color: var(--td-text-color-secondary);
  transition:
    transform 0.18s ease,
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    background-color 0.18s ease,
    color 0.18s ease;
}

:deep(.t-button--variant-outline:hover) {
  background: color-mix(in srgb, var(--td-bg-color-container) 90%, transparent);
  border-color: color-mix(in srgb, var(--td-brand-color) 14%, var(--td-component-stroke));
  color: var(--td-text-color-primary);
  transform: translateY(-1px);
}

:deep(.theme-workbench-dock__main.t-button) {
  align-items: center;
  border-radius: 999px;
  display: inline-flex;
  font-weight: 600;
  height: 48px;
  justify-content: center;
  overflow: hidden;
  padding-inline: 0;
  transition:
    min-width 0.22s ease,
    padding-inline 0.22s ease,
    background-color 0.18s ease,
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    color 0.18s ease,
    transform 0.18s ease;
}

:deep(.theme-workbench-dock__main .t-button__content) {
  align-items: center;
  display: inline-flex;
  height: 100%;
  justify-content: center;
  width: 100%;
}

:deep(.theme-workbench-dock__main .t-button__prefix) {
  align-items: center;
  display: inline-flex;
  height: 20px;
  justify-content: center;
  line-height: 1;
  margin-right: 0;
  width: 20px;
}

:deep(.theme-workbench-dock__main .t-icon) {
  display: block;
  flex: 0 0 auto;
}

:deep(.theme-workbench-dock__main .t-button__text) {
  margin-left: 0;
  max-width: 0;
  opacity: 0;
  overflow: hidden;
  transition:
    max-width 0.22s ease,
    margin-left 0.22s ease,
    opacity 0.16s ease;
  white-space: nowrap;
}

:deep(.theme-workbench-dock__main--active.t-button) {
  min-width: 118px;
  padding-inline: 16px;
}

:deep(.theme-workbench-dock__main--active .t-button__prefix) {
  margin-right: 8px;
}

:deep(.theme-workbench-dock__main--active .t-button__text) {
  max-width: 88px;
  opacity: 1;
}

@media (width <= 768px) {
  .theme-workbench-dock {
    bottom: 16px;
    padding: 6px;
  }

  :deep(.theme-workbench-dock__main.t-button) {
    height: 44px;
    min-width: 44px;
  }

  :deep(.theme-workbench-dock__main--active.t-button) {
    min-width: 108px;
    padding-inline: 14px;
  }
}
</style>
