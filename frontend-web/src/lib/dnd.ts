// Payload do arrastar-e-soltar entre o acervo e a fila.

export const DND_MIME = 'application/x-batradio'

export type DragPayload =
  | { type: 'library'; file: string; title: string }
  | { type: 'queue'; pos: number; title: string }

export function setDragPayload(e: React.DragEvent, payload: DragPayload) {
  e.dataTransfer.setData(DND_MIME, JSON.stringify(payload))
  e.dataTransfer.effectAllowed = payload.type === 'library' ? 'copy' : 'move'
}

export function hasDragPayload(e: React.DragEvent): boolean {
  return e.dataTransfer.types.includes(DND_MIME)
}

export function getDragPayload(e: React.DragEvent): DragPayload | null {
  const raw = e.dataTransfer.getData(DND_MIME)
  if (!raw) return null
  try {
    return JSON.parse(raw) as DragPayload
  } catch {
    return null
  }
}

// Destino do move do MPD ao soltar uma faixa da fila em `insertAt` (índice de
// inserção "antes de"). O MPD remove a faixa da origem antes de reinserir,
// então destinos abaixo da origem deslocam uma posição para cima. Retorna
// null quando o drop não muda nada (soltar sobre si mesma).
export function moveTargetFor(from: number, insertAt: number): number | null {
  const to = insertAt > from ? insertAt - 1 : insertAt
  return to === from ? null : to
}
