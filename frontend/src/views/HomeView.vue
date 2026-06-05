<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n }               from 'vue-i18n'
import { servicesApi }           from '@/api/services'
import { mastersApi }            from '@/api/masters'
import { galleryApi }            from '@/api/gallery'
import { useBookingStore }       from '@/stores/booking'
import type { Service, Master }  from '@/types'

const { t, tm, locale } = useI18n()
const router    = useRouter()
const store     = useBookingStore()

// Локализованное название услуги
function svcName(svc: any): string {
  if (locale.value === 'kz' && svc.name_kz) return svc.name_kz
  if (locale.value === 'en' && svc.name_en) return svc.name_en
  return svc.name
}

const uslugi       = ref<Service[]>([])
const mastera      = ref<Master[]>([])
const galereya     = ref<any[]>([])
const publicOtzivi = ref<any[]>([])
const zagruzka     = ref(true)

const howSteps = computed(
  () => tm('home.how_steps') as { num: string; title: string; desc: string }[],
)

onMounted(async () => {
  const [s, m, g, r] = await Promise.all([
    servicesApi.list(),
    mastersApi.list(),
    galleryApi.list().catch(() => []),
    fetch('/api/v1/reviews/public').then(r => r.json()).catch(() => []),
  ])
  uslugi.value       = s ?? []
  mastera.value      = m ?? []
  galereya.value     = g ?? []
  publicOtzivi.value = r ?? []
  zagruzka.value     = false
  startMasterSlider()
  startReviewAutoScroll()
})

onUnmounted(() => {
  stopMasterSlider()
})

// ── СЛАЙДЕР МАСТЕРОВ ─────────────────────────────
const currentMaster = ref(0)
let sliderTimer: ReturnType<typeof setInterval> | null = null
const sliderPaused = ref(false)

function startMasterSlider() {
  sliderTimer = setInterval(() => {
    if (!sliderPaused && mastera.value.length > 1) {
      currentMaster.value = (currentMaster.value + 1) % mastera.value.length
    }
  }, 7000)
}
function stopMasterSlider() {
  if (sliderTimer) clearInterval(sliderTimer)
}
function prevMaster() {
  currentMaster.value = (currentMaster.value - 1 + mastera.value.length) % mastera.value.length
}
function nextMaster() {
  currentMaster.value = (currentMaster.value + 1) % mastera.value.length
}
function goMaster(i: number) { currentMaster.value = i }

// ── КАРУСЕЛЬ ОТЗЫВОВ ─────────────────────────────
const reviewsEl = ref<HTMLElement | null>(null)
const reviewIdx = ref(0)

function startReviewAutoScroll() {}

function scrollReview(dir: 'prev' | 'next') {
  const max = publicOtzivi.value.length - 1
  if (dir === 'next') reviewIdx.value = Math.min(reviewIdx.value + 1, max)
  else               reviewIdx.value = Math.max(reviewIdx.value - 1, 0)
  const el = reviewsEl.value
  if (el) {
    const card = el.querySelector('.review-card') as HTMLElement
    const w = card ? card.offsetWidth + 24 : 320
    el.scrollTo({ left: reviewIdx.value * w, behavior: 'smooth' })
  }
}

// ── УТИЛИТЫ ──────────────────────────────────────
function bookService(svc: Service) {
  store.reset(); store.selectedServices = [svc]; router.push('/book')
}
function bookMaster(m: Master) {
  store.reset(); store.selectedMaster = m; router.push('/book')
}
function getCatLabel(k: string) { return t(`category.${k}`) }
const catColor: Record<string, string> = {
  hair:'#8B5CF6', nails:'#E8442A', face:'#F59E0B',
  combo:'#C9956A', beard:'#3B82F6', care:'#10B981',
}
function getCatColor(k: string) { return catColor[k] || '#8A8680' }
</script>

<template>
  <div class="home">

    <!-- ═══ HERO — FULLSCREEN ════════════════════════ -->
    <section class="hero-full">
      <div class="hero-full__bg">
        <img src="/images/мэйн.jpg" alt="Beauty Dana" class="hero-full__img" />
        <div class="hero-full__overlay" />
      </div>
      <div class="page-container hero-full__inner">
        <div class="hero-full__eyebrow">
          <span class="hero-full__dash">—</span>
          {{ t('home.hero_label') }} · Nº 01
        </div>
        <h1 class="hero-full__title">
          {{ t('home.hero_line1') }}<br/>
          <em>{{ t('home.hero_line2') }}</em><br/>
          {{ t('home.hero_line3') }}<span class="red-dot">.</span>
        </h1>
        <p class="hero-full__sub">{{ t('home.hero_sub') }}</p>
        <div class="hero-full__btns">
          <button class="btn btn-primary btn-lg" @click="router.push('/book')">
            {{ t('home.book_btn') }} →
          </button>
          <RouterLink to="/masters" class="btn btn-outline-light btn-lg">
            {{ t('home.masters_btn') }}
          </RouterLink>
        </div>
        <div class="hero-full__perks">
          <span>✓ {{ t('home.perk_discount') }}</span>
          <span class="hero-full__perks-dot">·</span>
          <span>✓ {{ t('home.perk_points') }}</span>
          <span class="hero-full__perks-dot">·</span>
          <span>✓ {{ t('home.perk_reminders') }}</span>
        </div>
      </div>
    </section>

    <!-- ═══ STATS BAR — отдельный блок ══════════════ -->
    <div class="stats-bar">
      <div class="stats-bar__inner">
        <div class="stats-bar__item">
          <span class="stats-bar__num">{{ uslugi.length || 8 }}</span>
          <span class="stats-bar__label">{{ t('home.stats_services') }}</span>
        </div>
        <div class="stats-bar__divider" />
        <div class="stats-bar__item">
          <span class="stats-bar__num">{{ mastera.length || 3 }}</span>
          <span class="stats-bar__label">{{ t('home.stats_masters') }}</span>
        </div>
        <div class="stats-bar__divider" />
        <div class="stats-bar__item">
          <span class="stats-bar__num">13</span>
          <span class="stats-bar__label">{{ t('home.stats_experience') }}</span>
        </div>
        <div class="stats-bar__divider" />
        <div class="stats-bar__item">
          <span class="stats-bar__num">24/7</span>
          <span class="stats-bar__label">{{ t('home.stats_online') }}</span>
        </div>
      </div>
    </div>

    <!-- ═══ УСЛУГИ ═══════════════════════════════════ -->
    <section class="home-section section-services" id="services">
      <div class="page-container">
        <div class="home-section__head">
          <div>
            <div class="section-label">{{ t('home.services_eyebrow') }} · Nº 02</div>
            <h2 class="home-section__title">
              {{ t('home.services_title') }}<span class="red-dot">.</span>
            </h2>
          </div>
          <button class="arrow-link" @click="router.push('/book')">
            Все услуги →
          </button>
        </div>
        <div class="services-grid">
          <div
            v-for="(svc, i) in uslugi"
            :key="svc.id"
            class="svc-card"
            :style="`--cat-color: ${getCatColor(svc.category)}`"
            @click="bookService(svc)"
          >
            <div class="svc-card__accent" />
            <div class="svc-card__body">
              <div class="svc-card__top">
                <span class="svc-card__num">{{ String(i+1).padStart(2,'0') }}</span>
                <span class="svc-card__cat">{{ getCatLabel(svc.category) }}</span>
              </div>
              <h3 class="svc-card__name">{{ svcName(svc) }}</h3>
              <div class="svc-card__footer">
                <span class="svc-card__dur">{{ svc.duration_min }} {{ t('common.min') }}</span>
                <span class="svc-card__price">{{ svc.price.toLocaleString() }} {{ t('common.currency') }}</span>
              </div>
            </div>
            <div class="svc-card__arrow">→</div>
          </div>
        </div>
      </div>
    </section>

    <!-- ═══ МАСТЕРА — СЛАЙДЕР ════════════════════════ -->
    <section
      class="home-section section-masters"
      id="masters"
      @mouseenter="sliderPaused = true"
      @mouseleave="sliderPaused = false"
    >
      <div class="page-container">
        <div class="home-section__head">
          <div>
            <div class="section-label">{{ t('home.masters_eyebrow') }} · Nº 03</div>
            <h2 class="home-section__title">
              {{ t('home.masters_title') }}<span class="red-dot">.</span>
            </h2>
          </div>
          <RouterLink to="/masters" class="arrow-link">
            Все мастера →
          </RouterLink>
        </div>

        <div v-if="mastera.length > 0" class="masters-slider">
          <!-- фото -->
          <div class="masters-slider__photo-wrap">
            <Transition name="master-fade" mode="out-in">
              <div :key="currentMaster" class="masters-slider__photo">
                <img
                  v-if="mastera[currentMaster]?.avatar_url"
                  :src="mastera[currentMaster].avatar_url"
                  :alt="mastera[currentMaster].full_name"
                />
                <div v-else class="masters-slider__initials">
                  {{ mastera[currentMaster]?.full_name[0] }}
                </div>
              </div>
            </Transition>
            <!-- нижний номер -->
            <span class="masters-slider__fig">
              FIG. {{ String(currentMaster + 1).padStart(2,'0') }} / МАСТЕР
            </span>
          </div>

          <!-- инфо -->
          <div class="masters-slider__body">
            <Transition name="master-fade" mode="out-in">
              <div :key="currentMaster" class="masters-slider__info">
                <p class="masters-slider__exp">
                  {{ mastera[currentMaster]?.experience_years }} {{ t('home.masters_years_exp') }}
                </p>
                <h3 class="masters-slider__name">
                  {{ mastera[currentMaster]?.full_name }}
                </h3>
                <p class="masters-slider__bio">{{ mastera[currentMaster]?.bio }}</p>
                <div class="masters-slider__tags">
                  <span
                    v-for="s in mastera[currentMaster]?.services?.slice(0, 4)"
                    :key="s.id"
                    class="masters-slider__tag"
                  >{{ s.name }}</span>
                </div>
                <button
                  class="btn btn-primary"
                  @click="bookMaster(mastera[currentMaster])"
                >
                  {{ t('home.masters_book_to') }} {{ mastera[currentMaster]?.full_name?.split(' ')[0] }} →
                </button>
              </div>
            </Transition>

            <!-- навигация -->
            <div class="masters-slider__nav">
              <button class="masters-slider__arrow" @click="prevMaster" aria-label="Назад">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                  <path d="M13 4l-6 6 6 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
              </button>
              <div class="masters-slider__dots">
                <button
                  v-for="(_, i) in mastera"
                  :key="i"
                  class="masters-slider__dot"
                  :class="{ 'masters-slider__dot--active': i === currentMaster }"
                  @click="goMaster(i)"
                  :aria-label="`Мастер ${i+1}`"
                />
              </div>
              <button class="masters-slider__arrow" @click="nextMaster" aria-label="Вперёд">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                  <path d="M7 4l6 6-6 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
              </button>
              <span class="masters-slider__counter">
                {{ String(currentMaster + 1).padStart(2,'0') }} / {{ String(mastera.length).padStart(2,'0') }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ═══ КАК РАБОТАЕТ ═════════════════════════════ -->
    <section class="home-section section-how">
      <div class="page-container">
        <div class="home-section__head">
          <div>
            <div class="section-label">METHOD · LOOP · 04 STAGES</div>
            <h2 class="home-section__title">
              {{ t('home.how_title') }}<span class="red-dot">.</span>
            </h2>
          </div>
        </div>
        <div class="how-steps">
          <div v-for="(step, i) in howSteps" :key="i" class="how-step">
            <div class="how-step__num">{{ step.num }}</div>
            <h4 class="how-step__title">{{ step.title }}</h4>
            <p class="how-step__desc">{{ step.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ═══ ГАЛЕРЕЯ ═══════════════════════════════════ -->
    <section class="home-section section-gallery" v-if="galereya.length > 0">
      <div class="page-container">
        <div class="home-section__head">
          <div>
            <div class="section-label">{{ t('home.gallery_eyebrow') }}</div>
            <h2 class="home-section__title">
              {{ t('gallery.title') }}<span class="red-dot">.</span>
            </h2>
          </div>
          <RouterLink to="/gallery" class="arrow-link">{{ t('home.gallery_open') }}</RouterLink>
        </div>
        <div class="gallery-blog">
          <RouterLink to="/gallery" class="gallery-blog__featured">
            <div class="gallery-blog__img-wrap">
              <img :src="galereya[0]?.image_url" :alt="galereya[0]?.title || 'Work'" loading="lazy" />
              <div class="gallery-blog__overlay"><span class="gallery-blog__arrow">→</span></div>
              <span class="gallery-blog__no">Nº 01</span>
            </div>
            <div class="gallery-blog__info">
              <p class="gallery-blog__master">{{ galereya[0]?.master_name }}</p>
              <p class="gallery-blog__service">{{ galereya[0]?.service_name }}</p>
            </div>
          </RouterLink>
          <div class="gallery-blog__grid">
            <RouterLink
              v-for="(w, i) in galereya.slice(1, 5)" :key="w.id"
              to="/gallery" class="gallery-blog__item"
            >
              <div class="gallery-blog__img-wrap">
                <img :src="w.image_url" :alt="w.title || 'Work'" loading="lazy" />
                <div class="gallery-blog__overlay"><span class="gallery-blog__arrow">→</span></div>
                <span class="gallery-blog__no">Nº {{ String(i+2).padStart(2,'0') }}</span>
              </div>
              <div class="gallery-blog__info">
                <p class="gallery-blog__master">{{ w.master_name }}</p>
                <p class="gallery-blog__service">{{ w.service_name }}</p>
              </div>
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- ═══ ОТЗЫВЫ — КАРУСЕЛЬ ════════════════════════ -->
    <section class="home-section section-reviews" v-if="publicOtzivi.length > 0">
      <div class="page-container">
        <div class="home-section__head">
          <div>
            <div class="section-label">{{ t('home.reviews_eyebrow') }}</div>
            <h2 class="home-section__title">
              {{ t('home.reviews_title') }}<span class="red-dot">.</span>
            </h2>
          </div>
          <!-- стрелки -->
          <div class="reviews-nav">
            <button
              class="reviews-nav__btn"
              :disabled="reviewIdx === 0"
              @click="scrollReview('prev')"
              aria-label="Назад"
            >
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                <path d="M10 3L5 8l5 5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              </svg>
            </button>
            <button
              class="reviews-nav__btn"
              :disabled="reviewIdx >= publicOtzivi.length - 1"
              @click="scrollReview('next')"
              aria-label="Вперёд"
            >
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                <path d="M6 3l5 5-5 5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              </svg>
            </button>
          </div>
        </div>
        <div class="reviews-track" ref="reviewsEl">
          <div
            v-for="rv in publicOtzivi"
            :key="rv.id"
            class="review-card"
          >
            <div class="review-card__quote">"</div>
            <p class="review-card__text">{{ rv.comment }}</p>
            <div class="review-card__footer">
              <span class="review-card__stars">
                {{ '★'.repeat(rv.rating) }}{{ '☆'.repeat(5 - rv.rating) }}
              </span>
              <span class="review-card__name">— {{ rv.client_name }}</span>
              <span class="review-card__service">{{ rv.service_name }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ═══ О НАС (краткая версия) ══════════════════ -->
    <section class="home-section section-about" id="about">
      <div class="page-container">
        <div class="home-section__head">
          <div>
            <div class="section-label">{{ t('home.about_eyebrow') }} · EST. MMXII</div>
            <h2 class="home-section__title">
              {{ t('home.about_brand_title') }}<span class="red-dot">.</span>
            </h2>
          </div>
          <RouterLink to="/about" class="arrow-link">{{ t('home.about_more') }}</RouterLink>
        </div>
        <div class="about-inner">
          <div class="about-text">
            <p class="about-p">{{ t('home.about_p1') }}</p>
            <p class="about-p">{{ t('home.about_p2') }}</p>
            <div class="about-values">
              <div class="about-value">
                <span class="about-value__icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="var(--accent)">
                    <path d="M12 2L2 9l10 13L22 9 12 2z"/>
                  </svg>
                </span>
                <div>
                  <p class="about-value__title">{{ t('home.value_quality_title') }}</p>
                  <p class="about-value__desc">{{ t('home.value_quality_desc') }}</p>
                </div>
              </div>
              <div class="about-value">
                <span class="about-value__icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="var(--accent)">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                  </svg>
                </span>
                <div>
                  <p class="about-value__title">{{ t('home.value_trust_title') }}</p>
                  <p class="about-value__desc">{{ t('home.value_trust_desc') }}</p>
                </div>
              </div>
              <div class="about-value">
                <span class="about-value__icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="var(--accent)">
                    <path d="M7 2v11h3v9l7-12h-4l4-8z"/>
                  </svg>
                </span>
                <div>
                  <p class="about-value__title">{{ t('home.value_innov_title') }}</p>
                  <p class="about-value__desc">{{ t('home.value_innov_desc') }}</p>
                </div>
              </div>
            </div>
          </div>
          <div class="about-imgs">
            <div class="about-img about-img--main">
              <img src="/images/общий.jpg" alt="Салон" />
            </div>
            <div class="about-img about-img--sm">
              <img src="/images/стулья.jpg" alt="Салон" />
            </div>
            <div class="about-img about-img--sm">
              <img src="/images/мэйн.jpg" alt="Салон" />
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ═══ CTA ═══════════════════════════════════════ -->
    <section class="section-cta">
      <div class="page-container">
        <div class="cta-inner">
          <div class="section-label">{{ t('home.cta_eyebrow') }}</div>
          <h2 class="t-h1 cta-title">
            {{ t('home.cta_title') }}<span class="red-dot">.</span>
          </h2>
          <button class="btn btn-primary btn-lg" @click="router.push('/book')">
            {{ t('home.cta_btn') }} →
          </button>
        </div>
      </div>
    </section>

  </div>
</template>

<style scoped>
/* ── HERO FULLSCREEN ─────────────────────────────── */
.hero-full {
  position: relative;
  min-height: 560px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.hero-full__bg {
  position: absolute; inset: 0;
}
.hero-full__img {
  width: 100%; height: 100%; object-fit: cover;
  object-position: center 30%;
}
.hero-full__overlay {
  position: absolute; inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(26,24,20,0.2) 0%,
    rgba(26,24,20,0.65) 50%,
    rgba(26,24,20,0.88) 100%
  );
}
.hero-full__inner {
  position: relative; z-index: 1;
  flex: 1;
  padding: 60px 0 56px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  border-top: 1px solid rgba(240,236,228,0.1);
  border-bottom: 1px solid rgba(240,236,228,0.1);
}
.hero-full__eyebrow {
  font-size: 11px; font-weight: 500; letter-spacing: 0.12em;
  text-transform: uppercase; color: rgba(240,236,228,0.5);
  display: flex; align-items: center; gap: 8px;
  margin-bottom: 20px; justify-content: center;
}
.hero-full__dash { color: var(--accent); }
.hero-full__title {
  font-family: var(--font-serif);
  font-size: clamp(48px, 8vw, 96px);
  font-weight: 900; line-height: 0.92;
  letter-spacing: -0.025em; color: #fff;
  margin-bottom: 20px;
}
.hero-full__title em { font-style: italic; color: rgba(240,236,228,0.6); }
.hero-full__sub {
  font-size: 16px; color: rgba(240,236,228,0.6);
  line-height: 1.7; max-width: 520px; margin-bottom: 36px;
}
.hero-full__btns {
  display: flex; gap: 12px; flex-wrap: wrap;
  margin-bottom: 28px; justify-content: center;
}
.btn-outline-light {
  border: 1px solid rgba(240,236,228,0.4);
  color: rgba(240,236,228,0.9);
  background: transparent;
  padding: 0 24px; height: 48px;
  display: inline-flex; align-items: center;
  font-size: 14px; font-weight: 500;
  border-radius: var(--radius-sm);
  text-decoration: none;
  transition: all 0.2s; cursor: pointer;
}
.btn-outline-light:hover {
  border-color: rgba(240,236,228,0.8);
  background: rgba(240,236,228,0.1);
}
.hero-full__perks {
  display: flex; align-items: center; gap: 10px; flex-wrap: wrap;
  font-size: 11px; font-weight: 500; letter-spacing: 0.06em;
  color: rgba(240,236,228,0.4); justify-content: center;
}
.hero-full__perks-dot { opacity: 0.3; }

/* ── STATS BAR — отдельный блок ─────────────────── */
.stats-bar {
  background: var(--bg-dark, #1a1814);
  border-bottom: 1px solid var(--border);
}
.stats-bar__inner {
  max-width: 1200px; margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr auto 1fr auto 1fr auto 1fr;
  align-items: stretch;
}
.stats-bar__item {
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  gap: 6px; padding: 24px 16px;
}
.stats-bar__num {
  font-family: var(--font-serif); font-size: 32px; font-weight: 700;
  color: rgba(240,236,228,0.92); line-height: 1;
}
.stats-bar__label {
  font-size: 9px; font-weight: 500; letter-spacing: 0.14em;
  text-transform: uppercase; color: rgba(240,236,228,0.35);
  white-space: nowrap;
}
.stats-bar__divider {
  width: 1px;
  background: rgba(240,236,228,0.08);
  align-self: stretch;
  margin: 12px 0;
}

/* ── УНИФИЦИРОВАННЫЕ СЕКЦИИ ──────────────────────── */
.home-section { padding: 80px 0; border-bottom: 1px solid var(--border); }
.home-section__head {
  display: flex; justify-content: space-between;
  align-items: flex-end; margin-bottom: 48px;
}
.home-section__title {
  font-family: var(--font-serif);
  font-size: clamp(32px, 4vw, 52px);
  font-weight: 900; line-height: 1.0;
  letter-spacing: -0.02em; color: var(--text);
  margin-top: 8px;
}

/* ── УСЛУГИ — КАРТОЧКИ ───────────────────────────── */
.services-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}
.svc-card {
  position: relative;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer; overflow: hidden;
  transition: all 0.2s ease;
  display: flex; align-items: stretch;
}
.svc-card:hover {
  border-color: var(--cat-color, var(--accent));
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.06);
}
.svc-card:hover .svc-card__arrow { opacity: 1; color: var(--cat-color, var(--accent)); }
.svc-card:hover .svc-card__name  { color: var(--cat-color, var(--accent)); }

/* цветная полоска слева */
.svc-card__accent {
  width: 3px; flex-shrink: 0;
  background: var(--cat-color, var(--accent));
  opacity: 0.7;
  transition: opacity 0.2s;
}
.svc-card:hover .svc-card__accent { opacity: 1; }

.svc-card__body {
  flex: 1; padding: 18px 16px;
  display: flex; flex-direction: column; gap: 10px;
}
.svc-card__top {
  display: flex; justify-content: space-between;
  align-items: center;
}
.svc-card__num {
  font-family: var(--font-serif); font-size: 13px; font-weight: 700;
  color: var(--text); opacity: 0.15; line-height: 1;
}
.svc-card__cat {
  font-size: 9px; font-weight: 600; letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--cat-color, var(--accent));
  opacity: 0.8;
}
.svc-card__name {
  font-family: var(--font-serif); font-size: 17px; font-weight: 700;
  color: var(--text); line-height: 1.2;
  transition: color 0.2s;
  flex: 1;
}
.svc-card__footer {
  display: flex; justify-content: space-between; align-items: baseline;
  margin-top: auto;
}
.svc-card__dur   { font-size: 11px; color: var(--text-3); }
.svc-card__price {
  font-family: var(--font-serif); font-size: 16px; font-weight: 700; color: var(--text);
}
.svc-card__arrow {
  position: absolute; bottom: 14px; right: 14px;
  font-size: 14px; color: var(--text-3); opacity: 0.3;
  transition: all 0.2s;
}

/* ── СЛАЙДЕР МАСТЕРОВ ────────────────────────────── */
.section-masters { background: var(--bg-muted); }
.masters-slider {
  display: grid; grid-template-columns: 420px 1fr;
  gap: 0; border: 1px solid var(--border); border-radius: var(--radius-sm);
  overflow: hidden; min-height: 480px;
}

/* фото */
.masters-slider__photo-wrap {
  position: relative; background: var(--bg-card); overflow: hidden;
}
.masters-slider__photo {
  width: 100%; height: 100%; min-height: 480px;
}
.masters-slider__photo img {
  width: 100%; height: 100%; object-fit: cover; object-position: center top;
  display: block;
}
.masters-slider__initials {
  width: 100%; height: 100%; min-height: 480px;
  display: flex; align-items: center; justify-content: center;
  font-family: var(--font-serif); font-size: 96px; font-weight: 700;
  color: var(--accent); background: var(--bg-muted);
}
.masters-slider__fig {
  position: absolute; bottom: 16px; left: 16px;
  font-size: 9px; font-weight: 500; letter-spacing: 0.14em;
  text-transform: uppercase; color: rgba(255,255,255,0.6);
  background: rgba(26,24,20,0.45); padding: 4px 10px; border-radius: 1px;
}

/* инфо */
.masters-slider__body {
  padding: 48px 40px;
  display: flex; flex-direction: column; justify-content: space-between;
  border-left: 1px solid var(--border);
}
.masters-slider__info { flex: 1; }
.masters-slider__exp {
  font-size: 10px; font-weight: 500; letter-spacing: 0.1em;
  text-transform: uppercase; color: var(--accent); margin-bottom: 12px;
}
.masters-slider__name {
  font-family: var(--font-serif); font-size: clamp(28px, 3vw, 40px);
  font-weight: 900; color: var(--text); line-height: 1.0;
  letter-spacing: -0.02em; margin-bottom: 20px;
}
.masters-slider__bio {
  font-size: 15px; color: var(--text-2); line-height: 1.7;
  margin-bottom: 24px;
  display: -webkit-box; -webkit-line-clamp: 4; line-clamp: 4;
  -webkit-box-orient: vertical; overflow: hidden;
}
.masters-slider__tags {
  display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 32px;
}
.masters-slider__tag {
  font-size: 10px; font-weight: 500; letter-spacing: 0.06em;
  padding: 4px 12px; border: 1px solid var(--border-strong);
  border-radius: var(--radius-sm); color: var(--text-3);
}

/* навигация слайдера */
.masters-slider__nav {
  display: flex; align-items: center; gap: 16px;
  padding-top: 24px; border-top: 1px solid var(--border);
  margin-top: 24px;
}
.masters-slider__arrow {
  width: 36px; height: 36px;
  display: flex; align-items: center; justify-content: center;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm); color: var(--text-2);
  cursor: pointer; transition: all 0.15s; flex-shrink: 0;
}
.masters-slider__arrow:hover { border-color: var(--accent); color: var(--accent); }
.masters-slider__dots {
  display: flex; gap: 8px; flex: 1;
}
.masters-slider__dot {
  width: 6px; height: 6px; border-radius: 50%;
  background: var(--border-strong);
  border: none; cursor: pointer; padding: 0;
  transition: all 0.2s;
}
.masters-slider__dot--active {
  background: var(--accent); transform: scale(1.4);
}
.masters-slider__counter {
  font-family: var(--font-serif); font-size: 13px;
  font-weight: 700; color: var(--text-3);
  letter-spacing: 0.04em;
}

/* анимация смены мастера */
.master-fade-enter-active { transition: opacity 0.4s ease, transform 0.4s ease; }
.master-fade-leave-active { transition: opacity 0.25s ease, transform 0.25s ease; }
.master-fade-enter-from  { opacity: 0; transform: translateX(20px); }
.master-fade-leave-to    { opacity: 0; transform: translateX(-20px); }

/* ── КАК РАБОТАЕТ ────────────────────────────────── */
.how-steps {
  display: grid; grid-template-columns: repeat(4, 1fr); gap: 0; position: relative;
}
.how-step { padding: 0 32px 0 0; position: relative; }
.how-step:not(:last-child)::after {
  content: ''; position: absolute;
  top: 28px; right: 16px; width: calc(100% - 60px);
  height: 1px; background: var(--border);
}
.how-step:hover .how-step__num   { color: var(--accent); opacity: 1; }
.how-step:hover .how-step__title { transform: translateY(-2px); }
.how-step__num {
  font-family: var(--font-serif); font-size: 48px; font-weight: 900;
  color: var(--text); opacity: 0.12; line-height: 1; margin-bottom: 16px;
  transition: color 0.2s, opacity 0.2s;
}
.how-step__title {
  font-family: var(--font-serif); font-size: 18px; font-weight: 700;
  color: var(--text); margin-bottom: 8px; transition: transform 0.2s ease;
}
.how-step__desc { font-size: 13px; color: var(--text-3); line-height: 1.6; }

/* ── ГАЛЕРЕЯ ─────────────────────────────────────── */
.section-gallery { background: var(--bg-muted); }
.gallery-blog { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.gallery-blog__featured { text-decoration: none; }
.gallery-blog__grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.gallery-blog__img-wrap {
  position: relative; overflow: hidden;
  border-radius: var(--radius-sm); background: var(--bg-card);
}
.gallery-blog__featured .gallery-blog__img-wrap { aspect-ratio: 4/3; }
.gallery-blog__item     .gallery-blog__img-wrap { aspect-ratio: 1/1; }
.gallery-blog__img-wrap img {
  width: 100%; height: 100%; object-fit: cover; transition: transform 0.3s ease;
}
.gallery-blog__featured:hover img,
.gallery-blog__item:hover img { transform: scale(1.04); }
.gallery-blog__overlay {
  position: absolute; inset: 0; background: rgba(26,24,20,0.3);
  display: flex; align-items: center; justify-content: center;
  opacity: 0; transition: opacity 0.2s;
}
.gallery-blog__featured:hover .gallery-blog__overlay,
.gallery-blog__item:hover     .gallery-blog__overlay { opacity: 1; }
.gallery-blog__arrow { font-size: 24px; color: #fff; font-weight: 300; }
.gallery-blog__no {
  position: absolute; top: 10px; left: 10px;
  font-size: 9px; font-weight: 500; letter-spacing: 0.12em;
  text-transform: uppercase; color: rgba(255,255,255,0.7);
  background: rgba(26,24,20,0.4); padding: 3px 7px; border-radius: 1px;
}
.gallery-blog__info    { padding: 10px 0 0; }
.gallery-blog__master  { font-size: 12px; font-weight: 600; color: var(--text); }
.gallery-blog__service { font-size: 11px; color: var(--text-3); margin-top: 2px; }
.gallery-blog__item    { text-decoration: none; }

/* ── ОТЗЫВЫ — КАРУСЕЛЬ ───────────────────────────── */
.reviews-nav { display: flex; gap: 8px; }
.reviews-nav__btn {
  width: 36px; height: 36px;
  display: flex; align-items: center; justify-content: center;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm); color: var(--text-2);
  cursor: pointer; transition: all 0.15s;
}
.reviews-nav__btn:hover:not(:disabled) { border-color: var(--accent); color: var(--accent); }
.reviews-nav__btn:disabled { opacity: 0.3; cursor: not-allowed; }

.reviews-track {
  display: flex; gap: 24px;
  overflow-x: auto; scroll-snap-type: x mandatory;
  scrollbar-width: none; -ms-overflow-style: none;
  padding-bottom: 4px;
}
.reviews-track::-webkit-scrollbar { display: none; }
.review-card {
  min-width: 320px; max-width: 360px;
  scroll-snap-align: start; flex-shrink: 0;
  padding: 28px 24px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--bg-card);
  transition: border-color 0.2s;
}
.review-card:hover { border-color: var(--border-strong); }
.review-card__quote {
  font-family: var(--font-serif); font-size: 56px; font-weight: 700;
  color: var(--accent); line-height: 0.7; margin-bottom: 14px; opacity: 0.4;
}
.review-card__text {
  font-size: 14px; color: var(--text-2); line-height: 1.7;
  font-style: italic; margin-bottom: 20px;
  display: -webkit-box; -webkit-line-clamp: 4; line-clamp: 4;
  -webkit-box-orient: vertical; overflow: hidden;
}
.review-card__footer {
  display: flex; flex-direction: column; gap: 4px;
  padding-top: 16px; border-top: 1px solid var(--border);
}
.review-card__stars  { font-size: 13px; color: var(--accent); letter-spacing: 2px; }
.review-card__name   { font-size: 13px; font-weight: 500; color: var(--text); }
.review-card__service { font-size: 11px; color: var(--text-3); }

/* ── О НАС (краткая) ─────────────────────────────── */
.about-inner {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 60px; align-items: start;
}
.about-p { font-size: 15px; color: var(--text-2); line-height: 1.8; margin-bottom: 16px; }
.about-values { margin-top: 28px; display: flex; flex-direction: column; gap: 18px; }
.about-value  { display: flex; align-items: flex-start; gap: 14px; }
.about-value__icon {
  width: 40px; height: 40px; flex-shrink: 0;
  background: var(--accent-bg); border-radius: var(--radius-sm);
  display: flex; align-items: center; justify-content: center;
}
.about-value__title { font-size: 14px; font-weight: 600; color: var(--text); margin-bottom: 3px; }
.about-value__desc  { font-size: 13px; color: var(--text-3); }
.about-imgs {
  display: grid; grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto; gap: 10px;
}
.about-img { border-radius: var(--radius-sm); overflow: hidden; }
.about-img img { width: 100%; height: 100%; object-fit: cover; display: block; }
.about-img--main { grid-row: span 2; }
.about-img--main img { aspect-ratio: 2/3; }
.about-img--sm img   { aspect-ratio: 1/1; }

/* ── CTA ─────────────────────────────────────────── */
.section-cta { padding: 100px 0; }
.cta-inner {
  text-align: center; display: flex;
  flex-direction: column; align-items: center; gap: 24px;
}
.cta-title { max-width: 600px; }

/* ── АДАПТИВ 900px ───────────────────────────────── */
@media (max-width: 900px) {
  .hero-full       { min-height: 480px; }
  .hero-full__inner { padding: 40px 0 36px; }
  .stats-bar__num  { font-size: 24px; }
  .stats-bar__item { padding: 18px 10px; }

  .home-section    { padding: 56px 0; }
  .home-section__head { margin-bottom: 32px; flex-wrap: wrap; gap: 12px; }

  /* услуги */
  .services-grid   { grid-template-columns: repeat(2, 1fr); }

  /* слайдер мастеров */
  .masters-slider  { grid-template-columns: 1fr; }
  .masters-slider__photo-wrap { height: 320px; }
  .masters-slider__photo      { min-height: 320px; }
  .masters-slider__initials   { min-height: 320px; font-size: 64px; }
  .masters-slider__body       { padding: 28px 24px; }

  /* как работает */
  .how-steps       { grid-template-columns: 1fr 1fr; gap: 32px; }
  .how-step::after { display: none; }
  .how-step        { padding: 0; }

  /* галерея */
  .gallery-blog     { grid-template-columns: 1fr; }
  .gallery-blog__grid { grid-template-columns: 1fr 1fr; }

  /* о нас */
  .about-inner     { grid-template-columns: 1fr; gap: 32px; }
  .about-imgs      { grid-template-columns: 1fr 1fr; }
  .about-img--main { grid-row: auto; }
  .about-img--main img { aspect-ratio: 16/9; }
}

/* ── АДАПТИВ 600px ───────────────────────────────── */
@media (max-width: 600px) {
  .hero-full        { min-height: 400px; }
  .hero-full__inner { padding: 32px 0 28px; }
  .hero-full__title { font-size: clamp(38px, 11vw, 56px); }
  .hero-full__btns  { flex-direction: column; }
  .hero-full__btns .btn,
  .hero-full__btns .btn-outline-light { width: 100%; justify-content: center; }
  .hero-full__sub   { font-size: 14px; margin-bottom: 24px; }
  .hero-full__perks { display: none; }

  /* stats — 2x2 на мобиле */
  .stats-bar__inner {
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
  }
  .stats-bar__divider { display: none; }
  .stats-bar__item {
    padding: 16px 8px;
    border-right: 1px solid rgba(240,236,228,0.08);
    border-bottom: 1px solid rgba(240,236,228,0.08);
  }
  .stats-bar__item:nth-child(2n) { border-right: none; }
  .stats-bar__item:nth-last-child(-n+2) { border-bottom: none; }
  .stats-bar__num   { font-size: 22px; }
  .stats-bar__label { font-size: 8px; }

  .home-section    { padding: 40px 0; }
  .home-section__head { flex-direction: column; align-items: flex-start; gap: 10px; }
  .home-section__title { font-size: clamp(26px, 8vw, 36px); }

  /* услуги */
  .services-grid   { grid-template-columns: repeat(2, 1fr); gap: 8px; }
  .svc-card__name  { font-size: 15px; }
  .svc-card__price { font-size: 14px; }

  /* слайдер */
  .masters-slider__photo-wrap { height: 260px; }
  .masters-slider__photo      { min-height: 260px; }
  .masters-slider__initials   { min-height: 260px; }
  .masters-slider__body       { padding: 20px 16px; }
  .masters-slider__name       { font-size: 24px; }
  .masters-slider__bio        { -webkit-line-clamp: 3; line-clamp: 3; }
  .masters-slider__tags       { margin-bottom: 20px; }

  /* как работает */
  .how-steps       { grid-template-columns: 1fr; gap: 24px; }
  .how-step__num   { font-size: 36px; }
  .how-step__title { font-size: 16px; }

  /* галерея */
  .gallery-blog__grid { grid-template-columns: 1fr; }

  /* отзывы */
  .review-card     { min-width: 280px; }

  /* о нас */
  .about-imgs      { display: none; }

  /* cta */
  .section-cta     { padding: 56px 0; }
  .cta-title       { font-size: clamp(26px, 8vw, 36px); }
}
</style>