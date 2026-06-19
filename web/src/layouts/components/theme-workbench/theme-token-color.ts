// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

export interface ParsedThemeTokenColor {
  alpha: number;
  blue: number;
  css: string;
  green: number;
  hex: string;
  red: number;
}

const FALLBACK_HEX = '#0052d9';

let sharedCanvasContext: CanvasRenderingContext2D | null | undefined;

function getCanvasContext() {
  if (sharedCanvasContext !== undefined) {
    return sharedCanvasContext;
  }

  if (typeof document === 'undefined') {
    sharedCanvasContext = null;
    return sharedCanvasContext;
  }

  const canvas = document.createElement('canvas');
  sharedCanvasContext = canvas.getContext('2d');
  return sharedCanvasContext;
}

function toHexSegment(value: number) {
  return Math.max(0, Math.min(255, Math.round(value)))
    .toString(16)
    .padStart(2, '0');
}

function normalizeHexInput(value: string) {
  const normalized = value.trim().replace(/^#?/, '#');
  const compact = normalized.slice(1);

  if (/^[0-9a-fA-F]{3}$/.test(compact)) {
    return `#${compact
      .split('')
      .map((item) => `${item}${item}`)
      .join('')
      .toLowerCase()}`;
  }

  if (/^[0-9a-fA-F]{4}$/.test(compact)) {
    const [red, green, blue] = compact
      .slice(0, 3)
      .split('')
      .map((item) => `${item}${item}`);
    return `#${red}${green}${blue}`.toLowerCase();
  }

  if (/^[0-9a-fA-F]{6}$/.test(compact)) {
    return normalized.toLowerCase();
  }

  if (/^[0-9a-fA-F]{8}$/.test(compact)) {
    return `#${compact.slice(0, 6)}`.toLowerCase();
  }

  return null;
}

function parseHexColor(value: string): ParsedThemeTokenColor | null {
  const normalized = value.trim().replace(/^#?/, '#');
  const compact = normalized.slice(1);

  if (!/^[0-9a-fA-F]{3,8}$/.test(compact) || ![3, 4, 6, 8].includes(compact.length)) {
    return null;
  }

  const expanded =
    compact.length === 3 || compact.length === 4
      ? compact
          .split('')
          .map((item) => `${item}${item}`)
          .join('')
      : compact;
  const red = Number.parseInt(expanded.slice(0, 2), 16);
  const green = Number.parseInt(expanded.slice(2, 4), 16);
  const blue = Number.parseInt(expanded.slice(4, 6), 16);
  const alpha = expanded.length === 8 ? Number.parseInt(expanded.slice(6, 8), 16) / 255 : 1;

  return {
    alpha,
    blue,
    css: alpha >= 0.999 ? `#${expanded.slice(0, 6).toLowerCase()}` : `rgba(${red}, ${green}, ${blue}, ${alpha})`,
    green,
    hex: `#${expanded.slice(0, 6).toLowerCase()}`,
    red,
  };
}

function parseRgbColor(value: string): ParsedThemeTokenColor | null {
  const matched = value.match(
    /^rgba?\(\s*(\d+(?:\.\d+)?)\s*,\s*(\d+(?:\.\d+)?)\s*,\s*(\d+(?:\.\d+)?)\s*(?:,\s*(\d*(?:\.\d+)?)\s*)?\)$/i,
  );

  if (!matched) {
    return null;
  }

  const red = Number(matched[1]);
  const green = Number(matched[2]);
  const blue = Number(matched[3]);
  const alpha = matched[4] === undefined ? 1 : Number(matched[4]);

  if ([red, green, blue].some((item) => Number.isNaN(item) || item < 0 || item > 255)) {
    return null;
  }

  if (Number.isNaN(alpha) || alpha < 0 || alpha > 1) {
    return null;
  }

  return {
    alpha,
    blue,
    css: alpha >= 0.999 ? `#${toHexSegment(red)}${toHexSegment(green)}${toHexSegment(blue)}` : value,
    green,
    hex: `#${toHexSegment(red)}${toHexSegment(green)}${toHexSegment(blue)}`,
    red,
  };
}

function normalizeCssColor(value: string) {
  if (typeof document === 'undefined') {
    return null;
  }

  const style = document.createElement('span').style;
  style.color = '';
  style.color = value.trim();

  if (!style.color) {
    return null;
  }

  const context = getCanvasContext();

  if (!context) {
    return null;
  }

  context.fillStyle = '#010203';

  try {
    context.fillStyle = style.color;
  } catch {
    return null;
  }

  return context.fillStyle;
}

export function isThemeTokenColorKey(tokenKey: string) {
  return /color|background|border/i.test(tokenKey) && !/shadow/i.test(tokenKey);
}

export function parseThemeTokenColor(value: string) {
  const trimmed = value.trim();

  if (!trimmed) {
    return null;
  }

  const hexColor = parseHexColor(trimmed);

  if (hexColor) {
    return hexColor;
  }

  const rgbColor = parseRgbColor(trimmed);

  if (rgbColor) {
    return rgbColor;
  }

  const normalized = normalizeCssColor(trimmed);

  if (!normalized) {
    return null;
  }

  return parseHexColor(normalized) ?? parseRgbColor(normalized);
}

export function buildThemeTokenColorValue(hexValue: string, opacityValue: number) {
  const normalizedHex = normalizeHexInput(hexValue);

  if (!normalizedHex) {
    return null;
  }

  const parsed = parseHexColor(normalizedHex);

  if (!parsed) {
    return null;
  }

  const opacity = Math.max(0, Math.min(100, Math.round(opacityValue)));

  if (opacity >= 100) {
    return parsed.hex;
  }

  return `rgba(${parsed.red}, ${parsed.green}, ${parsed.blue}, ${opacity / 100})`;
}

export function formatThemeTokenSummaryValue(tokenKey: string, value: string) {
  const trimmed = value.trim();

  if (!trimmed) {
    return '--';
  }

  if (!isThemeTokenColorKey(tokenKey)) {
    return trimmed;
  }

  const parsed = parseThemeTokenColor(trimmed);

  if (!parsed) {
    return trimmed;
  }

  const opacity = Math.round(parsed.alpha * 100);
  return opacity >= 100 ? parsed.hex : `${parsed.hex} / ${opacity}%`;
}

export function resolveThemeTokenPreviewHex(value: string) {
  return parseThemeTokenColor(value)?.hex ?? FALLBACK_HEX;
}
