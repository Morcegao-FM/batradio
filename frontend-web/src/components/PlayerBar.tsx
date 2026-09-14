import { useEffect, useRef, useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { post } from '../lib/api'
import type { Status } from '../lib/types'
import { useServerInfo, useStatus } from '../hooks/useStatus'
import logo from '../assets/logo.png'
import styles from './PlayerBar.module.css'

function Equalizer({ playing }: { playing: boolean }) {
  return (
    <span className={`${styles.eq} ${playing ? styles.eqOn : ''}`} aria-hidden>
      <i />
      <i />
      <i />
      <i />
    </span>
  )
}

// Monitor local: toca o stream do Icecast no browser (não afeta a transmissão).
function StreamMonitor({ url }: { url: string }) {
  const audioRef = useRef<HTMLAudioElement>(null)
  const [monitoring, setMonitoring] = useState(false)
  const [volume, setVolume] = useState(0.8)

  useEffect(() => {
    const audio = audioRef.current
    if (!audio) return
    audio.volume = volume
    if (monitoring) {
      audio.play().catch(() => setMonitoring(false))
    } else {
      audio.pause()
    }
  }, [monitoring, volume])

  return (
    <div className={styles.monitor}>
      <audio ref={audioRef} src={monitoring ? url : undefined} preload="none" />
      <button
        className={`${styles.monitorButton} ${monitoring ? styles.monitorOn : ''}`}
        onClick={() => setMonitoring((m) => !m)}
        title={monitoring ? 'Parar de ouvir a rádio' : 'Ouvir a rádio (monitor local)'}
      >
        {monitoring ? '🔊' : '🔈'}
      </button>
      <input
        type="range"
        min={0}
        max={1}
        step={0.05}
        value={volume}
        onChange={(e) => setVolume(Number(e.target.value))}
        className={styles.volume}
        aria-label="Volume do monitor"
      />
    </div>
  )
}

// Capa da faixa, vinda do catálogo do site. Sem capa — ou com URL que falha —
// cai no logo: o painel nunca fica com buraco na tela por causa do catálogo.
function Capa({ src }: { src?: string }) {
  const [falhou, setFalhou] = useState(false)
  useEffect(() => setFalhou(false), [src])
  return (
    <img
      src={src && !falhou ? src : logo}
      alt=""
      className={styles.thumb}
      onError={() => setFalhou(true)}
    />
  )
}

export default function PlayerBar() {
  const { status } = useStatus()
  const { info } = useServerInfo()
  const queryClient = useQueryClient()
  const toggle = useMutation({
    mutationFn: () => post<Status>('/api/player/toggle'),
    onSuccess: (st) => queryClient.setQueryData(['status'], st),
  })

  const playing = status?.state === 'play'
  const current = status?.current

  return (
    <footer className={styles.player}>
      <button
        className={styles.playButton}
        onClick={() => toggle.mutate()}
        disabled={toggle.isPending || !status}
        title={playing ? 'Pausar' : 'Tocar'}
      >
        {playing ? '❚❚' : '▶'}
      </button>
      <Capa src={current?.imageUrl} />
      <div className={styles.trackInfo}>
        <span className={styles.nowLabel}>
          Tocando agora <Equalizer playing={!!playing} />
        </span>
        <span className={styles.track}>
          <strong>{current?.displayName ?? current?.title ?? '—'}</strong>
          {!current?.displayName && current?.artist ? (
            <span className={styles.artist}> — {current.artist}</span>
          ) : null}
        </span>
      </div>
      <div className={styles.right}>
        <span className={`${styles.liveBadge} ${info?.nodeOk ? '' : styles.liveOff}`}>
          ● AO VIVO
        </span>
        {info?.streamUrl ? <StreamMonitor url={info.streamUrl} /> : null}
      </div>
    </footer>
  )
}
