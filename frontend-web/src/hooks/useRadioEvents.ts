import { useEffect, useRef, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import type { Status } from '../lib/types'

// Mantém uma conexão SSE com o gateway: eventos "status" alimentam o cache da
// query ['status']; eventos "playlist" invalidam a fila e as playlists salvas.
// Reconecta com backoff exponencial (1s → 30s).
export function useRadioEvents(): { connected: boolean } {
  const queryClient = useQueryClient()
  const [connected, setConnected] = useState(false)
  const backoffRef = useRef(1000)

  useEffect(() => {
    let source: EventSource | null = null
    let retryTimer: ReturnType<typeof setTimeout> | undefined
    let disposed = false

    const connect = () => {
      if (disposed) return
      source = new EventSource('/api/events')

      source.onopen = () => {
        backoffRef.current = 1000
        setConnected(true)
      }
      source.addEventListener('status', (e) => {
        const status = JSON.parse((e as MessageEvent).data) as Status
        queryClient.setQueryData(['status'], status)
      })
      source.addEventListener('playlist', () => {
        queryClient.invalidateQueries({ queryKey: ['playlist'] })
        queryClient.invalidateQueries({ queryKey: ['playlists'] })
      })
      source.onerror = () => {
        setConnected(false)
        source?.close()
        retryTimer = setTimeout(connect, backoffRef.current)
        backoffRef.current = Math.min(backoffRef.current * 2, 30000)
      }
    }

    connect()
    return () => {
      disposed = true
      clearTimeout(retryTimer)
      source?.close()
    }
  }, [queryClient])

  return { connected }
}
