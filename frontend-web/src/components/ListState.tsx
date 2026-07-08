import { apiErrorStatus, isTransientError } from '../lib/api'
import styles from './Panel.module.css'

// Estado vazio das listas: diferencia "carregando", "servidor ocupado/fora
// do ar (tentando de novo a cada 3s)" e "nada encontrado".
export default function ListState({
  error,
  isFetching,
  emptyMessage,
}: {
  error?: unknown
  isFetching: boolean
  emptyMessage: string
}) {
  if (isTransientError(error)) {
    return (
      <div className={`${styles.emptyState} ${styles.emptyBusy}`} role="status">
        <span className={styles.spinner} aria-hidden />
        {apiErrorStatus(error) === 503
          ? 'Servidor da rádio ocupado — tentando novamente…'
          : 'Servidor da rádio inacessível — tentando novamente…'}
      </div>
    )
  }
  if (error) {
    return (
      <div className={styles.emptyState} role="alert">
        Falha ao carregar: {(error as Error).message}
      </div>
    )
  }
  if (isFetching) {
    return (
      <div className={styles.emptyState} role="status">
        <span className={styles.spinner} aria-hidden />
        Carregando…
      </div>
    )
  }
  return <div className={styles.emptyState}>{emptyMessage}</div>
}
