import axios from 'axios'

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',    //http://localhost:8080/api/v1
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers = config.headers ?? {}
    config.headers['Authorization'] = `Bearer ${token}`
  }
  console.log('[API]', config.method?.toUpperCase(), config.url, 
              token ? 'token present' : 'no token')
  return config
})

api.interceptors.response.use(
  (res) => res,
  (err) => {
    console.error('[API] error:', err.response?.status, err.response?.data)
    if (err.response?.status === 401) {
      localStorage.removeItem('access_token')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)