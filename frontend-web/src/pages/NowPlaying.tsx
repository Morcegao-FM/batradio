import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import Layout from '../components/Layout'
import LibraryPanel from '../components/LibraryPanel'
import QueuePanel from '../components/QueuePanel'
import AddManyModal from '../components/AddManyModal'
import { post } from '../lib/api'
import { useStatus } from '../hooks/useStatus'
import type { QueueItem, Song } from '../lib/types'
import styles from './NowPlaying.module.css'

export default function NowPlaying() {
  const [selectedSong, setSelectedSong] = useState<Song>()
  const [selectedQueueItem, setSelectedQueueItem] = useState<QueueItem>()
  const [addManySong, setAddManySong] = useState<Song>()
  const [toast, setToast] = useState<string>()
  const queryClient = useQueryClient()
  const { status } = useStatus()

  const add = useMutation({
    mutationFn: ({ file, position }: { file: string; position: number }) =>
      post('/api/playlist/add', { files: [file], position }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['playlist'] }),
  })

  // Sem linha da fila marcada, a referência é a música que está tocando: é o
  // único "acima/abaixo de quê" que faz sentido quando nada foi selecionado, e
  // deixa os botões utilizáveis já no primeiro clique do acervo.
  const posicaoReferencia = selectedQueueItem?.pos ?? status?.song

  const handleAdd = (place: 'above' | 'below') => {
    if (!selectedSong || posicaoReferencia === undefined) return
    const position = place === 'above' ? posicaoReferencia : posicaoReferencia + 1
    add.mutate({ file: selectedSong.file, position })
  }

  const showToast = (msg: string) => {
    setToast(msg)
    setTimeout(() => setToast(undefined), 4000)
  }

  return (
    <Layout title="Bat Radio" subtitle="Controle da programação ao vivo">
      <div className={styles.columns}>
        <LibraryPanel
          selected={selectedSong}
          onSelect={setSelectedSong}
          canAdd={!!selectedSong && posicaoReferencia !== undefined && !add.isPending}
          referenciaFila={!!selectedQueueItem}
          onAdd={handleAdd}
          onAddMany={setAddManySong}
        />
        <QueuePanel
          selected={selectedQueueItem}
          onSelect={setSelectedQueueItem}
          onAddFile={(file, position) => add.mutate({ file, position })}
        />
      </div>

      {addManySong && (
        <AddManyModal
          song={addManySong}
          onClose={() => setAddManySong(undefined)}
          onDone={(n) => {
            setAddManySong(undefined)
            showToast(`Inserida em ${n} ${n === 1 ? 'posição' : 'posições'} da programação.`)
          }}
        />
      )}
      {toast && <div className={styles.toast}>{toast}</div>}
    </Layout>
  )
}
