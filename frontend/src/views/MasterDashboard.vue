<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { bookingsApi } from '@/api/bookings'
import { useLocaleFormatter } from '@/composables/useLocaleFormatter'
import { useAuthStore } from '@/stores/auth'
import { createClient } from '@supabase/supabase-js'
import { loyaltyApi }   from '@/api/loyalty'
import AppBadge from '@/components/ui/AppBadge.vue'
import type { Booking } from '@/types'
import { galleryApi } from '@/api/gallery'
import { servicesApi } from '@/api/services'
import { reviewsApi } from '@/api/reviews'
import { masterApi } from '@/api/admin'
import { watch } from 'vue'
import { Moon, Camera, GraduationCap, Zap, X, Check, CheckCheck, ArrowUpDown, Star } from 'lucide-vue-next'

async function udalitOtzyv(id: string) {
  if (!confirm('Удалить отзыв?')) return
  await reviewsApi.deleteMaster(id)
  reviews.value = reviews.value.filter((r: any) => r.id !== id)
}

const { t } = useI18n()
const { fmtDate, fmtTime, dateLocaleTag } = useLocaleFormatter()
const auth     = useAuthStore()
const bookings = ref<Booking[]>([])

const tabs = computed(() => [
  { key: 'bookings' as const, label: t('master.tab_bookings') },
  { key: 'history'  as const, label: 'История' },
  { key: 'stats'    as const, label: t('master.tab_stats') },
  { key: 'reviews'  as const, label: t('master.tab_reviews') },
  { key: 'gallery'  as const, label: 'Галерея' },
  { key: 'services' as const, label: 'Услуги' }
])

const allServices   = ref<any[]>([])
const masterSvcLoad = ref(false)

async function loadMyServices() {
  masterSvcLoad.value = true
  allServices.value   = await masterApi.getMyServices().catch(() => [])
  masterSvcLoad.value = false
}

const loading   = ref(true)
const stats     = ref<any>(null)
const reviews   = ref<any[]>([])
const sosList   = ref<any[]>([])
const activeTab = ref<'bookings' | 'history' | 'stats' | 'reviews' | 'gallery' | 'services'>('bookings')

// мобильное меню табов
const tabsOpen = ref(false)
const activeTabLabel = computed(() => tabs.value.find(t => t.key === activeTab.value)?.label ?? '')

function selectTab(key: typeof activeTab.value) {
  activeTab.value = key
  tabsOpen.value  = false
}

watch(activeTab, (tab) => {
  if (tab === 'services') loadMyServices()
  if (tab === 'gallery')  loadGalereya()
  if (tab === 'stats')    loadStats()
})

const showReschedule      = ref(false)
const rescheduleTarget    = ref<Booking | null>(null)
const rescheduleDate      = ref('')
const rescheduleTime      = ref('')
const rescheduleReason    = ref('')
const rescheduleSlots     = ref<any[]>([])
const rescheduling        = ref(false)
const rescheduleError     = ref('')
const rescheduleNewMasterID = ref('')
const rescheduleSlotISO   = ref('')
const allMastersForReschedule = ref<any[]>([])
const today = new Date().toISOString().split('T')[0]


const activeBookings = computed(() =>
  bookings.value.filter(b => ['pending', 'confirmed'].includes(b.status))
    .sort((a, b) => new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime())
)
const historyBookings = computed(() =>
  bookings.value.filter(b => ['completed', 'cancelled', 'no_show'].includes(b.status))
    .sort((a, b) => new Date(b.starts_at).getTime() - new Date(a.starts_at).getTime())
)

const vibeIcon: Record<string, any> = {
  beauty_school: GraduationCap,
  nap_time:      Moon,
  insta_vibe:    Camera,
  turbo:         Zap,
}

const galereya        = ref<any[]>([])
const zagruzkaFoto    = ref(false)
const showModalka     = ref(false)
const novoeNazvanie   = ref('')
const vybranayaUsluga = ref('')
const fileInput       = ref<HTMLInputElement | null>(null)
const vybraniyFayl    = ref<File | null>(null)
const uslugiSpisok    = ref<any[]>([])

async function loadGalereya() {
  galereya.value = await galleryApi.list(auth.user!.id).catch(() => [])
}
async function loadUslugi() {
  uslugiSpisok.value = await servicesApi.list().catch(() => [])
}
function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) vybraniyFayl.value = target.files[0]
}
async function zagruzitRabotu() {
  if (!vybraniyFayl.value) return
  zagruzkaFoto.value = true
  try {
    const formData = new FormData()
    formData.append('image', vybraniyFayl.value)
    formData.append('title', novoeNazvanie.value)
    formData.append('master_id', auth.user!.id)
    if (vybranayaUsluga.value) formData.append('service_id', vybranayaUsluga.value)
    await galleryApi.upload(formData)
    showModalka.value     = false
    novoeNazvanie.value   = ''
    vybranayaUsluga.value = ''
    vybraniyFayl.value    = null
    await loadGalereya()
  } catch (e) { console.error('upload error', e) }
  finally { zagruzkaFoto.value = false }
}
async function udalitRabotu(id: string) {
  if (!confirm('Удалить эту работу?')) return
  await galleryApi.delete(id)
  galereya.value = galereya.value.filter((w: any) => w.id !== id)
}

async function load() {
  loading.value = true
  try { bookings.value = await bookingsApi.masterList() }
  catch (e) { console.error('load error', e); bookings.value = [] }
  finally { loading.value = false }
}
async function loadStats() {
  try {
    const [s, r] = await Promise.all([
      bookingsApi.getMasterStats().catch((e) => { console.error('[stats error]', e); return null }),
      bookingsApi.getMasterReviews().catch(() => []),
    ])
    console.log('[stats response]', s)
    stats.value   = s ?? {
      total_bookings: 0,
      completed_bookings: 0,
      total_revenue: 0,
      this_month_revenue: 0,
      avg_price: 0,
      avg_rating: 0,
      total_reviews: 0,
      this_month_bookings: 0,
      top_service: null,
    }
    reviews.value = r ?? []
  } catch (e) { console.error('loadStats error', e) }
}
async function loadSOS() {
  try { sosList.value = await loyaltyApi.getMasterSOS().catch(() => []) }
  catch { sosList.value = [] }
}

async function respondSOS(id: string, status: 'accepted' | 'declined') {
  await loyaltyApi.respondSOS(id, status === 'accepted', '')
  await loadSOS()
}

onMounted(async () => {
  await Promise.all([load(), loadStats(), loadSOS(), loadGalereya(), loadUslugi(), loadMyServices()])
})

async function updateStatus(id: string, status: Booking['status']) {
  await bookingsApi.updateStatus(id, status)
  await load()
}
async function toggleService(svcId: string, isAssigned: boolean) {
  if (isAssigned) await masterApi.removeService(svcId)
  else            await masterApi.addService(svcId)
  await loadMyServices()
}
async function updateSvcPrice(svcId: string, price: string) {
  const num = parseFloat(price)
  await masterApi.updateServicePrice(svcId, isNaN(num) ? null : num)
  await loadMyServices()
}

async function openReschedule(b: Booking) {
  rescheduleTarget.value    = b
  rescheduleDate.value      = ''
  rescheduleTime.value      = ''
  rescheduleSlotISO.value   = ''
  rescheduleReason.value    = ''
  rescheduleSlots.value     = []
  rescheduleError.value     = ''
  rescheduleNewMasterID.value = ''
  if (allMastersForReschedule.value.length === 0) {
    const { mastersApi } = await import('@/api/masters')
    allMastersForReschedule.value = await mastersApi.list()
  }
  showReschedule.value = true
}
async function loadRescheduleSlots() {
  if (!rescheduleTarget.value || !rescheduleDate.value) return
  const { mastersApi } = await import('@/api/masters')
  const masterID = rescheduleNewMasterID.value || rescheduleTarget.value.master_id
  rescheduleSlots.value = await mastersApi.getSlots(
    masterID, [rescheduleTarget.value.service_id], rescheduleDate.value,
  )
}
async function confirmReschedule() {
  if (!rescheduleTarget.value || !rescheduleDate.value || !rescheduleSlotISO.value) return
  rescheduling.value    = true
  rescheduleError.value = ''
  try {
    const payload: any = { new_starts_at: rescheduleSlotISO.value, reason: rescheduleReason.value }
    if (rescheduleNewMasterID.value && rescheduleNewMasterID.value !== rescheduleTarget.value.master_id)
      payload.new_master_id = rescheduleNewMasterID.value
    await bookingsApi.masterReschedule(rescheduleTarget.value.id, payload)
    showReschedule.value = false
    await load()
  } catch (e: any) {
    rescheduleError.value = e.response?.data?.error === 'slot_unavailable'
      ? t('master.slot_busy_short') : t('master.reschedule_fail')
  } finally { rescheduling.value = false }
}

function formatTime(iso: string) { return fmtTime(iso) }
function formatDate(iso: string) {
  return fmtDate(iso, { weekday: 'long', day: 'numeric', month: 'long' })
}
function groupByDate(list: Booking[]) {
  const groups: Record<string, Booking[]> = {}
  for (const b of list) {
    const date = formatDate(b.starts_at)
    if (!groups[date]) groups[date] = []
    groups[date].push(b)
  }
  return groups
}
function fmtPrice(n: number) {
  return Math.round(n).toLocaleString(dateLocaleTag.value) + ' ' + t('common.currency')
}

const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL,
  import.meta.env.VITE_SUPABASE_ANON_KEY,
)
let channel: any
onMounted(async () => {
  channel = supabase.channel('master-bookings')
    .on('postgres_changes', {
      event: '*', schema: 'public', table: 'bookings',
      filter: `master_id=eq.${auth.user!.id}`,
    }, () => load())
    .subscribe()
})
onUnmounted(() => channel?.unsubscribe())
</script>

<template>
  <div class="dash-wrap">
    <div class="dash-container">

      <div class="dash-header">
        <div>
          <p class="t-label" style="margin-bottom:6px">{{ t('master.eyebrow') }}</p>
          <h1 class="t-h1 dash-title">{{ auth.user?.full_name }}</h1>
        </div>
        <button class="btn btn-outline btn-sm" @click="load">{{ t('master.refresh') }}</button>
      </div>

      <div class="tabs tabs--desktop" style="margin-bottom:28px">
        <button
          v-for="tabItem in tabs"
          :key="tabItem.key"
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === tabItem.key }"
          @click="activeTab = tabItem.key"
        >{{ tabItem.label }}</button>
      </div>

      <div class="tabs-mobile">
        <button class="tabs-mobile__btn" @click="tabsOpen = !tabsOpen">
          <span>{{ activeTabLabel }}</span>
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none"
            :style="{ transform: tabsOpen ? 'rotate(180deg)' : 'none' }"
            style="transition: transform 0.2s; flex-shrink:0"
          >
            <path d="M2 4l4 4 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
        <Transition name="fade">
          <div v-if="tabsOpen" class="tabs-mobile__dropdown">
            <button
              v-for="tabItem in tabs"
              :key="tabItem.key"
              class="tabs-mobile__item"
              :class="{ 'tabs-mobile__item--active': activeTab === tabItem.key }"
              @click="selectTab(tabItem.key)"
            >{{ tabItem.label }}</button>
          </div>
        </Transition>
      </div>

      <div v-if="activeTab === 'bookings'">
        <div v-if="sosList.length > 0" style="margin-bottom:24px">
          <h3 style="font-size:15px;font-weight:600;margin-bottom:10px">
            {{ t('master.sos_requests', { n: sosList.length }) }}
          </h3>
          <div v-for="sos in sosList" :key="sos.id" class="card" style="margin-bottom:10px;padding:16px">
            <div style="display:flex;justify-content:space-between;align-items:center">
              <div>
                <p style="font-weight:600;font-size:14px">{{ sos.service_name || '—' }}</p>
                <p style="font-size:12px;color:var(--text-3);margin-top:4px">
                  {{ sos.client_name || 'Клиент' }} · {{ formatDate(sos.preferred_range_start) }}
                </p>
                <p v-if="sos.client_note" style="font-size:12px;color:var(--text-2);margin-top:4px;font-style:italic">
                  {{ sos.client_note }}
                </p>
              </div>
              <div style="display:flex;gap:8px">
                <button class="btn btn-sm" style="background:var(--success-bg);color:var(--success)" @click="respondSOS(sos.id, 'accepted')">Принять</button>
                <button class="btn btn-danger btn-sm" @click="respondSOS(sos.id, 'declined')">Отклонить</button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="loading" class="state-msg">{{ t('common.loading') }}</div>
        <div v-else-if="activeBookings.length === 0" class="state-msg">{{ t('master.no_bookings') }}</div>

        <template v-else>
          <div v-for="(group, date) in groupByDate(activeBookings)" :key="date" class="date-group">
            <div class="date-group__title">{{ date }}</div>
            <div v-for="b in group" :key="b.id" class="booking-row">
              <div class="booking-row__time">{{ formatTime(b.starts_at) }}</div>
              <div class="booking-row__body">
                <p class="booking-row__service">{{ b.service_name || '—' }}</p>
                <p class="booking-row__client">{{ b.client_name || '' }}</p>
                <span v-if="(b as any).vibe_mode" class="vibe-badge">
                  <component :is="vibeIcon[(b as any).vibe_mode]" :size="14" />
                </span>
                <div class="booking-row__actions-mobile">
                  <AppBadge :status="b.status" />
                  <span class="booking-row__price-mobile">{{ fmtPrice(b.price_paid) }}</span>
                  <div class="booking-row__btns">
                    <button
                      v-if="b.status === 'pending'"
                      class="btn btn-sm"
                      style="background:var(--success-bg);color:var(--success)"
                      @click="updateStatus(b.id, 'confirmed')"
                    ><Check :size="16" /></button>
                    <button
                      v-if="b.status === 'confirmed'"
                      class="btn btn-sm"
                      style="background:var(--success-bg);color:var(--success)"
                      @click="updateStatus(b.id, 'completed')"
                    ><CheckCheck :size="16" /></button>
                    <button
                      v-if="['pending','confirmed'].includes(b.status)"
                      class="btn btn-outline btn-sm"
                      @click="openReschedule(b)"
                    ><ArrowUpDown :size="16" /></button>
                    <button
                      v-if="['pending','confirmed'].includes(b.status)"
                      class="btn btn-danger btn-sm"
                      @click="updateStatus(b.id, 'cancelled')"
                    ><X :size="16" /></button>
                  </div>
                </div>
              </div>
              <AppBadge :status="b.status" class="booking-row__badge-desktop" />
              <div class="booking-row__price booking-row__price-desktop">{{ fmtPrice(b.price_paid) }}</div>
              <div class="booking-row__actions booking-row__actions-desktop">
                <button
                  v-if="b.status === 'pending'"
                  class="btn btn-sm"
                  style="background:var(--success-bg);color:var(--success)"
                  @click="updateStatus(b.id, 'confirmed')"
                ><Check :size="16" /> {{ t('master.accept') }}</button>
                <button
                  v-if="b.status === 'confirmed'"
                  class="btn btn-sm"
                  style="background:var(--success-bg);color:var(--success)"
                  @click="updateStatus(b.id, 'completed')"
                ><CheckCheck :size="16" /> {{ t('master.complete') }}</button>
                <button
                  v-if="['pending','confirmed'].includes(b.status)"
                  class="btn btn-outline btn-sm"
                  @click="openReschedule(b)"
                ><ArrowUpDown :size="16" /> {{ t('master.reschedule') }}</button>
                <button
                  v-if="['pending','confirmed'].includes(b.status)"
                  class="btn btn-danger btn-sm"
                  @click="updateStatus(b.id, 'cancelled')"
                ><X :size="16" /> {{ t('master.cancel') }}</button>
              </div>
            </div>
          </div>
        </template>
      </div>

      <div v-else-if="activeTab === 'history'">
        <div v-if="historyBookings.length === 0" class="state-msg">История пуста</div>
        <template v-else>
          <div class="section-label" style="margin-bottom:20px">
            Обработанные записи · {{ historyBookings.length }}
          </div>
          <div class="bookings-list">
            <div
              v-for="b in historyBookings"
              :key="b.id"
              class="booking-row booking-row--past"
            >
              <div class="booking-row__time">
                <span class="booking-row__hour">{{ formatTime(b.starts_at) }}</span>
                <span class="booking-row__date">{{ formatDate(b.starts_at) }}</span>
              </div>
              <div class="booking-row__body">
                <p class="booking-row__service">{{ b.service_name || '—' }}</p>
                <p class="booking-row__client">{{ b.client_name || '—' }}</p>
                <p v-if="b.notes" class="booking-row__notes">{{ b.notes }}</p>
              </div>
              <div class="booking-row__right">
                <AppBadge :status="b.status" />
                <span class="booking-row__price">{{ fmtPrice(b.price_paid) }}</span>
              </div>
            </div>
          </div>
        </template>
      </div>

      <div v-else-if="activeTab === 'stats'">
        <div v-if="!stats" class="state-msg">{{ t('common.loading') }}</div>
        <template v-else>
          <div class="kpi-grid" style="margin-bottom:24px">
            <div class="kpi-card">
              <p class="kpi-label">{{ t('master.stats_total') }}</p>
              <p class="kpi-val">{{ stats.total_bookings }}</p>
            </div>
            <div class="kpi-card">
              <p class="kpi-label">{{ t('master.stats_completed') }}</p>
              <p class="kpi-val" style="color:var(--success)">{{ stats.completed_bookings }}</p>
            </div>
            <div class="kpi-card">
              <p class="kpi-label">{{ t('master.stats_revenue') }}</p>
              <p class="kpi-val">{{ fmtPrice(stats.total_revenue) }}</p>
            </div>
            <div class="kpi-card">
              <p class="kpi-label">{{ t('master.stats_month') }}</p>
              <p class="kpi-val" style="color:var(--accent)">{{ fmtPrice(Number(stats.this_month_revenue) || 0) }}</p>
            </div>
            <div class="kpi-card">
              <p class="kpi-label">{{ t('master.stats_avg') }}</p>
              <p class="kpi-val">{{ fmtPrice(stats.avg_price) }}</p>
            </div>
            <div class="kpi-card">
              <p class="kpi-label">{{ t('master.stats_rating') }}</p>
              <p class="kpi-val">
                {{ stats.avg_rating > 0 ? stats.avg_rating.toFixed(1) : '—' }}
                <span style="font-size:16px;color:var(--warning)"><Star :size="16" /></span>
              </p>
            </div>
          </div>
          <div class="card">
            <p class="t-label" style="margin-bottom:12px">{{ t('master.stats_extra') }}</p>
            <div class="stat-row">
              <span style="color:var(--text-2)">{{ t('master.stats_reviews') }}</span>
              <span>{{ stats.total_reviews }}</span>
            </div>
            <div class="stat-row">
              <span style="color:var(--text-2)">{{ t('master.stats_month_bookings') }}</span>
              <span>{{ stats.this_month_bookings }}</span>
            </div>
            <div class="stat-row" style="border:none">
              <span style="color:var(--text-2)">{{ t('master.stats_top') }}</span>
              <span>{{ stats.top_service || '—' }}</span>
            </div>
          </div>
        </template>
      </div>

      <div v-else-if="activeTab === 'reviews'">
        <div v-if="reviews.length === 0" class="state-msg">{{ t('master.no_reviews') }}</div>
        <div v-else class="reviews-list">
          <div v-for="rv in reviews" :key="rv.id" class="review-card card">
            <div class="review-card__header">
              <span class="review-stars" :style="rv.rating >= 4 ? 'color:var(--warning)' : 'color:var(--text-3)'">
                <span class="review-card__stars">
                  <template v-for="n in 5" :key="n">
                    <Star
                      :size="16"
                      :fill="n <= rv.rating ? 'var(--warning)' : 'none'"
                      :stroke="n <= rv.rating ? 'var(--warning)' : 'var(--border-strong)'"
                    />
                  </template>
                </span>
              </span>
              <div style="display:flex;align-items:center;gap:8px">
                <span style="font-size:12px;color:var(--text-3)">{{ fmtDate(rv.created_at) }}</span>
                <button
                  class="btn btn-danger btn-sm"
                  style="height:28px;padding:0 10px;font-size:12px"
                  @click="udalitOtzyv(rv.id)"
                >Удалить</button>
              </div>
            </div>
            <p v-if="rv.comment" style="font-size:14px;color:var(--text-2);margin-top:8px;line-height:1.5">
              {{ rv.comment }}
            </p>
          </div>
        </div>
      </div>

      <div v-else-if="activeTab === 'gallery'">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:20px">
          <p style="font-size:15px;color:var(--text-2)">{{ galereya.length }} работ</p>
          <button class="btn btn-accent btn-sm" @click="showModalka = true">+ Добавить</button>
        </div>
        <div v-if="galereya.length === 0" class="state-msg">Работы пока не добавлены</div>
        <div v-else class="galereya-grid">
          <div v-for="rabota in galereya" :key="rabota.id" class="galereya-item">
            <img :src="rabota.image_url" :alt="rabota.title || 'Work'" loading="lazy" />
            <div class="galereya-item__overlay">
              <p v-if="rabota.title" class="galereya-item__title">{{ rabota.title }}</p>
              <p v-if="rabota.service_name" style="font-size:11px;color:rgba(255,255,255,0.7)">
                {{ rabota.service_name }}
              </p>
              <button class="galereya-item__delete" @click.stop="udalitRabotu(rabota.id)"><X :size="16" /></button>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="activeTab === 'services'">
        <div v-if="masterSvcLoad" class="state-msg">Загрузка...</div>
        <div v-else class="svc-mgmt-list">
          <div
            v-for="s in allServices"
            :key="s.id"
            class="svc-mgmt-row"
            :class="{ 'svc-mgmt-row--active': s.is_assigned }"
          >
            <div class="svc-mgmt-row__info">
              <p class="svc-mgmt-row__name">{{ s.name }}</p>
              <p class="svc-mgmt-row__meta">{{ s.duration_min }} мин · {{ s.price.toLocaleString() }} ₸</p>
              <div v-if="s.is_assigned" class="svc-mgmt-row__price">
                <label style="font-size:11px;color:var(--text-3);letter-spacing:0.06em;text-transform:uppercase">
                  Ваша цена (₸)
                </label>
                <input
                  type="number" class="field-input"
                  style="width:140px;height:36px;padding:4px 10px;font-size:14px"
                  :value="s.effective_price" :placeholder="String(s.price)"
                  @change="updateSvcPrice(s.id, ($event.target as HTMLInputElement).value)"
                />
              </div>
            </div>
            <button
              class="btn btn-sm"
              :class="s.is_assigned ? 'btn-danger' : 'btn-outline'"
              @click="toggleService(s.id, s.is_assigned)"
            >{{ s.is_assigned ? 'Убрать' : 'Добавить' }}</button>
          </div>
        </div>
      </div>

    </div>
  </div>

  <Teleport to="body">
    <Transition name="fade">
      <div v-if="showReschedule" class="modal-overlay" @click.self="showReschedule = false">
        <div class="modal-box card">
          <div class="modal-header">
            <h3 class="t-h3">{{ t('master.reschedule_title') }}</h3>
            <button class="btn btn-ghost btn-icon" @click="showReschedule = false"><X :size="16" /></button>
          </div>
          <div style="background:var(--success-bg);padding:12px 16px;border-radius:var(--radius-md);margin-bottom:16px">
            <p style="font-size:13px;color:var(--success);font-weight:500">{{ t('master.reschedule_discount') }}</p>
          </div>
          <div style="display:flex;flex-direction:column;gap:16px">
            <div class="field">
              <label class="field-label">Мастер</label>
              <select v-model="rescheduleNewMasterID" class="field-input field-select"
                @change="rescheduleSlots = []; rescheduleTime = ''; rescheduleSlotISO = ''"
              >
                <option value="">Оставить себе</option>
                <option v-for="m in allMastersForReschedule" :key="m.id" :value="m.id">{{ m.full_name }}</option>
              </select>
            </div>
            <div class="field">
              <label class="field-label">{{ t('master.reschedule_date') }}</label>
              <input v-model="rescheduleDate" type="date" :min="today" class="field-input"
                @change="rescheduleTime = ''; rescheduleSlotISO = ''; loadRescheduleSlots()" />
            </div>
            <div v-if="rescheduleDate">
              <label class="field-label" style="display:block;margin-bottom:8px">{{ t('master.reschedule_time') }}</label>
              <div v-if="rescheduleSlots.length === 0" style="font-size:13px;color:var(--text-3)">{{ t('master.no_slots') }}</div>
              <div v-else class="slots-grid">
                <button
                  v-for="slot in rescheduleSlots" :key="slot.starts_at"
                  class="slot-chip"
                  :class="{ 'slot-chip--active': rescheduleSlotISO === slot.starts_at }"
                  @click="rescheduleTime = formatTime(slot.starts_at); rescheduleSlotISO = slot.starts_at"
                >{{ formatTime(slot.starts_at) }}</button>
              </div>
            </div>
            <div class="field">
              <label class="field-label">{{ t('master.reschedule_reason') }}</label>
              <input v-model="rescheduleReason" type="text" class="field-input"
                :placeholder="t('master.reschedule_reason_placeholder')" />
            </div>
            <div v-if="rescheduleError" style="color:var(--danger);font-size:13px;padding:8px 12px;background:var(--danger-bg);border-radius:var(--radius-sm)">
              {{ rescheduleError }}
            </div>
            <div class="modal-footer-btns">
              <button class="btn btn-outline" @click="showReschedule = false">{{ t('common.cancel') }}</button>
              <button class="btn btn-accent" :disabled="!rescheduleDate || !rescheduleSlotISO || rescheduling" @click="confirmReschedule">
                <span v-if="rescheduling" class="btn__spinner" />
                {{ t('master.confirm_reschedule') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <Teleport to="body">
    <Transition name="fade">
      <div v-if="showModalka" class="modal-overlay" @click.self="showModalka = false">
        <div class="modal-box card">
          <div class="modal-header">
            <h3 class="t-h3">Добавить работу</h3>
            <button class="btn btn-ghost btn-icon" @click="showModalka = false"><X :size="16" /></button>
          </div>
          <div style="display:flex;flex-direction:column;gap:16px;margin-top:8px">
            <div class="field">
              <label class="field-label">Фото</label>
              <div class="upload-zone" @click="fileInput?.click()">
                <div v-if="!vybraniyFayl">
                  <p style="font-size:14px;color:var(--text-2)">Нажмите чтобы выбрать фото</p>
                  <p style="font-size:12px;color:var(--text-3);margin-top:4px">JPG, PNG, WEBP</p>
                </div>
                <div v-else style="display:flex;align-items:center;gap:10px">
                  <span style="font-size:14px;color:var(--success);font-weight:700">Файл выбран</span>
                  <span style="font-size:14px;color:var(--text)">{{ vybraniyFayl.name }}</span>
                </div>
              </div>
              <input ref="fileInput" type="file" accept="image/*" style="display:none" @change="onFileChange" />
            </div>
            <div class="field">
              <label class="field-label">Название (необязательно)</label>
              <input v-model="novoeNazvanie" type="text" class="field-input" placeholder="Например: окрашивание балаяж" />
            </div>
            <div class="field">
              <label class="field-label">Услуга (необязательно)</label>
              <select v-model="vybranayaUsluga" class="field-input field-select">
                <option value="">Не выбрана</option>
                <option v-for="u in uslugiSpisok" :key="u.id" :value="u.id">{{ u.name }}</option>
              </select>
            </div>
            <div class="modal-footer-btns">
              <button class="btn btn-outline" @click="showModalka = false">{{ t('common.cancel') }}</button>
              <button class="btn btn-primary" :disabled="!vybraniyFayl || zagruzkaFoto" @click="zagruzitRabotu">
                <span v-if="zagruzkaFoto" class="btn__spinner" />
                Загрузить
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.dash-wrap      { min-height: calc(100vh - 56px); background: var(--bg); }
.dash-container {
  max-width: 900px; margin: 0 auto;
  padding: 40px 40px 80px;
}

.dash-header {
  display: flex; justify-content: space-between;
  align-items: flex-end; margin-bottom: 28px;
}
.dash-title { font-size: clamp(24px, 5vw, 42px); }

.tabs--desktop {
  display: flex; gap: 4px;
  border-bottom: 1px solid var(--border);
  overflow-x: auto; scrollbar-width: none;
}
.tabs--desktop::-webkit-scrollbar { display: none; }
.tab-btn {
  padding: 10px 20px; border: none; background: transparent;
  font-size: 14px; font-weight: 500; color: var(--text-2);
  border-bottom: 2px solid transparent; margin-bottom: -1px;
  cursor: pointer; transition: all 0.15s; white-space: nowrap; flex-shrink: 0;
}
.tab-btn--active { color: var(--accent); border-bottom-color: var(--accent); }

.tabs-mobile { display: none; margin-bottom: 24px; position: relative; }
.tabs-mobile__btn {
  display: flex; align-items: center; justify-content: space-between;
  width: 100%; padding: 12px 16px; height: 44px;
  background: var(--bg-card); border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-family: var(--font); font-size: 14px; font-weight: 500;
  color: var(--text); cursor: pointer; transition: border-color 0.15s;
}
.tabs-mobile__btn:hover { border-color: var(--border-strong); }
.tabs-mobile__dropdown {
  position: absolute; top: calc(100% + 4px); left: 0; right: 0;
  z-index: 50; background: var(--bg-card);
  border: 1px solid var(--border); border-radius: var(--radius-sm);
  overflow: hidden; box-shadow: var(--shadow-md);
}
.tabs-mobile__item {
  display: block; width: 100%; text-align: left;
  padding: 13px 16px; background: none; border: none;
  font-family: var(--font); font-size: 14px; font-weight: 500;
  color: var(--text-2); cursor: pointer; transition: background 0.1s;
  border-bottom: 1px solid var(--border);
}
.tabs-mobile__item:last-child { border-bottom: none; }
.tabs-mobile__item:hover        { background: var(--bg-hover); color: var(--text); }
.tabs-mobile__item--active      { color: var(--accent); background: var(--accent-bg); }

.date-group        { margin-bottom: 32px; }
.date-group__title {
  font-size: 13px; font-weight: 600; color: var(--text-2);
  padding: 10px 0; border-bottom: 1px solid var(--border);
  text-transform: capitalize;
}
.booking-row {
  display: grid;
  grid-template-columns: 70px 1fr auto auto auto;
  gap: 16px; align-items: center;
  padding: 16px 0; border-bottom: 1px solid var(--border);
}
.booking-row__time    { font-family: var(--font-serif); font-size: 26px; font-weight: 300; color: var(--text); }
.booking-row__service { font-size: 14px; font-weight: 600; color: var(--text); margin-bottom: 2px; }
.booking-row__client  { font-size: 12px; color: var(--text-3); }
.vibe-badge           { font-size: 14px; margin-top: 2px; display: inline-block; }
.booking-row__price   { font-size: 15px; font-weight: 700; white-space: nowrap; }
.booking-row__actions { display: flex; gap: 6px; flex-wrap: wrap; justify-content: flex-end; }

.booking-row__actions-mobile { display: none; }

.booking-row {
  display: grid;
  grid-template-columns: 110px 1fr auto;
  gap: 16px; align-items: start;
  padding: 16px 0; border-top: 1px solid var(--border);
}
.bookings-list .booking-row:last-child { border-bottom: 1px solid var(--border); }
.booking-row__time    { display: flex; flex-direction: column; gap: 2px; }
.booking-row__hour    { font-family: var(--font-serif); font-size: 22px; font-weight: 700; color: var(--text); line-height: 1; }
.booking-row__date    { font-size: 10px; color: var(--text-3); text-transform: capitalize; }
.booking-row__service { font-size: 15px; font-weight: 600; color: var(--text); margin-bottom: 3px; }
.booking-row__client  { font-size: 12px; color: var(--text-3); }
.booking-row__notes   { font-size: 11px; color: var(--text-3); font-style: italic; margin-top: 4px; }
.booking-row__right   { display: flex; flex-direction: column; align-items: flex-end; gap: 6px; }
.booking-row__price   { font-family: var(--font-serif); font-size: 15px; font-weight: 700; color: var(--text); }

.kpi-grid  { display: grid; grid-template-columns: repeat(auto-fill,minmax(140px,1fr)); gap: 12px; }
.kpi-card  {
  background: var(--bg-card); border: 1px solid var(--border);
  border-radius: var(--radius-md); padding: 14px 16px;
}
.kpi-label { font-size: 11px; font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; color: var(--text-3); margin-bottom: 8px; }
.kpi-val   { font-size: 22px; font-weight: 700; color: var(--text); }
.stat-row  { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid var(--border); font-size: 14px; }

.reviews-list          { display: flex; flex-direction: column; gap: 10px; }
.review-card__header   { display: flex; justify-content: space-between; align-items: center; }
.review-stars          { font-size: 18px; letter-spacing: 2px; }

.galereya-grid {
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px;
}
.galereya-item {
  aspect-ratio: 1/1; border-radius: var(--radius-md);
  overflow: hidden; position: relative; background: var(--bg-muted); cursor: pointer;
}
.galereya-item img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.3s; }
.galereya-item:hover img { transform: scale(1.05); }
.galereya-item__overlay {
  position: absolute; inset: 0;
  background: linear-gradient(transparent 50%, rgba(0,0,0,0.7));
  padding: 10px; opacity: 0; transition: opacity 0.2s;
  display: flex; flex-direction: column; justify-content: flex-end;
}
.galereya-item:hover .galereya-item__overlay { opacity: 1; }
.galereya-item__title   { font-size: 12px; color: #fff; font-weight: 500; }
.galereya-item__delete  {
  position: absolute; top: 8px; right: 8px;
  width: 26px; height: 26px; border-radius: 50%;
  background: rgba(220,38,38,0.85); color: #fff; border: none;
  cursor: pointer; font-size: 12px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  transition: background 0.15s;
}
.galereya-item__delete:hover { background: var(--danger); }

.svc-mgmt-list { display: flex; flex-direction: column; }
.svc-mgmt-row {
  display: flex; justify-content: space-between; align-items: center;
  gap: 16px; padding: 16px 0; border-top: 1px solid var(--border); transition: all 0.15s;
}
.svc-mgmt-row:last-child { border-bottom: 1px solid var(--border); }
.svc-mgmt-row--active    { border-left: 2px solid var(--accent); padding-left: 12px; }
.svc-mgmt-row__name  { font-family: var(--font-serif); font-size: 17px; font-weight: 700; color: var(--text); margin-bottom: 3px; }
.svc-mgmt-row__meta  { font-size: 12px; color: var(--text-3); }
.svc-mgmt-row__price { margin-top: 8px; display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }

.state-msg { text-align: center; padding: 60px; color: var(--text-3); font-size: 14px; }
.slots-grid { display: grid; grid-template-columns: repeat(auto-fill,minmax(80px,1fr)); gap: 8px; }
.slot-chip  {
  padding: 10px 8px; border-radius: var(--radius-md);
  border: 1.5px solid var(--border); background: var(--bg-card);
  color: var(--text-2); font-size: 14px; font-weight: 500; cursor: pointer; transition: all 0.15s;
}
.slot-chip:hover   { border-color: var(--accent); color: var(--accent); }
.slot-chip--active { background: var(--accent); border-color: var(--accent); color: #fff; }

.modal-overlay { position: fixed; inset: 0; z-index: 1000; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; padding: 16px; }
.modal-box     { width: 100%; max-width: 480px; max-height: 90vh; overflow-y: auto; padding: 24px; }
.modal-header  { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.modal-footer-btns { display: flex; gap: 10px; justify-content: flex-end; padding-top: 8px; border-top: 1px solid var(--border); }

.upload-zone {
  border: 2px dashed var(--border-strong); border-radius: var(--radius-md);
  padding: 32px 20px; text-align: center; cursor: pointer;
  transition: all 0.15s; background: var(--bg);
}
.upload-zone:hover { border-color: var(--accent); background: var(--accent-bg); }

@media (max-width: 700px) {
  .dash-container { padding: 20px 16px 60px; }

  .dash-header { margin-bottom: 20px; flex-wrap: wrap; gap: 12px; }
  .dash-header .btn { width: 100%; justify-content: center; }

  .tabs--desktop { display: none; }
  .tabs-mobile   { display: block; }

  .booking-row {
    grid-template-columns: 60px 1fr;
    gap: 12px;
  }
  .booking-row__badge-desktop   { display: none; }
  .booking-row__price-desktop   { display: none; }
  .booking-row__actions-desktop { display: none; }

  .booking-row__actions-mobile {
    display: flex; flex-direction: column; gap: 8px; margin-top: 10px;
  }
  .booking-row__price-mobile {
    font-family: var(--font-serif); font-size: 16px;
    font-weight: 700; color: var(--text);
  }
  .booking-row__btns { display: flex; gap: 6px; flex-wrap: wrap; }
  .booking-row__time { font-size: 20px; }

  .kpi-grid  { grid-template-columns: repeat(2, 1fr); }
  .kpi-val   { font-size: 18px; }

  .galereya-grid { grid-template-columns: repeat(2, 1fr); }
  .galereya-item__overlay { opacity: 1; }

  .svc-mgmt-row { flex-direction: column; align-items: flex-start; gap: 10px; }
  .svc-mgmt-row .btn { width: 100%; justify-content: center; }
  .svc-mgmt-row__price { width: 100%; }
  .svc-mgmt-row__price input { width: 100% !important; }

  .modal-overlay  { align-items: flex-end; padding: 0; }
  .modal-box      {
    max-width: 100%; border-radius: var(--radius-sm) var(--radius-sm) 0 0;
    max-height: 92vh;
  }
  .modal-footer-btns { flex-direction: column-reverse; }
  .modal-footer-btns .btn { width: 100%; justify-content: center; }

  .slot-chip { height: 44px; font-size: 15px; }
}
</style>