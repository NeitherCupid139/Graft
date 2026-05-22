<template>
  <section class="management-stats-grid">
    <article v-for="item in items" :key="item.label" class="management-stats-grid__item">
      <div class="management-stats-grid__head">
        <span class="management-stats-grid__label">{{ item.label }}</span>
        <span v-if="item.tip" class="management-stats-grid__tip">{{ item.tip }}</span>
      </div>
      <strong class="management-stats-grid__value">{{ item.value }}</strong>
      <p v-if="item.description" class="management-stats-grid__description">{{ item.description }}</p>
    </article>
  </section>
</template>
<script setup lang="ts">
export type ManagementStatItem = {
  label: string;
  value: string | number;
  description?: string;
  tip?: string;
};

defineProps<{
  items: ManagementStatItem[];
}>();
</script>
<style scoped lang="less">
.management-stats-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.management-stats-grid__item {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-large);
  box-shadow: var(--td-shadow-1);
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 132px;
  padding: 18px 20px;
}

.management-stats-grid__head {
  color: var(--td-text-color-secondary);
  display: flex;
  font: var(--td-font-body-small);
  gap: 8px;
  justify-content: space-between;
}

.management-stats-grid__tip,
.management-stats-grid__description {
  color: var(--td-text-color-placeholder);
}

.management-stats-grid__value {
  color: var(--td-text-color-primary);
  font: var(--td-font-headline-large);
  line-height: 1.1;
}

.management-stats-grid__description {
  font: var(--td-font-body-small);
  margin: 0;
}
</style>
