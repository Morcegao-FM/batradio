import type { ReactNode } from 'react'
import Sidebar from './Sidebar'
import HeaderBar from './HeaderBar'
import PlayerBar from './PlayerBar'
import { useRadioEvents } from '../hooks/useRadioEvents'
import { useStatus } from '../hooks/useStatus'
import { apiErrorStatus } from '../lib/api'
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
  const nodeDown = apiErrorStatus(error) === 502
  const nodeBusy = apiErrorStatus(error) === 503

  return (
    <div className={styles.app}>
      <Sidebar />
      <div className={styles.main}>
        <HeaderBar title={title} subtitle={subtitle} />
        {(nodeDown || nodeBusy || !connected) && (
          <div className={styles.banner} role="alert">
            {nodeDown
              ? 'Servidor da rádio inacessível. Verifique se o servidor de streaming está ativo.'
              : nodeBusy
                ? 'Servidor da rádio ocupado — tentando novamente…'
                : 'Reconectando ao gateway…'}
          </div>
        )}
        <main className={styles.content}>{children}</main>
        <PlayerBar />
      </div>
    </div>
  )
}
