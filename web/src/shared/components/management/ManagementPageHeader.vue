<template>
  <header class="management-page-header">
    <div class="management-page-header__copy">
      <div v-if="$slots.eyebrow" class="management-page-header__eyebrow">
        <slot name="eyebrow" />
      </div>
      <h1 class="management-page-header__title">{{ title }}</h1>
      <p v-if="description" class="management-page-header__description">{{ description }}</p>
    </div>
    <div class="management-page-header__side">
      <div v-if="$slots.meta" class="management-page-header__meta">
        <slot name="meta" />
      </div>
      <div v-if="$slots.actions" class="management-page-header__actions">
        <slot name="actions" />
      </div>
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
@import './card-surface.less';

.management-page-header {
  .management-card-surface();

  display: flex;
  gap: var(--graft-density-gap-16);
  justify-content: space-between;
  padding: var(--graft-density-gap-18) var(--graft-density-gap-20);
}

.management-page-header__side,
.management-page-header__actions {
  display: flex;
  gap: var(--graft-density-gap-12);
}

.management-page-header__copy {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: var(--graft-density-gap-4);
  min-width: 0;
}

.management-page-header__eyebrow {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-small);
}

.management-page-header__title {
  color: var(--td-text-color-primary);
  font: var(--td-font-title-large);
  margin: 0;
}

.management-page-header__description {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-medium);
  margin: 0;
  max-width: 720px;
}

.management-page-header__side {
  align-items: flex-start;
  flex: 0 0 auto;
  flex-direction: column;
}

.management-page-header__meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--graft-density-gap-12);
}

.management-page-header__actions {
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
}

@media (width <= 768px) {
  .management-page-header {
    flex-direction: column;
    padding: var(--graft-density-gap-16);
  }

  .management-page-header__side,
  .management-page-header__actions {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
