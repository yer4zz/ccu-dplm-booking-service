<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter }    from 'vue-router'
import { useI18n }      from 'vue-i18n'
import { bookingsApi }  from '@/api/bookings'
import { loyaltyApi }   from '@/api/loyalty'
import { mastersApi }   from '@/api/masters'
import { api }          from '@/api'
import AppBadge         from '@/components/ui/AppBadge.vue'
import type { Booking, WaitlistEntry, Master } from '@/types'

const router   = useRouter()
const { t }    = useI18n()
const zagruzka = ref(true)

const bookings        = ref<Booking[]>([])
const waitlist        = ref<WaitlistEntry[]>([])
const autoReschedules = ref<any[]>([])
const cancelling      = ref<string | null>(null)

const showReschedule   = ref(false)
const rescheduleTarget = ref<Booking | null>(null)
const newDate          = ref('')
const newSlotStartsAt  = ref('')
const newTime          = ref('')
const newMasterID      = ref('')
const rescheduling     = ref(false)
const rescheduleError  = ref('')
const rescheduleSlots  = ref<any[]>([])
const allMasters       = ref<Master[]>([])
const today            = new Date().toISOString().split('T')[0]

const showReview    = ref(false)
const reviewTarget  = ref<Booking | null>(null)
const reviewRating  = ref(5)
const reviewComment = ref('')
const reviewSending = ref(false)
const hoverRating   = ref(0)

onMounted(async () => {
  const [b, w, ar] = await Promise.all([
    bookingsApi.my(),
    loyaltyApi.getWaitlist().catch(() => []),
    bookingsApi.getAutoReschedules().catch(() => []),
  ])
  bookings.value        = b ?? []
  waitlist.value        = w ?? []
  autoReschedules.value = ar ?? []
  zagruzka.value        = false
})

const upcoming = computed(() =>
  bookings.value.filter(b => !['cancelled','completed','no_show'].includes(b.status))
)
const past = computed(() =>
  bookings.value.filter(b => ['cancelled','completed','no_show'].includes(b.status))
)

async function cancel(id: string) {
  cancelling.value = id
  await bookingsApi.cancel(id)
  bookings.value = bookings.value.map(b =>
    b.id === id ? { ...b, status: 'cancelled' as const } : b
  )
  cancelling.value = null
}

async function leaveWaitlist(id: string) {
  await loyaltyApi.removeFromWaitlist(id)
  waitlist.value = waitlist.value.filter(w => w.id !== id)
}

async function cancelAutoReschedule(id: string) {
  await bookingsApi.declineAutoReschedule(id)
  autoReschedules.value = autoReschedules.value.filter((r: any) => r.id !== id)
  bookings.value = await bookingsApi.my()
}

async function openReschedule(booking: Booking) {
  rescheduleTarget.value = booking
  newMasterID.value      = booking.master_id
  newDate.value          = ''
  newTime.value          = ''
  newSlotStartsAt.value  = ''
  rescheduleError.value  = ''
  rescheduleSlots.value  = []
  if (allMasters.value.length === 0) {
    allMasters.value = await mastersApi.list()
  }
  showReschedule.value = true
}

async function loadRescheduleSlots() {
  if (!newDate.value || !rescheduleTarget.value) return
  rescheduleSlots.value = await mastersApi.getSlots(
    newMasterID.value || rescheduleTarget.value.master_id,
    rescheduleTarget.value.service_id,
    newDate.value,
  )
}

async function confirmReschedule() {
  if (!rescheduleTarget.value || !newDate.value || !newSlotStartsAt.value) return
  rescheduling.value    = true
  rescheduleError.value = ''
  try {
    const payload: any = { new_starts_at: newSlotStartsAt.value }
    if (newMasterID.value && newMasterID.value !== rescheduleTarget.value.master_id) {
      payload.new_master_id = newMasterID.value
    }
    await api.patch(`/bookings/${rescheduleTarget.value.id}/reschedule`, payload)
    bookings.value       = await bookingsApi.my()
    showReschedule.value = false
  } catch (e: any) {
    rescheduleError.value = e.response?.data?.error === 'slot_unavailable'
      ? t('myBookings.slot_busy')
      : t('myBookings.generic_error')
  } finally {
    rescheduling.value = false
  }
}

function openReview(b: Booking) {
  reviewTarget.value  = b
  reviewRating.value  = 5
  reviewComment.value = ''
  showReview.value    = true
}

async function submitReview() {
  if (!reviewTarget.value) return
  reviewSending.value = true
  try {
    await bookingsApi.createReview(reviewTarget.value.id, reviewRating.value, reviewComment.value)
    bookings.value = bookings.value.map(b =>
      b.id === reviewTarget.value!.id ? { ...b, has_review: true } as any : b
    )
    showReview.value = false
  } finally {
    reviewSending.value = false
  }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('ru-RU', {
    weekday: 'short', day: 'numeric', month: 'long'
  })
}
function formatTime(iso: string) {
  return new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
}
function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString('ru-RU', {
    day: 'numeric', month: 'long', hour: '2-digit', minute: '2-digit'
  })
}
</script>

<template>
  <div class="my-page">
    <div class="my-container">

      <div class="page-head">
        <div>
          <div class="section-label" style="margin-bottom:8px">{{ t('myBookings.eyebrow') }}</div>
          <h1 class="t-h1">{{ t('myBookings.title') }}<span class="red-dot">.</span></h1>
        </div>
        <button class="btn btn-outline btn-sm" @click="router.push('/book')">
          + {{ t('myBookings.book_btn') }}
        </button>
      </div>

      <div class="deco-line" />

      <div v-if="zagruzka" class="page-state">{{ t('common.loading') }}</div>

      <template v-else>

        <div v-if="autoReschedules.length > 0" class="banners-wrap">
          <div
            v-for="ar in autoReschedules"
            :key="ar.id"
            class="info-banner info-banner--blue"
          >
            <div class="info-banner__icon">🔄</div>
            <div class="info-banner__body">
              <p class="info-banner__title">{{ t('myBookings.auto_reschedule_title') }}</p>
              <p class="info-banner__detail">
                <b>{{ ar.service_name }}</b> ·
                <b>{{ ar.new_master_name }}</b> ·
                {{ formatDateTime(ar.new_starts_at) }}
              </p>
              <p class="info-banner__price">{{ Number(ar.price_paid).toLocaleString() }} ₸</p>
            </div>
            <button class="btn btn-danger btn-sm banner-btn" @click="cancelAutoReschedule(ar.id)">
              {{ t('myBookings.reschedule_cancel') }}
            </button>
          </div>
        </div>

        <div v-if="waitlist.length > 0" class="waitlist-wrap">
          <div class="section-label" style="margin-bottom:16px">{{ t('myBookings.waitlist') }}</div>
          <div v-for="w in waitlist" :key="w.id" class="waitlist-row">
            <div>
              <p style="font-size:14px;font-weight:500;color:var(--text)">{{ t('myBookings.waitlist_waiting') }}</p>
              <p style="font-size:11px;color:var(--text-3);margin-top:2px">
                {{ t('myBookings.waitlist_added') }}: {{ formatDate(w.created_at) }}
              </p>
            </div>
            <button class="btn btn-ghost btn-sm" @click="leaveWaitlist(w.id)">
              {{ t('myBookings.waitlist_leave') }}
            </button>
          </div>
          <div class="deco-line" style="margin-top:20px" />
        </div>

        <div v-if="upcoming.length > 0" class="bookings-section">
          <div class="section-label" style="margin-bottom:20px">
            {{ t('myBookings.upcoming') }} · {{ upcoming.length }}
          </div>
          <div class="bookings-list">
            <div v-for="b in upcoming" :key="b.id" class="booking-item">

              <div class="booking-item__date">
                <span class="booking-item__time">{{ formatTime(b.starts_at) }}</span>
                <span class="booking-item__day">{{ formatDate(b.starts_at) }}</span>
              </div>

              <div class="booking-item__body">
                <h3 class="booking-item__service">{{ b.service_name || '—' }}</h3>
                <p class="booking-item__master">{{ b.master_name }}</p>
                <p v-if="b.notes" class="booking-item__notes">{{ b.notes }}</p>
                <div class="booking-item__actions-mobile">
                  <AppBadge :status="b.status" />
                  <span class="booking-item__price-mobile">{{ b.price_paid.toLocaleString() }} ₸</span>
                  <button
                    v-if="['pending','confirmed'].includes(b.status)"
                    class="btn btn-danger btn-sm"
                    :disabled="cancelling === b.id"
                    @click="cancel(b.id)"
                  >
                    <span v-if="cancelling === b.id" class="btn__spinner" />
                    {{ t('myBookings.cancel_btn') }}
                  </button>
                </div>
              </div>

              <div class="booking-item__right">
                <AppBadge :status="b.status" />
                <span class="booking-item__price">{{ b.price_paid.toLocaleString() }} ₸</span>
                <div class="booking-item__actions">
                  <button
                    v-if="b.status === 'cancelled' && autoReschedules.some((r: any) => r.original_booking_id === b.id)"
                    class="btn btn-outline btn-sm"
                    style="border-color:var(--accent);color:var(--accent)"
                    @click="openReschedule(b)"
                  >🔄 {{ t('myBookings.reschedule_btn') }}</button>
                  <button
                    v-if="['pending','confirmed'].includes(b.status)"
                    class="btn btn-danger btn-sm"
                    :disabled="cancelling === b.id"
                    @click="cancel(b.id)"
                  >
                    <span v-if="cancelling === b.id" class="btn__spinner" />
                    {{ t('myBookings.cancel_btn') }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="past.length > 0">
          <div class="section-label" style="margin-bottom:20px">
            {{ t('myBookings.history') }} · {{ past.length }}
          </div>
          <div class="bookings-list">
            <div v-for="b in past" :key="b.id" class="booking-item booking-item--past">
              <div class="booking-item__date">
                <span class="booking-item__time">{{ formatTime(b.starts_at) }}</span>
                <span class="booking-item__day">{{ formatDate(b.starts_at) }}</span>
              </div>
              <div class="booking-item__body">
                <h3 class="booking-item__service">{{ b.service_name || '—' }}</h3>
                <p class="booking-item__master">{{ b.master_name }}</p>
                <div class="booking-item__actions-mobile">
                  <AppBadge :status="b.status" />
                  <button
                    v-if="b.status === 'completed' && !(b as any).has_review"
                    class="btn btn-outline btn-sm"
                    @click="openReview(b)"
                  >{{ t('myBookings.review_btn') }}</button>
                </div>
              </div>
              <div class="booking-item__right">
                <AppBadge :status="b.status" />
                <span class="booking-item__price">{{ b.price_paid.toLocaleString() }} ₸</span>
                <button
                  v-if="b.status === 'completed' && !(b as any).has_review"
                  class="btn btn-outline btn-sm"
                  @click="openReview(b)"
                >{{ t('myBookings.review_btn') }}</button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="bookings.length === 0" class="empty-state">
          <div class="empty-state__num">00</div>
          <div class="section-label" style="margin-bottom:12px">{{ t('myBookings.empty_label') }}</div>
          <p class="empty-state__title">{{ t('myBookings.empty_title') }}</p>
          <p class="empty-state__sub">{{ t('myBookings.empty_sub') }}</p>
          <button class="btn btn-primary" style="margin-top:20px" @click="router.push('/book')">
            {{ t('myBookings.book_btn') }} →
          </button>
        </div>

      </template>
    </div>
  </div>

  <Teleport to="body">
    <Transition name="fade">
      <div v-if="showReschedule" class="modal-overlay" @click.self="showReschedule = false">
        <div class="modal-box card">
          <div class="modal-header">
            <div>
              <div class="section-label" style="margin-bottom:4px">{{ t('myBookings.reschedule_modal_title') }}</div>
              <h3 class="t-h3">{{ t('myBookings.reschedule_btn') }}<span class="red-dot">.</span></h3>
            </div>
            <button class="btn btn-ghost btn-icon" @click="showReschedule = false">✕</button>
          </div>

          <div class="reschedule-discount">
            <span>🎁</span>
            <div>
              <p style="font-size:13px;font-weight:600;color:var(--success)">
                {{ t('myBookings.reschedule_discount_title') }}
              </p>
              <p style="font-size:11px;color:var(--text-3);margin-top:1px">
                {{ t('myBookings.reschedule_discount_sub') }}
              </p>
            </div>
          </div>

          <div style="display:flex;flex-direction:column;gap:16px;margin-top:16px">
            <div class="field">
              <label class="field-label">{{ t('myBookings.reschedule_master') }}</label>
              <select
                v-model="newMasterID"
                class="field-input field-select"
                @change="rescheduleSlots = []; newTime = ''"
              >
                <option v-for="m in allMasters" :key="m.id" :value="m.id">{{ m.full_name }}</option>
              </select>
            </div>
            <div class="field">
              <label class="field-label">{{ t('master.reschedule_date') }}</label>
              <input
                v-model="newDate" type="date" :min="today" class="field-input"
                @change="newTime = ''; loadRescheduleSlots()"
              />
            </div>
            <div v-if="newDate">
              <label class="field-label" style="display:block;margin-bottom:10px">{{ t('master.reschedule_time') }}</label>
              <div v-if="rescheduleSlots.length === 0" style="font-size:13px;color:var(--text-3)">
                Нет свободных слотов
              </div>
              <div v-else class="slots-grid">
                <button
                  v-for="slot in rescheduleSlots"
                  :key="slot.starts_at"
                  class="slot-chip"
                  :class="{ 'slot-chip--active': newTime === formatTime(slot.starts_at) }"
                  @click="newTime = formatTime(slot.starts_at); newSlotStartsAt = slot.starts_at"
                >{{ formatTime(slot.starts_at) }}</button>
              </div>
            </div>
            <div
              v-if="rescheduleError"
              style="color:var(--danger);font-size:13px;padding:8px 12px;background:var(--danger-bg);border-radius:var(--radius-sm)"
            >{{ rescheduleError }}</div>
            <div class="modal-footer-btns">
              <button class="btn btn-outline" @click="showReschedule = false">{{ t('common.cancel') }}</button>
              <button
                class="btn btn-primary"
                :disabled="!newDate || !newTime || rescheduling"
                @click="confirmReschedule"
              >
                <span v-if="rescheduling" class="btn__spinner" />
                Подтвердить →
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <Teleport to="body">
    <Transition name="fade">
      <div v-if="showReview" class="modal-overlay" @click.self="showReview = false">
        <div class="modal-box card">
          <div class="modal-header">
            <div>
              <div class="section-label" style="margin-bottom:4px">{{ t('myBookings.review_modal_title') }}</div>
              <h3 class="t-h3">{{ t('myBookings.review_title') }}<span class="red-dot">.</span></h3>
            </div>
            <button class="btn btn-ghost btn-icon" @click="showReview = false">✕</button>
          </div>

          <p style="font-size:13px;color:var(--text-2);margin-bottom:24px">
            {{ reviewTarget?.service_name }} · {{ reviewTarget?.master_name }}
          </p>

          <div style="margin-bottom:24px">
            <label class="field-label" style="display:block;margin-bottom:12px">
              {{ t('myBookings.review_rating') }}
            </label>
            <div class="stars-row">
              <button
                v-for="n in 5"
                :key="n"
                class="star-btn"
                :style="{
                  color: n <= (hoverRating || reviewRating) ? '#E8442A' : 'var(--border-strong)',
                  transform: n <= (hoverRating || reviewRating) ? 'scale(1.1)' : 'scale(1)'
                }"
                @click="reviewRating = n"
                @mouseenter="hoverRating = n"
                @mouseleave="hoverRating = 0"
              >★</button>
            </div>
          </div>

          <div class="field" style="margin-bottom:24px">
            <label class="field-label">{{ t('myBookings.review_comment') }}</label>
            <textarea
              v-model="reviewComment"
              class="field-input"
              rows="3"
              :placeholder="t('myBookings.review_placeholder')"
              style="resize:none"
            />
          </div>

          <div class="modal-footer-btns">
            <button class="btn btn-outline" @click="showReview = false">{{ t('common.cancel') }}</button>
            <button class="btn btn-primary" :disabled="reviewSending" @click="submitReview">
              <span v-if="reviewSending" class="btn__spinner" />
              {{ t('myBookings.review_submit') }} →
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.my-page { min-height: calc(100vh - 56px); background: var(--bg); }
.my-container {
  max-width: 800px; margin: 0 auto;
  padding: 48px 40px 80px;
}

.page-head {
  display: flex; justify-content: space-between;
  align-items: flex-end; padding-bottom: 32px;
}
.page-state {
  text-align: center; padding: 80px;
  color: var(--text-3); font-size: 14px;
}

.banners-wrap { margin-bottom: 32px; }
.info-banner {
  display: flex; align-items: center;
  gap: 14px; padding: 16px 20px;
  border-radius: var(--radius-sm); border: 1px solid;
  margin-bottom: 12px;
}
.info-banner--blue    { background: var(--info-bg); border-color: var(--info); }
.info-banner__icon    { font-size: 22px; flex-shrink: 0; }
.info-banner__body    { flex: 1; min-width: 0; }
.info-banner__title   { font-size: 14px; font-weight: 600; color: var(--text); margin-bottom: 3px; }
.info-banner__detail  { font-size: 13px; color: var(--text-2); }
.info-banner__price   { font-size: 12px; color: var(--text-3); margin-top: 3px; }
.banner-btn           { flex-shrink: 0; }

.waitlist-wrap { margin-bottom: 32px; }
.waitlist-row {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14px 0; border-top: 1px solid var(--border);
}
.waitlist-row:last-child { border-bottom: 1px solid var(--border); }

.bookings-section { margin-bottom: 48px; }
.bookings-list    { display: flex; flex-direction: column; }

.booking-item {
  display: grid;
  grid-template-columns: 120px 1fr auto;
  gap: 20px; align-items: start;
  padding: 20px 0; border-top: 1px solid var(--border);
}
.bookings-list .booking-item:last-child { border-bottom: 1px solid var(--border); }
.booking-item--past       { opacity: 0.55; transition: opacity 0.15s; }
.booking-item--past:hover { opacity: 0.8; }

.booking-item__date { display: flex; flex-direction: column; gap: 3px; padding-top: 2px; }
.booking-item__time {
  font-family: var(--font-serif); font-size: 28px;
  font-weight: 700; color: var(--text); line-height: 1;
}
.booking-item__day {
  font-size: 11px; font-weight: 500; letter-spacing: 0.04em;
  color: var(--text-3); text-transform: capitalize;
}

.booking-item__service { font-family: var(--font-serif); font-size: 18px; font-weight: 700; color: var(--text); margin-bottom: 4px; }
.booking-item__master  { font-size: 13px; color: var(--accent); font-weight: 500; }
.booking-item__notes   { font-size: 12px; color: var(--text-3); margin-top: 6px; font-style: italic; }

.booking-item__actions-mobile { display: none; }

.booking-item__right {
  display: flex; flex-direction: column; align-items: flex-end; gap: 8px;
}
.booking-item__price {
  font-family: var(--font-serif); font-size: 18px;
  font-weight: 700; color: var(--text);
}
.booking-item__actions { display: flex; gap: 6px; flex-wrap: wrap; justify-content: flex-end; }

.reschedule-discount {
  display: flex; align-items: center; gap: 12px;
  padding: 12px 14px; background: var(--success-bg);
  border: 1px solid rgba(45,106,63,0.2); border-radius: var(--radius-sm);
}

.empty-state {
  text-align: center; padding: 80px 0;
  display: flex; flex-direction: column; align-items: center;
}
.empty-state__num {
  font-family: var(--font-serif); font-size: 100px; font-weight: 900;
  color: var(--text); opacity: 0.04; line-height: 1; margin-bottom: -20px;
}
.empty-state__title { font-family: var(--font-serif); font-size: 22px; font-weight: 700; color: var(--text); }
.empty-state__sub   { font-size: 14px; color: var(--text-3); margin-top: 6px; }

.modal-footer-btns {
  display: flex; gap: 10px;
  justify-content: flex-end;
  padding-top: 12px; border-top: 1px solid var(--border);
}
.stars-row { display: flex; gap: 6px; }
.star-btn  {
  font-size: 32px; background: none; border: none;
  cursor: pointer; transition: transform 0.1s; line-height: 1; padding: 0;
  min-width: 44px; min-height: 44px;
  display: flex; align-items: center; justify-content: center;
}

@media (max-width: 600px) {
  .my-container { padding: 24px 16px 60px; }

  .page-head { flex-direction: column; align-items: flex-start; gap: 12px; padding-bottom: 20px; }
  .page-head .btn { width: 100%; justify-content: center; }

  .booking-item {
    grid-template-columns: 80px 1fr;
    gap: 14px; padding: 16px 0;
  }
  .booking-item__right { display: none; }

  .booking-item__actions-mobile {
    display: flex; align-items: center;
    gap: 8px; flex-wrap: wrap; margin-top: 10px;
    padding-top: 10px; border-top: 1px solid var(--border);
  }
  .booking-item__price-mobile {
    font-family: var(--font-serif); font-size: 16px;
    font-weight: 700; color: var(--text); margin-left: auto;
  }

  .booking-item__time { font-size: 22px; }
  .booking-item__service { font-size: 16px; }

  .info-banner { flex-wrap: wrap; gap: 10px; padding: 14px 16px; }
  .banner-btn  { width: 100%; justify-content: center; }

  .waitlist-row { flex-wrap: wrap; gap: 10px; }
  .waitlist-row .btn { width: 100%; justify-content: center; }

  .modal-box {
    position: fixed !important;
    bottom: 0 !important;
    left: 0 !important;
    right: 0 !important;
    top: auto !important;
    max-width: 100% !important;
    border-radius: var(--radius-sm) var(--radius-sm) 0 0 !important;
    max-height: 90vh;
    overflow-y: auto;
  }
  .modal-footer-btns { flex-direction: column-reverse; }
  .modal-footer-btns .btn { width: 100%; justify-content: center; }

  .slot-chip { height: 44px; font-size: 15px; }

  .empty-state { padding: 48px 0; }
  .empty-state__num { font-size: 70px; }
  .empty-state .btn { width: 100%; justify-content: center; }
}
</style>