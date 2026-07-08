// @vitest-environment jsdom
import { QueryClient } from '@tanstack/react-query'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// EventSource fake controlável
class FakeEventSource {
  static instances: FakeEventSource[] = []
  url: string
  listeners = new Map<string, ((e: MessageEvent) => void)[]>()
  onopen: (() => void) | null = null
  onerror: (() => void) | null = null
  closed = false

  constructor(url: string) {
    this.url = url
    FakeEventSource.instances.push(this)
  }

  addEventListener(name: string, fn: (e: MessageEvent) => void) {
    this.listeners.set(name, [...(this.listeners.get(name) ?? []), fn])
  }

  emit(name: string, data: unknown) {
    for (const fn of this.listeners.get(name) ?? []) {
      fn(new MessageEvent(name, { data: JSON.stringify(data) }))
    }
  }

  close() {
    this.closed = true
  }
}

vi.stubGlobal('EventSource', FakeEventSource)

// Renderiza o hook sem @testing-library/react-hooks: componente mínimo.
import { act, createElement } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClientProvider } from '@tanstack/react-query'
import { useRadioEvents } from './useRadioEvents'

function Probe() {
  useRadioEvents()
  return null
}

describe('useRadioEvents', () => {
  let queryClient: QueryClient
  let container: HTMLDivElement
  let root: ReturnType<typeof createRoot>

  beforeEach(() => {
    FakeEventSource.instances = []
    queryClient = new QueryClient()
    container = document.createElement('div')
    document.body.appendChild(container)
    root = createRoot(container)
  })

  afterEach(() => {
    root.unmount()
    container.remove()
  })

  it('alimenta o cache de status e invalida a fila', async () => {
    const invalidate = vi.spyOn(queryClient, 'invalidateQueries')
    await act(async () => {
      root.render(
        createElement(QueryClientProvider, { client: queryClient }, createElement(Probe)),
      )
    })

    const es = FakeEventSource.instances[0]
    expect(es).toBeDefined()
    expect(es.url).toBe('/api/events')

    es.emit('status', { state: 'play', song: 3 })
    expect(queryClient.getQueryData(['status'])).toMatchObject({ state: 'play', song: 3 })

    es.emit('playlist', { version: 7 })
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ['playlist'] })
    expect(invalidate).toHaveBeenCalledWith({ queryKey: ['playlists'] })
  })
})
