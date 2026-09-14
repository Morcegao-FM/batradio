# BatRadio Web em produção: domínio próprio, tudo autenticado, capa vinda do website

Data: 2026-09-13

## Problema

O BatRadio Web (gateway Go + SPA React) está pronto no repositório e nunca foi
publicado. Queremos ele em `https://batradio.morcegaofm.com.br`, com login
Google restrito a `aguergolet@gmail.com` e `morcegaofm@gmail.com` — a mesma
allowlist da Batcaverna — e com os dados de exibição de faixa (capa, artista,
título, ano) vindos do catálogo que o website já mantém, em vez de cada cliente
resolver capa por conta própria.

Duas restrições moldam tudo:

1. **O painel é o controle da rádio no ar.** Quem entra nele muda o que está
   tocando. Nenhum endpoint pode ficar sem autenticação.
2. **O cliente Windows do Renato não pode parar.** Ele é o operador de hoje.
   O spec anterior (`2026-07-08-batradio-web-design.md`) declarava o cliente
   Windows aposentado e a compatibilidade abandonada. **Esta premissa está
   revogada**: o cliente continua em produção por tempo indeterminado, e o
   caminho que ele usa é intocável.

## Estado atual (levantado no código, 2026-09-13)

### O que já existe e funciona

- `backend-gateway/` — OAuth Google server-side com allowlist (`internal/auth/`),
  sessão assinada em cookie `HttpOnly`, e **todas** as rotas `/api/*` já atrás de
  `auth.RequireSession` (`internal/httpapi/api.go:66`). A allowlist é revalidada
  a cada request, então remover um e-mail da env revoga acesso mesmo com cookie
  válido.
- `.env.example` já lista os dois e-mails corretos.
- `Dockerfile.gateway` já builda frontend + Go e gera imagem distroless.
- Website: tabela `musica_catalogo` (`Morcegao.Domain/Entities/MusicaCatalogo.cs`)
  com arquivo → tipo/artista/título/ano/`ImagemUrl`/`NomeExibicao`, enriquecida
  pelo `MusicaCatalogoService` via `ITunesCapaClient` e corrigível na Batcaverna
  (`/api/batcaverna/musicas`, `[Authorize(Roles="admin")]`). O comentário da
  entidade já diz que ela existe para "o site e as demais plataformas" — o
  BatRadio é a plataforma que faltava.

### O que está errado ou perigoso hoje

- **`docker-compose.yml` sobe um `node-backend` próprio** e o tira da rede
  externa (`expose`, sem `ports`). Aplicado no servidor da rádio, isso cria um
  segundo cliente MPD e/ou tira a 9320 do ar — exatamente o cenário que derruba
  o Renato. Esse compose serve para dev/local e **não** para produção.
- **A porta 9320 está na internet em HTTP puro**, protegida só pelo header
  `batradio-apikey` (`backend-server/server.js:209`). O cliente Windows tem
  `http://` fixo em `BatRadio.UI/BatRadioClient.cs:124`, então não fala HTTPS sem
  recompilar. Decisão desta rodada: **não tocar**. Fica registrado como dívida
  conhecida, para resolver numa janela combinada com o Renato.
- `Secure` do cookie depende de `X-Forwarded-Proto` (`internal/auth/oauth.go:64`).
  Atrás do Traefik funciona, mas confia num header que vem de fora do processo.
- Sem rate limit no login, sem headers de segurança, sem log de auditoria de
  quem alterou a fila.

### `batbelt` ≠ node do BatRadio

São serviços diferentes no mesmo host, e confundi-los quebraria os dois:

| | batbelt | node do BatRadio |
|---|---|---|
| Porta | 8003 | 9320 |
| Auth | nenhuma | header `batradio-apikey` |
| Playlist | `GET /playlist` | `POST /playlist` |
| Consumidor | API do website | cliente Windows do Renato |

## Arquitetura proposta

```
Internet ─443─► Traefik (rede externa telogreug)
                  ├─ www.morcegaofm.com.br ──► web(nginx) ─► api .NET ─► MySQL
                  └─ batradio.morcegaofm.com.br ──► batradio-gateway (Go)
                                                      │
                              ┌───────────────────────┴────────────────────┐
                              ▼                                            ▼
                   node 9320 (EXISTENTE, intocado)          api .NET: POST /api/servico/musicas/lookup
                              │                                    (header X-Servico-Key)
                              ▼
                             MPD ──► Icecast

Renato (Windows) ──HTTP 9320──► o mesmo node de sempre  [caminho preservado byte a byte]
```

O gateway é o único componente novo em produção. Ele não publica porta no host:
sai o `ports: 8080:8080`, entra `expose` + labels do Traefik.

## Componentes

### 1. Compose de produção (`docker-compose.main.yml`, novo)

Declara **somente** o serviço `gateway`. Não declara, não referencia e não
constrói `node-backend`. `NODE_BACKEND_URL` aponta para o node que já roda no
host; o endereço exato (nome de container na `telogreug`, ou IP do host via
`extra_hosts`) é definido pelo levantamento descrito em "Verificações antes do
primeiro deploy".

Labels do Traefik espelhando o padrão do website: router em
`Host(batradio.morcegaofm.com.br)`, `entrypoints=websecure,web`, `tls=true`,
`certresolver=letsencrypt`, `loadbalancer.server.port=8080`,
`traefik.docker.network=telogreug`. Sem router de apex, sem redirect regex —
é um subdomínio só.

O `docker-compose.yml` da raiz continua como está, documentado no README como
ambiente de desenvolvimento. Ganha um comentário no topo dizendo em uma linha
por que ele não serve para produção.

### 2. Endurecimento do gateway

Mudanças em `backend-gateway/`:

- **Cookie**: `Secure` sempre verdadeiro quando `COOKIE_SECURE=true` (padrão em
  produção), sem consultar `X-Forwarded-Proto`. Nome do cookie de sessão ganha o
  prefixo `__Host-`, que exige `Secure`, `Path=/` e proíbe `Domain` — o
  navegador passa a recusar o cookie se alguma dessas condições cair.
- **Falha fechada no boot**: se `DEV_MODE=true` e `OAUTH_REDIRECT_URL` for
  `https://`, o processo aborta com erro claro em vez de subir um painel de
  rádio com login falso.
- **Headers de segurança** em middleware global: `Strict-Transport-Security`,
  `Content-Security-Policy` restritiva (SPA e assets são self-hosted; a única
  origem externa é a `STREAM_URL` do Icecast, em `media-src`),
  `X-Content-Type-Options: nosniff`, `Referrer-Policy: same-origin`,
  `X-Frame-Options: DENY`.
- **Rate limit** por IP em `/auth/login` e `/auth/callback`, com resposta 429.
  Não se aplica a `/api/*`, que já exige sessão válida.
- **Log de auditoria**: toda rota que muda estado (`/playlist/add`,
  `/playlist/add-many`, `/playlist/move`, `/playlist/remove`, `/player/*`,
  `/playlists*`) registra e-mail, ação e alvo. Uma linha por escrita, sem
  segredo e sem corpo de request.
- **`GET /healthz`**: pública, sem sessão, responde só com status do processo
  (sem versão de dependência, sem estado do MPD, sem contagem de acervo). É a
  única rota pública nova, e existe para o `up -d --wait` e o rollback.

### 3. Endpoint de serviço no website

`POST /api/servico/musicas/lookup` na API .NET.

- Corpo: `{"arquivos": ["Rock/ACDC/back-in-black.mp3", ...]}`, no máximo 200 por
  chamada; acima disso, 400.
- Resposta: lista de `{arquivo, tipo, artista, titulo, ano, imagemUrl,
  nomeExibicao}` **apenas para os arquivos que já existem** no catálogo.
  Arquivo desconhecido simplesmente não aparece na resposta — não é erro.
- **Somente leitura.** Não insere linha, não chama o iTunes. O acervo inteiro do
  MPD não deve virar linha na Batcaverna: o catálogo é dos que já tocaram.
- Autorização: policy `servico`, que lê o header `X-Servico-Key` e compara em
  tempo constante com `Servico__Chave`. **Sem chave configurada, a policy nega
  tudo** — mesmo padrão fail-closed já usado na validação de JWT.
- A policy `servico` não concede nada além deste endpoint. Não substitui nem
  encosta no `[Authorize(Roles="admin")]` da Batcaverna.
- Não expõe PII: só metadado de faixa.

### 4. Cliente de capa no gateway

Pacote novo `internal/catalogo`:

- Cache em memória com TTL (~1h) e **negative cache** (arquivo que o website não
  conhece não é reperguntado a cada render).
- Consulta em lote, só dos arquivos visíveis: tocando agora + janela da fila.
  O acervo de ~50 mil faixas nunca é enviado ao website.
- Timeout curto e **degradação silenciosa**: website indisponível, lento ou
  respondendo erro → a UI mostra placeholder de capa e o controle do MPD segue
  funcionando. O catálogo é enfeite; o painel é operação.
- A chave de serviço vem de env, nunca aparece em log nem em mensagem de erro
  devolvida ao browser.

### 5. Pipeline

`.github/workflows/ci.yml` — em todo push e PR: `go test ./...` em
`backend-gateway/`, `npm test` e `npm run build` em `frontend-web/`.

`.github/workflows/deploy.yml` — no push da `main`, depois do CI passar:

1. Build de `Dockerfile.gateway` (uma imagem só: o Dockerfile já embute o SPA).
2. Trivy na imagem; `CRITICAL`/`HIGH` bloqueia, sem bypass.
3. Push para `ghcr.io/Morcegao-FM/batradio-gateway:<sha>` — tag imutável, nunca
   `:latest`.
4. SSH pelo GitHub Environment `production` (mesmos secrets de host/usuário/chave
   do padrão já usado no website): `git pull` do compose, `docker compose pull`,
   `up -d --wait`.
5. Healthcheck falhou → rollback automático para a tag anterior.

**A primeira subida é manual e combinada com o Renato**, para confirmar com ele
na linha que a 9320 não piscou antes de o pipeline passar a rodar sozinho.

### 6. Skill `batradio-limites`

`.claude/skills/batradio-limites/SKILL.md` no repo do BatRadio, com as
limitações que não estão óbvias no código:

- O contrato do node 9320 (rotas, headers, formato) é a interface do cliente
  Windows. Não mudar, não renomear, não "arrumar".
- Nunca subir um segundo `node-backend` em produção; o compose de produção não
  declara esse serviço.
- `frontend-windows-client/` não é código morto. Não apagar, não modernizar sem
  janela combinada com o Renato.
- Toda rota nova em `/api` nasce atrás de `RequireSession`. Exceção só com teste
  que a documente; hoje a lista de exceções é `/healthz` e `/auth/*`.
- `DEV_MODE` jamais em produção.
- O BatRadio não escreve no banco do website — só lê pelo endpoint de serviço.
- Segredos (`NODE_API_KEY`, chave de serviço, `SESSION_SECRET`) nunca em log,
  mensagem de erro ou código.
- `batbelt:8003` e o node `:9320` são serviços diferentes (ver tabela acima).

## Testes

**Go (`backend-gateway/`)**
- Teste de tabela que varre as rotas registradas no chi router e **falha** se
  alguma rota fora da allowlist explícita (`/healthz`, `/auth/*`, estáticos)
  responder sem cookie de sessão. É a rede que pega rota nova desprotegida.
- Config: `DEV_MODE=true` + redirect `https://` → erro no boot.
- Cliente de capa contra servidor fake: 200, resposta parcial, 404, timeout e
  500 — todos devem degradar para placeholder sem derrubar a request do painel.
- Cache: negative cache não reconsulta; TTL expira e reconsulta.

**.NET (`Morcegao.Api.Tests/`)**
- `lookup` com chave correta: 200 e só os arquivos existentes.
- Chave errada, ausente ou **não configurada no servidor**: 401.
- Mais de 200 arquivos: 400.
- Um arquivo desconhecido não cria linha em `musica_catalogo` (garante que
  continua somente leitura).

**Fumaça pós-deploy** (roteiro no README, rodado à mão na primeira subida)
- `GET /healthz` → 200.
- `GET /api/status` sem cookie → 401.
- `GET http://<host>:9320/status` com o `batradio-apikey` do Renato → 200,
  provando que o caminho dele continua de pé.
- Login com um e-mail fora da allowlist → 403.

## Verificações antes do primeiro deploy

A serem rodadas por SSH no host de produção e registradas no plano de
implementação — o desenho do compose depende das respostas:

1. O node 9320 é container ou processo bare? Se container: nome, rede e se está
   na `telogreug`. Define se `NODE_BACKEND_URL` usa nome de serviço ou
   `extra_hosts` com o IP do host.
2. Nada no host já ocupa o que o compose novo vai querer (nome de container,
   rede, volume).
3. O `letsencrypt` do Traefik compartilhado emite certificado para subdomínio
   novo do mesmo jeito que emitiu para o `www`.
4. Qual `batradio-apikey` o cliente do Renato usa hoje, para o gateway usar a
   **mesma** e nada no node precisar mudar.

## Fora de escopo

- Tirar a 9320 da internet, pôr TLS nela ou recompilar o cliente Windows.
  Registrado como dívida; exige janela com o Renato.
- Enriquecimento de capa sob demanda a partir do acervo do BatRadio.
- Aposentar o cliente Windows.
- Papéis ou multiusuário além da allowlist de dois e-mails.
- Ambiente de homologação do BatRadio.
