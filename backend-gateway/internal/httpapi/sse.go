package httpapi

import (
	"fmt"
	"net/http"
	"time"
)

// handleEvents mantém uma conexão SSE aberta, repassando os eventos do poller.
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	fl, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, http.StatusInternalServerError, "internal", "streaming não suportado")
		return
	}
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	ch, cancel := s.poller.Subscribe()
	defer cancel()

	// estado atual imediatamente, se disponível
	if st, ok := s.poller.Last(); ok {
		writeSSE(w, Event{Name: "status", Data: mustJSON(st)})
		fl.Flush()
	}

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-heartbeat.C:
			fmt.Fprint(w, ": ping\n\n")
			fl.Flush()
		case e := <-ch:
			writeSSE(w, e)
			fl.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, e Event) {
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Name, e.Data)
}
