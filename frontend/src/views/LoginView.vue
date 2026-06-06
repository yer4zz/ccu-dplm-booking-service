<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n }      from 'vue-i18n'
import { Check, ArrowRight } from 'lucide-vue-next'

const router = useRouter()
const route  = useRoute()
const auth   = useAuthStore()
const { t }  = useI18n()

const rezhim   = ref<'login' | 'register'>('login')
const email    = ref('')
const password = ref('')
const fullName = ref('')
const oshibka  = ref('')
const zagruzka = ref(false)

async function submit() {
  oshibka.value  = ''
  zagruzka.value = true
  try {
    if (rezhim.value === 'login') {
      await auth.login(email.value, password.value)
    } else {
      await auth.register(email.value, password.value, fullName.value)
    }
    const redirect = route.query.redirect as string
    router.push(redirect || (auth.isMaster ? '/master' : auth.isAdmin ? '/admin' : '/'))
  } catch (e: any) {
    oshibka.value = e.response?.data?.error_description ?? t('auth.error')
  } finally {
    zagruzka.value = false
  }
}
</script>

<template>
  <div class="login-page">

    <div class="login-panel dot-grid">
      <div class="login-panel__inner">
        <div class="section-label" style="margin-bottom:32px">BEAUTY DANA</div>
        <h2 class="login-panel__title">
          Красота<br/>
          <em>{{ t('login.hero_italic') }}</em><br/>
          здесь<span class="red-dot">.</span>
        </h2>
        <p class="login-panel__sub">
          Записывайтесь к лучшим мастерам онлайн в любое время суток.
        </p>
        <div class="login-panel__stats">
          <div class="login-panel__stat">
            <span class="login-panel__stat-num">30%</span>
            <span class="login-panel__stat-label">{{ t('login.perk1_label') }}</span>
          </div>
          <div class="login-panel__stat">
            <span class="login-panel__stat-num">24/7</span>
            <span class="login-panel__stat-label">{{ t('login.perk2_label') }}</span>
          </div>
        </div>
      </div>
      <div class="login-panel__deco">
        <div class="login-panel__circle" />
        <p class="login-panel__deco-text">BEAUTY · CARE · STYLE</p>
      </div>
    </div>

    <div class="login-form-wrap">
      <div class="login-form-inner">

        <div class="login-mobile-head">
          <p class="login-mobile-head__logo">Beauty Dana</p>
          <p class="login-mobile-head__sub">{{ t('login.mobile_sub') }}</p>
          <div class="login-mobile-head__perks">
            <span><Check :size="14" /> {{ t('login.mobile_perk1') }}</span>
            <span class="login-perk-dot">·</span>
            <span><Check :size="14" /> {{ t('login.mobile_perk2') }}</span>
          </div>
        </div>

        <div class="login-tabs">
          <button
            class="login-tab"
            :class="{ 'login-tab--active': rezhim === 'login' }"
            @click="rezhim = 'login'"
          >{{ t('auth.login') }}</button>
          <button
            class="login-tab"
            :class="{ 'login-tab--active': rezhim === 'register' }"
            @click="rezhim = 'register'"
          >{{ t('auth.register') }}</button>
        </div>

        <div class="login-head">
          <div class="section-label" style="margin-bottom:8px">
            {{ rezhim === 'login' ? t('login.mode_login') : t('login.mode_register') }}
          </div>
          <h1 class="t-h2">
            {{ rezhim === 'login' ? t('auth.login') : t('auth.register') }}<span class="red-dot">.</span>
          </h1>
        </div>

        <div class="login-fields">
          <div v-if="rezhim === 'register'" class="field">
            <label class="field-label">{{ t('auth.name') }}</label>
            <input
              v-model="fullName"
              type="text"
              class="field-input"
              placeholder="Имя Фамилия"
              autofocus
            />
          </div>
          <div class="field">
            <label class="field-label">{{ t('auth.email') }}</label>
            <input
              v-model="email"
              type="email"
              class="field-input"
              placeholder="your@gmail.com"
              :autofocus="rezhim === 'login'"
              autocomplete="email"
            />
          </div>
          <div class="field">
            <label class="field-label">{{ t('auth.password') }}</label>
            <input
              v-model="password"
              type="password"
              class="field-input"
              placeholder="••••••••"
              autocomplete="current-password"
              @keyup.enter="submit"
            />
          </div>
        </div>

        <div v-if="oshibka" class="login-error">{{ oshibka }}</div>

        <button
          class="btn btn-primary btn-full"
          style="margin-top:8px"
          :disabled="zagruzka"
          @click="submit"
        >
          <span v-if="zagruzka" class="btn__spinner" />
          {{ rezhim === 'login' ? t('auth.login_btn') : t('auth.register_btn') }} <ArrowRight :size="16" />
        </button>

        <p class="login-switch">
          <template v-if="rezhim === 'login'">
            {{ t('auth.no_account') }}
            <button class="login-switch__btn" @click="rezhim = 'register'">
              {{ t('auth.register') }} <ArrowRight :size="14"/>
            </button>
          </template>
          <template v-else>
            {{ t('auth.has_account') }}
            <button class="login-switch__btn" @click="rezhim = 'login'">
              {{ t('auth.login') }} <ArrowRight :size="14"/>
            </button>
          </template>
        </p>

      </div>
    </div>

  </div>
</template>

<style scoped>
.login-page {
  min-height: calc(100vh - 56px);
  display: grid;
  grid-template-columns: 1fr 1fr;
}

.login-panel {
  background: var(--bg-muted);
  border-right: 1px solid var(--border);
  padding: 60px 48px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}
.login-panel__inner { position: relative; z-index: 1; }
.login-panel__title {
  font-family: var(--font-serif);
  font-size: clamp(36px, 4vw, 52px);
  font-weight: 900;
  line-height: 0.95;
  letter-spacing: -0.02em;
  color: var(--text);
  margin-bottom: 24px;
}
.login-panel__title em { font-style: italic; font-weight: 700; color: var(--text-2); }
.login-panel__sub {
  font-size: 14px; color: var(--text-2);
  line-height: 1.7; max-width: 320px; margin-bottom: 40px;
}
.login-panel__stats {
  display: flex; gap: 32px;
  padding-top: 28px; border-top: 1px solid var(--border);
}
.login-panel__stat      { display: flex; flex-direction: column; gap: 4px; }
.login-panel__stat-num  {
  font-family: var(--font-serif); font-size: 28px;
  font-weight: 700; color: var(--text); line-height: 1;
}
.login-panel__stat-label {
  font-size: 10px; font-weight: 500;
  letter-spacing: 0.08em; text-transform: uppercase; color: var(--text-3);
}
.login-panel__deco {
  display: flex; align-items: center;
  gap: 16px; position: relative; z-index: 1;
}
.login-panel__circle {
  width: 48px; height: 48px;
  border: 1px solid var(--border-strong);
  border-radius: 50%; flex-shrink: 0;
}
.login-panel__deco-text {
  font-size: 9px; font-weight: 500;
  letter-spacing: 0.16em; text-transform: uppercase; color: var(--text-3);
}

.login-form-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 48px;
  background: var(--bg);
}
.login-form-inner { width: 100%; max-width: 400px; }

.login-mobile-head { display: none; }

.login-tabs {
  display: flex;
  border-bottom: 1px solid var(--border);
  margin-bottom: 32px;
}
.login-tab {
  padding: 8px 20px 10px;
  background: none; border: none;
  font-family: var(--font); font-size: 12px; font-weight: 500;
  letter-spacing: 0.06em; text-transform: uppercase;
  color: var(--text-3); cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px; transition: all 0.15s;
}
.login-tab--active { color: var(--text); border-bottom-color: var(--accent); }
.login-tab:hover:not(.login-tab--active) { color: var(--text-2); }

.login-head   { margin-bottom: 28px; }
.login-fields { display: flex; flex-direction: column; gap: 18px; }

.login-error {
  margin-top: 14px; padding: 10px 14px;
  background: var(--danger-bg); color: var(--danger);
  border-left: 2px solid var(--danger);
  border-radius: var(--radius-sm); font-size: 13px;
}

.login-switch {
  margin-top: 20px; text-align: center;
  font-size: 13px; color: var(--text-3);
}
.login-switch__btn {
  background: none; border: none;
  font-family: var(--font); font-size: 13px; font-weight: 500;
  color: var(--accent); cursor: pointer; padding: 0;
  transition: opacity 0.15s;
  display: inline-flex; align-items: center; gap: 6px;
}
.login-switch__btn:hover { opacity: 0.7; }

@media (max-width: 600px) {
  .login-page {
    grid-template-columns: 1fr;
    min-height: calc(100vh - 56px);
  }

  .login-panel { display: none; }

  .login-form-wrap {
    align-items: flex-start;
    padding: 24px 20px 40px;
  }
  .login-form-inner { max-width: 100%; }

  .login-mobile-head {
    display: block;
    text-align: center;
    padding: 24px 0 28px;
    margin-bottom: 8px;
    border-bottom: 1px solid var(--border);
    margin-bottom: 24px;
  }
  .login-mobile-head__logo {
    font-family: var(--font-serif);
    font-size: 32px;
    font-weight: 900;
    color: var(--text);
    letter-spacing: -0.02em;
    line-height: 1;
    margin-bottom: 6px;
  }
  .login-mobile-head__sub {
    font-family: var(--font);
    font-size: 9px;
    font-weight: 500;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--text-3);
    margin-bottom: 14px;
  }
  .login-mobile-head__perks {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    font-size: 11px;
    color: var(--text-3);
    font-weight: 500;
    flex-wrap: wrap;
  }
  .login-perk-dot { opacity: 0.4; }

  .login-tabs { margin-bottom: 24px; }
  .login-tab  { flex: 1; text-align: center; padding: 10px 8px 12px; font-size: 11px; }

  .login-head { margin-bottom: 20px; }

  .login-fields { gap: 14px; }

  .field-input {
    height: 48px;
    font-size: 16px;
    padding: 12px 16px;
  }

  .btn-primary { height: 50px; font-size: 15px; }

  .login-switch { font-size: 14px; margin-top: 24px; }
}
</style>