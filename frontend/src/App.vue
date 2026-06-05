<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { useAuthStore }    from '@/stores/auth'
import { useTheme }        from '@/composables/useTheme'
import { useI18n }         from 'vue-i18n'
import { i18n }            from '@/i18n'
import AppFooter           from '@/components/AppFooter.vue'
import ChatBot             from '@/components/ChatBot.vue'
import { useBookingStore } from '@/stores/booking'
import { onMounted }       from 'vue'

const auth   = useAuthStore()
const router = useRouter()
const { theme, toggle } = useTheme()
const { t }  = useI18n()

const langs = ['ru', 'kz', 'en']

const menuOpen = ref(false)
function toggleMenu() { menuOpen.value = !menuOpen.value }
function closeMenu()  { menuOpen.value = false }

router.afterEach(() => closeMenu())

const metaDate = computed(() => {
  const d   = new Date()
  const vol = String(d.getMonth() + 1).padStart(2, '0')
  return `VOL. ${vol} · ${d.getFullYear()}`
})

const bookingStore = useBookingStore()
onMounted(() => { bookingStore.init().catch(() => {}) })

function setLang(lang: string) {
  i18n.global.locale.value = lang as 'ru' | 'kz' | 'en'
  localStorage.setItem('lang', lang)
}

function logout() {
  auth.logout()
  router.push('/login')
  closeMenu()
}
</script>

<template>
  <div class="app-layout">

    <div class="meta-bar">
      <span class="meta-bar__left">BEAUTY DANA</span>
      <span class="meta-bar__center">{{ metaDate }}</span>
      <div class="meta-bar__right">
        <button
          v-for="lang in langs"
          :key="lang"
          class="meta-lang"
          :class="{ 'meta-lang--active': i18n.global.locale.value === lang }"
          @click="setLang(lang)"
        >{{ lang.toUpperCase() }}</button>
      </div>
    </div>

    <header class="navbar">
      <div class="navbar__inner">

        <RouterLink to="/" class="navbar__logo" @click="closeMenu">
          <span class="navbar__logo-name">Beauty Dana</span>
          <span class="navbar__logo-sub">СТУДИЯ КРАСОТЫ</span>
        </RouterLink>

        <nav class="navbar__nav">
          <RouterLink to="/book"    class="navbar__link">{{ t('nav.book') }}</RouterLink>
          <RouterLink to="/gallery" class="navbar__link">{{ t('nav.gallery') }}</RouterLink>
          <RouterLink to="/masters" class="navbar__link">{{ t('nav.masters') }}</RouterLink>
          <RouterLink to="/about"   class="navbar__link">{{ t('nav.about') }}</RouterLink>
          <RouterLink
            v-if="auth.isLoggedIn && !auth.isMaster && !auth.isAdmin"
            to="/my" class="navbar__link"
          >{{ t('nav.myBookings') }}</RouterLink>
          <RouterLink
            v-if="auth.isLoggedIn && !auth.isMaster && !auth.isAdmin"
            to="/profile" class="navbar__link"
          >{{ t('nav.profile') }}</RouterLink>
          <RouterLink v-if="auth.isMaster" to="/master" class="navbar__link">
            {{ t('nav.cabinet') }}
          </RouterLink>
          <RouterLink v-if="auth.isAdmin" to="/admin" class="navbar__link">
            {{ t('nav.admin') }}
          </RouterLink>
        </nav>

        <div class="navbar__actions">
          <button class="navbar__theme-btn" @click="toggle">
            <svg v-if="theme === 'dark'" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5"/>
              <line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
              <line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/>
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
            </svg>
            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
            </svg>
          </button>
          <RouterLink v-if="!auth.isLoggedIn" to="/login" class="btn btn-outline btn-sm">
            {{ t('nav.login') }}
          </RouterLink>
          <button v-else class="btn btn-ghost btn-sm" @click="logout">
            {{ t('nav.logout') }}
          </button>
        </div>

        <div class="navbar__mobile-actions">
          <button class="navbar__theme-btn" @click="toggle">
            <svg v-if="theme === 'dark'" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5"/>
              <line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/>
            </svg>
            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
            </svg>
          </button>
          <button
            class="burger-btn"
            :class="{ 'burger-btn--open': menuOpen }"
            @click="toggleMenu"
            aria-label="Меню"
          >
            <span class="burger-line" />
            <span class="burger-line" />
            <span class="burger-line" />
          </button>
        </div>

      </div>
    </header>

    <Transition name="menu-slide">
      <div v-if="menuOpen" class="mobile-overlay" @click.self="closeMenu">
        <div class="mobile-menu">

          <div class="mobile-menu__head">
            <div class="mobile-menu__logo">
              <span class="navbar__logo-name">Beauty Dana</span>
              <span class="navbar__logo-sub">СТУДИЯ КРАСОТЫ</span>
            </div>
            <button class="mobile-menu__close" @click="closeMenu">✕</button>
          </div>

          <nav class="mobile-menu__nav">
            <RouterLink to="/book" class="mobile-link" @click="closeMenu">
              <span class="mobile-link__num">01</span>
              {{ t('nav.book') }}
              <span class="mobile-link__arrow">→</span>
            </RouterLink>
            <RouterLink to="/gallery" class="mobile-link" @click="closeMenu">
              <span class="mobile-link__num">02</span>
              {{ t('nav.gallery') }}
              <span class="mobile-link__arrow">→</span>
            </RouterLink>
            <RouterLink to="/masters" class="mobile-link" @click="closeMenu">
              <span class="mobile-link__num">03</span>
              {{ t('nav.masters') }}
              <span class="mobile-link__arrow">→</span>
            </RouterLink>
            <RouterLink to="/about" class="mobile-link" @click="closeMenu">
              <span class="mobile-link__num">04</span>
              {{ t('nav.about') }}
              <span class="mobile-link__arrow">→</span>
            </RouterLink>
            <template v-if="auth.isLoggedIn && !auth.isMaster && !auth.isAdmin">
              <RouterLink to="/my" class="mobile-link" @click="closeMenu">
                <span class="mobile-link__num">05</span>
                {{ t('nav.myBookings') }}
                <span class="mobile-link__arrow">→</span>
              </RouterLink>
              <RouterLink to="/profile" class="mobile-link" @click="closeMenu">
                <span class="mobile-link__num">06</span>
                {{ t('nav.profile') }}
                <span class="mobile-link__arrow">→</span>
              </RouterLink>
            </template>
            <RouterLink v-if="auth.isMaster" to="/master" class="mobile-link" @click="closeMenu">
              <span class="mobile-link__num">05</span>
              {{ t('nav.cabinet') }}
              <span class="mobile-link__arrow">→</span>
            </RouterLink>
            <RouterLink v-if="auth.isAdmin" to="/admin" class="mobile-link" @click="closeMenu">
              <span class="mobile-link__num">05</span>
              {{ t('nav.admin') }}
              <span class="mobile-link__arrow">→</span>
            </RouterLink>
          </nav>

          <div class="mobile-menu__footer">
            <div class="mobile-menu__langs">
              <button
                v-for="lang in langs"
                :key="lang"
                class="mobile-lang"
                :class="{ 'mobile-lang--active': i18n.global.locale.value === lang }"
                @click="setLang(lang)"
              >{{ lang.toUpperCase() }}</button>
            </div>
            <RouterLink
              v-if="!auth.isLoggedIn"
              to="/login"
              class="btn btn-primary btn-full"
              @click="closeMenu"
            >{{ t('nav.login') }} →</RouterLink>
            <button v-else class="btn btn-outline btn-full" @click="logout">
              {{ t('nav.logout') }}
            </button>
          </div>

        </div>
      </div>
    </Transition>

    <main class="main-content">
      <RouterView />
    </main>

    <AppFooter />
    <ChatBot />
  </div>
</template>

<style>
.app-layout  { min-height: 100vh; display: flex; flex-direction: column; }
.main-content { flex: 1; }
</style>

<style scoped>
.meta-bar {
  height: 32px;
  background: var(--bg-dark);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
}
.meta-bar__left,
.meta-bar__center,
.meta-bar__right {
  font-family: var(--font);
  font-size: 10px;
  font-weight: 500;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: rgba(240,236,228,0.45);
}
.meta-bar__right { display: flex; align-items: center; gap: 2px; }
.meta-lang {
  background: none; border: none;
  font-family: var(--font); font-size: 10px; font-weight: 500;
  letter-spacing: 0.1em; color: rgba(240,236,228,0.35);
  cursor: pointer; padding: 4px 6px; transition: color 0.15s;
}
.meta-lang:hover { color: rgba(240,236,228,0.8); }
.meta-lang--active { color: rgba(240,236,228,0.9); }
.meta-lang:not(:last-child)::after { content: ' ·'; color: rgba(240,236,228,0.25); margin-left: 2px; }

.navbar {
  background: var(--bg);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 200;
}
.navbar__inner {
  max-width: 1200px; margin: 0 auto;
  padding: 0 40px; height: 64px;
  display: flex; align-items: center; gap: 40px;
}
.navbar__logo {
  display: flex; flex-direction: column; gap: 1px;
  text-decoration: none; flex-shrink: 0;
}
.navbar__logo-name {
  font-family: var(--font-serif); font-size: 18px; font-weight: 700;
  color: var(--text); letter-spacing: -0.01em; line-height: 1;
}
.navbar__logo-sub {
  font-family: var(--font); font-size: 8px; font-weight: 500;
  letter-spacing: 0.12em; text-transform: uppercase;
  color: var(--text-3); line-height: 1;
}
.navbar__nav { display: flex; align-items: center; flex: 1; }
.navbar__link {
  padding: 8px 14px;
  font-family: var(--font); font-size: 13px; font-weight: 400;
  color: var(--text-3); text-decoration: none;
  letter-spacing: 0.02em; position: relative; transition: color 0.15s;
  white-space: nowrap;
}
.navbar__link::after {
  content: ''; position: absolute;
  bottom: 4px; left: 14px; right: 14px;
  height: 1px; background: var(--accent);
  transform: scaleX(0); transform-origin: left;
  transition: transform 0.2s ease;
}
.navbar__link:hover { color: var(--text); }
.navbar__link:hover::after { transform: scaleX(1); }
.navbar__link.router-link-active { color: var(--text); }
.navbar__link.router-link-active::after { transform: scaleX(1); }

.navbar__actions {
  display: flex; align-items: center; gap: 8px; margin-left: auto;
}
.navbar__mobile-actions {
  display: none; align-items: center; gap: 8px; margin-left: auto;
}
.navbar__theme-btn {
  width: 32px; height: 32px;
  display: flex; align-items: center; justify-content: center;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm); color: var(--text-3);
  cursor: pointer; transition: all 0.15s; flex-shrink: 0;
}
.navbar__theme-btn:hover { border-color: var(--border-strong); color: var(--text); }

.burger-btn {
  width: 40px; height: 40px;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  gap: 5px; cursor: pointer; padding: 0;
  transition: border-color 0.15s;
}
.burger-btn:hover { border-color: var(--border-strong); }
.burger-line {
  display: block; width: 18px; height: 1.5px;
  background: var(--text); border-radius: 2px;
  transition: all 0.25s cubic-bezier(0.4,0,0.2,1);
  transform-origin: center;
}
.burger-btn--open .burger-line:nth-child(1) { transform: translateY(6.5px) rotate(45deg); }
.burger-btn--open .burger-line:nth-child(2) { opacity: 0; transform: scaleX(0); }
.burger-btn--open .burger-line:nth-child(3) { transform: translateY(-6.5px) rotate(-45deg); }

.mobile-overlay {
  position: fixed; inset: 0; z-index: 300;
  background: rgba(26,24,20,0.5);
  backdrop-filter: blur(4px);
}
.mobile-menu {
  position: absolute; top: 0; right: 0; bottom: 0;
  width: min(320px, 88vw);
  background: var(--bg);
  border-left: 1px solid var(--border);
  display: flex; flex-direction: column;
  overflow-y: auto;
}
.mobile-menu__head {
  display: flex; align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.mobile-menu__logo { display: flex; flex-direction: column; gap: 2px; }
.mobile-menu__close {
  width: 36px; height: 36px;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-2); font-size: 14px; cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: all 0.15s; flex-shrink: 0;
}
.mobile-menu__close:hover { border-color: var(--accent); color: var(--accent); }

.mobile-menu__nav { flex: 1; padding: 4px 0; }
.mobile-link {
  display: flex; align-items: center; gap: 14px;
  padding: 16px 24px;
  font-family: var(--font-serif); font-size: 20px; font-weight: 700;
  color: var(--text); text-decoration: none;
  border-bottom: 1px solid var(--border);
  transition: all 0.15s;
}
.mobile-link:hover           { background: var(--bg-hover); padding-left: 30px; }
.mobile-link.router-link-active { color: var(--accent); }
.mobile-link__num {
  font-family: var(--font); font-size: 10px; font-weight: 500;
  letter-spacing: 0.1em; color: var(--text-3); min-width: 20px;
}
.mobile-link__arrow {
  margin-left: auto; font-size: 16px; color: var(--text-3);
  transition: transform 0.2s;
}
.mobile-link:hover .mobile-link__arrow { transform: translateX(4px); color: var(--accent); }

.mobile-menu__footer {
  padding: 20px 24px; border-top: 1px solid var(--border);
  display: flex; flex-direction: column; gap: 12px; flex-shrink: 0;
}
.mobile-menu__langs { display: flex; gap: 8px; }
.mobile-lang {
  flex: 1; padding: 8px;
  background: transparent; border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-family: var(--font); font-size: 11px; font-weight: 500;
  letter-spacing: 0.1em; color: var(--text-3);
  cursor: pointer; transition: all 0.15s; text-align: center;
}
.mobile-lang:hover          { border-color: var(--border-strong); color: var(--text); }
.mobile-lang--active        { border-color: var(--accent); color: var(--accent); background: var(--accent-bg); }

.menu-slide-enter-active,
.menu-slide-leave-active { transition: opacity 0.2s ease; }
.menu-slide-enter-from,
.menu-slide-leave-to     { opacity: 0; }
.menu-slide-enter-active .mobile-menu { transition: transform 0.25s cubic-bezier(0.4,0,0.2,1); }
.menu-slide-leave-active .mobile-menu { transition: transform 0.2s  cubic-bezier(0.4,0,0.2,1); }
.menu-slide-enter-from .mobile-menu   { transform: translateX(100%); }
.menu-slide-leave-to   .mobile-menu   { transform: translateX(100%); }

@media (max-width: 900px) {
  .meta-bar              { display: none; }
  .navbar__nav           { display: none; }
  .navbar__actions       { display: none; }
  .navbar__mobile-actions { display: flex; }
  .navbar__inner         { padding: 0 16px; gap: 12px; height: 56px; }
}
</style>