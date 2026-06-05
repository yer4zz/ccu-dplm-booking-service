import { api } from './index'

export const loyaltyApi = {
  getAccount: () =>
    api.get('/loyalty').then(r => r.data),

  getTransactions: () =>
    api.get('/loyalty/transactions').then(r => r.data),

  calcPrice: (payload: {
    master_id: string
    service_ids: string[]
    points_used?: number
  }) => api.post('/bookings/calc-price', payload).then(r => r.data),

  addToWaitlist: (master_id: string, service_id: string) =>
    api.post('/waitlist', { master_id, service_id }).then(r => r.data),

  getWaitlist: () =>
    api.get('/waitlist').then(r => r.data),

  removeFromWaitlist: (id: string) =>
    api.delete(`/waitlist/${id}`).then(r => r.data),

  getStreak: () =>
    api.get('/streak').then(r => r.data),

  createSOS: (payload: {
    master_id:       string
    service_id:      string
    preferred_start: string
    preferred_end:   string
    client_note?:    string
  }) => api.post('/sos', payload).then(r => r.data),

  getClientSOS: () =>
    api.get('/sos').then(r => r.data),

  getMasterSOS: () =>
    api.get('/master/sos').then(r => r.data),

  respondSOS: (id: string, accept: boolean, note?: string, startsAt?: string) =>
    api.patch(`/master/sos/${id}`, {
      accept,
      note:      note     ?? '',
      starts_at: startsAt ?? null,
    }).then(r => r.data),
}