<template>
  <section class="management-table-card">
    <div v-if="$slots.head" class="management-table-card__head">
      <slot name="head" />
    </div>
    <div v-if="$slots.batch" class="management-table-card__batch">
      <slot name="batch" />
    </div>
    <div class="management-table-card__body">
      <slot />
    </div>
    <div v-if="$slots.footer" class="management-table-card__footer">
      <slot name="footer" />
    </div>
  </section>
</template>
<style scoped lang="less">
.management-table-card {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-large);
  box-shadow: var(--td-shadow-1);
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px;
}

.management-table-card__head,
.management-table-card__batch,
.management-table-card__footer {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: space-between;
}

.management-table-card__batch {
  background: var(--td-bg-color-secondarycontainer);
  border: 1px solid var(--td-component-border);
  border-radius: var(--td-radius-medium);
  padding: 12px 14px;
}

.management-table-card__body {
  min-width: 0;
}

@media (width <= 768px) {
  .management-table-card {
    padding: 16px;
  }

  .management-table-card__head,
  .management-table-card__batch,
  .management-table-card__footer {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
