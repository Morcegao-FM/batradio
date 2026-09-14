// BatRadio Gateway: autentica via Google OAuth, expõe a API nova na frente do
// backend Node/MPD e serve o frontend React embutido.
package main

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/catalogo"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/config"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/httpapi"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/httpsec"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/node"
	"github.com/Morcegao-FM/batradio/backend-gateway/web"
)

func main() {
	// Healthcheck do container: a imagem é distroless (sem shell, sem curl),
	// então o próprio binário faz o GET e traduz em código de saída.
	if len(os.Args) > 1 && os.Args[1] == "--health-check" {
		porta := os.Getenv("PORT")
		if porta == "" {
			porta = "8080"
		}
		cliente := &http.Client{Timeout: 3 * time.Second}
		resp, err := cliente.Get("http://127.0.0.1:" + porta + "/healthz")
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		resp.Body.Close()
		os.Exit(0)
	}

	// Fora do Docker, carrega o .env do diretório atual ou da raiz do repo
	// (variáveis já exportadas no ambiente têm precedência).
	for _, path := range []string{".env", "../.env"} {
		loaded, err := config.LoadDotEnv(path)
		if err != nil {
			log.Fatalf("falha ao ler %s: %v", path, err)
		}
		if loaded {
			log.Printf("variáveis carregadas de %s", path)
			break
		}
	}

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		log.Fatalf("configuração inválida: %v", err)
	}
	if cfg.DevMode {
		log.Println("ATENÇÃO: DEV_MODE ativo — login sem Google, NÃO use em produção")
	}

	cat := catalogo.New(cfg.CatalogoURL, cfg.CatalogoChave, time.Hour)
	if cfg.CatalogoURL == "" {
		log.Println("catálogo do site não configurado — o painel roda sem capa")
	}

	nodeClient := node.New(cfg.NodeBackendURL, cfg.NodeAPIKey)
	server := httpapi.NewServer(cfg, nodeClient, cat)
	oauth := auth.NewOAuth(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Poller de status (SSE) e refresh inicial do acervo com retry: o Node pode
	// demorar a subir junto no docker compose.
	go server.Poller().Run(ctx)
	go func() {
		for {
			songs, err := server.Node().GetFiles(ctx, true)
			if err == nil {
				server.Library().Set(songs, time.Now())
				log.Printf("acervo indexado: %d faixas", len(songs))
				return
			}
			log.Printf("acervo indisponível (%v), tentando de novo em 5s", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	// CookieSecure é o proxy de "estou atrás de https": é exatamente a condição
	// em que o HSTS deve sair.
	r.Use(httpsec.Headers(cfg.StreamURL, cfg.CookieSecure))
	r.Mount("/", server.Routes(oauth))
	r.NotFound(spaHandler())

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("gateway ouvindo em :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("servidor: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("encerrando...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

// spaHandler serve os estáticos do build do React com fallback para
// index.html (rotas client-side como /playlists).
func spaHandler() http.HandlerFunc {
	dist, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		log.Fatalf("embed do frontend: %v", err)
	}
	fileServer := http.FileServerFS(dist)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/auth/") {
			http.NotFound(w, r)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path != "" {
			if _, err := fs.Stat(dist, path); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
		}
		// fallback SPA
		index, err := fs.ReadFile(dist, "index.html")
		if err != nil {
			http.Error(w, "frontend não embarcado", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(index)
	}
}
