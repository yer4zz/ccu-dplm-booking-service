import { api } from './index'

export const adminApi = {
    getServices: () =>
        api.get('/admin/services/all').then(r => r.data),

    createService: (data: {
        name: string; category: string; description?: string
        duration_min: number; price: number
    }) => api.post('/admin/services', data).then(r => r.data),

    updateService: (id: string, data: Partial<{
        name: string; category: string; description: string
        duration_min: number; price: number; is_active: boolean
    }>) => api.patch(`/admin/services/${id}`, data).then(r => r.data),

    deleteService: (id: string) =>
        api.delete(`/admin/services/${id}`).then(r => r.data),

    updateMaster: (id: string, data: Partial<{
        bio: string; experience_years: number
        instagram: string; is_active: boolean; full_name: string
    }>) => api.patch(`/admin/masters/${id}`, data).then(r => r.data),

    deleteMaster: (id: string) =>
        api.delete(`/admin/masters/${id}`).then(r => r.data),
}

export const masterApi = {
    getMyServices: () =>
        api.get('/master/my-services').then(r => r.data),
  
    addService: (service_id: string, custom_price?: number) =>
        api.post('/master/my-services', { service_id, custom_price }).then(r => r.data),
  
    removeService: (service_id: string) =>
        api.delete(`/master/my-services/${service_id}`).then(r => r.data),
  
    updateServicePrice: (service_id: string, custom_price: number | null) =>
        api.patch(`/master/my-services/${service_id}`, { custom_price }).then(r => r.data),
}