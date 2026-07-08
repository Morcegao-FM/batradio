// Cliente HTTP do gateway: JSON, erros tipados e redirect para /login em 401.

export class ApiError extends Error {
  code: string
  status: number

  constructor(status: number, code: string, message: string) {
    super(message)
    this.code = code
    this.status = status
  }
}

// Status HTTP de um erro, por forma e não por instanceof: com HMR o módulo
// pode ser duplicado e `instanceof ApiError` falhar entre cópias.
export function apiErrorStatus(error: unknown): number | undefined {
  if (typeof error === 'object' && error !== null && 'status' in error) {
    const status = (error as { status: unknown }).status
    return typeof status === 'number' ? status : undefined
  }
  return undefined
}

// 502 (fora do ar) e 503 (ocupado/indexando) são transitórios: vale re-tentar.
export function isTransientError(error: unknown): boolean {
  const status = apiErrorStatus(error)
  return status === 502 || status === 503
}

export async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    headers: init?.body ? { 'Content-Type': 'application/json' } : undefined,
    ...init,
  })
  if (res.status === 401 && !path.startsWith('/api/me')) {
    window.location.assign('/login')
  }
  if (!res.ok) {
    let code = 'unknown'
    let message = `HTTP ${res.status}`
    try {
      const body = await res.json()
      if (body?.error) {
        code = body.error.code ?? code
        message = body.error.message ?? message
      }
    } catch {
      // corpo não-JSON: mantém o padrão
    }
    throw new ApiError(res.status, code, message)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return (await res.json()) as T
}

export function post<T>(path: string, body?: unknown): Promise<T> {
  return api<T>(path, {
    method: 'POST',
    body: body === undefined ? '{}' : JSON.stringify(body),
  })
}
