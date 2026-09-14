package httpsec

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// Limiter é uma janela fixa por IP, em memória. O gateway é um processo só num
// host só; não vale trazer dependência de Redis para proteger duas contas.
type Limiter struct {
	limite int
	janela time.Duration

	mu    sync.Mutex
	visto map[string]janelaIP

	now func() time.Time
}

type janelaIP struct {
	inicio time.Time
	contou int
}

func NewLimiter(limite int, janela time.Duration) *Limiter {
	return &Limiter{
		limite: limite,
		janela: janela,
		visto:  make(map[string]janelaIP),
		now:    time.Now,
	}
}

// Permitir conta uma tentativa e diz se ela cabe na janela.
func (l *Limiter) Permitir(chave string) bool {
	agora := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	// Limpeza preguiçosa: sem isso o mapa cresce para sempre.
	for k, j := range l.visto {
		if agora.Sub(j.inicio) > l.janela {
			delete(l.visto, k)
		}
	}

	j, ok := l.visto[chave]
	if !ok || agora.Sub(j.inicio) > l.janela {
		l.visto[chave] = janelaIP{inicio: agora, contou: 1}
		return true
	}
	if j.contou >= l.limite {
		return false
	}
	j.contou++
	l.visto[chave] = j
	return true
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !l.Permitir(ip) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Muitas tentativas. Tente de novo em um minuto.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
