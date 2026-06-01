<template>
  <log-filter-builder
    :active-preset="activePreset"
    :add-filter-label="`+ ${t('accessLog.actions.addFilter')}`"
    :add-sorter-label="t('accessLog.actions.addSorter')"
    :builder-hint="t('accessLog.builder.hint')"
    :builder-title="t('accessLog.builder.title')"
    :field-values="fieldValues"
    :fields="definitions"
    :filters-group-label="t('accessLog.builder.groups.filters')"
    :keyword="modelValue.keyword"
    :keyword-placeholder="t('accessLog.page.searchPlaceholder')"
    :loading="loading"
    :move-down-label="t('accessLog.actions.moveSorterDown')"
    :move-up-label="t('accessLog.actions.moveSorterUp')"
    :preset-label="t('accessLog.presets.label')"
    :presets="presets"
    :remove-sorter-label="t('accessLog.actions.removeSorter')"
    :reset-label="t('accessLog.actions.reset')"
    :search-label="t('accessLog.actions.search')"
    :selected-field-key="selectedFieldKey"
    :sort-add-disabled="sortAddDisabled"
    :sort-direction-options="sortOrderOptions"
    :sort-direction-placeholder="t('accessLog.sort.directionPlaceholder')"
    :sort-field-key="'sorterBuilder'"
    :sort-field-options-by-index="sortFieldOptionsByIndex"
    :sort-field-placeholder="t('accessLog.sort.fieldPlaceholder')"
    :sort-move-down-disabled="sortMoveDownDisabled"
    :sort-move-up-disabled="sortMoveUpDisabled"
    :sorters="normalizedSorters"
    :tags="activeFilterTags"
    :time-field-key="'timeRange'"
    :time-fields="timeFields"
    v-on="builderListeners"
  />
</template>
<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import {
  appendSorterToState,
  moveSorterInState,
  prependSorterTags,
  removeSorterFromState,
  withSorterDirectionFromInput,
  withSorterFieldFromInput,
} from '@/shared/observability';
import type {
  LogFilterFieldDefinition,
  LogFilterTag,
  LogTimeRangeField,
} from '@/shared/observability/log-filter-builder';
import {
  createLogBuilderListeners,
  createSortDirection,
  useLogSorterUiState,
} from '@/shared/observability/log-filter-builder-helpers';
import LogFilterBuilder from '@/shared/observability/LogFilterBuilder.vue';

import type {
  AccessLogFilterState,
  AccessLogPathMatch,
  AccessLogSortBy,
  AccessLogSortOrder,
} from '../types/access-log';

type AccessLogPresetKey =
  | 'all'
  | 'todayErrors'
  | 'status4xx'
  | 'status5xx'
  | 'slowRequests'
  | 'currentUser'
  | 'lastHour';
type AccessTimeRangeKey = 'startedRange' | 'occurredRange';
type FilterKey = Exclude<
  keyof AccessLogFilterState,
  'keyword' | 'pathMatch' | 'route' | 'sorters' | 'startedRange' | 'occurredRange'
>;
type BuilderFieldKey = 'timeRange' | 'sorterBuilder' | FilterKey;
type TagKey = FilterKey | AccessTimeRangeKey | `sorter:${number}`;

const props = defineProps<{
  activePreset: AccessLogPresetKey;
  loading?: boolean;
  modelValue: AccessLogFilterState;
  presets: { key: AccessLogPresetKey; title: string }[];
}>();

const emit = defineEmits<{
  (e: 'apply-preset', preset: AccessLogPresetKey): void;
  (e: 'reset'): void;
  (e: 'search'): void;
  (e: 'update:modelValue', value: AccessLogFilterState): void;
}>();

const { t } = useI18n();

const selectedFieldKey = ref<BuilderFieldKey>('timeRange');

const methodOptions = computed(() =>
  ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'].map((value) => ({ label: value, value })),
);
const sortByOptions = computed(() => [
  { label: t('accessLog.filters.sortStartedAt'), value: 'started_at' },
  { label: t('accessLog.filters.sortOccurredAt'), value: 'occurred_at' },
  { label: t('accessLog.filters.sortDuration'), value: 'duration_ms' },
  { label: t('accessLog.filters.sortStatusCode'), value: 'status_code' },
]);
const sortOrderOptions = computed(() => [
  { label: t('accessLog.filters.sortDesc'), value: 'desc' },
  { label: t('accessLog.filters.sortAsc'), value: 'asc' },
]);
const { normalizedSorters, sortFieldOptionsByIndex, sortAddDisabled, sortMoveUpDisabled, sortMoveDownDisabled } =
  useLogSorterUiState(
    () => props.modelValue.sorters,
    () => sortByOptions.value,
  );

const definitions = computed<LogFilterFieldDefinition[]>(() => [
  { key: 'timeRange', kind: 'special', label: t('accessLog.builder.fields.timeRange') },
  { key: 'sorterBuilder', kind: 'special', label: t('accessLog.builder.fields.sorterBuilder') },
  {
    key: 'requestId',
    kind: 'text',
    label: t('accessLog.builder.fields.requestId'),
    placeholder: t('accessLog.filters.requestId'),
  },
  {
    key: 'userId',
    kind: 'text',
    label: t('accessLog.builder.fields.userId'),
    placeholder: t('accessLog.filters.userId'),
  },
  {
    key: 'username',
    kind: 'text',
    label: t('accessLog.builder.fields.username'),
    placeholder: t('accessLog.filters.username'),
  },
  {
    key: 'method',
    kind: 'select',
    label: t('accessLog.builder.fields.method'),
    placeholder: t('accessLog.filters.method'),
    options: methodOptions.value,
  },
  {
    key: 'path',
    kind: 'text',
    label: t('accessLog.builder.fields.path'),
    placeholder: t('accessLog.filters.path'),
  },
  {
    key: 'statusCode',
    kind: 'text',
    label: t('accessLog.builder.fields.statusCode'),
    placeholder: t('accessLog.filters.statusCode'),
  },
  {
    key: 'durationMinMs',
    kind: 'text',
    label: t('accessLog.builder.fields.durationMinMs'),
    placeholder: t('accessLog.filters.durationMin'),
  },
  {
    key: 'durationMaxMs',
    kind: 'text',
    label: t('accessLog.builder.fields.durationMaxMs'),
    placeholder: t('accessLog.filters.durationMax'),
  },
]);

const fieldValues = computed<Record<string, string | string[]>>(() => ({
  requestId: props.modelValue.requestId,
  userId: props.modelValue.userId,
  username: props.modelValue.username,
  method: props.modelValue.method,
  path: props.modelValue.path,
  statusCode: props.modelValue.statusCode,
  durationMinMs: props.modelValue.durationMinMs,
  durationMaxMs: props.modelValue.durationMaxMs,
}));

const timeFields = computed<LogTimeRangeField[]>(() => [
  {
    key: 'startedRange',
    label: t('accessLog.filters.startedRange'),
    value: props.modelValue.startedRange,
    placeholder: [t('accessLog.filters.startedRange'), t('accessLog.filters.startedRange')],
  },
  {
    key: 'occurredRange',
    label: t('accessLog.filters.occurredRange'),
    value: props.modelValue.occurredRange,
    placeholder: [t('accessLog.filters.occurredRange'), t('accessLog.filters.occurredRange')],
  },
]);

const builderListeners = createAccessLogBuilderListeners();

const activeFilterTags = computed<LogFilterTag[]>(() => {
  const filterTags = definitions.value
    .filter((definition) => definition.kind !== 'special')
    .map((definition) => {
      const rawValue = props.modelValue[definition.key as FilterKey];
      const value = typeof rawValue === 'string' ? rawValue.trim() : rawValue;
      if (!value) {
        return null;
      }
      const label =
        definition.kind === 'select'
          ? definition.options?.find((option) => option.value === value)?.label || String(value)
          : String(value);
      return { key: definition.key, label: `${definition.label}：${label}` };
    })
    .filter((item): item is LogFilterTag => Boolean(item));

  const timeTags: LogFilterTag[] = [];
  const startedTimeTag = buildTimeTag('startedRange', t('accessLog.filters.startedRange'));
  if (startedTimeTag) {
    timeTags.push({ key: 'startedRange', label: startedTimeTag });
  }
  const occurredTimeTag = buildTimeTag('occurredRange', t('accessLog.filters.occurredRange'));
  if (occurredTimeTag) {
    timeTags.push({ key: 'occurredRange', label: occurredTimeTag });
  }

  return prependSorterTags(
    [...timeTags, ...filterTags],
    normalizedSorters.value,
    sortByOptions.value,
    t('accessLog.sort.tagPrefix'),
  );
});

function updateField<Key extends keyof AccessLogFilterState>(key: Key, value: AccessLogFilterState[Key]) {
  emit('update:modelValue', {
    ...props.modelValue,
    [key]: typeof value === 'string' ? value.trim() : value,
  });
}

function handleFieldUpdate(payload: { key: string; value: string | string[] }) {
  updateField(payload.key as keyof AccessLogFilterState, payload.value as never);
}

function createAccessLogBuilderListeners() {
  return createLogBuilderListeners<AccessLogPresetKey, BuilderFieldKey, { key: string; value: string[] }>({
    handleFieldUpdate,
    selectedFieldKey,
    updateKeyword: (value) => updateField('keyword', value),
    emitSearch: () => emit('search'),
    addSorter,
    updateSortField,
    clearTag: (key) => clearTag(key as TagKey),
    moveSorterUp,
    emitReset: () => emit('reset'),
    updateTimeField: ({ key, value }) => updateTimeField(key as AccessTimeRangeKey, value),
    removeSorter,
    updateSortDirection,
    moveSorterDown,
    emitApplyPreset: (preset) => emit('apply-preset', preset),
  });
}

function clearField(key: FilterKey) {
  if (key === 'method') {
    updateField('method', '');
    return;
  }
  updateField(key, '');
}

function clearTag(key: TagKey) {
  if (key === 'startedRange' || key === 'occurredRange') {
    updateTimeField(key, []);
    return;
  }

  if (key.startsWith('sorter:')) {
    removeSorter(Number(key.split(':')[1] || 0));
    return;
  }

  clearField(key as FilterKey);
}

function normalizeSortBy(value: string): AccessLogSortBy {
  return value === 'occurred_at' || value === 'duration_ms' || value === 'status_code' ? value : 'started_at';
}

function normalizeSortOrder(value: string): AccessLogSortOrder {
  return createSortDirection(value);
}

function addSorter() {
  emit('update:modelValue', appendSorterToState(props.modelValue, sortByOptions.value));
}

function removeSorter(index: number) {
  emit('update:modelValue', removeSorterFromState(props.modelValue, index, sortByOptions.value));
}

function moveSorterUp(index: number) {
  emit('update:modelValue', moveSorterInState(props.modelValue, index, -1, sortByOptions.value));
}

function moveSorterDown(index: number) {
  emit('update:modelValue', moveSorterInState(props.modelValue, index, 1, sortByOptions.value));
}

function updateSortField(index: number, value: string | number | Array<string | number> | undefined) {
  emit(
    'update:modelValue',
    withSorterFieldFromInput(props.modelValue, index, value, normalizeSortBy, sortByOptions.value, 'desc'),
  );
}

function updateSortDirection(index: number, value: string | number | Array<string | number> | undefined) {
  emit(
    'update:modelValue',
    withSorterDirectionFromInput(props.modelValue, index, value, normalizeSortOrder, sortByOptions.value),
  );
}

function updateTimeField(key: AccessTimeRangeKey, value: string[]) {
  emit('update:modelValue', {
    ...props.modelValue,
    [key]: value,
  });
}

function buildTimeTag(key: AccessTimeRangeKey, label: string) {
  const range = props.modelValue[key];
  if (!range.length) {
    return '';
  }

  return `${label}：${range.filter(Boolean).join(' ~ ')}`;
}

void (null as unknown as AccessLogPathMatch);
</script>
