import { useEffect, useState } from 'react'
import { DataTable, type TableColumn } from '../components/DataTable'
import { serviceService } from '../api/services/serviceService'
import type { ServiceDto, UpsertServiceRequestDto } from '../types/dto/entities'
import { formatCurrency } from '../shared/formatters'

const initialForm: UpsertServiceRequestDto = {
  name: '',
  description: '',
  price: 1800,
  duration_minutes: 60,
}

const serviceColumns = (
  onEdit: (service: ServiceDto) => void,
  onDelete: (service: ServiceDto) => void,
): Array<TableColumn<ServiceDto>> => [
  { key: 'name', header: 'Услуга', render: (row) => row.name },
  { key: 'description', header: 'Описание', render: (row) => row.description || '—' },
  { key: 'price', header: 'Цена', render: (row) => formatCurrency(row.price) },
  { key: 'duration', header: 'Длительность', render: (row) => `${row.duration_minutes} мин.` },
  {
    key: 'materials',
    header: 'Материалы',
    render: (row) => row.materials?.map((item) => `${item.material?.name || item.material_id} x${item.quantity_per_use}`).join(', ') || '—',
  },
  {
    key: 'actions',
    header: 'Действия',
    render: (row) => (
      <div className="table-actions">
        <button type="button" className="ghost-button button-small" onClick={() => onEdit(row)}>
          Изменить
        </button>
        <button type="button" className="danger-button button-small" onClick={() => onDelete(row)}>
          Удалить
        </button>
      </div>
    ),
  },
]

export function ServicesPage() {
  const [form, setForm] = useState(initialForm)
  const [services, setServices] = useState<ServiceDto[]>([])
  const [editingId, setEditingId] = useState<number | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  async function loadServices() {
    setLoading(true)
    try {
      setServices(await serviceService.getAll())
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось загрузить услуги.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadServices()
  }, [])

  async function handleSubmit(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitting(true)
    setError('')
    setMessage('')

    try {
      const payload = {
        ...form,
        price: Number(form.price),
        duration_minutes: Number(form.duration_minutes),
      }

      if (editingId) {
        await serviceService.update(editingId, payload)
        setMessage(`Услуга #${editingId} обновлена.`)
      } else {
        const created = await serviceService.create(payload)
        setMessage(`Услуга "${created.name}" создана.`)
      }

      setForm(initialForm)
      setEditingId(null)
      await loadServices()
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось сохранить услугу.')
    } finally {
      setSubmitting(false)
    }
  }

  async function handleDelete(service: ServiceDto) {
    if (!window.confirm(`Удалить услугу "${service.name}"?`)) {
      return
    }

    try {
      await serviceService.remove(service.id)
      setMessage(`Услуга "${service.name}" удалена.`)
      await loadServices()
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось удалить услугу.')
    }
  }

  return (
    <section className="page-section">
      <div className="page-header">
        <p className="eyebrow">Прайс и выполнение</p>
        <h2>Услуги</h2>
        <p className="section-description">
          Экран для поддержки прайс-листа и связанных с ним параметров длительности, стоимости и расходников.
        </p>
      </div>

      {message ? <div className="alert alert-success">{message}</div> : null}
      {error ? <div className="alert alert-error">{error}</div> : null}

      <form className="card-form card-form-grid" onSubmit={handleSubmit}>
        <label className="field">
          <span>Название услуги</span>
          <input value={form.name} onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))} placeholder="Мужская стрижка" required />
        </label>
        <label className="field">
          <span>Цена</span>
          <input type="number" min="0" value={form.price} onChange={(event) => setForm((current) => ({ ...current, price: Number(event.target.value) }))} />
        </label>
        <label className="field">
          <span>Длительность, минут</span>
          <input type="number" min="15" step="5" value={form.duration_minutes} onChange={(event) => setForm((current) => ({ ...current, duration_minutes: Number(event.target.value) }))} />
        </label>
        <label className="field field-wide">
          <span>Описание</span>
          <textarea rows={4} value={form.description} onChange={(event) => setForm((current) => ({ ...current, description: event.target.value }))} placeholder="Стрижка с мытьём головы и укладкой." />
        </label>
        <div className="button-row field-wide">
          <button type="submit" className="primary-button" disabled={submitting}>
            {submitting ? 'Сохраняем...' : editingId ? 'Обновить услугу' : 'Создать услугу'}
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
        caption={loading ? 'Загружаем услуги...' : 'Справочник услуг'}
        columns={serviceColumns(
          (service) => {
            setEditingId(service.id)
            setForm({
              name: service.name,
              description: service.description,
              price: service.price,
              duration_minutes: service.duration_minutes,
            })
          },
          handleDelete,
        )}
        rows={services}
        emptyText="Услуги пока отсутствуют."
      />
    </section>
  )
}
