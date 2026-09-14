import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { afterEach, describe, expect, it } from 'vitest'
import NowPlaying from './NowPlaying'

function montar() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } })
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <NowPlaying />
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

afterEach(cleanup)

const abrir = () => screen.getByRole('button', { name: /adicionar música/i })
const acervo = () => screen.queryByRole('heading', { name: /buscar músicas/i })

describe('Tocando Agora', () => {
  it('mostra só a fila; o acervo entra sob demanda', () => {
    montar()
    expect(screen.getByRole('heading', { name: /fila de hoje/i })).toBeTruthy()
    expect(acervo()).toBeNull()

    fireEvent.click(abrir())
    expect(acervo()).toBeTruthy()
  })

  it('Esc fecha o acervo', () => {
    montar()
    fireEvent.click(abrir())
    expect(acervo()).toBeTruthy()

    fireEvent.keyDown(window, { key: 'Escape' })
    expect(acervo()).toBeNull()
  })

  it('o botão de fechar também fecha', () => {
    montar()
    fireEvent.click(abrir())
    fireEvent.click(screen.getByRole('button', { name: /fechar acervo/i }))
    expect(acervo()).toBeNull()
  })
})
