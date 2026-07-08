import { useMutation, useQueryClient } from '@tanstack/react-query'
import Layout from '../components/Layout'
import Toggle from '../components/Toggle'
import { post } from '../lib/api'
import type { Status } from '../lib/types'
import { useServerInfo, useStatus } from '../hooks/useStatus'
import { formatDate } from '../lib/format'
import styles from './Settings.module.css'

export default function Settings() {
  const { status } = useStatus()
  const { info } = useServerInfo()
  const queryClient = useQueryClient()

  // Um useMutation por ação; a ordem de chamada dos hooks é estável.
  const usePlayerAction = (action: string) =>
    // eslint-disable-next-line react-hooks/rules-of-hooks
    useMutation({
      mutationFn: () => post<Status>(`/api/player/${action}`),
      onSuccess: (st) => queryClient.setQueryData(['status'], st),
    })

  const toggle = usePlayerAction('toggle')
  const shuffle = usePlayerAction('shuffle')
  const repeat = usePlayerAction('repeat')
  const crossfade = usePlayerAction('crossfade')

  const refresh = useMutation({
    mutationFn: () => post<{ count: number }>('/api/library/refresh'),
    onSuccess: () => queryClient.invalidateQueries(),
  })

  return (
    <Layout title="Configurações" subtitle="Ajustes da rádio">
      <div className={styles.grid}>
        <section className={`${styles.card} ${styles.accentAmber}`}>
          <span className={`${styles.kicker} ${styles.kickerAmber}`}>⚡ Reprodução</span>
          <h2 className={styles.cardTitle}>Como a rádio toca</h2>

          <Toggle
            label="Tocar playlist"
            description="Reproduzir a programação continuamente"
            checked={status?.state === 'play'}
            disabled={!status || toggle.isPending}
            onChange={() => toggle.mutate()}
          />
          <Toggle
            label="Shuffle"
            description="Ordem aleatória das faixas"
            checked={!!status?.random}
            disabled={!status || shuffle.isPending}
            onChange={() => shuffle.mutate()}
          />
          <Toggle
            label="Repetir playlist"
            description="Recomeçar ao chegar no fim"
            checked={!!status?.repeat}
            disabled={!status || repeat.isPending}
            onChange={() => repeat.mutate()}
          />
          <Toggle
            label="Transição entre músicas"
            description="Crossfade suave entre faixas"
            checked={!!status?.crossfade}
            disabled={!status || crossfade.isPending}
            onChange={() => crossfade.mutate()}
          />
        </section>

        <section className={`${styles.card} ${styles.accentBlue}`}>
          <span className={`${styles.kicker} ${styles.kickerBlue}`}>⚡ Servidor</span>
          <h2 className={styles.cardTitle}>Conexão de streaming</h2>

          <dl className={styles.infoList}>
            <div>
              <dt>Servidor da rádio</dt>
              <dd>{info?.nodeHost ?? '—'}</dd>
            </div>
            <div>
              <dt>Conexão</dt>
              <dd className={info?.nodeOk ? styles.ok : styles.down}>
                {info ? (info.nodeOk ? '● Conectado' : '● Fora do ar') : '—'}
              </dd>
            </div>
            <div>
              <dt>Faixas no acervo</dt>
              <dd>{info ? info.libraryCount.toLocaleString('pt-BR') : '—'}</dd>
            </div>
            <div>
              <dt>Última atualização</dt>
              <dd>{formatDate(info?.libraryRefreshedAt)}</dd>
            </div>
          </dl>
          <p className={styles.note}>
            Endereço, porta e chave de API são configurados por variáveis de ambiente do gateway.
          </p>
          <button
            className={styles.refreshButton}
            disabled={refresh.isPending}
            onClick={() => refresh.mutate()}
          >
            ↻ {refresh.isPending ? 'Atualizando dados do servidor…' : 'Atualizar dados do servidor'}
          </button>
          {refresh.isSuccess && (
            <p className={styles.success}>Acervo atualizado: {refresh.data.count.toLocaleString('pt-BR')} faixas.</p>
          )}
        </section>
      </div>
    </Layout>
  )
}
