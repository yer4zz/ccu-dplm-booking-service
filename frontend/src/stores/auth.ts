import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/api'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user  = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('access_token'))

  const isLoggedIn = computed(() => !!token.value)
  const isMaster   = computed(() => user.value?.role === 'master')
  const isAdmin    = computed(() => user.value?.role === 'admin')

  async function login(email: string, password: string) {
    const res = await api.post('/auth/login', { email, password })

    const access_token = res.data.access_token
    if (!access_token) throw new Error('no token in response')

    token.value = access_token
    localStorage.setItem('access_token', access_token)

    const meRes = await api.get('/me', {
      headers: { Authorization: `Bearer ${access_token}` }
    })

    const userData = meRes.data
    user.value = {
      id:        userData.id,
      email:     userData.email,
      full_name: userData.user_metadata?.full_name ?? '',
      phone:     userData.user_metadata?.phone     ?? '',
      role:      userData.app_metadata?.role ?? 'client',
    }
  }

  async function register(email: string, password: string, fullName: string) {
    await api.post('/auth/register', {
      email,
      password,
      full_name: fullName,
    })
    await login(email, password)
  }

  function logout() {
    token.value = null
    user.value  = null
    localStorage.removeItem('access_token')
  }

  let meCache: any = null
  let meCacheTime  = 0

  async function restoreSession() {
    const savedToken = localStorage.getItem('access_token')
    if (!savedToken) return
    
    if (meCache && Date.now() - meCacheTime < 5 * 60 * 1000) {
      token.value = savedToken
      user.value  = meCache
      return
    }

    try {
      const meRes = await api.get('/me')
      token.value = savedToken

      const userData = {
        id:        meRes.data.id,
        email:     meRes.data.email,
        full_name: meRes.data.profile_full_name || meRes.data.user_metadata?.full_name || '',
        phone:     meRes.data.profile_phone     || '',
        role:      meRes.data.app_metadata?.role ?? 'client',
      }
      user.value   = userData
      meCache      = userData
      meCacheTime  = Date.now()
    } catch {
      logout()
    }
  }

  return { user, token, isLoggedIn, isMaster, isAdmin, login, register, logout, restoreSession }
})