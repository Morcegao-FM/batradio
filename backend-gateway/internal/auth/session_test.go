package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var secret = []byte(strings.Repeat("k", 32))

func TestSessionRoundTrip(t *testing.T) {
	exp := time.Now().Add(time.Hour)
	tok := SignSession("aguergolet@gmail.com", exp, secret)
	email, ok := VerifySession(tok, secret, time.Now())
	if !ok || email != "aguergolet@gmail.com" {
		t.Fatalf("ok=%v email=%q", ok, email)
	}
}

func TestSessionExpired(t *testing.T) {
	tok := SignSession("a@b.com", time.Now().Add(-time.Minute), secret)
	if _, ok := VerifySession(tok, secret, time.Now()); ok {
		t.Fatal("expected expired session to fail")
	}
}

func TestSessionTampered(t *testing.T) {
	tok := SignSession("a@b.com", time.Now().Add(time.Hour), secret)
	// troca um byte no meio
	b := []byte(tok)
	b[2] ^= 0x01
	if _, ok := VerifySession(string(b), secret, time.Now()); ok {
		t.Fatal("expected tampered token to fail")
	}
}

func TestSessionWrongSecret(t *testing.T) {
	tok := SignSession("a@b.com", time.Now().Add(time.Hour), secret)
	other := []byte(strings.Repeat("x", 32))
	if _, ok := VerifySession(tok, other, time.Now()); ok {
		t.Fatal("expected wrong secret to fail")
	}
}

func TestSessionGarbage(t *testing.T) {
	for _, tok := range []string{"", "abc", "a.b", "a.b.c.d", "!!!.123.???"} {
		if _, ok := VerifySession(tok, secret, time.Now()); ok {
			t.Fatalf("expected garbage %q to fail", tok)
		}
	}
}

func TestIsAllowed(t *testing.T) {
	allowed := []string{"aguergolet@gmail.com", "morcegaofm@gmail.com"}
	if !IsAllowed("Aguergolet@Gmail.COM", allowed) {
		t.Error("case-insensitive match failed")
	}
	if IsAllowed("intruso@gmail.com", allowed) {
		t.Error("unlisted email allowed")
	}
	if IsAllowed("", allowed) {
		t.Error("empty email allowed")
	}
}

func protectedHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := EmailFromContext(r.Context())
		if email == "" {
			t.Error("email missing from context")
		}
		w.WriteHeader(200)
	})
}

func doRequest(h http.Handler, cookie string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/api/x", nil)
	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: SessionCookie, Value: cookie})
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestMiddleware(t *testing.T) {
	allowed := []string{"aguergolet@gmail.com"}
	mw := RequireSession(secret, allowed)
	h := mw(protectedHandler(t))

	// sem cookie → 401
	if rec := doRequest(h, ""); rec.Code != 401 {
		t.Errorf("no cookie: code=%d", rec.Code)
	}
	// cookie inválido → 401
	if rec := doRequest(h, "lixo"); rec.Code != 401 {
		t.Errorf("bad cookie: code=%d", rec.Code)
	}
	// e-mail assinado mas fora da allowlist (removido depois) → 403
	tok := SignSession("removido@gmail.com", time.Now().Add(time.Hour), secret)
	rec := doRequest(h, tok)
	if rec.Code != 403 {
		t.Errorf("delisted email: code=%d", rec.Code)
	}
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Error.Code != "email_not_allowed" {
		t.Errorf("error code=%q", body.Error.Code)
	}
	// sessão válida → 200 com email no context
	tok = SignSession("aguergolet@gmail.com", time.Now().Add(time.Hour), secret)
	if rec := doRequest(h, tok); rec.Code != 200 {
		t.Errorf("valid: code=%d", rec.Code)
	}
}

func TestCookieNamePorModo(t *testing.T) {
	if got := CookieName(true); got != "__Host-batradio_session" {
		t.Errorf("com Secure esperava prefixo __Host-, veio %q", got)
	}
	// O prefixo __Host- exige Secure; em dev (http://localhost) o navegador
	// recusaria o cookie e o login nunca fecharia.
	if got := CookieName(false); got != "batradio_session" {
		t.Errorf("sem Secure esperava nome simples, veio %q", got)
	}
}

func TestLerSessaoAceitaOsDoisNomes(t *testing.T) {
	for _, nome := range []string{SessionCookieHost, SessionCookie} {
		r := httptest.NewRequest("GET", "/api/status", nil)
		r.AddCookie(&http.Cookie{Name: nome, Value: "abc"})
		valor, ok := LerSessao(r)
		if !ok || valor != "abc" {
			t.Errorf("%s: esperava ler abc, veio %q ok=%v", nome, valor, ok)
		}
	}

	r := httptest.NewRequest("GET", "/api/status", nil)
	if _, ok := LerSessao(r); ok {
		t.Error("sem cookie não pode reportar sessão")
	}
}
