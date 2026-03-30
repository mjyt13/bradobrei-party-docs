import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { tokenStorage } from '../api/services/tokenStorage'

export function ProtectedRoute() {
  const location = useLocation()

  if (!tokenStorage.has()) {
    return <Navigate to="/auth" replace state={{ from: location.pathname }} />
  }

  return <Outlet />
}
