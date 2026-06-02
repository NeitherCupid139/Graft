import type { ThemePresetDefinition, ThemeTokenDefinition, ThemeWorkbenchGroupDefinition } from '@/types/theme';

export const DEFAULT_THEME_PRESET_ID = 'tdesign-default';

export const THEME_WORKBENCH_GROUPS: ThemeWorkbenchGroupDefinition[] = [
  { key: 'overview', labelKey: 'layout.setting.workbench.groups.overview' },
  { key: 'appearance', labelKey: 'layout.setting.workbench.groups.appearance' },
  { key: 'layout', labelKey: 'layout.setting.workbench.groups.layout' },
  { key: 'typography', labelKey: 'layout.setting.workbench.groups.typography' },
  { key: 'style', labelKey: 'layout.setting.workbench.groups.style' },
  { key: 'advanced', labelKey: 'layout.setting.workbench.groups.advanced' },
];

export const THEME_TOKEN_DEFINITIONS: ThemeTokenDefinition[] = [
  { key: '--td-brand-color', group: 'brand', label: 'Primary Brand', followsBrandColor: true },
  { key: '--td-brand-color-1', group: 'brand', label: 'Brand Scale 1', followsBrandColor: true },
  { key: '--td-brand-color-2', group: 'brand', label: 'Brand Scale 2', followsBrandColor: true },
  { key: '--td-brand-color-3', group: 'brand', label: 'Brand Scale 3', followsBrandColor: true },
  { key: '--td-brand-color-4', group: 'brand', label: 'Brand Scale 4', followsBrandColor: true },
  { key: '--td-brand-color-5', group: 'brand', label: 'Brand Scale 5', followsBrandColor: true },
  { key: '--td-brand-color-6', group: 'brand', label: 'Brand Scale 6', followsBrandColor: true },
  { key: '--td-brand-color-7', group: 'brand', label: 'Brand Scale 7', followsBrandColor: true },
  { key: '--td-brand-color-8', group: 'brand', label: 'Brand Scale 8', followsBrandColor: true },
  { key: '--td-brand-color-9', group: 'brand', label: 'Brand Scale 9', followsBrandColor: true },
  { key: '--td-brand-color-10', group: 'brand', label: 'Brand Scale 10', followsBrandColor: true },
  { key: '--td-text-color-primary', group: 'text', label: 'Primary Text' },
  { key: '--td-text-color-secondary', group: 'text', label: 'Secondary Text' },
  { key: '--td-text-color-placeholder', group: 'text', label: 'Placeholder Text' },
  { key: '--td-bg-color-page', group: 'background', label: 'Page Background' },
  { key: '--td-bg-color-container', group: 'background', label: 'Container Background' },
  { key: '--td-bg-color-container-hover', group: 'background', label: 'Container Hover Background' },
  { key: '--td-component-border', group: 'border', label: 'Component Border' },
  { key: '--td-component-stroke', group: 'border', label: 'Component Stroke' },
  { key: '--td-border-level-1-color', group: 'border', label: 'Level 1 Border' },
  { key: '--td-scrollbar-color', group: 'border', label: 'Scrollbar' },
  { key: '--td-success-color', group: 'component', label: 'Success Feedback' },
  { key: '--td-error-color', group: 'component', label: 'Error Feedback' },
  { key: '--td-warning-color', group: 'component', label: 'Warning Feedback' },
  { key: '--td-font-family', group: 'component', label: 'Font Family' },
  { key: '--td-font-body-medium', group: 'component', label: 'Body Medium' },
  { key: '--td-font-title-medium', group: 'component', label: 'Title Medium' },
  { key: '--td-font-title-large', group: 'component', label: 'Title Large' },
  { key: '--td-radius-small', group: 'component', label: 'Small Radius' },
  { key: '--td-radius-default', group: 'component', label: 'Default Radius' },
  { key: '--td-radius-medium', group: 'component', label: 'Medium Radius' },
  { key: '--td-radius-large', group: 'component', label: 'Large Radius' },
  { key: '--td-radius-extraLarge', group: 'component', label: 'Extra Large Radius' },
  { key: '--td-radius-circle', group: 'component', label: 'Circle Radius' },
  { key: '--td-shadow-1', group: 'component', label: 'Shadow Level 1' },
  { key: '--td-shadow-2', group: 'component', label: 'Shadow Level 2' },
  { key: '--td-shadow-3', group: 'component', label: 'Shadow Level 3' },
  { key: '--td-comp-size-xs', group: 'component', label: 'Component Size XS' },
  { key: '--td-comp-size-s', group: 'component', label: 'Component Size S' },
  { key: '--td-comp-size-m', group: 'component', label: 'Component Size M' },
  { key: '--td-comp-size-l', group: 'component', label: 'Component Size L' },
  { key: '--td-comp-size-xl', group: 'component', label: 'Component Size XL' },
];

// 推荐主题保持为纯前端本地配置，先服务当前工作台预览，不与后端契约耦合。
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
    description: '收敛到更贴近腾讯云控制台的亮色主色，不改写页面中性色。',
    brandTheme: '#0064FF',
  },
  {
    id: 'mountain-green',
    label: 'Mountain Green',
    description: '用更温和的绿色主色驱动按钮、选中态与状态强调，不改写页面中性色。',
    brandTheme: '#2BA471',
  },
  {
    id: 'midnight-blue',
    label: 'Midnight Blue',
    description: '默认以暗色模式启动，保持 TDesign 中性深色工作区。',
    brandTheme: '#3B82F6',
    mode: 'dark',
  },
];
