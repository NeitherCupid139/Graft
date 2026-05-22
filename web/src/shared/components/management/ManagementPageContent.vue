<template>
  <div class="management-page-content">
    <slot />
  </div>
</template>
<style scoped lang="less">
.management-page-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: 0 auto;
  padding-bottom: max(var(--td-comp-margin-xl), calc(var(--graft-floating-toolbar-safe-area) - 8px));
  width: min(100%, 1320px);
}

@media (width <= 768px) {
  .management-page-content {
    gap: 12px;
    padding-bottom: var(--td-comp-margin-xl);
  }
}
</style>
