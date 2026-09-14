// Permite usar React.act nos testes de hooks.
;(globalThis as Record<string, unknown>).IS_REACT_ACT_ENVIRONMENT = true

// jsdom não implementa EventSource, e o Layout abre um SSE ao montar. Stub
// mínimo: registra handlers e nunca emite nada — os testes de componente não
// dependem de evento do servidor.
class EventSourceFake {
  static readonly CONNECTING = 0
  static readonly OPEN = 1
  static readonly CLOSED = 2
  readyState = EventSourceFake.CONNECTING
  onopen: (() => void) | null = null
  onerror: (() => void) | null = null
  onmessage: (() => void) | null = null
  constructor(public url: string) {}
  addEventListener() {}
  removeEventListener() {}
  close() {
    this.readyState = EventSourceFake.CLOSED
  }
}
;(globalThis as Record<string, unknown>).EventSource = EventSourceFake
