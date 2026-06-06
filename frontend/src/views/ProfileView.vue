<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore }   from '@/stores/auth'
import { loyaltyApi }     from '@/api/loyalty'
import { api }            from '@/api'
import { useI18n }        from 'vue-i18n'
import type { LoyaltyAccount, PointTransaction, BeautyStreak } from '@/types'
import { Star, ArrowRight } from 'lucide-vue-next'

const auth    = useAuthStore()
const { t }   = useI18n()
const zagruzka = ref(false)
const sohranenie = ref(false)
const uspeh      = ref(false)
const editMode   = ref(false)

const profileData = ref({ full_name: '', phone: '' })
const editName    = ref('')
const editPhone   = ref('')

const loyalty = ref<LoyaltyAccount | null>(null)
const txns    = ref<PointTransaction[]>([])
const streak  = ref<BeautyStreak | null>(null)

onMounted(async () => {
  zagruzka.value = true
  try {
    const [acc, tr, str, meRes] = await Promise.all([
      loyaltyApi.getAccount().catch(() => null),
      loyaltyApi.getTransactions().catch(() => []),
      loyaltyApi.getStreak().catch(() => null),
      api.get('/me'),
    ])
    loyalty.value = acc
    txns.value    = tr ?? []
    streak.value  = str

    profileData.value = {
      full_name: meRes.data.profile_full_name
                 || meRes.data.user_metadata?.full_name
                 || auth.user?.full_name || '',
      phone:     meRes.data.profile_phone
                 || auth.user?.phone || '',
    }
  } finally {
    zagruzka.value = false
  }
})

function startEdit() {
  editName.value  = profileData.value.full_name
  editPhone.value = profileData.value.phone
  editMode.value  = true
}
function cancelEdit() { editMode.value = false }

async function save() {
  sohranenie.value = true
  uspeh.value      = false
  try {
    await api.patch('/profile', {
      full_name: editName.value,
      phone:     editPhone.value,
    })
    profileData.value.full_name = editName.value
    profileData.value.phone     = editPhone.value
    if (auth.user) auth.user.full_name = editName.value
    editMode.value = false
    uspeh.value    = true
    setTimeout(() => { uspeh.value = false }, 3000)
  } finally {
    sohranenie.value = false
  }
}

const levelNames: Record<string, string> = {
  bronze: 'Bronze', silver: 'Silver', gold: 'Gold', platinum: 'Platinum'
}
const levelEmoji: Record<string, string> = {
  bronze: '🥉', silver: '🥈', gold: '🥇', platinum: '💎'
}

function streakPhrase(): string {
  if (!streak.value) return ''
  const n = streak.value.current_streak
  if (n >= 16) return t('profile.streak_phrases.platinum')
  if (n >= 9)  return t('profile.streak_phrases.gold')
  if (n >= 4)  return t('profile.streak_phrases.silver').replace('{n}', String(n))
  if (n > 0)   return t('profile.streak_phrases.bronze')
  return t('profile.streak_phrases.start')
}

function txnIcon(type: string)  { return type === 'earn' ? '+' : '−' }
function txnColor(type: string) { return type === 'earn' ? 'var(--success)' : 'var(--danger)' }
function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'short', year: 'numeric'
  })
}

function getInitials(): string {
  const name = profileData.value.full_name || auth.user?.email || 'U'
  return name[0].toUpperCase()
}
</script>

<template>
  <div class="page-container" style="max-width:760px; padding-top:48px; padding-bottom:80px">

    <div class="profile-head">
      <div class="section-label" style="margin-bottom:8px">{{ t('profile.eyebrow') }}</div>
      <h1 class="t-h1">{{ t('profile.title') }}<span class="red-dot">.</span></h1>
    </div>

    <div class="deco-line" />

    <div v-if="zagruzka" style="text-align:center;padding:80px;color:var(--text-3)">
      {{ t('common.loading') }}
    </div>

    <template v-else>

      <section class="profile-section">
        <div class="profile-card">

          <div v-if="!editMode" class="profile-view">
            <div class="profile-avatar">{{ getInitials() }}</div>
            <div class="profile-info">
              <h2 class="profile-name">
                {{ profileData.full_name || t('profile.name_missing') }}
              </h2>
              <p class="profile-email">{{ auth.user?.email }}</p>
              <p class="profile-phone" v-if="profileData.phone">
                📱 {{ profileData.phone }}
              </p>
              <p class="profile-phone profile-phone--empty" v-else>
                {{ t('profile.no_phone') }}
              </p>
            </div>
            <button class="btn btn-outline btn-sm" @click="startEdit">
              {{ t('profile.edit') }}
            </button>
          </div>

          <div v-else class="profile-edit">
            <div class="section-label" style="margin-bottom:16px">
              {{ t('profile.edit_title') }}
            </div>
            <div style="display:flex;flex-direction:column;gap:16px">
              <div class="field">
                <label class="field-label">{{ t('profile.email') }}</label>
                <input
                  :value="auth.user?.email"
                  class="field-input"
                  disabled
                  style="background:var(--bg-muted);color:var(--text-3);cursor:not-allowed"
                />
              </div>
              <div class="field">
                <label class="field-label">{{ t('profile.name') }}</label>
                <input v-model="editName" type="text" class="field-input" autofocus />
              </div>
              <div class="field">
                <label class="field-label">{{ t('profile.phone') }}</label>
                <input v-model="editPhone" type="tel" class="field-input"
                  :placeholder="t('profile.phone_placeholder')" />
              </div>
            </div>
            <div style="display:flex;align-items:center;gap:12px;margin-top:20px">
              <button class="btn btn-primary btn-sm" :disabled="sohranenie" @click="save">
                <span v-if="sohranenie" class="btn__spinner" />
                {{ t('profile.save') }} <ArrowRight :size="14" />
              </button>
              <button class="btn btn-ghost btn-sm" @click="cancelEdit">{{ t('profile.cancel') }}</button>
              <Transition name="fade">
                <span v-if="uspeh" style="font-size:12px;color:var(--success);font-weight:500;letter-spacing:0.04em">
                  {{ t('profile.saved') }}
                </span>
              </Transition>
            </div>
          </div>
        </div>
      </section>

      <section class="profile-section" v-if="streak">
        <div class="section-label" style="margin-bottom:20px">{{ t('profile.streak_title') }}</div>

        <div class="streak-card">
          <div class="streak-card__top">
            <div class="streak-fire">
              <span class="streak-fire__emoji">🔥</span>
              <div>
                <div style="display:flex;align-items:baseline;gap:10px">
                  <span class="streak-num">{{ streak.current_streak }}</span>
                  <span class="streak-unit">{{ t('profile.streak_visits') }}</span>
                </div>
                <p style="font-size:11px;color:var(--text-3);margin-top:4px;letter-spacing:0.04em">
                  {{ t('profile.streak_record') }}: {{ streak.longest_streak }}
                  <span v-if="streak.streak_discount > 0" style="color:var(--accent);margin-left:8px">
                    · {{ t('profile.streak_discount') }} {{ streak.streak_discount }}% {{ t('profile.streak_active') }}
                  </span>
                </p>
              </div>
            </div>
            <div class="streak-level" :class="`streak-level--${streak.level}`">
              {{ levelEmoji[streak.level] }} {{ levelNames[streak.level] }}
            </div>
          </div>

          <div v-if="streak.level !== 'platinum'" class="streak-progress">
            <div class="streak-progress__track">
              <div
                class="streak-progress__fill"
                :style="`width: ${Math.min(streak.progress, 100)}%`"
              />
            </div>
            <p style="font-size:11px;color:var(--text-3);margin-top:6px;letter-spacing:0.04em">
              {{ streak.current_streak }} / {{ streak.next_level_at }}
              {{ t('profile.streak_to_next') }}
            </p>
          </div>

          <p class="streak-phrase">{{ streakPhrase() }}</p>

          <div
            v-if="streak.days_until_deadline !== undefined"
            class="streak-deadline"
            :style="streak.days_until_deadline <= 3
              ? 'color:var(--danger);background:var(--danger-bg)'
              : ''"
          >
            {{ t('profile.streak_next') }}
            <b>{{ streak.days_until_deadline }} {{ t('profile.streak_days') }}</b>
          </div>
        </div>
      </section>

      <section class="profile-section" v-if="loyalty">
        <div class="section-label" style="margin-bottom:20px">{{ t('profile.loyalty_title') }}</div>

        <div class="loyalty-card">
          <div class="loyalty-card__top">
            <div>
              <div style="display:flex;align-items:baseline;gap:10px">
                <span class="loyalty-num">{{ loyalty.balance }}</span>
                <span class="loyalty-unit">{{ t('profile.loyalty_points') }}</span>
              </div>
              <p style="font-size:12px;color:var(--text-3);margin-top:6px;letter-spacing:0.04em">
                {{ t('profile.loyalty_discount_hint', { amt: (loyalty.balance * 10).toLocaleString() + ' ₸' }) }}
                · {{ t('profile.loyalty_earned') }}: {{ loyalty.total_earned }}
              </p>
            </div>
            <div class="loyalty-star"><Star :size="24" /></div>
          </div>

          <div v-if="txns.length > 0" class="loyalty-history">
            <div class="section-label" style="margin-bottom:12px">{{ t('profile.loyalty_history') }}</div>
            <div v-for="tx in txns.slice(0, 5)" :key="tx.id" class="txn-row">
              <span class="txn-amount" :style="`color: ${txnColor(tx.type)}`">
                {{ txnIcon(tx.type) }}{{ Math.abs(tx.amount) }}
              </span>
              <span class="txn-desc">{{ tx.description }}</span>
              <span class="txn-date">{{ formatDate(tx.created_at) }}</span>
            </div>
          </div>
          <p v-else style="font-size:13px;color:var(--text-3);margin-top:16px;font-style:italic">
            {{ t('profile.loyalty_empty') }}
          </p>
        </div>
      </section>

    </template>
  </div>
</template>

<style scoped>
.profile-head  { padding-bottom: 24px; }

.profile-section { margin-bottom: 48px; }

.profile-card {
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 28px;
  background: var(--bg-card);
}

.profile-view {
  display: flex;
  align-items: center;
  gap: 20px;
}
.profile-avatar {
  width: 64px; height: 64px;
  border-radius: 50%;
  background: var(--accent-bg);
  border: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-serif);
  font-size: 26px;
  font-weight: 700;
  color: var(--accent);
  flex-shrink: 0;
}
.profile-info   { flex: 1; }
.profile-name   {
  font-family: var(--font-serif);
  font-size: 22px;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 3px;
}
.profile-email  { font-size: 12px; color: var(--text-3); margin-bottom: 4px; letter-spacing: 0.02em; }
.profile-phone  { font-size: 13px; color: var(--text-2); }
.profile-phone--empty { color: var(--text-3); font-style: italic; }

.profile-edit { animation: fadeUp 0.2s ease; }
@keyframes fadeUp {
  from { opacity: 0; transform: translateY(6px); }
  to   { opacity: 1; transform: translateY(0); }
}

.streak-card {
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 28px;
  background: var(--bg-card);
}
.streak-card__top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}
.streak-fire {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}
.streak-fire__emoji { font-size: 36px; line-height: 1; flex-shrink: 0; }
.streak-num {
  font-family: var(--font-serif);
  font-size: 52px;
  font-weight: 900;
  color: var(--text);
  line-height: 1;
}
.streak-unit { font-size: 14px; color: var(--text-2); }

.streak-level {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0.04em;
  flex-shrink: 0;
}
.streak-level--bronze   { background: #FEF3C7; color: #92400E; border: 1px solid rgba(146,64,14,0.2); }
.streak-level--silver   { background: #F3F4F6; color: #374151; border: 1px solid rgba(55,65,81,0.2); }
.streak-level--gold     { background: #FEF9C3; color: #713F12; border: 1px solid rgba(113,63,18,0.2); }
.streak-level--platinum { background: #EDE9FE; color: #4C1D95; border: 1px solid rgba(76,29,149,0.2); }

.streak-progress { margin-bottom: 16px; }
.streak-progress__track {
  height: 4px;
  background: var(--bg-muted);
  border-radius: 99px;
  overflow: hidden;
}
.streak-progress__fill {
  height: 100%;
  background: linear-gradient(90deg, var(--accent), #F59E0B);
  border-radius: 99px;
  transition: width 0.5s cubic-bezier(0.4,0,0.2,1);
}

.streak-phrase {
  font-size: 13px;
  color: var(--text-2);
  font-style: italic;
  padding: 12px 16px;
  background: var(--bg-muted);
  border-radius: var(--radius-sm);
  border-left: 2px solid var(--accent);
  margin-bottom: 12px;
}
.streak-deadline {
  display: inline-block;
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0.04em;
  color: var(--text-3);
  background: var(--bg-muted);
  padding: 5px 12px;
  border-radius: var(--radius-sm);
}

.loyalty-card {
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 28px;
  background: var(--bg-card);
}
.loyalty-card__top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}
.loyalty-num {
  font-family: var(--font-serif);
  font-size: 52px;
  font-weight: 900;
  color: var(--text);
  line-height: 1;
}
.loyalty-unit { font-size: 14px; color: var(--text-2); }
.loyalty-star { font-size: 36px; }

.loyalty-history {
  border-top: 1px solid var(--border);
  padding-top: 20px;
  margin-top: 4px;
}
.txn-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 10px 0;
  border-bottom: 1px solid var(--border);
}
.txn-row:last-child { border-bottom: none; }
.txn-amount {
  font-family: var(--font-serif);
  font-size: 16px;
  font-weight: 700;
  min-width: 46px;
  flex-shrink: 0;
}
.txn-desc { flex: 1; font-size: 13px; color: var(--text-2); }
.txn-date { font-size: 11px; color: var(--text-3); flex-shrink: 0; letter-spacing: 0.02em; }

@media (max-width: 600px) {
  .profile-view { flex-wrap: wrap; }
  .streak-card__top { flex-direction: column; gap: 16px; }
  .loyalty-card__top { flex-direction: column; gap: 12px; }
}
</style>
