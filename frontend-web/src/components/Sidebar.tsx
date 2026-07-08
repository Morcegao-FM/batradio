import { NavLink } from 'react-router-dom'
import logo from '../assets/logo.png'
import { useServerInfo } from '../hooks/useStatus'
import styles from './Sidebar.module.css'

const NAV = [
  { to: '/', label: 'Tocando Agora', icon: '▶' },
  { to: '/playlists', label: 'Playlist e Arquivos', icon: '≡' },
  { to: '/config', label: 'Configurações', icon: '⚙' },
]

export default function Sidebar() {
  const { info } = useServerInfo()
  return (
    <aside className={styles.sidebar}>
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
            <span aria-hidden>{icon}</span>
            {label}
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
