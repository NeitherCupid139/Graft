<!--
  Copyright (c) 2025-2026 GeWuYou
  SPDX-License-Identifier: Apache-2.0
-->

<template>
  <section class="log-viewer">
    <div class="log-viewer__toolbar">
      <div class="log-viewer__toolbar-group log-viewer__toolbar-group--primary">
        <t-button theme="primary" :loading="loading" @click="$emit('refresh')">
          {{ refreshLabel }}
        </t-button>
        <t-button theme="default" variant="outline" :disabled="!displayLines.length" @click="copyContent">
          {{ copyLabel }}
        </t-button>
        <t-button theme="default" variant="outline" :disabled="!displayLines.length" @click="downloadContent">
          {{ downloadLabel }}
        </t-button>
        <span class="log-viewer__select-wrap">
          <t-select
            v-if="lineLimitOptions.length"
            v-model:value="selectedLineLimit"
            class="log-viewer__limit"
            :options="lineLimitOptions"
            size="small"
            @change="emitLimit"
          />
        </span>
        <t-select v-model:value="selectedLevel" class="log-viewer__level-filter" :options="levelOptions" size="small" />
      </div>

      <div class="log-viewer__toolbar-group log-viewer__toolbar-group--tools">
        <t-input
          v-model:value="searchKeyword"
          class="log-viewer__search"
          clearable
          type="search"
          :placeholder="searchPlaceholder"
        />
        <span v-if="normalizedSearchKeyword" class="log-viewer__match-count">
          {{ matchCountLabel.replace('{count}', String(searchMatchCount)) }}
        </span>
        <label class="log-viewer__switch">
          <span>{{ wrapLabel }}</span>
          <t-switch v-model:value="wrapLines" size="small" />
        </label>
        <label class="log-viewer__switch">
          <span>{{ refreshScrollLabel }}</span>
          <t-switch v-model:value="scrollAfterRefresh" size="small" />
        </label>
      </div>
    </div>

    <t-alert v-if="error" theme="error" :title="error">
      <template #operation>
        <t-button size="small" theme="danger" variant="text" @click="$emit('refresh')">
          {{ retryLabel }}
        </t-button>
      </template>
    </t-alert>
    <t-alert v-if="truncated" theme="warning" :title="truncatedLabel" />

    <div ref="viewport" :class="['log-viewer__viewport', { 'log-viewer__viewport--wrap': wrapLines }]">
      <t-skeleton v-if="loading && !displayLines.length" animation="gradient" :row-col="skeletonRows" />
      <ol v-else-if="displayLines.length" class="log-viewer__lines">
        <li
          v-for="line in displayLines"
          :key="line.lineNo"
          tabindex="0"
          :class="[
            'log-viewer__line',
            `log-viewer__line--${line.tone}`,
            { 'log-viewer__line--expanded': isExpanded(line.lineNo) },
          ]"
          @click="toggleLine(line.lineNo)"
          @keydown.enter.prevent="toggleLine(line.lineNo)"
          @keydown.space.prevent="toggleLine(line.lineNo)"
        >
          <span class="log-viewer__line-number">{{ line.lineNo }}</span>
          <div class="log-viewer__line-body">
            <div class="log-viewer__line-main">
              <t-tooltip v-if="line.timestamp" :content="line.timestamp" placement="top-left" theme="light">
                <time class="log-viewer__timestamp">{{ shortTimestamp(line.timestamp) }}</time>
              </t-tooltip>
              <span v-else class="log-viewer__timestamp">-</span>
              <t-tag class="log-viewer__level" :theme="levelTheme(line.level)" size="small" variant="light-outline">
                {{ line.level ?? 'LOG' }}
              </t-tag>
              <t-tooltip v-if="line.source" :content="line.source" placement="top-left" theme="light">
                <span class="log-viewer__source">{{ line.sourceShort || line.source }}</span>
              </t-tooltip>
              <span v-else class="log-viewer__source">-</span>
              <code class="log-viewer__message">
                <span
                  v-for="(token, tokenIndex) in line.messageTokens"
                  :key="`${line.lineNo}-message-${tokenIndex}`"
                  :class="tokenClass(token)"
                  >{{ token.text }}</span
                >
              </code>
              <div v-if="line.metadata" class="log-viewer__metadata-tags" @click.stop>
                <t-tag
                  v-for="[key, value] in visibleMetadataTags(line)"
                  :key="`${line.lineNo}-${key}`"
                  size="small"
                  theme="default"
                  variant="light"
                >
                  {{ key }}={{ formatMetadataValue(value) }}
                </t-tag>
                <t-button
                  v-if="hiddenMetadataCount(line)"
                  size="small"
                  theme="default"
                  variant="text"
                  @click="toggleLine(line.lineNo)"
                >
                  +{{ hiddenMetadataCount(line) }}
                </t-button>
              </div>
              <div class="log-viewer__row-actions" @click.stop>
                <t-tooltip :content="isExpanded(line.lineNo) ? collapseDetailLabel : viewDetailLabel" theme="light">
                  <t-button
                    :aria-label="isExpanded(line.lineNo) ? collapseDetailLabel : viewDetailLabel"
                    class="log-viewer__icon-action"
                    shape="square"
                    size="small"
                    theme="default"
                    variant="text"
                    @click="toggleLine(line.lineNo)"
                  >
                    <template #icon>
                      <browse-icon />
                    </template>
                  </t-button>
                </t-tooltip>
                <t-tooltip :content="copyLineLabel" theme="light">
                  <t-button
                    :aria-label="copyLineLabel"
                    class="log-viewer__icon-action"
                    shape="square"
                    size="small"
                    theme="default"
                    variant="text"
                    @click="copyLine(line.raw)"
                  >
                    <template #icon>
                      <copy-icon />
                    </template>
                  </t-button>
                </t-tooltip>
              </div>
            </div>

            <div v-if="isExpanded(line.lineNo)" class="log-viewer__line-detail" @click.stop>
              <div class="log-viewer__line-detail-actions">
                <t-button size="small" theme="default" variant="outline" @click.stop="copyMessage(line.message)">
                  {{ copyMessageLabel }}
                </t-button>
                <t-button size="small" theme="default" variant="outline" @click.stop="copyLine(line.raw)">
                  {{ copyLineLabel }}
                </t-button>
                <t-button
                  v-if="line.metadata"
                  size="small"
                  theme="default"
                  variant="outline"
                  @click.stop="copyJson(line.metadata)"
                >
                  {{ copyJsonLabel }}
                </t-button>
              </div>
              <section class="log-viewer__detail-section">
                <strong>{{ messageLabel }}</strong>
                <pre class="log-viewer__message-full">{{ line.message }}</pre>
              </section>
              <section v-if="line.metadata" class="log-viewer__detail-section">
                <strong>{{ metadataLabel }}</strong>
                <pre class="log-viewer__json">{{ formatJson(line.metadata) }}</pre>
              </section>
              <section class="log-viewer__detail-section">
                <strong>{{ rawLabel }}</strong>
                <pre class="log-viewer__raw"><span
                  v-for="(token, tokenIndex) in line.rawTokens"
                  :key="`${line.lineNo}-raw-${tokenIndex}`"
                  :class="tokenClass(token)"
                  >{{ token.text }}</span
                ></pre>
              </section>
            </div>
          </div>
        </li>
      </ol>
      <t-empty v-else size="small" :description="emptyLabel" />
    </div>
  </section>
</template>
<script setup lang="ts">
import { BrowseIcon, CopyIcon } from 'tdesign-icons-vue-next';
import type { SelectProps } from 'tdesign-vue-next';
import { MessagePlugin } from 'tdesign-vue-next/es/message';
import { computed, nextTick, ref, watch } from 'vue';

import { copyText } from './copy';
import type { LogLevel, LogToken } from './log-highlight';
import {
  buildDisplayLogLine,
  type DisplayLogLine,
  formatLogMetadataValue,
  type ParsedLogMetadata,
  parseLogLines,
  summarizeMetadata,
} from './log-parser';

const props = withDefaults(
  defineProps<{
    lines: string[];
    loading?: boolean;
    error?: string;
    truncated?: boolean;
    lineLimit?: number;
    lineLimits?: number[];
    refreshLabel: string;
    copyLabel: string;
    downloadLabel: string;
    retryLabel: string;
    searchPlaceholder: string;
    wrapLabel: string;
    refreshScrollLabel: string;
    levelFilterLabel: string;
    allLevelsLabel: string;
    matchCountLabel: string;
    emptyLabel: string;
    truncatedLabel: string;
    viewDetailLabel: string;
    collapseDetailLabel: string;
    metadataLabel: string;
    messageLabel: string;
    rawLabel: string;
    copyMessageLabel: string;
    copyLineLabel: string;
    copyJsonLabel: string;
    copySuccessLabel: string;
    copyErrorLabel: string;
  }>(),
  {
    loading: false,
    error: '',
    truncated: false,
    lineLimit: 200,
    lineLimits: () => [100, 200, 500, 1000],
  },
);

const emit = defineEmits<{
  refresh: [];
  'update:lineLimit': [value: number];
}>();

type SelectOption = NonNullable<SelectProps['options']>[number];
type LevelFilter = 'ALL' | LogLevel;

const searchKeyword = ref('');
const wrapLines = ref(true);
const scrollAfterRefresh = ref(true);
const selectedLineLimit = ref(props.lineLimit);
const selectedLevel = ref<LevelFilter>('ALL');
const viewport = ref<HTMLElement | null>(null);
const expandedLineNos = ref<Set<number>>(new Set());

const skeletonRows = [
  { height: '22px', width: '96%' },
  { height: '22px', width: '88%' },
  { height: '22px', width: '92%' },
  { height: '22px', width: '76%' },
  { height: '22px', width: '84%' },
];
const levelOptions = computed<SelectOption[]>(() => [
  { label: `${props.levelFilterLabel}: ${props.allLevelsLabel}`, value: 'ALL' },
  { label: 'ERROR', value: 'ERROR' },
  { label: 'WARN', value: 'WARN' },
  { label: 'INFO', value: 'INFO' },
  { label: 'DEBUG', value: 'DEBUG' },
  { label: 'TRACE', value: 'TRACE' },
]);
const lineLimitOptions = computed<SelectOption[]>(() =>
  props.lineLimits.map((value) => ({ label: String(value), value })),
);
const normalizedSearchKeyword = computed(() => searchKeyword.value.trim());
const parsedLines = computed(() => parseLogLines(props.lines));
const visibleRawLines = computed(() => parsedLines.value.slice(-selectedLineLimit.value));
const displayLines = computed(() =>
  visibleRawLines.value
    .filter((line) => selectedLevel.value === 'ALL' || line.level === selectedLevel.value)
    .map((line) => buildDisplayLogLine(line, normalizedSearchKeyword.value)),
);
const searchMatchCount = computed(() => displayLines.value.reduce((total, line) => total + line.searchMatchCount, 0));

watch(
  () => props.lineLimit,
  (value) => {
    selectedLineLimit.value = value;
  },
);

watch(
  () => [props.lines.length, scrollAfterRefresh.value],
  () => {
    if (!scrollAfterRefresh.value) return;
    void nextTick(scrollToBottom);
  },
  { flush: 'post' },
);

function emitLimit(value: SelectProps['value']) {
  if (typeof value === 'number') {
    emit('update:lineLimit', value);
  }
}

async function copyContent() {
  await copyTextWithFeedback(displayLines.value.map((line) => line.raw).join('\n'));
}

async function copyLine(raw: string) {
  await copyTextWithFeedback(raw);
}

async function copyMessage(message: string) {
  await copyTextWithFeedback(message);
}

async function copyJson(metadata: ParsedLogMetadata) {
  await copyTextWithFeedback(formatJson(metadata));
}

function downloadContent() {
  const blob = new Blob([displayLines.value.map((line) => line.raw).join('\n')], { type: 'text/plain;charset=utf-8' });
  const link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = `container-logs-${new Date().toISOString().replace(/[:.]/g, '-')}.log`;
  link.click();
  URL.revokeObjectURL(link.href);
}

function scrollToBottom() {
  const node = viewport.value;
  if (node) {
    node.scrollTop = node.scrollHeight;
  }
}

function toggleLine(lineNo: number) {
  const next = new Set(expandedLineNos.value);
  if (next.has(lineNo)) {
    next.delete(lineNo);
  } else {
    next.add(lineNo);
  }
  expandedLineNos.value = next;
}

function isExpanded(lineNo: number) {
  return expandedLineNos.value.has(lineNo);
}

function visibleMetadataTags(line: DisplayLogLine) {
  return summarizeMetadata(line.metadata).tags;
}

function hiddenMetadataCount(line: DisplayLogLine) {
  return summarizeMetadata(line.metadata).hiddenCount;
}

function formatMetadataValue(value: unknown) {
  return formatLogMetadataValue(value);
}

function formatJson(value: unknown) {
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

function shortTimestamp(timestamp: string) {
  const timeMatch = /(?:T|\s)(\d{2}:\d{2}:\d{2}(?:[.,]\d+)?)/.exec(timestamp);
  return timeMatch?.[1] ?? timestamp;
}

function levelTheme(level: LogLevel | null) {
  if (level === 'ERROR' || level === 'FATAL') return 'danger';
  if (level === 'WARN') return 'warning';
  if (level === 'INFO') return 'primary';
  return 'default';
}

function tokenClass(token: LogToken) {
  return [
    'log-viewer__token',
    `log-viewer__token--${token.type}`,
    token.level ? `log-viewer__token--level-${token.level.toLowerCase()}` : '',
  ];
}

async function copyTextWithFeedback(value: string) {
  try {
    const copied = await copyText(value);
    if (!copied) {
      MessagePlugin.error(props.copyErrorLabel);
      return;
    }
    MessagePlugin.success(props.copySuccessLabel);
  } catch {
    MessagePlugin.error(props.copyErrorLabel);
  }
}
</script>
<style scoped lang="less">
.log-viewer {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-10);
  min-width: 0;
}

.log-viewer__toolbar {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: var(--graft-density-gap-8);
  justify-content: space-between;
  min-width: 0;
  position: sticky;
  top: 0;
  z-index: 1;
}

.log-viewer__toolbar-group,
.log-viewer__switch {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: var(--graft-density-gap-6);
  min-width: 0;
}

.log-viewer__toolbar-group--primary {
  flex: 1 1 520px;
}

.log-viewer__toolbar-group--tools {
  flex: 1 1 520px;
  justify-content: flex-end;
}

.log-viewer__limit {
  width: 96px;
}

.log-viewer__level-filter {
  width: 128px;
}

.log-viewer__search {
  width: min(340px, 100%);
}

.log-viewer__match-count,
.log-viewer__switch {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
}

.log-viewer__viewport {
  background: color-mix(in srgb, var(--td-bg-color-page) 92%, black 8%);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-medium);
  color: var(--td-text-color-primary);
  max-height: min(62vh, 680px);
  min-height: 360px;
  min-width: 0;
  overflow: auto;
  padding: var(--graft-density-gap-8) var(--graft-density-gap-12);
  scrollbar-color: var(--td-scrollbar-color) transparent;
  scrollbar-gutter: stable;
  scrollbar-width: thin;
}

.log-viewer__viewport::-webkit-scrollbar {
  background: transparent;
  height: 8px;
  width: 8px;
}

.log-viewer__viewport::-webkit-scrollbar-track {
  background: transparent;
}

.log-viewer__viewport::-webkit-scrollbar-thumb {
  background-clip: content-box;
  background-color: var(--td-scrollbar-color);
  border: 2px solid transparent;
  border-radius: 6px;
}

.log-viewer__lines {
  counter-reset: none;
  list-style: none;
  margin: 0;
  min-width: max-content;
  padding: 0;
}

.log-viewer__line {
  border-left: var(--graft-density-gap-2) solid transparent;
  border-radius: var(--td-radius-small);
  display: grid;
  grid-template-columns: 48px minmax(0, 1fr);
  min-height: 36px;
  padding: var(--graft-density-gap-4) var(--graft-density-gap-6) var(--graft-density-gap-4) 0;
}

.log-viewer__line:hover,
.log-viewer__line:focus-within {
  background: color-mix(in srgb, var(--td-bg-color-container-hover) 72%, transparent);
  outline: none;
}

.log-viewer__line-number {
  color: var(--td-text-color-placeholder);
  font-family: var(--td-font-family-monospace);
  padding: var(--graft-density-gap-6) var(--graft-density-gap-8);
  text-align: right;
  user-select: none;
}

.log-viewer__line-body {
  min-width: 0;
}

.log-viewer__line-main {
  align-items: center;
  display: grid;
  gap: var(--graft-density-gap-8);
  grid-template-columns: 88px 56px 156px minmax(420px, 1fr) minmax(0, 260px) 56px;
  min-width: max-content;
}

.log-viewer__timestamp,
.log-viewer__source,
.log-viewer__message,
.log-viewer__raw,
.log-viewer__json {
  font-family: var(--td-font-family-monospace);
}

.log-viewer__timestamp,
.log-viewer__source {
  color: var(--td-text-color-placeholder);
  padding-top: 0;
}

.log-viewer__source {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.log-viewer__message {
  color: var(--td-text-color-primary);
  line-height: var(--td-line-height-body-medium);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: pre;
}

.log-viewer__metadata-tags {
  display: flex;
  flex-wrap: nowrap;
  gap: var(--graft-density-gap-2);
  max-width: 100%;
  min-width: 0;
  opacity: 0.74;
  overflow: hidden;
}

.log-viewer__metadata-tags :deep(.t-tag) {
  max-width: 148px;
  overflow: hidden;
  padding: 0 var(--graft-density-gap-6);
  text-overflow: ellipsis;
  white-space: nowrap;
}

.log-viewer__row-actions {
  align-items: center;
  display: flex;
  gap: var(--graft-density-gap-2);
  justify-content: flex-end;
  opacity: 0;
  transition: opacity 0.16s ease;
}

.log-viewer__line:hover .log-viewer__row-actions,
.log-viewer__line:focus-within .log-viewer__row-actions,
.log-viewer__line--expanded .log-viewer__row-actions {
  opacity: 1;
}

.log-viewer__icon-action {
  color: var(--td-text-color-secondary);
}

.log-viewer__viewport--wrap .log-viewer__lines,
.log-viewer__viewport--wrap .log-viewer__line-main {
  min-width: 0;
}

.log-viewer__viewport--wrap .log-viewer__line-main {
  grid-template-columns: 76px 56px 148px minmax(0, 1fr) 56px;
  min-width: 0;
}

.log-viewer__viewport--wrap .log-viewer__message {
  -webkit-box-orient: vertical;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  overflow: hidden;
  overflow-wrap: anywhere;
  white-space: pre-wrap;
}

.log-viewer__viewport--wrap .log-viewer__metadata-tags {
  grid-column: 4 / 5;
  max-height: 22px;
}

.log-viewer__viewport--wrap .log-viewer__row-actions {
  grid-column: 5 / 6;
  grid-row: 1 / 2;
}

.log-viewer__line--expanded .log-viewer__message {
  display: block;
  -webkit-line-clamp: unset;
  overflow: visible;
}

.log-viewer__line--danger {
  background: color-mix(in srgb, var(--td-error-color-5) 4%, transparent);
  border-left-color: var(--td-error-color);
}

.log-viewer__line--warning {
  background: color-mix(in srgb, var(--td-warning-color-5) 4%, transparent);
  border-left-color: var(--td-warning-color);
}

.log-viewer__line--muted {
  color: var(--td-text-color-secondary);
}

.log-viewer__line-detail {
  background: color-mix(in srgb, var(--td-bg-color-container) 82%, transparent);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-medium);
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-10);
  margin: var(--graft-density-gap-8) 0 var(--graft-density-gap-8);
  max-width: 100%;
  padding: var(--graft-density-gap-10);
}

.log-viewer__line-detail-actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--graft-density-gap-8);
}

.log-viewer__detail-section {
  display: flex;
  flex-direction: column;
  gap: var(--graft-density-gap-6);
  min-width: 0;
}

.log-viewer__detail-section strong {
  color: var(--td-text-color-secondary);
}

.log-viewer__json,
.log-viewer__raw,
.log-viewer__message-full {
  background: var(--td-bg-color-container-hover);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-small);
  margin: 0;
  max-height: 260px;
  overflow: auto;
  padding: var(--graft-density-gap-10);
  white-space: pre-wrap;
}

.log-viewer__token--keyword {
  background: color-mix(in srgb, var(--td-warning-color-5) 34%, transparent);
  border-radius: var(--td-radius-small);
  color: var(--td-text-color-primary);
}

.log-viewer__token--field-key {
  color: var(--td-brand-color);
}

.log-viewer__token--field-value {
  color: var(--td-text-color-secondary);
}

.log-viewer__token--level-error,
.log-viewer__token--level-fatal {
  color: var(--td-error-color);
  font-weight: 600;
}

.log-viewer__token--level-warn {
  color: var(--td-warning-color);
  font-weight: 600;
}

.log-viewer__token--level-info {
  color: var(--td-brand-color);
  font-weight: 600;
}

.log-viewer__token--level-debug,
.log-viewer__token--level-trace {
  color: var(--td-text-color-placeholder);
}

@media (width <= 1024px) {
  .log-viewer__toolbar-group--tools {
    justify-content: flex-start;
  }

  .log-viewer__viewport--wrap .log-viewer__line-main {
    grid-template-columns: 72px 56px minmax(0, 1fr) 56px;
  }

  .log-viewer__viewport--wrap .log-viewer__source,
  .log-viewer__viewport--wrap .log-viewer__metadata-tags,
  .log-viewer__viewport--wrap .log-viewer__row-actions {
    grid-column: 3 / 4;
  }
}
</style>
