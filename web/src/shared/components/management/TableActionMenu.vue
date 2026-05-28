<template>
  <div class="table-action-menu">
    <t-button
      v-if="primaryAction"
      :data-testid="primaryAction.testId"
      :disabled="primaryAction.disabled"
      size="small"
      theme="default"
      variant="outline"
      @click="handlePrimaryClick"
    >
      {{ primaryAction.label }}
    </t-button>
    <t-dropdown v-if="menuOptions.length > 0" :options="menuOptions" trigger="click" @click="handleMenuClick">
      <t-button size="small" theme="default" variant="outline">
        {{ moreLabel }}
      </t-button>
    </t-dropdown>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue';

type ActionOption = {
  disabled?: boolean;
  label: string;
  testId?: string;
  value: string;
};

const props = withDefaults(
  defineProps<{
    actions: ActionOption[];
    moreLabel?: string;
  }>(),
  {
    moreLabel: 'More',
  },
);

const emit = defineEmits<{
  action: [value: string];
}>();

const primaryAction = computed(() => props.actions[0] ?? null);
const menuOptions = computed(() =>
  props.actions.slice(1).map((action) => ({
    content: action.label,
    disabled: action.disabled,
    testId: action.testId,
    value: action.value,
  })),
);

function handlePrimaryClick() {
  if (!primaryAction.value?.disabled) {
    emit('action', primaryAction.value.value);
  }
}

function handleMenuClick(payload: { value?: string | number | Record<string, unknown> } | string | number) {
  const action = typeof payload === 'object' && payload ? payload.value : payload;
  if (typeof action === 'string') {
    emit('action', action);
  }
}
</script>
<style scoped lang="less">
.table-action-menu {
  align-items: center;
  display: inline-flex;
  gap: 8px;
  justify-content: center;
  width: 100%;
}

.table-action-menu :deep(.t-button) {
  min-width: 64px;
  white-space: nowrap;
}

.table-action-menu :deep(.t-dropdown) {
  flex: none;
}
</style>
