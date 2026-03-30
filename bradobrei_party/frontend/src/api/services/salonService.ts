import { apiRequest } from '../client'
import type {
  EmployeeProfileSummaryDto,
  SalonDto,
  UpsertSalonRequestDto,
} from '../../types/dto/entities'

export const salonService = {
  getAll() {
    return apiRequest<SalonDto[]>('/salons')
  },
  getById(id: number) {
    return apiRequest<SalonDto>(`/salons/${id}`)
  },
  getMasters(id: number) {
    return apiRequest<EmployeeProfileSummaryDto[]>(`/salons/${id}/masters`)
  },
  create(payload: UpsertSalonRequestDto) {
    return apiRequest<SalonDto>('/salons', {
      method: 'POST',
      body: payload,
    })
  },
  update(id: number, payload: UpsertSalonRequestDto) {
    return apiRequest<SalonDto>(`/salons/${id}`, {
      method: 'PUT',
      body: payload,
    })
  },
  remove(id: number) {
    return apiRequest<{ message: string }>(`/salons/${id}`, {
      method: 'DELETE',
    })
  },
}
