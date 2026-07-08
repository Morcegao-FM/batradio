import { useEffect, useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { post } from '../lib/api'
import styles from './HeaderBar.module.css'

const WEEKDAYS = ['DOM', 'SEG', 'TER', 'QUA', 'QUI', 'SEX', 'SÁB']

function Clock() {
  const [now, setNow] = useState(() => new Date())
  useEffect(() => {
    const t = setInterval(() => setNow(new Date()), 1000)
    return () => clearInterval(t)
  }, [])
  const hh = String(now.getHours()).padStart(2, '0')
  const mm = String(now.getMinutes()).padStart(2, '0')
  return (
    <span className={styles.clock}>
      {hh}:{mm} · {WEEKDAYS[now.getDay()]}
    </span>
  )
}

export default function HeaderBar({ title, subtitle }: { title: string; subtitle: string }) {
  const queryClient = useQueryClient()
  const refresh = useMutation({
    mutationFn: () => post<{ count: number }>('/api/library/refresh'),
    onSuccess: () => queryClient.invalidateQueries(),
  })

  return (
    <header className={styles.header}>
      <div className={styles.titles}>
        <h1 className={styles.title}>{title}</h1>
        <span className={styles.subtitle}>{subtitle}</span>
      </div>
      <div className={styles.right}>
        <Clock />
        <button
          className={styles.refresh}
          onClick={() => refresh.mutate()}
          disabled={refresh.isPending}
        >
          ↻ {refresh.isPending ? 'Atualizando…' : 'Atualizar'}
        </button>
      </div>
    </header>
  )
}
