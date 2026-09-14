# Deploy do BatRadio Web

## O que este deploy NÃO faz

Não encosta no backend Node da porta 9320. Ele já roda neste host e é a
interface do cliente Windows do Renato. `docker-compose.main.yml` declara
somente o gateway, de propósito. Ver `.claude/skills/batradio-limites/`.

## Antes da primeira subida

Levantar no servidor, por SSH:

```bash
# 1. Como o Node roda hoje: container ou processo bare?
docker ps --format '{{.Names}}\t{{.Ports}}' | grep -i 9320 || ss -lntp | grep 9320

# 2. Se for container, em que redes ele está:
docker inspect -f '{{range $k,$v := .NetworkSettings.Networks}}{{$k}} {{end}}' <nome>

# 3. A rede do Traefik existe:
docker network ls | grep telogreug

# 4. Nada já ocupa o nome do serviço novo:
docker ps -a --format '{{.Names}}' | grep -i batradio
```

`NODE_BACKEND_URL` sai daí:

- Node em container na `telogreug` → `http://<nome-do-container>:9320`.
- Node como processo no host → `http://host.docker.internal:9320`, e acrescentar
  ao serviço `gateway` em `docker-compose.main.yml`:
  ```yaml
      extra_hosts:
        - "host.docker.internal:host-gateway"
  ```

`NODE_API_KEY` tem de ser **a mesma chave que o cliente do Renato usa**
(`backend-server/config.json` no servidor, campo `webserver.apiKey`). Assim nada
no Node precisa mudar.

## Google OAuth

No Google Cloud Console → APIs & Services → Credentials, no OAuth client ID do
tipo Web application, acrescentar em *Authorized redirect URIs*:

```
https://batradio.morcegaofm.com.br/auth/callback
```

## DNS

`batradio.morcegaofm.com.br` → IP do servidor. O certresolver `letsencrypt` do
Traefik emite o certificado sozinho quando o router sobe.

## Chave de serviço do catálogo

A mesma string nos dois lados:

- Website: `SERVICO_CHAVE` no `.env` de `~/apps/morcegao-fm-main`, que o
  `docker-compose.main.yml` de lá injeta como `Servico__Chave`.
- BatRadio: `CATALOGO_CHAVE` no `.env` de `~/apps/batradio`.

Gerar com `openssl rand -hex 32`. Menos de 32 caracteres e a API recusa
autenticar — de propósito.

## Primeira subida (manual, com o Renato avisado)

```bash
mkdir -p ~/apps/batradio && cd ~/apps/batradio
git clone -b main git@github.com:Morcegao-FM/batradio.git .
cp backend-gateway/.env.example .env   # preencher com os valores reais
echo "TAG=<sha do commit>" >> .env
docker compose -f docker-compose.main.yml up -d --wait
```

## Fumaça (rodar à mão na primeira subida)

```bash
# Gateway de pé
curl -fsS https://batradio.morcegaofm.com.br/healthz

# API fechada sem sessão
curl -s -o /dev/null -w '%{http_code}\n' https://batradio.morcegaofm.com.br/api/status
# esperado: 401

# O caminho do Renato continua de pé — o teste que mais importa
curl -s -o /dev/null -w '%{http_code}\n' -H "batradio-apikey: <a chave dele>" \
  http://<host>:9320/status
# esperado: 200
```

Depois disso, avisar o Renato para conferir o cliente dele antes de o pipeline
passar a fazer deploy sozinho.

## Secrets do GitHub

Em Settings → Environments → `production` (mesmo padrão do website):

| Secret | Valor |
|---|---|
| `DEPLOY_SSH_HOST` | IP do servidor do playout |
| `DEPLOY_SSH_USER` | usuário SSH |
| `DEPLOY_SSH_KEY` | chave privada dedicada ao Actions (não a pessoal) |
| `DEPLOY_SSH_PORT` | porta SSH, se não for 22 (opcional) |

A imagem vai para o GHCR com o `GITHUB_TOKEN` do próprio workflow — nenhuma
credencial de registry precisa ser criada.
