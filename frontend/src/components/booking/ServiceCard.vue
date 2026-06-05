<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Service } from '@/types'

defineProps<{
  service:  Service
  selected: boolean
  index?:   number
}>()
defineEmits<{ toggle: [service: Service] }>()

const { t } = useI18n()

const catColor: Record<string, string> = {
  hair:  '#8B5CF6',
  nails: '#EC4899',
  face:  '#F59E0B',
  combo: '#C9956A',
  beard: '#3B82F6',
  care:  '#10B981',
}

function getCatLabel(kategoriya: string): string {
  try {
    const perevod = t(`category.${kategoriya}`)
    return perevod
  } catch {
    return kategoriya
  }
}

function getCatColor(kategoriya: string): string {
  return catColor[kategoriya] || '#C9956A'
}
</script>

<template>
  <div
    class="svc card card-hover"
    :class="{ 'card-active': selected }"
    @click="$emit('toggle', service)"
  >
    <div class="svc__top">
      <span
        class="svc__cat"
        :style="`background:${getCatColor(service.category)}20; color:${getCatColor(service.category)}`"
      >
        {{ getCatLabel(service.category) }}
      </span>

      <div class="svc__check" :class="{ 'svc__check--active': selected }">
        <svg v-if="selected" width="14" height="14" viewBox="0 0 14 14" fill="none">
          <path d="M3 7l3 3 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </div>
    </div>

    <h3 class="svc__name">{{ service.name }}</h3>
    <p class="svc__desc">{{ service.description }}</p>

    <div class="svc__footer">
      <span class="svc__dur">
        {{ service.duration_min }} {{ t('common.min') }}
      </span>
      <span class="svc__price">
        {{ service.price.toLocaleString() }} {{ t('common.currency') }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.svc { display: flex; flex-direction: column; gap: 10px; user-select: none; }

.svc__top {
  display: flex; justify-content: space-between; align-items: center;
}
.svc__cat {
  font-size: 11px; font-weight: 600; letter-spacing: 0.04em;
  padding: 3px 10px; border-radius: var(--radius-full);
}
.svc__check {
  width: 22px; height: 22px; border-radius: 50%;
  border: 2px solid var(--border-strong);
  display: flex; align-items: center; justify-content: center;
  transition: all 0.15s; flex-shrink: 0;
  color: #fff;
}
.svc__check--active {
  background: var(--accent); border-color: var(--accent);
}

.svc__name  { font-size: 16px; font-weight: 600; color: var(--text); line-height: 1.3; }
.svc__desc  { font-size: 13px; color: var(--text-2); line-height: 1.6; flex: 1; }
.svc__footer {
  display: flex; justify-content: space-between; align-items: center;
  padding-top: 12px; border-top: 1px solid var(--border); margin-top: auto;
}
.svc__dur   { font-size: 13px; color: var(--text-3); }
.svc__price { font-size: 17px; font-weight: 700; color: var(--text); }
</style>