package httpapi

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/node"
)

// Event é um evento SSE já serializado.
type Event struct {
	Name string
	Data []byte
}

// Poller consulta o status do Node uma vez por intervalo e distribui aos
// assinantes SSE. Também detecta mudanças externas na fila (playlistlength).
type Poller struct {
	node     *node.Client
	interval time.Duration

	mu      sync.Mutex
	last    model.Status
	hasLast bool
	lastLen int
	lastErr error
	version int64
	subs    map[chan Event]struct{}
}

func NewPoller(n *node.Client, interval time.Duration) *Poller {
	return &Poller{
		node:     n,
		interval: interval,
		lastLen:  -1,
		subs:     make(map[chan Event]struct{}),
	}
}

func (p *Poller) Run(ctx context.Context) {
	t := time.NewTicker(p.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			p.Poll(ctx)
		}
	}
}

// Poll faz uma consulta única ao Node e publica os eventos resultantes.
func (p *Poller) Poll(ctx context.Context) {
	st, err := p.node.GetStatus(ctx)
	p.mu.Lock()
	if err != nil {
		p.lastErr = err
		p.hasLast = false
		p.mu.Unlock()
		return
	}
	p.lastErr = nil
	bumped := false
	if p.lastLen >= 0 && st.PlaylistLength != p.lastLen {
		p.version++
		bumped = true
	}
	p.lastLen = st.PlaylistLength
	p.last = st
	p.hasLast = true
	version := p.version
	p.mu.Unlock()

	if data, err := json.Marshal(st); err == nil {
		p.broadcast(Event{Name: "status", Data: data})
	}
	if bumped {
		p.broadcastVersion(version)
	}
}

// BumpPlaylist sinaliza que a fila mudou por uma mutação local.
func (p *Poller) BumpPlaylist() {
	p.mu.Lock()
	p.version++
	version := p.version
	p.mu.Unlock()
	p.broadcastVersion(version)
}

// NoteQueueLen registra o tamanho atual da fila para que o próximo Poll não
// interprete uma mutação local como mudança externa.
func (p *Poller) NoteQueueLen(n int) {
	p.mu.Lock()
	p.lastLen = n
	p.mu.Unlock()
}

func (p *Poller) broadcastVersion(v int64) {
	data, _ := json.Marshal(map[string]int64{"version": v})
	p.broadcast(Event{Name: "playlist", Data: data})
}

func (p *Poller) Version() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.version
}

func (p *Poller) Last() (model.Status, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.last, p.hasLast
}

func (p *Poller) LastError() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lastErr
}

func (p *Poller) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 16)
	p.mu.Lock()
	p.subs[ch] = struct{}{}
	p.mu.Unlock()
	return ch, func() {
		p.mu.Lock()
		delete(p.subs, ch)
		p.mu.Unlock()
	}
}

func (p *Poller) broadcast(e Event) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for ch := range p.subs {
		select {
		case ch <- e:
		default: // assinante lento: descarta em vez de travar o poller
		}
	}
}
