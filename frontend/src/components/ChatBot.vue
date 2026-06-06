<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { X } from 'lucide-vue-next'

const auth = useAuthStore()

interface Message {
  role: 'user' | 'bot'
  text: string
  time: string
  quickReplies?: string[]
  typing?: boolean
  source?: string
}

interface HistoryItem {
  role: 'user' | 'bot'
  message: string
}

const STORAGE_KEY = 'beauty_dana_chat'
const HISTORY_TTL = 24 * 60 * 60 * 1000 // 24 часа

const isOpen     = ref(false)
const messages   = ref<Message[]>([])
const history    = ref<HistoryItem[]>([])
const input      = ref('')
const loading    = ref(false)
const messagesEl = ref<HTMLElement | null>(null)
const sessionId      = 'session_' + Math.random().toString(36).slice(2)
const currentSession = ref<any>({ step: 'idle' })
const statusText = ref('Онлайн помощник')
const statusType = ref<'online' | 'typing' | 'thinking'>('online')

function getTime() {
  return new Date().toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const data = JSON.parse(saved)
      if (data.timestamp && Date.now() - data.timestamp < HISTORY_TTL) {
        messages.value = data.messages ?? []
        history.value  = data.history ?? []
      } else {
        localStorage.removeItem(STORAGE_KEY)
      }
    }
  } catch {}
})

function saveToStorage() {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({
      messages:  messages.value.filter(m => !m.typing),
      history:   history.value,
      timestamp: Date.now(),
    }))
  } catch {}
}

function clearChat() {
  messages.value = []
  history.value  = []
  localStorage.removeItem(STORAGE_KEY)
  nextTick(() => {
    messages.value.push({
      role: 'bot',
      text: 'История очищена. Чем могу помочь?',
      time: getTime(),
      quickReplies: ['Записаться', 'Услуги и цены', 'Наши мастера'],
    })
  })
}

// ── Открытие чата ─────────────────────────────────────
function open() {
  isOpen.value = true
  if (messages.value.length === 0) {
    messages.value.push({
      role: 'bot',
      text: 'Привет! 👋 Я ИИ-помощник Beauty Dana. Могу рассказать об услугах, ценах, мастерах или помочь с записью. Что вас интересует?',
      time: getTime(),
      quickReplies: ['Записаться', 'Услуги и цены', 'Наши мастера', 'Скидки'],
    })
  }
  nextTick(scrollToBottom)
}

// ── Отправка сообщения ────────────────────────────────
async function send(text?: string) {
  const msgText = (text ?? input.value).trim()
  if (!msgText || loading.value) return

  // Спецкнопка — войти
  if (msgText === '🔑 Войти') {
    isOpen.value = false
    window.location.href = '/login?redirect=/book'
    return
  }

  messages.value.push({ role: 'user', text: msgText, time: getTime() })
  history.value.push({ role: 'user', message: msgText })
  input.value   = ''
  loading.value = true
  await nextTick(scrollToBottom)

  // Показываем анимацию — сначала "typing", потом если долго — "thinking"
  const typingMsg: Message = { role: 'bot', text: '', time: getTime(), typing: true }
  messages.value.push(typingMsg)
  statusType.value = 'typing'
  statusText.value = 'Печатает...'

  const thinkingTimer = setTimeout(() => {
    statusType.value = 'thinking'
    statusText.value = 'Обращается к ИИ...'
  }, 1500)

  try {
    const res = await fetch('/api/v1/chat', {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message:    msgText,
        session_id: sessionId,
        user_id:    auth.user?.id ?? '',
        history:    history.value.slice(-8), // последние 8
      }),
    })
    const data = await res.json()

    clearTimeout(thinkingTimer)

    // Убираем typing bubble
    const idx = [...messages.value].reverse().findIndex(m => m.typing)
    const realIdx = idx !== -1 ? messages.value.length - 1 - idx : -1
    if (realIdx !== -1) messages.value.splice(realIdx, 1)

    const botText = data.message ?? 'Что-то пошло не так. Попробуйте ещё раз.'

    const botMsg: Message = {
      role:   'bot',
      text:   botText,
      time:   getTime(),
      source: data.source,
    }
    messages.value.push(botMsg)

    // Обновляем сессию бронирования
    if (data.session) {
      currentSession.value = data.session
    }

    // Quick replies если есть
    if (data.quick_replies?.length) {
      botMsg.quickReplies = data.quick_replies
    }
    // Если нужен логин — добавляем кнопку только если её ещё нет
    if (botText.includes('войти') || botText.includes('аккаунт')) {
      if (!auth.isLoggedIn && !botMsg.quickReplies?.includes('🔑 Войти')) {
        botMsg.quickReplies = [...(botMsg.quickReplies ?? []), '🔑 Войти']
      }
    }

    history.value.push({ role: 'bot', message: botText })
    saveToStorage()

  } catch {
    clearTimeout(thinkingTimer)
    const idx2 = [...messages.value].reverse().findIndex(m => m.typing)
    const realIdx2 = idx2 !== -1 ? messages.value.length - 1 - idx2 : -1
    if (realIdx2 !== -1) messages.value.splice(realIdx2, 1)
    messages.value.push({
      role: 'bot',
      text: 'Что-то пошло не так. Попробуйте ещё раз или свяжитесь с нами напрямую.',
      time: getTime(),
    })
  } finally {
    loading.value    = false
    statusType.value = 'online'
    statusText.value = 'Онлайн помощник'
    await nextTick(scrollToBottom)
  }
}

// Streaming — показываем текст по символам


function scrollToBottom() {
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight
  }
}

function formatText(text: string): string {
  if (!text) return ''
  try {
    return text
      .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
      .replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" style="color:var(--accent)">$1</a>')
      .replace(/\n/g, '<br/>')
  } catch {
    return text.replace(/\n/g, '<br/>')
  }
}
</script>

<template>
  <div class="chat-fab" @click="open">
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
      <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z"
            stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
    </svg>
    <span class="chat-fab__label">Чат</span>
  </div>

  <Transition name="chat-slide">
    <div v-if="isOpen" class="chat-window">

      <div class="chat-header">
        <div class="chat-header__left">
          <div class="chat-header__avatar">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14H9V8h2v8zm4 0h-2V8h2v8z"
                fill="rgba(255,255,255,0.9)"/>
            </svg>
          </div>
          <div>
            <p class="chat-header__name">Beauty Dana · ИИ</p>
            <p class="chat-header__status">
              <span
                class="chat-header__dot"
                :class="{
                  'chat-header__dot--online':   statusType === 'online',
                  'chat-header__dot--typing':   statusType === 'typing',
                  'chat-header__dot--thinking': statusType === 'thinking',
                }"
              />
              {{ statusText }}
            </p>
          </div>
        </div>
        <div class="chat-header__actions">
          <button
            class="chat-header__btn"
            @click="clearChat"
            title="Очистить историю"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6"/>
            </svg>
          </button>
          <button class="chat-header__btn" @click="isOpen = false"><X :size="16" /></button>
        </div>
      </div>

      <div class="chat-messages" ref="messagesEl">
        <div
          v-for="(msg, i) in messages"
          :key="i"
          class="chat-msg"
          :class="msg.role === 'user' ? 'chat-msg--user' : 'chat-msg--bot'"
        >
          <div v-if="msg.typing" class="chat-bubble chat-bubble--bot">
            <div class="chat-typing">
              <span /><span /><span />
            </div>
          </div>

          <template v-else>
            <div
              class="chat-bubble"
              :class="msg.role === 'user' ? 'chat-bubble--user' : 'chat-bubble--bot'"
              v-html="formatText(msg.text)"
            />
            <div class="chat-msg__meta">
              <span class="chat-msg__time">{{ msg.time }}</span>
              <span v-if="msg.source === 'gemini'" class="chat-msg__ai-badge">ИИ</span>
            </div>

            <div
              v-if="msg.role === 'bot' && i === messages.length - 1 && msg.quickReplies?.length"
              class="chat-replies"
            >
              <button
                v-for="reply in msg.quickReplies"
                :key="reply"
                class="chat-reply"
                @click="send(reply)"
              >{{ reply }}</button>
            </div>
          </template>
        </div>
      </div>

      <div class="chat-input-row">
        <input
          v-model="input"
          type="text"
          class="chat-input"
          placeholder="Напишите вопрос..."
          @keyup.enter="send()"
          :disabled="loading"
        />
        <button
          class="chat-send"
          :disabled="!input.trim() || loading"
          @click="send()"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <path d="M22 2L11 13M22 2L15 22l-4-9-9-4 20-7z"
                  stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </button>
      </div>

    </div>
  </Transition>
</template>

<style scoped>
.chat-fab {
  position: fixed; bottom: 24px; right: 24px; z-index: 200;
  display: flex; align-items: center; gap: 8px;
  background: var(--accent); color: #fff;
  padding: 12px 20px; border-radius: 100px;
  cursor: pointer; box-shadow: 0 4px 20px rgba(232,68,42,0.4);
  transition: all 0.2s; user-select: none;
}
.chat-fab:hover { transform: translateY(-2px); box-shadow: 0 6px 24px rgba(232,68,42,0.5); }
.chat-fab__label { font-size: 14px; font-weight: 600; }

.chat-window {
  position: fixed; bottom: 84px; right: 24px; z-index: 200;
  width: 380px; height: 540px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  display: flex; flex-direction: column;
  overflow: hidden;
  box-shadow: 0 8px 40px rgba(0,0,0,0.12);
}

.chat-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14px 16px;
  background: var(--accent);
  flex-shrink: 0;
}
.chat-header__left   { display: flex; align-items: center; gap: 10px; }
.chat-header__avatar {
  width: 36px; height: 36px; border-radius: 50%;
  background: rgba(255,255,255,0.15);
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
}
.chat-header__name   { font-size: 13px; font-weight: 600; color: #fff; }
.chat-header__status {
  display: flex; align-items: center; gap: 5px;
  font-size: 11px; color: rgba(255,255,255,0.75);
  margin-top: 1px;
}
.chat-header__dot {
  width: 6px; height: 6px; border-radius: 50%; flex-shrink: 0;
}
.chat-header__dot--online   { background: #4ade80; }
.chat-header__dot--typing   { background: #fbbf24; animation: blink 1s infinite; }
.chat-header__dot--thinking { background: #60a5fa; animation: blink 0.6s infinite; }

@keyframes blink {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.3; }
}

.chat-header__actions { display: flex; align-items: center; gap: 4px; }
.chat-header__btn {
  width: 28px; height: 28px;
  display: flex; align-items: center; justify-content: center;
  background: rgba(255,255,255,0.1); border: none;
  border-radius: 6px; color: rgba(255,255,255,0.8);
  cursor: pointer; transition: background 0.15s;
  font-size: 13px;
}
.chat-header__btn:hover { background: rgba(255,255,255,0.2); }

.chat-messages {
  flex: 1; overflow-y: auto; padding: 16px 14px;
  display: flex; flex-direction: column; gap: 10px;
  background: var(--bg);
  scrollbar-width: thin;
  scrollbar-color: var(--border) transparent;
}

.chat-msg { display: flex; flex-direction: column; max-width: 88%; gap: 3px; }
.chat-msg--user { align-self: flex-end; align-items: flex-end; }
.chat-msg--bot  { align-self: flex-start; align-items: flex-start; }

.chat-bubble {
  padding: 10px 14px; border-radius: 16px;
  font-size: 14px; line-height: 1.55; word-break: break-word;
}
.chat-bubble--user {
  background: var(--accent); color: #fff;
  border-bottom-right-radius: 4px;
}
.chat-bubble--bot {
  background: var(--bg-card); color: var(--text);
  border: 1px solid var(--border);
  border-bottom-left-radius: 4px;
}

.chat-msg__meta {
  display: flex; align-items: center; gap: 6px; padding: 0 4px;
}
.chat-msg__time { font-size: 10px; color: var(--text-3); }
.chat-msg__ai-badge {
  font-size: 9px; font-weight: 600; letter-spacing: 0.06em;
  padding: 1px 5px; background: var(--accent-bg);
  color: var(--accent); border-radius: 3px;
}

.chat-typing {
  display: flex; align-items: center; gap: 4px; padding: 4px 0;
}
.chat-typing span {
  width: 7px; height: 7px; border-radius: 50%;
  background: var(--text-3);
  animation: bounce 1.2s infinite;
}
.chat-typing span:nth-child(2) { animation-delay: 0.2s; }
.chat-typing span:nth-child(3) { animation-delay: 0.4s; }
@keyframes bounce {
  0%, 60%, 100% { transform: translateY(0); }
  30%           { transform: translateY(-6px); }
}

.chat-replies { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 6px; }
.chat-reply {
  font-size: 12px; padding: 5px 12px;
  background: var(--bg-card);
  border: 1.5px solid var(--accent);
  color: var(--accent); border-radius: 100px;
  cursor: pointer; transition: all 0.15s;
  white-space: nowrap;
}
.chat-reply:hover { background: var(--accent); color: #fff; }

.chat-input-row {
  display: flex; gap: 8px; padding: 12px 14px;
  border-top: 1px solid var(--border);
  background: var(--bg-card); flex-shrink: 0;
}
.chat-input {
  flex: 1; padding: 9px 14px;
  background: var(--bg); color: var(--text);
  border: 1.5px solid var(--border); border-radius: 100px;
  font-size: 14px; outline: none; transition: border-color 0.15s;
  font-family: var(--font);
}
.chat-input:focus  { border-color: var(--accent); }
.chat-send {
  width: 38px; height: 38px; border-radius: 50%;
  background: var(--accent); color: #fff; border: none;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; flex-shrink: 0; transition: all 0.15s;
}
.chat-send:hover:not(:disabled) { transform: scale(1.08); }
.chat-send:disabled { opacity: 0.4; cursor: not-allowed; }

.chat-slide-enter-active { transition: all 0.25s cubic-bezier(0.4,0,0.2,1); }
.chat-slide-leave-active { transition: all 0.2s  cubic-bezier(0.4,0,0.2,1); }
.chat-slide-enter-from   { opacity: 0; transform: translateY(20px) scale(0.96); }
.chat-slide-leave-to     { opacity: 0; transform: translateY(16px) scale(0.96); }

@media (max-width: 440px) {
  .chat-window { width: calc(100vw - 24px); right: 12px; bottom: 80px; height: 500px; }
  .chat-fab    { right: 16px; bottom: 16px; }
}
</style>