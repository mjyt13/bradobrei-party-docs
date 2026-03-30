import { useState } from 'react'
import { DataTable } from '../components/DataTable'
import { reportService } from '../api/services/reportService'
import { masterActivityColumns } from '../shared/reportTables'
import type { MasterActivityRowDto } from '../types/dto/reports'

function getDefaultPeriod() {
  const now = new Date()
  const from = new Date(now)
  from.setDate(now.getDate() - 30)

  return {
    from: from.toISOString().slice(0, 10),
    to: now.toISOString().slice(0, 10),
  }
}

export function MasterActivityReportPage() {
  const [filters, setFilters] = useState(getDefaultPeriod())
  const [rows, setRows] = useState<MasterActivityRowDto[]>([])
  const [error, setError] = useState('')
  const [periodLabel, setPeriodLabel] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await reportService.getMasterActivity(filters)
      setRows(response.data)
      if (response.period) {
        setPeriodLabel(`${response.period.from} — ${response.period.to}`)
      }
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось загрузить отчёт.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="page-section">
      <div className="page-header">
        <div>
          <p className="eyebrow">Отчёт 2.2.4</p>
          <h2>Активность мастеров</h2>
          <p className="section-description">
            Выполненные услуги, выручка и материальные затраты по мастерам за период.
          </p>
        </div>
      </div>

      <form className="filter-card" onSubmit={handleSubmit}>
        <label className="field">
          <span>С</span>
          <input
            type="date"
            value={filters.from}
            onChange={(event) => setFilters((current) => ({ ...current, from: event.target.value }))}
          />
        </label>
        <label className="field">
          <span>По</span>
          <input
            type="date"
            value={filters.to}
            onChange={(event) => setFilters((current) => ({ ...current, to: event.target.value }))}
          />
        </label>
        <button type="submit" className="primary-button">
          {loading ? 'Собираем...' : 'Показать отчёт'}
        </button>
      </form>

      {periodLabel ? <p className="period-label">Период: {periodLabel}</p> : null}
      {error ? <div className="alert alert-error">{error}</div> : null}

      <DataTable
        caption="Активность мастеров"
        columns={masterActivityColumns}
        rows={rows}
        emptyText="Нет данных за выбранный период."
      />
    </section>
  )
}
