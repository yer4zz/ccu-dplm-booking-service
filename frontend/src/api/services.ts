import { api } from './index'
import type { Service } from '@/types'

export const servicesApi = {
  list: () => api.get<Service[]>('/services').then(r => r.data),
}