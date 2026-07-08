import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import Layout from '../components/Layout'
import LibraryPanel from '../components/LibraryPanel'
import QueuePanel from '../components/QueuePanel'
import AddManyModal from '../components/AddManyModal'
import { post } from '../lib/api'
import type { QueueItem, Song } from '../lib/types'
import styles from './NowPlaying.module.css'

export default function NowPlaying() {
  const [selectedSong, setSelectedSong] = useState<Song>()
  const [selectedQueueItem, setSelectedQueueItem] = useState<QueueItem>()
  const [addManySong, setAddManySong] = useState<Song>()
  const [toast, setToast] = useState<string>()
  const queryClient = useQueryClient()

  const add = useMutation({
    mutationFn: ({ file, position }: { file: string; position: number }) =>
      post('/api/playlist/add', { files: [file], position }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['playlist'] }),
  })

  const handleAdd = (place: 'above' | 'below') => {
    if (!selectedSong || !selectedQueueItem) return
    const position = place === 'above' ? selectedQueueItem.pos : selectedQueueItem.pos + 1
    add.mutate({ file: selectedSong.file, position })
  }

  const showToast = (msg: string) => {
    setToast(msg)
    setTimeout(() => setToast(undefined), 4000)
  }

  return (
    <Layout title="Tocando Agora" subtitle="Controle da programação ao vivo">
      <div className={styles.columns}>
        <LibraryPanel
          selected={selectedSong}
          onSelect={setSelectedSong}
          canAdd={!!selectedSong && !!selectedQueueItem && !add.isPending}
          onAdd={handleAdd}
          onAddMany={setAddManySong}
        />
        <QueuePanel selected={selectedQueueItem} onSelect={setSelectedQueueItem} />
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
