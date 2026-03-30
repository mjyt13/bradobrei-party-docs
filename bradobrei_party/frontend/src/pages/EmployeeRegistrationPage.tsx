import { useState } from 'react'
import { employeeService } from '../api/services/employeeService'
import { employeeRoleOptions } from '../shared/options'

const initialForm = {
  username: '',
  password: '',
  full_name: '',
  phone: '',
  email: '',
  role: 'ADVANCED_MASTER' as const,
  specialization: '',
  expected_salary: 65000,
  work_schedule: '{"mon":"10:00-19:00","wed":"10:00-19:00"}',
  salon_id: 1,
}

export function EmployeeRegistrationPage() {
  const [form, setForm] = useState(initialForm)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitting(true)
    setMessage('')
    setError('')

    try {
      const response = await employeeService.hire({
        ...form,
        expected_salary: Number(form.expected_salary),
        salon_id: Number(form.salon_id),
      })
      setMessage(`Сотрудник создан. Профиль #${response.id} готов к дальнейшей настройке.`)
      setForm(initialForm)
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось создать сотрудника.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="page-section">
      <div className="page-header">
        <div>
          <p className="eyebrow">Кадровый контур</p>
          <h2>Регистрация нового сотрудника</h2>
          <p className="section-description">
            Форма использует DTO backend для `HireEmployeeRequest` и сразу готова
            к работе в локальной разработке.
          </p>
        </div>
      </div>

      {message ? <div className="alert alert-success">{message}</div> : null}
      {error ? <div className="alert alert-error">{error}</div> : null}

      <form className="card-form card-form-grid" onSubmit={handleSubmit}>
        <label className="field">
          <span>Логин</span>
          <input
            value={form.username}
            onChange={(event) => setForm((current) => ({ ...current, username: event.target.value }))}
            placeholder="master_ivan"
            required
          />
        </label>
        <label className="field">
          <span>Пароль</span>
          <input
            type="password"
            value={form.password}
            onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))}
            placeholder="password"
            required
          />
        </label>
        <label className="field field-wide">
          <span>ФИО</span>
          <input
            value={form.full_name}
            onChange={(event) => setForm((current) => ({ ...current, full_name: event.target.value }))}
            placeholder="Иван Барбер"
            required
          />
        </label>
        <label className="field">
          <span>Телефон</span>
          <input
            value={form.phone}
            onChange={(event) => setForm((current) => ({ ...current, phone: event.target.value }))}
            placeholder="+79990001122"
            required
          />
        </label>
        <label className="field">
          <span>Email</span>
          <input
            type="email"
            value={form.email}
            onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))}
            placeholder="ivan.barber@example.com"
            required
          />
        </label>
        <label className="field">
          <span>Роль</span>
          <select
            value={form.role}
            onChange={(event) =>
              setForm((current) => ({ ...current, role: event.target.value as typeof form.role }))
            }
          >
            {employeeRoleOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
        <label className="field">
          <span>Ожидаемый оклад</span>
          <input
            type="number"
            min="0"
            value={form.expected_salary}
            onChange={(event) =>
              setForm((current) => ({
                ...current,
                expected_salary: Number(event.target.value),
              }))
            }
          />
        </label>
        <label className="field field-wide">
          <span>Специализация</span>
          <input
            value={form.specialization}
            onChange={(event) => setForm((current) => ({ ...current, specialization: event.target.value }))}
            placeholder="Fade, beard styling"
          />
        </label>
        <label className="field">
          <span>ID салона</span>
          <input
            type="number"
            min="1"
            value={form.salon_id}
            onChange={(event) => setForm((current) => ({ ...current, salon_id: Number(event.target.value) }))}
          />
        </label>
        <label className="field field-wide">
          <span>График в JSON</span>
          <textarea
            rows={5}
            value={form.work_schedule}
            onChange={(event) => setForm((current) => ({ ...current, work_schedule: event.target.value }))}
          />
        </label>
        <button type="submit" className="primary-button field-wide" disabled={submitting}>
          {submitting ? 'Создаём сотрудника...' : 'Создать сотрудника'}
        </button>
      </form>
    </section>
  )
}
