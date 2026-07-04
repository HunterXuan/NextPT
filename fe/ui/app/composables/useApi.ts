import type { UseFetchOptions } from 'nuxt/app'

type ApiEnvelope<T> = {
  code: number
  message?: string
  data?: T
}

export class ApiError extends Error {
  public code: number
  public status?: number
  public data: unknown

  constructor(message: string, code: number, data?: unknown, status?: number) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.data = data
    this.status = status
  }
}

export interface ApiBlobResponse {
  blob: Blob
  filename?: string
}

const supportedLocales = new Set(['zh-CN', 'zh-TW', 'en-US'])
let currentAuthToken: string | null = null

export function setApiAuthToken(token: string | null | undefined) {
  currentAuthToken = token || null
}

function readCookie(name: string) {
  if (!import.meta.client) return null

  const item = document.cookie
    .split('; ')
    .find((row) => row.startsWith(`${encodeURIComponent(name)}=`))

  if (!item) return null

  return decodeURIComponent(item.slice(name.length + 1))
}

function currentLocale() {
  if (import.meta.client) {
    const htmlLang = document.documentElement.lang
    if (supportedLocales.has(htmlLang)) return htmlLang

    const pathLocale = window.location.pathname.split('/').filter(Boolean)[0] || ''
    if (supportedLocales.has(pathLocale)) return pathLocale

    const cookieLocale = readCookie('nextpt_locale')
    if (cookieLocale && supportedLocales.has(cookieLocale)) return cookieLocale
  }

  return 'zh-CN'
}

function useApiHeaders(extra?: HeadersInit) {
  const headers = new Headers(extra)
  const token = currentAuthToken || readCookie('nextpt_token')

  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  headers.set('Accept-Language', currentLocale())
  return headers
}

function clearSession() {
  setApiAuthToken(null)

  if (import.meta.client) {
    document.cookie = 'nextpt_token=; Max-Age=0; path=/; SameSite=Lax'
  }

  try {
    const token = useCookie<string | null>('nextpt_token')
    token.value = null
  } catch {
    // The caller may be outside setup; the browser cookie has already been cleared.
  }

  try {
    const user = useState('auth:user', () => null)
    user.value = null
  } catch {
    // Ignore context cleanup failures in low-level API handling.
  }
}

function unwrapEnvelope<T>(payload: ApiEnvelope<T> | T): T {
  if (payload && typeof payload === 'object' && 'code' in payload) {
    const envelope = payload as ApiEnvelope<T>
    if (envelope.code !== 0) {
      if (envelope.code === 401) {
        clearSession()
      }

      throw new ApiError(envelope.message || 'Request failed', envelope.code, envelope.data)
    }

    return envelope.data as T
  }

  return payload as T
}

function parseContentDispositionFilename(value: string | null) {
  if (!value) return undefined

  const encoded = value.match(/filename\*=UTF-8''([^;]+)/i)
  if (encoded?.[1]) {
    try {
      return decodeURIComponent(encoded[1].trim())
    } catch {
      return encoded[1].trim()
    }
  }

  const quoted = value.match(/filename="([^"]+)"/i)
  if (quoted?.[1]) return quoted[1].trim()

  const bare = value.match(/filename=([^;]+)/i)
  if (bare?.[1]) return bare[1].trim()

  return undefined
}

async function normalizeFetchError(error: any): Promise<ApiError> {
  const status = error?.response?.status
  let data = error?.response?._data

  if (status === 401) {
    clearSession()
  }

  if (data instanceof Blob) {
    try {
      const text = await data.text()
      data = JSON.parse(text)
    } catch {
      data = undefined
    }
  }

  if (data && typeof data === 'object') {
    const envelope = data as ApiEnvelope<unknown>
    return new ApiError(
      envelope.message || error.message || 'Request failed',
      Number(envelope.code || status || -1),
      envelope.data,
      status
    )
  }

  if (error instanceof ApiError) {
    return error
  }

  return new ApiError(error?.message || 'Request failed', Number(status || -1), undefined, status)
}

export function useApi<T = unknown>(
  request: Parameters<typeof useFetch>[0],
  opts?: UseFetchOptions<ApiEnvelope<T> | T, T>
) {
  const options = {
    ...opts,
    headers: useApiHeaders(opts?.headers as HeadersInit),
    transform: (payload) => unwrapEnvelope<T>(payload)
  } satisfies UseFetchOptions<ApiEnvelope<T> | T, T>

  return useFetch<ApiEnvelope<T> | T, ApiError, any, any, ApiEnvelope<T> | T, T>(request as any, options)
}

export async function fetchApi<T = unknown>(
  request: Parameters<typeof $fetch>[0],
  opts?: Parameters<typeof $fetch>[1]
): Promise<T> {
  try {
    const payload = await $fetch<ApiEnvelope<T> | T>(request, {
      ...opts,
      headers: useApiHeaders(opts?.headers as HeadersInit)
    })

    return unwrapEnvelope<T>(payload)
  } catch (error: any) {
    throw await normalizeFetchError(error)
  }
}

export async function fetchApiBlob(
  request: Parameters<typeof $fetch>[0],
  opts?: Parameters<typeof $fetch>[1]
): Promise<ApiBlobResponse> {
  try {
    const response = await $fetch.raw<Blob>(request, {
      ...opts,
      responseType: 'blob',
      headers: useApiHeaders(opts?.headers as HeadersInit)
    })
    const blob = response._data
    if (!blob) {
      throw new ApiError('Empty response', -1)
    }

    const contentType = response.headers.get('content-type') || blob.type || ''
    if (contentType.includes('application/json')) {
      const text = await blob.text()
      const payload = JSON.parse(text) as ApiEnvelope<unknown>
      unwrapEnvelope(payload)
      throw new ApiError('Unexpected JSON response', -1)
    }

    return {
      blob,
      filename: parseContentDispositionFilename(response.headers.get('content-disposition'))
    }
  } catch (error: any) {
    throw await normalizeFetchError(error)
  }
}
