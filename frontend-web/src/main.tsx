import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { focusManager, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import '@fontsource/anton/400.css'
import '@fontsource/oswald/400.css'
import '@fontsource/oswald/600.css'
import '@fontsource/barlow/400.css'
import '@fontsource/barlow/500.css'
import '@fontsource/barlow/600.css'
import '@fontsource/space-mono/400.css'
import './styles/tokens.css'
import './styles/base.css'
import App from './App'
import { isTransientError } from './lib/api'

// 502/503 vêm do gateway quando o servidor da rádio está ocupado ou fora do
// ar: transitório — tenta de novo a cada 3s (a UI mostra o estado de
// carregamento em vez de travar). Outros erros: 1 retry e desiste.
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: (failureCount, error) =>
        isTransientError(error) ? failureCount < 100 : failureCount < 1,
      retryDelay: 3000,
      // O padrão ('online') PAUSA fetches/retries quando navigator.onLine é
      // false — heurística que falha fácil (VPN, bridges, Linux) e deixa a
      // tela congelada sem erro. A API é o próprio host: sempre tentar.
      networkMode: 'always',
    },
    mutations: {
      networkMode: 'always',
    },
  },
})

// Painel de rádio costuma ficar num monitor secundário, sem foco: o TanStack
// pausa retries de abas desfocadas, o que congelaria o "tentando novamente".
focusManager.setFocused(true)

if (import.meta.env.DEV) {
  // inspeção no console durante desenvolvimento
  ;(window as unknown as Record<string, unknown>).__queryClient = queryClient
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </StrictMode>,
)
