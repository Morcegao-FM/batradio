package httpapi

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

func playlistName(w http.ResponseWriter, r *http.Request) (string, bool) {
	raw := chi.URLParam(r, "name")
	name, err := url.PathUnescape(raw)
	if err != nil || name == "" || strings.ContainsAny(name, "/\\") {
		writeErr(w, http.StatusBadRequest, "bad_request", "nome de playlist inválido")
		return "", false
	}
	return name, true
}

func (s *Server) handleListPlaylists(w http.ResponseWriter, r *http.Request) {
	pls, err := s.node.GetPlaylists(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": pls})
}

func (s *Server) handleSavePlaylist(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	name := strings.TrimSpace(body.Name)
	if name == "" || strings.ContainsAny(name, "/\\") {
		writeErr(w, http.StatusBadRequest, "bad_request", "nome de playlist inválido")
		return
	}
	if err := s.node.SavePlaylist(r.Context(), name); err != nil {
		writeNodeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleLoadPlaylist(w http.ResponseWriter, r *http.Request) {
	name, ok := playlistName(w, r)
	if !ok {
		return
	}
	var body struct {
		PlayAtRandom bool `json:"playAtRandom"`
	}
	if r.ContentLength > 0 && !decodeBody(w, r, &body) {
		return
	}
	songs, err := s.node.LoadPlaylist(r.Context(), name)
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	s.setQueue(songs, true)

	if body.PlayAtRandom && len(songs) > 1 {
		// Comportamento do cliente Windows: começa de uma posição aleatória
		// (nunca a primeira, para não repetir sempre a abertura da playlist).
		pos := 1 + s.randInt(len(songs)-1)
		st, err := s.node.Play(r.Context(), pos)
		if err != nil {
			writeNodeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, st)
		return
	}
	st, err := s.currentStatus(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}

func (s *Server) handleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	name, ok := playlistName(w, r)
	if !ok {
		return
	}
	if err := s.node.RemovePlaylist(r.Context(), name); err != nil {
		writeNodeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
