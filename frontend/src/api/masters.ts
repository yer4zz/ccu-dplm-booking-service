import { api } from './index'
import type { Master } from '@/types'

export const mastersApi = {
  list: () =>
    api.get<Master[]>('/masters').then(r => r.data),

  getSlots: (masterID: string, serviceIDs: string | string[], date: string) => {
    const ids = Array.isArray(serviceIDs) ? serviceIDs : [serviceIDs]
    const params = new URLSearchParams()
    params.set('master_id', masterID)
    params.set('date', date)
    ids.forEach(id => params.append('service_ids', id))
    return api.get(`/slots?${params.toString()}`).then(r => r.data?.slots ?? [])
  },
}