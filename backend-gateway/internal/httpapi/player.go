package httpapi

import (
	"context"
	"net/http"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/node"
)

// playerAction adapta um método de toggle do cliente Node em um handler.
func (s *Server) playerAction(fn func(*node.Client, context.Context) (model.Status, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st, err := fn(s.node, r.Context())
		if err != nil {
			writeNodeErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, st)
	}
}

func (s *Server) handlePlay(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Position int `json:"position"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if body.Position < 0 {
		writeErr(w, http.StatusBadRequest, "bad_request", "position deve ser >= 0")
		return
	}
	st, err := s.node.Play(r.Context(), body.Position)
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}
