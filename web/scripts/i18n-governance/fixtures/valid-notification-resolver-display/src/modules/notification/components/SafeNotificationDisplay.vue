<template>
  <section>
    <strong>{{ resolveNotificationTitle(item, t) }}</strong>
    <span>{{ resolveNotificationMessage(item, t) }}</span>
    <span>{{ resolveNotificationResourceType(item, t) }}</span>
    <code>{{ formatNotificationDiagnosticValue(item.event_type, t) }}</code>
  </section>
</template>
<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import {
  formatNotificationDiagnosticValue,
  resolveNotificationMessage,
  resolveNotificationResourceType,
  resolveNotificationTitle,
} from '../shared/presentation';

const { t } = useI18n();
defineProps<{ item: { event_type: string; message: string; resource_type: string; title: string } }>();
</script>
