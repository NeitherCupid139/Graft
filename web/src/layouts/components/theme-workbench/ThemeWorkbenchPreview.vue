<template>
  <div class="theme-preview" :class="[`theme-preview--${themeState.densityPreset}`]">
    <div class="theme-preview__sticky">
      <t-card :bordered="false" class="theme-preview__panel">
        <template #header>
          <div class="theme-preview__header">
            <div>
              <div class="theme-preview__eyebrow">{{ t('layout.setting.workbench.preview.title') }}</div>
              <div class="theme-preview__title">{{ settingStore.effectiveThemeDisplayName }}</div>
            </div>
            <t-tag theme="primary" variant="light-outline">{{ themeState.mode }}</t-tag>
          </div>
        </template>

        <div class="preview-shell">
          <div class="preview-shell__nav">
            <div class="preview-shell__brand">{{ t('common.appName') }}</div>
            <div class="preview-shell__menu">
              <span class="menu-item menu-item--active">{{ t('menu.access_control.title') }}</span>
              <span class="menu-item">{{ t('menu.audit.title') }}</span>
              <span class="menu-item">{{ t('menu.server.title') }}</span>
            </div>
          </div>

          <div class="preview-shell__content">
            <div class="preview-shell__breadcrumb">Security / Audit / Access Log</div>
            <div class="preview-shell__page-header">
              <div>
                <div class="preview-shell__page-title">{{ t('layout.setting.workbench.preview.pageTitle') }}</div>
                <div class="preview-shell__page-desc">{{ t('layout.setting.workbench.preview.pageDesc') }}</div>
              </div>
              <div class="preview-shell__actions">
                <t-button variant="outline">{{ t('components.commonTable.reset') }}</t-button>
                <t-button theme="primary">{{ t('components.commonTable.query') }}</t-button>
              </div>
            </div>

            <div class="preview-toolbar">
              <t-input :model-value="t('layout.setting.workbench.preview.keyword')" />
              <t-input :model-value="t('layout.setting.workbench.preview.timeRange')" />
              <t-input :model-value="t('layout.setting.workbench.preview.status')" />
            </div>

            <div class="preview-stats">
              <t-card v-for="stat in stats" :key="stat.label" size="small" class="preview-stats__card">
                <div class="preview-stats__label">{{ stat.label }}</div>
                <div class="preview-stats__value">{{ stat.value }}</div>
              </t-card>
            </div>

            <t-alert theme="info" :message="t('layout.setting.workbench.preview.alert')" />

            <div class="preview-tags">
              <t-tag theme="success" variant="light-outline">{{ t('components.isSetup.on') }}</t-tag>
              <t-tag theme="warning" variant="light-outline">{{ t('layout.setting.workbench.preview.pending') }}</t-tag>
              <t-tag theme="danger" variant="light-outline">{{ t('layout.setting.workbench.preview.risk') }}</t-tag>
            </div>

            <t-card :title="t('layout.setting.workbench.preview.tableTitle')" size="small">
              <t-table row-key="name" size="small" :bordered="false" :data="tableData" :columns="tableColumns" />
            </t-card>

            <t-card :title="t('layout.setting.workbench.preview.logTitle')" size="small">
              <div class="preview-log">
                <div v-for="item in logs" :key="item.time" class="preview-log__item">
                  <div class="preview-log__time">{{ item.time }}</div>
                  <div class="preview-log__body">
                    <div class="preview-log__row">
                      <span class="preview-log__type">{{ item.type }}</span>
                      <span>{{ item.actor }}</span>
                    </div>
                    <div class="preview-log__detail">{{ item.detail }}</div>
                  </div>
                </div>
              </div>
            </t-card>
          </div>
        </div>
      </t-card>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue';

import { t } from '@/locales';
import { useSettingStore } from '@/store';

const settingStore = useSettingStore();

const themeState = computed(() => settingStore.effectiveThemeState);

const stats = computed(() => [
  { label: t('layout.setting.workbench.preview.totalUsers'), value: '12,480' },
  { label: t('layout.setting.workbench.preview.onlineUsers'), value: '832' },
  { label: t('layout.setting.workbench.preview.alertCount'), value: '14' },
]);

const tableColumns = computed(() => [
  { colKey: 'name', title: t('layout.setting.workbench.preview.userColumn') },
  { colKey: 'role', title: t('layout.setting.workbench.preview.roleColumn') },
  { colKey: 'status', title: t('layout.setting.workbench.preview.statusColumn') },
  { colKey: 'action', title: t('components.commonTable.operation') },
]);

const tableData = computed(() => [
  { name: 'alice', role: 'admin', status: t('components.isSetup.on'), action: t('components.manage') },
  {
    name: 'bruce',
    role: 'auditor',
    status: t('layout.setting.workbench.preview.pending'),
    action: t('components.commonTable.detail'),
  },
  { name: 'cathy', role: 'operator', status: t('components.isSetup.off'), action: t('components.manage') },
]);

const logs = computed(() => [
  {
    time: '2026-06-02 09:12',
    type: 'AUTH',
    actor: 'alice',
    detail: t('layout.setting.workbench.preview.logAuth'),
  },
  {
    time: '2026-06-02 09:20',
    type: 'AUDIT',
    actor: 'bruce',
    detail: t('layout.setting.workbench.preview.logAudit'),
  },
  {
    time: '2026-06-02 09:31',
    type: 'ACCESS',
    actor: 'gateway',
    detail: t('layout.setting.workbench.preview.logAccess'),
  },
]);
</script>
<style lang="less" scoped>
@import '../../../shared/components/management/card-surface.less';

.theme-preview {
  min-width: 0;
}

.theme-preview__sticky {
  position: sticky;
  top: 16px;
}

.theme-preview__panel {
  .management-card-surface();
}

.theme-preview__header {
  align-items: center;
  display: flex;
  gap: 12px;
  justify-content: space-between;
}

.theme-preview__eyebrow {
  color: var(--td-text-color-placeholder);
  font-size: 12px;
}

.theme-preview__title {
  color: var(--td-text-color-primary);
  font-size: 18px;
  font-weight: 700;
}

.preview-shell {
  display: grid;
  gap: 16px;
}

.preview-shell__nav {
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-large);
  padding: 14px 16px;
}

.preview-shell__brand {
  font-size: 16px;
  font-weight: 700;
}

.preview-shell__menu {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.menu-item {
  border-radius: var(--td-radius-default);
  color: var(--td-text-color-secondary);
  padding: 6px 10px;
}

.menu-item--active {
  background: color-mix(in srgb, var(--td-brand-color) 12%, transparent);
  color: var(--td-brand-color);
}

.preview-shell__content {
  display: grid;
  gap: 16px;
}

.preview-shell__breadcrumb,
.preview-shell__page-desc,
.preview-stats__label,
.preview-log__detail {
  color: var(--td-text-color-secondary);
}

.preview-shell__page-header {
  align-items: start;
  display: flex;
  gap: 12px;
  justify-content: space-between;
}

.preview-shell__page-title {
  font-size: 20px;
  font-weight: 700;
}

.preview-shell__actions,
.preview-toolbar,
.preview-tags {
  display: flex;
  gap: 12px;
}

.preview-toolbar > * {
  flex: 1;
}

.preview-stats {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.preview-stats__card {
  box-shadow: var(--td-shadow-1);
}

.preview-stats__value {
  font-size: 22px;
  font-weight: 700;
  margin-top: 8px;
}

.preview-log {
  display: grid;
  gap: 12px;
}

.preview-log__item {
  border-bottom: 1px solid var(--td-component-stroke);
  display: grid;
  gap: 10px;
  grid-template-columns: 132px 1fr;
  padding: 10px 0;
}

.preview-log__item:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.preview-log__time {
  color: var(--td-text-color-placeholder);
  font-size: 12px;
}

.preview-log__row {
  display: flex;
  font-weight: 600;
  gap: 10px;
}

.preview-log__type {
  color: var(--td-brand-color);
}

.theme-preview--compact {
  --preview-gap: 10px;
}

.theme-preview--compact .preview-shell,
.theme-preview--compact .preview-shell__content,
.theme-preview--compact .preview-toolbar,
.theme-preview--compact .preview-tags,
.theme-preview--compact .preview-stats,
.theme-preview--compact .preview-log {
  gap: 10px;
}

.theme-preview--comfortable .preview-shell,
.theme-preview--comfortable .preview-shell__content {
  gap: 20px;
}

@media (width <= 1024px) {
  .theme-preview__sticky {
    position: static;
  }

  .preview-stats {
    grid-template-columns: 1fr;
  }

  .preview-toolbar,
  .preview-shell__actions,
  .preview-tags,
  .preview-shell__menu,
  .preview-shell__page-header,
  .preview-log__item {
    flex-direction: column;
    grid-template-columns: 1fr;
  }
}
</style>
