import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'
import type { ServerInfo, Status } from '../lib/types'

// O SSE (useRadioEvents) mantém este cache atualizado; a query só faz o fetch
// inicial.
export function useStatus() {
  const { data, error } = useQuery({
    queryKey: ['status'],
    queryFn: () => api<Status>('/api/status'),
    staleTime: Infinity,
    retry: 1,
  })
  return { status: data, error }
}

export function useServerInfo() {
  const { data } = useQuery({
    queryKey: ['server-info'],
    queryFn: () => api<ServerInfo>('/api/server/info'),
    staleTime: 60_000,
  })
  return { info: data }
}
