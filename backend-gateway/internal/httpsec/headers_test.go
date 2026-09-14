package httpsec

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHeadersDefineCSPeHSTS(t *testing.T) {
	h := Headers("https://stream.morcegaofm.com.br/live", true)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Errorf("CSP sem frame-ancestors: %q", csp)
	}
	if !strings.Contains(csp, "https://stream.morcegaofm.com.br/live") {
		t.Errorf("CSP tem de liberar o stream em media-src: %q", csp)
	}
	if !strings.Contains(csp, "img-src 'self' data: https:") {
		t.Errorf("CSP tem de liberar capa https (iTunes e /uploads do site): %q", csp)
	}
	if rec.Header().Get("Strict-Transport-Security") == "" {
		t.Error("faltou HSTS")
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("faltou X-Frame-Options: DENY")
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("faltou nosniff")
	}
	if rec.Header().Get("Referrer-Policy") != "same-origin" {
		t.Error("faltou Referrer-Policy")
	}
}

func TestHeadersSemHSTSForaDeHTTPS(t *testing.T) {
	// Em dev (http://localhost) o HSTS travaria o navegador em https para sempre.
	h := Headers("", false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Header().Get("Strict-Transport-Security") != "" {
		t.Error("HSTS não pode sair quando o gateway não está em https")
	}
}
