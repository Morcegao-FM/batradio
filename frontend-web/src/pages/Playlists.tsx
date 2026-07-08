import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import Layout from '../components/Layout'
import ConfirmDialog from '../components/ConfirmDialog'
import { api, post } from '../lib/api'
import type { PlaylistInfo, Status } from '../lib/types'
import { formatDate } from '../lib/format'
import styles from './Playlists.module.css'

const CARD_ACCENTS = ['accentRed', 'accentBlue', 'accentAmber'] as const

export default function Playlists() {
  const queryClient = useQueryClient()
  const { data } = useQuery({
    queryKey: ['playlists'],
    queryFn: () => api<{ items: PlaylistInfo[] }>('/api/playlists'),
  })
  const playlists = data?.items ?? []

  const [selectedName, setSelectedName] = useState('')
  const [confirmLoad, setConfirmLoad] = useState<{ name: string; playAtRandom: boolean } | null>(
    null,
  )
  const [confirmDelete, setConfirmDelete] = useState<string | null>(null)
  const [newName, setNewName] = useState('')
  const [feedback, setFeedback] = useState<string>()

  const flash = (msg: string) => {
    setFeedback(msg)
    setTimeout(() => setFeedback(undefined), 4000)
  }

  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: ['playlists'] })
    queryClient.invalidateQueries({ queryKey: ['playlist'] })
  }

  const load = useMutation({
    mutationFn: ({ name, playAtRandom }: { name: string; playAtRandom: boolean }) =>
      post<Status>(`/api/playlists/${encodeURIComponent(name)}/load`, { playAtRandom }),
    onSuccess: (st, vars) => {
      queryClient.setQueryData(['status'], st)
      invalidate()
      flash(
        vars.playAtRandom
          ? `Playlist "${vars.name}" carregada e tocando de posição aleatória.`
          : `Playlist "${vars.name}" carregada.`,
      )
    },
    onError: (e) => flash(`Falha ao carregar: ${(e as Error).message}`),
  })

  const save = useMutation({
    mutationFn: (name: string) => post('/api/playlists', { name }),
    onSuccess: (_, name) => {
      invalidate()
      setNewName('')
      flash(`Fila atual salva como "${name}".`)
    },
    onError: (e) => flash(`Falha ao salvar: ${(e as Error).message}`),
  })

  const remove = useMutation({
    mutationFn: (name: string) =>
      api(`/api/playlists/${encodeURIComponent(name)}`, { method: 'DELETE' }),
    onSuccess: (_, name) => {
      invalidate()
      flash(`Playlist "${name}" excluída.`)
    },
    onError: (e) => flash(`Falha ao excluir: ${(e as Error).message}`),
  })

  return (
    <Layout title="Playlist e Arquivos" subtitle="Biblioteca do servidor">
      <div className={styles.stack}>
        <section className={styles.loaderCard}>
          <span className={styles.loaderLabel}>Carregar playlist do servidor</span>
          <div className={styles.loaderRow}>
            <select
              value={selectedName}
              onChange={(e) => setSelectedName(e.target.value)}
              className={styles.select}
            >
              <option value="">Selecione uma playlist…</option>
              {playlists.map((p) => (
                <option key={p.name} value={p.name}>
                  {p.name}
                </option>
              ))}
            </select>
            <button
              className={styles.loadPlay}
              disabled={!selectedName || load.isPending}
              onClick={() => setConfirmLoad({ name: selectedName, playAtRandom: true })}
            >
              ▶ Carregar e tocar
            </button>
            <span className={styles.loaderHint}>
              Substitui a fila atual pela playlist e começa a tocar de uma posição aleatória.
            </span>
          </div>
        </section>

        <span className={styles.sectionLabel}>⚡ Playlists disponíveis</span>
        <div className={styles.grid}>
          {playlists.map((p, i) => (
            <article key={p.name} className={`${styles.card} ${styles[CARD_ACCENTS[i % 3]]}`}>
              <header className={styles.cardHeader}>
                <h3 className={styles.cardTitle}>{p.name}</h3>
                <button
                  className={styles.cardDelete}
                  title="Excluir playlist"
                  onClick={() => setConfirmDelete(p.name)}
                >
                  ✕
                </button>
              </header>
              <span className={styles.cardMeta}>{formatDate(p.lastModified)}</span>
              <button
                className={styles.cardLoad}
                disabled={load.isPending}
                onClick={() => setConfirmLoad({ name: p.name, playAtRandom: false })}
              >
                Carregar esta
              </button>
            </article>
          ))}
          <article className={`${styles.card} ${styles.saveCard}`}>
            <h3 className={styles.cardTitle}>Salvar fila atual</h3>
            <input
              placeholder="Nome da playlist"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
            />
            <button
              className={styles.cardLoad}
              disabled={!newName.trim() || save.isPending}
              onClick={() => save.mutate(newName.trim())}
            >
              Salvar
            </button>
          </article>
        </div>
      </div>

      <ConfirmDialog
        open={confirmLoad !== null}
        title="Substituir a fila atual?"
        confirmLabel={confirmLoad?.playAtRandom ? 'Carregar e tocar' : 'Carregar'}
        onCancel={() => setConfirmLoad(null)}
        onConfirm={() => {
          if (confirmLoad) load.mutate(confirmLoad)
          setConfirmLoad(null)
        }}
      >
        A fila atual será substituída pela playlist <strong>{confirmLoad?.name}</strong>
        {confirmLoad?.playAtRandom ? ', tocando a partir de uma posição aleatória.' : '.'}
      </ConfirmDialog>

      <ConfirmDialog
        open={confirmDelete !== null}
        title="Excluir playlist?"
        confirmLabel="Excluir"
        onCancel={() => setConfirmDelete(null)}
        onConfirm={() => {
          if (confirmDelete) remove.mutate(confirmDelete)
          setConfirmDelete(null)
        }}
      >
        A playlist <strong>{confirmDelete}</strong> será excluída do servidor. Esta ação não pode
        ser desfeita.
      </ConfirmDialog>

      {feedback && <div className={styles.toast}>{feedback}</div>}
    </Layout>
  )
}
