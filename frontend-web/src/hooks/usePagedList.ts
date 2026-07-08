import { useQueries } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import type { Page } from '../lib/types'

export const PAGE_SIZE = 100

// Debounce genérico (busca server-side).
export function useDebounced<T>(value: T, ms = 300): T {
  const [debounced, setDebounced] = useState(value)
  useEffect(() => {
    const t = setTimeout(() => setDebounced(value), ms)
    return () => clearTimeout(t)
  }, [value, ms])
  return debounced
}

export interface VisibleRange {
  start: number
  end: number
}

// Lista paginada com acesso aleatório: busca apenas as páginas visíveis no
// virtualizador (+1 de prefetch). Permite pular direto para qualquer índice
// (ex.: "música atual") sem carregar as páginas anteriores.
export function usePagedList<T>(
  baseKey: readonly unknown[],
  fetchPage: (offset: number, limit: number) => Promise<Page<T>>,
  range: VisibleRange,
) {
  const firstPage = Math.max(0, Math.floor(range.start / PAGE_SIZE))
  const lastPage = Math.max(firstPage, Math.floor(range.end / PAGE_SIZE)) + 1 // prefetch
  const pageIndices: number[] = []
  for (let p = firstPage; p <= lastPage; p++) {
    pageIndices.push(p)
  }

  const results = useQueries({
    queries: pageIndices.map((p) => ({
      queryKey: [...baseKey, p],
      queryFn: () => fetchPage(p * PAGE_SIZE, PAGE_SIZE),
      staleTime: 30_000,
    })),
  })

  const firstLoaded = results.find((r) => r.data)?.data
  const total = firstLoaded?.total ?? 0

  const itemAt = (index: number): T | undefined => {
    const p = Math.floor(index / PAGE_SIZE)
    const slot = pageIndices.indexOf(p)
    if (slot < 0) return undefined
    return results[slot].data?.items[index - p * PAGE_SIZE]
  }

  return {
    total,
    itemAt,
    firstLoaded,
    isFetching: results.some((r) => r.isFetching),
  }
}
