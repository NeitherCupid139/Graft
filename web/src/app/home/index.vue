<template>
  <section class="app-home" data-page-type="overview-dashboard">
    <t-card class="app-home__panel" :bordered="true">
      <div class="app-home__content">
        <span class="app-home__eyebrow">{{ t('app.home.eyebrow') }}</span>
        <h1>{{ t('app.home.title') }}</h1>
        <p>{{ t('app.home.description') }}</p>
      </div>
    </t-card>
  </section>
</template>
<script setup lang="ts">
import { t } from '@/locales';
</script>
<style lang="less" scoped>
.app-home {
  display: flex;
  min-height: 320px;
}

.app-home__panel {
  width: 100%;
}

.app-home__content {
  display: flex;
  flex-direction: column;
  gap: var(--td-comp-margin-s);
  max-width: 640px;
}

.app-home__eyebrow {
  color: var(--td-brand-color);
  font: var(--td-font-body-small);
}

.app-home h1 {
  color: var(--td-text-color-primary);
  font: var(--td-font-headline-medium);
  margin: 0;
}

.app-home p {
  color: var(--td-text-color-secondary);
  font: var(--td-font-body-medium);
  margin: 0;
}
</style>
