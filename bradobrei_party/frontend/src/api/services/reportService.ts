import { apiRequest } from '../client'
import type {
  EmployeeReportResponseDto,
  MasterActivityReportResponseDto,
  ReportPeriodQuery,
  SalonActivityReportResponseDto,
} from '../../types/dto/reports'

export const reportService = {
  getEmployees() {
    return apiRequest<EmployeeReportResponseDto>('/reports/employees')
  },
  getSalonActivity(query: ReportPeriodQuery) {
    return apiRequest<SalonActivityReportResponseDto>('/reports/salon-activity', {
      query,
    })
  },
  getMasterActivity(query: ReportPeriodQuery) {
    return apiRequest<MasterActivityReportResponseDto>('/reports/master-activity', {
      query,
    })
  },
}
