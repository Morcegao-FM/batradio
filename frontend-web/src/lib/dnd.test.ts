import { describe, expect, it } from 'vitest'
import { moveTargetFor } from './dnd'

describe('moveTargetFor', () => {
  it('mover para baixo desconta a remoção da origem', () => {
    // faixa #2 solta antes da #10 → depois de remover a #2, insere na 9
    expect(moveTargetFor(2, 10)).toBe(9)
  })

  it('mover para cima usa o índice de inserção direto', () => {
    // faixa #10 solta antes da #2 → vai para a posição 2
    expect(moveTargetFor(10, 2)).toBe(2)
  })

  it('soltar sobre si mesma é no-op', () => {
    // antes de si mesma (insertAt = from) e depois de si mesma (insertAt = from+1)
    expect(moveTargetFor(5, 5)).toBeNull()
    expect(moveTargetFor(5, 6)).toBeNull()
  })

  it('mover para o fim da fila', () => {
    // fila de 8 itens: soltar abaixo da última (insertAt = 8) → posição 7
    expect(moveTargetFor(3, 8)).toBe(7)
  })

  it('mover para o topo', () => {
    expect(moveTargetFor(3, 0)).toBe(0)
  })

  it('vizinhos: descer uma posição', () => {
    // #3 solta depois da #4 (insertAt 5) → to 4
    expect(moveTargetFor(3, 5)).toBe(4)
  })
})
