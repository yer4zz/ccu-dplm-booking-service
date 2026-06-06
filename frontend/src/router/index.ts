import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',       name: 'home',    component: () => import('@/views/HomeView.vue') },
    { path: '/book',   name: 'book',    component: () => import('@/views/BookingFlow.vue') },
    { path: '/login',  name: 'login',   component: () => import('@/views/LoginView.vue') },
    { path: '/my',     name: 'my',      component: () => import('@/views/MyBookings.vue'),
      meta: { requiresAuth: true } },
    { path: '/master', name: 'master',  component: () => import('@/views/MasterDashboard.vue'),
      meta: { requiresAuth: true, role: 'master' } },
    { path: '/admin',  name: 'admin',   component: () => import('@/views/AdminPanel.vue'),
      meta: { requiresAuth: true, role: 'admin' } },
    { path: '/profile', name: 'profile', component: () => import('@/views/ProfileView.vue'),
      meta: { requiresAuth: true } },
    { path: '/gallery', name: 'gallery', component: () => import('@/views/GalleryView.vue') },
    { path: '/masters', component: () => import('@/views/MastersView.vue') },
    { path: '/about',   component: () => import('@/views/AboutView.vue') },
  ],
})

let sessionRestored = false

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  if (!sessionRestored) {
    sessionRestored = true
    await auth.restoreSession()
  }

  if (to.meta.requiresAuth && !auth.isLoggedIn) return '/login'
  if (to.meta.role === 'master' && !auth.isMaster) return '/'
  if (to.meta.role === 'admin'  && !auth.isAdmin)  return '/'
  return true
})

export default router