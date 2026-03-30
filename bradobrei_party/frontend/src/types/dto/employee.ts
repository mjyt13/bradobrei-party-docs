import type { UserDto } from './auth'
import type { SalonBriefDto, UserRole } from './common'
import type { ServiceDto } from './entities'
import type { EmployeeProfileDto } from './auth'

export interface HireEmployeeRequestDto {
  username: string
  password: string
  full_name: string
  phone: string
  email: string
  role: UserRole
  specialization: string
  expected_salary: number
  work_schedule: string
  salon_id: number
}

export type HireEmployeeResponseDto = EmployeeProfileDto

export interface UpdateEmployeeRequestDto {
  username: string
  full_name: string
  phone: string
  email: string
  role: UserRole
  specialization: string
  expected_salary: number
  work_schedule: string
  salon_ids: number[]
}

export interface EmployeeManagementDto {
  id: number
  user_id: number
  specialization: string
  expected_salary: number
  work_schedule?: string
  user?: UserDto
  salons?: SalonBriefDto[]
  services?: ServiceDto[]
}
