import { api } from './index'

export const galleryApi = {
  list: (masterID?: string, serviceID?: string) =>
    api.get('/gallery', { params: { master_id: masterID, service_id: serviceID } })
      .then(r => r.data),

    upload: (formData: FormData) =>
    api.post('/master/gallery', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }).then(r => r.data),

  delete: (id: string) =>
    api.delete(`/gallery/${id}`).then(r => r.data),
}