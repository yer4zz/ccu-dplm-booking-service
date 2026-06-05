<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Master } from '@/types'

defineProps<{ master: Master; selected: boolean }>()
defineEmits<{ select: [master: Master] }>()

const { t } = useI18n()
</script>

<template>
  <div
    class="mc card card-hover"
    :class="{ 'card-active': selected }"
    @click="$emit('select', master)"
  >
    <div class="mc__left">
      <div class="mc__avatar">
        <img v-if="master.avatar_url" :src="master.avatar_url" :alt="master.full_name" />
        <span v-else>{{ master.full_name[0] }}</span>
      </div>
    </div>
    <div class="mc__body">
      <div class="mc__top">
        <div>
          <h3 class="mc__name">{{ master.full_name }}</h3>
          <p class="mc__exp">{{ master.experience_years }} {{ t('common.experience') }}</p>
        </div>
        <div v-if="selected" class="mc__check">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <circle cx="10" cy="10" r="10" fill="var(--accent)"/>
            <path d="M6 10l3 3 5-5" stroke="white" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </div>
      </div>
      <p class="mc__bio">{{ master.bio }}</p>
      <div class="mc__tags">
        <span v-for="s in master.services.slice(0,3)" :key="s.id" class="mc__tag">
          {{ s.name }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mc { display: flex; gap: 16px; align-items: flex-start; }
.mc__avatar {
  width: 56px; height: 56px; border-radius: 50%;
  background: var(--accent-bg); border: 2px solid var(--border);
  display: flex; align-items: center; justify-content: center;
  font-size: 22px; font-weight: 700; color: var(--accent);
  flex-shrink: 0; overflow: hidden;
}
.mc__avatar img { width: 100%; height: 100%; object-fit: cover; }
.mc__body  { flex: 1; min-width: 0; }
.mc__top   { display: flex; justify-content: space-between; align-items: flex-start; gap: 8px; margin-bottom: 6px; }
.mc__name  { font-size: 16px; font-weight: 600; color: var(--text); }
.mc__exp   { font-size: 13px; color: var(--accent); font-weight: 500; }
.mc__bio   { font-size: 14px; color: var(--text-2); line-height: 1.5; margin-bottom: 10px; }
.mc__tags  { display: flex; flex-wrap: wrap; gap: 6px; }
.mc__tag   {
  font-size: 12px; padding: 3px 10px;
  background: var(--bg-muted); color: var(--text-2);
  border-radius: var(--radius-full); border: 1px solid var(--border);
}
</style>
