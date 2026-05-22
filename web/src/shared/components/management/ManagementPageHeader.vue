<template>
  <header class="management-page-header">
    <div class="management-page-header__main">
      <div class="management-page-header__copy">
        <div v-if="$slots.eyebrow" class="management-page-header__eyebrow">
          <slot name="eyebrow" />
        </div>
        <h1 class="management-page-header__title">{{ title }}</h1>
        <p v-if="description" class="management-page-header__description">{{ description }}</p>
      </div>
      <div v-if="$slots.meta" class="management-page-header__meta">
        <slot name="meta" />
      </div>
    </div>
    <div v-if="$slots.actions" class="management-page-header__actions">
      <slot name="actions" />
    </div>
  </header>
</template>
<script setup lang="ts">
defineProps<{
  title: string;
  description?: string;
}>();
</script>
<style scoped lang="less">
.management-page-header {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-large);
  box-shadow: var(--td-shadow-1);
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 24px;
}

.management-page-header__main,
.management-page-header__actions {
  display: flex;
  gap: 16px;
  justify-content: space-between;
}

.management-page-header__main {
  align-items: flex-start;
}

.management-page-header__copy {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.management-page-header__eyebrow {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
}

.management-page-header__title {
  color: var(--td-text-color-primary);
  font: var(--td-font-headline-medium);
  margin: 0;
}

.management-page-header__description {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-medium);
  margin: 0;
  max-width: 760px;
}

.management-page-header__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.management-page-header__actions {
  align-items: center;
  flex-wrap: wrap;
}

@media (width <= 768px) {
  .management-page-header {
    padding: 20px 16px;
  }

  .management-page-header__main,
  .management-page-header__actions {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
