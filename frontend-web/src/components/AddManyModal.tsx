import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { post } from '../lib/api'
import type { Song } from '../lib/types'
import styles from './AddManyModal.module.css'

const INTERVALS = [1, 2, 3, 4, 6, 8, 12]

function toLocalInputValue(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

export default function AddManyModal({
  song,
  onClose,
  onDone,
}: {
  song: Song
  onClose: () => void
  onDone: (insertedCount: number) => void
}) {
  const [timesPerDay, setTimesPerDay] = useState(6)
  const [intervalHours, setIntervalHours] = useState(4)
  const [usePeriod, setUsePeriod] = useState(false)
  const [from, setFrom] = useState(() => toLocalInputValue(new Date()))
  const [to, setTo] = useState(() =>
    toLocalInputValue(new Date(Date.now() + 24 * 60 * 60 * 1000)),
  )
  const queryClient = useQueryClient()

  const insert = useMutation({
    mutationFn: () =>
      post<{ insertedAt: number[] }>('/api/playlist/add-many', {
        file: song.file,
        timesPerDay,
        intervalHours,
        ...(usePeriod
          ? { from: new Date(from).toISOString(), to: new Date(to).toISOString() }
          : {}),
      }),
    onSuccess: (res) => {
      queryClient.invalidateQueries({ queryKey: ['playlist'] })
      onDone(res.insertedAt.length)
    },
  })

  return (
    <div className={styles.backdrop} onClick={onClose}>
      <div
        className={styles.dialog}
        role="dialog"
        aria-label="Inserir música várias vezes"
        onClick={(e) => e.stopPropagation()}
      >
        <span className={styles.kicker}>Programação</span>
        <h2 className={styles.title}>Inserir música várias vezes</h2>
        <p className={styles.subtitle}>A faixa será distribuída automaticamente pela programação.</p>

        <div className={styles.trackCard}>
          <span className={styles.trackLabel}>Faixa selecionada</span>
          <strong>
            {song.artist ? `${song.artist} — ` : ''}
            {song.title}
          </strong>
        </div>

        <div className={styles.fields}>
          <label className={styles.field}>
            <span>Vezes por dia</span>
            <input
              type="number"
              min={1}
              max={96}
              value={timesPerDay}
              onChange={(e) => setTimesPerDay(Number(e.target.value))}
            />
          </label>
          <label className={styles.field}>
            <span>Intervalo</span>
            <select
              value={intervalHours}
              onChange={(e) => setIntervalHours(Number(e.target.value))}
            >
              {INTERVALS.map((h) => (
                <option key={h} value={h}>
                  A cada {h} hora{h > 1 ? 's' : ''}
                </option>
              ))}
            </select>
          </label>
        </div>

        <label className={styles.periodToggle}>
          <input
            type="checkbox"
            checked={usePeriod}
            onChange={(e) => setUsePeriod(e.target.checked)}
          />
          Apenas por um período
        </label>

        <div className={`${styles.fields} ${usePeriod ? '' : styles.disabled}`}>
          <label className={styles.field}>
            <span>De</span>
            <input
              type="datetime-local"
              value={from}
              disabled={!usePeriod}
              onChange={(e) => setFrom(e.target.value)}
            />
          </label>
          <label className={styles.field}>
            <span>Até</span>
            <input
              type="datetime-local"
              value={to}
              disabled={!usePeriod}
              onChange={(e) => setTo(e.target.value)}
            />
          </label>
        </div>

        {insert.isError && (
          <p className={styles.error} role="alert">
            Falha ao inserir: {(insert.error as Error).message}
          </p>
        )}

        <div className={styles.actions}>
          <button className={styles.cancel} onClick={onClose}>
            Cancelar
          </button>
          <button
            className={styles.confirm}
            onClick={() => insert.mutate()}
            disabled={insert.isPending || timesPerDay < 1 || timesPerDay > 96}
          >
            {insert.isPending ? 'Inserindo…' : 'Confirmar inserção'}
          </button>
        </div>
      </div>
    </div>
  )
}
