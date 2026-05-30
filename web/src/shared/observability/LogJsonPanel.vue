<template>
  <section class="log-json-panel">
    <h4 class="log-json-panel__title">{{ title }}</h4>
    <t-collapse :value="expandedPanels" borderless expand-icon-placement="right">
      <t-collapse-panel value="json">
        <template #header>
          <div class="log-json-panel__header">
            <span>{{ toggleLabel }}</span>
            <t-button v-if="!isEmpty" size="small" theme="default" variant="text" @click.stop="copyJson">
              {{ copyLabel }}
            </t-button>
          </div>
        </template>
        <div v-if="isEmpty" class="log-json-panel__empty">{{ emptyText }}</div>
        <pre v-else class="log-json-panel__code">{{ formattedJson }}</pre>
      </t-collapse-panel>
    </t-collapse>
  </section>
</template>
<script setup lang="ts">
import { MessagePlugin } from 'tdesign-vue-next';
import { computed } from 'vue';

import { copyText } from './copy';

const props = defineProps<{
  title: string;
  toggleLabel: string;
  copyLabel: string;
  copySuccessLabel: string;
  copyFailLabel: string;
  emptyText: string;
  value: unknown;
}>();

const expandedPanels = ['json'];

const isEmpty = computed(() => {
  const value = props.value;
  if (value === null || value === undefined || value === '') {
    return true;
  }

  if (typeof value === 'object') {
    return Object.keys(value as Record<string, unknown>).length === 0;
  }

  return false;
});

const formattedJson = computed(() => {
  if (isEmpty.value) {
    return '';
  }

  try {
    return JSON.stringify(props.value, null, 2);
  } catch {
    return String(props.value);
  }
});

async function copyJson() {
  try {
    const copied = await copyText(formattedJson.value);
    if (!copied) {
      MessagePlugin.error(props.copyFailLabel);
      return;
    }
    MessagePlugin.success(props.copySuccessLabel);
  } catch {
    MessagePlugin.error(props.copyFailLabel);
  }
}
</script>
<style scoped lang="less">
.log-json-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.log-json-panel__title {
  margin: 0;
}

.log-json-panel__header {
  align-items: center;
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.log-json-panel__empty,
.log-json-panel__code {
  background: var(--td-bg-color-container-hover);
  border: 1px solid var(--td-component-stroke);
  border-radius: var(--td-radius-medium);
  margin: 0;
  overflow: auto;
  overflow-wrap: anywhere;
  padding: 12px;
  white-space: pre-wrap;
}

.log-json-panel__empty {
  color: var(--td-text-color-secondary);
}
</style>
