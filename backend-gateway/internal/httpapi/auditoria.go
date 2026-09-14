package httpapi

import (
	"net/http"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
)

// gravadorStatus lembra o status escrito, que o http.ResponseWriter não expõe.
type gravadorStatus struct {
	http.ResponseWriter
	status int
}

func (g *gravadorStatus) WriteHeader(status int) {
	g.status = status
	g.ResponseWriter.WriteHeader(status)
}

func (g *gravadorStatus) Write(b []byte) (int, error) {
	if g.status == 0 {
		g.status = http.StatusOK
	}
	return g.ResponseWriter.Write(b)
}

// auditoria registra uma linha por escrita: quem, o quê, com que resultado.
//
// Só escrita. O painel faz GET de status a cada 2 segundos; auditar leitura
// afogaria o log justamente no dia em que você for procurar quem tirou uma
// música da fila. Não registra corpo de request — não vale o risco de um
// segredo cair no log por acidente.
func auditoria(log func(formato string, args ...any)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}

			g := &gravadorStatus{ResponseWriter: w}
			next.ServeHTTP(g, r)

			email := auth.EmailFromContext(r.Context())
			if email == "" {
				email = "(sem sessão)"
			}
			if g.status == 0 {
				g.status = http.StatusOK
			}
			log("auditoria: %s %s %s -> %d", email, r.Method, r.URL.Path, g.status)
		})
	}
}
