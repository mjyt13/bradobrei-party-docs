export interface ErrorResponseDto {
  error: string
  code: number
  message?: string
}

export type UserRole =
  | 'CLIENT'
  | 'BASIC_MASTER'
  | 'ADVANCED_MASTER'
  | 'HR_SPECIALIST'
  | 'ACCOUNTANT'
  | 'NETWORK_MANAGER'
  | 'ADMINISTRATOR'

export interface SalonBriefDto {
  id: number
  name: string
  address: string
}
