package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/config"
)

// fakeGoogle simula os endpoints de token e userinfo.
func fakeGoogle(t *testing.T, email string, verified bool) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.FormValue("code") != "good-code" {
			w.WriteHeader(400)
			fmt.Fprint(w, `{"error":"invalid_grant"}`)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"access_token":"at-123","token_type":"Bearer","expires_in":3600}`)
	})
	mux.HandleFunc("GET /userinfo", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer at-123" {
			w.WriteHeader(401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"email": email, "email_verified": verified, "sub": "123",
		})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func testConfig(devMode bool) *config.Config {
	return &config.Config{
		Port:               "8080",
		GoogleClientID:     "cid",
		GoogleClientSecret: "csecret",
		OAuthRedirectURL:   "http://localhost:8080/auth/callback",
		AllowedEmails:      []string{"aguergolet@gmail.com", "morcegaofm@gmail.com"},
		SessionSecret:      []byte(strings.Repeat("k", 32)),
		DevMode:            devMode,
	}
}

func newTestOAuth(t *testing.T, email string, verified bool) *OAuth {
	g := fakeGoogle(t, email, verified)
	o := NewOAuth(testConfig(false))
	o.SetEndpoints(g.URL+"/auth", g.URL+"/token", g.URL+"/userinfo")
	return o
}

// startLogin faz GET /auth/login e retorna o state e o cookie de state.
func startLogin(t *testing.T, o *OAuth) (state string, stateCookie *http.Cookie) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/login", nil)
	o.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("login: code=%d", rec.Code)
	}
	loc, err := url.Parse(rec.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	state = loc.Query().Get("state")
	if state == "" {
		t.Fatal("state missing from redirect")
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == stateCookieName {
			stateCookie = c
		}
	}
	if stateCookie == nil {
		t.Fatal("state cookie not set")
	}
	return state, stateCookie
}

func callback(t *testing.T, o *OAuth, state string, stateCookie *http.Cookie, code string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/callback?state="+url.QueryEscape(state)+"&code="+code, nil)
	if stateCookie != nil {
		req.AddCookie(stateCookie)
	}
	o.Routes().ServeHTTP(rec, req)
	return rec
}

func TestCallbackHappyPath(t *testing.T) {
	o := newTestOAuth(t, "aguergolet@gmail.com", true)
	state, cookie := startLogin(t, o)
	rec := callback(t, o, state, cookie, "good-code")
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/" {
		t.Fatalf("code=%d loc=%q body=%s", rec.Code, rec.Header().Get("Location"), rec.Body.String())
	}
	var session string
	for _, c := range rec.Result().Cookies() {
		if c.Name == SessionCookie {
			session = c.Value
			if !c.HttpOnly {
				t.Error("session cookie must be HttpOnly")
			}
		}
	}
	if session == "" {
		t.Fatal("session cookie not set")
	}
	email, ok := VerifySession(session, testConfig(false).SessionSecret, time.Now())
	if !ok || email != "aguergolet@gmail.com" {
		t.Fatalf("session invalid: ok=%v email=%q", ok, email)
	}
}

func TestCallbackEmailNotAllowed(t *testing.T) {
	o := newTestOAuth(t, "intruso@gmail.com", true)
	state, cookie := startLogin(t, o)
	rec := callback(t, o, state, cookie, "good-code")
	if rec.Code != http.StatusFound {
		t.Fatalf("code=%d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, "/login?error=email_not_allowed") {
		t.Fatalf("loc=%q", loc)
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == SessionCookie && c.Value != "" {
			t.Fatal("session cookie must not be set")
		}
	}
}

func TestCallbackBadState(t *testing.T) {
	o := newTestOAuth(t, "aguergolet@gmail.com", true)
	_, cookie := startLogin(t, o)
	// state da URL não bate com o cookie
	rec := callback(t, o, "state-forjado", cookie, "good-code")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("forged state: code=%d", rec.Code)
	}
	// sem cookie de state
	state, _ := startLogin(t, o)
	rec = callback(t, o, state, nil, "good-code")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("missing cookie: code=%d", rec.Code)
	}
}

func TestCallbackUnverifiedEmail(t *testing.T) {
	o := newTestOAuth(t, "aguergolet@gmail.com", false)
	state, cookie := startLogin(t, o)
	rec := callback(t, o, state, cookie, "good-code")
	if rec.Code != http.StatusForbidden {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestCallbackBadCode(t *testing.T) {
	o := newTestOAuth(t, "aguergolet@gmail.com", true)
	state, cookie := startLogin(t, o)
	rec := callback(t, o, state, cookie, "bad-code")
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestDevModeLogin(t *testing.T) {
	o := NewOAuth(testConfig(true))
	rec := httptest.NewRecorder()
	o.Routes().ServeHTTP(rec, httptest.NewRequest("GET", "/auth/login", nil))
	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/" {
		t.Fatalf("code=%d loc=%q", rec.Code, rec.Header().Get("Location"))
	}
	var session string
	for _, c := range rec.Result().Cookies() {
		if c.Name == SessionCookie {
			session = c.Value
		}
	}
	email, ok := VerifySession(session, testConfig(true).SessionSecret, time.Now())
	if !ok || email != DevEmail {
		t.Fatalf("ok=%v email=%q", ok, email)
	}
}

func TestLogout(t *testing.T) {
	o := NewOAuth(testConfig(true))
	rec := httptest.NewRecorder()
	o.Routes().ServeHTTP(rec, httptest.NewRequest("POST", "/auth/logout", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("code=%d", rec.Code)
	}
	found := false
	for _, c := range rec.Result().Cookies() {
		if c.Name == SessionCookie && c.MaxAge < 0 {
			found = true
		}
	}
	if !found {
		t.Fatal("session cookie not expired")
	}
}
