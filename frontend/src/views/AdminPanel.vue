<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '@/api'
import { useLocaleFormatter } from '@/composables/useLocaleFormatter'
import { analyticsApi, type AnalyticsRange } from '@/api/analytics'
import AppBadge from '@/components/ui/AppBadge.vue'
import type { Booking, Master, Service } from '@/types'
import {
  Chart as ChartJS,
  CategoryScale, LinearScale, PointElement,
  LineElement, BarElement, ArcElement,
  Title, Tooltip, Legend, Filler,
} from 'chart.js'
import { Line, Bar, Doughnut } from 'vue-chartjs'

ChartJS.register(
  CategoryScale, LinearScale, PointElement,
  LineElement, BarElement, ArcElement,
  Title, Tooltip, Legend, Filler,
)

const { t } = useI18n()
const { fmtDate, fmtDateTime, dateLocaleTag } = useLocaleFormatter()

const activeTab = ref<'analytics' | 'bookings' | 'masters' | 'services'>('analytics')

const adminTabs = computed(() => [
  { key: 'analytics' as const, label: t('admin.tab_analytics') },
  { key: 'bookings'  as const, label: t('admin.tab_bookings') },
  { key: 'masters'   as const, label: t('admin.tab_masters') },
  { key: 'services'  as const, label: t('admin.tab_services') },
])

function rangeLabel(r: AnalyticsRange) {
  return ({
    week: t('admin.range_week'),
    month: t('admin.range_month'),
    quarter: t('admin.range_quarter'),
    year: t('admin.range_year'),
  } as const)[r]
}

function chartShortDate(day: string) {
  return fmtDate(day, { day: 'numeric', month: 'short' })
}
const range     = ref<AnalyticsRange>('month')
const analytics = ref<any>(null)
const bookings  = ref<Booking[]>([])
const masters   = ref<Master[]>([])
const services  = ref<Service[]>([])
const loading   = ref(true)
const filterStatus = ref('all')
const filterMaster = ref('all')

async function loadAnalytics() {
  loading.value = true
  analytics.value = await analyticsApi.get(range.value)
  loading.value   = false
}

async function loadAll() {
  loading.value = true
  const [b, m, s] = await Promise.all([
    api.get('/admin/bookings').then(r => r.data),
    api.get('/masters').then(r => r.data),
    api.get('/admin/services/all').then(r => r.data),
  ])
  bookings.value = b ?? []
  masters.value  = m ?? []
  services.value = s ?? []
  loading.value  = false
}

onMounted(async () => {
  await Promise.all([loadAnalytics(), loadAll()])
})

watch(range, loadAnalytics)

const filteredBookings = computed(() =>
  bookings.value.filter(b => {
    const statusOk = filterStatus.value === 'all' || b.status === filterStatus.value
    const masterOk = filterMaster.value === 'all' || b.master_id === filterMaster.value
    return statusOk && masterOk
  })
)

const isDark = computed(() =>
  document.documentElement.getAttribute('data-theme') === 'dark'
)

const textColor = computed(() => isDark.value ? '#A09890' : '#6B6560')
const gridColor = computed(() => isDark.value ? '#2E2C28' : '#F4F2EF')
const accentColor = '#C9956A'

const revenueChart = computed(() => {
  if (!analytics.value?.forecast) return null

  const forecast = analytics.value.forecast as Array<{
    day: string
    revenue: number
    is_forecast: boolean
    lower: number
    upper: number
  }>

  const histLabels     = forecast.filter(f => !f.is_forecast).map(f =>
    new Date(f.day).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
  )
  const forecastLabels = forecast.filter(f => f.is_forecast).map(f =>
    new Date(f.day).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
  )
  const allLabels = [...histLabels, ...forecastLabels]

  const histValues     = forecast.filter(f => !f.is_forecast).map(f => f.revenue)
  const forecastValues = forecast.filter(f => f.is_forecast).map(f => f.revenue)
  const lowerValues    = forecast.filter(f => f.is_forecast).map(f => f.lower)
  const upperValues    = forecast.filter(f => f.is_forecast).map(f => f.upper)

  const histPad     = new Array(histValues.length - 1).fill(null)
  const forecastPad = new Array(histValues.length - 1).fill(null)

  return {
    labels: allLabels,
    datasets: [
      {
        label: 'Выручка',
        data: [...histValues, ...new Array(forecastValues.length).fill(null)],
        borderColor: '#E8442A',
        backgroundColor: 'rgba(232,68,42,0.08)',
        fill: true,
        tension: 0.4,
        pointRadius: 2,
      },
      {
        label: 'Прогноз (Holt-Winters)',
        data: [...histPad, histValues[histValues.length-1], ...forecastValues],
        borderColor: '#E8442A',
        borderDash: [6, 4],
        backgroundColor: 'transparent',
        fill: false,
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: 'Верхняя граница',
        data: [...forecastPad, histValues[histValues.length-1], ...upperValues],
        borderColor: 'rgba(232,68,42,0.2)',
        backgroundColor: 'rgba(232,68,42,0.06)',
        fill: '+1',
        tension: 0.4,
        pointRadius: 0,
        borderWidth: 1,
      },
      {
        label: 'Нижняя граница',
        data: [...forecastPad, histValues[histValues.length-1], ...lowerValues],
        borderColor: 'rgba(232,68,42,0.2)',
        backgroundColor: 'transparent',
        fill: false,
        tension: 0.4,
        pointRadius: 0,
        borderWidth: 1,
      },
    ],
  }
})

const servicesChart = computed(() => {
  if (!analytics.value?.services) return null
  const sv = analytics.value.services
  const palette = ['#C9956A','#8B5CF6','#EC4899','#3B82F6','#10B981','#F59E0B','#EF4444']
  return {
    labels:   sv.map((s: any) => s.service_name),
    datasets: [{ data: sv.map((s: any) => s.count), backgroundColor: palette, borderWidth: 0 }],
  }
})

const bookingsBarChart = computed(() => {
  if (!analytics.value?.daily) return null
  const daily = analytics.value.daily
  return {
    labels: daily.map((d: any) => chartShortDate(d.day)),
    datasets: [{
      label: t('admin.chart_bookings_label'),
      data: daily.map((d: any) => d.bookings_count),
      backgroundColor: accentColor + 'AA',
      borderRadius: 6,
    }],
  }
})

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: (ctx: any) => {
          const v = ctx.raw
          if (v === null) return ''
          const revenueLbl = t('admin.chart_revenue_label')
          const forecastLbl = t('admin.chart_forecast_label')
          return ctx.dataset.label === revenueLbl || ctx.dataset.label === forecastLbl
            ? ` ${Number(v).toLocaleString(dateLocaleTag.value)} ${t('common.currency')}`
            : ` ${v}`
        }
      }
    }
  },
  scales: {
    x: { ticks: { color: textColor.value, font: { size: 11 } }, grid: { color: gridColor.value } },
    y: { ticks: { color: textColor.value, font: { size: 11 } }, grid: { color: gridColor.value } },
  },
}))

const doughnutOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { position: 'bottom' as const } },
  cutout: '65%',
}

async function updateBookingStatus(id: string, status: Booking['status']) {
  await api.patch(`/admin/bookings/${id}`, { status })
  bookings.value = bookings.value.map(b =>
    b.id === id ? { ...b, status } : b
  )
}
async function toggleMaster(id: string, isActive: boolean) {
  await api.patch(`/admin/masters/${id}`, { is_active: !isActive })
  masters.value = masters.value.map(m =>
    m.id === id ? { ...m, is_active: !isActive } : m
  )
}
async function toggleService(id: string, isActive: boolean) {
  await api.patch(`/admin/services/${id}`, { is_active: !isActive })
  services.value = services.value.map(s =>
    s.id === id ? { ...s, is_active: !isActive } : s
  )
}

function fmtPrice(n: number) {
  const num = isNaN(n) ? 0 : Math.round(n)
  return num.toLocaleString(dateLocaleTag.value) + ' ' + t('common.currency')
}

function masterName(id: string) {
  return masters.value.find(m => m.id === id)?.full_name ?? '—'
}
function serviceName(id: string) {
  return services.value.find(s => s.id === id)?.name ?? '—'
}
function formatDate(iso: string) {
  return fmtDateTime(iso, {
    day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit',
  })
}

const showAddMaster  = ref(false)
const showAddService = ref(false)

const newMaster = ref({
  email:            '',
  password:         '',
  full_name:        '',
  bio:              '',
  experience_years: 0,
  instagram:        '',
})
const savingMaster = ref(false)

async function createMaster() {
  if (!newMaster.value.email || !newMaster.value.full_name) return
  savingMaster.value = true
  try {
    await api.post('/admin/masters', newMaster.value)
    showAddMaster.value = false
    newMaster.value = {
      email:'', password:'', full_name:'',
      bio:'', experience_years:0, instagram:'',
    }
    await loadAll()
  } catch (e: any) {
    alert(e.response?.data?.error || 'Ошибка создания мастера')
  } finally {
    savingMaster.value = false
  }
}

const newService = ref({
  name: '', category: 'hair', description: '',
  duration_min: 60, price: 0
})

async function deleteMaster(id: string) {
  if (!confirm('Удалить мастера? Это действие необратимо.')) return
  await api.delete(`/admin/masters/${id}`)
  masters.value = masters.value.filter(m => m.id !== id)
}

async function deleteService(id: string) {
  if (!confirm('Удалить услугу?')) return
  await api.delete(`/admin/services/${id}`)
  services.value = services.value.filter(s => s.id !== id)
}

async function createService() {
  await api.post('/admin/services', newService.value)
  showAddService.value = false
  newService.value = { name:'', category:'hair', description:'', duration_min:60, price:0 }
  await loadAll()
}

</script>

<template>
  <div class="page-container" style="padding-top: 32px; padding-bottom: 80px;">

    <div class="admin-header">
      <div>
        <p class="t-label" style="margin-bottom: 6px;">{{ t('admin.eyebrow') }}</p>
        <h1 class="t-h1">{{ t('admin.title') }}</h1>
      </div>
      <button class="btn btn-outline btn-sm" @click="loadAll(); loadAnalytics()">
        {{ t('admin.refresh') }}
      </button>
    </div>

    <div class="tabs">
      <button
        v-for="tab in adminTabs"
        :key="tab.key"
        class="tab-btn"
        :class="{ 'tab-btn--active': activeTab === tab.key }"
        @click="activeTab = tab.key"
      >{{ tab.label }}</button>
    </div>

    <div v-if="loading" class="loading-state">{{ t('common.loading') }}</div>

    <div v-else-if="activeTab === 'analytics'">

      <div class="range-row">
        <button
          v-for="r in (['week','month','quarter','year'] as AnalyticsRange[])"
          :key="r"
          class="range-btn"
          :class="{ 'range-btn--active': range === r }"
          @click="range = r"
        >{{ rangeLabel(r) }}</button>
      </div>

      <div class="kpi-grid" v-if="analytics?.overview">
        <div class="kpi-card">
          <p class="kpi-card__label">{{ t('admin.kpi_revenue') }}</p>
          <p class="kpi-card__val">{{ fmtPrice(analytics.overview.revenue) }}</p>
        </div>
        <div class="kpi-card">
          <p class="kpi-card__label">{{ t('admin.kpi_total') }}</p>
          <p class="kpi-card__val">{{ analytics.overview.total }}</p>
        </div>
        <div class="kpi-card">
          <p class="kpi-card__label">{{ t('admin.kpi_completed') }}</p>
          <p class="kpi-card__val" style="color: var(--success)">{{ analytics.overview.completed }}</p>
        </div>
        <div class="kpi-card">
          <p class="kpi-card__label">{{ t('admin.kpi_avg') }}</p>
          <p class="kpi-card__val">{{ fmtPrice(analytics.overview.avg_price) }}</p>
        </div>
        <div class="kpi-card">
          <p class="kpi-card__label">{{ t('admin.kpi_cancel') }}</p>
          <p class="kpi-card__val" :style="analytics.overview.cancel_rate > 20 ? 'color:var(--danger)' : ''">
            {{ analytics.overview.cancel_rate }}%
          </p>
        </div>
        <div class="kpi-card">
          <p class="kpi-card__label">{{ t('admin.kpi_noshow') }}</p>
          <p class="kpi-card__val">{{ analytics.overview.no_show }}</p>
        </div>
      </div>

      <div class="charts-grid">

        <div class="chart-card card">
          <div class="chart-card__header">
            <h3 class="chart-card__title">{{ t('admin.chart_revenue') }}</h3>
            <span class="badge badge-accent">{{ t('admin.chart_forecast') }}</span>
          </div>
          <div class="chart-card__body">
            <Line v-if="revenueChart" :data="revenueChart" :options="chartOptions" />
          </div>
        </div>

        <div class="chart-card card">
          <div class="chart-card__header">
            <h3 class="chart-card__title">{{ t('admin.chart_bookings') }}</h3>
          </div>
          <div class="chart-card__body">
            <Bar v-if="bookingsBarChart" :data="bookingsBarChart" :options="chartOptions" />
          </div>
        </div>

        <div class="chart-card card">
          <div class="chart-card__header">
            <h3 class="chart-card__title">{{ t('admin.chart_services') }}</h3>
          </div>
          <div class="chart-card__body" style="height: 220px;">
            <Doughnut v-if="servicesChart" :data="servicesChart" :options="doughnutOptions" />
          </div>
        </div>

        <div class="chart-card card">
          <div class="chart-card__header">
            <h3 class="chart-card__title">{{ t('admin.chart_masters') }}</h3>
          </div>
          <div v-if="analytics?.masters?.length">
            <div
              v-for="m in analytics.masters"
              :key="m.master_name"
              class="master-stat-row"
            >
              <div class="master-stat-row__avatar">{{ m.master_name[0] }}</div>
              <div class="master-stat-row__body">
                <div class="master-stat-row__name">{{ m.master_name }}</div>
                <div class="master-stat-row__bar-wrap">
                  <div
                    class="master-stat-row__bar"
                    :style="`width: ${analytics.masters[0].revenue > 0
                      ? (m.revenue / analytics.masters[0].revenue * 100).toFixed(0)
                      : 0}%`"
                  />
                </div>
              </div>
              <div class="master-stat-row__stats">
                <span>{{ fmtPrice(m.revenue) }}</span>
                <span style="color:var(--text-3)">{{ m.completed }}/{{ m.total }}</span>
              </div>
            </div>
          </div>
          <p v-else style="color:var(--text-3);font-size:14px;padding:16px 0">{{ t('admin.no_data') }}</p>
        </div>

      </div>
    </div>

    <div v-else-if="activeTab === 'bookings'">
      <div class="filter-row">
        <select v-model="filterStatus" class="field-input field-select" style="width:auto">
          <option value="all">{{ t('admin.filter_all_statuses') }}</option>
          <option value="pending">{{ t('status.pending') }}</option>
          <option value="confirmed">{{ t('status.confirmed') }}</option>
          <option value="completed">{{ t('status.completed') }}</option>
          <option value="cancelled">{{ t('status.cancelled') }}</option>
        </select>
        <select v-model="filterMaster" class="field-input field-select" style="width:auto">
          <option value="all">{{ t('admin.filter_all_masters') }}</option>
          <option v-for="m in masters" :key="m.id" :value="m.id">{{ m.full_name }}</option>
        </select>
        <span style="font-size:13px;color:var(--text-3);margin-left:auto">
          {{ filteredBookings.length }} {{ t('admin.bookings_count') }}
        </span>
      </div>

      <div class="bookings-table card" style="padding: 0; overflow: hidden;">
        <table style="width:100%;border-collapse:collapse">
          <thead>
            <tr style="background:var(--bg-muted)">
              <th class="tbl-th">{{ t('admin.table_date') }}</th>
              <th class="tbl-th">{{ t('admin.table_service') }}</th>
              <th class="tbl-th">{{ t('admin.table_master') }}</th>
              <th class="tbl-th">{{ t('admin.table_amount') }}</th>
              <th class="tbl-th">{{ t('admin.table_status') }}</th>
              <th class="tbl-th">{{ t('admin.table_actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredBookings.length === 0">
              <td colspan="6" style="text-align:center;padding:40px;color:var(--text-3)">{{ t('admin.no_bookings') }}</td>
            </tr>
            <tr v-for="b in filteredBookings" :key="b.id" class="tbl-row">
              <td class="tbl-td">{{ formatDate(b.starts_at) }}</td>
              <td class="tbl-td" style="font-weight:500">{{ serviceName(b.service_id) }}</td>
              <td class="tbl-td">{{ masterName(b.master_id) }}</td>
              <td class="tbl-td" style="font-weight:600">{{ b.price_paid.toLocaleString() }} {{ t('common.currency') }}</td>
              <td class="tbl-td"><AppBadge :status="b.status" /></td>
              <td class="tbl-td">
                <div style="display:flex;gap:6px">
                  <button
                    v-if="b.status === 'pending'"
                    class="btn btn-sm"
                    style="background:var(--success-bg);color:var(--success);height:30px;padding:0 10px;font-size:12px"
                    @click="updateBookingStatus(b.id, 'confirmed')"
                  >✓</button>
                  <button
                    v-if="b.status === 'confirmed'"
                    class="btn btn-sm"
                    style="background:var(--info-bg,#EBF3FB);color:var(--info);height:30px;padding:0 10px;font-size:12px"
                    @click="updateBookingStatus(b.id, 'completed')"
                  >★</button>
                  <button
                    v-if="['pending','confirmed'].includes(b.status)"
                    class="btn btn-danger btn-sm"
                    style="height:30px;padding:0 10px;font-size:12px"
                    @click="updateBookingStatus(b.id, 'cancelled')"
                  >✕</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-else-if="activeTab === 'masters'" class="mgmt-list">
    <div style="margin-bottom:20px">
      <button class="btn btn-primary btn-sm" @click="showAddMaster = true">
        + Добавить мастера
      </button>
    </div>

    <div v-for="m in masters" :key="m.id" class="mgmt-row card">
      <div class="mgmt-row__avatar">{{ m.full_name[0] }}</div>
      <div class="mgmt-row__body">
        <p class="mgmt-row__name">{{ m.full_name }}</p>
        <p class="mgmt-row__sub">{{ m.experience_years }} лет · {{ m.instagram }}</p>
      </div>
      <div style="display:flex;gap:8px">
        <button
          class="btn btn-sm btn-outline"
          @click="toggleMaster(m.id, m.is_active)"
        >{{ m.is_active ? 'Деактивировать' : 'Активировать' }}</button>
        <button
          class="btn btn-sm btn-danger"
          @click="deleteMaster(m.id)"
        >Удалить</button>
      </div>
    </div>
    </div>

    <div v-else-if="activeTab === 'services'" class="mgmt-list">
    
    <div style="margin-bottom:20px">
      <button class="btn btn-primary btn-sm" @click="showAddService = true">
        + Добавить услугу
      </button>
    </div>

    <div v-for="s in services" :key="s.id" class="mgmt-row card">
      <div class="mgmt-row__body">
        <p class="mgmt-row__name">{{ s.name }}</p>
        <p class="mgmt-row__sub">{{ s.duration_min }} мин · {{ s.price.toLocaleString() }} ₸</p>
        <p style="font-size:13px;color:var(--text-3);margin-top:2px">{{ s.description }}</p>
      </div>
      <div style="display:flex;gap:8px;align-items:center">
        <span class="badge" :class="s.is_active ? 'badge-confirmed' : 'badge-cancelled'">
          {{ s.is_active ? 'Активна' : 'Скрыта' }}
        </span>
        <button
          class="btn btn-sm btn-outline"
          @click="toggleService(s.id, s.is_active)"
        >{{ s.is_active ? 'Скрыть' : 'Показать' }}</button>
        <button
          class="btn btn-sm btn-danger"
          @click="deleteService(s.id)"
        >Удалить</button>
      </div>
    </div>
    </div>

    <div v-else-if="activeTab === 'services'" class="mgmt-list">
      <div v-for="s in services" :key="s.id" class="mgmt-row card">
        <div class="mgmt-row__body">
          <p class="mgmt-row__name">{{ s.name }}</p>
          <p class="mgmt-row__sub">{{ s.duration_min }} {{ t('common.min') }} · {{ s.price.toLocaleString() }} {{ t('common.currency') }}</p>
          <p style="font-size:13px;color:var(--text-3);margin-top:2px">{{ s.description }}</p>
        </div>
        <button
          class="btn btn-sm"
          :class="s.is_active ? 'btn-outline' : 'btn-primary'"
          @click="toggleService(s.id, s.is_active)"
        >{{ s.is_active ? t('admin.hide') : t('admin.activate') }}</button>
      </div>
    </div>

  </div>

  <Teleport to="body">
    <Transition name="fade">
      <div v-if="showAddService" class="modal-overlay" @click.self="showAddService = false">
        <div class="modal-box card">
          <div class="modal-header">
            <h3 class="t-h3">Добавить услугу</h3>
            <button class="btn btn-ghost btn-icon" @click="showAddService = false">✕</button>
          </div>
          <div style="display:flex;flex-direction:column;gap:16px;margin-top:8px">
            <div class="field">
              <label class="field-label">Название</label>
              <input v-model="newService.name" type="text" class="field-input" />
            </div>
            <div class="field">
              <label class="field-label">Категория</label>
              <select v-model="newService.category" class="field-input field-select">
                <option value="hair">Волосы</option>
                <option value="nails">Ногти</option>
                <option value="face">Лицо</option>
                <option value="combo">Комплекс</option>
                <option value="beard">Борода</option>
                <option value="care">Уход</option>
              </select>
            </div>
            <div class="field">
              <label class="field-label">Описание</label>
              <textarea v-model="newService.description" class="field-input" rows="2" style="resize:none"/>
            </div>
            <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px">
              <div class="field">
                <label class="field-label">Длительность (мин)</label>
                <input v-model.number="newService.duration_min" type="number" class="field-input" />
              </div>
              <div class="field">
                <label class="field-label">Цена (₸)</label>
                <input v-model.number="newService.price" type="number" class="field-input" />
              </div>
            </div>
            <div style="display:flex;gap:10px;justify-content:flex-end;padding-top:8px;border-top:1px solid var(--border)">
              <button class="btn btn-outline" @click="showAddService = false">Отмена</button>
              <button class="btn btn-primary" :disabled="!newService.name" @click="createService">
                Создать →
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <Teleport to="body">
    <Transition name="fade">
      <div v-if="showAddMaster" class="modal-overlay" @click.self="showAddMaster = false">
        <div class="modal-box card">
          <div class="modal-header">
            <h3 class="t-h3">Добавить мастера<span class="red-dot">.</span></h3>
            <button class="btn btn-ghost btn-icon" @click="showAddMaster = false">✕</button>
          </div>

          <div style="display:flex;flex-direction:column;gap:16px;margin-top:8px">
            <div class="field">
              <label class="field-label">Email *</label>
              <input v-model="newMaster.email" type="email" class="field-input"
                placeholder="master@example.com" />
            </div>
            <div class="field">
              <label class="field-label">Пароль *</label>
              <input v-model="newMaster.password" type="password" class="field-input"
                placeholder="Минимум 6 символов" />
            </div>
            <div class="field">
              <label class="field-label">Имя и фамилия *</label>
              <input v-model="newMaster.full_name" type="text" class="field-input"
                placeholder="Имя Фамалия" />
            </div>
            <div class="field">
              <label class="field-label">О мастере</label>
              <textarea v-model="newMaster.bio" class="field-input" rows="2"
                style="resize:none" placeholder="Специализация, опыт..." />
            </div>
            <div style="display:grid;grid-template-columns:1fr 1fr;gap:12px">
              <div class="field">
                <label class="field-label">Опыт (лет)</label>
                <input v-model.number="newMaster.experience_years" type="number"
                  class="field-input" min="0" />
              </div>
              <div class="field">
                <label class="field-label">Instagram</label>
                <input v-model="newMaster.instagram" type="text" class="field-input"
                  placeholder="@username" />
              </div>
            </div>

            <div style="background:var(--info-bg);padding:12px 14px;border-radius:var(--radius-sm);font-size:13px;color:var(--info)">
              ℹ️ После создания мастер сможет войти через страницу входа с указанным email и паролем.
            </div>

            <div style="display:flex;gap:10px;justify-content:flex-end;padding-top:8px;border-top:1px solid var(--border)">
              <button class="btn btn-outline" @click="showAddMaster = false">Отмена</button>
              <button
                class="btn btn-primary"
                :disabled="!newMaster.email || !newMaster.full_name || !newMaster.password || savingMaster"
                @click="createMaster"
              >
                <span v-if="savingMaster" class="btn__spinner" />
                Создать мастера →
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.admin-header {
  display: flex; justify-content: space-between;
  align-items: flex-end; margin-bottom: 28px;
}

.tabs {
  display: flex; gap: 4px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 28px; overflow-x: auto;
}
.tab-btn {
  padding: 10px 20px; border: none; background: transparent;
  font-size: 14px; font-weight: 500; color: var(--text-2);
  border-bottom: 2px solid transparent; margin-bottom: -1px;
  cursor: pointer; transition: all 0.15s; white-space: nowrap;
}
.tab-btn--active { color: var(--accent); border-bottom-color: var(--accent); }
.tab-btn:hover:not(.tab-btn--active) { color: var(--text); }

.range-row {
  display: flex; gap: 6px; margin-bottom: 20px; flex-wrap: wrap;
}
.range-btn {
  padding: 7px 18px; border: 1.5px solid var(--border);
  border-radius: var(--radius-full); background: transparent;
  font-size: 13px; font-weight: 500; color: var(--text-2);
  cursor: pointer; transition: all 0.15s;
}
.range-btn--active {
  border-color: var(--accent); color: var(--accent);
  background: var(--accent-bg);
}

.kpi-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 12px; margin-bottom: 24px;
}
.kpi-card {
  background: var(--bg-card); border: 1px solid var(--border);
  border-radius: var(--radius-lg); padding: 16px 20px;
}
.kpi-card__label { font-size: 12px; color: var(--text-3); font-weight: 500;
                   letter-spacing: 0.04em; text-transform: uppercase; margin-bottom: 8px; }
.kpi-card__val   { font-size: 24px; font-weight: 700; color: var(--text); }

.charts-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
@media (max-width: 768px) { .charts-grid { grid-template-columns: 1fr; } }

.chart-card__header {
  display: flex; justify-content: space-between;
  align-items: center; margin-bottom: 16px;
}
.chart-card__title { font-size: 15px; font-weight: 600; color: var(--text); }
.chart-card__body  { height: 240px; position: relative; }

.master-stat-row {
  display: flex; align-items: center; gap: 12px;
  padding: 10px 0; border-bottom: 1px solid var(--border);
}
.master-stat-row:last-child { border-bottom: none; }
.master-stat-row__avatar {
  width: 36px; height: 36px; border-radius: 50%;
  background: var(--accent-bg); color: var(--accent);
  display: flex; align-items: center; justify-content: center;
  font-weight: 700; font-size: 15px; flex-shrink: 0;
}
.master-stat-row__body { flex: 1; min-width: 0; }
.master-stat-row__name { font-size: 14px; font-weight: 500; color: var(--text); }
.master-stat-row__bar-wrap {
  height: 4px; background: var(--bg-muted);
  border-radius: 99px; margin-top: 4px; overflow: hidden;
}
.master-stat-row__bar {
  height: 100%; background: var(--accent);
  border-radius: 99px; transition: width 0.4s;
}
.master-stat-row__stats {
  display: flex; flex-direction: column; align-items: flex-end;
  gap: 2px; font-size: 13px; font-weight: 600; color: var(--text);
}

.tbl-th {
  text-align: left; padding: 12px 16px;
  font-size: 12px; font-weight: 600; color: var(--text-3);
  letter-spacing: 0.04em; text-transform: uppercase;
  border-bottom: 1px solid var(--border);
}
.tbl-td { padding: 12px 16px; border-bottom: 1px solid var(--border); font-size: 14px; }
.tbl-row:last-child .tbl-td { border-bottom: none; }
.tbl-row:hover .tbl-td { background: var(--bg-muted); }

.filter-row {
  display: flex; gap: 10px; align-items: center;
  margin-bottom: 16px; flex-wrap: wrap;
}

.mgmt-list { display: flex; flex-direction: column; gap: 10px; }
.mgmt-row  { display: flex; align-items: center; gap: 16px; }
.mgmt-row__avatar {
  width: 44px; height: 44px; border-radius: 50%;
  background: var(--accent-bg); color: var(--accent);
  display: flex; align-items: center; justify-content: center;
  font-size: 18px; font-weight: 700; flex-shrink: 0;
}
.mgmt-row__body  { flex: 1; min-width: 0; }
.mgmt-row__name  { font-size: 15px; font-weight: 600; color: var(--text); }
.mgmt-row__sub   { font-size: 13px; color: var(--text-2); margin-top: 2px; }
.mgmt-row__tags  { display: flex; flex-wrap: wrap; gap: 4px; margin-top: 6px; }

.loading-state { text-align: center; padding: 80px; color: var(--text-3); font-size: 15px; }
</style>