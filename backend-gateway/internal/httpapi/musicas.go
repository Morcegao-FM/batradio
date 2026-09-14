package httpapi

import (
	"context"
	"net/http"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/catalogo"
)

// catalogoEditavel é a parte de escrita do catálogo. Separada de catalogoLookup
// porque a leitura é opcional e silenciosa, e a escrita não é.
type catalogoEditavel interface {
	Editar(ctx context.Context, arquivo, editadoPor string, e catalogo.Edicao) error
}

// handleEditarMusica corrige a faixa no catálogo do website.
//
// O autor vem SEMPRE da sessão, nunca do corpo: aceitar `editadoPor` do cliente
// deixaria qualquer um dos dois e-mails da allowlist assinar uma edição com o
// nome do outro, e é justamente a auditoria que sustenta o website aceitar
// escrita por chave de serviço.
func (s *Server) handleEditarMusica(w http.ResponseWriter, r *http.Request) {
	var body struct {
		File         string `json:"file"`
		Artista      string `json:"artista"`
		Titulo       string `json:"titulo"`
		Ano          int    `json:"ano"`
		NomeExibicao string `json:"nomeExibicao"`
		ImagemURL    string `json:"imagemUrl"`
		Tipo         string `json:"tipo"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if body.File == "" {
		writeErr(w, http.StatusBadRequest, "arquivo_obrigatorio", "Informe o arquivo da faixa.")
		return
	}

	editor, _ := s.catalogo.(catalogoEditavel)
	if editor == nil {
		writeErr(w, http.StatusServiceUnavailable, "catalogo_indisponivel",
			"O catálogo do site não está configurado neste gateway.")
		return
	}

	err := editor.Editar(r.Context(), body.File, auth.EmailFromContext(r.Context()), catalogo.Edicao{
		Artista:      body.Artista,
		Titulo:       body.Titulo,
		Ano:          body.Ano,
		NomeExibicao: body.NomeExibicao,
		ImagemURL:    body.ImagemURL,
		Tipo:         body.Tipo,
	})
	if err != nil {
		writeErr(w, http.StatusBadGateway, "catalogo_erro", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
