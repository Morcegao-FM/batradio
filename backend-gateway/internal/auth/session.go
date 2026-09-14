// Package auth implementa sessões em cookie assinado (HMAC-SHA256) e o fluxo
// de login com Google OAuth restrito a uma allowlist de e-mails.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Nomes do cookie de sessão. O prefixo __Host- é uma trava do navegador: ele só
// aceita o cookie se vier com Secure, Path=/ e sem Domain — então uma regressão
// que desligue o Secure quebra o login em vez de degradar em silêncio. Em dev
// (http://localhost) o prefixo é impossível, daí o nome simples.
const (
	SessionCookie     = "batradio_session"
	SessionCookieHost = "__Host-batradio_session"
)

// CookieName devolve o nome a gravar conforme o cookie vá ou não com Secure.
func CookieName(secure bool) string {
	if secure {
		return SessionCookieHost
	}
	return SessionCookie
}

// LerSessao procura o cookie de sessão nos dois nomes possíveis. Aceitar ambos
// evita deslogar todo mundo no deploy que liga o Secure.
func LerSessao(r *http.Request) (string, bool) {
	for _, nome := range []string{SessionCookieHost, SessionCookie} {
		if c, err := r.Cookie(nome); err == nil && c.Value != "" {
			return c.Value, true
		}
	}
	return "", false
}

// SessionTTL é a validade da sessão.
const SessionTTL = 7 * 24 * time.Hour

func sign(payload string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// SignSession gera "base64url(email).expUnix.assinatura".
func SignSession(email string, exp time.Time, secret []byte) string {
	payload := base64.RawURLEncoding.EncodeToString([]byte(email)) + "." + strconv.FormatInt(exp.Unix(), 10)
	return payload + "." + sign(payload, secret)
}

// VerifySession valida assinatura e expiração; retorna o e-mail.
func VerifySession(token string, secret []byte, now time.Time) (string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", false
	}
	payload := parts[0] + "." + parts[1]
	expected := sign(payload, secret)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return "", false
	}
	expUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || now.After(time.Unix(expUnix, 0)) {
		return "", false
	}
	emailBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false
	}
	return string(emailBytes), true
}

// IsAllowed compara o e-mail (case-insensitive, exato) com a allowlist.
func IsAllowed(email string, allowed []string) bool {
	if email == "" {
		return false
	}
	email = strings.ToLower(email)
	for _, a := range allowed {
		if email == a {
			return true
		}
	}
	return false
}
