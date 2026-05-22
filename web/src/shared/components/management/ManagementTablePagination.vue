<template>
  <div class="management-table-pagination">
    <div class="management-table-pagination__summary">
      <slot name="summary">
        <span>{{ summary }}</span>
      </slot>
    </div>
    <div class="management-table-pagination__controls">
      <slot />
    </div>
  </div>
</template>
<script setup lang="ts">
defineProps<{
  summary: string;
}>();
</script>
<style scoped lang="less">
.management-table-pagination {
  align-items: center;
  display: flex;
  gap: 16px;
  justify-content: space-between;
  min-height: 56px;
  width: 100%;
}

.management-table-pagination__summary {
  color: var(--td-text-color-secondary);
  flex: 0 0 auto;
  font: var(--td-font-body-small);
  min-width: 0;
  white-space: nowrap;
}

.management-table-pagination__controls {
  align-items: center;
  display: flex;
  flex: 1 1 auto;
  justify-content: flex-end;
  min-width: 0;
}

.management-table-pagination__controls :deep(.t-pagination) {
  align-items: center;
  display: flex;
  flex-wrap: nowrap;
  gap: 12px;
  justify-content: flex-end;
  width: auto;
}

.management-table-pagination__controls :deep(.t-pagination__total),
.management-table-pagination__controls :deep(.t-pagination__page-size),
.management-table-pagination__controls :deep(.t-pagination__list),
.management-table-pagination__controls :deep(.t-pagination__jumper) {
  margin: 0;
  white-space: nowrap;
}

.management-table-pagination__controls :deep(.t-pagination__list) {
  align-items: center;
  display: inline-flex;
}

@media (width <= 768px) {
  .management-table-pagination {
    align-items: stretch;
    flex-direction: column;
  }

  .management-table-pagination__summary,
  .management-table-pagination__controls {
    white-space: normal;
    width: 100%;
  }

  .management-table-pagination__controls :deep(.t-pagination) {
    flex-wrap: wrap;
    justify-content: flex-start;
  }
}
</style>
