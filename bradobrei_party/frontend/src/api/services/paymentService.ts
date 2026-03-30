import { apiRequest } from '../client'
import type { CreatePaymentRequestDto, PaymentDto } from '../../types/dto/entities'

export const paymentService = {
  getAll() {
    return apiRequest<PaymentDto[]>('/payments')
  },
  getById(id: number) {
    return apiRequest<PaymentDto>(`/payments/${id}`)
  },
  create(payload: CreatePaymentRequestDto) {
    return apiRequest<PaymentDto>('/payments', {
      method: 'POST',
      body: payload,
    })
  },
}
