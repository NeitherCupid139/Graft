<!-- Copyright (c) 2025-2026 GeWuYou -->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <ActionMenu :items="items" more-label-fallback="更多" />
</template>

<script setup lang="ts">
const items = [{ key: 'refresh', labelKey: 'demo.actions.refresh' }];
</script>
