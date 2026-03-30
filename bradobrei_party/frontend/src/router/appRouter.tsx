import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppShell } from '../components/AppShell'
import { ProtectedRoute } from '../components/ProtectedRoute'
import { AuthPage } from '../pages/AuthPage'
import { BookingsPage } from '../pages/BookingsPage'
import { EmployeesPage } from '../pages/EmployeesPage'
import { EmployeeRegistrationPage } from '../pages/EmployeeRegistrationPage'
import { EmployeesReportPage } from '../pages/EmployeesReportPage'
import { MaterialsPage } from '../pages/MaterialsPage'
import { MasterActivityReportPage } from '../pages/MasterActivityReportPage'
import { PaymentsPage } from '../pages/PaymentsPage'
import { SalonActivityReportPage } from '../pages/SalonActivityReportPage'
import { SalonsPage } from '../pages/SalonsPage'
import { ServicesPage } from '../pages/ServicesPage'

export const appRouter = createBrowserRouter([
  {
    path: '/auth',
    element: <AuthPage />,
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <AppShell />,
        children: [
          { index: true, element: <Navigate to="/salons" replace /> },
          { path: '/salons', element: <SalonsPage /> },
          { path: '/bookings', element: <BookingsPage /> },
          { path: '/services', element: <ServicesPage /> },
          { path: '/materials', element: <MaterialsPage /> },
          { path: '/payments', element: <PaymentsPage /> },
          { path: '/employees', element: <EmployeesPage /> },
          { path: '/reports/employees', element: <EmployeesReportPage /> },
          { path: '/reports/salon-activity', element: <SalonActivityReportPage /> },
          { path: '/reports/master-activity', element: <MasterActivityReportPage /> },
          { path: '/employees/new', element: <EmployeeRegistrationPage /> },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <Navigate to="/" replace />,
  },
])
