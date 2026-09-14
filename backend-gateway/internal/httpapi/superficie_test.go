package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// publicasPermitidas é a lista FECHADA de rotas que podem responder sem sessão.
// Acrescentar item aqui é uma decisão de segurança — não faça isso para calar
// um teste vermelho.
var publicasPermitidas = map[string]bool{
	"GET /healthz":       true,
	"GET /auth/login":    true,
	"GET /auth/callback": true,
	"POST /auth/logout":  true,
}

// TestNenhumaRotaNovaSemSessao percorre as rotas registradas no chi e exige que
// toda rota fora da allowlist responda 401 sem cookie de sessão. É a rede que
// pega a rota nova que alguém esqueceu de proteger.
func TestNenhumaRotaNovaSemSessao(t *testing.T) {
	srv := newTestServer(t, "http://127.0.0.1:1")
	router := srv.Routes(nil)

	var verificadas int
	err := chi.Walk(router, func(metodo, rota string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		chave := metodo + " " + rota
		if publicasPermitidas[chave] {
			return nil
		}

		caminho := strings.ReplaceAll(rota, "{name}", "qualquer")
		if strings.Contains(caminho, "{") {
			t.Errorf("%s: parâmetro de rota não previsto no teste (%s)", chave, caminho)
			return nil
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(metodo, caminho, nil))

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s respondeu %d sem sessão; esperava 401. "+
				"Rota nova precisa nascer atrás de RequireSession.", chave, rec.Code)
		}
		verificadas++
		return nil
	})
	if err != nil {
		t.Fatalf("chi.Walk: %v", err)
	}
	if verificadas == 0 {
		t.Fatal("nenhuma rota verificada — o teste não está enxergando o router")
	}
	t.Logf("%d rotas verificadas atrás de sessão", verificadas)
}

func TestHealthzEhPublicoENaoVazaNada(t *testing.T) {
	srv := newTestServer(t, "http://127.0.0.1:1")
	router := srv.Routes(nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("GET", "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("healthz respondeu %d", rec.Code)
	}
	corpo := rec.Body.String()
	for _, proibido := range []string{"nodeHost", "libraryCount", "version", "9320"} {
		if strings.Contains(corpo, proibido) {
			t.Errorf("healthz é público e não pode expor %q; corpo: %s", proibido, corpo)
		}
	}
}
