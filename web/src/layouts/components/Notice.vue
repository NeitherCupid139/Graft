<!--
  Copyright (c) 2025-2026 GeWuYou
  SPDX-License-Identifier: Apache-2.0
-->

<template>
  <component :is="announcementHeaderEntry" />
  <component :is="notificationHeaderWidget" />
</template>
<script setup lang="ts">
import { announcementHeaderEntry } from '@/modules/announcement';
import { notificationHeaderWidget } from '@/modules/notification';
</script>
