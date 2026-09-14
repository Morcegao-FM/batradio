package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ctxKey int

const emailKey ctxKey = 0

// EmailFromContext retorna o e-mail autenticado ou "".
func EmailFromContext(ctx context.Context) string {
	email, _ := ctx.Value(emailKey).(string)
	return email
}

// ContextComEmail injeta o e-mail autenticado. Existe para os testes montarem
// um request já autenticado sem forjar cookie.
func ContextComEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailKey, email)
}

func writeAuthError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}

// RequireSession valida o cookie de sessão e a allowlist; injeta o e-mail no
// context. A allowlist é reavaliada a cada request: remover um e-mail da env
// revoga o acesso mesmo com sessão assinada válida.
func RequireSession(secret []byte, allowed []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			valor, ok := LerSessao(r)
			if !ok {
				writeAuthError(w, http.StatusUnauthorized, "unauthenticated", "Sessão ausente. Faça login.")
				return
			}
			email, ok := VerifySession(valor, secret, time.Now())
			if !ok {
				writeAuthError(w, http.StatusUnauthorized, "unauthenticated", "Sessão inválida ou expirada. Faça login novamente.")
				return
			}
			if !IsAllowed(email, allowed) {
				writeAuthError(w, http.StatusForbidden, "email_not_allowed", "Este e-mail não tem permissão de acesso.")
				return
			}
			next.ServeHTTP(w, r.WithContext(ContextComEmail(r.Context(), email)))
		})
	}
}
