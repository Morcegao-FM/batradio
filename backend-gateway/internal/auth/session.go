// Package auth implementa sessões em cookie assinado (HMAC-SHA256) e o fluxo
// de login com Google OAuth restrito a uma allowlist de e-mails.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

const SessionCookie = "batradio_session"

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
