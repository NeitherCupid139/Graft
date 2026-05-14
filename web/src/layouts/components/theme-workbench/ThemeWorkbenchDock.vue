<template>
  <div class="theme-workbench-dock">
    <t-button class="theme-workbench-dock__main" size="large" variant="outline" @click="openGroup('overview')">
      <template #icon>
        <t-icon name="palette" />
      </template>
      {{ t('layout.setting.workbench.dock.title') }}
    </t-button>
    <t-button shape="circle" variant="outline" @click="openGroup('brand')">
      <t-icon name="edit-1" />
    </t-button>
    <t-button shape="circle" variant="outline" @click="openGroup('font')">
      <t-icon name="textformat" />
    </t-button>
    <t-button shape="circle" variant="outline" @click="openGroup('radius')">
      <t-icon name="layers" />
    </t-button>
    <t-button shape="circle" variant="outline" @click="resetWorkbench">
      <t-icon name="rollback" />
    </t-button>
  </div>
</template>
<script setup lang="ts">
import { t } from '@/locales';
import { useSettingStore } from '@/store';
import type { ThemeWorkbenchGroupKey } from '@/types/theme';

const settingStore = useSettingStore();

const openGroup = (group: ThemeWorkbenchGroupKey) => {
  settingStore.openThemeWorkbench(group);
};

const resetWorkbench = () => {
  settingStore.resetThemeWorkbench();
};
</script>
<style lang="less" scoped>
.theme-workbench-dock {
  align-items: center;
  backdrop-filter: blur(14px);
  background: color-mix(in srgb, var(--td-bg-color-container) 92%, transparent);
  border: 1px solid color-mix(in srgb, var(--td-component-stroke) 70%, transparent);
  border-radius: 999px;
  bottom: 28px;
  box-shadow: 0 14px 30px rgb(15 23 42 / 14%);
  display: inline-flex;
  gap: 10px;
  left: 50%;
  padding: 10px 12px;
  position: fixed;
  transform: translateX(-50%);
  z-index: 1100;
}

.theme-workbench-dock__main {
  min-width: 156px;
}

:deep(.t-button--variant-outline) {
  backdrop-filter: blur(10px);
}
</style>
