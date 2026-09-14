---
name: batradio-limites
description: Use ao mexer em QUALQUER coisa do BatRadio — gateway Go, SPA, compose, deploy, backend Node ou cliente Windows. Define o que não pode ser tocado (a porta 9320 do operador), a regra de autenticação de toda rota nova e a fronteira com o website. Dispara em "batradio", "gateway", "9320", "node-backend", "MPD", "deploy", "compose", "Renato", "cliente Windows", "capa", "catálogo", "batbelt".
---

# Limites do BatRadio

Três regras que o código não conta sozinho, e cujo custo de violar é a rádio
fora do ar ou o painel de controle exposto.

## 1. A porta 9320 é de outra pessoa

O backend Node (`backend-server/`) roda em produção no host do playout, na porta
9320, autenticado pelo header `batradio-apikey` em HTTP puro. **O cliente
Windows do operador depende dele**, e o cliente tem `http://` fixo no código
(`frontend-windows-client/BatRadio.UI/BatRadioClient.cs:124`) — ou seja, não
fala HTTPS sem ser recompilado e reinstalado na máquina dele.

Portanto:

- **Nunca** mude rota, header ou formato de resposta de `backend-server/`.
- **Nunca** declare um serviço `node-backend` no compose de produção. Dois
  processos falando com o mesmo MPD quebram a operação. `docker-compose.yml` é
  de desenvolvimento; produção é `docker-compose.main.yml`, só com o gateway.
- **Nunca** mude a porta, a chave ou a exposição da 9320 sem janela combinada
  com o operador. O gateway usa **a mesma** `batradio-apikey` que ele usa.
- `frontend-windows-client/` **não é código morto.** Não apagar, não
  "modernizar", não deixar o Dependabot mudá-lo sozinho.

O spec de 2026-07-08 dizia que o cliente Windows estava aposentado. **Isso está
revogado** — ver `docs/superpowers/specs/2026-09-13-batradio-deploy-seguro-design.md`.

Que a 9320 esteja na internet em HTTP é dívida conhecida, registrada no mesmo
spec. Resolver isso é um projeto com o operador na sala, não um refactor de
passagem.

## 2. Toda rota nasce autenticada

O BatRadio controla o que está no ar. Quem entra nele muda a transmissão.

- Toda rota nova em `/api` vai atrás de `auth.RequireSession`. Sem exceção
  "temporária".
- A lista de rotas públicas é fechada: `/healthz`, `/auth/login`,
  `/auth/callback`, `/auth/logout` e os estáticos do SPA. Ela vive em
  `internal/httpapi/superficie_test.go`, e `TestNenhumaRotaNovaSemSessao`
  percorre o router do chi e falha se alguém acrescentar rota desprotegida.
  **Se esse teste ficar vermelho, proteja a rota — não relaxe a allowlist.**
- `/healthz` é público e deliberadamente vazio (`{"status":"ok"}`). Não
  acrescente versão, host do Node nem contagem de acervo nele; há teste que
  barra isso.
- `DEV_MODE=true` é login falso, assinando sessão como `dev@localhost`. O
  `config.Load` recusa subir se ele vier com `OAUTH_REDIRECT_URL` https. Não
  contorne.
- Acesso é a allowlist de `ALLOWED_EMAILS`, revalidada a cada request. Não
  existe papel, convite nem cadastro.
- Escrita é auditada (`internal/httpapi/auditoria.go`): uma linha com o e-mail
  por ação que muda estado. Rota de escrita nova entra nesse caminho de graça,
  desde que esteja sob `/api`.

## 3. O website é dono dos dados de faixa; o BatRadio só lê

Capa, artista, título e ano vivem na tabela `musica_catalogo` da API .NET, e são
corrigidos à mão na Batcaverna.

- O BatRadio lê por `POST /api/servico/musicas/lookup` e corrige por
  `PUT /api/servico/musicas/por-arquivo`, ambos com header `X-Servico-Key`.
- **Corrige linha que existe; nunca cria.** Faixa sem linha no catálogo devolve
  404. O acervo do MPD tem dezenas de milhares de faixas e não pode entrar no
  catálogo por ação do BatRadio — quem cria linha é o sync do próprio site, a
  partir do que tocou.
- **Toda escrita carrega autor.** O gateway manda o e-mail da sessão como
  `editadoPor`, e o website grava em `musica_catalogo.editado_por`. O autor vem
  SEMPRE da sessão, nunca do corpo do request: é a auditoria que sustenta o
  website aceitar escrita por chave de serviço. Se você acrescentar outra rota
  de escrita, ela carrega autor também.
- **Editar trava o enriquecimento** daquela faixa (`editado_em`), senão o iTunes
  desfaria a correção no ciclo seguinte. Consequência: edição errada não se
  conserta sozinha — só pelo botão Restaurar da Batcaverna.
- O BatRadio não chama o iTunes e não escreve em nenhuma outra tabela.
- Capa é enfeite; o painel é operação. `internal/catalogo` nunca devolve erro:
  site fora do ar vira placeholder, e o controle do MPD segue.
- Só a faixa atual, a próxima e a página visível da fila são enriquecidas. Não
  enriqueça `/api/library`: são dezenas de milhares de faixas e quase nenhuma
  tem linha no catálogo.

## Confusão que já custou tempo

`batbelt` e o node do BatRadio são **serviços diferentes** no mesmo host:

| | batbelt | node do BatRadio |
|---|---|---|
| Porta | 8003 | 9320 |
| Auth | nenhuma | header `batradio-apikey` |
| Playlist | `GET /playlist` | `POST /playlist` |
| Consumidor | API do website | cliente Windows do operador |

## Segredos

`SESSION_SECRET`, `NODE_API_KEY`, `CATALOGO_CHAVE` e `Servico:Chave` vêm de env.
Nunca em código, em log, em teste ou em mensagem de erro devolvida ao browser —
inclusive na hora de registrar falha de autenticação.
