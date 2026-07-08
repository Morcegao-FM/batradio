// Package config carrega e valida a configuração do gateway a partir de
// variáveis de ambiente.
package config

import (
	"fmt"
	"strings"
)

type Config struct {
	Port               string
	NodeBackendURL     string
	NodeAPIKey         string
	GoogleClientID     string
	GoogleClientSecret string
	OAuthRedirectURL   string
	AllowedEmails      []string
	SessionSecret      []byte
	StreamURL          string
	DevMode            bool
}

func Load(getenv func(string) string) (*Config, error) {
	c := &Config{
		Port:               getenv("PORT"),
		NodeBackendURL:     strings.TrimRight(getenv("NODE_BACKEND_URL"), "/"),
		NodeAPIKey:         getenv("NODE_API_KEY"),
		GoogleClientID:     getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getenv("GOOGLE_CLIENT_SECRET"),
		OAuthRedirectURL:   getenv("OAUTH_REDIRECT_URL"),
		SessionSecret:      []byte(getenv("SESSION_SECRET")),
		StreamURL:          getenv("STREAM_URL"),
		DevMode:            strings.EqualFold(getenv("DEV_MODE"), "true"),
	}
	if c.Port == "" {
		c.Port = "8080"
	}
	if c.NodeBackendURL == "" {
		return nil, fmt.Errorf("NODE_BACKEND_URL é obrigatório")
	}
	if c.NodeAPIKey == "" {
		return nil, fmt.Errorf("NODE_API_KEY é obrigatório")
	}
	if len(c.SessionSecret) < 32 {
		return nil, fmt.Errorf("SESSION_SECRET deve ter pelo menos 32 bytes")
	}
	for _, e := range strings.Split(getenv("ALLOWED_EMAILS"), ",") {
		e = strings.ToLower(strings.TrimSpace(e))
		if e != "" {
			c.AllowedEmails = append(c.AllowedEmails, e)
		}
	}
	if len(c.AllowedEmails) == 0 {
		return nil, fmt.Errorf("ALLOWED_EMAILS deve conter pelo menos um e-mail")
	}
	if c.DevMode {
		// O login de desenvolvimento assina sessão como dev@localhost; precisa
		// estar na allowlist para passar no middleware.
		c.AllowedEmails = append(c.AllowedEmails, "dev@localhost")
	} else {
		if c.GoogleClientID == "" || c.GoogleClientSecret == "" || c.OAuthRedirectURL == "" {
			return nil, fmt.Errorf("GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET e OAUTH_REDIRECT_URL são obrigatórios fora do DEV_MODE")
		}
	}
	return c, nil
}
