import { apiRequest } from '../client'
import type { MaterialDto, UpsertMaterialRequestDto } from '../../types/dto/entities'

export const materialService = {
  getAll() {
    return apiRequest<MaterialDto[]>('/materials')
  },
  getById(id: number) {
    return apiRequest<MaterialDto>(`/materials/${id}`)
  },
  create(payload: UpsertMaterialRequestDto) {
    return apiRequest<MaterialDto>('/materials', {
      method: 'POST',
      body: payload,
    })
  },
  update(id: number, payload: UpsertMaterialRequestDto) {
    return apiRequest<MaterialDto>(`/materials/${id}`, {
      method: 'PUT',
      body: payload,
    })
  },
  remove(id: number) {
    return apiRequest<{ message: string }>(`/materials/${id}`, {
      method: 'DELETE',
    })
  },
}
