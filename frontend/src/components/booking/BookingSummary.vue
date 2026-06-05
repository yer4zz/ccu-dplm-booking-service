<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useBookingStore } from '@/stores/booking'
import { useLocaleFormatter } from '@/composables/useLocaleFormatter'
import PriceBreakdown from './PriceBreakdown.vue'

const { t } = useI18n()
const { fmtDateTime } = useLocaleFormatter()
const store = useBookingStore()

function formatDate(iso: string) {
  return fmtDateTime(iso, {
    weekday: 'long', day: 'numeric', month: 'long',
    hour: '2-digit', minute: '2-digit',
  })
}

function totalDuration() {
  return store.selectedServices.reduce((s, svc) => s + svc.duration_min, 0)
}
</script>

<template>
  <div class="summary">
    <div class="summary__card card">
      <div class="summary__row">
        <span class="summary__lbl">{{ t('booking.summary_master') }}</span>
        <span class="summary__val">{{ store.selectedMaster?.full_name }}</span>
      </div>
      <div class="summary__row">
        <span class="summary__lbl">{{ t('booking.summary_services') }}</span>
        <div class="summary__services">
          <span
            v-for="svc in store.selectedServices"
            :key="svc.id"
            class="badge badge-accent"
          >{{ svc.name }}</span>
        </div>
      </div>
      <div class="summary__row">
        <span class="summary__lbl">{{ t('booking.summary_datetime') }}</span>
        <span class="summary__val" style="text-transform: capitalize;">
          {{ formatDate(store.selectedSlot!.starts_at) }}
        </span>
      </div>
      <div class="summary__row" style="border: none;">
        <span class="summary__lbl">{{ t('booking.summary_duration') }}</span>
        <span class="summary__val">{{ totalDuration() }} {{ t('common.min') }}</span>
      </div>
    </div>

    <div v-if="store.priceCalc" style="margin-top: 16px;">
      <PriceBreakdown
        :calc="store.priceCalc"
        :points-available="store.priceCalc.points_available"
        :points-used="store.pointsToUse"
        @update:points-used="(v) => { store.pointsToUse = v; store.calcPrice() }"
      />
    </div>

    <div class="field" style="margin-top: 16px;">
      <label class="field-label">{{ t('booking.wishes') }}</label>
      <textarea
        v-model="store.notes"
        class="field-input"
        :placeholder="t('booking.wishes_placeholder')"
        rows="3"
        style="resize: none;"
      />
    </div>

    <div v-if="store.error" class="summary__error">{{ store.error }}</div>
  </div>
</template>

<style scoped>
.summary__card { padding: 0; overflow: hidden; }
.summary__row  {
  display: flex; justify-content: space-between;
  align-items: flex-start; gap: 16px;
  padding: 14px 20px; border-bottom: 1px solid var(--border);
}
.summary__lbl  { font-size: 13px; color: var(--text-2); flex-shrink: 0; padding-top: 2px; }
.summary__val  { font-size: 14px; font-weight: 500; color: var(--text); text-align: right; }
.summary__services { display: flex; flex-wrap: wrap; gap: 6px; justify-content: flex-end; }
.summary__error {
  margin-top: 12px; padding: 12px 16px;
  background: var(--danger-bg); color: var(--danger);
  border-radius: var(--radius-md); font-size: 14px;
}
</style>
