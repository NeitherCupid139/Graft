<template>
  <div class="preset-grid">
    <button
      v-for="preset in presets"
      :key="preset.id"
      type="button"
      class="preset-card"
      :class="{ 'preset-card--active': activePresetId === preset.id }"
      @click="$emit('select', preset.id)"
    >
      <span class="preset-card__swatch" :style="{ background: preset.brandTheme }" />
      <span class="preset-card__title">{{ preset.label }}</span>
      <span class="preset-card__desc">{{ preset.description }}</span>
    </button>
  </div>
</template>
<script setup lang="ts">
import type { ThemePresetDefinition } from '@/types/theme';

defineProps<{
  presets: ThemePresetDefinition[];
  activePresetId: string | null;
}>();

defineEmits<{
  select: [presetId: string];
}>();
</script>
