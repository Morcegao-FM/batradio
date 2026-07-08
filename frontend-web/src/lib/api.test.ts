import { afterEach, describe, expect, it, vi } from 'vitest'
import { api, ApiError, post } from './api'

const fetchMock = vi.fn()
vi.stubGlobal('fetch', fetchMock)

afterEach(() => {
  fetchMock.mockReset()
})

function jsonResponse(status: number, body: unknown) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('api', () => {
  it('retorna o JSON tipado em 200', async () => {
    fetchMock.mockResolvedValue(jsonResponse(200, { email: 'a@b.com' }))
    const out = await api<{ email: string }>('/api/me')
    expect(out.email).toBe('a@b.com')
  })

  it('converte erro JSON do gateway em ApiError', async () => {
    fetchMock.mockResolvedValue(
      jsonResponse(502, { error: { code: 'node_unavailable', message: 'Node fora do ar' } }),
    )
    const err = (await api('/api/status').catch((e) => e)) as ApiError
    expect(err).toBeInstanceOf(ApiError)
    expect(err.code).toBe('node_unavailable')
    expect(err.status).toBe(502)
    expect(err.message).toBe('Node fora do ar')
  })

  it('trata 204 sem corpo', async () => {
    fetchMock.mockResolvedValue(new Response(null, { status: 204 }))
    await expect(post('/api/playlist/move', { from: 1, to: 2 })).resolves.toBeUndefined()
    const [, init] = fetchMock.mock.calls[0]
    expect(init.method).toBe('POST')
    expect(JSON.parse(init.body)).toEqual({ from: 1, to: 2 })
  })
})
