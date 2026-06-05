<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBookingStore } from '@/stores/booking'
import { loyaltyApi } from '@/api/loyalty'
import { useLocaleFormatter } from '@/composables/useLocaleFormatter'

const { t } = useI18n()
const { fmtTime } = useLocaleFormatter()

const showSOS    = ref(false)
const sosStart   = ref('')
const sosEnd     = ref('')
const sosNote    = ref('')
const sosSending = ref(false)

const store = useBookingStore()
const today = new Date().toISOString().split('T')[0]

const slotsLabel = computed(() =>
  t('booking.slots_available_label', { n: store.slots.length }),
)

function formatTime(iso: string) {
  return fmtTime(iso)
}

async function submitSOS() {
  if (!store.selectedMaster || !store.selectedServices[0]) return
  sosSending.value = true
  try {
    await loyaltyApi.createSOS({
      master_id:       store.selectedMaster.id,
      service_id:      store.selectedServices[0].id,
      preferred_start: new Date(sosStart.value).toISOString(),
      preferred_end:   new Date(sosEnd.value).toISOString(),
      client_note:     sosNote.value,
    })
    showSOS.value = false
    alert(t('booking.sos_alert_ok'))
  } catch {
    alert(t('booking.sos_alert_fail'))
  } finally {
    sosSending.value = false
  }
}

watch(() => store.selectedDate, () => {
  store.selectedSlot = null
  store.loadSlots()
})
</script>

<template>
  <div class="picker">
    <div class="picker__header">
      <div>
        <p class="t-label" style="margin-bottom: 4px;">{{ t('booking.select_date') }}</p>
        <p v-if="store.selectedMaster && store.selectedServices.length > 0"
           style="font-size:15px;font-weight:500;color:var(--text)">
          {{ store.selectedMaster.full_name }} ·
          {{ store.selectedServices.map(s => s.name).join(', ') }}
        </p>
      </div>
      <input
        v-model="store.selectedDate"
        type="date"
        :min="today"
        class="field-input"
        style="width: auto;"
      />
    </div>

    <div v-if="store.slotsLoading" class="picker__state">
      <div class="picker__spinner"></div>
      {{ t('booking.loading_slots') }}
    </div>
    <div v-else-if="!store.selectedDate" class="picker__state picker__state--hint">
      {{ t('booking.slots_hint_calendar') }}
    </div>

    <div v-else-if="store.slots.length === 0 && store.selectedDate" class="picker__empty">
      <p style="font-size:15px;font-weight:500;color:var(--text);margin-bottom:6px">
        {{ t('booking.no_slots') }}
      </p>
      <p style="font-size:13px;color:var(--text-2);margin-bottom:16px">
        {{ t('booking.try_other_day_sos') }}
      </p>
      <button class="btn btn-danger btn-sm" @click="showSOS = true">
        {{ t('booking.sos_short_button') }}
      </button>

      <div v-if="showSOS" class="sos-form card" style="margin-top:16px;text-align:left">
        <h4 style="font-size:15px;font-weight:600;margin-bottom:4px">{{ t('booking.sos_form_heading') }}</h4>
        <p style="font-size:13px;color:var(--text-2);margin-bottom:14px" v-html="t('booking.sos_schedule_note_html')" />
        <div class="field" style="margin-bottom:12px">
          <label class="field-label">{{ t('booking.sos_start') }}</label>
          <input v-model="sosStart" type="datetime-local" class="field-input" />
        </div>
        <div class="field" style="margin-bottom:14px">
          <label class="field-label">{{ t('booking.sos_end') }}</label>
          <input v-model="sosEnd" type="datetime-local" class="field-input" />
        </div>
        <div class="field" style="margin-bottom:14px">
          <label class="field-label">{{ t('booking.sos_note') }}</label>
          <input v-model="sosNote" type="text" class="field-input" :placeholder="t('booking.sos_placeholder_note')" />
        </div>
        <button class="btn btn-danger btn-full" :disabled="!sosStart || !sosEnd || sosSending" @click="submitSOS">
          <span v-if="sosSending" class="btn__spinner" />
          {{ t('booking.sos_send') }}
        </button>
      </div>
    </div>

    <div v-else>
      <p class="t-label" style="margin-bottom: 12px;">{{ slotsLabel }}</p>
      <div class="picker__grid">
        <button
          v-for="slot in store.slots"
          :key="slot.starts_at"
          class="slot-btn"
          :class="{ 'slot-btn--active': store.selectedSlot?.starts_at === slot.starts_at }"
          @click="store.selectedSlot = slot"
        >
          {{ formatTime(slot.starts_at) }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.picker__header {
  display: flex; justify-content: space-between;
  align-items: flex-end; gap: 16px;
  margin-bottom: 24px;
}
.picker__state {
  display: flex; align-items: center; gap: 12px;
  padding: 48px 24px; text-align: center; justify-content: center;
  font-size: 15px; color: var(--text-2);
  background: var(--bg-muted); border-radius: var(--radius-md);
}
.picker__state--hint { color: var(--text-3); }
.picker__spinner {
  width: 18px; height: 18px;
  border: 2px solid var(--border-strong);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  flex-shrink: 0;
}
@keyframes spin { to { transform: rotate(360deg); } }

.picker__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(90px, 1fr));
  gap: 8px;
}
.slot-btn {
  padding: 12px 8px; border-radius: var(--radius-md);
  border: 1.5px solid var(--border);
  background: var(--bg-card); color: var(--text-2);
  font-size: 15px; font-weight: 500;
  transition: all 0.15s; cursor: pointer;
}
.slot-btn:hover { border-color: var(--accent); color: var(--accent); background: var(--accent-bg); }
.slot-btn--active {
  border-color: var(--accent); background: var(--accent);
  color: #fff;
}
</style>
