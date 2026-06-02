<template>
  <t-drawer
    v-model:visible="drawerVisible"
    class="theme-workbench-panel"
    destroy-on-close
    placement="right"
    size="620px"
    :close-btn="false"
    :footer="false"
    :header="false"
  >
    <div class="theme-workbench-panel__shell">
      <header class="theme-workbench-panel__header">
        <div>
          <div class="panel-title">{{ t('layout.setting.workbench.title') }}</div>
          <div class="panel-subtitle">{{ t('layout.setting.workbench.subtitle') }}</div>
        </div>
        <t-button shape="square" variant="text" @click="closeWorkbench">
          <t-icon name="close" />
        </t-button>
      </header>

      <div class="theme-workbench-panel__body">
        <aside class="theme-workbench-panel__nav">
          <button
            v-for="group in groups"
            :key="group.key"
            type="button"
            class="nav-item"
            :class="{ 'nav-item--active': group.key === activeGroup }"
            @click="openGroup(group.key)"
          >
            <span class="nav-item__icon">
              <t-icon :name="groupIconMap[group.key]" />
            </span>
            <span class="nav-item__text">{{ t(group.labelKey) }}</span>
          </button>
        </aside>

        <section class="theme-workbench-panel__content">
          <template v-if="activeGroup === 'overview'">
            <div class="section">
              <div class="section-title">{{ t('layout.setting.workbench.groups.overview') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.overview.description') }}</div>
              <div class="summary-grid">
                <div class="summary-card">
                  <div class="summary-card__label">{{ t('layout.setting.workbench.overview.currentTheme') }}</div>
                  <div class="summary-card__value">{{ themeIdentityTitle }}</div>
                </div>
                <div class="summary-card">
                  <div class="summary-card__label">{{ t('layout.setting.navigationLayout') }}</div>
                  <div class="summary-card__value">{{ layoutLabel }}</div>
                </div>
                <div class="summary-card">
                  <div class="summary-card__label">{{ t('layout.setting.theme.mode') }}</div>
                  <div class="summary-card__value">{{ modeLabel }}</div>
                </div>
                <div class="summary-card">
                  <div class="summary-card__label">{{ t('layout.setting.workbench.overview.modifiedCount') }}</div>
                  <div class="summary-card__value">{{ themeIdentity.modifiedCount }}</div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="activeGroup === 'appearance'">
            <div class="section">
              <div class="section-title">{{ t('layout.setting.theme.mode') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.appearance.description') }}</div>
              <div class="choice-grid choice-grid--mode">
                <button
                  v-for="item in modeOptions"
                  :key="item.type"
                  type="button"
                  class="choice-card choice-card--mode"
                  :class="{ 'choice-card--active': effectiveTheme.mode === item.type }"
                  @click="settingStore.updateThemeDraftAppearance({ mode: item.type })"
                >
                  <span class="choice-card__check">
                    <t-icon v-if="effectiveTheme.mode === item.type" name="check" />
                  </span>
                  <span class="mode-thumbnail" :class="`mode-thumbnail--${item.type}`">
                    <component :is="item.icon" class="mode-thumbnail__icon" />
                  </span>
                  <span class="choice-card__title">{{ item.text }}</span>
                </button>
              </div>
            </div>

            <div class="section">
              <div class="section-title">{{ t('layout.setting.theme.color') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.appearance.colorDescription') }}</div>
              <div class="brand-palette">
                <button
                  v-for="color in brandOptions"
                  :key="color"
                  type="button"
                  class="brand-palette__item"
                  :class="{ 'brand-palette__item--active': effectiveTheme.brandTheme === color }"
                  :style="{ background: color }"
                  @click="settingStore.setCustomBrandTheme(color)"
                >
                  <t-icon v-if="effectiveTheme.brandTheme === color" name="check" />
                </button>
              </div>
              <div class="brand-input">
                <input
                  type="color"
                  :value="effectiveTheme.brandTheme"
                  @input="settingStore.setCustomBrandTheme(($event.target as HTMLInputElement).value)"
                />
                <t-input
                  :model-value="effectiveTheme.brandTheme"
                  @update:model-value="(value) => settingStore.setCustomBrandTheme(String(value ?? ''))"
                />
              </div>
            </div>
          </template>

          <template v-else-if="activeGroup === 'layout'">
            <div class="section">
              <div class="section-title">{{ t('layout.setting.navigationLayout') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.layout.tip') }}</div>
              <div class="choice-grid">
                <button
                  v-for="item in layoutOptions"
                  :key="item.value"
                  type="button"
                  class="choice-card"
                  :class="{ 'choice-card--active': settingStore.layout === item.value }"
                  @click="settingStore.updateConfig({ layout: item.value })"
                >
                  <span class="choice-card__check">
                    <t-icon v-if="settingStore.layout === item.value" name="check" />
                  </span>
                  <span class="layout-thumbnail" :class="`layout-thumbnail--${item.value}`">
                    <span class="layout-thumbnail__header" />
                    <span class="layout-thumbnail__sidebar" />
                    <span class="layout-thumbnail__content" />
                  </span>
                  <span class="choice-card__title">{{ item.label }}</span>
                </button>
              </div>
            </div>

            <div class="section">
              <div class="section-title">{{ t('layout.setting.workbench.layout.navigationBehavior') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.layout.behaviorDescription') }}</div>
              <div class="switch-list">
                <div class="switch-item">
                  <div>
                    <div class="switch-item__label">{{ t('layout.setting.splitMenuShort') }}</div>
                    <div class="switch-item__hint">{{ splitMenuHint }}</div>
                  </div>
                  <t-switch
                    :model-value="settingStore.splitMenu"
                    :disabled="!splitMenuAvailable"
                    @update:model-value="(value) => settingStore.updateConfig({ splitMenu: value })"
                  />
                </div>
                <div class="switch-item">
                  <div>
                    <div class="switch-item__label">{{ t('layout.setting.fixedSidebar') }}</div>
                    <div class="switch-item__hint">{{ fixedSidebarHint }}</div>
                  </div>
                  <t-switch
                    :model-value="settingStore.isSidebarFixed"
                    :disabled="!fixedSidebarAvailable"
                    @update:model-value="(value) => settingStore.updateConfig({ isSidebarFixed: value })"
                  />
                </div>
                <div class="switch-item">
                  <div>
                    <div class="switch-item__label">{{ t('layout.setting.element.menuAutoCollapsed') }}</div>
                    <div class="switch-item__hint">
                      {{ t('layout.setting.workbench.layout.menuAutoCollapsedHint') }}
                    </div>
                  </div>
                  <t-switch
                    :model-value="settingStore.menuAutoCollapsed"
                    @update:model-value="(value) => settingStore.updateConfig({ menuAutoCollapsed: value })"
                  />
                </div>
              </div>
            </div>

            <div class="section">
              <div class="section-title">{{ t('layout.setting.element.title') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.layout.elementsDescription') }}</div>
              <div class="switch-list">
                <div class="switch-item">
                  <div class="switch-item__label">{{ t('layout.setting.element.showHeader') }}</div>
                  <t-switch
                    :model-value="settingStore.showHeader"
                    @update:model-value="(value) => settingStore.updateConfig({ showHeader: value })"
                  />
                </div>
                <div class="switch-item">
                  <div class="switch-item__label">{{ t('layout.setting.element.showBreadcrumb') }}</div>
                  <t-switch
                    :model-value="settingStore.showBreadcrumb"
                    @update:model-value="(value) => settingStore.updateConfig({ showBreadcrumb: value })"
                  />
                </div>
                <div v-if="footerOptionVisible" class="switch-item">
                  <div class="switch-item__label">{{ t('layout.setting.element.showFooter') }}</div>
                  <t-switch
                    :model-value="settingStore.showFooter"
                    @update:model-value="(value) => settingStore.updateConfig({ showFooter: value })"
                  />
                </div>
                <div class="switch-item">
                  <div class="switch-item__label">{{ t('layout.setting.element.useTagTabs') }}</div>
                  <t-switch
                    :model-value="settingStore.isUseTabsRouter"
                    @update:model-value="(value) => settingStore.updateConfig({ isUseTabsRouter: value })"
                  />
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="activeGroup === 'typography'">
            <div class="section">
              <div class="section-title">{{ t('layout.setting.workbench.typography.fontFamily') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.typography.description') }}</div>
              <div
                class="font-option-list"
                role="radiogroup"
                :aria-label="t('layout.setting.workbench.typography.fontFamily')"
              >
                <label
                  v-for="item in fontFamilyOptions"
                  :key="item.value"
                  class="font-option"
                  :class="{ 'font-option--active': effectiveTheme.fontFamilyPreset === item.value }"
                >
                  <input
                    class="font-option__input"
                    type="radio"
                    name="theme-font-family"
                    :value="item.value"
                    :checked="effectiveTheme.fontFamilyPreset === item.value"
                    @change="
                      settingStore.updateThemeDraftAppearance({
                        fontFamilyPreset: item.value,
                      })
                    "
                  />
                  <span class="font-option__main">
                    <span class="font-option__title">{{ item.label }}</span>
                    <span class="font-option__preview" :style="{ fontFamily: item.previewFamily }">
                      {{ t('layout.setting.workbench.typography.previewLine') }}
                    </span>
                  </span>
                  <span class="font-option__check">
                    <t-icon v-if="effectiveTheme.fontFamilyPreset === item.value" name="check" />
                  </span>
                </label>
              </div>
              <div class="font-preview-strip" :style="{ fontFamily: activeFontPreviewFamily }">
                <span class="font-preview-strip__label">{{
                  t('layout.setting.workbench.typography.sampleTitle')
                }}</span>
                <span class="font-preview-strip__content">{{
                  t('layout.setting.workbench.typography.previewLine')
                }}</span>
              </div>
            </div>
          </template>

          <template v-else-if="activeGroup === 'style'">
            <div class="section">
              <div class="section-title">{{ t('layout.setting.workbench.style.radius') }}</div>
              <div class="section-desc">{{ t('layout.setting.workbench.style.description') }}</div>
              <div class="style-preview-grid">
                <button
                  v-for="item in radiusOptions"
                  :key="item.value"
                  type="button"
                  class="style-preview-card"
                  :class="{ 'style-preview-card--active': effectiveTheme.radiusPreset === item.value }"
                  @click="
                    settingStore.updateThemeDraftAppearance({
                      radiusPreset: item.value,
                    })
                  "
                >
                  <span class="style-preview-card__label">{{ item.label }}</span>
                  <span class="radius-preview" :class="`radius-preview--${item.value}`">
                    <span class="radius-preview__surface radius-preview__surface--main" />
                    <span class="radius-preview__surface radius-preview__surface--sub" />
                  </span>
                </button>
              </div>
            </div>

            <div class="section">
              <div class="section-title">{{ t('layout.setting.workbench.style.shadow') }}</div>
              <div class="style-preview-grid style-preview-grid--triple">
                <button
                  v-for="item in shadowOptions"
                  :key="item.value"
                  type="button"
                  class="style-preview-card"
                  :class="{ 'style-preview-card--active': effectiveTheme.shadowPreset === item.value }"
                  @click="
                    settingStore.updateThemeDraftAppearance({
                      shadowPreset: item.value,
                    })
                  "
                >
                  <span class="style-preview-card__label">{{ item.label }}</span>
                  <span class="shadow-preview" :class="`shadow-preview--${item.value}`">
                    <span class="shadow-preview__card shadow-preview__card--back" />
                    <span class="shadow-preview__card shadow-preview__card--front" />
                  </span>
                </button>
              </div>
            </div>

            <div class="section">
              <div class="section-title">{{ t('layout.setting.workbench.style.density') }}</div>
              <div class="style-preview-grid style-preview-grid--triple">
                <button
                  v-for="item in densityOptions"
                  :key="item.value"
                  type="button"
                  class="style-preview-card"
                  :class="{ 'style-preview-card--active': effectiveTheme.densityPreset === item.value }"
                  @click="
                    settingStore.updateThemeDraftAppearance({
                      densityPreset: item.value,
                    })
                  "
                >
                  <span class="style-preview-card__label">{{ item.label }}</span>
                  <span class="density-preview" :class="`density-preview--${item.value}`">
                    <span class="density-preview__row">
                      <span class="density-preview__bullet" />
                      <span class="density-preview__line density-preview__line--long" />
                    </span>
                    <span class="density-preview__row">
                      <span class="density-preview__bullet" />
                      <span class="density-preview__line" />
                    </span>
                    <span class="density-preview__row">
                      <span class="density-preview__bullet" />
                      <span class="density-preview__line density-preview__line--short" />
                    </span>
                  </span>
                </button>
              </div>
            </div>
          </template>

          <template v-else>
            <div class="advanced-settings-card">
              <div class="advanced-settings-card__content">
                <div class="section-title">{{ t('layout.setting.workbench.advanced.title') }}</div>
                <div class="section-desc advanced-settings-card__desc">
                  {{ t('layout.setting.workbench.advanced.description') }}
                </div>
              </div>
              <t-switch :model-value="advancedVisible" @update:model-value="toggleAdvancedVisible" />
            </div>

            <div v-if="advancedVisible" class="advanced-groups">
              <section
                v-for="group in tokenGroups"
                :key="group.value"
                class="advanced-group"
                :class="{ 'advanced-group--expanded': expandedAdvancedGroups.includes(group.value) }"
              >
                <button type="button" class="advanced-group__header" @click="toggleAdvancedGroup(group.value)">
                  <span class="advanced-group__title-block">
                    <span class="advanced-group__title">{{ group.label }}</span>
                    <span class="advanced-group__mode-tag">
                      {{ activeTokenEditorModeLabel }}
                    </span>
                  </span>
                  <span class="advanced-group__header-actions">
                    <span class="advanced-group__count">{{ tokenDefinitionsByGroup[group.value].length }}</span>
                    <t-icon name="chevron-down" class="advanced-group__chevron" />
                  </span>
                </button>
                <div v-if="expandedAdvancedGroups.includes(group.value)" class="advanced-group__body">
                  <theme-token-editor
                    :group-key="group.value"
                    :token-definitions="tokenDefinitionsByGroup[group.value]"
                    @mode-change="handleTokenEditorModeChange"
                  />
                </div>
              </section>
            </div>
          </template>

          <theme-workbench-preset-section
            v-if="showPresetSection"
            :title="t('layout.setting.workbench.presets.title')"
            :presets="presetDefinitions"
            :active-preset-id="effectivePresetId"
            @select="settingStore.selectThemePreset"
          />
        </section>
      </div>

      <footer class="theme-workbench-panel__footer">
        <t-button variant="outline" @click="settingStore.resetThemeDraftToDefault()">
          {{ t('layout.setting.workbench.actions.reset') }}
        </t-button>
        <div class="theme-workbench-panel__footer-actions">
          <t-button variant="outline" @click="settingStore.cancelThemeDraft()">
            {{ t('layout.setting.workbench.actions.cancel') }}
          </t-button>
          <t-button
            :theme="hasPendingChanges ? 'primary' : 'default'"
            :disabled="!hasPendingChanges"
            @click="settingStore.applyThemeDraft()"
          >
            {{ t('layout.setting.workbench.actions.apply') }}
          </t-button>
        </div>
      </footer>
    </div>
  </t-drawer>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';

import SettingAutoIcon from '@/assets/assets-setting-auto.svg';
import SettingDarkIcon from '@/assets/assets-setting-dark.svg';
import SettingLightIcon from '@/assets/assets-setting-light.svg';
import { DEFAULT_COLOR_OPTIONS } from '@/config/color';
import { t } from '@/locales';
import { useSettingStore } from '@/store';
import type { ThemeTokenGroupKey, ThemeWorkbenchGroupKey } from '@/types/theme';

import ThemeTokenEditor from './ThemeTokenEditor.vue';
import ThemeWorkbenchPresetSection from './ThemeWorkbenchPresetSection.vue';

const settingStore = useSettingStore();
const route = useRoute();
const advancedVisible = ref(false);
const expandedAdvancedGroups = ref<Array<string | number>>(['brand']);

const groupIconMap: Record<ThemeWorkbenchGroupKey, string> = {
  overview: 'app',
  appearance: 'fill-color',
  layout: 'view-list',
  typography: 'text',
  style: 'chart-bubble',
  advanced: 'sliders',
};

const groups = computed(() => settingStore.themeWorkbenchGroups);
const presetDefinitions = computed(() => settingStore.themePresetDefinitions);
const effectiveTheme = computed(() => settingStore.effectiveThemeState);
const effectivePresetId = computed(() => settingStore.effectiveThemeState.selectedThemePresetId);
const activeGroup = computed(() => settingStore.activeThemeWorkbenchGroup);
const themeIdentity = computed(() => settingStore.themeIdentitySummary);
const themeDiffItems = computed(() => settingStore.themeAuthorityDiff);
const brandOptions = DEFAULT_COLOR_OPTIONS;

const modeOptions = [
  { type: 'light' as const, text: t('layout.setting.theme.options.light'), icon: SettingLightIcon },
  { type: 'dark' as const, text: t('layout.setting.theme.options.dark'), icon: SettingDarkIcon },
  { type: 'auto' as const, text: t('layout.setting.theme.options.auto'), icon: SettingAutoIcon },
];

const fontFamilyOptions = [
  {
    value: 'system',
    label: t('layout.setting.workbench.typography.system'),
    previewFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif',
  },
  {
    value: 'harmonyos',
    label: 'HarmonyOS Sans',
    previewFamily: '"HarmonyOS Sans SC", "HarmonyOS Sans", "PingFang SC", "Microsoft YaHei", sans-serif',
  },
  {
    value: 'inter',
    label: 'Inter',
    previewFamily: '"Inter", "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif',
  },
  {
    value: 'source-han-sans',
    label: 'Source Han Sans',
    previewFamily: '"Source Han Sans SC", "Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif',
  },
] as const;

const radiusOptions = [
  { value: 'business', label: t('layout.setting.workbench.style.business') },
  { value: 'standard', label: t('layout.setting.workbench.style.standard') },
  { value: 'rounded', label: t('layout.setting.workbench.style.rounded') },
  { value: 'capsule', label: t('layout.setting.workbench.style.capsule') },
] as const;

const shadowOptions = [
  { value: 'flat', label: t('layout.setting.workbench.style.flat') },
  { value: 'standard', label: t('layout.setting.workbench.style.standard') },
  { value: 'floating', label: t('layout.setting.workbench.style.floating') },
] as const;

const densityOptions = [
  { value: 'compact', label: t('layout.setting.workbench.style.compact') },
  { value: 'standard', label: t('layout.setting.workbench.style.standard') },
  { value: 'comfortable', label: t('layout.setting.workbench.style.comfortable') },
] as const;

const tokenGroups: Array<{ value: ThemeTokenGroupKey; label: string }> = [
  { value: 'brand', label: t('layout.setting.workbench.groups.brand') },
  { value: 'text', label: t('layout.setting.workbench.groups.text') },
  { value: 'background', label: t('layout.setting.workbench.groups.background') },
  { value: 'border', label: t('layout.setting.workbench.groups.border') },
  { value: 'component', label: t('layout.setting.workbench.groups.component') },
];

const tokenDefinitionsByGroup = computed(() =>
  tokenGroups.reduce(
    (acc, group) => {
      acc[group.value] = settingStore.themeTokenDefinitions.filter((item) => item.group === group.value);
      return acc;
    },
    {} as Record<ThemeTokenGroupKey, typeof settingStore.themeTokenDefinitions>,
  ),
);

const layoutOptions = computed(
  () =>
    [
      { value: 'side', label: t('layout.setting.workbench.layout.side') },
      { value: 'top', label: t('layout.setting.workbench.layout.top') },
      { value: 'mix', label: t('layout.setting.workbench.layout.mix') },
    ] as const,
);

const modeLabel = computed(() => {
  const matched = modeOptions.find((item) => item.type === effectiveTheme.value.mode);
  return matched?.text ?? effectiveTheme.value.mode;
});

const layoutLabel = computed(() => {
  const matched = layoutOptions.value.find((item) => item.value === settingStore.layout);
  return matched?.label ?? settingStore.layout;
});

const themeIdentityTitle = computed(() =>
  themeIdentity.value.sourceType === 'customized'
    ? t('layout.setting.workbench.identity.customized', { name: themeIdentity.value.currentLabel })
    : themeIdentity.value.currentLabel,
);

const activeFontPreviewFamily = computed(() => {
  const matched = fontFamilyOptions.find((item) => item.value === effectiveTheme.value.fontFamilyPreset);
  return matched?.previewFamily ?? fontFamilyOptions[0].previewFamily;
});

const hasPendingChanges = computed(() => themeDiffItems.value.length > 0);
const showPresetSection = computed(() => activeGroup.value === 'overview' || activeGroup.value === 'appearance');
const splitMenuAvailable = computed(() => settingStore.layout === 'mix');
const splitMenuHint = computed(() =>
  splitMenuAvailable.value
    ? t('layout.setting.workbench.layout.splitMenuHint')
    : t('layout.setting.workbench.layout.onlyMix'),
);
const fixedSidebarAvailable = computed(() => settingStore.layout !== 'top');
const fixedSidebarHint = computed(() =>
  fixedSidebarAvailable.value
    ? t('layout.setting.workbench.layout.fixedSidebarHint')
    : t('layout.setting.workbench.layout.notIntegrated'),
);
const footerOptionVisible = computed(() => route.meta.footer !== false);
const activeTokenEditorMode = ref<'light' | 'dark'>(settingStore.displayMode);
const activeTokenEditorModeLabel = computed(() =>
  activeTokenEditorMode.value === 'dark'
    ? t('layout.setting.workbench.token.dark')
    : t('layout.setting.workbench.token.light'),
);

const drawerVisible = computed({
  get: () => settingStore.showThemeWorkbench,
  set: (visible: boolean) => {
    if (!visible) {
      settingStore.cancelThemeDraft();
      return;
    }

    settingStore.openThemeWorkbench(settingStore.activeThemeWorkbenchGroup);
  },
});

const openGroup = (group: ThemeWorkbenchGroupKey) => {
  settingStore.setActiveThemeWorkbenchGroup(group);
};

const closeWorkbench = () => {
  settingStore.cancelThemeDraft();
};

const toggleAdvancedVisible = (value: boolean) => {
  advancedVisible.value = value;
};

const toggleAdvancedGroup = (group: ThemeTokenGroupKey) => {
  if (expandedAdvancedGroups.value.includes(group)) {
    expandedAdvancedGroups.value = expandedAdvancedGroups.value.filter((item) => item !== group);
    return;
  }

  expandedAdvancedGroups.value = [...expandedAdvancedGroups.value, group];
};

const handleTokenEditorModeChange = (mode: 'light' | 'dark') => {
  activeTokenEditorMode.value = mode;
};
</script>
<style lang="less" scoped>
@import '../../../shared/components/management/card-surface.less';
@import './theme-surface.less';

.theme-workbench-panel {
  :deep(.t-drawer) {
    background: transparent;
    max-width: calc(100vw - 16px);
  }

  :deep(.t-drawer__content) {
    background: var(--td-bg-color-page);
  }

  :deep(.t-drawer__body) {
    padding: 0;
  }
}

.theme-workbench-panel__shell {
  background: var(--td-bg-color-page);
  display: flex;
  flex-direction: column;
  height: 100%;
}

.theme-workbench-panel__header {
  align-items: flex-start;
  background: var(--td-bg-color-container);
  border-bottom: 1px solid var(--td-component-stroke);
  display: flex;
  justify-content: space-between;
  padding: 18px 20px 16px;
}

.panel-title {
  color: var(--td-text-color-primary);
  font-size: 20px;
  font-weight: 700;
}

.panel-subtitle {
  color: var(--td-text-color-secondary);
  font-size: 13px;
  margin-top: 4px;
}

.theme-workbench-panel__body {
  display: grid;
  flex: 1;
  gap: 16px;
  grid-template-columns: 82px minmax(0, 1fr);
  min-height: 0;
  padding: 16px;
}

.theme-workbench-panel__nav {
  display: grid;
  gap: 8px;
  height: fit-content;
}

.theme-workbench-selectable-card() {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: 16px;
  color: var(--td-text-color-primary);
  cursor: pointer;
}

.nav-item {
  appearance: none;
  .theme-workbench-selectable-card();

  color: var(--td-text-color-secondary);
  display: grid;
  gap: 8px;
  min-height: 56px;
  padding: 12px 8px;
  place-items: center;
}

.nav-item--active {
  border-color: var(--td-brand-color);
  box-shadow: var(--td-shadow-1);
  color: var(--td-brand-color);
}

.nav-item__icon {
  font-size: 20px;
}

.nav-item__text {
  font-size: 12px;
  line-height: 1.3;
  max-width: 100%;
  overflow-wrap: anywhere;
  text-align: center;
}

.theme-workbench-panel__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 0;
  overflow: auto;
  padding-bottom: 104px;
  padding-right: 4px;
  scrollbar-color: color-mix(in srgb, var(--td-brand-color) 18%, var(--td-scrollbar-color)) transparent;
}

.section {
  .theme-workbench-surface();

  gap: 14px;
  padding: 16px;
}

.section--compact {
  align-items: center;
  gap: 12px 16px;
  grid-template-columns: minmax(0, 1fr) auto;
}

.section-title {
  color: var(--td-text-color-primary);
  font-size: 16px;
  font-weight: 700;
}

.section-desc {
  color: var(--td-text-color-secondary);
  font-size: 13px;
}

.section--compact .section-desc {
  margin-right: auto;
}

.advanced-settings-card {
  align-items: center;
  background: color-mix(in srgb, var(--td-bg-color-container) 96%, var(--td-bg-color-page));
  border: 1px solid var(--td-component-stroke);
  border-radius: 16px;
  display: grid;
  gap: 16px;
  grid-template-columns: minmax(0, 1fr) auto;
  padding: 14px 16px;
}

.advanced-settings-card__content {
  display: grid;
  gap: 6px;
}

.advanced-settings-card__desc {
  color: color-mix(in srgb, var(--td-text-color-secondary) 90%, var(--td-text-color-primary));
}

.summary-grid,
.choice-grid,
.brand-palette {
  display: grid;
  gap: 12px;
}

.summary-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.summary-card {
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-large);
  padding: 12px;
}

.summary-card__label {
  color: var(--td-text-color-secondary);
  font-size: 12px;
}

.summary-card__value {
  color: var(--td-text-color-primary);
  font-weight: 600;
  margin-top: 6px;
}

.choice-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.choice-card {
  appearance: none;
  .theme-workbench-selectable-card();

  display: grid;
  gap: 10px;
  padding: 12px;
  position: relative;
  text-align: left;
}

.choice-card--active {
  border-color: var(--td-brand-color);
  box-shadow: var(--td-shadow-1);
}

.style-preview-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.style-preview-grid--triple {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.style-preview-card {
  align-items: stretch;
  appearance: none;
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: 14px;
  color: var(--td-text-color-primary);
  display: grid;
  gap: 12px;
  padding: 12px;
  text-align: left;
  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    transform 0.18s ease,
    background-color 0.18s ease;
}

.style-preview-card:hover {
  border-color: color-mix(in srgb, var(--td-brand-color) 24%, var(--td-component-stroke));
  transform: translateY(-1px);
}

.style-preview-card--active {
  background: color-mix(in srgb, var(--td-brand-color) 5%, var(--td-bg-color-page));
  border-color: color-mix(in srgb, var(--td-brand-color) 44%, var(--td-component-stroke));
  box-shadow:
    0 0 0 1px color-mix(in srgb, var(--td-brand-color) 14%, transparent),
    var(--td-shadow-1);
}

.style-preview-card__label {
  font-size: 13px;
  font-weight: 600;
}

.radius-preview,
.shadow-preview,
.density-preview {
  background:
    linear-gradient(180deg, color-mix(in srgb, var(--td-bg-color-container) 88%, transparent), transparent),
    color-mix(in srgb, var(--td-bg-color-container) 65%, var(--td-bg-color-page));
  border: 1px solid color-mix(in srgb, var(--td-component-stroke) 86%, transparent);
  border-radius: 12px;
  min-height: 88px;
  padding: 14px;
}

.radius-preview {
  align-items: end;
  display: flex;
  gap: 12px;
}

.radius-preview__surface {
  background: color-mix(in srgb, var(--td-brand-color) 12%, var(--td-bg-color-container));
  border: 1px solid color-mix(in srgb, var(--td-brand-color) 18%, var(--td-component-stroke));
  display: block;
}

.radius-preview__surface--main {
  flex: 1;
  height: 42px;
}

.radius-preview__surface--sub {
  height: 28px;
  width: 32%;
}

.radius-preview--business .radius-preview__surface--main,
.radius-preview--business .radius-preview__surface--sub {
  border-radius: 6px;
}

.radius-preview--standard .radius-preview__surface--main,
.radius-preview--standard .radius-preview__surface--sub {
  border-radius: 12px;
}

.radius-preview--rounded .radius-preview__surface--main,
.radius-preview--rounded .radius-preview__surface--sub {
  border-radius: 20px;
}

.radius-preview--capsule .radius-preview__surface--main,
.radius-preview--capsule .radius-preview__surface--sub {
  border-radius: 999px;
}

.shadow-preview {
  display: grid;
  place-items: center;
  position: relative;
}

.shadow-preview__card {
  background: color-mix(in srgb, var(--td-bg-color-container) 94%, white 3%);
  border: 1px solid color-mix(in srgb, var(--td-component-stroke) 86%, transparent);
  border-radius: 14px;
  display: block;
  position: absolute;
}

.shadow-preview__card--back {
  height: 34px;
  transform: translate(-10px, -8px);
  width: 56px;
}

.shadow-preview__card--front {
  height: 42px;
  transform: translate(10px, 8px);
  width: 74px;
}

.shadow-preview--flat .shadow-preview__card {
  box-shadow: none;
}

.shadow-preview--standard .shadow-preview__card--back {
  box-shadow: 0 6px 16px rgb(15 23 42 / 10%);
}

.shadow-preview--standard .shadow-preview__card--front {
  box-shadow: 0 12px 24px rgb(15 23 42 / 12%);
}

.shadow-preview--floating .shadow-preview__card--back {
  box-shadow: 0 10px 22px rgb(15 23 42 / 14%);
}

.shadow-preview--floating .shadow-preview__card--front {
  box-shadow: 0 18px 40px rgb(15 23 42 / 18%);
}

[theme-mode='dark'] .shadow-preview--standard .shadow-preview__card--back {
  box-shadow: 0 6px 18px rgb(0 0 0 / 26%);
}

[theme-mode='dark'] .shadow-preview--standard .shadow-preview__card--front {
  box-shadow: 0 14px 28px rgb(0 0 0 / 32%);
}

[theme-mode='dark'] .shadow-preview--floating .shadow-preview__card--back {
  box-shadow: 0 10px 26px rgb(0 0 0 / 34%);
}

[theme-mode='dark'] .shadow-preview--floating .shadow-preview__card--front {
  box-shadow: 0 18px 42px rgb(0 0 0 / 40%);
}

.density-preview {
  display: grid;
}

.density-preview--compact {
  gap: 6px;
}

.density-preview--standard {
  gap: 10px;
}

.density-preview--comfortable {
  gap: 14px;
}

.density-preview__row {
  align-items: center;
  display: grid;
  gap: 10px;
  grid-template-columns: 10px 1fr;
}

.density-preview__bullet {
  background: color-mix(in srgb, var(--td-brand-color) 20%, var(--td-bg-color-container));
  border-radius: 999px;
  display: block;
  height: 10px;
  width: 10px;
}

.density-preview__line {
  background: color-mix(in srgb, var(--td-brand-color) 10%, var(--td-text-color-placeholder));
  border-radius: 999px;
  display: block;
  height: 8px;
  width: 84%;
}

.density-preview__line--long {
  width: 100%;
}

.density-preview__line--short {
  width: 62%;
}

.choice-card__check {
  color: var(--td-brand-color);
  position: absolute;
  right: 10px;
  top: 10px;
}

.choice-card__title {
  font-weight: 600;
}

.mode-thumbnail,
.layout-thumbnail {
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: 12px;
  display: block;
  height: 92px;
  overflow: hidden;
  position: relative;
}

.mode-thumbnail {
  align-items: center;
  display: flex;
  justify-content: center;
  padding: 12px;
}

.mode-thumbnail--light {
  background: linear-gradient(180deg, #fff, #f3f5f9);
}

.mode-thumbnail--dark {
  background: linear-gradient(180deg, #1f2937, #111827);
}

.mode-thumbnail--auto {
  background: linear-gradient(135deg, #fff 0 50%, #111827 50% 100%);
}

.mode-thumbnail__icon {
  max-width: 92px;
  width: 100%;
}

.layout-thumbnail {
  display: grid;
  grid-template-columns: 24px 1fr;
  grid-template-rows: 18px 1fr;
  padding: 10px;
}

.layout-thumbnail__header,
.layout-thumbnail__sidebar,
.layout-thumbnail__content {
  border-radius: 8px;
  display: block;
}

.layout-thumbnail__header {
  background: color-mix(in srgb, var(--td-brand-color) 18%, var(--td-bg-color-container));
  grid-column: 1 / span 2;
}

.layout-thumbnail__sidebar {
  background: color-mix(in srgb, var(--td-brand-color) 10%, var(--td-bg-color-container));
  grid-row: 2;
  margin-top: 8px;
}

.layout-thumbnail__content {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  grid-row: 2;
  margin-left: 8px;
  margin-top: 8px;
}

.layout-thumbnail--top {
  grid-template-columns: 1fr;
  grid-template-rows: 18px 1fr;
}

.layout-thumbnail--top .layout-thumbnail__header {
  grid-column: 1;
}

.layout-thumbnail--top .layout-thumbnail__sidebar {
  display: none;
}

.layout-thumbnail--top .layout-thumbnail__content {
  grid-column: 1;
  margin-left: 0;
}

.layout-thumbnail--mix::after {
  background: color-mix(in srgb, var(--td-brand-color) 16%, var(--td-bg-color-container));
  border-radius: 6px;
  content: '';
  height: 38px;
  left: 42px;
  position: absolute;
  top: 34px;
  width: 16px;
}

.switch-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.switch-item {
  align-items: center;
  display: flex;
  gap: 16px;
  justify-content: space-between;
}

.switch-item__label {
  color: var(--td-text-color-primary);
  font-weight: 500;
}

.switch-item__hint {
  color: var(--td-text-color-secondary);
  font-size: 12px;
  margin-top: 4px;
}

.brand-palette {
  grid-template-columns: repeat(8, minmax(0, 1fr));
}

.brand-palette__item {
  align-items: center;
  appearance: none;
  border: 2px solid transparent;
  border-radius: 14px;
  color: #fff;
  cursor: pointer;
  display: inline-flex;
  height: 40px;
  justify-content: center;
}

.brand-palette__item--active {
  border-color: var(--td-text-color-primary);
}

.brand-input {
  align-items: center;
  display: flex;
  gap: 12px;
}

.brand-input input[type='color'] {
  appearance: none;
  background: transparent;
  border: 0;
  cursor: pointer;
  height: 40px;
  padding: 0;
  width: 48px;
}

.brand-input :deep(.t-input) {
  flex: 1;
}

.font-option-list {
  display: grid;
  gap: 10px;
}

.font-option {
  align-items: center;
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: 14px;
  cursor: pointer;
  display: grid;
  gap: 12px;
  grid-template-columns: minmax(0, 1fr) auto;
  padding: 12px 14px;
  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease,
    background-color 0.18s ease;
}

.font-option:hover {
  background: color-mix(in srgb, var(--td-bg-color-container) 82%, var(--td-bg-color-page));
  border-color: color-mix(in srgb, var(--td-brand-color) 20%, var(--td-component-stroke));
}

.font-option--active {
  border-color: color-mix(in srgb, var(--td-brand-color) 40%, var(--td-component-stroke));
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--td-brand-color) 12%, transparent);
}

.font-option__input {
  opacity: 0;
  pointer-events: none;
  position: absolute;
}

.font-option__main {
  display: grid;
  gap: 4px;
  min-width: 0;
}

.font-option__title {
  color: var(--td-text-color-primary);
  font-size: 14px;
  font-weight: 600;
}

.font-option__preview {
  color: var(--td-text-color-secondary);
  display: block;
  font-size: 13px;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.font-option__check {
  color: var(--td-brand-color);
}

.font-preview-strip {
  align-items: center;
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: 14px;
  display: grid;
  gap: 6px;
  padding: 12px 14px;
}

.font-preview-strip__label {
  color: var(--td-text-color-secondary);
  font-size: 12px;
}

.font-preview-strip__content {
  color: var(--td-text-color-primary);
  font-size: 16px;
  line-height: 1.4;
}

.advanced-groups {
  display: grid;
  gap: 12px;
}

.advanced-group {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: 16px;
  overflow: hidden;
}

.advanced-group--expanded {
  border-color: color-mix(in srgb, var(--td-brand-color) 20%, var(--td-component-stroke));
  box-shadow: var(--td-shadow-1);
}

.advanced-group__header {
  align-items: center;
  appearance: none;
  background: transparent;
  color: inherit;
  display: grid;
  gap: 12px;
  grid-template-columns: minmax(0, 1fr) auto;
  padding: 14px 16px;
  text-align: left;
  width: 100%;
}

.advanced-group__title-block {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  min-width: 0;
}

.advanced-group__title {
  color: var(--td-text-color-primary);
  font-size: 14px;
  font-weight: 700;
}

.advanced-group__mode-tag,
.advanced-group__count {
  background: color-mix(in srgb, var(--td-brand-color) 10%, var(--td-bg-color-page));
  border: 1px solid color-mix(in srgb, var(--td-brand-color) 16%, var(--td-component-stroke));
  border-radius: 999px;
  color: var(--td-text-color-secondary);
  font-size: 12px;
  line-height: 1;
  padding: 5px 8px;
}

.advanced-group__header-actions {
  align-items: center;
  display: inline-flex;
  gap: 10px;
}

.advanced-group__chevron {
  color: var(--td-text-color-secondary);
  transition: transform 0.18s ease;
}

.advanced-group--expanded .advanced-group__chevron {
  transform: rotate(180deg);
}

.advanced-group__body {
  border-top: 1px solid color-mix(in srgb, var(--td-component-stroke) 80%, transparent);
  padding: 0 16px 16px;
}

.theme-workbench-panel__footer {
  align-items: center;
  background: var(--td-bg-color-container);
  border-top: 1px solid var(--td-component-stroke);
  bottom: 0;
  display: flex;
  gap: 12px;
  justify-content: space-between;
  left: 0;
  padding: 12px 16px calc(12px + env(safe-area-inset-bottom, 0px));
  position: sticky;
  z-index: 2;
}

.theme-workbench-panel__footer-actions {
  display: flex;
  gap: 12px;
}

@media (width <= 768px) {
  .theme-workbench-panel__body {
    grid-template-columns: 1fr;
  }

  .theme-workbench-panel__nav {
    grid-template-columns: repeat(5, minmax(0, 1fr));
  }

  .summary-grid,
  .choice-grid,
  .style-preview-grid,
  .style-preview-grid--triple {
    grid-template-columns: 1fr;
  }

  .section--compact {
    grid-template-columns: 1fr;
  }

  .advanced-settings-card {
    grid-template-columns: 1fr;
  }
}
</style>
