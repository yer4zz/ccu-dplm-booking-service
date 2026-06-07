<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n }       from 'vue-i18n'
import { galleryApi }    from '@/api/gallery'
import { mastersApi }    from '@/api/masters'
import { servicesApi }   from '@/api/services'
import type { Master, Service } from '@/types'
import { ArrowRight, X } from 'lucide-vue-next'

const { t } = useI18n()

const raboty   = ref<any[]>([])
const mastera  = ref<Master[]>([])
const uslugi   = ref<Service[]>([])
const zagruzka = ref(true)
const vybrana  = ref<any | null>(null)

const filterMaster  = ref('all')
const filterService = ref('all')

const filtersOpen = ref(false)

onMounted(async () => {
  const [r, m, s] = await Promise.all([
    galleryApi.list(),
    mastersApi.list(),
    servicesApi.list(),
  ])
  raboty.value   = r ?? []
  mastera.value  = m ?? []
  uslugi.value   = s ?? []
  zagruzka.value = false
})

const filteredRaboty = computed(() => {
  return raboty.value.filter(w => {
    const masterOk  = filterMaster.value  === 'all' || w.master_id  === filterMaster.value
    const serviceOk = filterService.value === 'all' || w.service_id === filterService.value
    return masterOk && serviceOk
  })
})

function setFilter(type: 'master' | 'service', id: string) {
  if (type === 'master') { filterMaster.value = id; filterService.value = 'all' }
  else                   { filterService.value = id; filterMaster.value = 'all' }
  filtersOpen.value = false
}

function resetFilters() {
  filterMaster.value  = 'all'
  filterService.value = 'all'
  filtersOpen.value   = false
}

const activeFilterLabel = computed(() => {
  if (filterMaster.value !== 'all') {
    return mastera.value.find(m => m.id === filterMaster.value)?.full_name ?? t('gallery.filter_all')
  }
  if (filterService.value !== 'all') {
    return uslugi.value.find(s => s.id === filterService.value)?.name ?? t('gallery.filter_all')
  }
  return null
})

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'long', year: 'numeric'
  })
}
</script>

<template>
  <div class="gallery-page">
    <div class="page-container">

      <div class="gallery-head">
        <div class="gallery-head__text">
          <div class="section-label">{{ t('gallery.eyebrow') }} · РАБОТЫ МАСТЕРОВ</div>
          <h1 class="t-h1">
            {{ t('gallery.title') }}<span class="red-dot">.</span>
          </h1>
          <p class="gallery-head__desc">
            Избранные работы наших мастеров. Каждая работа — отражение профессионализма и внимания к деталям.
          </p>
        </div>
        <div class="gallery-head__count">
          <span class="gallery-head__num">{{ filteredRaboty.length }}</span>
          <span class="gallery-head__label">{{ t('gallery.works_count') }}</span>
        </div>
      </div>

      <div class="deco-line" />

      <div class="gallery-filters gallery-filters--desktop">
        <button
          class="filter-btn"
          :class="{ 'filter-btn--active': filterMaster === 'all' && filterService === 'all' }"
          @click="resetFilters"
        >{{ t('gallery.all') }}</button>
        <button
          v-for="m in mastera"
          :key="m.id"
          class="filter-btn"
          :class="{ 'filter-btn--active': filterMaster === m.id }"
          @click="setFilter('master', m.id)"
        >{{ m.full_name }}</button>
        <div class="filter-sep" />
        <button
          v-for="s in uslugi.slice(0, 5)"
          :key="s.id"
          class="filter-btn filter-btn--sm"
          :class="{ 'filter-btn--active': filterService === s.id }"
          @click="setFilter('service', s.id)"
        >{{ s.name }}</button>
      </div>

      <div class="gallery-filters-mobile">
        <button
          class="filter-toggle"
          :class="{ 'filter-toggle--active': activeFilterLabel }"
          @click="filtersOpen = !filtersOpen"
        >
          <span>{{ activeFilterLabel ?? t('gallery.all') }}</span>
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none"
            :style="{ transform: filtersOpen ? 'rotate(180deg)' : 'rotate(0)' }"
            style="transition: transform 0.2s; flex-shrink:0"
          >
            <path d="M2 4l4 4 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>

        <button v-if="activeFilterLabel" class="filter-clear" @click="resetFilters">
          <X :size="20" /> Все
        </button>

        <Transition name="fade">
          <div v-if="filtersOpen" class="filter-dropdown">
            <p class="filter-dropdown__label">{{ t('gallery.filter_masters') }}</p>
            <button
              class="filter-dropdown__item"
              :class="{ 'filter-dropdown__item--active': filterMaster === 'all' && filterService === 'all' }"
              @click="resetFilters"
            >{{ t('gallery.all') }}</button>
            <button
              v-for="m in mastera"
              :key="m.id"
              class="filter-dropdown__item"
              :class="{ 'filter-dropdown__item--active': filterMaster === m.id }"
              @click="setFilter('master', m.id)"
            >{{ m.full_name }}</button>

            <div class="filter-dropdown__sep" />

            <p class="filter-dropdown__label">{{ t('gallery.filter_services') }}</p>
            <button
              v-for="s in uslugi.slice(0, 5)"
              :key="s.id"
              class="filter-dropdown__item"
              :class="{ 'filter-dropdown__item--active': filterService === s.id }"
              @click="setFilter('service', s.id)"
            >{{ s.name }}</button>
          </div>
        </Transition>
      </div>

      <div class="deco-line" />

      <div v-if="zagruzka" class="gallery-state">
        <div class="gallery-spinner" />
        {{ t('common.loading') }}
      </div>

      <div v-else-if="filteredRaboty.length === 0" class="gallery-state">
        <div class="section-label" style="margin-bottom:16px">{{ t('gallery.empty_label') }}</div>
        <p style="font-size:22px;font-family:var(--font-serif);color:var(--text)">{{ t('gallery.empty') }}</p>
      </div>

      <div v-else class="gallery-grid-blog">

        <div v-if="filteredRaboty.length >= 1" class="gallery-row gallery-row--featured">
          <div class="gallery-card gallery-card--lg" @click="vybrana = filteredRaboty[0]">
            <div class="gallery-card__img-wrap">
              <img :src="filteredRaboty[0].image_url" :alt="filteredRaboty[0].title || 'Work'" loading="lazy" />
              <div class="gallery-card__overlay"><span class="gallery-card__zoom"><ArrowRight :size="14"/></span></div>
              <span class="gallery-card__no">Nº 01</span>
            </div>
            <div class="gallery-card__info">
              <p v-if="filteredRaboty[0].title" class="gallery-card__title">{{ filteredRaboty[0].title }}</p>
              <div class="gallery-card__meta">
                <span class="gallery-card__master">{{ filteredRaboty[0].master_name }}</span>
                <span class="gallery-card__dot">·</span>
                <span class="gallery-card__service">{{ filteredRaboty[0].service_name }}</span>
                <span class="gallery-card__dot gallery-card__dot--hide">·</span>
                <span class="gallery-card__date gallery-card__date--hide">{{ formatDate(filteredRaboty[0].created_at) }}</span>
              </div>
            </div>
          </div>

          <div class="gallery-sidebar" v-if="filteredRaboty.length >= 3">
            <div
              v-for="(w, i) in filteredRaboty.slice(1, 3)"
              :key="w.id"
              class="gallery-card gallery-card--md"
              @click="vybrana = w"
            >
              <div class="gallery-card__img-wrap">
                <img :src="w.image_url" :alt="w.title || 'Work'" loading="lazy" />
                <div class="gallery-card__overlay"><span class="gallery-card__zoom"><ArrowRight :size="14"/></span></div>
                <span class="gallery-card__no">Nº {{ String(i+2).padStart(2,'0') }}</span>
              </div>
              <div class="gallery-card__info gallery-card__info--sm">
                <span class="gallery-card__master">{{ w.master_name }}</span>
                <span class="gallery-card__dot"> · </span>
                <span class="gallery-card__service">{{ w.service_name }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="deco-line" v-if="filteredRaboty.length > 3" />

        <div class="gallery-grid-3" v-if="filteredRaboty.length > 3">
          <div
            v-for="(w, i) in filteredRaboty.slice(3)"
            :key="w.id"
            class="gallery-card"
            @click="vybrana = w"
          >
            <div class="gallery-card__img-wrap">
              <img :src="w.image_url" :alt="w.title || 'Work'" loading="lazy" />
              <div class="gallery-card__overlay"><span class="gallery-card__zoom"><ArrowRight :size="14"/></span></div>
              <span class="gallery-card__no">Nº {{ String(i+4).padStart(2,'0') }}</span>
            </div>
            <div class="gallery-card__info">
              <p v-if="w.title" class="gallery-card__title">{{ w.title }}</p>
              <div class="gallery-card__meta">
                <span class="gallery-card__master">{{ w.master_name }}</span>
                <span class="gallery-card__dot">·</span>
                <span class="gallery-card__service">{{ w.service_name }}</span>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>

    <Teleport to="body">
      <Transition name="fade">
        <div v-if="vybrana" class="lightbox" @click.self="vybrana = null">
          <div class="lightbox__box">
            <button class="lightbox__close btn btn-ghost btn-sm" @click="vybrana = null">
              <X :size="20" /> Закрыть
            </button>
            <img :src="vybrana.image_url" :alt="vybrana.title || 'Work'" class="lightbox__img" />
            <div class="lightbox__info">
              <div class="section-label" style="margin-bottom:8px">{{ t('gallery.by') }}</div>
              <h3 class="lightbox__title" v-if="vybrana.title">{{ vybrana.title }}</h3>
              <div class="lightbox__meta">
                <span>{{ vybrana.master_name }}</span>
                <span class="lightbox__dot">·</span>
                <span>{{ vybrana.service_name }}</span>
                <span class="lightbox__dot">·</span>
                <span>{{ formatDate(vybrana.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

  </div>
</template>

<style scoped>
.gallery-page { min-height: calc(100vh - 56px); padding: 60px 0 100px; }

.gallery-head {
  display: flex; justify-content: space-between;
  align-items: flex-end; padding-bottom: 40px;
}
.gallery-head__desc {
  font-size: 15px; color: var(--text-2);
  margin-top: 12px; max-width: 480px; line-height: 1.7;
}
.gallery-head__count { display: flex; flex-direction: column; align-items: flex-end; flex-shrink: 0; }
.gallery-head__num {
  font-family: var(--font-serif); font-size: 56px; font-weight: 900;
  color: var(--text); opacity: 0.12; line-height: 1;
}
.gallery-head__label {
  font-size: 10px; font-weight: 500; letter-spacing: 0.1em;
  text-transform: uppercase; color: var(--text-3);
}

.gallery-filters--desktop {
  display: flex; align-items: center; gap: 6px;
  flex-wrap: wrap; padding: 20px 0;
}
.filter-btn {
  font-family: var(--font); font-size: 12px; font-weight: 500;
  letter-spacing: 0.06em; padding: 6px 14px;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm); color: var(--text-3);
  cursor: pointer; transition: all 0.15s;
}
.filter-btn:hover    { border-color: var(--border-strong); color: var(--text); }
.filter-btn--active  { border-color: var(--accent); color: var(--accent); background: var(--accent-bg); }
.filter-btn--sm      { font-size: 11px; padding: 5px 10px; }
.filter-sep          { width: 1px; height: 20px; background: var(--border); margin: 0 4px; }

.gallery-filters-mobile { display: none; padding: 16px 0; position: relative; }
.filter-toggle {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 16px; height: 38px;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-family: var(--font); font-size: 13px; font-weight: 500;
  color: var(--text-2); cursor: pointer; transition: all 0.15s;
}
.filter-toggle--active { border-color: var(--accent); color: var(--accent); background: var(--accent-bg); }
.filter-clear {
  margin-left: 8px; padding: 8px 14px; height: 38px;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-family: var(--font); font-size: 12px; font-weight: 500;
  color: var(--text-3); cursor: pointer; transition: all 0.15s;
}
.filter-clear:hover { border-color: var(--danger); color: var(--danger); }

.filter-dropdown {
  position: absolute; top: calc(100% - 8px); left: 0;
  z-index: 50; min-width: 220px;
  background: var(--bg-card); border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 8px 0;
  box-shadow: var(--shadow-lg);
}
.filter-dropdown__label {
  font-size: 9px; font-weight: 600; letter-spacing: 0.12em;
  text-transform: uppercase; color: var(--text-3);
  padding: 8px 16px 4px;
}
.filter-dropdown__item {
  display: block; width: 100%; text-align: left;
  padding: 10px 16px;
  background: none; border: none;
  font-family: var(--font); font-size: 14px;
  color: var(--text-2); cursor: pointer;
  transition: background 0.1s;
}
.filter-dropdown__item:hover          { background: var(--bg-hover); color: var(--text); }
.filter-dropdown__item--active        { color: var(--accent); font-weight: 500; }
.filter-dropdown__sep { height: 1px; background: var(--border); margin: 6px 0; }

.gallery-state {
  display: flex; align-items: center; justify-content: center;
  flex-direction: column; gap: 16px; padding: 80px 0; color: var(--text-2);
}
.gallery-spinner {
  width: 20px; height: 20px;
  border: 1.5px solid var(--border-strong);
  border-top-color: var(--accent);
  border-radius: 50%; animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.gallery-grid-blog { margin-top: 0; }
.gallery-row--featured {
  display: grid; grid-template-columns: 1fr 400px; gap: 16px;
}
.gallery-sidebar { display: flex; flex-direction: column; gap: 16px; }

.gallery-card { cursor: pointer; }
.gallery-card__img-wrap {
  position: relative; overflow: hidden;
  border-radius: var(--radius-sm); background: var(--bg-muted);
}
.gallery-card--lg .gallery-card__img-wrap  { aspect-ratio: 4/3; }
.gallery-card--md .gallery-card__img-wrap  { aspect-ratio: 16/9; }
.gallery-card:not(.gallery-card--lg):not(.gallery-card--md) .gallery-card__img-wrap {
  aspect-ratio: 4/3;
}
.gallery-card__img-wrap img {
  width: 100%; height: 100%; object-fit: cover;
  transition: transform 0.4s cubic-bezier(0.4,0,0.2,1);
}
.gallery-card:hover .gallery-card__img-wrap img { transform: scale(1.04); }
.gallery-card__overlay {
  position: absolute; inset: 0;
  background: rgba(26,24,20,0.35);
  display: flex; align-items: center; justify-content: center;
  opacity: 0; transition: opacity 0.2s;
}
.gallery-card:hover .gallery-card__overlay { opacity: 1; }
.gallery-card__zoom  { font-size: 28px; color: #fff; font-weight: 200; }
.gallery-card__no {
  position: absolute; top: 12px; left: 12px;
  font-size: 9px; font-weight: 500; letter-spacing: 0.14em;
  text-transform: uppercase; color: rgba(255,255,255,0.7);
  background: rgba(26,24,20,0.45); padding: 3px 8px; border-radius: 1px;
}
.gallery-card__info    { padding: 12px 0 0; }
.gallery-card__info--sm { padding: 8px 0 0; }
.gallery-card__title {
  font-family: var(--font-serif); font-size: 18px; font-weight: 700;
  color: var(--text); margin-bottom: 6px; line-height: 1.3;
}
.gallery-card__meta {
  display: flex; align-items: center; gap: 6px; flex-wrap: wrap;
}
.gallery-card__master { font-size: 12px; font-weight: 600; color: var(--text); }
.gallery-card__service { font-size: 12px; color: var(--text-3); }
.gallery-card__date    { font-size: 11px; color: var(--text-3); }
.gallery-card__dot     { color: var(--border-strong); font-size: 10px; }

.gallery-grid-3 {
  display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px 20px;
}

.lightbox {
  position: fixed; inset: 0; z-index: 1000;
  background: rgba(26,24,20,0.92);
  backdrop-filter: blur(6px);
  display: flex; align-items: center; justify-content: center;
  padding: 20px;
}
.lightbox__box {
  position: relative; max-width: 900px; width: 100%;
  background: var(--bg-card); border-radius: var(--radius-sm); overflow: hidden;
}
.lightbox__close {
  position: absolute; top: 14px; right: 14px; z-index: 10;
  background: rgba(26,24,20,0.5) !important;
  color: #fff !important; border-color: transparent !important;
  font-size: 11px; letter-spacing: 0.06em;
}
.lightbox__img {
  width: 100%; max-height: 65vh;
  object-fit: contain; display: block; background: var(--bg-muted);
}
.lightbox__info  { padding: 20px 24px; border-top: 1px solid var(--border); }
.lightbox__title {
  font-family: var(--font-serif); font-size: 22px; font-weight: 700;
  color: var(--text); margin-bottom: 8px;
}
.lightbox__meta {
  display: flex; align-items: center; gap: 8px;
  font-size: 13px; color: var(--text-2); flex-wrap: wrap;
}
.lightbox__dot { color: var(--border-strong); }

@media (max-width: 900px) {
  .gallery-row--featured { grid-template-columns: 1fr; }
  .gallery-sidebar       { flex-direction: row; }
  .gallery-sidebar .gallery-card--md .gallery-card__img-wrap { aspect-ratio: 1/1; }
  .gallery-grid-3        { grid-template-columns: repeat(2, 1fr); }
  .gallery-head          { flex-direction: column; align-items: flex-start; gap: 16px; }
  .gallery-head__count   { align-items: flex-start; }
  .gallery-head__desc    { display: none; }
}

@media (max-width: 600px) {
  .gallery-page          { padding: 32px 0 60px; }

  .gallery-head          { padding-bottom: 20px; }
  .gallery-head__num     { font-size: 40px; }

  .gallery-filters--desktop { display: none; }
  .gallery-filters-mobile   { display: flex; align-items: center; }

  .gallery-row--featured { grid-template-columns: 1fr; }
  .gallery-sidebar       { flex-direction: column; }
  .gallery-sidebar .gallery-card--md .gallery-card__img-wrap { aspect-ratio: 16/9; }
  .gallery-grid-3        { grid-template-columns: 1fr; gap: 16px; }

  .gallery-card__dot--hide  { display: none; }
  .gallery-card__date--hide { display: none; }

  .lightbox          { padding: 0; align-items: flex-end; }
  .lightbox__box     { border-radius: var(--radius-sm) var(--radius-sm) 0 0; max-height: 95vh; overflow-y: auto; }
  .lightbox__img     { max-height: 55vh; }
  .lightbox__info    { padding: 16px; }
  .lightbox__title   { font-size: 18px; }
  .lightbox__meta    { font-size: 12px; gap: 6px; }
}
</style>