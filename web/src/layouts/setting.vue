<template>
  <theme-workbench-dock v-if="showFloatingWorkbench" />
  <theme-workbench-panel />
</template>
<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';

import { useSettingStore } from '@/store';

import ThemeWorkbenchDock from './components/theme-workbench/ThemeWorkbenchDock.vue';
import ThemeWorkbenchPanel from './components/theme-workbench/ThemeWorkbenchPanel.vue';

const route = useRoute();
const settingStore = useSettingStore();

const showFloatingWorkbench = computed(() => route.path !== '/login');

onMounted(() => {
  // 统一在工作台组件挂载时恢复主题状态，避免后台壳和登录页各自维护一套入口状态。
  settingStore.initializeThemeWorkbenchRuntime();
  if (settingStore.showSettingPanel && !settingStore.showThemeWorkbench) {
    settingStore.openThemeWorkbench('overview');
  }
});
</script>
