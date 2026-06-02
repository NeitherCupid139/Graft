import type STYLE_CONFIG from '@/config/style';
import type { ModeType } from '@/utils/types';

export type SettingStyleConfig = typeof STYLE_CONFIG;

export type ThemeWorkbenchGroupKey = 'overview' | 'appearance' | 'typography' | 'style' | 'advanced';

export type ThemeTokenGroupKey = 'brand' | 'semantic' | 'neutral' | 'font' | 'radius' | 'shadow' | 'size';

export type ThemeSourceType = 'preset' | 'customized';

export type ThemeTokenMap = Record<string, string>;

export interface ThemeModeTokenState {
  light: ThemeTokenMap;
  dark: ThemeTokenMap;
}

export interface ThemeWorkbenchGroupDefinition {
  key: ThemeWorkbenchGroupKey;
  labelKey: string;
  descriptionKey?: string;
}

export interface ThemeTokenDefinition {
  key: string;
  group: ThemeTokenGroupKey;
  label: string;
  followsBrandColor?: boolean;
}

export interface ThemePresetDefinition {
  id: string;
  label: string;
  description: string;
  brandTheme: string;
  mode?: ModeType | 'auto';
  tokenOverrides?: Partial<ThemeModeTokenState>;
}

export interface ThemeAuthorityState {
  mode: ModeType | 'auto';
  brandTheme: string;
  selectedThemePresetId: string | null;
  themeSource: ThemeSourceType;
  fontFamilyPreset: 'system' | 'harmonyos' | 'inter' | 'source-han-sans';
  radiusPreset: 'business' | 'standard' | 'rounded' | 'capsule';
  shadowPreset: 'flat' | 'standard' | 'floating';
  densityPreset: 'compact' | 'standard' | 'comfortable';
  themeTokenOverrides: ThemeModeTokenState;
}
