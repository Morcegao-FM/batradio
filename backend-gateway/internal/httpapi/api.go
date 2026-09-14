// Package httpapi expõe a API REST + SSE do gateway.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/catalogo"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/config"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/library"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/node"
)

// catalogoLookup é o que o gateway usa do catálogo do website. Interface (e não
// o tipo concreto) para o teste injetar um fake.
type catalogoLookup interface {
	Lookup(ctx context.Context, arquivos []string) map[string]catalogo.Info
}

type Server struct {
	cfg      *config.Config
	node     *node.Client
	lib      *library.Index
	poller   *Poller
	catalogo catalogoLookup

	// injetáveis em teste
	now     func() time.Time
	randInt func(n int) int

	qmu         sync.RWMutex
	queue       []model.Song
	queueLoaded bool
}

func NewServer(cfg *config.Config, n *node.Client, cat catalogoLookup) *Server {
	return &Server{
		cfg:      cfg,
		node:     n,
		lib:      library.NewIndex(),
		poller:   NewPoller(n, 2*time.Second),
		catalogo: cat,
		now:      time.Now,
		randInt:  rand.Intn,
	}
}

// Poller expõe o poller para o main iniciar o loop.
func (s *Server) Poller() *Poller { return s.poller }

// Library expõe o índice para o refresh inicial no main.
func (s *Server) Library() *library.Index { return s.lib }

// Node expõe o cliente do Node para o main (refresh inicial).
func (s *Server) Node() *node.Client { return s.node }

// Routes monta as rotas de auth (se fornecidas) e a API protegida.
func (s *Server) Routes(o *auth.OAuth) chi.Router {
	r := chi.NewRouter()
	if o != nil {
		o.Register(r)
	}
	// Público de propósito, e deliberadamente vazio: é só o sinal de vida que o
	// `docker compose up -d --wait` e o rollback do deploy leem. Não expõe
	// versão, host do Node nem estado do acervo — isso seria superfície de graça
	// para quem estiver sondando.
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(api chi.Router) {
		api.Use(auth.RequireSession(s.cfg.SessionSecret, s.cfg.AllowedEmails))
		// Depois do RequireSession de propósito: a auditoria precisa do e-mail
		// já no contexto.
		api.Use(auditoria(log.Printf))
		api.Get("/me", s.handleMe)
		api.Get("/status", s.handleStatus)
		api.Get("/events", s.handleEvents)
		api.Get("/library", s.handleLibrary)
		api.Post("/library/refresh", s.handleLibraryRefresh)
		api.Get("/server/info", s.handleServerInfo)

		api.Get("/playlist", s.handleGetPlaylist)
		api.Get("/playlist/current", s.handlePlaylistCurrent)
		api.Post("/playlist/add", s.handleAdd)
		api.Post("/playlist/add-many", s.handleAddMany)
		api.Post("/playlist/move", s.handleMove)
		api.Post("/playlist/remove", s.handleRemove)

		api.Post("/player/toggle", s.playerAction((*node.Client).PlayOrPause))
		api.Post("/player/shuffle", s.playerAction((*node.Client).Shuffle))
		api.Post("/player/repeat", s.playerAction((*node.Client).Repeat))
		api.Post("/player/crossfade", s.playerAction((*node.Client).Crossfade))
		api.Post("/player/play", s.handlePlay)

		api.Get("/playlists", s.handleListPlaylists)
		api.Post("/playlists", s.handleSavePlaylist)
		api.Post("/playlists/{name}/load", s.handleLoadPlaylist)
		api.Delete("/playlists/{name}", s.handleDeletePlaylist)
	})
	return r
}

// ---- helpers ----

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}

// writeNodeErr mapeia erros do cliente Node para respostas HTTP. "Busy" é
// transitório: responde 503 rápido para o frontend mostrar loading e tentar
// de novo, em vez de segurar a requisição.
func writeNodeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, node.ErrNodeBusy):
		w.Header().Set("Retry-After", "3")
		writeErr(w, http.StatusServiceUnavailable, "node_busy", "Servidor da rádio ocupado, tente novamente em instantes.")
	case errors.Is(err, node.ErrNodeUnavailable):
		writeErr(w, http.StatusBadGateway, "node_unavailable", err.Error())
	default:
		writeErr(w, http.StatusInternalServerError, "internal", err.Error())
	}
}

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_request", "JSON inválido: "+err.Error())
		return false
	}
	return true
}

func pageParams(r *http.Request) (q string, offset, limit int) {
	q = r.URL.Query().Get("q")
	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	return q, offset, limit
}

// ---- fila em cache ----

func (s *Server) setQueue(songs []model.Song, bump bool) {
	s.qmu.Lock()
	s.queue = songs
	s.queueLoaded = true
	s.qmu.Unlock()
	s.poller.NoteQueueLen(len(songs))
	if bump {
		s.poller.BumpPlaylist()
	}
}

// currentQueue retorna a fila em cache, buscando do Node na primeira vez.
func (s *Server) currentQueue(ctx context.Context) ([]model.Song, error) {
	s.qmu.RLock()
	if s.queueLoaded {
		q := s.queue
		s.qmu.RUnlock()
		return q, nil
	}
	s.qmu.RUnlock()
	songs, err := s.node.GetPlaylist(ctx, true)
	if err != nil {
		return nil, err
	}
	s.setQueue(songs, false)
	return songs, nil
}

// currentStatus prefere o cache do poller e cai para consulta direta.
func (s *Server) currentStatus(ctx context.Context) (model.Status, error) {
	st, err := s.node.GetStatus(ctx)
	if err == nil {
		return st, nil
	}
	if last, ok := s.poller.Last(); ok {
		return last, nil
	}
	return model.Status{}, err
}

// ---- handlers básicos ----

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"email": auth.EmailFromContext(r.Context())})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	st, err := s.currentStatus(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	s.enriquecerStatus(r.Context(), &st)
	writeJSON(w, http.StatusOK, st)
}

// enriquecer sobrepõe os dados de exibição do catálogo nas faixas dadas,
// in-place. Silencioso por natureza: catálogo desligado ou fora do ar deixa as
// faixas exatamente como vieram do MPD.
func (s *Server) enriquecer(ctx context.Context, musicas []model.Song) {
	if s.catalogo == nil || len(musicas) == 0 {
		return
	}

	arquivos := make([]string, 0, len(musicas))
	for _, m := range musicas {
		if m.File != "" {
			arquivos = append(arquivos, m.File)
		}
	}

	info := s.catalogo.Lookup(ctx, arquivos)
	for i := range musicas {
		dados, ok := info[musicas[i].File]
		if !ok {
			continue
		}
		musicas[i].ImageURL = dados.ImagemURL
		musicas[i].DisplayName = dados.NomeExibicao
		musicas[i].Year = dados.Ano
		musicas[i].Kind = dados.Tipo
	}
}

// enriquecerStatus completa só a faixa atual e a próxima. Nunca a fila inteira:
// o painel chama /api/status a cada 2 segundos.
func (s *Server) enriquecerStatus(ctx context.Context, st *model.Status) {
	var visiveis []model.Song
	if st.Current != nil {
		visiveis = append(visiveis, *st.Current)
	}
	if st.Next != nil {
		visiveis = append(visiveis, *st.Next)
	}
	if len(visiveis) == 0 {
		return
	}

	s.enriquecer(ctx, visiveis)

	i := 0
	if st.Current != nil {
		st.Current = &visiveis[i]
		i++
	}
	if st.Next != nil {
		st.Next = &visiveis[i]
	}
}

func (s *Server) handleLibrary(w http.ResponseWriter, r *http.Request) {
	// Antes do primeiro refresh o índice está vazio: um 200 com total 0 seria
	// cacheado pelo frontend como sucesso. Responde "ocupado" para ele
	// continuar tentando até o acervo ficar pronto.
	if _, refreshedAt := s.lib.Stats(); refreshedAt.IsZero() {
		w.Header().Set("Retry-After", "3")
		writeErr(w, http.StatusServiceUnavailable, "node_busy", "Acervo ainda sendo indexado, tente novamente em instantes.")
		return
	}
	q, offset, limit := pageParams(r)
	items, total := s.lib.Search(q, offset, limit)
	if limit <= 0 {
		limit = 50
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": items, "total": total, "offset": offset, "limit": limit,
	})
}

func (s *Server) handleLibraryRefresh(w http.ResponseWriter, r *http.Request) {
	songs, err := s.node.GetFiles(r.Context(), true)
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	s.lib.Set(songs, s.now())
	count, _ := s.lib.Stats()
	writeJSON(w, http.StatusOK, map[string]any{"count": count})
}

func (s *Server) handleServerInfo(w http.ResponseWriter, r *http.Request) {
	count, refreshedAt := s.lib.Stats()
	_, hasStatus := s.poller.Last()
	nodeOk := hasStatus && s.poller.LastError() == nil
	host := s.cfg.NodeBackendURL
	if u, err := url.Parse(s.cfg.NodeBackendURL); err == nil && u.Host != "" {
		host = u.Host
	}
	var refreshed any
	if !refreshedAt.IsZero() {
		refreshed = refreshedAt
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"nodeOk":             nodeOk,
		"nodeHost":           host,
		"libraryCount":       count,
		"libraryRefreshedAt": refreshed,
		"streamUrl":          s.cfg.StreamURL,
	})
}
