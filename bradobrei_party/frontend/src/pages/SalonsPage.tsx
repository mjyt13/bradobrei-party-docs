import { useEffect, useState } from 'react'
import { DataTable, type TableColumn } from '../components/DataTable'
import { salonService } from '../api/services/salonService'
import type { EmployeeProfileSummaryDto, SalonDto, UpsertSalonRequestDto } from '../types/dto/entities'
import { formatCurrency, formatJsonPreview, formatRole, formatSalonStatus } from '../shared/formatters'
import { salonStatusOptions } from '../shared/options'

const initialForm: UpsertSalonRequestDto = {
  name: '',
  address: '',
  location: '',
  working_hours: '{"mon":"10:00-20:00","tue":"10:00-20:00"}',
  status: 'OPEN',
  max_staff: 8,
  base_hourly_rate: 1400,
}

const salonColumns = (
  onEdit: (salon: SalonDto) => void,
  onDelete: (salon: SalonDto) => void,
  onShowMasters: (salon: SalonDto) => void,
): Array<TableColumn<SalonDto>> => [
  { key: 'name', header: 'Салон', render: (row) => row.name },
  { key: 'address', header: 'Адрес', render: (row) => row.address },
  {
    key: 'status',
    header: 'Статус',
    render: (row) => <span className="status-pill">{formatSalonStatus(row.status)}</span>,
  },
  { key: 'staff', header: 'Штат', render: (row) => `${row.max_staff} чел.` },
  { key: 'rate', header: 'Базовая ставка', render: (row) => formatCurrency(row.base_hourly_rate) },
  { key: 'hours', header: 'Часы', render: (row) => formatJsonPreview(row.working_hours) },
  {
    key: 'actions',
    header: 'Действия',
    render: (row) => (
      <div className="table-actions">
        <button type="button" className="ghost-button button-small" onClick={() => onEdit(row)}>
          Изменить
        </button>
        <button type="button" className="ghost-button button-small" onClick={() => onShowMasters(row)}>
          Мастера
        </button>
        <button type="button" className="danger-button button-small" onClick={() => onDelete(row)}>
          Удалить
        </button>
      </div>
    ),
  },
]

const masterColumns: Array<TableColumn<EmployeeProfileSummaryDto>> = [
  { key: 'name', header: 'Сотрудник', render: (row) => row.user?.full_name || `Профиль #${row.id}` },
  { key: 'role', header: 'Роль', render: (row) => formatRole(row.user?.role || '—') },
  { key: 'spec', header: 'Специализация', render: (row) => row.specialization || '—' },
  { key: 'services', header: 'Услуги', render: (row) => row.services?.map((service) => service.name).join(', ') || '—' },
]

export function SalonsPage() {
  const [form, setForm] = useState(initialForm)
  const [salons, setSalons] = useState<SalonDto[]>([])
  const [masters, setMasters] = useState<EmployeeProfileSummaryDto[]>([])
  const [selectedSalon, setSelectedSalon] = useState<SalonDto | null>(null)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  async function loadSalons() {
    setLoading(true)
    try {
      setSalons(await salonService.getAll())
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось загрузить салоны.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadSalons()
  }, [])

  async function handleSubmit(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    setMessage('')

    try {
      const payload = {
        ...form,
        max_staff: Number(form.max_staff),
        base_hourly_rate: Number(form.base_hourly_rate),
      }

      if (editingId) {
        await salonService.update(editingId, payload)
        setMessage(`Салон #${editingId} обновлён.`)
      } else {
        const created = await salonService.create(payload)
        setMessage(`Салон "${created.name}" создан.`)
      }

      setForm(initialForm)
      setEditingId(null)
      await loadSalons()
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось сохранить салон.')
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete(salon: SalonDto) {
    if (!window.confirm(`Удалить салон "${salon.name}"?`)) {
      return
    }

    setError('')
    setMessage('')
    try {
      await salonService.remove(salon.id)
      setMessage(`Салон "${salon.name}" удалён.`)
      if (selectedSalon?.id === salon.id) {
        setSelectedSalon(null)
        setMasters([])
      }
      await loadSalons()
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось удалить салон.')
    }
  }

  async function handleShowMasters(salon: SalonDto) {
    setSelectedSalon(salon)
    setError('')
    try {
      setMasters(await salonService.getMasters(salon.id))
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось загрузить мастеров салона.')
      setMasters([])
    }
  }

  return (
    <section className="page-section">
      <div className="page-header">
        <p className="eyebrow">Операционный контур</p>
        <h2>Салоны</h2>
        <p className="section-description">
          Здесь можно вести справочник филиалов, обновлять параметры работы и быстро смотреть мастеров по выбранному салону.
        </p>
      </div>

      {message ? <div className="alert alert-success">{message}</div> : null}
      {error ? <div className="alert alert-error">{error}</div> : null}

      <form className="card-form card-form-grid" onSubmit={handleSubmit}>
        <label className="field">
          <span>Название салона</span>
          <input value={form.name} onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))} placeholder="Bradobrei Party Center" required />
        </label>
        <label className="field">
          <span>Статус</span>
          <select value={form.status} onChange={(event) => setForm((current) => ({ ...current, status: event.target.value as UpsertSalonRequestDto['status'] }))}>
            {salonStatusOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
        <label className="field field-wide">
          <span>Адрес</span>
          <input value={form.address} onChange={(event) => setForm((current) => ({ ...current, address: event.target.value }))} placeholder="Екатеринбург, ул. Малышева, 12" required />
        </label>
        <label className="field">
          <span>Координаты</span>
          <input value={form.location || ''} onChange={(event) => setForm((current) => ({ ...current, location: event.target.value }))} placeholder="58.0141, 56.2230" />
        </label>
        <label className="field">
          <span>Макс. персонал</span>
          <input type="number" min="1" value={form.max_staff} onChange={(event) => setForm((current) => ({ ...current, max_staff: Number(event.target.value) }))} />
        </label>
        <label className="field">
          <span>Базовая ставка</span>
          <input type="number" min="0" value={form.base_hourly_rate} onChange={(event) => setForm((current) => ({ ...current, base_hourly_rate: Number(event.target.value) }))} />
        </label>
        <label className="field field-wide">
          <span>Часы работы JSON</span>
          <textarea rows={4} value={form.working_hours || ''} onChange={(event) => setForm((current) => ({ ...current, working_hours: event.target.value }))} />
        </label>
        <div className="button-row field-wide">
          <button type="submit" className="primary-button" disabled={submitting}>
            {submitting ? 'Сохраняем...' : editingId ? 'Обновить салон' : 'Создать салон'}
          </button>
          {editingId ? (
            <button
              type="button"
              className="ghost-button"
              onClick={() => {
                setEditingId(null)
                setForm(initialForm)
              }}
            >
              Сбросить редактирование
            </button>
          ) : null}
        </div>
      </form>

      <DataTable
        caption={loading ? 'Загружаем салоны...' : 'Справочник салонов'}
        columns={salonColumns(
          (salon) => {
            setEditingId(salon.id)
            setForm({
              name: salon.name,
              address: salon.address,
              location: salon.location || '',
              working_hours: salon.working_hours || '',
              status: salon.status,
              max_staff: salon.max_staff,
              base_hourly_rate: salon.base_hourly_rate,
            })
          },
          handleDelete,
          handleShowMasters,
        )}
        rows={salons}
        emptyText="Салоны пока не созданы или недоступны для текущей роли."
      />

      {selectedSalon ? (
        <DataTable
          caption={`Мастера салона: ${selectedSalon.name}`}
          columns={masterColumns}
          rows={masters}
          emptyText="Для этого салона пока нет закреплённых мастеров."
        />
      ) : null}
    </section>
  )
}
