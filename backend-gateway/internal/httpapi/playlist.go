package httpapi

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/library"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/schedule"
)

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func (s *Server) handleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	queue, err := s.currentQueue(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	st, err := s.currentStatus(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	items := schedule.ComputeTimes(queue, st, s.now())

	q, offset, limit := pageParams(r)
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	filtered := items
	if q != "" {
		filtered = filtered[:0:0]
		for _, it := range items {
			if library.MatchesQuery(it.Song, q) {
				filtered = append(filtered, it)
			}
		}
	}
	total := len(filtered)
	end := offset + limit
	if offset > total {
		offset = total
	}
	if end > total {
		end = total
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items":      filtered[offset:end],
		"total":      total,
		"offset":     offset,
		"limit":      limit,
		"currentPos": st.Song,
		"version":    s.poller.Version(),
	})
}

func (s *Server) handlePlaylistCurrent(w http.ResponseWriter, r *http.Request) {
	queue, err := s.currentQueue(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	st, err := s.currentStatus(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	index := -1
	for i, song := range queue {
		if song.Pos == st.Song {
			index = i
			break
		}
	}
	writeJSON(w, http.StatusOK, map[string]int{"pos": st.Song, "index": index})
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Files    []string `json:"files"`
		Position int      `json:"position"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if len(body.Files) == 0 {
		writeErr(w, http.StatusBadRequest, "bad_request", "files não pode ser vazio")
		return
	}
	if body.Position < 0 {
		writeErr(w, http.StatusBadRequest, "bad_request", "position deve ser >= 0")
		return
	}
	songs, err := s.node.Add(r.Context(), body.Files, body.Position)
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	s.setQueue(songs, true)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleAddMany(w http.ResponseWriter, r *http.Request) {
	var body struct {
		File          string     `json:"file"`
		TimesPerDay   int        `json:"timesPerDay"`
		IntervalHours int        `json:"intervalHours"`
		From          *time.Time `json:"from"`
		To            *time.Time `json:"to"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	switch {
	case body.File == "":
		writeErr(w, http.StatusBadRequest, "bad_request", "file é obrigatório")
		return
	case body.TimesPerDay < 1 || body.TimesPerDay > 96:
		writeErr(w, http.StatusBadRequest, "bad_request", "timesPerDay deve estar entre 1 e 96")
		return
	case body.IntervalHours < 1 || body.IntervalHours > 24:
		writeErr(w, http.StatusBadRequest, "bad_request", "intervalHours deve estar entre 1 e 24")
		return
	}
	queue, err := s.currentQueue(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	st, err := s.currentStatus(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	items := schedule.ComputeTimes(queue, st, s.now())
	positions := schedule.PlanInsertions(items, body.TimesPerDay,
		time.Duration(body.IntervalHours)*time.Hour, body.From, body.To, s.now())

	// posições em ordem decrescente: inserir de trás pra frente preserva índices
	for _, pos := range positions {
		if _, err := s.node.Add(r.Context(), []string{body.File}, pos); err != nil {
			writeNodeErr(w, err)
			return
		}
	}
	songs, err := s.node.GetPlaylist(r.Context(), true)
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	s.setQueue(songs, true)
	if positions == nil {
		positions = []int{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"insertedAt": positions})
}

func (s *Server) handleMove(w http.ResponseWriter, r *http.Request) {
	var body struct {
		From int `json:"from"`
		To   int `json:"to"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if body.From < 0 || body.To < 0 {
		writeErr(w, http.StatusBadRequest, "bad_request", "posições devem ser >= 0")
		return
	}
	songs, err := s.node.Move(r.Context(), body.From, body.To)
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	s.setQueue(songs, true)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleRemove(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Positions []int `json:"positions"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if len(body.Positions) == 0 {
		writeErr(w, http.StatusBadRequest, "bad_request", "positions não pode ser vazio")
		return
	}
	// Remove da maior posição para a menor, uma por chamada: o endpoint /delete
	// do Node dispara os comandos de forma assíncrona, então mandar o lote todo
	// numa chamada não garante ordem.
	positions := append([]int(nil), body.Positions...)
	sort.Sort(sort.Reverse(sort.IntSlice(positions)))
	for _, pos := range positions {
		out, err := s.node.Delete(r.Context(), []int{pos})
		if err != nil {
			writeNodeErr(w, err)
			return
		}
		s.setQueue(out, false)
	}
	s.poller.BumpPlaylist()
	w.WriteHeader(http.StatusNoContent)
}
