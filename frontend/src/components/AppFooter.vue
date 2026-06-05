<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const year = new Date().getFullYear()
</script>

<template>
  <footer class="footer">

    <!-- верхняя секция с большим лого -->
    <div class="footer__hero">
      <div class="footer__hero-inner page-container">
        <div class="footer__brand">
          <h2 class="footer__logo">Beauty Dana<span style="color: var(--accent)">.</span></h2>
          <p class="footer__sub">СТУДИЯ КРАСОТЫ · ЕСИК · EST. MMXII</p>
        </div>
        <RouterLink to="/book" class="footer__cta">
          {{ t('footer.cta') }}
        </RouterLink>
      </div>
    </div>

    <div class="footer__divider" />

    <!-- основная сетка -->
    <div class="footer__main page-container">

      <div class="footer__col">
        <p class="footer__col-title">{{ t('footer.nav_title') }}</p>
        <RouterLink to="/"        class="footer__link">{{ t('footer.home') }}</RouterLink>
        <RouterLink to="/book"    class="footer__link">{{ t('nav.book') }}</RouterLink>
        <RouterLink to="/gallery" class="footer__link">{{ t('nav.gallery') }}</RouterLink>
        <a href="/masters" class="footer__link">{{ t('nav.masters') }}</a>
        <a href="/book" class="footer__link">{{ t('footer.services') }}</a>
        <a href="/about" class="footer__link">{{ t('nav.about') }}</a>
      </div>

      <div class="footer__col">
        <p class="footer__col-title">{{ t('footer.services_title') }}</p>
        <a href="/book" class="footer__link">{{ t('footer.svc_hair') }}</a>
        <a href="/book" class="footer__link">{{ t('footer.svc_color') }}</a>
        <a href="/book" class="footer__link">{{ t('footer.svc_nails') }}</a>
        <a href="/book" class="footer__link">{{ t('footer.svc_brows') }}</a>
        <a href="/book" class="footer__link">{{ t('footer.svc_skin') }}</a>
      </div>

      <div class="footer__col">
        <p class="footer__col-title">{{ t('footer.contacts_title') }}</p>
        <div class="footer__contact">
          <span class="footer__contact-icon">📍</span>
          <span>{{ t('footer.address') }}</span>
        </div>
        <div class="footer__contact">
          <span class="footer__contact-icon">📞</span>
          <span>+7 777 000 0000</span>
        </div>
        <div class="footer__contact">
          <span class="footer__contact-icon">✉️</span>
          <span>beauty@dana.kz</span>
        </div>
        <div class="footer__contact">
          <span class="footer__contact-icon">🕐</span>
          <span>{{ t('footer.hours') }}</span>
        </div>
      </div>

      <div class="footer__col">
        <p class="footer__col-title">{{ t('footer.socials_title') }}</p>
        <a href="#" class="footer__social-link">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="2" y="2" width="20" height="20" rx="5"/>
            <circle cx="12" cy="12" r="4"/>
            <circle cx="17.5" cy="6.5" r="1" fill="currentColor"/>
          </svg>
          Instagram
        </a>
        <a href="#" class="footer__social-link">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07A19.5 19.5 0 013.07 9.8a19.79 19.79 0 01-3.07-8.63A2 2 0 012 .18h3a2 2 0 012 1.72c.127.96.361 1.903.7 2.81a2 2 0 01-.45 2.11L6.91 7.91a16 16 0 006.18 6.18l1.29-1.29a2 2 0 012.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0122 16.92z"/>
          </svg>
          WhatsApp
        </a>
        <a href="#" class="footer__social-link">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 2H3v16h5v4l4-4h5l4-4V2z"/>
          </svg>
          Telegram
        </a>
        <a href="#" class="footer__social-link">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="10" r="4"/>
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/>
          </svg>
          2GIS
        </a>
      </div>

    </div>

    <div class="footer__divider" />

    <!-- нижняя строка как у open-design -->
    <div class="footer__bottom page-container">
      <p class="footer__copy">© {{ year }} Beauty Dana. {{ t('footer.rights') }}</p>
      <p class="footer__made">RESEARCH · BEAUTY · CARE · REPEAT</p>
      <p class="footer__origin">MADE IN YESIK · KZ</p>
    </div>

  </footer>
</template>

<style scoped>
.footer {
  background: var(--bg-dark);
  color: rgba(240,236,228,0.6);
  margin-top: auto;
}

/* верхняя секция с лого */
.footer__hero { padding: 60px 0 40px; }
.footer__hero-inner {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}
.footer__logo {
  font-family: var(--font-serif);
  font-size: clamp(40px, 6vw, 64px);
  font-weight: 900;
  color: rgba(240,236,228,0.9);
  letter-spacing: -0.02em;
  line-height: 1;
  margin-bottom: 10px;
}
.footer__sub {
  font-family: var(--font);
  font-size: 10px;
  font-weight: 500;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: rgba(240,236,228,0.35);
}
.footer__cta {
  font-family: var(--font);
  font-size: 13px;
  font-weight: 500;
  letter-spacing: 0.06em;
  color: rgba(240,236,228,0.5);
  text-decoration: none;
  border-bottom: 1px solid rgba(240,236,228,0.15);
  padding-bottom: 2px;
  transition: all 0.2s;
  flex-shrink: 0;
  margin-bottom: 14px;
}
.footer__cta:hover {
  color: var(--accent);
  border-color: var(--accent);
}

/* разделитель */
.footer__divider {
  height: 1px;
  background: rgba(240,236,228,0.08);
}

/* основная сетка */
.footer__main {
  display: grid;
  grid-template-columns: 1.5fr 1.5fr 1.5fr 1fr;
  gap: 40px;
  padding-top: 48px;
  padding-bottom: 48px;
}

.footer__col-title {
  font-family: var(--font);
  font-size: 9px;
  font-weight: 600;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: rgba(240,236,228,0.3);
  margin-bottom: 20px;
}

.footer__link {
  display: block;
  font-family: var(--font);
  font-size: 13px;
  color: rgba(240,236,228,0.55);
  text-decoration: none;
  margin-bottom: 10px;
  transition: color 0.15s;
  letter-spacing: 0.01em;
}
.footer__link:hover { color: rgba(240,236,228,0.9); }

.footer__contact {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: rgba(240,236,228,0.55);
  margin-bottom: 10px;
}
.footer__contact-icon { font-size: 12px; }

.footer__social-link {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: rgba(240,236,228,0.55);
  text-decoration: none;
  margin-bottom: 12px;
  transition: color 0.15s;
}
.footer__social-link:hover { color: var(--accent); }

/* нижняя строка */
.footer__bottom {
  padding-top: 20px;
  padding-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}
.footer__copy,
.footer__made,
.footer__origin {
  font-family: var(--font);
  font-size: 10px;
  font-weight: 500;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: rgba(240,236,228,0.25);
}

@media (max-width: 900px) {
  .footer__main { grid-template-columns: 1fr 1fr; }
  .footer__hero-inner { flex-direction: column; align-items: flex-start; gap: 20px; }
  .footer__bottom { justify-content: flex-start; }
}
@media (max-width: 480px) {
  .footer__main { grid-template-columns: 1fr; }
  .footer__logo { font-size: 36px; }
}
</style>