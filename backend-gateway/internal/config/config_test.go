package config

import (
	"strings"
	"testing"
)

func fakeEnv(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func validEnv() map[string]string {
	return map[string]string{
		"NODE_BACKEND_URL":     "http://node-backend:9320",
		"NODE_API_KEY":         "secret-key",
		"GOOGLE_CLIENT_ID":     "cid",
		"GOOGLE_CLIENT_SECRET": "csecret",
		"OAUTH_REDIRECT_URL":   "https://painel.morcegaofm.com.br/auth/callback",
		"ALLOWED_EMAILS":       " Aguergolet@Gmail.com , morcegaofm@gmail.com ",
		"SESSION_SECRET":       strings.Repeat("s", 32),
	}
}

func TestLoadValid(t *testing.T) {
	c, err := Load(fakeEnv(validEnv()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Port != "8080" {
		t.Errorf("default port: got %q", c.Port)
	}
	want := []string{"aguergolet@gmail.com", "morcegaofm@gmail.com"}
	if len(c.AllowedEmails) != 2 || c.AllowedEmails[0] != want[0] || c.AllowedEmails[1] != want[1] {
		t.Errorf("emails not normalized: %v", c.AllowedEmails)
	}
	if len(c.SessionSecret) != 32 {
		t.Errorf("session secret length: %d", len(c.SessionSecret))
	}
}

func TestLoadMissingNodeURL(t *testing.T) {
	env := validEnv()
	delete(env, "NODE_BACKEND_URL")
	if _, err := Load(fakeEnv(env)); err == nil {
		t.Fatal("expected error for missing NODE_BACKEND_URL")
	}
}

func TestLoadShortSecret(t *testing.T) {
	env := validEnv()
	env["SESSION_SECRET"] = "short"
	if _, err := Load(fakeEnv(env)); err == nil {
		t.Fatal("expected error for short SESSION_SECRET")
	}
}

func TestLoadMissingGoogleOutsideDevMode(t *testing.T) {
	env := validEnv()
	delete(env, "GOOGLE_CLIENT_ID")
	if _, err := Load(fakeEnv(env)); err == nil {
		t.Fatal("expected error for missing google credentials")
	}
	env["DEV_MODE"] = "true"
	// Em dev o login é local; o redirect https só existe em produção, e agora
	// Load recusa a combinação DEV_MODE + https (login falso no ar).
	env["OAUTH_REDIRECT_URL"] = ""
	c, err := Load(fakeEnv(env))
	if err != nil {
		t.Fatalf("dev mode should allow missing google creds: %v", err)
	}
	if !c.DevMode {
		t.Error("DevMode not set")
	}
	found := false
	for _, e := range c.AllowedEmails {
		if e == "dev@localhost" {
			found = true
		}
	}
	if !found {
		t.Error("dev@localhost should be allowlisted in dev mode")
	}
}

func TestLoadCustomPortAndStream(t *testing.T) {
	env := validEnv()
	env["PORT"] = "9000"
	env["STREAM_URL"] = "https://stream.morcegaofm.com.br/live"
	c, err := Load(fakeEnv(env))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Port != "9000" || c.StreamURL != "https://stream.morcegaofm.com.br/live" {
		t.Errorf("got port=%q stream=%q", c.Port, c.StreamURL)
	}
}

func TestLoadNoEmails(t *testing.T) {
	env := validEnv()
	env["ALLOWED_EMAILS"] = " , "
	if _, err := Load(fakeEnv(env)); err == nil {
		t.Fatal("expected error for empty allowlist")
	}
}

func TestLoadRecusaDevModeComRedirectHTTPS(t *testing.T) {
	env := validEnv()
	env["DEV_MODE"] = "true"
	// validEnv já traz OAUTH_REDIRECT_URL https.

	if _, err := Load(fakeEnv(env)); err == nil {
		t.Fatal("esperava erro: DEV_MODE com redirect https é login falso em produção")
	} else if !strings.Contains(err.Error(), "DEV_MODE") {
		t.Fatalf("erro deveria citar DEV_MODE, veio: %v", err)
	}
}

func TestCookieSecurePadraoPorModo(t *testing.T) {
	prod, err := Load(fakeEnv(validEnv()))
	if err != nil {
		t.Fatalf("produção: %v", err)
	}
	if !prod.CookieSecure {
		t.Error("fora do DEV_MODE o cookie tem de ser Secure por padrão")
	}

	dev := validEnv()
	dev["DEV_MODE"] = "true"
	dev["OAUTH_REDIRECT_URL"] = ""
	dc, err := Load(fakeEnv(dev))
	if err != nil {
		t.Fatalf("dev: %v", err)
	}
	if dc.CookieSecure {
		t.Error("em DEV_MODE (http://localhost) o cookie Secure impediria o login")
	}

	forcado := validEnv()
	forcado["COOKIE_SECURE"] = "false"
	fc, err := Load(fakeEnv(forcado))
	if err != nil {
		t.Fatalf("forçado: %v", err)
	}
	if fc.CookieSecure {
		t.Error("COOKIE_SECURE=false deve vencer o padrão")
	}
}

func TestCatalogoEhOpcional(t *testing.T) {
	c, err := Load(fakeEnv(validEnv()))
	if err != nil {
		t.Fatalf("catálogo não configurado não pode impedir o gateway de subir: %v", err)
	}
	if c.CatalogoURL != "" || c.CatalogoChave != "" {
		t.Error("sem env, o catálogo tem de vir desligado")
	}
}
