import { api } from './index'

export type AnalyticsRange = 'week' | 'month' | 'quarter' | 'year'

export const analyticsApi = {
  get: (range: AnalyticsRange) =>
    api.get('/admin/analytics', { params: { range } }).then(r => r.data),
}