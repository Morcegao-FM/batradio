import { useSearchParams } from 'react-router-dom'
import logo from '../assets/logo.png'
import styles from './Login.module.css'

export default function Login() {
  const [params] = useSearchParams()
  const denied = params.get('error') === 'email_not_allowed'
  const email = params.get('email')

  return (
    <main className={styles.page}>
      <div className={styles.card}>
        <img src={logo} alt="Morcegão FM" className={styles.logo} />
        <h1 className={styles.title}>BATRADIO</h1>
        <p className={styles.subtitle}>Controle da programação ao vivo</p>

        {denied && (
          <div className={styles.denied} role="alert">
            <strong>Acesso negado.</strong>
            <p>
              {email ? (
                <>
                  A conta <code>{email}</code> não tem permissão.
                </>
              ) : (
                'Esta conta Google não tem permissão.'
              )}{' '}
              Peça acesso ao administrador da rádio.
            </p>
          </div>
        )}

        <a className={styles.googleButton} href="/auth/login">
          Entrar com Google
        </a>
        <p className={styles.footer}>Morcegão FM · Só rock clássico e blues.</p>
      </div>
    </main>
  )
}
