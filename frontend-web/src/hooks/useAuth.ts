import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'

export function useAuth(): { email?: string; loading: boolean } {
  const { data, isLoading } = useQuery({
    queryKey: ['me'],
    queryFn: () => api<{ email: string }>('/api/me'),
    retry: false,
    staleTime: Infinity,
  })
  return { email: data?.email, loading: isLoading }
}
