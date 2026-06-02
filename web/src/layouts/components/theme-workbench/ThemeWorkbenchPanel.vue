<template>
  <t-drawer
    v-model:visible="drawerVisible"
    class="theme-workbench-panel"
    destroy-on-close
    placement="right"
    :size="drawerSize"
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
            :class="{ 'nav-item--active': group.key === settingStore.activeThemeWorkbenchGroup }"
            @click="openGroup(group.key)"
          >
            <span class="nav-item__icon">
              <t-icon :name="groupIconMap[group.key]" />
            </span>
            <span class="nav-item__text">{{ t(group.labelKey) }}</span>
          </button>
        </aside>

        <div class="theme-workbench-panel__main">
          <section class="theme-workbench-panel__content">
            <template v-if="activeGroup === 'overview'">
              <div class="section">
                <div class="section-title">{{ t('layout.setting.workbench.overview.currentTheme') }}</div>
                <div class="theme-summary">
                  <div class="theme-summary__name">{{ settingStore.effectiveThemeDisplayName }}</div>
                  <div class="theme-summary__meta">
                    <t-tag theme="primary" variant="light-outline">{{ modeLabel }}</t-tag>
                    <t-tag variant="light-outline">{{ radiusLabel }}</t-tag>
                    <t-tag variant="light-outline">{{ densityLabel }}</t-tag>
                  </div>
                </div>
              </div>

              <component
                :is="ThemeWorkbenchPresetSection"
                :title="t('layout.setting.workbench.presets.title')"
                :presets="presetDefinitions"
                :active-preset-id="effectivePresetId"
                @select="settingStore.selectThemePreset"
              />

              <div class="section">
                <div class="section-title">{{ t('layout.setting.workbench.overview.quickActions') }}</div>
                <div class="quick-actions">
                  <t-button variant="outline" @click="toggleDarkMode">{{
                    t('layout.setting.workbench.overview.toggleDark')
                  }}</t-button>
                  <t-button variant="outline" @click="settingStore.setCustomBrandTheme(brandOptions[0])">
                    {{ t('layout.setting.workbench.overview.resetColor') }}
                  </t-button>
                  <t-button theme="primary" @click="settingStore.resetThemeDraftToDefault()">
                    {{ t('layout.setting.workbench.actions.reset') }}
                  </t-button>
                </div>
              </div>

              <div class="section">
                <div class="section-title">{{ t('layout.setting.workbench.layoutPreference.title') }}</div>
                <div class="section-desc">{{ t('layout.setting.workbench.layoutPreference.description') }}</div>
                <div class="switch-list">
                  <div class="switch-item">
                    <span>{{ t('layout.setting.element.showHeader') }}</span>
                    <t-switch
                      :model-value="settingStore.showHeader"
                      @update:model-value="(value) => settingStore.updateConfig({ showHeader: value })"
                    />
                  </div>
                  <div class="switch-item">
                    <span>{{ t('layout.setting.element.showBreadcrumb') }}</span>
                    <t-switch
                      :model-value="settingStore.showBreadcrumb"
                      @update:model-value="(value) => settingStore.updateConfig({ showBreadcrumb: value })"
                    />
                  </div>
                  <div class="switch-item">
                    <span>{{ t('layout.setting.element.showFooter') }}</span>
                    <t-switch
                      :model-value="settingStore.showFooter"
                      @update:model-value="(value) => settingStore.updateConfig({ showFooter: value })"
                    />
                  </div>
                </div>
              </div>
            </template>

            <template v-else-if="activeGroup === 'appearance'">
              <div class="section">
                <div class="section-title">{{ t('layout.setting.theme.mode') }}</div>
                <div class="mode-grid">
                  <button
                    v-for="item in modeOptions"
                    :key="item.type"
                    type="button"
                    class="mode-card"
                    :class="{ 'mode-card--active': effectiveTheme.mode === item.type }"
                    @click="settingStore.updateThemeDraftAppearance({ mode: item.type })"
                  >
                    <span class="mode-card__preview">
                      <component :is="item.icon" class="mode-card__icon" />
                    </span>
                    <span class="mode-card__label">{{ item.text }}</span>
                  </button>
                </div>
              </div>

              <div class="section">
                <div class="section-title">{{ t('layout.setting.theme.color') }}</div>
                <div class="brand-palette">
                  <button
                    v-for="color in brandOptions"
                    :key="color"
                    type="button"
                    class="brand-palette__item"
                    :class="{ 'brand-palette__item--active': effectiveTheme.brandTheme === color }"
                    :style="{ background: color }"
                    @click="settingStore.setCustomBrandTheme(color)"
                  />
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

              <theme-workbench-preset-section
                :title="t('layout.setting.workbench.presets.title')"
                :presets="presetDefinitions"
                :active-preset-id="effectivePresetId"
                @select="settingStore.selectThemePreset"
              />
            </template>

            <template v-else-if="activeGroup === 'typography'">
              <div class="section">
                <div class="section-title">{{ t('layout.setting.workbench.typography.fontFamily') }}</div>
                <t-radio-group
                  :model-value="effectiveTheme.fontFamilyPreset"
                  variant="default-filled"
                  @update:model-value="
                    (value) =>
                      settingStore.updateThemeDraftAppearance({
                        fontFamilyPreset: value as typeof effectiveTheme.fontFamilyPreset,
                      })
                  "
                >
                  <t-radio-button v-for="item in fontFamilyOptions" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </t-radio-button>
                </t-radio-group>
              </div>
            </template>

            <template v-else-if="activeGroup === 'style'">
              <div class="section">
                <div class="section-title">{{ t('layout.setting.workbench.style.radius') }}</div>
                <t-radio-group
                  :model-value="effectiveTheme.radiusPreset"
                  variant="default-filled"
                  @update:model-value="
                    (value) =>
                      settingStore.updateThemeDraftAppearance({
                        radiusPreset: value as typeof effectiveTheme.radiusPreset,
                      })
                  "
                >
                  <t-radio-button v-for="item in radiusOptions" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </t-radio-button>
                </t-radio-group>
              </div>

              <div class="section">
                <div class="section-title">{{ t('layout.setting.workbench.style.shadow') }}</div>
                <t-radio-group
                  :model-value="effectiveTheme.shadowPreset"
                  variant="default-filled"
                  @update:model-value="
                    (value) =>
                      settingStore.updateThemeDraftAppearance({
                        shadowPreset: value as typeof effectiveTheme.shadowPreset,
                      })
                  "
                >
                  <t-radio-button v-for="item in shadowOptions" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </t-radio-button>
                </t-radio-group>
              </div>

              <div class="section">
                <div class="section-title">{{ t('layout.setting.workbench.style.density') }}</div>
                <t-radio-group
                  :model-value="effectiveTheme.densityPreset"
                  variant="default-filled"
                  @update:model-value="
                    (value) =>
                      settingStore.updateThemeDraftAppearance({
                        densityPreset: value as typeof effectiveTheme.densityPreset,
                      })
                  "
                >
                  <t-radio-button v-for="item in densityOptions" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </t-radio-button>
                </t-radio-group>
              </div>
            </template>

            <template v-else>
              <div class="section section--compact">
                <div class="section-title">{{ t('layout.setting.workbench.advanced.title') }}</div>
                <div class="section-desc">{{ t('layout.setting.workbench.advanced.description') }}</div>
                <t-button theme="warning" variant="outline" @click="advancedVisible = !advancedVisible">
                  {{
                    advancedVisible
                      ? t('layout.setting.workbench.advanced.hide')
                      : t('layout.setting.workbench.advanced.enter')
                  }}
                </t-button>
              </div>

              <template v-if="advancedVisible">
                <div class="section">
                  <div class="section-title">{{ activeTokenGroupLabel }}</div>
                  <div class="token-tab-list">
                    <t-button
                      v-for="group in tokenGroups"
                      :key="group.value"
                      variant="outline"
                      size="small"
                      :theme="settingStore.activeThemeTokenGroup === group.value ? 'primary' : 'default'"
                      @click="settingStore.activeThemeTokenGroup = group.value"
                    >
                      {{ group.label }}
                    </t-button>
                  </div>
                </div>
                <theme-token-editor
                  :group-key="settingStore.activeThemeTokenGroup"
                  :token-definitions="activeTokenDefinitions"
                />
              </template>
            </template>
          </section>

          <theme-workbench-preview class="theme-workbench-panel__preview" />
        </div>
      </div>
      <footer class="theme-workbench-panel__footer">
        <t-button variant="outline" @click="settingStore.resetThemeDraftToDefault()">
          {{ t('layout.setting.workbench.actions.reset') }}
        </t-button>
        <div class="theme-workbench-panel__footer-actions">
          <t-button variant="outline" @click="settingStore.cancelThemeDraft()">
            {{ t('layout.setting.workbench.actions.cancel') }}
          </t-button>
          <t-button theme="primary" @click="settingStore.applyThemeDraft()">
            {{ t('layout.setting.workbench.actions.apply') }}
          </t-button>
        </div>
      </footer>
    </div>
  </t-drawer>
</template>
<script setup lang="ts">
import { computed, ref } from 'vue';

import SettingAutoIcon from '@/assets/assets-setting-auto.svg';
import SettingDarkIcon from '@/assets/assets-setting-dark.svg';
import SettingLightIcon from '@/assets/assets-setting-light.svg';
import { DEFAULT_COLOR_OPTIONS } from '@/config/color';
import { t } from '@/locales';
import { useSettingStore } from '@/store';
import type { ThemeTokenGroupKey, ThemeWorkbenchGroupKey } from '@/types/theme';

import ThemeTokenEditor from './ThemeTokenEditor.vue';
import ThemeWorkbenchPresetSection from './ThemeWorkbenchPresetSection.vue';
import ThemeWorkbenchPreview from './ThemeWorkbenchPreview.vue';

const settingStore = useSettingStore();
const advancedVisible = ref(false);

const groupIconMap: Record<ThemeWorkbenchGroupKey, string> = {
  overview: 'app',
  appearance: 'fill-color',
  typography: 'text',
  style: 'chart-bubble',
  advanced: 'settings',
};

const groups = computed(() => settingStore.themeWorkbenchGroups);
const presetDefinitions = computed(() => settingStore.themePresetDefinitions);
const effectiveTheme = computed(() => settingStore.effectiveThemeState);
const effectivePresetId = computed(() => settingStore.effectiveThemeState.selectedThemePresetId);
const activeGroup = computed(() => settingStore.activeThemeWorkbenchGroup);
const brandOptions = DEFAULT_COLOR_OPTIONS;
const drawerSize = 'min(1280px, calc(100vw - 24px))';

const modeOptions = [
  { type: 'light' as const, text: t('layout.setting.theme.options.light'), icon: SettingLightIcon },
  { type: 'dark' as const, text: t('layout.setting.theme.options.dark'), icon: SettingDarkIcon },
  { type: 'auto' as const, text: t('layout.setting.theme.options.auto'), icon: SettingAutoIcon },
];

const fontFamilyOptions = [
  { value: 'system', label: t('layout.setting.workbench.typography.system') },
  { value: 'harmonyos', label: 'HarmonyOS Sans' },
  { value: 'inter', label: 'Inter' },
  { value: 'source-han-sans', label: 'Source Han Sans' },
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
  { value: 'semantic', label: t('layout.setting.workbench.groups.semantic') },
  { value: 'neutral', label: t('layout.setting.workbench.groups.neutral') },
  { value: 'font', label: t('layout.setting.workbench.groups.font') },
  { value: 'radius', label: t('layout.setting.workbench.groups.radius') },
  { value: 'shadow', label: t('layout.setting.workbench.groups.shadow') },
  { value: 'size', label: t('layout.setting.workbench.groups.size') },
];

const activeTokenDefinitions = computed(() =>
  settingStore.themeTokenDefinitions.filter((item) => item.group === settingStore.activeThemeTokenGroup),
);

const activeTokenGroupLabel = computed(() => {
  const matched = tokenGroups.find((item) => item.value === settingStore.activeThemeTokenGroup);
  return matched?.label ?? t('layout.setting.workbench.advanced.title');
});

const modeLabel = computed(() => {
  const matched = modeOptions.find((item) => item.type === effectiveTheme.value.mode);
  return matched?.text ?? effectiveTheme.value.mode;
});

const radiusLabel = computed(() => {
  const matched = radiusOptions.find((item) => item.value === effectiveTheme.value.radiusPreset);
  return matched?.label ?? effectiveTheme.value.radiusPreset;
});

const densityLabel = computed(() => {
  const matched = densityOptions.find((item) => item.value === effectiveTheme.value.densityPreset);
  return matched?.label ?? effectiveTheme.value.densityPreset;
});

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

const toggleDarkMode = () => {
  settingStore.updateThemeDraftAppearance({ mode: effectiveTheme.value.mode === 'dark' ? 'light' : 'dark' });
};
</script>
<style lang="less" scoped>
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
  background: linear-gradient(135deg, color-mix(in srgb, var(--td-brand-color) 84%, #fff), var(--td-brand-color));
  color: #fff;
  display: flex;
  justify-content: space-between;
  padding: 18px 20px 16px;
}

.theme-workbench-panel__header :deep(.t-button) {
  color: #fff;
}

.panel-title {
  font-size: 22px;
  font-weight: 700;
}

.panel-subtitle {
  font-size: 13px;
  margin-top: 4px;
  opacity: 0.78;
}

.theme-workbench-panel__body {
  display: grid;
  flex: 1;
  gap: 16px;
  grid-template-columns: 72px minmax(0, 1fr);
  min-height: 0;
  padding: 16px;
}

.theme-workbench-panel__nav {
  display: grid;
  gap: 8px;
  height: fit-content;
}

.nav-item {
  appearance: none;
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: 16px;
  color: var(--td-text-color-secondary);
  cursor: pointer;
  display: grid;
  gap: 5px;
  padding: 10px 4px;
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
  line-height: 1.2;
  text-align: center;
}

.theme-workbench-panel__main {
  display: grid;
  gap: 16px;
  grid-template-columns: minmax(0, 1fr) minmax(360px, 420px);
  min-height: 0;
}

.theme-workbench-panel__content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 0;
  overflow: auto;
  padding-right: 4px;
}

.theme-workbench-panel__preview {
  min-width: 0;
}

.section {
  .theme-workbench-surface();

  gap: 14px;
  padding: 16px;
}

.section--compact {
  gap: 8px;
}

.section-title {
  color: var(--td-text-color-primary);
  font-size: 16px;
  font-weight: 700;
}

.section-desc {
  color: var(--td-text-color-secondary);
}

.theme-summary__name {
  font-size: 20px;
  font-weight: 700;
}

.theme-summary__meta,
.quick-actions,
.switch-list,
.switch-item,
.token-tab-list,
.brand-input {
  display: flex;
  gap: 12px;
}

.theme-summary__meta,
.quick-actions,
.token-tab-list {
  flex-wrap: wrap;
}

.switch-list {
  flex-direction: column;
}

.switch-item {
  align-items: center;
  justify-content: space-between;
}

.preset-grid,
.mode-grid {
  display: grid;
  gap: 12px;
}

.preset-grid {
  grid-template-columns: repeat(auto-fit, minmax(132px, 1fr));
}

.mode-grid {
  grid-template-columns: repeat(auto-fit, minmax(108px, 1fr));
}

.preset-card,
.mode-card {
  appearance: none;
  background: var(--td-bg-color-page);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-large);
  cursor: pointer;
  display: grid;
  gap: 8px;
  padding: 12px;
  text-align: left;
}

.preset-card--active,
.mode-card--active {
  border-color: var(--td-brand-color);
  box-shadow: var(--td-shadow-1);
}

.mode-card {
  justify-items: center;
  text-align: center;
}

.mode-card__preview {
  display: grid;
  place-items: center;
}

.mode-card__icon {
  max-width: 88px;
  width: 100%;
}

.preset-card__swatch {
  border-radius: 12px;
  display: block;
  height: 48px;
}

.preset-card__title {
  font-weight: 700;
}

.preset-card__desc {
  color: var(--td-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}

.brand-palette {
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(8, minmax(0, 1fr));
}

.brand-palette__item {
  appearance: none;
  border: 2px solid transparent;
  border-radius: 14px;
  cursor: pointer;
  height: 40px;
}

.brand-palette__item--active {
  border-color: var(--td-text-color-primary);
}

.brand-input {
  align-items: center;
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

.theme-workbench-panel__footer {
  align-items: center;
  border-top: 1px solid var(--td-component-stroke);
  display: flex;
  justify-content: space-between;
  padding: 12px 16px;
}

.theme-workbench-panel__footer-actions {
  display: flex;
  gap: 12px;
}

@media (width <= 1024px) {
  .theme-workbench-panel__body {
    grid-template-columns: 1fr;
  }

  .theme-workbench-panel__nav {
    grid-template-columns: repeat(5, minmax(0, 1fr));
  }

  .theme-workbench-panel__main {
    grid-template-columns: 1fr;
  }
}
</style>
