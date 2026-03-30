import type { TableColumn } from '../components/DataTable'
import type { UserDto } from '../types/dto/auth'
import type { MasterActivityRowDto, SalonActivityRowDto } from '../types/dto/reports'
import { formatCurrency, formatRole } from './formatters'

export const employeeReportColumns: Array<TableColumn<UserDto>> = [
  { key: 'name', header: 'Сотрудник', render: (row) => row.full_name },
  { key: 'role', header: 'Роль', render: (row) => formatRole(row.role) },
  {
    key: 'specialization',
    header: 'Специализация',
    render: (row) => row.employee_profile?.specialization || '—',
  },
  {
    key: 'salary',
    header: 'Ожидаемый оклад',
    render: (row) =>
      row.employee_profile?.expected_salary
        ? formatCurrency(row.employee_profile.expected_salary)
        : '—',
  },
  {
    key: 'salons',
    header: 'Филиалы',
    render: (row) =>
      row.employee_profile?.salons?.length
        ? row.employee_profile.salons.map((salon) => salon.address).join(', ')
        : 'Не назначены',
  },
]

export const salonActivityColumns: Array<TableColumn<SalonActivityRowDto>> = [
  { key: 'address', header: 'Филиал', render: (row) => row.address },
  { key: 'clients', header: 'Клиенты', render: (row) => row.client_count },
  { key: 'services', header: 'Услуги', render: (row) => row.service_count },
  {
    key: 'revenue',
    header: 'Выручка',
    render: (row) => formatCurrency(row.total_revenue),
  },
]

export const masterActivityColumns: Array<TableColumn<MasterActivityRowDto>> = [
  { key: 'master', header: 'Мастер', render: (row) => row.full_name },
  { key: 'count', header: 'Услуг выполнено', render: (row) => row.service_count },
  { key: 'revenue', header: 'Выручка', render: (row) => formatCurrency(row.revenue) },
  {
    key: 'materials',
    header: 'Материалы',
    render: (row) => formatCurrency(row.material_cost),
  },
]
