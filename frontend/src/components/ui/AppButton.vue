<script setup lang="ts">
defineProps<{
  variant?: 'primary' | 'accent' | 'outline' | 'ghost' | 'danger'
  size?: 'sm' | 'lg'
  disabled?: boolean
  loading?: boolean
  full?: boolean
}>()
</script>

<template>
  <button
    class="btn"
    :class="[
      `btn-${variant ?? 'primary'}`,
      size ? `btn-${size}` : '',
      full ? 'btn-full' : '',
    ]"
    :disabled="disabled || loading"
  >
    <span v-if="loading" class="btn__spinner" />
    <slot />
  </button>
</template>