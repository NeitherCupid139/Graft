import type { ThemePresetDefinition, ThemeTokenDefinition, ThemeWorkbenchGroupDefinition } from '@/types/theme';

export const DEFAULT_THEME_PRESET_ID = 'tdesign-default';

export const THEME_WORKBENCH_GROUPS: ThemeWorkbenchGroupDefinition[] = [
  { key: 'overview', labelKey: 'layout.setting.workbench.groups.overview' },
  { key: 'brand', labelKey: 'layout.setting.workbench.groups.brand' },
  { key: 'semantic', labelKey: 'layout.setting.workbench.groups.semantic' },
  { key: 'neutral', labelKey: 'layout.setting.workbench.groups.neutral' },
  { key: 'font', labelKey: 'layout.setting.workbench.groups.font' },
  { key: 'radius', labelKey: 'layout.setting.workbench.groups.radius' },
  { key: 'shadow', labelKey: 'layout.setting.workbench.groups.shadow' },
  { key: 'size', labelKey: 'layout.setting.workbench.groups.size' },
];

export const THEME_TOKEN_DEFINITIONS: ThemeTokenDefinition[] = [
  { key: '--td-brand-color', group: 'brand', label: 'brand-color', followsBrandColor: true },
  { key: '--td-brand-color-1', group: 'brand', label: 'brand-color-1', followsBrandColor: true },
  { key: '--td-brand-color-2', group: 'brand', label: 'brand-color-2', followsBrandColor: true },
  { key: '--td-brand-color-3', group: 'brand', label: 'brand-color-3', followsBrandColor: true },
  { key: '--td-brand-color-4', group: 'brand', label: 'brand-color-4', followsBrandColor: true },
  { key: '--td-brand-color-5', group: 'brand', label: 'brand-color-5', followsBrandColor: true },
  { key: '--td-brand-color-6', group: 'brand', label: 'brand-color-6', followsBrandColor: true },
  { key: '--td-brand-color-7', group: 'brand', label: 'brand-color-7', followsBrandColor: true },
  { key: '--td-brand-color-8', group: 'brand', label: 'brand-color-8', followsBrandColor: true },
  { key: '--td-brand-color-9', group: 'brand', label: 'brand-color-9', followsBrandColor: true },
  { key: '--td-brand-color-10', group: 'brand', label: 'brand-color-10', followsBrandColor: true },
  { key: '--td-success-color', group: 'semantic', label: 'success-color' },
  { key: '--td-error-color', group: 'semantic', label: 'error-color' },
  { key: '--td-warning-color', group: 'semantic', label: 'warning-color' },
  { key: '--td-text-color-primary', group: 'neutral', label: 'text-color-primary' },
  { key: '--td-text-color-secondary', group: 'neutral', label: 'text-color-secondary' },
  { key: '--td-text-color-placeholder', group: 'neutral', label: 'text-color-placeholder' },
  { key: '--td-bg-color-page', group: 'neutral', label: 'bg-color-page' },
  { key: '--td-bg-color-container', group: 'neutral', label: 'bg-color-container' },
  { key: '--td-bg-color-container-hover', group: 'neutral', label: 'bg-color-container-hover' },
  { key: '--td-component-border', group: 'neutral', label: 'component-border' },
  { key: '--td-component-stroke', group: 'neutral', label: 'component-stroke' },
  { key: '--td-border-level-1-color', group: 'neutral', label: 'border-level-1-color' },
  { key: '--td-scrollbar-color', group: 'neutral', label: 'scrollbar-color' },
  { key: '--td-font-family', group: 'font', label: 'font-family' },
  { key: '--td-font-body-medium', group: 'font', label: 'font-body-medium' },
  { key: '--td-font-title-medium', group: 'font', label: 'font-title-medium' },
  { key: '--td-font-title-large', group: 'font', label: 'font-title-large' },
  { key: '--td-radius-small', group: 'radius', label: 'radius-small' },
  { key: '--td-radius-default', group: 'radius', label: 'radius-default' },
  { key: '--td-radius-medium', group: 'radius', label: 'radius-medium' },
  { key: '--td-radius-large', group: 'radius', label: 'radius-large' },
  { key: '--td-radius-extraLarge', group: 'radius', label: 'radius-extra-large' },
  { key: '--td-radius-circle', group: 'radius', label: 'radius-circle' },
  { key: '--td-shadow-1', group: 'shadow', label: 'shadow-1' },
  { key: '--td-shadow-2', group: 'shadow', label: 'shadow-2' },
  { key: '--td-shadow-3', group: 'shadow', label: 'shadow-3' },
  { key: '--td-comp-size-xs', group: 'size', label: 'comp-size-xs' },
  { key: '--td-comp-size-s', group: 'size', label: 'comp-size-s' },
  { key: '--td-comp-size-m', group: 'size', label: 'comp-size-m' },
  { key: '--td-comp-size-l', group: 'size', label: 'comp-size-l' },
  { key: '--td-comp-size-xl', group: 'size', label: 'comp-size-xl' },
];

export const THEME_PRESET_DEFINITIONS: ThemePresetDefinition[] = [
  {
    id: 'tdesign-default',
    label: 'TDesign Original',
    description: '保留 starter 默认主题基线。',
    brandTheme: '#0052D9',
  },
  {
    id: 'tencent-cloud',
    label: 'Tencent Cloud',
    description: '收敛到更贴近腾讯云控制台的亮色主色。',
    brandTheme: '#0064FF',
    tokenOverrides: {
      light: {
        '--td-bg-color-page': '#f5f7ff',
        '--td-bg-color-container-hover': '#eef3ff',
      },
      dark: {
        '--td-bg-color-page': '#111827',
        '--td-bg-color-container': '#1b2333',
      },
    },
  },
];
