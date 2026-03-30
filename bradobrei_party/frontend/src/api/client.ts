import { apiConfig } from './config'
import { tokenStorage } from './services/tokenStorage'
import type { ErrorResponseDto } from '../types/dto/common'

type Primitive = string | number | boolean | null | undefined
type QueryParams = Record<string, Primitive>

interface RequestOptions extends Omit<RequestInit, 'body'> {
  auth?: boolean
  body?: unknown
  query?: QueryParams
}

export class ApiError extends Error {
  status: number
  payload?: ErrorResponseDto

  constructor(status: number, message: string, payload?: ErrorResponseDto) {
    super(message)
    this.status = status
    this.payload = payload
  }
}

function buildUrl(path: string, query?: QueryParams) {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  const url = new URL(`${apiConfig.apiBaseUrl}${normalizedPath}`, window.location.origin)

  if (query) {
    for (const [key, value] of Object.entries(query)) {
      if (value !== undefined && value !== null && value !== '') {
        url.searchParams.set(key, String(value))
      }
    }
  }

  return url.toString()
}

async function parseResponse(response: Response) {
  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    return response.json()
  }
  return response.text()
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { auth = true, body, headers, query, ...rest } = options
  const token = tokenStorage.get()
  const requestHeaders = new Headers(headers)

  if (!requestHeaders.has('Accept')) {
    requestHeaders.set('Accept', 'application/json')
  }

  if (body !== undefined) {
    requestHeaders.set('Content-Type', 'application/json')
  }

  if (auth && token) {
    requestHeaders.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(buildUrl(path, query), {
    ...rest,
    headers: requestHeaders,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  const payload = await parseResponse(response)

  if (!response.ok) {
    if (response.status === 401) {
      tokenStorage.clear()
    }

    const errorPayload = typeof payload === 'object' && payload !== null ? (payload as ErrorResponseDto) : undefined
    const message = errorPayload?.message || errorPayload?.error || response.statusText
    throw new ApiError(response.status, message, errorPayload)
  }

  return payload as T
}
