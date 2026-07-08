import { useRef, useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useVirtualizer } from '@tanstack/react-virtual'
import { api, post } from '../lib/api'
import type { QueueItem, QueuePage, Status } from '../lib/types'
import { formatClock, formatDuration } from '../lib/format'
import { useDebounced, usePagedList, type VisibleRange } from '../hooks/usePagedList'
import { useStatus } from '../hooks/useStatus'
import GenreChip from './GenreChip'
import ConfirmDialog from './ConfirmDialog'
import styles from './Panel.module.css'

const ROW_HEIGHT = 56

export default function QueuePanel({
  selected,
  onSelect,
}: {
  selected?: QueueItem
  onSelect: (item: QueueItem) => void
}) {
  const [query, setQuery] = useState('')
  const q = useDebounced(query)
  const [range, setRange] = useState<VisibleRange>({ start: 0, end: 40 })
  const { status } = useStatus()
  const queryClient = useQueryClient()

  const { total, itemAt, firstLoaded } = usePagedList<QueueItem>(
    ['playlist', q],
    (offset, limit) =>
      api<QueuePage>(
        `/api/playlist?q=${encodeURIComponent(q)}&offset=${offset}&limit=${limit}`,
      ),
    range,
  )
  const currentPos = (firstLoaded as QueuePage | undefined)?.currentPos ?? status?.song

  const parentRef = useRef<HTMLDivElement>(null)
  const virtualizer = useVirtualizer({
    count: total,
    getScrollElement: () => parentRef.current,
    estimateSize: () => ROW_HEIGHT,
    overscan: 10,
    onChange: (v) => {
      const items = v.getVirtualItems()
      if (items.length === 0) return
      const start = items[0].index
      const end = items[items.length - 1].index
      setRange((r) => (r.start === start && r.end === end ? r : { start, end }))
    },
  })

  const invalidateQueue = () => queryClient.invalidateQueries({ queryKey: ['playlist'] })

  const move = useMutation({
    mutationFn: ({ from, to }: { from: number; to: number }) =>
      post('/api/playlist/move', { from, to }),
    onSuccess: invalidateQueue,
  })
  const remove = useMutation({
    mutationFn: (pos: number) => post('/api/playlist/remove', { positions: [pos] }),
    onSuccess: invalidateQueue,
  })
  const play = useMutation({
    mutationFn: (pos: number) => post<Status>('/api/player/play', { position: pos }),
    onSuccess: (st) => {
      queryClient.setQueryData(['status'], st)
      invalidateQueue()
    },
  })

  // Dupla confirmação para tocar (comportamento do cliente Windows).
  const [confirmPlay, setConfirmPlay] = useState<{ item: QueueItem; step: 1 | 2 } | null>(null)

  const goToCurrent = async () => {
    setQuery('')
    const { index } = await api<{ pos: number; index: number }>('/api/playlist/current')
    if (index >= 0) {
      virtualizer.scrollToIndex(index, { align: 'center' })
    }
  }

  return (
    <section className={`${styles.panel} ${styles.accentRed}`}>
      <header className={styles.panelHeader}>
        <div>
          <span className={`${styles.kicker} ${styles.kickerRed}`}>⚡ No ar · Programação</span>
          <h2 className={styles.panelTitle}>Fila de Hoje</h2>
        </div>
        <span className={styles.count}>{total.toLocaleString('pt-BR')} faixas</span>
      </header>

      <div className={styles.searchRow}>
        <input
          className={styles.search}
          placeholder="Pesquisar na fila..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button className={styles.buttonAccent} onClick={goToCurrent}>
          » Música atual
        </button>
      </div>

      <div ref={parentRef} className={styles.list}>
        <div style={{ height: virtualizer.getTotalSize(), position: 'relative' }}>
          {virtualizer.getVirtualItems().map((row) => {
            const item = itemAt(row.index)
            const isCurrent = item !== undefined && item.pos === currentPos
            const isSelected = item !== undefined && selected?.pos === item.pos
            return (
              <div
                key={row.key}
                className={`${styles.row} ${isCurrent ? styles.rowCurrent : ''} ${isSelected ? styles.rowSelected : ''}`}
                style={{ transform: `translateY(${row.start}px)`, height: ROW_HEIGHT }}
                onClick={() => item && onSelect(item)}
              >
                {item ? (
                  <>
                    <span className={styles.timeCol}>
                      <strong>{formatClock(item.nextPresentation)}</strong>
                      <small>#{item.pos}</small>
                    </span>
                    <span className={styles.avatar}>
                      {(item.artist || item.title).charAt(0).toUpperCase()}
                    </span>
                    <span className={styles.trackText}>
                      <strong>{item.title}</strong>
                      <small>
                        {item.artist || '—'} · {formatDuration(item.time)}
                      </small>
                    </span>
                    <GenreChip genre={item.genre} />
                    <span className={styles.rowActions}>
                      <button
                        title="Mover para cima"
                        disabled={item.pos === 0 || move.isPending}
                        onClick={(e) => {
                          e.stopPropagation()
                          move.mutate({ from: item.pos, to: item.pos - 1 })
                        }}
                      >
                        ▲
                      </button>
                      <button
                        title="Mover para baixo"
                        disabled={item.pos >= total - 1 || move.isPending}
                        onClick={(e) => {
                          e.stopPropagation()
                          move.mutate({ from: item.pos, to: item.pos + 1 })
                        }}
                      >
                        ▼
                      </button>
                      <button
                        title="Tocar agora"
                        className={styles.playAction}
                        onClick={(e) => {
                          e.stopPropagation()
                          setConfirmPlay({ item, step: 1 })
                        }}
                      >
                        ▶
                      </button>
                      <button
                        title="Remover da fila"
                        className={styles.removeAction}
                        disabled={remove.isPending}
                        onClick={(e) => {
                          e.stopPropagation()
                          remove.mutate(item.pos)
                        }}
                      >
                        ✕
                      </button>
                    </span>
                  </>
                ) : (
                  <span className={styles.placeholder}>Carregando…</span>
                )}
              </div>
            )
          })}
        </div>
      </div>

      <ConfirmDialog
        open={confirmPlay?.step === 1}
        title="Tem certeza?"
        confirmLabel="Sim, tocar"
        onCancel={() => setConfirmPlay(null)}
        onConfirm={() => setConfirmPlay((c) => (c ? { ...c, step: 2 } : c))}
      >
        Você tem certeza de que quer tocar <strong>{confirmPlay?.item.title}</strong>?
      </ConfirmDialog>
      <ConfirmDialog
        open={confirmPlay?.step === 2}
        title="TEM CERTEZA???"
        confirmLabel="Tocar agora"
        onCancel={() => setConfirmPlay(null)}
        onConfirm={() => {
          if (confirmPlay) play.mutate(confirmPlay.item.pos)
          setConfirmPlay(null)
        }}
      >
        Você vai parar de tocar <strong>{status?.current?.title ?? 'a música atual'}</strong>, tem
        CERTEZA?
      </ConfirmDialog>
    </section>
  )
}
