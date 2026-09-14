package catalogo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func servidorFake(t *testing.T, chaveEsperada string, resposta any, status int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Servico-Key") != chaveEsperada {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resposta)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestLookupDevolveOQueOSiteConhece(t *testing.T) {
	srv := servidorFake(t, "chave", []map[string]any{
		{"arquivo": "a.mp3", "tipo": "musica", "artista": "AC/DC",
			"titulo": "Back in Black", "ano": 1980, "imagemUrl": "https://capa/512.jpg"},
	}, http.StatusOK)

	c := New(srv.URL, "chave", time.Hour)
	got := c.Lookup(context.Background(), []string{"a.mp3", "b.mp3"})

	if len(got) != 1 {
		t.Fatalf("esperava 1 resultado, veio %d", len(got))
	}
	if got["a.mp3"].ImagemURL != "https://capa/512.jpg" || got["a.mp3"].Ano != 1980 {
		t.Errorf("resultado errado: %+v", got["a.mp3"])
	}
}

func TestLookupDegradaEmFalha(t *testing.T) {
	for _, caso := range []struct {
		nome   string
		status int
		chave  string
	}{
		{"500 do site", http.StatusInternalServerError, "chave"},
		{"401 chave errada", http.StatusOK, "chave-errada"},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			srv := servidorFake(t, "chave", nil, caso.status)
			c := New(srv.URL, caso.chave, time.Hour)
			if got := c.Lookup(context.Background(), []string{"a.mp3"}); len(got) != 0 {
				t.Errorf("falha do site tem de virar mapa vazio, veio %+v", got)
			}
		})
	}
}

func TestLookupDegradaComSiteFora(t *testing.T) {
	// Porta fechada: o painel não pode travar por causa de capa.
	c := New("http://127.0.0.1:1", "chave", time.Hour)
	feito := make(chan map[string]Info, 1)
	go func() { feito <- c.Lookup(context.Background(), []string{"a.mp3"}) }()

	select {
	case got := <-feito:
		if len(got) != 0 {
			t.Errorf("esperava mapa vazio, veio %+v", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Lookup travou; tem de ter timeout curto")
	}
}

func TestLookupDesligadoSemConfiguracao(t *testing.T) {
	c := New("", "", time.Hour)
	if got := c.Lookup(context.Background(), []string{"a.mp3"}); len(got) != 0 {
		t.Errorf("sem URL/chave o cliente fica desligado, veio %+v", got)
	}
}

func TestCacheNegativoNaoReconsulta(t *testing.T) {
	var chamadas int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chamadas++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New(srv.URL, "chave", time.Hour)
	c.Lookup(context.Background(), []string{"desconhecida.mp3"})
	c.Lookup(context.Background(), []string{"desconhecida.mp3"})

	if chamadas != 1 {
		t.Errorf("arquivo desconhecido não pode ser reperguntado a cada render; chamadas=%d", chamadas)
	}
}

func TestCacheExpiraNoTTL(t *testing.T) {
	var chamadas int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chamadas++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"arquivo":"a.mp3","tipo":"musica","imagemUrl":"https://capa/1.jpg"}]`))
	}))
	defer srv.Close()

	agora := time.Unix(1_700_000_000, 0)
	c := New(srv.URL, "chave", time.Hour)
	c.now = func() time.Time { return agora }

	c.Lookup(context.Background(), []string{"a.mp3"})
	c.Lookup(context.Background(), []string{"a.mp3"})
	if chamadas != 1 {
		t.Fatalf("segunda consulta dentro do TTL deveria vir do cache; chamadas=%d", chamadas)
	}

	agora = agora.Add(time.Hour + time.Minute)
	c.Lookup(context.Background(), []string{"a.mp3"})
	if chamadas != 2 {
		t.Errorf("depois do TTL deveria reconsultar; chamadas=%d", chamadas)
	}
}

func TestLookupRespeitaTetoDe200PorChamada(t *testing.T) {
	var lotes []int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var corpo struct {
			Arquivos []string `json:"arquivos"`
		}
		json.NewDecoder(r.Body).Decode(&corpo)
		lotes = append(lotes, len(corpo.Arquivos))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	arquivos := make([]string, 250)
	for i := range arquivos {
		arquivos[i] = fmt.Sprintf("faixa-%d.mp3", i)
	}

	c := New(srv.URL, "chave", time.Hour)
	c.Lookup(context.Background(), arquivos)

	if len(lotes) != 2 {
		t.Fatalf("250 arquivos deveriam virar 2 lotes, veio %v", lotes)
	}
	for _, n := range lotes {
		if n > 200 {
			t.Errorf("lote de %d arquivos estoura o teto do endpoint (200)", n)
		}
	}
}

func TestEditarMandaChaveEAutorEInvalidaCache(t *testing.T) {
	var recebido struct {
		Arquivo    string `json:"arquivo"`
		EditadoPor string `json:"editadoPor"`
		Titulo     string `json:"titulo"`
	}
	var chave string
	var lookups int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chave = r.Header.Get("X-Servico-Key")
		if r.Method == http.MethodPut {
			json.NewDecoder(r.Body).Decode(&recebido)
			w.WriteHeader(http.StatusOK)
			return
		}
		lookups++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"arquivo":"a.mp3","titulo":"antigo"}]`))
	}))
	defer srv.Close()

	c := New(srv.URL, "chave", time.Hour)

	// Popula o cache.
	c.Lookup(context.Background(), []string{"a.mp3"})
	if lookups != 1 {
		t.Fatalf("esperava 1 lookup, veio %d", lookups)
	}

	if err := c.Editar(context.Background(), "a.mp3", "aguergolet@gmail.com",
		Edicao{Titulo: "Back in Black"}); err != nil {
		t.Fatalf("Editar: %v", err)
	}

	if chave != "chave" {
		t.Errorf("a chave de serviço não foi enviada no PUT: %q", chave)
	}
	if recebido.Arquivo != "a.mp3" || recebido.EditadoPor != "aguergolet@gmail.com" {
		t.Errorf("corpo errado: %+v", recebido)
	}

	// Sem invalidar, a correção só apareceria daqui a uma hora.
	c.Lookup(context.Background(), []string{"a.mp3"})
	if lookups != 2 {
		t.Errorf("editar tem de invalidar o cache do arquivo; lookups=%d", lookups)
	}
}

func TestEditarPropagaErroDoSite(t *testing.T) {
	// Diferente de Lookup, Editar devolve erro: o usuário clicou em salvar e
	// precisa saber que não salvou.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := New(srv.URL, "chave", time.Hour)
	if err := c.Editar(context.Background(), "a.mp3", "a@b.com", Edicao{}); err == nil {
		t.Error("404 do site tem de virar erro, não silêncio")
	}
}

func TestEditarDesligadoSemConfiguracao(t *testing.T) {
	c := New("", "", time.Hour)
	if err := c.Editar(context.Background(), "a.mp3", "a@b.com", Edicao{}); err == nil {
		t.Error("sem catálogo configurado, editar tem de falhar explicitamente")
	}
}
