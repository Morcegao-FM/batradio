# BatRadio Web — Design (Go gateway + React)

**Data:** 2026-07-08 · **Status:** aprovado pelo usuário

## Objetivo

Substituir o cliente Windows Forms (`frontend-windows-client/BatRadio.UI`) por uma
aplicação web: frontend React e backend em Go. O backend Node existente
(`backend-server/server.js`) continua sendo o único processo que fala com o MPD e passa a
rodar apenas na rede interna; o Go atua como **gateway** na frente dele.

O acervo de músicas é gigantesco: **nunca** é carregado inteiro no browser. Busca e
paginação acontecem no servidor; o frontend usa listas virtualizadas.

## Decisões aprovadas

1. **Go na frente do Node** — o Go não fala MPD diretamente; consome a API do Node
   (header `batradio-apikey`) pela rede interna do Docker.
2. **API nova, limpa** — o cliente Windows é aposentado; sem compatibilidade com os
   endpoints antigos (parâmetros via headers).
3. **Autenticação Google OAuth** (Authorization Code, server-side) com allowlist de
   e-mails exatos: `aguergolet@gmail.com` e `morcegaofm@gmail.com` (env `ALLOWED_EMAILS`).
4. **Deploy via Docker Compose** no servidor da rádio; o binário Go serve o build do
   React (`go:embed`); TLS fica no proxy do host.
5. **Adaptações do wireframe** (aprovadas):
   - Card "Conexão de Streaming" das Configurações vira **status somente leitura**
     (Node/MPD alcançável, tamanho do acervo, última atualização) + botão "Atualizar
     dados do servidor". Endereço/porta/chave são env vars do gateway.
   - O slider de **volume** do player controla o monitoramento local: o browser toca o
     stream Icecast (env `STREAM_URL`); play/pause grande continua controlando o MPD.

## Arquitetura

```
Browser (React SPA, design system Morcegão FM)
   │ HTTPS + cookie de sessão (HttpOnly, assinado)
   ▼
Go Gateway  ── batradio-apikey ──►  Node backend (interno)  ──►  MPD ──► Icecast
 • OAuth Google + allowlist
 • Índice em memória do acervo (busca/paginação)
 • Cálculo de horários da fila (port do TimeCalculation)
 • SSE de status (1 poll no Node → N browsers)
 • Serve estáticos do React (go:embed)
```

## Backend Go (`backend-gateway/`)

Go 1.22+, router `chi`, sem framework pesado. Pacotes:

- `internal/config` — env vars: `PORT`, `NODE_BACKEND_URL`, `NODE_API_KEY`,
  `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `ALLOWED_EMAILS`
  (CSV), `SESSION_SECRET`, `STREAM_URL` (opcional), `DEV_MODE`.
- `internal/auth` — fluxo OAuth (endpoints `/auth/login`, `/auth/callback`,
  `/auth/logout`, `GET /api/me`); sessão = cookie HMAC-assinado (email + exp, 7 dias);
  middleware que exige sessão válida em `/api/*`; e-mail fora da allowlist → 403 com
  motivo (`email_not_allowed`) para a tela de acesso negado. State do OAuth em cookie
  de curta duração.
- `internal/node` — cliente HTTP do backend Node (timeout 30s): `GetStatus`,
  `GetPlaylist(update)`, `GetList(type, update)`, `PlayOrPause`, `Play(pos)`,
  `AddToPlaylist(files, pos)`, `Delete(positions)`, `Move(from,to)`, `Repeat`,
  `Shuffle`, `FadeIn`, `LoadPlaylist(name)`, `SavePlaylist(name)`,
  `RemovePlaylist(name)`. Erros do Node → `502` com mensagem.
- `internal/library` — índice em memória do acervo (RWMutex): normalização
  lowercase de `artist/title/album/file`; busca com semântica do cliente Windows
  (todos os termos devem bater); paginação `offset/limit`; `Refresh()` re-busca do
  Node (`update=true`); guarda `lastRefresh` e contagem.
- `internal/schedule` — port do `TimeCalculation` do C#: dado o status (song,
  elapsed, repeat) e a fila, calcula `nextPresentation` de cada faixa. Também o
  cálculo do "inserir várias vezes": dado `timesPerDay`/`interval`/período opcional,
  gera timestamps alvo e resolve as posições de inserção pela grade de horários.
- `internal/httpapi` — handlers + SSE. Poller de status: a cada 2s consulta o Node
  **uma vez** e publica para todos os assinantes SSE; incrementa `playlistVersion`
  quando a fila muda (qualquer mutação local ou mudança de `playlistlength`).

### API (`/api`, JSON, autenticada por cookie exceto `/auth/*`)

| Método/rota | Descrição |
|---|---|
| `GET /api/me` | e-mail logado ou 401 |
| `GET /api/status` | status atual (estado, faixa, elapsed, flags shuffle/repeat/xfade) |
| `GET /api/events` | SSE: `status` (2s) e `playlist` (versão p/ invalidação) |
| `GET /api/library?q=&offset=&limit=` | acervo paginado `{items,total}` (limit ≤ 200) |
| `POST /api/library/refresh` | re-indexa do Node (`update=true`) |
| `GET /api/server/info` | status da conexão Node, contagem do acervo, última atualização, `streamUrl` |
| `GET /api/playlist?q=&offset=&limit=` | fila paginada com `nextPresentation` e `total`; `q` busca na fila |
| `GET /api/playlist/current` | posição atual + janela em torno dela (para "Música atual") |
| `POST /api/player/toggle` \| `/shuffle` \| `/repeat` \| `/crossfade` | toggles (espelham Node) |
| `POST /api/player/play {position}` | toca posição específica |
| `POST /api/playlist/add {files[], position}` | adiciona acima/abaixo |
| `POST /api/playlist/add-many {file, timesPerDay, intervalHours, from?, to?}` | inserção distribuída |
| `POST /api/playlist/move {from, to}` | move faixa |
| `POST /api/playlist/remove {positions[]}` | remove faixas |
| `GET /api/playlists` | playlists salvas (nome, data; contagem se disponível) |
| `POST /api/playlists/{name}/load {playAtRandom?}` | carrega; se `playAtRandom`, toca posição aleatória (comportamento do botão WinForms) |
| `POST /api/playlists {name}` | salva fila atual |
| `DELETE /api/playlists/{name}` | remove playlist |

Erros: `{error: {code, message}}`; 401 sem sessão, 403 allowlist, 502 Node
indisponível, 400 validação.

## Frontend React (`frontend-web/`)

Vite + React 18 + TypeScript. Estado servidor: TanStack Query; listas: TanStack
Virtual com paginação infinita (páginas de ~100). Busca com debounce 300ms,
server-side. SSE via `EventSource` alimenta o cache do status e invalida a fila
pela versão.

Design system Morcegão FM portado como CSS variables (tokens extraídos do wireframe:
neutros quentes `--bat-*`, marca `--red-500 #E4261D`, destaque `--amber-500 #F2A007`,
info `--blue-500 #2A9DF4`, live `#34C759`; fontes Anton/Oswald/Barlow/Space Mono
via `@fontsource`; radii, sombras, glows e durações conforme tokens).

### Telas

- **Login** (`/login`): logo, botão "Entrar com Google"; variante de acesso negado
  mostrando o e-mail rejeitado.
- **Tocando Agora** (`/`): dois painéis.
  - *Acervo* (borda azul): busca, lista virtualizada (avatar-letra, título,
    artista·duração, chip de gênero), botões "↑ Adicionar acima"/"↓ Adicionar abaixo"
    da faixa selecionada na fila, menu "..." com "Inserir várias vezes".
  - *Fila de Hoje* (borda vermelha): contagem de faixas, busca na fila, botão
    "Música atual" (scroll até a atual), linhas com horário previsto (Space Mono) +
    posição, faixa atual com borda âmbar; ações por linha: subir/descer/tocar
    (dupla confirmação, como no WinForms)/remover.
  - *Player bar* fixo: play/pause (MPD), faixa atual, indicador AO VIVO,
    volume do monitor local (stream Icecast via `<audio>`, se `STREAM_URL` configurado).
- **Playlist e Arquivos** (`/playlists`): dropdown + "Carregar e tocar" (substitui a
  fila, toca de posição aleatória, com confirmação), cards das playlists salvas
  (nome, contagem·data, "Carregar esta"), salvar fila atual como playlist.
- **Configurações** (`/config`): toggles Tocar playlist / Shuffle / Repetir /
  Transição (crossfade); card servidor somente leitura + "Atualizar dados do servidor".
- **Modal Inserir várias vezes**: vezes por dia, intervalo (select), checkbox
  "Apenas por um período" com De/Até; confirma → `add-many`.

Cabeçalho comum: título da tela, relógio (hora · dia da semana), botão "Atualizar".
Sidebar: logo, badge AO VIVO + servidor, navegação (3 itens), rodapé versão.

## Tratamento de erros

- Node fora do ar: gateway responde 502; frontend exibe banner persistente
  "Servidor da rádio inacessível" e pausa o polling visível, sem derrubar sessão.
- Sessão expirada: 401 → redirect a `/login`.
- Mutações da fila retornam a fila atualizada (janela atual) ou disparam invalidação
  via SSE `playlist`.

## Testes

- **Go**: unit — allowlist/middleware de sessão, busca/paginação do índice,
  `schedule` (TimeCalculation com repeat on/off; distribuição do add-many),
  handlers com Node mockado (`httptest`).
- **Frontend**: vitest — hooks de busca/paginação e reducer de seleção; smoke de
  render das telas com MSW.

## Fora do escopo

Papéis/multiusuário além da allowlist, histórico de reprodução, upload/edição de
músicas, refresh token Google (sessão própria de 7 dias basta), edição do card de
servidor, controle de volume do MPD.
