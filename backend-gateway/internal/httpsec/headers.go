// Package httpsec reúne os middlewares de borda do gateway: headers de
// segurança e limite de tentativas no login.
package httpsec

import (
	"net/http"
	"strings"
)

// Headers aplica os headers de segurança da resposta.
//
// O SPA e os assets são servidos pelo próprio binário (go:embed), então a CSP
// pode ser fechada: nada de script externo. As exceções são a capa (https:,
// porque vem do iTunes ou do /uploads do site) e o stream do Icecast.
func Headers(streamURL string, hsts bool) func(http.Handler) http.Handler {
	mediaSrc := "'self'"
	if u := strings.TrimSpace(streamURL); u != "" {
		mediaSrc += " " + u
	}

	csp := strings.Join([]string{
		"default-src 'self'",
		"img-src 'self' data: https:",
		"media-src " + mediaSrc,
		"script-src 'self'",
		// O Vite injeta estilo inline no build; sem unsafe-inline o painel fica sem CSS.
		"style-src 'self' 'unsafe-inline'",
		"connect-src 'self'",
		"font-src 'self' data:",
		"object-src 'none'",
		"frame-ancestors 'none'",
		"base-uri 'self'",
		"form-action 'self'",
	}, "; ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("Content-Security-Policy", csp)
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "same-origin")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			if hsts {
				h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			next.ServeHTTP(w, r)
		})
	}
}
