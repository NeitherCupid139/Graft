<template>
  <header class="login-header">
    <button type="button" class="login-header__brand" :aria-label="t('layout.header.home')" @click="goHome">
      <graft-brand-logo variant="wordmark" />
    </button>
    <div class="operations-container">
      <shell-header-resource-actions />
      <language-switcher />
      <t-tooltip placement="bottom" :content="t('layout.header.setting')">
        <t-button theme="default" shape="square" variant="text" @click="toggleSettingPanel">
          <t-icon name="setting" class="icon" />
        </t-button>
      </t-tooltip>
    </div>
  </header>
</template>
<script setup lang="ts">
import { useShellNavigation } from '@/layouts/useShellNavigation';
import { t } from '@/locales';
import GraftBrandLogo from '@/shared/components/GraftBrandLogo.vue';
import LanguageSwitcher from '@/shared/components/LanguageSwitcher.vue';
import ShellHeaderResourceActions from '@/shared/components/ShellHeaderResourceActions.vue';
import { useSettingStore } from '@/store';

const settingStore = useSettingStore();
const { goHome } = useShellNavigation();

const toggleSettingPanel = () => {
  settingStore.openThemeWorkbench('overview');
};
</script>
<style lang="less" scoped>
.login-header {
  align-items: center;
  backdrop-filter: blur(10px);
  color: var(--td-text-color-primary);
  display: flex;
  height: var(--td-comp-size-xxxl);
  justify-content: space-between;
  padding: 0 var(--td-comp-paddingLR-xl);

  &__brand {
    align-items: center;
    background: transparent;
    border: 0;
    color: var(--td-brand-color);
    cursor: pointer;
    display: inline-flex;
    max-width: 11.5rem;
    padding: 0;
    transition: color 0.2s ease, opacity 0.2s ease;

    &:hover {
      color: var(--td-brand-color-hover);
      opacity: 0.92;
    }

    &:focus-visible {
      outline: 2px solid var(--td-brand-color-focus);
      outline-offset: 2px;
    }
  }

  .operations-container {
    align-items: center;
    display: flex;
    gap: var(--td-comp-margin-xs);

    .t-button {
      margin-left: 0;
    }
  }
}
</style>