import { useEffect, useState } from 'react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { authService } from '../api/services/authService'
import { tokenStorage } from '../api/services/tokenStorage'
import { apiConfig } from '../api/config'
import type { UserDto } from '../types/dto/auth'
import { formatRole } from '../shared/formatters'

const navigation = [
  { to: '/salons', label: 'Салоны' },
  { to: '/bookings', label: 'Бронирования' },
  { to: '/services', label: 'Услуги' },
  { to: '/materials', label: 'Материалы' },
  { to: '/payments', label: 'Платежи' },
  { to: '/employees', label: 'Сотрудники' },
  { to: '/reports/employees', label: 'Отчёт 2.2.1' },
  { to: '/reports/salon-activity', label: 'Отчёт 2.2.2' },
  { to: '/reports/master-activity', label: 'Отчёт 2.2.4' },
  { to: '/employees/new', label: 'Новый сотрудник' },
]

export function AppShell() {
  const [currentUser, setCurrentUser] = useState<UserDto | null>(null)
  const [error, setError] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    let cancelled = false

    authService
      .getMe()
      .then((response) => {
        if (!cancelled) {
          setCurrentUser(response.user)
        }
      })
      .catch((requestError: Error) => {
        if (!cancelled) {
          setError(requestError.message)
          tokenStorage.clear()
          navigate('/auth', { replace: true })
        }
      })

    return () => {
      cancelled = true
    }
  }, [navigate])

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand-block">
          <p className="eyebrow">Bradobrei Party</p>
          <h1>Панель операций и отчётов</h1>
          <p className="lede">
            Один контур для авторизации, справочников, бронирований, платежей, кадровых операций и отчётов,
            уже подключённых к текущему backend API.
          </p>
        </div>

        <nav className="navigation">
          {navigation.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) =>
                isActive ? 'nav-link nav-link-active' : 'nav-link'
              }
            >
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="sidebar-footer">
          <a href={apiConfig.docsUrl} target="_blank" rel="noreferrer">
            Swagger backend
          </a>
        </div>
      </aside>

      <main className="content">
        <header className="topbar">
          <div>
            <p className="topbar-label">Текущий пользователь</p>
            {currentUser ? (
              <>
                <strong>{currentUser.full_name}</strong>
                <span className="topbar-meta">{formatRole(currentUser.role)}</span>
              </>
            ) : (
              <span className="topbar-meta">Загрузка профиля...</span>
            )}
          </div>

          <button
            type="button"
            className="ghost-button"
            onClick={() => {
              tokenStorage.clear()
              navigate('/auth', { replace: true })
            }}
          >
            Выйти
          </button>
        </header>

        {error ? <div className="alert alert-error">{error}</div> : null}
        <Outlet context={{ currentUser }} />
      </main>
    </div>
  )
}
