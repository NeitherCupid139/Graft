<!-- Copyright (c) 2025-2026 GeWuYou -->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
const condition = true;
const userAction = {
  key: 'toggleUser',
  fallbackLabel: condition ? '启用用户' : '禁用用户',
};
</script>

<template>
  <button>{{ userAction.key }}</button>
</template>
