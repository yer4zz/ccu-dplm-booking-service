import { ref, watch, onMounted } from 'vue'

const theme = ref<'light' | 'dark'>('light')

export function useTheme() {
  onMounted(() => {
    const saved = localStorage.getItem('theme') as 'light' | 'dark' | null
    const preferred = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    theme.value = saved ?? preferred
    applyTheme(theme.value)
  })

  watch(theme, (val) => {
    applyTheme(val)
    localStorage.setItem('theme', val)
  })

  function applyTheme(val: string) {
    document.documentElement.setAttribute('data-theme', val)
  }

  function toggle() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  return { theme, toggle }
}