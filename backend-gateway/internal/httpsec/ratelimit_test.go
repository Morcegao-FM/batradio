package httpsec

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimiterBloqueiaAcimaDoLimite(t *testing.T) {
	agora := time.Unix(1_700_000_000, 0)
	l := NewLimiter(3, time.Minute)
	l.now = func() time.Time { return agora }

	h := l.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	pedir := func() int {
		r := httptest.NewRequest("GET", "/auth/login", nil)
		r.RemoteAddr = "203.0.113.7:12345"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}

	for i := 0; i < 3; i++ {
		if got := pedir(); got != http.StatusOK {
			t.Fatalf("tentativa %d deveria passar, veio %d", i+1, got)
		}
	}
	if got := pedir(); got != http.StatusTooManyRequests {
		t.Fatalf("4ª tentativa deveria dar 429, veio %d", got)
	}

	// Passada a janela, libera de novo.
	agora = agora.Add(time.Minute + time.Second)
	if got := pedir(); got != http.StatusOK {
		t.Fatalf("depois da janela deveria liberar, veio %d", got)
	}
}

func TestLimiterSeparaPorIP(t *testing.T) {
	l := NewLimiter(1, time.Minute)
	h := l.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	pedir := func(ip string) int {
		r := httptest.NewRequest("GET", "/auth/login", nil)
		r.RemoteAddr = ip + ":1111"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}

	pedir("203.0.113.1")
	if got := pedir("203.0.113.2"); got != http.StatusOK {
		t.Fatalf("IP diferente não pode herdar o limite do vizinho, veio %d", got)
	}
}
