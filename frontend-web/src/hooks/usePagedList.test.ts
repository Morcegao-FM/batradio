import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest'
import { act, createElement } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useDebounced } from './usePagedList'

function DebounceProbe({ value, out }: { value: string; out: (v: string) => void }) {
  out(useDebounced(value, 300))
  return null
}

describe('useDebounced', () => {
  let container: HTMLDivElement
  let root: ReturnType<typeof createRoot>
  const queryClient = new QueryClient()

  beforeEach(() => {
    vi.useFakeTimers()
    container = document.createElement('div')
    document.body.appendChild(container)
    root = createRoot(container)
  })

  afterEach(() => {
    act(() => root.unmount())
    container.remove()
    vi.useRealTimers()
  })

  it('só propaga o valor após 300ms sem digitação', async () => {
    let latest = ''
    const render = (value: string) =>
      act(async () => {
        root.render(
          createElement(
            QueryClientProvider,
            { client: queryClient },
            createElement(DebounceProbe, { value, out: (v) => (latest = v) }),
          ),
        )
      })

    await render('z')
    expect(latest).toBe('z')
    // digitação rápida: valores intermediários não propagam
    await render('ze')
    await act(async () => vi.advanceTimersByTime(100))
    await render('zep')
    await act(async () => vi.advanceTimersByTime(100))
    expect(latest).toBe('z')
    // 300ms depois do último caractere, propaga o final
    await act(async () => vi.advanceTimersByTime(300))
    expect(latest).toBe('zep')
  })
})
