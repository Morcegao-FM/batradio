import type { ReactNode } from 'react'
import Sidebar from './Sidebar'
import HeaderBar from './HeaderBar'
import PlayerBar from './PlayerBar'
import { useRadioEvents } from '../hooks/useRadioEvents'
import { useStatus } from '../hooks/useStatus'
import { ApiError } from '../lib/api'
import styles from './Layout.module.css'

export default function Layout({
  title,
  subtitle,
  children,
}: {
  title: string
  subtitle: string
  children: ReactNode
}) {
  const { connected } = useRadioEvents()
  const { error } = useStatus()
  const nodeDown = error instanceof ApiError && error.status === 502

  return (
    <div className={styles.app}>
      <Sidebar />
      <div className={styles.main}>
        <HeaderBar title={title} subtitle={subtitle} />
        {(nodeDown || !connected) && (
          <div className={styles.banner} role="alert">
            {nodeDown
              ? 'Servidor da rádio inacessível. Verifique se o servidor de streaming está ativo.'
              : 'Reconectando ao gateway…'}
          </div>
        )}
        <main className={styles.content}>{children}</main>
        <PlayerBar />
      </div>
    </div>
  )
}
