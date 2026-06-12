<!--
  Copyright (c) 2025-2026 GeWuYou
  SPDX-License-Identifier: Apache-2.0
-->

<template>
  <t-badge :count="unreadCount" :max-count="99" :offset="[4, 4]">
    <t-button
      theme="default"
      shape="square"
      variant="text"
      :loading="loading"
      :aria-label="t('announcement.header.open')"
      :title="t('announcement.header.open')"
      @click="openAnnouncements"
    >
      <t-icon name="notification" />
    </t-button>
  </t-badge>
</template>
<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { getAnnouncementUnreadCount } from '../api/announcement';
import { ANNOUNCEMENT_ROUTE_PATH } from '../contract/paths';
import { ANNOUNCEMENT_HEADER_REFRESH_EVENT } from '../contract/refresh';

const { t } = useI18n();
const router = useRouter();

const loading = ref(false);
const unreadCount = ref(0);

onMounted(() => {
  void refreshUnreadCount();
  window.addEventListener(ANNOUNCEMENT_HEADER_REFRESH_EVENT, refreshUnreadCount);
});

onBeforeUnmount(() => {
  window.removeEventListener(ANNOUNCEMENT_HEADER_REFRESH_EVENT, refreshUnreadCount);
});

async function refreshUnreadCount() {
  loading.value = true;
  try {
    const response = await getAnnouncementUnreadCount();
    unreadCount.value = response.count;
  } catch {
    unreadCount.value = 0;
  } finally {
    loading.value = false;
  }
}

function openAnnouncements() {
  void router.push(ANNOUNCEMENT_ROUTE_PATH.USER_LIST);
}
</script>
