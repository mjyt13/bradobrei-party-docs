import type { SalonBriefDto, UserRole } from './common'

export interface LoginRequestDto {
  username: string
  password: string
}

export interface LoginResponseDto {
  token: string
}

export interface RegisterRequestDto {
  username: string
  password: string
  full_name: string
  phone: string
  email: string
  role?: UserRole
}

export interface EmployeeProfileDto {
  id: number
  user_id: number
  specialization: string
  expected_salary: number
  work_schedule?: string
  salons?: SalonBriefDto[]
}

export interface UserDto {
  id: number
  username: string
  full_name: string
  phone: string
  email?: string
  role: UserRole
  created_at: string
  updated_at: string
  employee_profile?: EmployeeProfileDto
}

export interface UserResponseDto {
  user: UserDto
}
