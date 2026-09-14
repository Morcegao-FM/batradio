import { cleanup, render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { afterEach, describe, expect, it, vi } from 'vitest'
import LibraryPanel from './LibraryPanel'
import type { Song } from '../lib/types'

const FAIXA: Song = {
  file: 'a.mp3', artist: 'AC/DC', title: 'Back in Black', time: 255, pos: 0, id: 1,
}

function montar(props: Partial<Parameters<typeof LibraryPanel>[0]> = {}) {
  // O painel busca o acervo pelo TanStack Query; sem retry o teste não espera
  // por rede que não existe aqui.
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  })
  return render(
    <QueryClientProvider client={client}>
    <LibraryPanel
      selected={undefined}
      onSelect={vi.fn()}
      canAdd={false}
      referenciaFila={false}
      onAdd={vi.fn()}
      onAddMany={vi.fn()}
      {...props}
    />
    </QueryClientProvider>,
  )
}

// O setup do projeto não instala os matchers do jest-dom nem o cleanup
// automático, então isto fica explícito aqui.
afterEach(cleanup)

const acima = () => screen.getByRole('button', { name: /adicionar acima/i }) as HTMLButtonElement
const abaixo = () => screen.getByRole('button', { name: /adicionar abaixo/i }) as HTMLButtonElement

describe('LibraryPanel: botões de adicionar', () => {
  it('sem faixa escolhida, ficam desabilitados e a dica diz o que falta', () => {
    montar()
    expect(acima().disabled).toBe(true)
    expect(abaixo().disabled).toBe(true)
    // O botão morto sem explicação era a reclamação: a dica passa a dizer o motivo.
    expect(screen.getByText(/escolha uma faixa na lista abaixo/i)).toBeTruthy()
  })

  it('com faixa escolhida, habilitam sem exigir seleção na fila', () => {
    montar({ selected: FAIXA, canAdd: true, referenciaFila: false })
    expect(acima().disabled).toBe(false)
    expect(abaixo().disabled).toBe(false)
    // Referência padrão: a música no ar.
    expect(screen.getByText(/da música que está tocando/i)).toBeTruthy()
  })

  it('com linha da fila marcada, a referência passa a ser ela', () => {
    montar({ selected: FAIXA, canAdd: true, referenciaFila: true })
    expect(screen.getByText(/da faixa selecionada na fila/i)).toBeTruthy()
  })
})
