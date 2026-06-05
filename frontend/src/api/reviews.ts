import { api } from './index'

export const reviewsApi = {
  deleteMaster: (id: string) =>
    api.delete(`/master/reviews/${id}`).then(r => r.data),

  deleteAdmin: (id: string) =>
    api.delete(`/admin/reviews/${id}`).then(r => r.data),
}