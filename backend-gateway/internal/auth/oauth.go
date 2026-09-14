package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/config"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/httpsec"
)

// DevEmail é a identidade usada pelo login de desenvolvimento (DEV_MODE=true).
const DevEmail = "dev@localhost"

const stateCookieName = "batradio_oauth_state"

type OAuth struct {
	cfg         *config.Config
	oauth       *oauth2.Config
	userinfoURL string
}

func NewOAuth(c *config.Config) *OAuth {
	return &OAuth{
		cfg: c,
		oauth: &oauth2.Config{
			ClientID:     c.GoogleClientID,
			ClientSecret: c.GoogleClientSecret,
			RedirectURL:  c.OAuthRedirectURL,
			Scopes:       []string{"openid", "email"},
			Endpoint:     endpoints.Google,
		},
		userinfoURL: "https://openidconnect.googleapis.com/v1/userinfo",
	}
}

// SetEndpoints troca os endpoints do Google por fakes (testes).
func (o *OAuth) SetEndpoints(authURL, tokenURL, userinfoURL string) {
	o.oauth.Endpoint = oauth2.Endpoint{AuthURL: authURL, TokenURL: tokenURL}
	o.userinfoURL = userinfoURL
}

// Register adiciona as rotas de autenticação a um router existente.
func (o *OAuth) Register(r chi.Router) {
	// 10 tentativas por minuto e por IP: folga para quem erra a conta, teto
	// para quem está sondando o callback.
	limite := httpsec.NewLimiter(10, time.Minute)
	r.With(limite.Middleware).Get("/auth/login", o.handleLogin)
	r.With(limite.Middleware).Get("/auth/callback", o.handleCallback)
	r.Post("/auth/logout", o.handleLogout)
}

func (o *OAuth) Routes() chi.Router {
	r := chi.NewRouter()
	o.Register(r)
	return r
}

func (o *OAuth) setSessionCookie(w http.ResponseWriter, r *http.Request, email string) {
	token := SignSession(email, time.Now().Add(SessionTTL), o.cfg.SessionSecret)
	secure := o.cfg.CookieSecure
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName(secure),
		Value:    token,
		Path:     "/",
		MaxAge:   int(SessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (o *OAuth) handleLogin(w http.ResponseWriter, r *http.Request) {
	if o.cfg.DevMode {
		o.setSessionCookie(w, r, DevEmail)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		writeAuthError(w, http.StatusInternalServerError, "internal", "Falha ao gerar state.")
		return
	}
	state := hex.EncodeToString(buf)
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     "/auth",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   o.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	url := o.oauth.AuthCodeURL(state, oauth2.SetAuthURLParam("prompt", "select_account"))
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *OAuth) handleCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil || stateCookie.Value == "" || r.URL.Query().Get("state") != stateCookie.Value {
		writeAuthError(w, http.StatusBadRequest, "invalid_state", "State do OAuth ausente ou inválido. Tente fazer login novamente.")
		return
	}
	// invalida o cookie de state
	http.SetCookie(w, &http.Cookie{Name: stateCookieName, Path: "/auth", MaxAge: -1})

	code := r.URL.Query().Get("code")
	token, err := o.oauth.Exchange(r.Context(), code)
	if err != nil {
		writeAuthError(w, http.StatusBadGateway, "oauth_exchange_failed", fmt.Sprintf("Falha na troca do código OAuth: %v", err))
		return
	}

	resp, err := o.oauth.Client(r.Context(), token).Get(o.userinfoURL)
	if err != nil {
		writeAuthError(w, http.StatusBadGateway, "userinfo_failed", "Falha ao consultar dados do usuário no Google.")
		return
	}
	defer resp.Body.Close()
	var info struct {
		Email         string `json:"email"`
		EmailVerified any    `json:"email_verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil || info.Email == "" {
		writeAuthError(w, http.StatusBadGateway, "userinfo_failed", "Resposta inválida do Google.")
		return
	}
	if !truthy(info.EmailVerified) {
		writeAuthError(w, http.StatusForbidden, "email_not_verified", "O e-mail da conta Google não é verificado.")
		return
	}
	if !IsAllowed(info.Email, o.cfg.AllowedEmails) {
		http.Redirect(w, r, "/login?error=email_not_allowed&email="+url.QueryEscape(info.Email), http.StatusFound)
		return
	}
	o.setSessionCookie(w, r, info.Email)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (o *OAuth) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Expira os dois nomes: um cookie gravado antes de o Secure entrar não pode
	// sobreviver ao logout.
	for _, nome := range []string{SessionCookieHost, SessionCookie} {
		http.SetCookie(w, &http.Cookie{
			Name:     nome,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   o.cfg.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}

// truthy aceita bool ou string "true" (o Google já retornou ambos).
func truthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t == "true"
	}
	return false
}
