import { api } from './index'
import type { Booking } from '@/types'

export const bookingsApi = {
  create: (payload: {
    master_id:   string
    service_ids: string[]
    starts_at:   string
    notes?:      string
    points_used?: number
    vibe_mode?:   string
    vibe_note?:   string
  }) => api.post<Booking>('/bookings', payload).then(r => r.data),

  my: () =>
    api.get<Booking[]>('/bookings/my').then(r => r.data),

  cancel: (id: string) =>
    api.patch(`/bookings/${id}/cancel`).then(r => r.data),

  masterList: () =>
    api.get<Booking[]>('/master/bookings').then(r => r.data),

  updateStatus: (id: string, status: Booking['status']) =>
    api.patch(`/master/bookings/${id}`, { status }).then(r => r.data),

  getAutoReschedules: () =>
    api.get('/auto-reschedules').then(r => r.data),
  
  declineAutoReschedule: (id: string) =>
    api.delete(`/auto-reschedules/${id}`).then(r => r.data),


  getReschedules: () =>
    api.get('/reschedules').then(r => r.data),
  
  declineReschedule: (id: string) =>
    api.delete(`/reschedules/${id}`).then(r => r.data),
  
  createReview: (bookingId: string, rating: number, comment: string) =>
    api.post(`/bookings/${bookingId}/review`, { rating, comment }).then(r => r.data),
  
  masterReschedule: (bookingId: string, payload: {
    new_master_id?: string
    new_starts_at:  string
    reason?:        string
  }) => api.post(`/master/bookings/${bookingId}/reschedule`, payload).then(r => r.data),
  
  getMasterStats: () =>
    api.get('/master/stats').then(r => r.data),
  
  getMasterReviews: () =>
    api.get('/master/reviews').then(r => r.data),
}

