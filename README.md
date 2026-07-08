# Bat Radio by Morcegão FM

Gerenciador de programação para web rádios baseadas em MPD + Icecast 2, criado
para a [Morcegão FM](http://www.morcegaofm.com.br) — só rock clássico e blues.

## Arquitetura

```
Browser (React SPA)
   │  HTTPS + cookie de sessão (login Google)
   ▼
Gateway Go (backend-gateway/)  ──►  Backend Node (backend-server/)  ──►  MPD ──► Icecast
 • OAuth Google + allowlist          • único processo que fala MPD
 • API REST nova + SSE               • rede interna, protegido por API key
 • busca/paginação do acervo
 • serve o build do React (embed)
```

- **`backend-gateway/`** — gateway em Go: autenticação Google (allowlist de
  e-mails), índice em memória do acervo (busca e paginação server-side — o
  acervo nunca vai inteiro ao browser), cálculo dos horários da fila, SSE de
  status e o build do React embutido no binário.
- **`frontend-web/`** — SPA em React + TypeScript com o design system Morcegão
  FM: Tocando Agora (acervo + fila virtualizadas), Playlist e Arquivos,
  Configurações e o modal "Inserir música várias vezes".
- **`backend-server/`** — backend Node legado que conversa com o MPD. Continua
  como está, mas passa a rodar apenas na rede interna do Docker.
- **`frontend-windows-client/`** — cliente Windows Forms **descontinuado**,
  substituído pelo frontend web.

## Rodando em produção (Docker Compose)

1. **Credenciais Google OAuth** — no
   [Google Cloud Console](https://console.cloud.google.com/apis/credentials):
   crie um *OAuth client ID* do tipo **Web application**, com
   *Authorized redirect URI* = `https://SEU_DOMINIO/auth/callback`.
2. **Configuração do gateway** — copie `backend-gateway/.env.example` para
   `.env` na raiz e preencha (client id/secret, `ALLOWED_EMAILS`,
   `SESSION_SECRET` com `openssl rand -hex 32`, `STREAM_URL` do Icecast).
3. **Configuração do Node** — copie `backend-server/config.template.json` para
   `backend-server/config.json` (endereço do MPD e `apiKey` igual ao
   `NODE_API_KEY` do `.env`).
4. Suba: `docker compose up -d --build`.
5. Aponte seu proxy reverso com HTTPS (Caddy, Nginx, Traefik) para a porta
   `8080`. O domínio precisa ser o mesmo do `OAUTH_REDIRECT_URL`.

Somente os e-mails em `ALLOWED_EMAILS` conseguem entrar (hoje:
`aguergolet@gmail.com` e `morcegaofm@gmail.com`).

## Desenvolvimento local

Sem MPD por perto? Use o backend fake com 50 mil faixas sintéticas:

```bash
# 1. Backend Node fake (porta 9320, apikey "dev")
node backend-gateway/hack/fake-node.mjs

# 2. Gateway em modo dev (login sem Google, usuário dev@localhost)
cd backend-gateway
PORT=8090 NODE_BACKEND_URL=http://127.0.0.1:9320 NODE_API_KEY=dev \
DEV_MODE=true ALLOWED_EMAILS=aguergolet@gmail.com \
SESSION_SECRET=$(openssl rand -hex 32) go run .
# (ou crie um .env — no diretório atual ou na raiz do repo — que o gateway
#  carrega sozinho; variáveis já exportadas têm precedência)

# 3. Frontend com hot reload (proxy /api e /auth → :8090)
cd frontend-web && npm install && npm run dev
```

Testes: `go test ./...` em `backend-gateway/`; `npm test` em `frontend-web/`.

## API do gateway (resumo)

Tudo em `/api`, JSON, autenticado por cookie de sessão. Erros:
`{"error":{"code","message"}}`; `401` sem sessão, `403` e-mail fora da
allowlist, `502` Node/MPD indisponível.

| Rota | Descrição |
|---|---|
| `GET /api/status` · `GET /api/events` (SSE) | status do player (poll único no Node, fan-out) |
| `GET /api/library?q=&offset=&limit=` | acervo paginado; busca por termos (todos precisam bater) |
| `POST /api/library/refresh` · `GET /api/server/info` | re-indexa o acervo · diagnóstico |
| `GET /api/playlist?q=&offset=&limit=` | fila paginada com horário previsto de cada faixa |
| `POST /api/player/{toggle,shuffle,repeat,crossfade,play}` | controles |
| `POST /api/playlist/{add,add-many,move,remove}` | mutações da fila |
| `GET/POST /api/playlists` · `POST /api/playlists/{name}/load` · `DELETE /api/playlists/{name}` | playlists salvas |

## Histórico — Windows Frontend (descontinuado)

O cliente Windows (v1.2.0.x) foi substituído pelo frontend web. O código segue
em `frontend-windows-client/` como referência.
