import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { put } from '../lib/api'
import type { Song } from '../lib/types'
import styles from './EditTrackModal.module.css'

// Os mesmos tipos que a Batcaverna usa, para as duas telas não divergirem.
const TIPOS = [
  { valor: 'musica', rotulo: 'Música' },
  { valor: 'vinheta', rotulo: 'Vinheta' },
  { valor: 'campanha', rotulo: 'Campanha' },
  { valor: 'abertura', rotulo: 'Abertura' },
  { valor: 'bg', rotulo: 'BG' },
]

export default function EditTrackModal({
  song,
  onClose,
}: {
  song: Song
  onClose: () => void
}) {
  const [artista, setArtista] = useState(song.artist ?? '')
  const [titulo, setTitulo] = useState(song.title ?? '')
  const [ano, setAno] = useState(song.year ? String(song.year) : '')
  const [nomeExibicao, setNomeExibicao] = useState(song.displayName ?? '')
  const [imagemUrl, setImagemUrl] = useState(song.imageUrl ?? '')
  const [tipo, setTipo] = useState(song.kind || 'musica')
  const queryClient = useQueryClient()

  const salvar = useMutation({
    mutationFn: () =>
      put<void>('/api/musicas', {
        file: song.file,
        artista,
        titulo,
        ano: ano ? Number(ano) : 0,
        nomeExibicao,
        imagemUrl,
        tipo,
      }),
    onSuccess: () => {
      // O gateway já invalidou o cache do catálogo; aqui é só refazer as telas.
      queryClient.invalidateQueries({ queryKey: ['status'] })
      queryClient.invalidateQueries({ queryKey: ['playlist'] })
      onClose()
    },
  })

  return (
    <div className={styles.backdrop} onClick={onClose}>
      <div
        className={styles.dialog}
        role="dialog"
        aria-label="Corrigir dados da faixa"
        onClick={(e) => e.stopPropagation()}
      >
        <span className={styles.kicker}>Catálogo</span>
        <h2 className={styles.title}>Corrigir a faixa</h2>
        <p className={styles.subtitle}>
          Vale para o site e para o app também. Depois de corrigir, a busca automática de
          capa para de mexer nesta faixa.
        </p>

        <div className={styles.arquivo} title={song.file}>
          {song.file}
        </div>

        <div className={styles.fields}>
          <label className={styles.field}>
            <span>Artista</span>
            <input value={artista} onChange={(e) => setArtista(e.target.value)} maxLength={255} />
          </label>

          <label className={styles.field}>
            <span>Título</span>
            <input value={titulo} onChange={(e) => setTitulo(e.target.value)} maxLength={255} />
          </label>

          <div className={styles.linha}>
            <label className={styles.field}>
              <span>Ano</span>
              <input
                type="number"
                inputMode="numeric"
                value={ano}
                onChange={(e) => setAno(e.target.value)}
                min={1900}
                max={2100}
              />
            </label>

            <label className={styles.field}>
              <span>Tipo</span>
              <select value={tipo} onChange={(e) => setTipo(e.target.value)}>
                {TIPOS.map((t) => (
                  <option key={t.valor} value={t.valor}>
                    {t.rotulo}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <label className={styles.field}>
            <span>Nome de exibição</span>
            <input
              value={nomeExibicao}
              onChange={(e) => setNomeExibicao(e.target.value)}
              maxLength={300}
              placeholder="Só para vinheta com nome de arquivo feio"
            />
          </label>

          <label className={styles.field}>
            <span>URL da capa</span>
            <input
              value={imagemUrl}
              onChange={(e) => setImagemUrl(e.target.value)}
              maxLength={500}
              placeholder="https://..."
            />
          </label>
        </div>

        {imagemUrl ? (
          <img className={styles.previa} src={imagemUrl} alt="" />
        ) : null}

        {salvar.isError ? (
          <p className={styles.error} role="alert">
            {salvar.error instanceof Error ? salvar.error.message : 'Não deu para salvar.'}
          </p>
        ) : null}

        <div className={styles.actions}>
          <button type="button" className={styles.secondary} onClick={onClose}>
            Cancelar
          </button>
          <button
            type="button"
            className={styles.primary}
            onClick={() => salvar.mutate()}
            disabled={salvar.isPending}
          >
            {salvar.isPending ? 'Salvando…' : 'Salvar'}
          </button>
        </div>
      </div>
    </div>
  )
}
