import { apiRequest } from '../client'
import type { ServiceDto, UpsertServiceRequestDto } from '../../types/dto/entities'

export const serviceService = {
  getAll() {
    return apiRequest<ServiceDto[]>('/services')
  },
  getById(id: number) {
    return apiRequest<ServiceDto>(`/services/${id}`)
  },
  create(payload: UpsertServiceRequestDto) {
    return apiRequest<ServiceDto>('/services', {
      method: 'POST',
      body: payload,
    })
  },
  update(id: number, payload: UpsertServiceRequestDto) {
    return apiRequest<ServiceDto>(`/services/${id}`, {
      method: 'PUT',
      body: payload,
    })
  },
  remove(id: number) {
    return apiRequest<{ message: string }>(`/services/${id}`, {
      method: 'DELETE',
    })
  },
}
