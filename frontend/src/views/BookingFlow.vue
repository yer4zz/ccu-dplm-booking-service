<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter }               from 'vue-router'
import { useI18n }                 from 'vue-i18n'
import { useAuthStore }            from '@/stores/auth'
import { useBookingStore }         from '@/stores/booking'
import { loyaltyApi }             from '@/api/loyalty'
import VibeCheck                   from '@/components/booking/VibeCheck.vue'
import { ArrowRight } from 'lucide-vue-next'

const { t, locale }  = useI18n()
const router = useRouter()
const auth   = useAuthStore()
const store  = useBookingStore()

function svcName(svc: any): string {
  if (locale.value === 'kz' && svc.name_kz) return svc.name_kz
  if (locale.value === 'en' && svc.name_en) return svc.name_en
  return svc.name
}
function svcDesc(svc: any): string {
  if (locale.value === 'kz' && svc.desc_kz) return svc.desc_kz
  if (locale.value === 'en' && svc.desc_en) return svc.desc_en
  return svc.description || ''
}

const today = new Date().toISOString().split('T')[0]

onMounted(async () => {
  if (!auth.isLoggedIn) {
    router.push({ path: '/login', query: { redirect: '/book' } })
    return
  }
  await store.init()
})

const stepKeys   = ['service', 'master', 'slot', 'confirm'] as const
const stepNums   = ['01', '02', '03', '04']
const stepLabels = computed(() => [
  t('booking.step_service'),
  t('booking.step_master'),
  t('booking.step_time'),
  t('booking.step_confirm'),
])

const currentStepIndex = computed(() => stepKeys.indexOf(store.step as any))

const catColors: Record<string, string> = {
  hair:'#8B5CF6', nails:'#E8442A', face:'#F59E0B',
  combo:'#C9956A', beard:'#3B82F6', care:'#10B981',
}
function getCatColor(k: string) { return catColors[k] || '#8A8680' }
function getCatLabel(k: string) { return t(`category.${k}`) }

function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('ru-RU', { hour:'2-digit', minute:'2-digit' })
}
function totalDuration() {
  return store.selectedServices.reduce((s, svc) => s + svc.duration_min, 0)
}
function formatDateFull(iso: string) {
  return new Date(iso).toLocaleString('ru-RU', {
    weekday:'long', day:'numeric', month:'long',
    hour:'2-digit', minute:'2-digit',
  })
}

const showSOS    = ref(false)
const sosStart   = ref('')
const sosEnd     = ref('')
const sosNote    = ref('')
const sosSending = ref(false)

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
    alert(t('booking.sos_sent'))
  } catch {
    alert(t('booking.sos_alert_fail'))
  } finally {
    sosSending.value = false
  }
}
</script>

<template>
  <div class="booking-page">
    <div class="booking-container">

      <div v-if="store.step === 'done'" class="done">
        <div class="done__mark">
          <svg width="28" height="28" viewBox="0 0 28 28" fill="none">
            <path d="M6 14l6 6 10-10" stroke="var(--accent)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="section-label" style="justify-content:center">
          ЗАПИСЬ ОФОРМЛЕНА · Nº {{ new Date().getDate() }}
        </div>
        <h1 class="t-h1 done__title">
          {{ t('booking.success_title') }}<span class="red-dot">.</span>
        </h1>
        <p class="done__sub">{{ t('booking.success_sub') }}</p>
        <div class="done__btns">
          <button class="btn btn-primary" @click="router.push('/my')">
            {{ t('booking.my_bookings') }} <ArrowRight :size="16" />
          </button>
          <button class="btn btn-outline" @click="store.reset()">
            {{ t('booking.new_booking') }}
          </button>
        </div>
      </div>

      <template v-else>

        <div class="step-indicator">
          <div
            v-for="(label, i) in stepLabels"
            :key="i"
            class="step-item"
            :class="{
              'step-item--active': i === currentStepIndex,
              'step-item--done':   i < currentStepIndex,
            }"
          >
            <span class="step-item__num">{{ stepNums[i] }}</span>
            <span class="step-item__label">{{ label }}</span>
          </div>
        </div>
        <div class="step-line">
          <div
            class="step-line__fill"
            :style="`width: ${(currentStepIndex / (stepKeys.length - 1)) * 100}%`"
          />
        </div>

        <div class="step-header">
          <div class="section-label">
            {{ stepNums[currentStepIndex] }} / {{ stepLabels[currentStepIndex] }}
          </div>
          <h2 class="t-h2" v-if="store.step === 'service'">
            {{ t('booking.choose_services') }}<span class="red-dot">.</span>
          </h2>
          <h2 class="t-h2" v-else-if="store.step === 'master'">
            {{ t('booking.choose_master') }}<span class="red-dot">.</span>
          </h2>
          <h2 class="t-h2" v-else-if="store.step === 'slot'">
            {{ t('booking.choose_time') }}<span class="red-dot">.</span>
          </h2>
          <h2 class="t-h2" v-else-if="store.step === 'confirm'">
            {{ t('booking.confirm_heading') }}<span class="red-dot">.</span>
          </h2>
        </div>

        <div
          v-if="store.selectedServices.length > 0 && store.step !== 'service'"
          class="selected-pills"
        >
          <span
            v-for="svc in store.selectedServices"
            :key="svc.id"
            class="selected-pill"
          >{{ svcName(svc) }}</span>
        </div>

        <div v-if="store.step === 'service'">
          <p class="step-hint">{{ t('booking.hint_multi') }}</p>
          <div class="svc-list">
            <div
              v-for="(svc, i) in store.services"
              :key="svc.id"
              class="svc-row"
              :class="{ 'svc-row--active': store.selectedServices.some(s => s.id === svc.id) }"
              @click="store.toggleService(svc)"
            >
              <span class="svc-row__num">{{ String(i+1).padStart(2,'0') }}</span>
              <div class="svc-row__main">
                <h3 class="svc-row__name">{{ svcName(svc) }}</h3>
                <span class="svc-row__cat" :style="`color:${getCatColor(svc.category)}`">
                  {{ getCatLabel(svc.category) }}
                </span>
                <span class="svc-row__price-mobile">
                  {{ svc.duration_min }} {{ t('common.min') }} · {{ svc.price.toLocaleString() }} {{ t('common.currency') }}
                </span>
              </div>
              <p class="svc-row__desc">{{ svcDesc(svc) }}</p>
              <div class="svc-row__right">
                <span class="svc-row__dur">{{ svc.duration_min }} {{ t('common.min') }}</span>
                <span class="svc-row__price">{{ svc.price.toLocaleString() }} {{ t('common.currency') }}</span>
              </div>
              <div class="svc-row__check">
                <svg v-if="store.selectedServices.some(s => s.id === svc.id)" width="14" height="14" viewBox="0 0 14 14" fill="none">
                  <path d="M3 7l3 3 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="store.step === 'master'" class="masters-list">
          <div
            v-for="(m, i) in store.mastersForServices"
            :key="m.id"
            class="master-row-book"
            :class="{ 'master-row-book--active': store.selectedMaster?.id === m.id }"
            @click="store.selectedMaster = m"
          >
            <span class="master-row-book__num">{{ String(i+1).padStart(2,'0') }}</span>
            <div class="master-row-book__photo">
              <img v-if="m.avatar_url" :src="m.avatar_url" :alt="m.full_name" />
              <span v-else>{{ m.full_name[0] }}</span>
            </div>
            <div class="master-row-book__body">
              <h3 class="master-row-book__name">{{ m.full_name }}</h3>
              <p class="master-row-book__exp">{{ m.experience_years }} {{ t('home.masters_exp') }}</p>
              <p class="master-row-book__bio">{{ m.bio }}</p>
              <div class="master-row-book__tags">
                <span v-for="s in m.services.slice(0,3)" :key="s.id" class="master-row-book__tag">
                  {{ s.name }}
                </span>
              </div>
            </div>
            <div class="master-row-book__check" v-if="store.selectedMaster?.id === m.id">
              <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                <path d="M3 7l3 3 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
          </div>
        </div>

        <div v-else-if="store.step === 'slot'" class="slot-picker">
          <div class="slot-picker__header">
            <div class="field">
              <label class="field-label">{{ t('booking.select_date') }}</label>
              <input
                v-model="store.selectedDate"
                type="date"
                :min="today"
                class="field-input date-input"
                @change="store.selectedSlot = null; store.loadSlots()"
              />
            </div>
          </div>

          <div v-if="store.slotsLoading" class="slot-state">
            <div class="slot-spinner" />
            {{ t('booking.loading_slots') }}
          </div>
          <div v-else-if="!store.selectedDate" class="slot-state slot-state--hint">
            {{ t('booking.slots_hint_calendar') }}
          </div>
          <div v-else-if="store.slots.length === 0" class="slot-state">
            <div>
              <p style="font-size:15px;font-weight:500;margin-bottom:6px">{{ t('booking.no_slots') }}</p>
              <p style="font-size:13px;color:var(--text-3);margin-bottom:16px">
                {{ t('booking.try_other_day_sos') }}
              </p>
              <button
                class="btn btn-outline btn-sm"
                style="border-color:var(--accent);color:var(--accent)"
                @click="showSOS = true"
              >🆘 SOS (+30%)</button>
            </div>
          </div>
          <div v-else>
            <p class="t-label" style="margin-bottom:16px">
              {{ t('booking.slots_available_label', { n: store.slots.length }) }}
            </p>
            <div class="slots-grid">
              <button
                v-for="slot in store.slots"
                :key="slot.starts_at"
                class="slot-chip"
                :class="{ 'slot-chip--active': store.selectedSlot?.starts_at === slot.starts_at }"
                @click="store.selectedSlot = slot"
              >{{ formatTime(slot.starts_at) }}</button>
            </div>
          </div>

          <div v-if="showSOS" class="sos-form">
            <div class="section-label" style="margin-bottom:12px">SOS · СРОЧНАЯ ЗАПИСЬ</div>
            <p class="sos-form__desc" v-html="t('booking.sos_schedule_note_html')" />
            <div style="display:flex;flex-direction:column;gap:12px;margin-top:16px">
              <div class="field">
                <label class="field-label">{{ t('booking.sos_start') }}</label>
                <input v-model="sosStart" type="datetime-local" class="field-input" />
              </div>
              <div class="field">
                <label class="field-label">{{ t('booking.sos_end') }}</label>
                <input v-model="sosEnd" type="datetime-local" class="field-input" />
              </div>
              <div class="field">
                <label class="field-label">{{ t('booking.sos_note') }}</label>
                <input v-model="sosNote" type="text" class="field-input" :placeholder="t('booking.sos_placeholder_note')" />
              </div>
              <div style="display:flex;gap:10px;flex-wrap:wrap">
                <button class="btn btn-outline btn-sm" @click="showSOS = false">
                  {{ t('common.cancel') }}
                </button>
                <button
                  class="btn btn-primary btn-sm"
                  style="background:var(--accent);border-color:var(--accent)"
                  :disabled="!sosStart || !sosEnd || sosSending"
                  @click="submitSOS"
                >
                  <span v-if="sosSending" class="btn__spinner" />
                  {{ t('booking.sos_send') }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="store.step === 'confirm'" class="confirm-step">
          <div class="receipt">
            <div class="receipt__header">
              <span class="t-label">{{ t('booking.details_label') }}</span>
            </div>
            <div class="receipt__row">
              <span class="receipt__label">{{ t('booking.summary_master') }}</span>
              <span class="receipt__val">{{ store.selectedMaster?.full_name }}</span>
            </div>
            <div class="receipt__row">
              <span class="receipt__label">{{ t('booking.summary_services') }}</span>
              <div class="receipt__val receipt__val--tags">
                <span
                  v-for="svc in store.selectedServices"
                  :key="svc.id"
                  class="badge badge-accent"
                >{{ svcName(svc) }}</span>
              </div>
            </div>
            <div class="receipt__row">
              <span class="receipt__label">{{ t('booking.summary_datetime') }}</span>
              <span class="receipt__val" style="text-transform:capitalize">
                {{ formatDateFull(store.selectedSlot!.starts_at) }}
              </span>
            </div>
            <div class="receipt__row" style="border:none">
              <span class="receipt__label">{{ t('booking.summary_duration') }}</span>
              <span class="receipt__val">{{ totalDuration() }} {{ t('common.min') }}</span>
            </div>

            <template v-if="store.priceCalc">
              <div class="receipt__divider" />
              <div class="receipt__row receipt__row--sm">
                <span class="receipt__label">{{ t('booking.price_base') }}</span>
                <span class="receipt__val">{{ store.priceCalc.base_price.toLocaleString() }} ₸</span>
              </div>
              <div v-if="store.priceCalc.is_first_booking" class="receipt__row receipt__row--sm receipt__row--discount">
                <span class="receipt__label">
                  <span class="badge badge-accent" style="font-size:9px">{{ t('booking.price_first_visit') }}</span>
                  {{ t('booking.price_discount_30') }}
                </span>
                <span class="receipt__val">−{{ store.priceCalc.first_discount.toLocaleString() }} ₸</span>
              </div>
              <div v-if="store.priceCalc.multi_discount > 0" class="receipt__row receipt__row--sm receipt__row--discount">
                <span class="receipt__label">{{ t('booking.price_multi_badge', { n: store.priceCalc.service_count - 1 }) }}</span>
                <span class="receipt__val">−{{ store.priceCalc.multi_discount.toLocaleString() }} ₸</span>
              </div>

              <div v-if="store.priceCalc.points_available > 0" class="receipt__points">
                <div class="receipt__points-header">
                  <div>
                    <p style="font-size:13px;font-weight:500;color:var(--text)">{{ t('booking.price_loyalty') }}</p>
                    <p style="font-size:11px;color:var(--text-3);margin-top:2px">
                      {{ t('booking.price_available_points', {
                        n: store.priceCalc.points_available,
                        amt: (store.priceCalc.points_available * 10).toLocaleString() + ' ₸'
                      }) }}
                    </p>
                  </div>
                  <label style="display:flex;align-items:center;gap:6px;cursor:pointer">
                    <input
                      type="checkbox"
                      :checked="store.pointsToUse > 0"
                      @change="(e) => {
                        store.pointsToUse = (e.target as HTMLInputElement).checked
                          ? store.priceCalc!.points_available : 0
                        store.calcPrice()
                      }"
                    />
                    <span style="font-size:12px;font-weight:500;color:var(--accent)">
                      {{ t('booking.price_use_points') }}
                    </span>
                  </label>
                </div>
                <div v-if="store.pointsToUse > 0" class="receipt__row receipt__row--sm receipt__row--discount" style="margin-top:8px">
                  <span class="receipt__label">{{ t('booking.price_redeem', { n: store.pointsToUse }) }}</span>
                  <span class="receipt__val">−{{ store.priceCalc.points_discount.toLocaleString() }} ₸</span>
                </div>
              </div>

              <div class="receipt__divider" />
              <div class="receipt__row receipt__row--total">
                <span class="receipt__label-total">{{ t('booking.price_total_due') }}</span>
                <span class="receipt__price-total">{{ store.priceCalc.final_price.toLocaleString() }} ₸</span>
              </div>
              <div v-if="store.priceCalc.total_discount > 0" class="receipt__saved">
                {{ t('booking.price_you_save', { amt: store.priceCalc.total_discount.toLocaleString() + ' ₸' }) }}
              </div>
            </template>
          </div>

          <div class="field" style="margin-top:24px">
            <label class="field-label">{{ t('booking.wishes') }}</label>
            <textarea
              v-model="store.notes"
              class="field-input"
              :placeholder="t('booking.wishes_placeholder')"
              rows="3"
              style="resize:none"
            />
          </div>

          <div style="margin-top:24px">
            <VibeCheck
              v-model="store.vibeMode"
              :vibe-note="store.vibeNote"
              @update:vibeNote="store.vibeNote = $event"
            />
          </div>

          <div v-if="store.error" class="booking-error">{{ store.error }}</div>
        </div>

        <div class="step-nav">
          <button
            v-if="store.step !== 'service'"
            class="btn btn-outline"
            @click="store.prevStep()"
          >← {{ t('booking.back') }}</button>
          <div style="flex:1" />
          <button
            v-if="store.step !== 'confirm'"
            class="btn btn-primary"
            :disabled="!store.canProceed"
            @click="store.nextStep()"
          >{{ t('booking.next') }} <ArrowRight :size="16" /></button>
          <button
            v-else
            class="btn btn-primary"
            :disabled="store.loading"
            @click="store.confirm()"
          >
            <span v-if="store.loading" class="btn__spinner" />
            {{ t('booking.book') }} <ArrowRight :size="16" />
          </button>
        </div>

      </template>
    </div>
  </div>
</template>

<style scoped>
.booking-page { min-height: calc(100vh - 56px); background: var(--bg); }

.booking-container {
  max-width: 860px;
  margin: 0 auto;
  padding: 40px 40px 80px;
}

.done {
  padding: 60px 0;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}
.done__mark {
  width: 56px; height: 56px;
  border: 2px solid var(--accent);
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
}
.done__title { text-align: center; }
.done__sub   { font-size: 15px; color: var(--text-2); text-align: center; max-width: 400px; }
.done__btns  { display: flex; gap: 12px; flex-wrap: wrap; justify-content: center; }

.step-indicator {
  display: flex;
  overflow-x: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
  margin-bottom: 0;
  gap: 0;
}
.step-indicator::-webkit-scrollbar { display: none; }

.step-item {
  display: flex; align-items: center; gap: 6px;
  padding: 12px 16px 12px 0;
  flex-shrink: 0; transition: all 0.2s;
}
.step-item__num {
  font-family: var(--font-serif); font-size: 13px; font-weight: 700;
  color: var(--text-3); transition: color 0.2s;
}
.step-item__label {
  font-family: var(--font); font-size: 11px; font-weight: 500;
  letter-spacing: 0.08em; text-transform: uppercase;
  color: var(--text-3); transition: color 0.2s; white-space: nowrap;
}
.step-item--active .step-item__num   { color: var(--accent); }
.step-item--active .step-item__label { color: var(--text); }
.step-item--done   .step-item__num   { color: var(--text-2); }
.step-item--done   .step-item__label { color: var(--text-2); }

.step-line {
  height: 1px; background: var(--border);
  margin-bottom: 32px; position: relative;
}
.step-line__fill {
  position: absolute; top: 0; left: 0; height: 100%;
  background: var(--accent);
  transition: width 0.4s cubic-bezier(0.4,0,0.2,1);
}

.step-header { margin-bottom: 24px; }
.step-hint {
  font-size: 13px; color: var(--text-3); font-style: italic;
  margin-bottom: 20px; padding: 10px 14px;
  border-left: 2px solid var(--accent); background: var(--accent-bg);
}

.selected-pills {
  display: flex; flex-wrap: wrap; gap: 8px;
  margin-bottom: 20px; padding-bottom: 16px;
  border-bottom: 1px solid var(--border);
}
.selected-pill {
  font-size: 11px; font-weight: 500; letter-spacing: 0.06em;
  padding: 4px 12px;
  border: 1px solid var(--accent);
  border-radius: var(--radius-sm);
  color: var(--accent); background: var(--accent-bg);
}

.svc-list { display: flex; flex-direction: column; }
.svc-row {
  display: grid;
  grid-template-columns: 48px 1fr 1fr auto 32px;
  gap: 20px; align-items: center;
  padding: 18px 8px;
  border-top: 1px solid var(--border);
  cursor: pointer; transition: all 0.18s ease;
}
.svc-row:last-child { border-bottom: 1px solid var(--border); }
.svc-row:hover      { background: var(--bg-hover); padding-left: 14px; }
.svc-row:hover .svc-row__name { transform: translateX(4px); }
.svc-row--active {
  border-left: 2px solid var(--accent);
  background: var(--accent-bg); padding-left: 14px;
}
.svc-row__num {
  font-family: var(--font-serif); font-size: 24px; font-weight: 700;
  color: var(--text); opacity: 0.1; line-height: 1;
}
.svc-row--active .svc-row__num { opacity: 0.4; }
.svc-row__main { display: flex; flex-direction: column; gap: 3px; }
.svc-row__name {
  font-family: var(--font-serif); font-size: 19px; font-weight: 700;
  color: var(--text); line-height: 1.2; transition: transform 0.2s;
}
.svc-row__cat  { font-size: 10px; font-weight: 500; letter-spacing: 0.08em; text-transform: uppercase; }
.svc-row__desc { font-size: 13px; color: var(--text-3); line-height: 1.5; }
.svc-row__right {
  display: flex; flex-direction: column; align-items: flex-end; gap: 3px;
}
.svc-row__dur   { font-size: 11px; color: var(--text-3); }
.svc-row__price {
  font-family: var(--font-serif); font-size: 18px; font-weight: 700; color: var(--text);
}
.svc-row__price-mobile { display: none; }
.svc-row__check {
  width: 24px; height: 24px;
  border: 1.5px solid var(--border-strong);
  border-radius: var(--radius-sm);
  display: flex; align-items: center; justify-content: center;
  color: var(--accent); transition: all 0.15s;
}
.svc-row--active .svc-row__check {
  border-color: var(--accent); background: var(--accent); color: #fff;
}

.masters-list { display: flex; flex-direction: column; }
.master-row-book {
  display: grid;
  grid-template-columns: 40px 64px 1fr 32px;
  gap: 20px; align-items: center;
  padding: 20px 8px;
  border-top: 1px solid var(--border);
  cursor: pointer; transition: all 0.18s ease;
}
.master-row-book:last-child { border-bottom: 1px solid var(--border); }
.master-row-book:hover      { transform: translateX(4px); }
.master-row-book--active {
  border-left: 2px solid var(--accent);
  background: var(--accent-bg); padding-left: 14px;
}
.master-row-book__num {
  font-family: var(--font-serif); font-size: 20px; font-weight: 700; opacity: 0.1;
}
.master-row-book__photo {
  width: 56px; height: 56px; border-radius: 50%; overflow: hidden;
  background: var(--bg-muted); border: 1px solid var(--border);
  display: flex; align-items: center; justify-content: center;
  font-family: var(--font-serif); font-size: 20px; font-weight: 700;
  color: var(--accent); flex-shrink: 0;
}
.master-row-book__photo img { width: 100%; height: 100%; object-fit: cover; }
.master-row-book__name {
  font-family: var(--font-serif); font-size: 20px; font-weight: 700;
  color: var(--text); margin-bottom: 2px;
}
.master-row-book__exp {
  font-size: 10px; font-weight: 500; letter-spacing: 0.1em;
  text-transform: uppercase; color: var(--accent); margin-bottom: 5px;
}
.master-row-book__bio {
  font-size: 13px; color: var(--text-3); line-height: 1.5; margin-bottom: 8px;
  display: -webkit-box; -webkit-line-clamp: 2; line-clamp: 2;
  -webkit-box-orient: vertical; overflow: hidden;
}
.master-row-book__tags { display: flex; flex-wrap: wrap; gap: 5px; }
.master-row-book__tag {
  font-size: 10px; padding: 2px 8px;
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-sm); color: var(--text-3);
}
.master-row-book__check {
  width: 24px; height: 24px;
  border: 1.5px solid var(--accent); background: var(--accent);
  border-radius: var(--radius-sm);
  display: flex; align-items: center; justify-content: center; color: #fff;
}

.slot-picker__header { margin-bottom: 20px; }
.date-input { width: auto; min-width: 200px; }
.slot-state {
  display: flex; align-items: center; justify-content: center;
  gap: 12px; padding: 40px 20px;
  background: var(--bg-muted); border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 14px; color: var(--text-2); text-align: center;
}
.slot-state--hint { color: var(--text-3); }
.slot-spinner {
  width: 16px; height: 16px;
  border: 1.5px solid var(--border-strong);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.6s linear infinite; flex-shrink: 0;
}
@keyframes spin { to { transform: rotate(360deg); } }

.sos-form {
  margin-top: 24px; padding: 20px;
  border: 1px solid var(--accent);
  border-radius: var(--radius-sm); background: var(--accent-bg);
}
.sos-form__desc { font-size: 14px; color: var(--text-2); line-height: 1.6; }

.receipt { border: 1px solid var(--border); border-radius: var(--radius-sm); overflow: hidden; }
.receipt__header {
  padding: 14px 20px; background: var(--bg-muted); border-bottom: 1px solid var(--border);
}
.receipt__row {
  display: flex; justify-content: space-between; align-items: center;
  gap: 12px; padding: 14px 20px; border-bottom: 1px solid var(--border);
}
.receipt__row--sm      { padding: 10px 20px; }
.receipt__row--discount { color: var(--success); }
.receipt__row--discount .receipt__val { font-weight: 600; }
.receipt__row--total   { padding: 18px 20px; background: var(--bg-muted); border-bottom: none; }
.receipt__label        { font-size: 13px; color: var(--text-2); flex-shrink: 0; max-width: 50%; }
.receipt__val          { font-size: 14px; font-weight: 500; color: var(--text); text-align: right; }
.receipt__val--tags    { display: flex; flex-wrap: wrap; gap: 6px; justify-content: flex-end; }
.receipt__label-total  { font-size: 15px; font-weight: 600; color: var(--text); }
.receipt__price-total  {
  font-family: var(--font-serif); font-size: 28px; font-weight: 700; color: var(--text);
}
.receipt__divider { height: 1px; background: var(--border); }
.receipt__saved {
  padding: 10px 20px; font-size: 12px; font-weight: 500;
  letter-spacing: 0.04em; color: var(--success);
  background: var(--success-bg); text-align: center;
}
.receipt__points {
  padding: 14px 20px; border-bottom: 1px solid var(--border); background: var(--accent-bg);
}
.receipt__points-header {
  display: flex; justify-content: space-between; align-items: center; gap: 12px;
}

.step-nav {
  display: flex; align-items: center; gap: 12px;
  padding-top: 24px; border-top: 1px solid var(--border); margin-top: 28px;
}
.booking-error {
  margin-top: 16px; padding: 12px 16px;
  background: var(--danger-bg); color: var(--danger);
  border-radius: var(--radius-sm); font-size: 14px;
  border-left: 2px solid var(--danger);
}


@media (max-width: 600px) {
  .booking-container { padding: 20px 16px 60px; }

  .step-item       { padding: 10px 12px 10px 0; gap: 5px; }
  .step-item__label { display: none; }
  .step-item__num  { font-size: 15px; }

  .step-line       { margin-bottom: 20px; }
  .step-header     { margin-bottom: 16px; }
  .step-header .t-h2 { font-size: clamp(22px, 7vw, 30px); }

  .svc-row {
    grid-template-columns: 32px 1fr 28px;
    gap: 12px; padding: 14px 4px;
  }
  .svc-row__desc          { display: none; }
  .svc-row__right         { display: none; }
  .svc-row__price-mobile  {
    display: block;
    font-size: 11px; color: var(--text-3);
    margin-top: 2px; letter-spacing: 0.02em;
  }
  .svc-row__num  { font-size: 18px; }
  .svc-row__name { font-size: 16px; }

  .master-row-book {
    grid-template-columns: 52px 1fr 28px;
    gap: 12px; padding: 16px 4px;
  }
  .master-row-book__num { display: none; }
  .master-row-book__photo { width: 48px; height: 48px; font-size: 18px; }
  .master-row-book__name  { font-size: 17px; }
  .master-row-book__tags  { display: none; }

  .date-input { width: 100%; min-width: unset; }

  .slots-grid { grid-template-columns: repeat(auto-fill, minmax(72px, 1fr)); gap: 10px; }
  .slot-chip  { height: 44px; font-size: 15px; }

  .receipt__row  { padding: 12px 16px; gap: 8px; }
  .receipt__row--sm { padding: 8px 16px; }
  .receipt__row--total { padding: 14px 16px; }
  .receipt__header    { padding: 12px 16px; }
  .receipt__label     { font-size: 12px; max-width: 45%; }
  .receipt__val       { font-size: 13px; }
  .receipt__price-total { font-size: 22px; }
  .receipt__points    { padding: 12px 16px; }
  .receipt__points-header { flex-direction: column; align-items: flex-start; gap: 10px; }

  .step-nav     { flex-wrap: wrap; gap: 10px; }
  .step-nav .btn { flex: 1; min-width: 120px; justify-content: center; }

  .done        { padding: 40px 0; }
  .done__btns  { flex-direction: column; width: 100%; }
  .done__btns .btn { width: 100%; justify-content: center; }
}

@media (max-width: 380px) {
  .booking-container { padding: 16px 12px 60px; }
  .svc-row__name  { font-size: 15px; }
  .slot-chip      { font-size: 14px; }
  .step-item__num { font-size: 13px; }
}
</style>