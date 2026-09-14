import { useEffect, useState } from 'react'
import { NavLink } from 'react-router-dom'
import logo from '../assets/logo-horizontal.png'
import { useServerInfo } from '../hooks/useStatus'
import styles from './Sidebar.module.css'

const CHAVE_COLAPSO = 'batradio:menu-colapsado'

const NAV = [
  { to: '/', label: 'Tocando Agora', icon: '▶' },
  { to: '/playlists', label: 'Playlist e Arquivos', icon: '≡' },
  { to: '/config', label: 'Configurações', icon: '⚙' },
]

export default function Sidebar() {
  const { info } = useServerInfo()
  // Lembra a escolha entre sessões: quem trabalha com o menu fechado não quer
  // reabrir a cada visita. localStorage pode estourar (aba anônima), daí o try.
  const [colapsado, setColapsado] = useState(() => {
    try {
      return localStorage.getItem(CHAVE_COLAPSO) === '1'
    } catch {
      return false
    }
  })

  useEffect(() => {
    try {
      localStorage.setItem(CHAVE_COLAPSO, colapsado ? '1' : '0')
    } catch {
      // sem persistência, o menu só não lembra — não é motivo para quebrar
    }
  }, [colapsado])

  return (
    <aside className={`${styles.sidebar} ${colapsado ? styles.colapsado : ''}`}>
      <button
        type="button"
        className={styles.toggle}
        onClick={() => setColapsado((c) => !c)}
        aria-expanded={!colapsado}
        aria-label={colapsado ? 'Abrir menu' : 'Recolher menu'}
        title={colapsado ? 'Abrir menu' : 'Recolher menu'}
      >
        {colapsado ? '»' : '«'}
      </button>

      <img src={logo} alt="Morcegão FM" className={styles.logo} />

      <div className={styles.liveCard}>
        <span className={`${styles.liveBadge} ${info?.nodeOk ? styles.liveOk : styles.liveDown}`}>
          <span className={styles.liveDot} />
          AO VIVO
        </span>
        <span className={styles.server}>
          <span className={styles.serverLabel}>Servidor</span>
          {info?.nodeHost ?? '—'}
        </span>
      </div>

      <nav className={styles.nav}>
        <span className={styles.navLabel}>Navegação</span>
        {NAV.map(({ to, label, icon }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) => `${styles.navItem} ${isActive ? styles.active : ''}`}
          >
            <span className={styles.navIcon} aria-hidden>
              {icon}
            </span>
            <span className={styles.navTexto}>{label}</span>
          </NavLink>
        ))}
      </nav>

      <footer className={styles.footer}>
        <span>BATRADIO · V2.0</span>
        <span>Só rock clássico e blues.</span>
      </footer>
    </aside>
  )
}
