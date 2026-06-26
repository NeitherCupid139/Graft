<template>
  <span :class="rootCls" aria-hidden="true">
    <svg class="graft-brand-logo__mark" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M16 5C9.925 5 5 9.925 5 16C5 22.075 9.925 27 16 27C19.87 27 23.3 24.87 25.1 21.5H19.5V18.5H26.5C26.5 11.5 21.5 5 16 5ZM16 9C12.134 9 9 12.134 9 16C9 19.866 12.134 23 16 23C18.21 23 20.21 21.88 21.5 20.2V17H16V13H21.5C21.5 10.5 19 9 16 9Z"
        fill="currentColor"
      />
    </svg>
    <span v-if="variant === 'wordmark'" class="graft-brand-logo__wordmark">Graft</span>
  </span>
</template>
<script setup lang="ts">
import { computed } from 'vue';

const { variant = 'mark' } = defineProps<{
  variant?: 'mark' | 'wordmark';
}>();

const rootCls = computed(() => ['graft-brand-logo', `graft-brand-logo--${variant}`]);
</script>
<style lang="less" scoped>
.graft-brand-logo {
  align-items: center;
  display: inline-flex;
  gap: var(--td-comp-margin-s);
  line-height: 1;
  max-width: 100%;
  min-width: 0;
  transition: color 0.2s ease;

  &__mark {
    display: block;
    flex: 0 0 auto;
    height: 1.625rem;
    width: 1.625rem;
  }

  &--mark &__mark {
    height: var(--td-comp-size-s);
    width: var(--td-comp-size-s);
  }

  &__wordmark {
    color: var(--td-text-color-primary);
    font: var(--td-font-title-medium);
    font-weight: 600;
    letter-spacing: 0.02em;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &--wordmark &__mark {
    color: var(--td-brand-color);
  }
}
</style>