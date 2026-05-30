<template>
  <management-toolbar>
    <template #filters>
      <div class="access-log-filters">
        <t-input
          :model-value="modelValue.requestId"
          :placeholder="t('accessLog.filters.requestId')"
          clearable
          @update:model-value="updateField('requestId', $event)"
        />
        <t-input
          :model-value="modelValue.userId"
          :placeholder="t('accessLog.filters.userId')"
          clearable
          @update:model-value="updateField('userId', $event)"
        />
        <t-input
          :model-value="modelValue.username"
          :placeholder="t('accessLog.filters.username')"
          clearable
          @update:model-value="updateField('username', $event)"
        />
        <t-select
          :model-value="modelValue.method"
          clearable
          :options="methodOptions"
          :placeholder="t('accessLog.filters.method')"
          @update:model-value="updateField('method', normalizeSelectValue($event))"
        />
        <t-input
          :model-value="modelValue.path"
          :placeholder="t('accessLog.filters.path')"
          clearable
          @update:model-value="updateField('path', $event)"
        />
        <t-select
          :model-value="modelValue.pathMatch"
          :options="pathMatchOptions"
          :placeholder="t('accessLog.filters.pathMatch')"
          @update:model-value="updateField('pathMatch', normalizePathMatch($event))"
        />
        <t-input
          :model-value="modelValue.route"
          :placeholder="t('accessLog.filters.route')"
          clearable
          @update:model-value="updateField('route', $event)"
        />
        <t-input
          :model-value="modelValue.statusCode"
          :placeholder="t('accessLog.filters.statusCode')"
          clearable
          @update:model-value="updateField('statusCode', $event)"
        />
        <t-input
          :model-value="modelValue.durationMinMs"
          :placeholder="t('accessLog.filters.durationMin')"
          clearable
          @update:model-value="updateField('durationMinMs', $event)"
        />
        <t-input
          :model-value="modelValue.durationMaxMs"
          :placeholder="t('accessLog.filters.durationMax')"
          clearable
          @update:model-value="updateField('durationMaxMs', $event)"
        />
        <t-date-range-picker
          :model-value="modelValue.occurredRange"
          allow-input
          clearable
          enable-time-picker
          format="YYYY-MM-DD HH:mm:ss"
          @update:model-value="updateField('occurredRange', $event)"
        />
        <t-select
          :model-value="modelValue.sortBy"
          :options="sortByOptions"
          :placeholder="t('accessLog.filters.sortBy')"
          @update:model-value="updateField('sortBy', normalizeSortBy($event))"
        />
        <t-select
          :model-value="modelValue.sortOrder"
          :options="sortOrderOptions"
          :placeholder="t('accessLog.filters.sortOrder')"
          @update:model-value="updateField('sortOrder', normalizeSortOrder($event))"
        />
        <div class="access-log-filters__actions">
          <t-button theme="primary" :loading="loading" @click="$emit('search')">{{
            t('accessLog.actions.search')
          }}</t-button>
          <t-button theme="default" variant="outline" @click="$emit('reset')">{{
            t('accessLog.actions.reset')
          }}</t-button>
        </div>
      </div>
    </template>
  </management-toolbar>
</template>
<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { ManagementToolbar } from '@/shared/components/management';

import type {
  AccessLogFilterState,
  AccessLogPathMatch,
  AccessLogSortBy,
  AccessLogSortOrder,
} from '../types/access-log';

const props = defineProps<{
  loading?: boolean;
  modelValue: AccessLogFilterState;
}>();

const emit = defineEmits<{
  (e: 'reset'): void;
  (e: 'search'): void;
  (e: 'update:modelValue', value: AccessLogFilterState): void;
}>();

const { t } = useI18n();

const methodOptions = computed(() =>
  ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'].map((value) => ({ label: value, value })),
);
const pathMatchOptions = computed(() => [
  { label: t('accessLog.filters.pathMatchExact'), value: 'exact' },
  { label: t('accessLog.filters.pathMatchPrefix'), value: 'prefix' },
]);
const sortByOptions = computed(() => [
  { label: t('accessLog.filters.sortOccurredAt'), value: 'occurred_at' },
  { label: t('accessLog.filters.sortDuration'), value: 'duration_ms' },
  { label: t('accessLog.filters.sortStatusCode'), value: 'status_code' },
]);
const sortOrderOptions = computed(() => [
  { label: t('accessLog.filters.sortDesc'), value: 'desc' },
  { label: t('accessLog.filters.sortAsc'), value: 'asc' },
]);

function updateField<Key extends keyof AccessLogFilterState>(key: Key, value: AccessLogFilterState[Key]) {
  emit('update:modelValue', {
    ...props.modelValue,
    [key]: typeof value === 'string' ? value.trim() : value,
  });
}

function normalizeSelectValue(value: unknown) {
  return typeof value === 'string' ? value : '';
}

function normalizePathMatch(value: unknown): AccessLogPathMatch {
  return value === 'prefix' ? 'prefix' : 'exact';
}

function normalizeSortBy(value: unknown): AccessLogSortBy {
  return value === 'duration_ms' || value === 'status_code' ? value : 'occurred_at';
}

function normalizeSortOrder(value: unknown): AccessLogSortOrder {
  return value === 'asc' ? 'asc' : 'desc';
}
</script>
<style scoped lang="less">
.access-log-filters {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.access-log-filters__actions {
  display: flex;
  gap: 12px;
}

@media (width <= 1024px) {
  .access-log-filters {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (width <= 768px) {
  .access-log-filters {
    grid-template-columns: 1fr;
  }
}
</style>
