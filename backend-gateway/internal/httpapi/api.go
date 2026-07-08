// Package httpapi expõe a API REST + SSE do gateway.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/config"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/library"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/node"
)

type Server struct {
	cfg    *config.Config
	node   *node.Client
	lib    *library.Index
	poller *Poller

	// injetáveis em teste
	now     func() time.Time
	randInt func(n int) int

	qmu         sync.RWMutex
	queue       []model.Song
	queueLoaded bool
}

func NewServer(cfg *config.Config, n *node.Client) *Server {
	return &Server{
		cfg:     cfg,
		node:    n,
		lib:     library.NewIndex(),
		poller:  NewPoller(n, 2*time.Second),
		now:     time.Now,
		randInt: rand.Intn,
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
	r.Route("/api", func(api chi.Router) {
		api.Use(auth.RequireSession(s.cfg.SessionSecret, s.cfg.AllowedEmails))
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

// writeNodeErr mapeia erros do cliente Node para respostas HTTP.
func writeNodeErr(w http.ResponseWriter, err error) {
	if errors.Is(err, node.ErrNodeUnavailable) {
		writeErr(w, http.StatusBadGateway, "node_unavailable", err.Error())
		return
	}
	writeErr(w, http.StatusInternalServerError, "internal", err.Error())
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
	writeJSON(w, http.StatusOK, st)
}

func (s *Server) handleLibrary(w http.ResponseWriter, r *http.Request) {
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
