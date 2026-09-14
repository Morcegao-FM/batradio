import { useRef, useState } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { api } from '../lib/api'
import type { Page, Song } from '../lib/types'
import { formatDuration } from '../lib/format'
import { useDebounced, usePagedList, type VisibleRange } from '../hooks/usePagedList'
import { setDragPayload } from '../lib/dnd'
import GenreChip from './GenreChip'
import ListState from './ListState'
import styles from './Panel.module.css'

const ROW_HEIGHT = 56

export default function LibraryPanel({
  selected,
  onSelect,
  canAdd,
  referenciaFila,
  onAdd,
  onAddMany,
}: {
  selected?: Song
  onSelect: (song: Song) => void
  canAdd: boolean
  referenciaFila: boolean
  onAdd: (place: 'above' | 'below') => void
  onAddMany: (song: Song) => void
}) {
  const [query, setQuery] = useState('')
  const q = useDebounced(query)
  const [range, setRange] = useState<VisibleRange>({ start: 0, end: 40 })

  const { total, itemAt, isFetching, error } = usePagedList<Song>(
    ['library', q],
    (offset, limit) =>
      api<Page<Song>>(
        `/api/library?q=${encodeURIComponent(q)}&offset=${offset}&limit=${limit}`,
      ),
    range,
  )

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

  return (
    <section className={`${styles.panel} ${styles.accentBlue}`}>
      <header className={styles.panelHeader}>
        <div>
          <span className={`${styles.kicker} ${styles.kickerBlue}`}>⚡ Acervo</span>
          <h2 className={styles.panelTitle}>Buscar Músicas</h2>
        </div>
        <span className={styles.count}>
          {isFetching ? '…' : `${total.toLocaleString('pt-BR')} faixas`}
        </span>
      </header>

      <input
        className={styles.search}
        placeholder="Banda, música, álbum..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
      />

      <div className={styles.actionsRow}>
        <button className={styles.buttonNeutral} disabled={!canAdd} onClick={() => onAdd('above')}>
          ↑ Adicionar acima
        </button>
        <button className={styles.buttonBrand} disabled={!canAdd} onClick={() => onAdd('below')}>
          ↓ Adicionar abaixo
        </button>
        <span className={styles.hint}>
          {!selected
            ? 'escolha uma faixa na lista abaixo'
            : referenciaFila
              ? 'da faixa selecionada na fila'
              : 'da música que está tocando'}
        </span>
      </div>

      <div ref={parentRef} className={styles.list}>
        {total === 0 && (
          <ListState error={error} isFetching={isFetching} emptyMessage="Nenhuma faixa encontrada" />
        )}
        <div style={{ height: virtualizer.getTotalSize(), position: 'relative' }}>
          {virtualizer.getVirtualItems().map((row) => {
            const song = itemAt(row.index)
            return (
              <div
                key={row.key}
                className={`${styles.row} ${selected?.file === song?.file ? styles.rowSelected : ''}`}
                style={{ transform: `translateY(${row.start}px)`, height: ROW_HEIGHT }}
                onClick={() => song && onSelect(song)}
                draggable={!!song}
                onDragStart={(e) => {
                  if (!song) return
                  onSelect(song)
                  setDragPayload(e, { type: 'library', file: song.file, title: song.title })
                }}
                title={song ? 'Arraste para a fila para agendar' : undefined}
              >
                {song ? (
                  <>
                    <span className={styles.avatar}>
                      {(song.artist || song.title).charAt(0).toUpperCase()}
                    </span>
                    <span className={styles.trackText}>
                      <strong>{song.title}</strong>
                      <small>
                        {song.artist || '—'} · {formatDuration(song.time)}
                      </small>
                    </span>
                    <GenreChip genre={song.genre} />
                    <button
                      className={styles.moreButton}
                      title="Inserir várias vezes"
                      onClick={(e) => {
                        e.stopPropagation()
                        onAddMany(song)
                      }}
                    >
                      ⋯
                    </button>
                  </>
                ) : (
                  <span className={styles.placeholder}>Carregando…</span>
                )}
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
