<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useLocaleFormatter } from '@/composables/useLocaleFormatter'
import type { BookingPriceCalc } from '@/types'

defineProps<{
  calc: BookingPriceCalc
  pointsAvailable: number
  pointsUsed: number
}>()
defineEmits<{ 'update:pointsUsed': [v: number] }>()

const { t } = useI18n()
const { dateLocaleTag } = useLocaleFormatter()

function fmt(n: number) {
  return n.toLocaleString(dateLocaleTag.value) + ' ' + t('common.currency')
}
</script>

<template>
  <div class="breakdown card">
    <h4 class="breakdown__title">{{ t('booking.price_title') }}</h4>

    <div class="breakdown__rows">
      <div class="breakdown__row">
        <span>{{ t('booking.price_base') }}</span>
        <span>{{ fmt(calc.base_price) }}</span>
      </div>

      <div v-if="calc.is_first_booking" class="breakdown__row breakdown__row--discount">
        <span>
          <span class="badge badge-accent">{{ t('booking.price_first_visit') }}</span>
          {{ t('booking.price_discount_30') }}
        </span>
        <span>−{{ fmt(calc.first_discount) }}</span>
      </div>

      <div v-if="calc.multi_discount > 0" class="breakdown__row breakdown__row--discount">
        <span>
          <span class="badge badge-accent">{{ t('booking.price_multi_badge', { n: calc.service_count - 1 }) }}</span>
          {{ t('booking.price_multi_discount_word', { pct: (calc.service_count - 1) * 5 }) }}
        </span>
        <span>−{{ fmt(calc.multi_discount) }}</span>
      </div>
    </div>

    <div v-if="calc.points_available > 0" class="breakdown__points">
      <div class="breakdown__points-header">
        <div>
          <p class="breakdown__points-label">{{ t('booking.price_loyalty') }}</p>
          <p class="breakdown__points-sub">
            {{ t('booking.price_available_points', {
              n: calc.points_available,
              amt: fmt(calc.points_available * 10),
            }) }}
          </p>
        </div>
        <label class="breakdown__points-toggle">
          <input
            type="checkbox"
            :checked="pointsUsed > 0"
            @change="(e) => $emit('update:pointsUsed',
              (e.target as HTMLInputElement).checked ? calc.points_available : 0)"
          />
          <span>{{ t('booking.price_use_points') }}</span>
        </label>
      </div>
      <div v-if="pointsUsed > 0" class="breakdown__row breakdown__row--discount" style="margin-top:8px">
        <span>{{ t('booking.price_redeem', { n: pointsUsed }) }}</span>
        <span>−{{ fmt(calc.points_discount) }}</span>
      </div>
    </div>

    <div class="breakdown__divider" />

    <div class="breakdown__total">
      <span>{{ t('booking.price_total_due') }}</span>
      <span class="breakdown__total-price">{{ fmt(calc.final_price) }}</span>
    </div>

    <div v-if="calc.total_discount > 0" class="breakdown__saved">
      {{ t('booking.price_you_save', { amt: fmt(calc.total_discount) }) }}
    </div>
  </div>
</template>

<style scoped>
.breakdown { padding: 20px; }
.breakdown__title { font-size: 15px; font-weight: 600; margin-bottom: 16px; color: var(--text); }

.breakdown__rows  { display: flex; flex-direction: column; gap: 10px; margin-bottom: 16px; }
.breakdown__row   {
  display: flex; justify-content: space-between; align-items: center;
  gap: 8px; font-size: 14px; color: var(--text-2);
}
.breakdown__row--discount { color: var(--success); }
.breakdown__row--discount span:last-child { font-weight: 600; }

.breakdown__points {
  background: var(--accent-bg);
  border: 1px solid var(--accent);
  border-radius: var(--radius-md);
  padding: 14px; margin-bottom: 16px;
}
.breakdown__points-header {
  display: flex; justify-content: space-between;
  align-items: flex-start; gap: 12px;
}
.breakdown__points-label { font-size: 14px; font-weight: 500; color: var(--text); }
.breakdown__points-sub   { font-size: 12px; color: var(--text-2); margin-top: 2px; }
.breakdown__points-toggle {
  display: flex; align-items: center; gap: 6px;
  font-size: 13px; font-weight: 500; color: var(--accent);
  cursor: pointer; flex-shrink: 0;
}

.breakdown__divider { height: 1px; background: var(--border); margin: 12px 0; }
.breakdown__total   {
  display: flex; justify-content: space-between; align-items: center;
  font-size: 15px; font-weight: 600; color: var(--text);
}
.breakdown__total-price { font-size: 22px; font-weight: 700; color: var(--text); }
.breakdown__saved {
  margin-top: 10px; text-align: center;
  font-size: 13px; font-weight: 500;
  color: var(--success);
  background: var(--success-bg);
  padding: 6px 12px; border-radius: var(--radius-full);
}
</style>
