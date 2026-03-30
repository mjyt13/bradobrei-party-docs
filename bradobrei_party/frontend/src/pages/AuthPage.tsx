import { startTransition, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { authService } from '../api/services/authService'
import { tokenStorage } from '../api/services/tokenStorage'

const initialLoginForm = {
  username: '',
  password: '',
}

const initialRegisterForm = {
  username: '',
  password: '',
  full_name: '',
  phone: '',
  email: '',
}

export function AuthPage() {
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [loginForm, setLoginForm] = useState(initialLoginForm)
  const [registerForm, setRegisterForm] = useState(initialRegisterForm)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()

  const redirectTo = (location.state as { from?: string } | null)?.from || '/reports/employees'

  async function handleLoginSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitting(true)
    setMessage('')
    setError('')

    try {
      const response = await authService.login(loginForm)
      tokenStorage.set(response.token)
      startTransition(() => {
        navigate(redirectTo, { replace: true })
      })
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось войти.')
    } finally {
      setSubmitting(false)
    }
  }

  async function handleRegisterSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitting(true)
    setMessage('')
    setError('')

    try {
      await authService.register(registerForm)
      setMode('login')
      setMessage('Пользователь создан. Теперь можно выполнить вход.')
      setRegisterForm(initialRegisterForm)
    } catch (requestError) {
      setError(requestError instanceof Error ? requestError.message : 'Не удалось зарегистрироваться.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="auth-layout">
      <section className="auth-panel auth-panel-highlight">
        <p className="eyebrow">Локальная разработка</p>
        <h1>Один вход для отчётов, найма и дальнейших операций</h1>
        <p className="lede">
          Токен сохраняется в <code>localStorage</code>, а затем автоматически
          подставляется в сервисы API. Для разработки можно быстро
          зарегистрировать клиента и тут же войти.
        </p>
        <div className="tip-list">
          <div className="tip-card">
            <strong>Регистрация</strong>
            <span>Создаёт обычного пользователя через backend DTO.</span>
          </div>
          <div className="tip-card">
            <strong>Вход</strong>
            <span>Сохраняет JWT и переводит в раздел отчётов.</span>
          </div>
          <div className="tip-card">
            <strong>Swagger</strong>
            <span>Бэкенд остаётся доступен отдельно, а фронтенд ходит через proxy Vite.</span>
          </div>
        </div>
      </section>

      <section className="auth-panel">
        <div className="tab-row">
          <button
            type="button"
            className={mode === 'login' ? 'tab-button tab-button-active' : 'tab-button'}
            onClick={() => setMode('login')}
          >
            Вход
          </button>
          <button
            type="button"
            className={mode === 'register' ? 'tab-button tab-button-active' : 'tab-button'}
            onClick={() => setMode('register')}
          >
            Регистрация
          </button>
        </div>

        {message ? <div className="alert alert-success">{message}</div> : null}
        {error ? <div className="alert alert-error">{error}</div> : null}

        {mode === 'login' ? (
          <form className="stack-form" onSubmit={handleLoginSubmit}>
            <label className="field">
              <span>Логин</span>
              <input
                value={loginForm.username}
                onChange={(event) =>
                  setLoginForm((current) => ({ ...current, username: event.target.value }))
                }
                placeholder="admin"
                required
              />
            </label>

            <label className="field">
              <span>Пароль</span>
              <input
                type="password"
                value={loginForm.password}
                onChange={(event) =>
                  setLoginForm((current) => ({ ...current, password: event.target.value }))
                }
                placeholder="password"
                required
              />
            </label>

            <button type="submit" className="primary-button" disabled={submitting}>
              {submitting ? 'Выполняем вход...' : 'Войти'}
            </button>
          </form>
        ) : (
          <form className="stack-form" onSubmit={handleRegisterSubmit}>
            <label className="field">
              <span>Логин</span>
              <input
                value={registerForm.username}
                onChange={(event) =>
                  setRegisterForm((current) => ({ ...current, username: event.target.value }))
                }
                placeholder="client_ivan"
                required
              />
            </label>

            <label className="field">
              <span>Пароль</span>
              <input
                type="password"
                value={registerForm.password}
                onChange={(event) =>
                  setRegisterForm((current) => ({ ...current, password: event.target.value }))
                }
                placeholder="password"
                required
              />
            </label>

            <label className="field">
              <span>ФИО</span>
              <input
                value={registerForm.full_name}
                onChange={(event) =>
                  setRegisterForm((current) => ({ ...current, full_name: event.target.value }))
                }
                placeholder="Иван Петров"
                required
              />
            </label>

            <label className="field">
              <span>Телефон</span>
              <input
                value={registerForm.phone}
                onChange={(event) =>
                  setRegisterForm((current) => ({ ...current, phone: event.target.value }))
                }
                placeholder="+79991234567"
                required
              />
            </label>

            <label className="field">
              <span>Email</span>
              <input
                type="email"
                value={registerForm.email}
                onChange={(event) =>
                  setRegisterForm((current) => ({ ...current, email: event.target.value }))
                }
                placeholder="client@example.com"
                required
              />
            </label>

            <button type="submit" className="primary-button" disabled={submitting}>
              {submitting ? 'Создаём пользователя...' : 'Зарегистрироваться'}
            </button>
          </form>
        )}
      </section>
    </div>
  )
}
