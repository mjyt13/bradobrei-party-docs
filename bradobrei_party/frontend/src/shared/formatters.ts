export function formatCurrency(value: number) {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    maximumFractionDigits: 0,
  }).format(value)
}

export function formatDateTime(value?: string | null) {
  if (!value) {
    return '—'
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat('ru-RU', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

export function formatRole(role: string) {
  const labels: Record<string, string> = {
    CLIENT: 'Клиент',
    BASIC_MASTER: 'Базовый мастер',
    ADVANCED_MASTER: 'Старший мастер',
    HR_SPECIALIST: 'HR-специалист',
    ACCOUNTANT: 'Бухгалтер',
    NETWORK_MANAGER: 'Менеджер сети',
    ADMINISTRATOR: 'Администратор',
  }

  return labels[role] || role
}

export function formatBookingStatus(status: string) {
  const labels: Record<string, string> = {
    PENDING: 'Ожидает',
    CONFIRMED: 'Подтверждено',
    IN_PROGRESS: 'В работе',
    COMPLETED: 'Завершено',
    CANCELLED: 'Отменено',
  }

  return labels[status] || status
}

export function formatPaymentStatus(status: string) {
  const labels: Record<string, string> = {
    PENDING: 'Ожидает оплаты',
    SUCCESS: 'Оплачен',
    FAILED: 'Ошибка оплаты',
    REFUNDED: 'Возврат',
  }

  return labels[status] || status
}

export function formatSalonStatus(status: string) {
  const labels: Record<string, string> = {
    OPEN: 'Открыт',
    CLOSED: 'Закрыт',
  }

  return labels[status] || status
}

export function formatJsonPreview(value?: string | null) {
  if (!value) {
    return '—'
  }

  return value.length > 48 ? `${value.slice(0, 48)}...` : value
}
