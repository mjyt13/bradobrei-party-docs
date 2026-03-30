import { apiRequest } from '../client'
import type { BookingDto, CreateBookingRequestDto } from '../../types/dto/entities'

export const bookingService = {
  getAll() {
    return apiRequest<BookingDto[]>('/bookings')
  },
  getMy() {
    return apiRequest<BookingDto[]>('/bookings/my')
  },
  getMaster() {
    return apiRequest<BookingDto[]>('/bookings/master')
  },
  getById(id: number) {
    return apiRequest<BookingDto>(`/bookings/${id}`)
  },
  create(payload: CreateBookingRequestDto) {
    return apiRequest<BookingDto>('/bookings', {
      method: 'POST',
      body: payload,
    })
  },
  confirm(id: number) {
    return apiRequest<BookingDto>(`/bookings/${id}/confirm`, {
      method: 'POST',
    })
  },
  cancel(id: number) {
    return apiRequest<{ message: string }>(`/bookings/${id}/cancel`, {
      method: 'POST',
    })
  },
}
