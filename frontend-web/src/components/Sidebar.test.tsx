import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import Sidebar from './Sidebar'

function montar() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } })
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <Sidebar />
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

const botao = () => screen.getByRole('button', { name: /menu/i })

afterEach(cleanup)
beforeEach(() => localStorage.clear())

describe('Sidebar colapsável', () => {
  it('começa aberta e alterna ao clicar', () => {
    const { container } = montar()
    const aside = container.querySelector('aside')!

    expect(botao().getAttribute('aria-expanded')).toBe('true')
    const abertaClasses = aside.className

    fireEvent.click(botao())
    expect(botao().getAttribute('aria-expanded')).toBe('false')
    expect(aside.className).not.toBe(abertaClasses)

    fireEvent.click(botao())
    expect(botao().getAttribute('aria-expanded')).toBe('true')
  })

  it('lembra a escolha entre sessões', () => {
    montar()
    fireEvent.click(botao())
    cleanup()

    montar()
    expect(botao().getAttribute('aria-expanded')).toBe('false')
  })

  it('mantém os destinos acessíveis mesmo recolhida', () => {
    // O rótulo some por CSS, não do DOM: leitor de tela continua anunciando.
    montar()
    fireEvent.click(botao())
    expect(screen.getByRole('link', { name: /tocando agora/i })).toBeTruthy()
    expect(screen.getByRole('link', { name: /playlist e arquivos/i })).toBeTruthy()
    expect(screen.getByRole('link', { name: /configurações/i })).toBeTruthy()
  })
})
