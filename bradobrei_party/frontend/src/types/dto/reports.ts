import type { UserDto } from './auth'

export interface PeriodDto {
  from: string
  to: string
}

export type ReportPeriodQuery = Record<'from' | 'to', string>

export interface ReportEnvelopeDto<T> {
  report: string
  data: T
  period?: PeriodDto
}

export interface SalonActivityRowDto {
  salon_id: number
  address: string
  client_count: number
  service_count: number
  total_revenue: number
}

export interface MasterActivityRowDto {
  master_id: number
  full_name: string
  service_count: number
  revenue: number
  material_cost: number
}

export type EmployeeReportResponseDto = ReportEnvelopeDto<UserDto[]>
export type SalonActivityReportResponseDto = ReportEnvelopeDto<SalonActivityRowDto[]>
export type MasterActivityReportResponseDto = ReportEnvelopeDto<MasterActivityRowDto[]>
