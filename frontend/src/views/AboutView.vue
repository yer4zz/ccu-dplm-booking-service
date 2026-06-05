<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router'
const router = useRouter()
const { t } = useI18n()
</script>

<template>
  <div class="about-page">

    <!-- ══ HERO ═══════════════════════════════════ -->
    <section class="about-hero">
      <div class="about-hero__bg">
        <img src="/images/мэйн.jpg" alt="Beauty Dana" class="about-hero__img" />
        <div class="about-hero__overlay" />
      </div>
      <div class="page-container about-hero__inner">
        <div class="about-hero__eyebrow">
          <span class="about-hero__dash">—</span>
          {{ t('about.eyebrow') }}
        </div>
        <h1 class="about-hero__title">
          Beauty<br/>
          <em>Dana</em><span class="red-dot">.</span>
        </h1>
        <p class="about-hero__sub">
          Студия красоты в городе Есик. С 2012 года мы создаём образы, которые вдохновляют.
        </p>
      </div>
    </section>

    <!-- ══ ИСТОРИЯ ════════════════════════════════ -->
    <section class="about-story">
      <div class="page-container about-story__inner">
        <div class="about-story__left">
          <div class="section-label" style="margin-bottom:20px">{{ t('about.history_label') }}</div>
          <h2 class="t-h2" style="margin-bottom:24px">
            {{ t('about.history_title') }}<span class="red-dot">.</span>
          </h2>
          <p class="about-p">
            Beauty Dana открылась в 2012 году в Есике. С первых дней мы поставили перед собой
            одну цель — сделать так, чтобы каждый клиент выходил от нас с улыбкой и уверенностью в себе.
          </p>
          <p class="about-p">
            За эти годы мы выросли из небольшого салона в полноценную студию красоты,
            где работают профессионалы с многолетним опытом. Мы постоянно учимся,
            следим за тенденциями и внедряем новые техники.
          </p>
          <p class="about-p">
            В 2026 году мы запустили онлайн-систему бронирования — теперь записаться
            к любому мастеру можно в любое время суток, без звонков и ожидания ответа.
          </p>
        </div>
        <div class="about-story__right">
          <div class="about-photos">
            <div class="about-photo about-photo--main">
              <img src="/images/стулья.jpg" alt="Интерьер салона" />
              <span class="about-photo__label">FIG. 01 · {{ t('about.fig_interior') }}</span>
            </div>
            <div class="about-photo about-photo--sm">
              <img src="/images/общий.jpg" alt="Салон" />
              <span class="about-photo__label">FIG. 02</span>
            </div>
            <div class="about-photo about-photo--sm">
              <img src="/images/мэйн.jpg" alt="Мастер" />
              <span class="about-photo__label">FIG. 03</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ══ ЦИФРЫ ══════════════════════════════════ -->
    <section class="about-numbers">
      <div class="page-container">
        <div class="section-label" style="margin-bottom:40px;text-align:center">{{ t('about.numbers_label') }}</div>
        <div class="numbers-grid">
          <div class="number-item">
            <span class="number-item__val">2012</span>
            <span class="number-item__label">{{ t('about.num1_label') }}</span>
            <p class="number-item__desc">{{ t('about.num1_desc') }}</p>
          </div>
          <div class="number-item">
            <span class="number-item__val">500+</span>
            <span class="number-item__label">{{ t('about.num2_label') }}</span>
            <p class="number-item__desc">{{ t('about.num2_desc') }}</p>
          </div>
          <div class="number-item">
            <span class="number-item__val">10+</span>
            <span class="number-item__label">{{ t('about.num3_label') }}</span>
            <p class="number-item__desc">{{ t('about.num3_desc') }}</p>
          </div>
          <div class="number-item">
            <span class="number-item__val">24/7</span>
            <span class="number-item__label">{{ t('about.num4_label') }}</span>
            <p class="number-item__desc">{{ t('about.num4_desc') }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ══ ЦЕННОСТИ ═══════════════════════════════ -->
    <section class="about-values-section">
      <div class="page-container">
        <div class="about-values-inner">
          <div class="about-values-left">
            <div class="section-label" style="margin-bottom:20px">{{ t('about.values_label') }}</div>
            <h2 class="t-h2" style="margin-bottom:16px">
              {{ t('about.values_title') }}<span class="red-dot">.</span>
            </h2>
            <p style="font-size:15px;color:var(--text-2);line-height:1.7">
              Мы верим что красота — это не просто внешность. Это уверенность, комфорт и хорошее настроение.
            </p>
          </div>
          <div class="values-list">
            <div class="value-item">
              <div class="value-item__icon">
                <!-- star / quality -->
                <svg width="22" height="22" viewBox="0 0 24 24" fill="var(--accent)">
                  <path d="M12 17.27L18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"/>
                </svg>
              </div>
              <div class="value-item__body">
                <h3 class="value-item__title">{{ t('about.val1_title') }}</h3>
                <p class="value-item__desc">{{ t('about.val1_desc') }}</p>
              </div>
            </div>
            <div class="value-item">
              <div class="value-item__icon">
                <!-- heart / trust -->
                <svg width="22" height="22" viewBox="0 0 24 24" fill="var(--accent)">
                  <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                </svg>
              </div>
              <div class="value-item__body">
                <h3 class="value-item__title">{{ t('about.val2_title') }}</h3>
                <p class="value-item__desc">{{ t('about.val2_desc') }}</p>
              </div>
            </div>
            <div class="value-item">
              <div class="value-item__icon">
                <!-- clock / time -->
                <svg width="22" height="22" viewBox="0 0 24 24" fill="var(--accent)">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm.5 5v5.25l4.5 2.67-.75 1.23L11 13V7h1.5z"/>
                </svg>
              </div>
              <div class="value-item__body">
                <h3 class="value-item__title">{{ t('about.val3_title') }}</h3>
                <p class="value-item__desc">{{ t('about.val3_desc') }}</p>
              </div>
            </div>
            <div class="value-item">
              <div class="value-item__icon">
                <!-- graduation / education -->
                <svg width="22" height="22" viewBox="0 0 24 24" fill="var(--accent)">
                  <path d="M5 13.18v4L12 21l7-3.82v-4L12 17l-7-3.82zM12 3L1 9l11 6 9-4.91V17h2V9L12 3z"/>
                </svg>
              </div>
              <div class="value-item__body">
                <h3 class="value-item__title">{{ t('about.val4_title') }}</h3>
                <p class="value-item__desc">{{ t('about.val4_desc') }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ══ ФОТО ГАЛЕРЕЯ ИНТЕРЬЕРА ══════════════════ -->
    <section class="about-gallery">
      <div class="page-container">
        <div class="section-label" style="margin-bottom:32px">{{ t('about.gallery_label') }}</div>
        <div class="about-gallery-grid">
          <div class="about-gallery-item about-gallery-item--wide">
            <img src="/images/общий.jpg" alt="Салон Beauty Dana" loading="lazy" />
            <div class="about-gallery-item__overlay">
              <span class="about-gallery-item__no">Nº 01 · {{ t('about.fig_general') }}</span>
            </div>
          </div>
          <div class="about-gallery-item">
            <img src="/images/стулья.jpg" alt="Рабочие места" loading="lazy" />
            <div class="about-gallery-item__overlay">
              <span class="about-gallery-item__no">Nº 02 · {{ t('about.fig_workplaces') }}</span>
            </div>
          </div>
          <div class="about-gallery-item">
            <img src="/images/мэйн.jpg" alt="Студия" loading="lazy" />
            <div class="about-gallery-item__overlay">
              <span class="about-gallery-item__no">Nº 03 · {{ t('about.fig_studio') }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ══ КОНТАКТЫ ════════════════════════════════ -->
    <section class="about-contacts">
      <div class="page-container about-contacts__inner">
        <div class="about-contacts__info">
          <div class="section-label" style="margin-bottom:20px">{{ t('about.contacts_label') }}</div>
          <h2 class="t-h2" style="margin-bottom:32px">
            Найдите нас<span class="red-dot">.</span>
          </h2>
          <div class="contact-list">
            <div class="contact-item">
              <span class="contact-item__icon">📍</span>
              <div>
                <p class="contact-item__label">{{ t('about.addr_label') }}</p>
                <p class="contact-item__val">{{ t('about.addr_val') }}</p>
              </div>
            </div>
            <div class="contact-item">
              <span class="contact-item__icon">📞</span>
              <div>
                <p class="contact-item__label">{{ t('about.phone_label') }}</p>
                <p class="contact-item__val">+7 777 000 00 00</p>
              </div>
            </div>
            <div class="contact-item">
              <span class="contact-item__icon">🕐</span>
              <div>
                <p class="contact-item__label">{{ t('about.hours_label') }}</p>
                <p class="contact-item__val">{{ t('about.hours_val') }}</p>
                <p class="contact-item__val">{{ t('about.hours_sun') }}</p>
              </div>
            </div>
            <div class="contact-item">
              <span class="contact-item__icon">📱</span>
              <div>
                <p class="contact-item__label">{{ t('about.socials_label') }}</p>
                <div style="display:flex;gap:12px;margin-top:4px">
                  <a href="#" class="contact-social">Instagram</a>
                  <a href="#" class="contact-social">WhatsApp</a>
                  <a href="#" class="contact-social">Telegram</a>
                </div>
              </div>
            </div>
          </div>
          <button class="btn btn-primary btn-lg" style="margin-top:32px" @click="router.push('/book')">
            Записаться онлайн →
          </button>
        </div>

        <!-- декоративный блок -->
        <div class="about-contacts__map">
          <div class="about-map-box">
            <div class="about-map-box__label">FIG. 06 · РАСПОЛОЖЕНИЕ</div>
            <div class="about-map-box__frame">
              <iframe
                src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d2912.0!2d77.461282!3d43.355135!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x0%3A0x0!2zNDPCsDIxJzE4LjUiTiA3N8KwMjcnNDAuNiJF!5e0!3m2!1sru!2skz!4v1700000000000!5m2!1sru!2skz"
                width="100%"
                height="100%"
                style="border:0"
                allowfullscreen
                loading="lazy"
                referrerpolicy="no-referrer-when-downgrade"
                title="Beauty Dana на карте"
              />
            </div>
            <div class="about-map-box__footer">
              <span class="about-map-box__addr">{{ t('about.addr_val') }}</span>
              <a
                href="https://www.google.com/maps/dir/?api=1&destination=43.355135,77.461282"
                target="_blank"
                rel="noopener"
                class="about-map-box__route"
              >{{ t('about.route_btn') }}</a>
            </div>
          </div>
        </div>
      </div>
    </section>

  </div>
</template>

<style scoped>
.about-page { min-height: calc(100vh - 56px); }

/* ── HERO ─────────────────────────────────────── */
.about-hero {
  position: relative; height: 500px;
  display: flex; align-items: flex-end; overflow: hidden;
}
.about-hero__bg { position: absolute; inset: 0; }
.about-hero__img { width: 100%; height: 100%; object-fit: cover; object-position: center 40%; }
.about-hero__overlay {
  position: absolute; inset: 0;
  background: linear-gradient(to bottom, rgba(26,24,20,0.1) 0%, rgba(26,24,20,0.75) 60%, rgba(26,24,20,0.95) 100%);
}
.about-hero__inner { position: relative; z-index: 1; padding-bottom: 56px; }
.about-hero__eyebrow {
  font-size: 11px; font-weight: 500; letter-spacing: 0.12em;
  text-transform: uppercase; color: rgba(240,236,228,0.55);
  display: flex; align-items: center; gap: 8px; margin-bottom: 16px;
}
.about-hero__dash { color: var(--accent); }
.about-hero__title {
  font-family: var(--font-serif);
  font-size: clamp(52px, 9vw, 100px);
  font-weight: 900; line-height: 0.9;
  letter-spacing: -0.025em; color: #fff; margin-bottom: 20px;
}
.about-hero__title em { font-style: italic; color: rgba(240,236,228,0.65); }
.about-hero__sub { font-size: 16px; color: rgba(240,236,228,0.6); line-height: 1.7; max-width: 500px; }

/* ── ИСТОРИЯ ──────────────────────────────────── */
.about-story { padding: 100px 0; border-bottom: 1px solid var(--border); }
.about-story__inner {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 80px; align-items: start;
}
.about-p { font-size: 15px; color: var(--text-2); line-height: 1.8; margin-bottom: 16px; }

.about-photos {
  display: grid; grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto; gap: 8px;
}
.about-photo         { position: relative; overflow: hidden; border-radius: var(--radius-sm); }
.about-photo--main   { grid-row: span 2; }
.about-photo img     { width: 100%; height: 100%; object-fit: cover; display: block; }
.about-photo--main img { aspect-ratio: 2/3; }
.about-photo--sm img   { aspect-ratio: 1/1; }
.about-photo__label {
  position: absolute; bottom: 10px; left: 10px;
  font-size: 9px; font-weight: 500; letter-spacing: 0.12em;
  text-transform: uppercase; color: rgba(255,255,255,0.6);
  background: rgba(26,24,20,0.45); padding: 3px 8px; border-radius: 1px;
}

/* ── ЦИФРЫ ────────────────────────────────────── */
.about-numbers {
  padding: 80px 0; background: var(--bg-muted); border-bottom: 1px solid var(--border);
}
.numbers-grid {
  display: grid; grid-template-columns: repeat(4, 1fr);
  gap: 0;
}
.number-item {
  padding: 32px 24px; border-right: 1px solid var(--border);
  text-align: center;
}
.number-item:last-child { border-right: none; }
.number-item__val {
  font-family: var(--font-serif); font-size: 48px; font-weight: 900;
  color: var(--text); line-height: 1; display: block; margin-bottom: 8px;
}
.number-item__label {
  font-size: 10px; font-weight: 500; letter-spacing: 0.1em;
  text-transform: uppercase; color: var(--accent);
  display: block; margin-bottom: 12px;
}
.number-item__desc { font-size: 13px; color: var(--text-3); line-height: 1.6; }

/* ── ЦЕННОСТИ ─────────────────────────────────── */
.about-values-section { padding: 100px 0; border-bottom: 1px solid var(--border); }
.about-values-inner {
  display: grid; grid-template-columns: 1fr 1fr; gap: 80px; align-items: start;
}
.values-list { display: flex; flex-direction: column; gap: 0; }
.value-item {
  display: flex; gap: 20px; padding: 24px 0;
  border-bottom: 1px solid var(--border);
}
.value-item:first-child { border-top: 1px solid var(--border); }
.value-item__icon  {
  width: 44px; height: 44px; flex-shrink: 0; margin-top: 2px;
  background: var(--accent-bg); border-radius: var(--radius-sm);
  display: flex; align-items: center; justify-content: center;
}
.value-item__title { font-size: 16px; font-weight: 600; color: var(--text); margin-bottom: 6px; }
.value-item__desc  { font-size: 13px; color: var(--text-3); line-height: 1.6; }

/* ── ГАЛЕРЕЯ ИНТЕРЬЕРА ────────────────────────── */
.about-gallery { padding: 80px 0; background: var(--bg-muted); border-bottom: 1px solid var(--border); }
.about-gallery-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto auto;
  gap: 8px;
}
.about-gallery-item {
  position: relative; overflow: hidden;
  border-radius: var(--radius-sm); background: var(--bg-card);
  cursor: pointer;
}
.about-gallery-item--wide { grid-row: span 2; }
.about-gallery-item img {
  width: 100%; height: 100%; object-fit: cover; display: block;
  transition: transform 0.4s cubic-bezier(0.4,0,0.2,1);
}
.about-gallery-item--wide img { aspect-ratio: 2/3; }
.about-gallery-item:not(.about-gallery-item--wide) img { aspect-ratio: 16/9; }
.about-gallery-item:hover img { transform: scale(1.03); }
.about-gallery-item__overlay {
  position: absolute; inset: 0;
  background: linear-gradient(transparent 60%, rgba(26,24,20,0.6));
  display: flex; align-items: flex-end; padding: 16px;
}
.about-gallery-item__no {
  font-size: 9px; font-weight: 500; letter-spacing: 0.12em;
  text-transform: uppercase; color: rgba(255,255,255,0.65);
}

/* ── КОНТАКТЫ ─────────────────────────────────── */
.about-contacts { padding: 100px 0; }
.about-contacts__inner {
  display: grid; grid-template-columns: 1fr 1fr;
  gap: 80px; align-items: start;
}
.contact-list  { display: flex; flex-direction: column; gap: 28px; }
.contact-item  { display: flex; gap: 16px; align-items: flex-start; }
.contact-item__icon { font-size: 20px; flex-shrink: 0; margin-top: 2px; }
.contact-item__label { font-size: 10px; font-weight: 500; letter-spacing: 0.1em; text-transform: uppercase; color: var(--text-3); margin-bottom: 4px; }
.contact-item__val   { font-size: 15px; color: var(--text); }
.contact-social {
  font-size: 13px; font-weight: 500; color: var(--accent);
  text-decoration: none; letter-spacing: 0.04em;
  transition: opacity 0.15s;
}
.contact-social:hover { opacity: 0.7; }

/* декоративная карточка */
.about-contacts__map { display: flex; align-items: flex-start; }
.about-map-box {
  width: 100%; max-width: 580px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
}
.about-map-box__label {
  padding: 10px 16px;
  font-size: 9px; font-weight: 500; letter-spacing: 0.14em;
  text-transform: uppercase; color: var(--text-3);
  border-bottom: 1px solid var(--border);
  background: var(--bg-muted);
}
.about-map-box__frame {
  height: 380px;
  position: relative;
}
.about-map-box__frame iframe {
  display: block;
  width: 100%; height: 100%;
  filter: grayscale(20%);
  transition: filter 0.3s;
}
.about-map-box__frame:hover iframe { filter: grayscale(0%); }
.about-map-box__footer {
  display: flex; justify-content: space-between; align-items: center;
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  background: var(--bg-muted);
}
.about-map-box__addr {
  font-size: 13px; color: var(--text-2);
}
.about-map-box__route {
  font-size: 12px; font-weight: 500;
  color: var(--accent); text-decoration: none;
  letter-spacing: 0.04em; transition: opacity 0.15s;
}
.about-map-box__route:hover { opacity: 0.7; }

/* ── АДАПТИВ ──────────────────────────────────── */
@media (max-width: 900px) {
  .about-story__inner  { grid-template-columns: 1fr; gap: 40px; }
  .about-values-inner  { grid-template-columns: 1fr; gap: 40px; }
  .about-contacts__inner { grid-template-columns: 1fr; gap: 40px; }
  .about-map-box { max-width: 100%; }
  .about-map-box__frame { height: 280px; }
  .numbers-grid        { grid-template-columns: repeat(2, 1fr); }
  .number-item         { border-bottom: 1px solid var(--border); }
  .number-item:nth-child(2n) { border-right: none; }
}

@media (max-width: 600px) {
  .about-hero          { height: 360px; }
  .about-hero__title   { font-size: clamp(44px, 12vw, 64px); }
  .about-story         { padding: 60px 0; }
  .about-photos        { grid-template-columns: 1fr; }
  .about-photo--main   { grid-row: auto; }
  .about-photo--main img { aspect-ratio: 16/9; }
  .about-gallery-grid  { grid-template-columns: 1fr; }
  .about-gallery-item--wide { grid-row: auto; }
  .about-gallery-item--wide img { aspect-ratio: 16/9; }
  .numbers-grid        { grid-template-columns: 1fr 1fr; }
  .number-item__val    { font-size: 36px; }
  .about-numbers       { padding: 48px 0; }
  .about-values-section { padding: 60px 0; }
  .about-contacts      { padding: 60px 0; }
}
</style>