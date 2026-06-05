<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

defineProps<{
  modelValue: string
  vibeNote:   string
}>()
const emit = defineEmits<{
  'update:modelValue': [v: string]
  'update:vibeNote':  [v: string]
}>()

const { t } = useI18n()

const vibes = computed(() => [
  {
    key:   'beauty_school',
    emoji: '🎓',
    title: t('vibe.beauty_school_title'),
    desc:  t('vibe.beauty_school_desc'),
  },
  {
    key:   'nap_time',
    emoji: '😴',
    title: t('vibe.nap_time_title'),
    desc:  t('vibe.nap_time_desc'),
  },
  {
    key:   'insta_vibe',
    emoji: '📸',
    title: t('vibe.insta_vibe_title'),
    desc:  t('vibe.insta_vibe_desc'),
  },
  {
    key:   'turbo',
    emoji: '⚡',
    title: t('vibe.turbo_title'),
    desc:  t('vibe.turbo_desc'),
  },
])
</script>

<template>
  <div class="vibe">
    <div class="vibe__header">
      <h4 class="vibe__title">{{ t('booking.vibe_title') }}</h4>
      <p class="vibe__sub">{{ t('booking.vibe_sub') }}</p>
    </div>

    <div class="vibe__grid">
      <button
        v-for="v in vibes"
        :key="v.key"
        class="vibe__card"
        :class="{ 'vibe__card--active': modelValue === v.key }"
        @click="emit('update:modelValue', modelValue === v.key ? '' : v.key)"
        type="button"
      >
        <span class="vibe__emoji">{{ v.emoji }}</span>
        <span class="vibe__name">{{ v.title }}</span>
        <span class="vibe__desc">{{ v.desc }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.vibe__header   { margin-bottom: 14px; }
.vibe__title    { font-size: 16px; font-weight: 600; color: var(--text); margin-bottom: 3px; }
.vibe__sub      { font-size: 13px; color: var(--text-2); }

.vibe__grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}
.vibe__card {
  display: flex; flex-direction: column; align-items: flex-start;
  gap: 4px; padding: 14px;
  border: 1.5px solid var(--border);
  border-radius: var(--radius-md);
  background: var(--bg-card);
  cursor: pointer; transition: all 0.15s;
  text-align: left;
}
.vibe__card:hover       { border-color: var(--accent); background: var(--accent-bg); }
.vibe__card--active     {
  border-color: var(--accent);
  background: var(--accent-bg);
  box-shadow: 0 0 0 3px var(--accent-bg);
}
.vibe__emoji { font-size: 22px; margin-bottom: 2px; }
.vibe__name  { font-size: 13px; font-weight: 600; color: var(--text); }
.vibe__desc  { font-size: 12px; color: var(--text-2); line-height: 1.4; }
</style>
