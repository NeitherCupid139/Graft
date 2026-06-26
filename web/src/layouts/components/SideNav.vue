<template>
  <div :class="sideNavCls">
    <t-menu
      :class="menuCls"
      :theme="theme"
      :value="active"
      :collapsed="collapsed"
      :expanded="expanded"
      :expand-mutex="menuAutoCollapsed"
      :width="menuWidth"
      @expand="onExpanded"
    >
      <template #logo>
        <button
          v-if="showLogo"
          type="button"
          :class="`${prefix}-side-nav-logo-wrapper`"
          :aria-label="t('layout.header.home')"
          @click="goHome"
        >
          <graft-brand-logo :variant="collapsed ? 'mark' : 'wordmark'" :class="logoCls" />
        </button>
      </template>
      <menu-content :nav-data="menu" />
      <template #operations>
        <div :class="versionCls">
          <span v-if="!collapsed" class="side-nav-meta__name">{{ t('common.appName') }}</span>
          <span class="side-nav-meta__version">{{ t('layout.sideNav.version', { version: appVersion }) }}</span>
        </div>
      </template>
    </t-menu>
    <div :class="`${prefix}-side-nav-placeholder${collapsed ? '-hidden' : ''}`"></div>
  </div>
</template>
<script setup lang="ts">
import difference from 'lodash/difference';
import remove from 'lodash/remove';
import union from 'lodash/union';
import type { MenuValue } from 'tdesign-vue-next';
import type { PropType } from 'vue';
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';

import { prefix } from '@/config/global';
import { useShellNavigation } from '@/layouts/useShellNavigation';
import { t } from '@/locales';
import { getActive } from '@/router';
import GraftBrandLogo from '@/shared/components/GraftBrandLogo.vue';
import { useSettingStore } from '@/store';
import type { MenuRoute, ModeType } from '@/utils/types';

import pgk from '../../../package.json';
import MenuContent from './MenuContent.vue';

const appVersion = 'version' in pgk ? String(pgk.version) : '';
const menuWidth = ['232px', '64px'];

const { menu, showLogo, isFixed, layout, theme, isCompact } = defineProps({
  menu: {
    type: Array as PropType<MenuRoute[]>,
    default: () => [],
  },
  showLogo: {
    type: Boolean as PropType<boolean>,
    default: true,
  },
  isFixed: {
    type: Boolean as PropType<boolean>,
    default: true,
  },
  layout: {
    type: String as PropType<string>,
    default: '',
  },
  headerHeight: {
    type: String as PropType<string>,
    default: '64px',
  },
  theme: {
    type: String as PropType<ModeType>,
    default: 'light',
  },
  isCompact: {
    type: Boolean as PropType<boolean>,
    default: false,
  },
});

const MIN_POINT = 992 - 1;

const collapsed = computed(() => useSettingStore().isSidebarCompact);
const menuAutoCollapsed = computed(() => useSettingStore().menuAutoCollapsed);

const active = computed(() => getActive());

const expanded = ref<MenuValue[]>([]);

const getExpanded = () => {
  const path = getActive();
  const parts = path.split('/').slice(1);
  const result = parts.map((_, index) => `/${parts.slice(0, index + 1).join('/')}`);

  expanded.value = menuAutoCollapsed.value ? result : union(result, expanded.value);
};

watch(
  () => active.value,
  () => {
    getExpanded();
  },
);

const onExpanded = (value: MenuValue[]) => {
  const currentOperationMenu = difference(expanded.value, value);
  const allExpanded = union(value, expanded.value);
  remove(allExpanded, (item) => currentOperationMenu.includes(item));
  expanded.value = allExpanded;
};

const sideMode = computed(() => {
  return theme === 'dark';
});
const sideNavCls = computed(() => {
  return [
    `${prefix}-sidebar-layout`,
    {
      [`${prefix}-sidebar-compact`]: isCompact,
    },
  ];
});
const logoCls = computed(() => {
  return [
    `${prefix}-side-nav-logo`,
    {
      [`${prefix}-side-nav-logo--compact`]: collapsed.value,
      [`${prefix}-side-nav-dark`]: sideMode.value,
    },
  ];
});
const versionCls = computed(() => {
  return [
    'side-nav-meta',
    {
      'side-nav-meta--compact': collapsed.value,
      [`${prefix}-side-nav-dark`]: sideMode.value,
    },
  ];
});
const menuCls = computed(() => {
  return [
    `${prefix}-side-nav`,
    {
      [`${prefix}-side-nav-no-logo`]: !showLogo,
      [`${prefix}-side-nav-no-fixed`]: !isFixed,
      [`${prefix}-side-nav-mix-fixed`]: layout === 'mix' && isFixed,
    },
  ];
});

const settingStore = useSettingStore();
const shellNavigation = useShellNavigation();

const autoCollapsed = () => {
  const isCompact = window.innerWidth <= MIN_POINT;
  settingStore.updateConfig({
    isSidebarCompact: isCompact,
  });
};

onMounted(() => {
  getExpanded();
  autoCollapsed();

  window.addEventListener('resize', autoCollapsed);
});

onUnmounted(() => {
  window.removeEventListener('resize', autoCollapsed);
});

const goHome = () => {
  void shellNavigation.goHome();
};

</script>
<style lang="less" scoped>
.side-nav-meta {
  align-items: flex-start;
  color: var(--td-text-color-secondary);
  display: flex;
  flex-direction: column;
  font: var(--td-font-body-small);
  gap: var(--td-comp-margin-xxs);
  line-height: 1.4;
  padding: var(--td-comp-paddingTB-xs) 0;
  width: 100%;

  &--compact {
    align-items: center;
    text-align: center;
  }

  &__name {
    color: var(--td-text-color-primary);
    font: var(--td-font-body-medium);
    opacity: 0.85;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 100%;
  }

  &__version {
    font-variant-numeric: tabular-nums;
    opacity: 0.55;
  }
}
</style>
