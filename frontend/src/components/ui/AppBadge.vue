<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
defineProps<{ status: string }>()
</script>

<template>
  <span class="badge" :class="`badge-${status}`">
    {{ t(`status.${status}`, status) }}
  </span>
</template>