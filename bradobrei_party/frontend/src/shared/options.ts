import type { UserRole } from '../types/dto/common'
import type { BookingStatus, PaymentStatus, SalonStatus } from '../types/dto/entities'

export const employeeRoleOptions: Array<{ label: string; value: UserRole }> = [
  { label: 'Базовый мастер', value: 'BASIC_MASTER' },
  { label: 'Старший мастер', value: 'ADVANCED_MASTER' },
  { label: 'HR-специалист', value: 'HR_SPECIALIST' },
  { label: 'Бухгалтер', value: 'ACCOUNTANT' },
  { label: 'Менеджер сети', value: 'NETWORK_MANAGER' },
  { label: 'Администратор', value: 'ADMINISTRATOR' },
]

export const salonStatusOptions: Array<{ label: string; value: SalonStatus }> = [
  { label: 'Открыт', value: 'OPEN' },
  { label: 'Закрыт', value: 'CLOSED' },
]

export const paymentStatusOptions: Array<{ label: string; value: PaymentStatus }> = [
  { label: 'Ожидает оплаты', value: 'PENDING' },
  { label: 'Оплачен', value: 'SUCCESS' },
  { label: 'Ошибка оплаты', value: 'FAILED' },
  { label: 'Возврат', value: 'REFUNDED' },
]

export const bookingStatusOptions: Array<{ label: string; value: BookingStatus }> = [
  { label: 'Ожидает', value: 'PENDING' },
  { label: 'Подтверждено', value: 'CONFIRMED' },
  { label: 'В работе', value: 'IN_PROGRESS' },
  { label: 'Завершено', value: 'COMPLETED' },
  { label: 'Отменено', value: 'CANCELLED' },
]

export const materialUnitOptions = ['мл', 'шт', 'гр', 'л', 'уп.']
