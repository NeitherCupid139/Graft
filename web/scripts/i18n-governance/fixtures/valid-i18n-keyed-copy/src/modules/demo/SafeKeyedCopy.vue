<!-- Copyright (c) 2025-2026 GeWuYou -->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
const { t } = useI18n();

const actions = [
  {
    key: 'edit',
    labelKey: 'demo.actions.edit',
  },
  {
    key: 'toggleUser',
    labelKey: 'demo.actions.toggleUser',
  },
];

const title = t('demo.title');
</script>

<template>
  <section>
    <h2>{{ title }}</h2>
    <button v-for="action in actions" :key="action.key">
      {{ t(action.labelKey) }}
    </button>
  </section>
</template>
