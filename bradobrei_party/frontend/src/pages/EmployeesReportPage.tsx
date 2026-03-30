import { useEffect, useState } from 'react'
import { DataTable } from '../components/DataTable'
import { reportService } from '../api/services/reportService'
import { employeeReportColumns } from '../shared/reportTables'
import type { UserDto } from '../types/dto/auth'

export function EmployeesReportPage() {
  const [rows, setRows] = useState<UserDto[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false

    reportService
      .getEmployees()
      .then((response) => {
        if (!cancelled) {
          setRows(response.data)
        }
      })
      .catch((requestError: Error) => {
        if (!cancelled) {
          setError(requestError.message)
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false)
        }
      })

    return () => {
      cancelled = true
    }
  }, [])

  return (
    <section className="page-section">
      <div className="page-header">
        <div>
          <p className="eyebrow">Отчёт 2.2.1</p>
          <h2>Реестр персонала</h2>
          <p className="section-description">
            Сводная таблица по сотрудникам, ролям, специализациям и филиалам.
          </p>
        </div>
      </div>

      {error ? <div className="alert alert-error">{error}</div> : null}

      <DataTable
        caption={loading ? 'Загрузка реестра персонала...' : 'Реестр персонала'}
        columns={employeeReportColumns}
        rows={rows}
        emptyText="Сотрудники пока не найдены."
      />
    </section>
  )
}
