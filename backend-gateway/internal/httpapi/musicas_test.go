package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/catalogo"
)

type catalogoEditavelFake struct {
	catalogoFake
	arquivo    string
	editadoPor string
	edicao     catalogo.Edicao
	err        error
}

func (f *catalogoEditavelFake) Editar(_ context.Context, arquivo, editadoPor string, e catalogo.Edicao) error {
	f.arquivo, f.editadoPor, f.edicao = arquivo, editadoPor, e
	return f.err
}

func TestEditarMusicaUsaEmailDaSessao(t *testing.T) {
	fake := &catalogoEditavelFake{}
	s := newTestServer(t, "http://127.0.0.1:1")
	s.catalogo = fake

	req := httptest.NewRequest("PUT", "/api/musicas",
		strings.NewReader(`{"file":"a.mp3","titulo":"Back in Black","ano":1980}`))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(auth.ContextComEmail(req.Context(), "aguergolet@gmail.com"))
	rec := httptest.NewRecorder()

	s.handleEditarMusica(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("esperava 204, veio %d: %s", rec.Code, rec.Body.String())
	}
	// O e-mail vem da sessão, nunca do corpo: senão qualquer um assinaria a
	// edição com o nome de outra pessoa.
	if fake.editadoPor != "aguergolet@gmail.com" {
		t.Errorf("autor errado: %q", fake.editadoPor)
	}
	if fake.arquivo != "a.mp3" || fake.edicao.Titulo != "Back in Black" || fake.edicao.Ano != 1980 {
		t.Errorf("payload errado: %s / %+v", fake.arquivo, fake.edicao)
	}
}

func TestEditarMusicaIgnoraAutorForjadoNoCorpo(t *testing.T) {
	fake := &catalogoEditavelFake{}
	s := newTestServer(t, "http://127.0.0.1:1")
	s.catalogo = fake

	req := httptest.NewRequest("PUT", "/api/musicas",
		strings.NewReader(`{"file":"a.mp3","editadoPor":"outra-pessoa@gmail.com"}`))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(auth.ContextComEmail(req.Context(), "aguergolet@gmail.com"))

	s.handleEditarMusica(httptest.NewRecorder(), req)

	if fake.editadoPor != "aguergolet@gmail.com" {
		t.Errorf("corpo não pode sobrescrever o autor da sessão, veio %q", fake.editadoPor)
	}
}

func TestEditarMusicaSemArquivoRetorna400(t *testing.T) {
	s := newTestServer(t, "http://127.0.0.1:1")
	s.catalogo = &catalogoEditavelFake{}

	req := httptest.NewRequest("PUT", "/api/musicas", strings.NewReader(`{"titulo":"sem arquivo"}`))
	req = req.WithContext(auth.ContextComEmail(req.Context(), "a@b.com"))
	rec := httptest.NewRecorder()

	s.handleEditarMusica(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, veio %d", rec.Code)
	}
}

func TestEditarMusicaSemCatalogoRetorna503(t *testing.T) {
	s := newTestServer(t, "http://127.0.0.1:1") // catalogo nil
	req := httptest.NewRequest("PUT", "/api/musicas", strings.NewReader(`{"file":"a.mp3"}`))
	req = req.WithContext(auth.ContextComEmail(req.Context(), "a@b.com"))
	rec := httptest.NewRecorder()

	s.handleEditarMusica(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("esperava 503, veio %d", rec.Code)
	}
}
