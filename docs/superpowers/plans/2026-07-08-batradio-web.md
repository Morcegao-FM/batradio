# BatRadio Web Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Substituir o cliente Windows Forms por um gateway Go (auth Google + API paginada na frente do backend Node/MPD) e um frontend React com o design system Morcegão FM.

**Architecture:** O Go é a única porta de entrada: autentica via Google OAuth (allowlist), mantém índice em memória do acervo, calcula horários da fila, publica status via SSE e serve o build do React. O Node existente fica na rede interna e continua sendo o único processo que fala MPD.

**Tech Stack:** Go 1.23+ (chi, golang.org/x/oauth2), React 18 + TypeScript + Vite, TanStack Query + TanStack Virtual, vitest, Docker Compose.

**Spec:** `docs/superpowers/specs/2026-07-08-batradio-web-design.md` (leia antes de executar).

## Global Constraints

- Allowlist inicial: `aguergolet@gmail.com,morcegaofm@gmail.com` (env `ALLOWED_EMAILS`, comparação case-insensitive de e-mails exatos).
- O acervo NUNCA é enviado inteiro ao browser: toda listagem tem `offset/limit` (limit máx. 200) e `total`.
- API JSON: erros no formato `{"error":{"code":"...","message":"..."}}`; 401 sem sessão, 403 e-mail fora da allowlist, 502 Node indisponível, 400 validação.
- Textos da UI em pt-BR, como no wireframe (ex.: "TOCANDO AGORA", "FILA DE HOJE", "CARREGAR E TOCAR").
- Cookies de sessão: HttpOnly, SameSite=Lax, Secure quando request via HTTPS; validade 7 dias.
- Node backend não é modificado.
- Go não instalado na máquina: Task 0 instala em `~/.local/go` (sem sudo). Docker indisponível localmente — arquivos são escritos e validados por revisão, não por build.
- Commits frequentes, mensagens em inglês, sufixo `Co-Authored-By: Claude Fable 5 <noreply@anthropic.com>`.

---

### Task 0: Toolchain Go

**Files:** nenhum no repo.

- [ ] **Step 1:** Baixar e extrair Go (linux-amd64, ~70 MB, fonte oficial go.dev):

```bash
mkdir -p ~/.local && cd ~/.local && \
curl -sSLo go.tgz https://go.dev/dl/go1.23.4.linux-amd64.tar.gz && \
rm -rf ~/.local/go && tar -xzf go.tgz && rm go.tgz
```

- [ ] **Step 2:** Verificar: `~/.local/go/bin/go version` → `go version go1.23.4 linux/amd64`. Exportar `PATH=$HOME/.local/go/bin:$PATH` em cada sessão de shell (o Bash tool não persiste env; prefixar comandos ou usar caminho completo).

---

### Task 1: Scaffold do módulo + config

**Files:**
- Create: `backend-gateway/go.mod`, `backend-gateway/internal/config/config.go`
- Test: `backend-gateway/internal/config/config_test.go`

**Interfaces (Produces):**

```go
package config

type Config struct {
    Port            string   // default "8080"
    NodeBackendURL  string   // obrigatório, ex. http://node-backend:9320
    NodeAPIKey      string   // obrigatório
    GoogleClientID  string
    GoogleClientSecret string
    OAuthRedirectURL string  // ex. https://painel.morcegaofm.com.br/auth/callback
    AllowedEmails   []string // lowercase, trim
    SessionSecret   []byte   // obrigatório, >=32 bytes
    StreamURL       string   // opcional
    DevMode         bool     // DEV_MODE=true: login fake sem Google
}
func Load(getenv func(string) string) (*Config, error)
```

- [ ] **Step 1:** `cd backend-gateway && go mod init github.com/Morcegao-FM/batradio/backend-gateway`
- [ ] **Step 2:** Teste falhando (`config_test.go`): `Load` com map fake — casos: completo OK (emails normalizados p/ lowercase, espaços trimados); sem `NODE_BACKEND_URL` → erro; `SESSION_SECRET` curto → erro; defaults (Port 8080). Rodar `go test ./internal/config/` → FAIL (não compila).
- [ ] **Step 3:** Implementar `Load` lendo `PORT`, `NODE_BACKEND_URL`, `NODE_API_KEY`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `ALLOWED_EMAILS` (split por vírgula), `SESSION_SECRET`, `STREAM_URL`, `DEV_MODE`. Em DevMode, Google vars podem faltar; fora dele são obrigatórias.
- [ ] **Step 4:** `go test ./internal/config/` → PASS.
- [ ] **Step 5:** Commit `feat(gateway): scaffold module and config loading`.

---

### Task 2: Tipos de domínio + cliente do Node

**Files:**
- Create: `backend-gateway/internal/model/model.go`, `backend-gateway/internal/node/client.go`
- Test: `backend-gateway/internal/node/client_test.go`, `backend-gateway/internal/model/model_test.go`

**Interfaces (Produces):**

```go
package model

// FlexInt/FlexFloat aceitam número OU string no JSON (MPD devolve strings).
type Song struct {
    File   string `json:"file"`
    Artist string `json:"artist"`
    Title  string `json:"title"`
    Album  string `json:"album,omitempty"`
    Genre  string `json:"genre,omitempty"`
    Time   int    `json:"time"`           // segundos
    Pos    int    `json:"pos"`
    ID     int    `json:"id"`
}
type Status struct {
    State          string  `json:"state"` // play|pause|stop
    Song           int     `json:"song"`  // pos da atual
    Elapsed        float64 `json:"elapsed"`
    Repeat         bool    `json:"repeat"`
    Random         bool    `json:"random"`
    Crossfade      bool    `json:"crossfade"`
    PlaylistLength int     `json:"playlistLength"`
    Current        *Song   `json:"current,omitempty"`
    Next           *Song   `json:"next,omitempty"`
}
type PlaylistInfo struct {
    Name         string `json:"name"`
    LastModified string `json:"lastModified,omitempty"`
}

package node
type Client struct { /* baseURL, apiKey, http.Client{Timeout: 35s} */ }
func New(baseURL, apiKey string) *Client
func (c *Client) GetStatus(ctx) (model.Status, error)
func (c *Client) GetPlaylist(ctx, update bool) ([]model.Song, error)     // POST /playlist
func (c *Client) GetFiles(ctx, update bool) ([]model.Song, error)        // GET /list type=file
func (c *Client) GetPlaylists(ctx) ([]model.PlaylistInfo, error)         // GET /list type=playlist
func (c *Client) PlayOrPause(ctx) (model.Status, error)
func (c *Client) Play(ctx, pos int) (model.Status, error)
func (c *Client) Repeat(ctx) (model.Status, error)                       // toggle
func (c *Client) Shuffle(ctx) (model.Status, error)
func (c *Client) Crossfade(ctx) (model.Status, error)                    // POST /fadein
func (c *Client) Add(ctx, files []string, pos int) ([]model.Song, error)
func (c *Client) Delete(ctx, positions []int) ([]model.Song, error)
func (c *Client) Move(ctx, from, to int) ([]model.Song, error)
func (c *Client) LoadPlaylist(ctx, name string) ([]model.Song, error)
func (c *Client) SavePlaylist(ctx, name string) error                    // POST /saveplaylist
func (c *Client) RemovePlaylist(ctx, name string) error
var ErrNodeUnavailable = errors.New(...)  // wraps rede/5xx
```

O cliente replica o protocolo legado do Node: parâmetros via **headers** (`batradio-apikey`, `update`, `type`, `files` JSON, `position`, `from-pos`, `to-pos`, `name`). Parse do JSON do Node: chaves MPD (`Artist`, `Title`, `Time`, `Pos`, `Id`, `Genre`, `Album`, `file`) com valores string → decodificar via structs intermediárias com `FlexInt`.

- [ ] **Step 1:** Testes de `model`: `FlexInt` decodifica `"3"` e `3`; `Song` decodifica payload MPD real (fixture com `{"file":"a.mp3","Artist":"AC/DC","Time":"352","Pos":"0","Id":"7"}`). FAIL.
- [ ] **Step 2:** Implementar `model` (UnmarshalJSON de FlexInt/FlexFloat; mapping das chaves via struct `mpdSong` interna com método `ToSong()`). PASS.
- [ ] **Step 3:** Testes de `node` com `httptest.Server` verificando: header `batradio-apikey` enviado; `GetFiles` manda `type: file`; `Add` manda `files` como JSON string e `position`; status 500 do Node → `ErrNodeUnavailable`; conexão recusada → `ErrNodeUnavailable`; status Node `{"state":"play","song":"2","elapsed":"10.5","repeat":"1","xfade":"0","playlistlength":"8","currentSong":{...}}` → `model.Status` correto. FAIL.
- [ ] **Step 4:** Implementar cliente. PASS em `go test ./...`.
- [ ] **Step 5:** Commit `feat(gateway): domain model and node backend client`.

---

### Task 3: Índice do acervo (busca/paginação)

**Files:**
- Create: `backend-gateway/internal/library/index.go`
- Test: `backend-gateway/internal/library/index_test.go`

**Interfaces (Produces):**

```go
package library
type Index struct { /* mu sync.RWMutex; songs []model.Song; keys []string; refreshedAt time.Time */ }
func NewIndex() *Index
func (i *Index) Set(songs []model.Song, at time.Time)      // pré-computa chave normalizada "file title album artist genre" em lowercase
func (i *Index) Search(q string, offset, limit int) (items []model.Song, total int)
func (i *Index) Stats() (count int, refreshedAt time.Time)
```

Semântica de busca = cliente Windows (`StatusSong.Matches`): split de `q` por espaço; faixa entra se **todos** os termos são substring (case-insensitive) da chave. `q` vazio → lista completa paginada, ordenada por `File`. `limit<=0` → 50; `limit>200` → 200.

- [ ] **Step 1:** Testes: busca multi-termo ("zeppelin whole" acha "Whole Lotta Love"; termo ausente exclui); acentos não normalizados (comportamento igual ao legado, só case-fold); paginação `offset/limit` + `total` correto com filtro; limite máx. 200; concorrência básica (`Set` durante `Search` com `-race`). FAIL.
- [ ] **Step 2:** Implementar. `go test -race ./internal/library/` PASS.
- [ ] **Step 3:** Commit `feat(gateway): in-memory library index with legacy search semantics`.

---

### Task 4: Cálculo de horários + inserção distribuída

**Files:**
- Create: `backend-gateway/internal/schedule/schedule.go`
- Test: `backend-gateway/internal/schedule/schedule_test.go`

**Interfaces (Produces):**

```go
package schedule
type QueueItem struct {
    model.Song
    NextPresentation time.Time `json:"nextPresentation"`
}
// Port fiel de BatRadioClient.TimeCalculation (frmMain): âncora = faixa atual
// (now - elapsed); para frente acumula durações; para trás: repeat=false subtrai,
// repeat=true continua acumulando após o fim (wrap).
func ComputeTimes(songs []model.Song, st model.Status, now time.Time) []QueueItem
// Inserção distribuída: gera timestamps alvo a partir de `from` (ou now) a cada
// `interval`, limitado a `times` ocorrências e (se to!=nil) ao período; para cada
// alvo escolhe a posição da primeira faixa com NextPresentation >= alvo (ou fim da
// fila). Retorna posições de inserção em ordem DECRESCENTE (inserir de trás pra
// frente preserva índices).
func PlanInsertions(queue []QueueItem, times int, interval time.Duration, from, to *time.Time, now time.Time) []int
```

- [ ] **Step 1:** Testes `ComputeTimes` com fila de 4 faixas (60s cada), atual pos=2, elapsed=30, now fixo:
  - repeat=false: pos2=now-30s; pos3=pos2+60s; pos1=pos2-60s; pos0=pos1-60s.
  - repeat=true: pos0 e pos1 vêm DEPOIS do fim (pos3+60s, +120s).
  - faixa atual ausente da fila (song=99) → NextPresentation zero para todos (retorna sem calcular, como o C#).
  Testes `PlanInsertions`: fila de 24h, times=6, interval=4h → 6 posições decrescentes; com `to` curto → menos posições; alvo além do fim → posição len(queue). FAIL.
- [ ] **Step 2:** Implementar. PASS.
- [ ] **Step 3:** Commit `feat(gateway): queue time calculation and distributed insertion planning`.

---

### Task 5: Sessões + middleware de auth

**Files:**
- Create: `backend-gateway/internal/auth/session.go`, `backend-gateway/internal/auth/middleware.go`
- Test: `backend-gateway/internal/auth/session_test.go`

**Interfaces (Produces):**

```go
package auth
// token = base64url(email) + "." + expUnix + "." + base64url(hmacSHA256(email+"."+exp, secret))
func SignSession(email string, exp time.Time, secret []byte) string
func VerifySession(token string, secret []byte, now time.Time) (email string, ok bool)
const SessionCookie = "batradio_session"
// Middleware: lê cookie, verifica assinatura+expiração+allowlist; injeta email no
// context (EmailFromContext(ctx) string); falha → 401 JSON.
func RequireSession(secret []byte, allowed []string) func(http.Handler) http.Handler
func IsAllowed(email string, allowed []string) bool // case-insensitive exato
```

- [ ] **Step 1:** Testes: round-trip sign/verify; expirado → !ok; assinatura adulterada (flip byte) → !ok; segredo diferente → !ok; `IsAllowed` case-insensitive e rejeita não listado; middleware 401 sem cookie, 401 cookie inválido, 403 e-mail válido fora da allowlist (assinado mas removido da lista), 200 com email no context. FAIL.
- [ ] **Step 2:** Implementar (stdlib: crypto/hmac, sha256, base64url). PASS.
- [ ] **Step 3:** Commit `feat(gateway): signed session cookies and auth middleware`.

---

### Task 6: Fluxo OAuth Google

**Files:**
- Create: `backend-gateway/internal/auth/oauth.go`
- Test: `backend-gateway/internal/auth/oauth_test.go`
- Modify: `backend-gateway/go.mod` (add `golang.org/x/oauth2`, `github.com/go-chi/chi/v5`)

**Interfaces (Produces):**

```go
package auth
type OAuth struct { /* cfg *oauth2.Config; secret []byte; allowed []string; devMode bool; userinfoURL string (override em teste) */ }
func NewOAuth(c *config.Config) *OAuth
// GET /auth/login    → set cookie "oauth_state" (10min, HMAC do valor aleatório) + redirect Google (scope "openid email", prompt=select_account)
// GET /auth/callback → confere state, troca code, GET userinfo, exige email_verified,
//                      IsAllowed? set session cookie + redirect "/" : redirect "/login?error=email_not_allowed&email=..."
// POST /auth/logout  → expira cookie, 204
// GET  /api/me       → {"email": "..."} (protegido pelo middleware)
// DevMode: /auth/login vira login imediato como "dev@localhost" (allowlist ignorada), sem Google.
func (o *OAuth) Routes() chi.Router
```

- [ ] **Step 1:** Testes com servidor OAuth fake (`httptest`: endpoints token e userinfo): callback feliz → Set-Cookie de sessão + redirect `/`; e-mail fora da lista → redirect `/login?error=email_not_allowed`; state ausente/errado → 400; `email_verified:false` → 403; devMode → sessão direta. FAIL.
- [ ] **Step 2:** Implementar com `golang.org/x/oauth2` (endpoint google, `AuthCodeURL(state)`, `Exchange`, userinfo `https://openidconnect.googleapis.com/v1/userinfo` com URL injetável). `go mod tidy`. PASS.
- [ ] **Step 3:** Commit `feat(gateway): google oauth login with email allowlist`.

---

### Task 7: API HTTP — status, acervo, SSE, server info

**Files:**
- Create: `backend-gateway/internal/httpapi/api.go`, `backend-gateway/internal/httpapi/sse.go`, `backend-gateway/internal/httpapi/poller.go`
- Test: `backend-gateway/internal/httpapi/api_test.go`

**Interfaces (Produces):**

```go
package httpapi
type Server struct { /* node *node.Client; lib *library.Index; cfg *config.Config; poller *Poller */ }
func NewServer(...) *Server
func (s *Server) Routes() chi.Router  // monta /api/* protegido + /auth/*
// GET  /api/status          → model.Status (do cache do poller; se vazio, consulta Node)
// GET  /api/library         → {"items":[Song],"total":N,"offset":o,"limit":l}
// POST /api/library/refresh → node.GetFiles(update=true) → lib.Set; {"count":N}
// GET  /api/server/info     → {"nodeOk":bool,"libraryCount":N,"libraryRefreshedAt":ts,"streamUrl":"..."}
// GET  /api/events          → SSE; eventos: "status" (JSON model.Status, a cada tick
//                             de 2s) e "playlist" (JSON {"version":N}) quando versão muda.
type Poller struct{}
func NewPoller(n *node.Client, interval time.Duration) *Poller
func (p *Poller) Run(ctx context.Context)
func (p *Poller) Last() (model.Status, bool)
func (p *Poller) Subscribe() (ch <-chan Event, cancel func())
func (p *Poller) BumpPlaylist() // mutações locais incrementam versão
// Poller também incrementa versão quando status.PlaylistLength muda.
```

- [ ] **Step 1:** Testes (Node fake via `httptest`): 401 sem cookie em `/api/status`; com cookie válido: status OK; library com `q/offset/limit` e total; refresh popula índice; server/info reflete Node caído (`nodeOk:false`, demais campos preenchidos); SSE: subscrever, injetar tick, receber evento `status` formatado (`event: status\ndata: {...}\n\n`); BumpPlaylist gera evento `playlist`. Erro Node → 502 `{"error":{"code":"node_unavailable"}}`. FAIL.
- [ ] **Step 2:** Implementar (helpers `writeJSON`, `writeErr`; SSE com `http.Flusher`, heartbeat comment a cada 15s, cleanup no context done). PASS `go test ./... -race`.
- [ ] **Step 3:** Commit `feat(gateway): status/library/sse http api`.

---

### Task 8: API HTTP — fila, player, playlists

**Files:**
- Modify: `backend-gateway/internal/httpapi/api.go` (novos handlers em `playlist.go`, `player.go`, `playlists.go` no mesmo pacote)
- Test: `backend-gateway/internal/httpapi/playlist_test.go`

**Interfaces (Produces):** rotas conforme spec:

```
GET  /api/playlist?q=&offset=&limit=  → {"items":[QueueItem],"total":N,"currentPos":P,"version":V}
     (busca na fila usa a mesma semântica Matches; horários via schedule.ComputeTimes
      com status do poller; cache da fila no Server, atualizado por mutações e por
      version bump — GET /playlist com cache frio busca do Node update=true)
GET  /api/playlist/current            → {"pos":P,"index":I} (I = índice na lista sem filtro)
POST /api/player/toggle|shuffle|repeat|crossfade → model.Status
POST /api/player/play {"position":N}  → model.Status (sem confirmação no backend; a dupla
                                        confirmação é do frontend)
POST /api/playlist/add {"files":[],"position":N}        → 204 + bump
POST /api/playlist/add-many {"file":"","timesPerDay":N,"intervalHours":H,"from":ts?,"to":ts?}
     → {"insertedAt":[pos...]} (usa schedule.PlanInsertions; insere de trás pra frente)
POST /api/playlist/move {"from":N,"to":N}               → 204 + bump
POST /api/playlist/remove {"positions":[N]}             → 204 + bump (ordena desc antes de deletar)
GET  /api/playlists                    → {"items":[{"name","lastModified"}]}
POST /api/playlists {"name":""}        → 204 (salva fila atual)
POST /api/playlists/{name}/load {"playAtRandom":bool} → model.Status
     (load; se playAtRandom: rand pos em [1,len-1] e Play — comportamento do botão WinForms)
DELETE /api/playlists/{name}           → 204
```

- [ ] **Step 1:** Testes com Node fake: fila paginada com horários corretos (status conhecido → nextPresentation determinístico com clock injetado); `q` filtra e `total` reflete filtro; add/move/remove chamam Node com headers certos e bumpam versão; remove com posições [2,5,1] deleta na ordem 5,2,1; add-many insere nas posições planejadas (Node fake registra chamadas); load com playAtRandom chama load + play com 1<=pos<len (rand injetável `randInt func(int) int`); validação: position<0 → 400, timesPerDay<1 ou >96 → 400, name vazio → 400 (e name com `/` → 400). FAIL.
- [ ] **Step 2:** Implementar. PASS `go test ./... -race`.
- [ ] **Step 3:** Commit `feat(gateway): queue, player and saved playlists api`.

---

### Task 9: main.go + embed do frontend + modo dev

**Files:**
- Create: `backend-gateway/main.go`, `backend-gateway/web/web.go` (embed), `backend-gateway/web/dist/.gitkeep`, `backend-gateway/.env.example`
- Test: smoke manual via curl

**Interfaces:**

```go
// main: config.Load(os.Getenv) → node.New → library.NewIndex → NewPoller →
// httpapi.NewServer → chi mount: /auth/*, /api/* (RequireSession), /* estáticos.
// Estáticos: embed.FS de web/dist com fallback SPA (qualquer 404 não-/api → index.html).
// Refresh inicial do acervo em goroutine (com retry 30s até Node responder).
// Graceful shutdown (SIGTERM, 10s).
package web
//go:embed all:dist
var Dist embed.FS
```

`.env.example` documenta todas as env vars com exemplos.

- [ ] **Step 1:** Implementar main.go + web/web.go (`fs.Sub`, `http.FileServerFS`, fallback lendo index.html; se dist vazio → 503 "frontend não embarcado" em `/`).
- [ ] **Step 2:** `go vet ./... && go build ./...` limpos; subir com DEV_MODE=true + Node fake (`go run` + curl): `curl -c jar -b jar localhost:8080/auth/login` (segue redirect) e depois `curl -b jar localhost:8080/api/me` → `{"email":"dev@localhost"}`.
- [ ] **Step 3:** Commit `feat(gateway): main entrypoint, spa embed and dev mode`.

---

### Task 10: Frontend — scaffold, tokens, API client, login

**Files:**
- Create: `frontend-web/` (Vite react-ts): `package.json`, `vite.config.ts` (proxy `/api` e `/auth` → `localhost:8080`; build outDir `../backend-gateway/web/dist`, emptyOutDir), `src/styles/tokens.css`, `src/styles/base.css`, `src/lib/api.ts`, `src/lib/types.ts`, `src/hooks/useAuth.ts`, `src/pages/Login.tsx`, `src/App.tsx`, `src/main.tsx`
- Test: `src/lib/api.test.ts` (vitest)

**Interfaces (Produces):**

```ts
// types.ts — espelha o JSON do gateway
export interface Song { file: string; artist: string; title: string; album?: string; genre?: string; time: number; pos: number; id: number }
export interface QueueItem extends Song { nextPresentation: string }
export interface Status { state: 'play'|'pause'|'stop'; song: number; elapsed: number; repeat: boolean; random: boolean; crossfade: boolean; playlistLength: number; current?: Song; next?: Song }
export interface Page<T> { items: T[]; total: number; offset: number; limit: number }

// api.ts
export class ApiError extends Error { code: string; status: number }
export async function api<T>(path: string, init?: RequestInit): Promise<T>
// 401 → window.location = '/login'; body de erro {error:{code,message}} → ApiError

// useAuth.ts
export function useAuth(): { email?: string; loading: boolean }  // GET /api/me
```

`tokens.css`: TODAS as 123 variáveis extraídas do wireframe (paleta `--bat-*`, `--red/amber/blue/live-*`, aliases semânticos `--bg-page`, `--surface-card`, `--brand`, `--accent`, `--highlight`, `--live`, fontes `--font-display:'Anton'...`, `--text-*`, `--space-*`, `--radius-*`, sombras/glows, `--header-h:72px`, `--player-h:84px`, `--dur-*`). Fontes via `@fontsource/anton`, `@fontsource/oswald`, `@fontsource/barlow`, `@fontsource/space-mono`.

Login: fundo `--bg-page` com grain sutil, logo (copiar `logo.png` do cliente Windows para `src/assets/`), título display "BATRADIO", botão "ENTRAR COM GOOGLE" (estilo brand); se `?error=email_not_allowed`, card de acesso negado com o e-mail e aviso "Peça acesso ao administrador".

- [ ] **Step 1:** `npm create vite@latest frontend-web -- --template react-ts`; instalar deps: `@tanstack/react-query @tanstack/react-virtual react-router-dom @fontsource/anton @fontsource/oswald @fontsource/barlow @fontsource/space-mono`; dev: `vitest @testing-library/react jsdom msw`.
- [ ] **Step 2:** Teste vitest de `api()`: resposta ok tipada; erro JSON vira `ApiError` com code; (mock `fetch`). FAIL → implementar → PASS.
- [ ] **Step 3:** tokens.css + base.css (reset, scrollbars escuras, ::selection âmbar) + Login + App com rotas (`/login` pública; demais atrás de `RequireAuth` que usa `useAuth`).
- [ ] **Step 4:** `npm run build` limpo; `npm test` PASS.
- [ ] **Step 5:** Commit `feat(web): scaffold with design tokens, api client and login`.

---

### Task 11: Frontend — layout (sidebar/header/player) + SSE

**Files:**
- Create: `src/components/Layout.tsx`, `src/components/Sidebar.tsx`, `src/components/HeaderBar.tsx`, `src/components/PlayerBar.tsx`, `src/hooks/useRadioEvents.ts`, `src/hooks/useStatus.ts`
- Test: `src/hooks/useRadioEvents.test.ts`

**Interfaces (Produces):**

```ts
// useRadioEvents: abre EventSource('/api/events');
//  - evento 'status' → queryClient.setQueryData(['status'], parsed)
//  - evento 'playlist' → queryClient.invalidateQueries({queryKey:['playlist']}) e ['playlists']
//  - desconexão → reconecta com backoff (1s..30s); expõe {connected: boolean}
export function useRadioEvents(): { connected: boolean }
export function useStatus(): { status?: Status }   // useQuery(['status'], GET /api/status, staleTime Infinity — SSE alimenta)
```

Layout do wireframe: sidebar 200px (logo, badge "AO VIVO" + host do servidor via `/api/server/info`, nav TOCANDO AGORA / PLAYLIST E ARQUIVOS / CONFIGURAÇÕES com item ativo em vermelho, rodapé "BATRADIO · V2.0"); header com título da página, subtítulo, relógio `HH:MM · SEG` (Space Mono, tick 1s) e botão "↻ ATUALIZAR" (refresh library + invalidate tudo); player bar fixa embaixo (84px): botão play/pause redondo vermelho (`POST /api/player/toggle`), thumb, "TOCANDO AGORA" + equalizer animado quando play, título — artista, badge "● AO VIVO" verde, ícone alto-falante + slider de volume controlando `<audio src={streamUrl}>` (renderiza só se `streamUrl` configurado; autoplay off). Banner global vermelho "Servidor da rádio inacessível" quando status query erra com 502 ou SSE desconectado por >10s.

- [ ] **Step 1:** Teste de `useRadioEvents` com EventSource mockado (classe fake global): evento status atualiza cache; playlist invalida. FAIL → implementar → PASS.
- [ ] **Step 2:** Implementar componentes com CSS modules (`Layout.module.css` etc.) usando os tokens.
- [ ] **Step 3:** `npm run build && npm test` PASS. Verificação visual: `npm run dev` + gateway DEV_MODE + Node fake script (`node backend-gateway/hack/fake-node.mjs`, criado aqui: servidor http nas rotas do Node com dataset sintético de 50k faixas) — checar telas no browser (preview tool).
- [ ] **Step 4:** Commit `feat(web): app shell with sidebar, header, player bar and sse`.

---

### Task 12: Frontend — Tocando Agora (acervo + fila virtualizados)

**Files:**
- Create: `src/pages/NowPlaying.tsx`, `src/components/LibraryPanel.tsx`, `src/components/QueuePanel.tsx`, `src/components/TrackRow.tsx`, `src/components/GenreChip.tsx`, `src/hooks/useLibrarySearch.ts`, `src/hooks/useQueue.ts`, `src/components/ConfirmDialog.tsx`
- Test: `src/hooks/useLibrarySearch.test.ts`

**Interfaces (Produces):**

```ts
// useLibrarySearch: debounce 300ms de q; useInfiniteQuery(['library', q], páginas de 100,
// getNextPageParam por offset+limit<total). Retorna {rows, total, fetchNextPage, isFetching}
// useQueue: idem para ['playlist', q]; expõe também currentPos (do status) e version.
```

Comportamentos (wireframe + WinForms):
- Acervo: input "Banda, música, álbum..."; mínimo 0 chars (server pagina de qualquer forma); linhas com avatar-letra (primeira letra do artista, fundo `--surface-raised`, borda por gênero), título, "artista · m:ss", chip de gênero (VINHETA âmbar, demais neutros), botão "..." abre menu (Inserir várias vezes — modal na Task 14).
- Seleção única no acervo; seleção única na fila. Botões "↑ ADICIONAR ACIMA" (neutro) / "↓ ADICIONAR ABAIXO" (vermelho) habilitados só com seleção em ambos painéis → `POST /api/playlist/add` na posição da seleção da fila (acima = pos, abaixo = pos+1, como o WinForms).
- Fila: header "FILA DE HOJE" + "N faixas"; busca local server-side na fila; botão "» MÚSICA ATUAL" (azul) → `GET /api/playlist/current` e scrolla o virtualizador até o índice; linha atual com borda âmbar + horário destacado; cada linha: horário previsto (HH:MM, Space Mono) + `#pos`, ações on-hover: ▲ ▼ (move ±1), ▶ (tocar: **duas** confirmações via ConfirmDialog, textos do WinForms), ✕ (remover).
- TanStack Virtual nos dois painéis (`estimateSize` 56px, overscan 10); scroll até 90% dispara `fetchNextPage`.

- [ ] **Step 1:** Teste `useLibrarySearch`: debounce (fake timers) — digitação rápida gera 1 fetch; mudança de q reseta páginas. FAIL → implementar → PASS.
- [ ] **Step 2:** Implementar componentes; estilos fiéis (painel acervo `border-top: 2px solid var(--accent)`, fila `var(--brand)`; cabeçalhos Oswald caps; horários Space Mono).
- [ ] **Step 3:** `npm run build && npm test`; verificação no browser com fake-node (50k faixas: scroll fluido, busca "zeppelin whole" filtra, adicionar acima/abaixo funciona, mover/remover/tocar com confirmações).
- [ ] **Step 4:** Commit `feat(web): now playing screen with virtualized library and queue`.

---

### Task 13: Frontend — Playlists + Configurações

**Files:**
- Create: `src/pages/Playlists.tsx`, `src/pages/Settings.tsx`, `src/components/Toggle.tsx`, `src/components/PlaylistCard.tsx`
- Test: nenhum novo (cobertura via build + verificação manual)

Playlists (wireframe pág. 2): card topo "CARREGAR PLAYLIST DO SERVIDOR" com select + botão vermelho "► CARREGAR E TOCAR" (ConfirmDialog: "Substitui a fila atual..."; chama `/api/playlists/{name}/load {playAtRandom:true}`); grid de cards (borda superior alternando red/blue/amber por índice) com nome, "N FAIXAS · DATA" (lastModified; contagem omitida se indisponível), botão "CARREGAR ESTA" (`playAtRandom:false`); card extra "Salvar fila atual" com input de nome + botão salvar; ação de excluir playlist no card ("..." → confirmar → DELETE).

Configurações: card "COMO A RÁDIO TOCA" (borda âmbar) com 4 Toggles ligados ao status (Tocar playlist→toggle, Shuffle→shuffle, Repetir→repeat, Transição→crossfade; otimista com rollback em erro); card "SERVIDOR" (borda azul) somente leitura: endereço do Node (host de `/api/server/info`), status conexão (ok/falha), faixas no acervo, última atualização, botão azul "↻ ATUALIZAR DADOS DO SERVIDOR" (`POST /api/library/refresh`, spinner enquanto roda).

- [ ] **Step 1:** Implementar páginas/componentes.
- [ ] **Step 2:** `npm run build && npm test`; verificação manual das duas telas com fake-node.
- [ ] **Step 3:** Commit `feat(web): playlists and settings screens`.

---

### Task 14: Frontend — modal "Inserir várias vezes"

**Files:**
- Create: `src/components/AddManyModal.tsx`
- Modify: `src/components/LibraryPanel.tsx` (menu "..." abre o modal)
- Test: `src/components/AddManyModal.test.tsx`

Modal fiel ao wireframe: header "PROGRAMAÇÃO / INSERIR MÚSICA VÁRIAS VEZES", subtítulo "A faixa será distribuída automaticamente pela programação.", card da faixa selecionada, campos "VEZES POR DIA" (number 1–96, default 6) e "INTERVALO" (select: A cada 1/2/3/4/6/8/12 horas, default 4), checkbox "Apenas por um período" habilitando De/Até (datetime-local, default agora → +1 dia, como o form C#); CANCELAR / "CONFIRMAR INSERÇÃO" (vermelho) → `POST /api/playlist/add-many`; toast com "Inserida em N posições".

- [ ] **Step 1:** Teste: render com faixa, submit chama fetch com payload correto (msw); período desabilitado omite from/to. FAIL → implementar → PASS.
- [ ] **Step 2:** `npm run build && npm test` PASS; verificação manual.
- [ ] **Step 3:** Commit `feat(web): add-many-times modal`.

---

### Task 15: Docker + docs + verificação final

**Files:**
- Create: `Dockerfile.gateway` (raiz), `docker-compose.yml` (raiz), `backend-server/Dockerfile` já existe (`backend-server/Dockerfile` — reutilizar)
- Modify: `README.md`

`Dockerfile.gateway` multi-stage: stage 1 `node:22-alpine` builda `frontend-web` (outDir dist local); stage 2 `golang:1.23-alpine` copia `backend-gateway/` + dist para `web/dist`, `CGO_ENABLED=0 go build`; stage 3 `gcr.io/distroless/static` com o binário, `EXPOSE 8080`. `docker-compose.yml`: serviço `node-backend` (build `./backend-server`, sem ports públicos, volume `./config.json`), `gateway` (build Dockerfile.gateway, ports `8080:8080`, env_file `.env`, depends_on node-backend). README: seção nova "BatRadio Web" — arquitetura, como criar credenciais OAuth no Google Cloud Console (origem + redirect URI), env vars, como rodar dev (fake-node + DEV_MODE) e produção (compose), nota de que o cliente Windows foi substituído.

- [ ] **Step 1:** Escrever Dockerfile.gateway + docker-compose.yml (sem como buildar localmente — validar por leitura cuidadosa; anotar no relatório final).
- [ ] **Step 2:** README atualizado.
- [ ] **Step 3:** Verificação final completa: `go vet ./... && go test ./... -race` no gateway; `npm test && npm run build` no frontend; build embed: `go build` com dist populado; smoke dev completo (gateway + fake-node + browser: login dev, as 3 telas, uma mutação de fila).
- [ ] **Step 4:** Commit `chore: docker compose deployment and docs for batradio web`.

---

## Self-review (feito na escrita)

- Cobertura do spec: auth/allowlist (T5-6), paginação/índice (T3, T7), horários e add-many (T4, T8), SSE (T7, T11), telas do wireframe (T10-14), adaptações aprovadas (T11 volume monitor, T13 card somente leitura), Docker/env (T9, T15). Login screen (T10). Erro Node→502+banner (T7, T11).
- Tipos consistentes entre tasks (Song/Status/QueueItem/Page definidos em T2/T4/T10 e referenciados idênticos).
- Sem placeholders: cada task tem arquivos exatos, contratos e casos de teste nomeados.
