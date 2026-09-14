package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
)

func TestAuditoriaRegistraEscritaComEmail(t *testing.T) {
	var linhas []string
	h := auditoria(func(f string, a ...any) { linhas = append(linhas, fmt.Sprintf(f, a...)) })(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))

	req := httptest.NewRequest("POST", "/api/playlist/remove", nil)
	req = req.WithContext(auth.ContextComEmail(context.Background(), "aguergolet@gmail.com"))
	h.ServeHTTP(httptest.NewRecorder(), req)

	if len(linhas) != 1 {
		t.Fatalf("esperava 1 linha de auditoria, veio %d: %v", len(linhas), linhas)
	}
	for _, esperado := range []string{"aguergolet@gmail.com", "POST", "/api/playlist/remove", "200"} {
		if !strings.Contains(linhas[0], esperado) {
			t.Errorf("linha de auditoria sem %q: %s", esperado, linhas[0])
		}
	}
}

func TestAuditoriaIgnoraLeitura(t *testing.T) {
	var linhas []string
	h := auditoria(func(f string, a ...any) { linhas = append(linhas, fmt.Sprintf(f, a...)) })(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// GET de status roda a cada 2 segundos; auditar leitura só afogaria o log.
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/api/status", nil))

	if len(linhas) != 0 {
		t.Errorf("leitura não deve ser auditada, veio: %v", linhas)
	}
}
