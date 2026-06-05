<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter }      from 'vue-router'
import { useI18n }        from 'vue-i18n'
import { mastersApi }     from '@/api/masters'
import { useBookingStore } from '@/stores/booking'
import type { Master }    from '@/types'

const { t } = useI18n()

const router = useRouter()
const store  = useBookingStore()

const mastera  = ref<Master[]>([])
const zagruzka = ref(true)
const vybran   = ref<Master | null>(null)

onMounted(async () => {
  mastera.value  = await mastersApi.list().catch(() => [])
  zagruzka.value = false
})

function bookMaster(m: Master) {
  store.reset()
  store.selectedMaster = m
  router.push('/book')
}

function openMaster(m: Master) {
  vybran.value = vybran.value?.id === m.id ? null : m
}
</script>

<template>
  <div class="masters-page">

    <section class="masters-hero">
      <div class="masters-hero__bg">
        <img src="/images/общий.jpg" alt="Салон" class="masters-hero__img" />
        <div class="masters-hero__overlay" />
      </div>
      <div class="page-container masters-hero__inner">
        <div class="masters-hero__eyebrow">
          <span class="masters-hero__dash">—</span>
          КОМАНДА · Nº 03
        </div>
        <h1 class="masters-hero__title">
          {{ t('masters.title_line1') }}<br/>
          <em>{{ t('masters.title_line2') }}</em><span class="red-dot">.</span>
        </h1>
        <p class="masters-hero__sub">
          Профессионалы своего дела с многолетним опытом. Каждый мастер — это история, стиль и внимание к деталям.
        </p>
      </div>
    </section>

    <section class="masters-stats">
      <div class="page-container masters-stats__inner">
        <div class="masters-stat">
          <span class="masters-stat__num">{{ mastera.length }}</span>
          <span class="masters-stat__label">{{ t('masters.stat1_label') }}</span>
        </div>
        <div class="masters-stat__sep">·</div>
        <div class="masters-stat">
          <span class="masters-stat__num">10+</span>
          <span class="masters-stat__label">{{ t('masters.stat2_label') }}</span>
        </div>
        <div class="masters-stat__sep">·</div>
        <div class="masters-stat">
          <span class="masters-stat__num">500+</span>
          <span class="masters-stat__label">{{ t('masters.stat3_label') }}</span>
        </div>
        <div class="masters-stat__sep">·</div>
        <div class="masters-stat">
          <span class="masters-stat__num">24/7</span>
          <span class="masters-stat__label">{{ t('masters.stat4_label') }}</span>
        </div>
      </div>
    </section>

    <section class="masters-list-section">
      <div class="page-container">

        <div v-if="zagruzka" class="masters-loading">
          <div class="masters-spinner" />
        </div>

        <div v-else class="masters-grid">
          <div
            v-for="(m, i) in mastera"
            :key="m.id"
            class="master-card"
            :class="{ 'master-card--open': vybran?.id === m.id }"
          >
            <div class="master-card__photo-wrap" @click="openMaster(m)">
              <img
                v-if="m.avatar_url"
                :src="m.avatar_url"
                :alt="m.full_name"
                class="master-card__photo"
              />
              <div v-else class="master-card__initials">
                {{ m.full_name[0] }}
              </div>
              <div class="master-card__photo-overlay">
                <span class="master-card__photo-icon">
                  {{ vybran?.id === m.id ? '✕' : '→' }}
                </span>
              </div>
              <span class="master-card__num">{{ String(i+1).padStart(2,'0') }}</span>
            </div>

            <div class="master-card__body">
              <div class="master-card__head">
                <div>
                  <h2 class="master-card__name">{{ m.full_name }}</h2>
                  <p class="master-card__exp">
                    {{ m.experience_years }} лет опыта
                  </p>
                </div>
                <button class="btn btn-primary btn-sm master-card__book" @click="bookMaster(m)">
                  Записаться →
                </button>
              </div>

              <p class="master-card__bio">{{ m.bio }}</p>

              <div class="master-card__tags">
                <span
                  v-for="s in m.services.slice(0, 4)"
                  :key="s.id"
                  class="master-card__tag"
                >{{ s.name }}</span>
                <span
                  v-if="m.services.length > 4"
                  class="master-card__tag master-card__tag--more"
                >+{{ m.services.length - 4 }}</span>
              </div>

              <Transition name="expand">
                <div v-if="vybran?.id === m.id" class="master-card__detail">
                  <div class="master-card__detail-divider" />
                  <p class="master-card__detail-label">{{ t('masters.all_services') }}</p>
                  <div class="master-card__services">
                    <div
                      v-for="s in m.services"
                      :key="s.id"
                      class="master-svc"
                    >
                      <span class="master-svc__name">{{ s.name }}</span>
                      <span class="master-svc__price">{{ s.price?.toLocaleString() }} ₸</span>
                    </div>
                  </div>
                  <button class="btn btn-primary" style="margin-top:16px;width:100%" @click="bookMaster(m)">
                    Записаться к {{ m.full_name.split(' ')[0] }} →
                  </button>
                </div>
              </Transition>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="masters-cta">
      <div class="page-container masters-cta__inner">
        <div class="section-label" style="margin-bottom:12px">{{ t('masters.cta_label') }}</div>
        <h2 class="t-h2 masters-cta__title">
          {{ t('masters.cta_title') }}<span class="red-dot">.</span>
        </h2>
        <button class="btn btn-primary btn-lg" @click="router.push('/book')">
          Записаться сейчас →
        </button>
      </div>
    </section>

  </div>
</template>

<style scoped>
.masters-page { min-height: calc(100vh - 56px); }

.masters-hero {
  position: relative;
  height: 480px;
  display: flex;
  align-items: flex-end;
  overflow: hidden;
}
.masters-hero__bg {
  position: absolute; inset: 0;
}
.masters-hero__img {
  width: 100%; height: 100%; object-fit: cover;
  object-position: center 30%;
}
.masters-hero__overlay {
  position: absolute; inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(26,24,20,0.2) 0%,
    rgba(26,24,20,0.7) 60%,
    rgba(26,24,20,0.92) 100%
  );
}
.masters-hero__inner {
  position: relative; z-index: 1;
  padding-bottom: 48px;
}
.masters-hero__eyebrow {
  font-family: var(--font); font-size: 11px; font-weight: 500;
  letter-spacing: 0.12em; text-transform: uppercase;
  color: rgba(240,236,228,0.6);
  display: flex; align-items: center; gap: 8px;
  margin-bottom: 16px;
}
.masters-hero__dash { color: var(--accent); }
.masters-hero__title {
  font-family: var(--font-serif);
  font-size: clamp(44px, 7vw, 80px);
  font-weight: 900; line-height: 0.95;
  letter-spacing: -0.025em;
  color: #fff; margin-bottom: 20px;
}
.masters-hero__title em { font-style: italic; color: rgba(240,236,228,0.7); }
.masters-hero__sub {
  font-size: 15px; color: rgba(240,236,228,0.65);
  line-height: 1.7; max-width: 480px;
}

.masters-stats { background: var(--bg-dark); border-bottom: 1px solid rgba(240,236,228,0.08); }
.masters-stats__inner {
  display: flex; align-items: center;
  gap: 16px; padding: 24px 40px;
  flex-wrap: wrap;
}
.masters-stat         { display: flex; align-items: baseline; gap: 10px; }
.masters-stat__num    {
  font-family: var(--font-serif); font-size: 28px; font-weight: 700;
  color: rgba(240,236,228,0.9);
}
.masters-stat__label  {
  font-size: 11px; font-weight: 500; letter-spacing: 0.1em;
  text-transform: uppercase; color: rgba(240,236,228,0.4);
}
.masters-stat__sep { color: rgba(240,236,228,0.2); font-size: 20px; }

.masters-list-section { padding: 80px 0; }
.masters-loading {
  display: flex; justify-content: center; padding: 80px 0;
}
.masters-spinner {
  width: 24px; height: 24px;
  border: 2px solid var(--border);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.masters-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 2px;
}

.master-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  transition: border-color 0.2s;
}
.master-card:hover        { border-color: var(--border-strong); }
.master-card--open        { border-color: var(--accent); }

.master-card__photo-wrap {
  position: relative; overflow: hidden;
  aspect-ratio: 4/3; cursor: pointer;
  background: var(--bg-muted); flex-shrink: 0;
}
.master-card__photo {
  width: 100%; height: 100%; object-fit: cover;
  transition: transform 0.4s cubic-bezier(0.4,0,0.2,1);
}
.master-card:hover .master-card__photo { transform: scale(1.04); }
.master-card__initials {
  width: 100%; height: 100%;
  display: flex; align-items: center; justify-content: center;
  font-family: var(--font-serif); font-size: 72px; font-weight: 700;
  color: var(--accent); background: var(--bg-muted);
}
.master-card__photo-overlay {
  position: absolute; inset: 0;
  background: rgba(26,24,20,0.4);
  display: flex; align-items: center; justify-content: center;
  opacity: 0; transition: opacity 0.2s;
}
.master-card:hover .master-card__photo-overlay { opacity: 1; }
.master-card__photo-icon {
  font-size: 28px; color: #fff; font-weight: 300;
}
.master-card__num {
  position: absolute; top: 16px; left: 16px;
  font-family: var(--font-serif); font-size: 11px; font-weight: 700;
  color: rgba(255,255,255,0.7);
  background: rgba(26,24,20,0.45); padding: 3px 8px; border-radius: 1px;
  letter-spacing: 0.1em;
}

.master-card__body { padding: 24px; flex: 1; }
.master-card__head {
  display: flex; justify-content: space-between;
  align-items: flex-start; gap: 12px; margin-bottom: 14px;
}
.master-card__name {
  font-family: var(--font-serif); font-size: 24px; font-weight: 700;
  color: var(--text); line-height: 1.1; margin-bottom: 4px;
}
.master-card__exp {
  font-size: 10px; font-weight: 500; letter-spacing: 0.1em;
  text-transform: uppercase; color: var(--accent);
}
.master-card__book { flex-shrink: 0; }
.master-card__bio {
  font-size: 14px; color: var(--text-2); line-height: 1.7; margin-bottom: 16px;
}
.master-card__tags { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 4px; }
.master-card__tag {
  font-size: 10px; font-weight: 500; letter-spacing: 0.06em;
  padding: 3px 10px; border: 1px solid var(--border);
  border-radius: var(--radius-sm); color: var(--text-3);
  white-space: nowrap;
}
.master-card__tag--more { color: var(--accent); border-color: var(--accent-bg); background: var(--accent-bg); }

.master-card__detail-divider {
  height: 1px; background: var(--border); margin: 20px 0;
}
.master-card__detail-label {
  font-size: 10px; font-weight: 500; letter-spacing: 0.1em;
  text-transform: uppercase; color: var(--text-3); margin-bottom: 12px;
}
.master-card__services { display: flex; flex-direction: column; gap: 0; }
.master-svc {
  display: flex; justify-content: space-between;
  padding: 10px 0; border-bottom: 1px solid var(--border);
  font-size: 14px;
}
.master-svc:last-child { border-bottom: none; }
.master-svc__name  { color: var(--text); }
.master-svc__price { font-weight: 600; color: var(--text); }

.expand-enter-active { transition: all 0.3s ease; overflow: hidden; }
.expand-leave-active { transition: all 0.2s ease; overflow: hidden; }
.expand-enter-from, .expand-leave-to { opacity: 0; transform: translateY(-8px); }

.masters-cta { padding: 100px 0; background: var(--bg-muted); border-top: 1px solid var(--border); }
.masters-cta__inner {
  display: flex; flex-direction: column;
  align-items: center; gap: 24px; text-align: center;
}
.masters-cta__title { max-width: 500px; }

@media (max-width: 768px) {
  .masters-hero  { height: 360px; }
  .masters-hero__title { font-size: clamp(36px, 10vw, 56px); }
  .masters-stats__inner { padding: 20px 16px; gap: 12px; }
  .masters-stat__num { font-size: 22px; }
  .masters-grid  { grid-template-columns: 1fr; }
  .masters-list-section { padding: 48px 0; }
  .master-card__head { flex-direction: column; gap: 12px; }
  .master-card__book { width: 100%; justify-content: center; }
  .masters-cta   { padding: 60px 0; }
}

@media (max-width: 480px) {
  .masters-hero  { height: 300px; }
  .master-card__body { padding: 16px; }
  .master-card__name { font-size: 20px; }
}
</style>