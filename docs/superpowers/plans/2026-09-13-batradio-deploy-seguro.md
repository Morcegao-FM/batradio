# BatRadio em produção: domínio próprio, tudo autenticado, capa do website — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Publicar o BatRadio Web em `https://batradio.morcegaofm.com.br` com login Google restrito a dois e-mails, toda rota `/api` autenticada, capa vinda do catálogo do website — sem tocar no caminho que o cliente Windows do Renato usa.

**Architecture:** O gateway Go entra como mais um serviço atrás do Traefik compartilhado da rede `telogreug`, no mesmo host do playout. Ele fala com o backend Node **que já roda** (porta 9320, intocado) e com a API .NET do website por um endpoint de serviço somente leitura, autenticado por chave dedicada. O compose de produção declara apenas o gateway.

**Tech Stack:** Go 1.25 + chi v5 (gateway), React 18 + TypeScript + Vite (SPA), .NET 10 + EF Core (API do website), Docker Compose + Traefik, GitHub Actions + GHCR + Trivy.

## Global Constraints

- **Dois repositórios.** Cada tarefa diz `[W]` para `~/code/morcegaofm/morcegao-fm-website` ou `[B]` para `~/code/morcegaofm/batradio`. Commits vão no repo da tarefa.
- **O contrato do node 9320 é imutável.** Nenhuma tarefa altera `backend-server/` ou `frontend-windows-client/`.
- **O compose de produção do BatRadio nunca declara `node-backend`.**
- Rotas públicas do gateway, lista completa e fechada: `/healthz`, `/auth/login`, `/auth/callback`, `/auth/logout` e os estáticos do SPA. Qualquer outra rota exige sessão.
- Allowlist de produção: `aguergolet@gmail.com,morcegaofm@gmail.com`.
- Domínio: `batradio.morcegaofm.com.br`. Redirect OAuth: `https://batradio.morcegaofm.com.br/auth/callback`.
- Segredos (`SESSION_SECRET`, `NODE_API_KEY`, `Servico:Chave`) nunca em código, log, teste ou mensagem de erro devolvida ao browser.
- Mensagens de erro visíveis ao usuário em português; comentários e commits em português, seguindo o estilo dos repositórios.
- Go: `go test ./...` em `backend-gateway/`. Frontend: `npm test` em `frontend-web/`. Website: `dotnet test backend/Morcegao.sln`.

---

### Task 1 [W]: Endpoint de serviço para lookup de capa

**Files:**
- Create: `backend/Morcegao.Api/Auth/ServicoAuthenticationHandler.cs`
- Create: `backend/Morcegao.Application/Radio/Catalogo/MusicaLookupDtos.cs`
- Create: `backend/Morcegao.Api/Controllers/ServicoMusicasController.cs`
- Modify: `backend/Morcegao.Api/Program.cs` (registro do esquema, junto ao `AddJwtBearer`)
- Modify: `backend/Morcegao.Api.Tests/ApiTestFactory.cs` (nova propriedade `ServicoChave`)
- Test: `backend/Morcegao.Api.Tests/ServicoMusicasLookupTests.cs`

**Interfaces:**
- Consumes: `IAppDbContext.MusicasCatalogo` (`DbSet<MusicaCatalogo>`); `MusicaCatalogoService.Hash(string arquivo)` — `public static string`, SHA-256 hex minúsculo.
- Produces: `POST /api/servico/musicas/lookup`, header `X-Servico-Key`. Corpo `{"arquivos":["a.mp3"]}`. Resposta `200` com lista de `{arquivo,tipo,artista,titulo,ano,imagemUrl,nomeExibicao}`. É o contrato que a Task 7 consome do lado Go.

- [ ] **Step 1: Escrever os testes que falham**

Criar `backend/Morcegao.Api.Tests/ServicoMusicasLookupTests.cs`:

```csharp
using System.Net;
using System.Net.Http.Json;
using Morcegao.Application.Radio.Catalogo;
using Morcegao.Domain.Entities;

namespace Morcegao.Api.Tests;

public class ServicoMusicasLookupTests
{
    private const string Chave = "chave-de-servico-de-teste-com-32-caracteres";

    private static ApiTestFactory Factory(string? chave = Chave) => new()
    {
        ServicoChave = chave,
        Seed = db =>
        {
            db.MusicasCatalogo.Add(new MusicaCatalogo
            {
                Arquivo = "Rock/ACDC/back-in-black.mp3",
                ArquivoHash = MusicaCatalogoService.Hash("Rock/ACDC/back-in-black.mp3"),
                Tipo = "musica",
                Artista = "AC/DC",
                Titulo = "Back in Black",
                Ano = 1980,
                ImagemUrl = "https://capa/512.jpg",
            });
        },
    };

    private static HttpClient ComChave(ApiTestFactory factory, string? chave)
    {
        var client = factory.CreateClient();
        if (chave is not null)
        {
            client.DefaultRequestHeaders.Add("X-Servico-Key", chave);
        }

        return client;
    }

    [Fact]
    public async Task Lookup_ComChaveCorreta_DevolveSomenteOsArquivosConhecidos()
    {
        using var factory = Factory();
        var client = ComChave(factory, Chave);

        var resposta = await client.PostAsJsonAsync("/api/servico/musicas/lookup", new
        {
            arquivos = new[] { "Rock/ACDC/back-in-black.mp3", "Vinhetas/nao-existe.mp3" },
        });

        Assert.Equal(HttpStatusCode.OK, resposta.StatusCode);
        var linhas = await resposta.Content.ReadFromJsonAsync<List<MusicaLookupDto>>();
        var linha = Assert.Single(linhas!);
        Assert.Equal("Rock/ACDC/back-in-black.mp3", linha.Arquivo);
        Assert.Equal("AC/DC", linha.Artista);
        Assert.Equal("https://capa/512.jpg", linha.ImagemUrl);
        Assert.Equal(1980, linha.Ano);
    }

    [Fact]
    public async Task Lookup_SemChaveNoHeader_Retorna401()
    {
        using var factory = Factory();
        var client = ComChave(factory, null);

        var resposta = await client.PostAsJsonAsync("/api/servico/musicas/lookup", new
        {
            arquivos = new[] { "Rock/ACDC/back-in-black.mp3" },
        });

        Assert.Equal(HttpStatusCode.Unauthorized, resposta.StatusCode);
    }

    [Fact]
    public async Task Lookup_ComChaveErrada_Retorna401()
    {
        using var factory = Factory();
        var client = ComChave(factory, "chave-errada-porem-com-32-caracteres-ok");

        var resposta = await client.PostAsJsonAsync("/api/servico/musicas/lookup", new
        {
            arquivos = new[] { "Rock/ACDC/back-in-black.mp3" },
        });

        Assert.Equal(HttpStatusCode.Unauthorized, resposta.StatusCode);
    }

    [Fact]
    public async Task Lookup_SemChaveConfiguradaNoServidor_Retorna401MesmoComHeader()
    {
        // Fail-closed: servidor sem Servico:Chave não autentica ninguém.
        using var factory = Factory(chave: null);
        var client = ComChave(factory, Chave);

        var resposta = await client.PostAsJsonAsync("/api/servico/musicas/lookup", new
        {
            arquivos = new[] { "Rock/ACDC/back-in-black.mp3" },
        });

        Assert.Equal(HttpStatusCode.Unauthorized, resposta.StatusCode);
    }

    [Fact]
    public async Task Lookup_AcimaDoTeto_Retorna400()
    {
        using var factory = Factory();
        var client = ComChave(factory, Chave);

        var arquivos = Enumerable.Range(0, 201).Select(i => $"faixa-{i}.mp3").ToArray();
        var resposta = await client.PostAsJsonAsync("/api/servico/musicas/lookup", new { arquivos });

        Assert.Equal(HttpStatusCode.BadRequest, resposta.StatusCode);
    }

    [Fact]
    public async Task Lookup_ArquivoDesconhecido_NaoCriaLinhaNoCatalogo()
    {
        // Garante que o endpoint é somente leitura: o acervo do MPD não pode
        // virar linha em musica_catalogo por consulta do BatRadio.
        using var factory = Factory();
        var client = ComChave(factory, Chave);

        await client.PostAsJsonAsync("/api/servico/musicas/lookup", new
        {
            arquivos = new[] { "Vinhetas/nao-existe.mp3" },
        });

        using var escopo = factory.Services.CreateScope();
        var db = escopo.ServiceProvider.GetRequiredService<IAppDbContext>();
        Assert.Equal(1, await db.MusicasCatalogo.CountAsync());
    }
}
```

Acrescentar no topo do arquivo os `using` que o projeto de teste já usa para escopo e EF: `using Microsoft.EntityFrameworkCore;`, `using Microsoft.Extensions.DependencyInjection;`, `using Morcegao.Application.Common;`.

- [ ] **Step 2: Rodar e ver falhar**

Run: `dotnet test backend/Morcegao.sln --filter ServicoMusicasLookupTests`
Expected: FAIL na compilação — `MusicaLookupDto` e `ServicoChave` não existem.

- [ ] **Step 3: Criar os DTOs**

`backend/Morcegao.Application/Radio/Catalogo/MusicaLookupDtos.cs`:

```csharp
namespace Morcegao.Application.Radio.Catalogo;

/// <summary>Corpo de POST /api/servico/musicas/lookup.</summary>
public sealed record MusicaLookupRequest(IReadOnlyList<string>? Arquivos);

/// <summary>
/// O que o catálogo sabe de uma faixa, para um consumidor de serviço (hoje, o
/// gateway do BatRadio). Sem PII: só metadado de faixa.
/// </summary>
public sealed record MusicaLookupDto(
    string Arquivo,
    string Tipo,
    string? Artista,
    string? Titulo,
    int? Ano,
    string? ImagemUrl,
    string? NomeExibicao);
```

- [ ] **Step 4: Criar o esquema de autenticação de serviço**

`backend/Morcegao.Api/Auth/ServicoAuthenticationHandler.cs`:

```csharp
using System.Security.Claims;
using System.Security.Cryptography;
using System.Text;
using System.Text.Encodings.Web;
using Microsoft.AspNetCore.Authentication;

namespace Morcegao.Api.Auth;

/// <summary>
/// Autentica chamadas máquina-a-máquina pelo header X-Servico-Key.
///
/// Fail-closed, igual à validação de JWT: sem `Servico:Chave` configurada (ou
/// com menos de 32 caracteres), nada autentica — nem em Development. Uma chave
/// placeholder no código-fonte seria uma porta aberta, então não existe.
/// </summary>
public sealed class ServicoAuthenticationHandler : AuthenticationHandler<AuthenticationSchemeOptions>
{
    public const string Scheme = "Servico";
    public const string Header = "X-Servico-Key";

    private readonly byte[] _chave;

    public ServicoAuthenticationHandler(
        IOptionsMonitor<AuthenticationSchemeOptions> options,
        ILoggerFactory logger,
        UrlEncoder encoder,
        IConfiguration configuration)
        : base(options, logger, encoder)
    {
        var chave = configuration["Servico:Chave"];
        _chave = string.IsNullOrEmpty(chave) || chave.Length < 32
            ? []
            : Encoding.UTF8.GetBytes(chave);
    }

    protected override Task<AuthenticateResult> HandleAuthenticateAsync()
    {
        if (_chave.Length == 0)
        {
            return Task.FromResult(AuthenticateResult.NoResult());
        }

        if (!Request.Headers.TryGetValue(Header, out var valores))
        {
            return Task.FromResult(AuthenticateResult.NoResult());
        }

        var enviada = Encoding.UTF8.GetBytes(valores.ToString());

        // FixedTimeEquals não vaza em qual caractere a comparação falhou.
        if (!CryptographicOperations.FixedTimeEquals(enviada, _chave))
        {
            // Sem ecoar a chave enviada: ela não vai para o log.
            Logger.LogWarning("Chave de serviço inválida em {Path}", Request.Path);
            return Task.FromResult(AuthenticateResult.Fail("Chave de serviço inválida."));
        }

        var identidade = new ClaimsIdentity([new Claim(ClaimTypes.Name, "servico")], Scheme);
        return Task.FromResult(AuthenticateResult.Success(
            new AuthenticationTicket(new ClaimsPrincipal(identidade), Scheme)));
    }
}
```

- [ ] **Step 5: Registrar o esquema no Program.cs**

Em `backend/Morcegao.Api/Program.cs`, logo depois do bloco `.AddJwtBearer(options => { ... })` (que termina por volta da linha 112, antes de `builder.Services.AddAuthorization();`), encadear:

```csharp
    .AddScheme<AuthenticationSchemeOptions, ServicoAuthenticationHandler>(
        ServicoAuthenticationHandler.Scheme, _ => { });
```

E acrescentar os `using` no topo: `using Morcegao.Api.Auth;` e `using Microsoft.AspNetCore.Authentication;`.

- [ ] **Step 6: Criar o controller**

`backend/Morcegao.Api/Controllers/ServicoMusicasController.cs`:

```csharp
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Morcegao.Api.Auth;
using Morcegao.Application.Common;
using Morcegao.Application.Radio.Catalogo;

namespace Morcegao.Api.Controllers;

/// <summary>
/// Leitura do catálogo de faixas para consumidores de serviço (hoje só o
/// gateway do BatRadio). Somente leitura de propósito: o acervo do MPD tem
/// dezenas de milhares de faixas e não pode virar linha em musica_catalogo por
/// consulta — o catálogo é das que já tocaram.
/// </summary>
[ApiController]
[Route("api/servico/musicas")]
[Authorize(AuthenticationSchemes = ServicoAuthenticationHandler.Scheme)]
public class ServicoMusicasController(IAppDbContext db) : ControllerBase
{
    /// <summary>Teto por chamada. Acima disso, 400 — evita varredura do acervo inteiro.</summary>
    public const int MaxArquivos = 200;

    [HttpPost("lookup")]
    public async Task<ActionResult<IReadOnlyList<MusicaLookupDto>>> Lookup(
        [FromBody] MusicaLookupRequest requisicao,
        CancellationToken ct)
    {
        var arquivos = (requisicao.Arquivos ?? [])
            .Where(a => !string.IsNullOrWhiteSpace(a))
            .Distinct(StringComparer.Ordinal)
            .ToList();

        if (arquivos.Count > MaxArquivos)
        {
            return BadRequest(new { erro = $"No máximo {MaxArquivos} arquivos por chamada." });
        }

        if (arquivos.Count == 0)
        {
            return Ok(Array.Empty<MusicaLookupDto>());
        }

        var porHash = arquivos.ToDictionary(MusicaCatalogoService.Hash, a => a, StringComparer.Ordinal);
        var hashes = porHash.Keys.ToList();

        var linhas = await db.MusicasCatalogo
            .AsNoTracking()
            .Where(m => hashes.Contains(m.ArquivoHash))
            .ToListAsync(ct);

        return Ok(linhas
            .Select(m => new MusicaLookupDto(
                porHash[m.ArquivoHash],
                m.Tipo,
                m.Artista,
                m.Titulo,
                m.Ano,
                m.ImagemUrl,
                m.NomeExibicao))
            .ToList());
    }
}
```

- [ ] **Step 7: Expor `ServicoChave` no ApiTestFactory**

Em `backend/Morcegao.Api.Tests/ApiTestFactory.cs`, junto às outras propriedades públicas (perto de `AgoraFixo`):

```csharp
    /// <summary>Chave do esquema Servico. Null = servidor sem chave (fail-closed).</summary>
    public string? ServicoChave { get; set; }
```

E dentro de `ConfigureAppConfiguration`, junto aos outros `if` condicionais (perto de `EnvioAudioHabilitado`):

```csharp
            if (ServicoChave is not null)
            {
                settings["Servico:Chave"] = ServicoChave;
            }
```

- [ ] **Step 8: Rodar os testes**

Run: `dotnet test backend/Morcegao.sln --filter ServicoMusicasLookupTests`
Expected: PASS, 6 testes.

- [ ] **Step 9: Rodar a suíte inteira**

Run: `dotnet test backend/Morcegao.sln`
Expected: PASS — o esquema novo não pode ter mexido em nada existente.

- [ ] **Step 10: Documentar a variável e commitar**

Em `.env.example`, na seção de segredos:

```
# Chave do endpoint de serviço /api/servico/musicas/lookup, consumido pelo
# gateway do BatRadio. >= 32 caracteres; sem ela o endpoint nega tudo.
SERVICO_CHAVE=
```

Em `docker-compose.main.yml`, dentro de `services.api.environment`:

```yaml
      Servico__Chave: "${SERVICO_CHAVE}"
```

```bash
git add backend/Morcegao.Api/Auth/ServicoAuthenticationHandler.cs \
        backend/Morcegao.Api/Controllers/ServicoMusicasController.cs \
        backend/Morcegao.Application/Radio/Catalogo/MusicaLookupDtos.cs \
        backend/Morcegao.Api/Program.cs \
        backend/Morcegao.Api.Tests/ApiTestFactory.cs \
        backend/Morcegao.Api.Tests/ServicoMusicasLookupTests.cs \
        .env.example docker-compose.main.yml
git commit -m "Serve o catálogo de faixas ao BatRadio por chave de serviço"
```

---

### Task 2 [B]: Config falha fechada — DEV_MODE e cookie seguro

**Files:**
- Modify: `backend-gateway/internal/config/config.go`
- Test: `backend-gateway/internal/config/config_test.go`

**Interfaces:**
- Produces: `config.Config.CookieSecure bool` — consumido pelas Tasks 3 e 4.

- [ ] **Step 1: Escrever os testes que falham**

Acrescentar em `backend-gateway/internal/config/config_test.go`:

```go
func TestLoadRecusaDevModeComRedirectHTTPS(t *testing.T) {
	env := map[string]string{
		"NODE_BACKEND_URL":   "http://node:9320",
		"NODE_API_KEY":       "k",
		"SESSION_SECRET":     strings.Repeat("s", 32),
		"ALLOWED_EMAILS":     "a@b.com",
		"DEV_MODE":           "true",
		"OAUTH_REDIRECT_URL": "https://batradio.morcegaofm.com.br/auth/callback",
	}
	_, err := Load(func(k string) string { return env[k] })
	if err == nil {
		t.Fatal("esperava erro: DEV_MODE com redirect https é login falso em produção")
	}
	if !strings.Contains(err.Error(), "DEV_MODE") {
		t.Fatalf("erro deveria citar DEV_MODE, veio: %v", err)
	}
}

func TestCookieSecurePadraoPorModo(t *testing.T) {
	base := map[string]string{
		"NODE_BACKEND_URL":     "http://node:9320",
		"NODE_API_KEY":         "k",
		"SESSION_SECRET":       strings.Repeat("s", 32),
		"ALLOWED_EMAILS":       "a@b.com",
		"GOOGLE_CLIENT_ID":     "id",
		"GOOGLE_CLIENT_SECRET": "secret",
		"OAUTH_REDIRECT_URL":   "https://batradio.morcegaofm.com.br/auth/callback",
	}

	prod, err := Load(func(k string) string { return base[k] })
	if err != nil {
		t.Fatalf("produção: %v", err)
	}
	if !prod.CookieSecure {
		t.Error("fora do DEV_MODE o cookie tem de ser Secure por padrão")
	}

	dev := map[string]string{}
	for k, v := range base {
		dev[k] = v
	}
	dev["DEV_MODE"] = "true"
	dev["OAUTH_REDIRECT_URL"] = ""
	dc, err := Load(func(k string) string { return dev[k] })
	if err != nil {
		t.Fatalf("dev: %v", err)
	}
	if dc.CookieSecure {
		t.Error("em DEV_MODE (http://localhost) o cookie Secure impediria o login")
	}

	forcado := map[string]string{}
	for k, v := range base {
		forcado[k] = v
	}
	forcado["COOKIE_SECURE"] = "false"
	fc, err := Load(func(k string) string { return forcado[k] })
	if err != nil {
		t.Fatalf("forçado: %v", err)
	}
	if fc.CookieSecure {
		t.Error("COOKIE_SECURE=false deve vencer o padrão")
	}
}
```

Garantir que `strings` está nos imports do arquivo de teste.

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd backend-gateway && go test ./internal/config/ -run 'DevMode|CookieSecure' -v`
Expected: FAIL — `prod.CookieSecure undefined`.

- [ ] **Step 3: Implementar**

Em `internal/config/config.go`, acrescentar o campo na struct, depois de `DevMode`:

```go
	// CookieSecure marca o cookie de sessão como Secure. Padrão: ligado fora do
	// DEV_MODE. COOKIE_SECURE=true/false força explicitamente.
	CookieSecure bool
```

E no fim de `Load`, imediatamente antes do `return c, nil`:

```go
	if c.DevMode && strings.HasPrefix(strings.ToLower(c.OAuthRedirectURL), "https://") {
		return nil, fmt.Errorf(
			"DEV_MODE=true com OAUTH_REDIRECT_URL https:// — isso subiria um painel " +
				"de rádio com login falso em produção; desligue DEV_MODE")
	}

	c.CookieSecure = !c.DevMode
	switch strings.ToLower(strings.TrimSpace(getenv("COOKIE_SECURE"))) {
	case "true":
		c.CookieSecure = true
	case "false":
		c.CookieSecure = false
	}
```

- [ ] **Step 4: Rodar os testes**

Run: `cd backend-gateway && go test ./internal/config/ -v`
Expected: PASS.

- [ ] **Step 5: Commitar**

```bash
git add backend-gateway/internal/config/config.go backend-gateway/internal/config/config_test.go
git commit -m "Recusa subir com DEV_MODE apontando para https, e decide o cookie Secure"
```

---

### Task 3 [B]: Cookie de sessão com prefixo `__Host-`

**Files:**
- Modify: `backend-gateway/internal/auth/session.go`
- Modify: `backend-gateway/internal/auth/oauth.go`
- Modify: `backend-gateway/internal/auth/middleware.go`
- Modify: `backend-gateway/internal/httpapi/api.go` (passa `cfg.CookieSecure` ao middleware)
- Test: `backend-gateway/internal/auth/session_test.go`

**Interfaces:**
- Consumes: `config.Config.CookieSecure` (Task 2).
- Produces: `auth.CookieName(secure bool) string`; `auth.LerSessao(r *http.Request) (string, bool)`; assinatura nova `auth.RequireSession(secret []byte, allowed []string) func(http.Handler) http.Handler` **inalterada** — a leitura passa a aceitar os dois nomes.

- [ ] **Step 1: Escrever os testes que falham**

Acrescentar em `backend-gateway/internal/auth/session_test.go`:

```go
func TestCookieNamePorModo(t *testing.T) {
	if got := CookieName(true); got != "__Host-batradio_session" {
		t.Errorf("com Secure esperava prefixo __Host-, veio %q", got)
	}
	// O prefixo __Host- exige Secure; em dev (http://localhost) o navegador
	// recusaria o cookie e o login nunca fecharia.
	if got := CookieName(false); got != "batradio_session" {
		t.Errorf("sem Secure esperava nome simples, veio %q", got)
	}
}

func TestLerSessaoAceitaOsDoisNomes(t *testing.T) {
	for _, nome := range []string{"__Host-batradio_session", "batradio_session"} {
		r := httptest.NewRequest("GET", "/api/status", nil)
		r.AddCookie(&http.Cookie{Name: nome, Value: "abc"})
		valor, ok := LerSessao(r)
		if !ok || valor != "abc" {
			t.Errorf("%s: esperava ler abc, veio %q ok=%v", nome, valor, ok)
		}
	}

	r := httptest.NewRequest("GET", "/api/status", nil)
	if _, ok := LerSessao(r); ok {
		t.Error("sem cookie não pode reportar sessão")
	}
}
```

Garantir os imports `net/http` e `net/http/httptest` no arquivo de teste.

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd backend-gateway && go test ./internal/auth/ -run 'CookieName|LerSessao' -v`
Expected: FAIL — `CookieName`/`LerSessao` não definidos.

- [ ] **Step 3: Implementar em session.go**

Em `internal/auth/session.go`, substituir a constante `SessionCookie` por:

```go
// Nomes do cookie de sessão. O prefixo __Host- é uma trava do navegador: ele só
// aceita o cookie se vier com Secure, Path=/ e sem Domain — então uma
// regressão que desligue o Secure quebra o login em vez de degradar em silêncio.
// Em dev (http://localhost) o prefixo é impossível, daí o nome simples.
const (
	SessionCookie     = "batradio_session"
	SessionCookieHost = "__Host-batradio_session"
)

// CookieName devolve o nome a gravar conforme o cookie vá ou não com Secure.
func CookieName(secure bool) string {
	if secure {
		return SessionCookieHost
	}
	return SessionCookie
}

// LerSessao procura o cookie de sessão nos dois nomes possíveis. Aceitar ambos
// evita deslogar todo mundo no deploy que liga o Secure.
func LerSessao(r *http.Request) (string, bool) {
	for _, nome := range []string{SessionCookieHost, SessionCookie} {
		if c, err := r.Cookie(nome); err == nil && c.Value != "" {
			return c.Value, true
		}
	}
	return "", false
}
```

Acrescentar `net/http` aos imports de `session.go`.

- [ ] **Step 4: Usar o nome certo no oauth.go**

Em `internal/auth/oauth.go`, em `setSessionCookie` (linha ~67), trocar a construção do cookie por:

```go
func (o *OAuth) setSessionCookie(w http.ResponseWriter, r *http.Request, email string) {
	token := SignSession(email, time.Now().Add(SessionTTL), o.cfg.SessionSecret)
	secure := o.cfg.CookieSecure
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName(secure),
		Value:    token,
		Path:     "/",
		MaxAge:   int(SessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}
```

Remover a função `secureCookie` (linha ~63) e seus usos: nos cookies de `state` (linha ~95) e no logout (linha ~151), trocar `Secure: secureCookie(r)` por `Secure: o.cfg.CookieSecure`.

No logout, expirar **os dois** nomes, senão um cookie antigo sobrevive:

```go
	for _, nome := range []string{SessionCookieHost, SessionCookie} {
		http.SetCookie(w, &http.Cookie{
			Name:     nome,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   o.cfg.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
	}
```

Se o parâmetro `r *http.Request` ficar sem uso em alguma dessas funções, mantê-lo na assinatura (os chamadores não mudam) e ignorá-lo com `_ = r` não é necessário em Go para parâmetros — parâmetro não usado não é erro de compilação.

- [ ] **Step 5: Ler pelos dois nomes no middleware**

Em `internal/auth/middleware.go`, substituir o bloco que lê o cookie:

```go
			valor, ok := LerSessao(r)
			if !ok {
				writeAuthError(w, http.StatusUnauthorized, "unauthenticated", "Sessão ausente. Faça login.")
				return
			}
			email, ok := VerifySession(valor, secret, time.Now())
```

- [ ] **Step 6: Rodar os testes**

Run: `cd backend-gateway && go test ./... `
Expected: PASS. Se algum teste existente referenciar `SessionCookie` diretamente ao montar request, ele continua válido — o middleware aceita os dois nomes.

- [ ] **Step 7: Commitar**

```bash
git add backend-gateway/internal/auth/
git commit -m "Cookie de sessão com prefixo __Host- quando vai com Secure"
```

---

### Task 4 [B]: Headers de segurança e rate limit no login

**Files:**
- Create: `backend-gateway/internal/httpsec/headers.go`
- Create: `backend-gateway/internal/httpsec/headers_test.go`
- Create: `backend-gateway/internal/httpsec/ratelimit.go`
- Create: `backend-gateway/internal/httpsec/ratelimit_test.go`
- Modify: `backend-gateway/main.go`
- Modify: `backend-gateway/internal/auth/oauth.go` (aplicar o limite nas rotas de login)

**Interfaces:**
- Produces: `httpsec.Headers(streamURL string, hsts bool) func(http.Handler) http.Handler`; `httpsec.NewLimiter(limite int, janela time.Duration) *httpsec.Limiter` com método `Middleware(next http.Handler) http.Handler`.

- [ ] **Step 1: Escrever o teste dos headers**

`backend-gateway/internal/httpsec/headers_test.go`:

```go
package httpsec

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHeadersDefineCSPeHSTS(t *testing.T) {
	h := Headers("https://stream.morcegaofm.com.br/live", true)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	csp := rec.Header().Get("Content-Security-Policy")
	if !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Errorf("CSP sem frame-ancestors: %q", csp)
	}
	if !strings.Contains(csp, "https://stream.morcegaofm.com.br/live") {
		t.Errorf("CSP tem de liberar o stream em media-src: %q", csp)
	}
	if !strings.Contains(csp, "img-src 'self' data: https:") {
		t.Errorf("CSP tem de liberar capa https (iTunes e /uploads do site): %q", csp)
	}
	if rec.Header().Get("Strict-Transport-Security") == "" {
		t.Error("faltou HSTS")
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("faltou X-Frame-Options: DENY")
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("faltou nosniff")
	}
	if rec.Header().Get("Referrer-Policy") != "same-origin" {
		t.Error("faltou Referrer-Policy")
	}
}

func TestHeadersSemHSTSForaDeHTTPS(t *testing.T) {
	// Em dev (http://localhost) o HSTS travaria o navegador em https para sempre.
	h := Headers("", false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Header().Get("Strict-Transport-Security") != "" {
		t.Error("HSTS não pode sair quando o gateway não está em https")
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd backend-gateway && go test ./internal/httpsec/ -v`
Expected: FAIL — pacote não existe.

- [ ] **Step 3: Implementar os headers**

`backend-gateway/internal/httpsec/headers.go`:

```go
// Package httpsec reúne os middlewares de borda do gateway: headers de
// segurança e limite de tentativas no login.
package httpsec

import (
	"net/http"
	"strings"
)

// Headers aplica os headers de segurança da resposta.
//
// O SPA e os assets são servidos pelo próprio binário (go:embed), então a CSP
// pode ser fechada: nada de script externo. As exceções são a capa (https:,
// porque vem do iTunes ou do /uploads do site) e o stream do Icecast.
func Headers(streamURL string, hsts bool) func(http.Handler) http.Handler {
	mediaSrc := "'self'"
	if u := strings.TrimSpace(streamURL); u != "" {
		mediaSrc += " " + u
	}

	csp := strings.Join([]string{
		"default-src 'self'",
		"img-src 'self' data: https:",
		"media-src " + mediaSrc,
		"script-src 'self'",
		// O Vite injeta estilo inline no build; sem unsafe-inline o painel fica sem CSS.
		"style-src 'self' 'unsafe-inline'",
		"connect-src 'self'",
		"font-src 'self' data:",
		"object-src 'none'",
		"frame-ancestors 'none'",
		"base-uri 'self'",
		"form-action 'self'",
	}, "; ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("Content-Security-Policy", csp)
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "same-origin")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			if hsts {
				h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

- [ ] **Step 4: Escrever o teste do rate limit**

`backend-gateway/internal/httpsec/ratelimit_test.go`:

```go
package httpsec

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimiterBloqueiaAcimaDoLimite(t *testing.T) {
	agora := time.Unix(1_700_000_000, 0)
	l := NewLimiter(3, time.Minute)
	l.now = func() time.Time { return agora }

	h := l.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	pedir := func() int {
		r := httptest.NewRequest("GET", "/auth/login", nil)
		r.RemoteAddr = "203.0.113.7:12345"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}

	for i := 0; i < 3; i++ {
		if got := pedir(); got != http.StatusOK {
			t.Fatalf("tentativa %d deveria passar, veio %d", i+1, got)
		}
	}
	if got := pedir(); got != http.StatusTooManyRequests {
		t.Fatalf("4ª tentativa deveria dar 429, veio %d", got)
	}

	// Passada a janela, libera de novo.
	agora = agora.Add(time.Minute + time.Second)
	if got := pedir(); got != http.StatusOK {
		t.Fatalf("depois da janela deveria liberar, veio %d", got)
	}
}

func TestLimiterSeparaPorIP(t *testing.T) {
	l := NewLimiter(1, time.Minute)
	h := l.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	pedir := func(ip string) int {
		r := httptest.NewRequest("GET", "/auth/login", nil)
		r.RemoteAddr = ip + ":1111"
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}

	pedir("203.0.113.1")
	if got := pedir("203.0.113.2"); got != http.StatusOK {
		t.Fatalf("IP diferente não pode herdar o limite do vizinho, veio %d", got)
	}
}
```

- [ ] **Step 5: Implementar o rate limit**

`backend-gateway/internal/httpsec/ratelimit.go`:

```go
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
```

- [ ] **Step 6: Rodar os testes**

Run: `cd backend-gateway && go test ./internal/httpsec/ -v`
Expected: PASS, 4 testes.

- [ ] **Step 8: Ligar no main.go**

Em `backend-gateway/main.go`, no bloco do router (depois de `r.Use(middleware.Recoverer)`):

```go
	r.Use(httpsec.Headers(cfg.StreamURL, cfg.CookieSecure))
```

`cfg.CookieSecure` serve de proxy para "estou atrás de https": é exatamente a condição em que o HSTS deve sair.

Acrescentar aos imports: `"github.com/Morcegao-FM/batradio/backend-gateway/internal/httpsec"`.

- [ ] **Step 8: Aplicar o limite nas rotas de login**

Em `internal/auth/oauth.go`, em `Register` (linha ~51):

```go
func (o *OAuth) Register(r chi.Router) {
	// 10 tentativas por minuto e por IP: folga para quem erra a conta, teto
	// para quem está sondando o callback.
	limite := httpsec.NewLimiter(10, time.Minute)
	r.With(limite.Middleware).Get("/auth/login", o.handleLogin)
	r.With(limite.Middleware).Get("/auth/callback", o.handleCallback)
	r.Post("/auth/logout", o.handleLogout)
}
```

Acrescentar o import de `httpsec` em `oauth.go`.

- [ ] **Step 9: Rodar tudo**

Run: `cd backend-gateway && go test ./...`
Expected: PASS.

- [ ] **Step 10: Commitar**

```bash
git add backend-gateway/internal/httpsec/ backend-gateway/main.go backend-gateway/internal/auth/oauth.go
git commit -m "Headers de segurança na borda e limite de tentativas no login"
```

---

### Task 5 [B]: `/healthz` e o teste que tranca a superfície pública

**Files:**
- Modify: `backend-gateway/internal/httpapi/api.go`
- Test: `backend-gateway/internal/httpapi/superficie_test.go` (novo)

**Interfaces:**
- Consumes: `Server.Routes(o *auth.OAuth) chi.Router` (existente).
- Produces: `GET /healthz` → `200 {"status":"ok"}`, público.

Esta é a tarefa mais importante do plano em termos de segurança: é a rede que pega a rota nova que alguém esqueceu de proteger.

- [ ] **Step 1: Escrever o teste que falha**

`backend-gateway/internal/httpapi/superficie_test.go`:

```go
package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// publicasPermitidas é a lista FECHADA de rotas que podem responder sem sessão.
// Acrescentar item aqui é uma decisão de segurança — não faça isso para calar
// um teste vermelho.
var publicasPermitidas = map[string]bool{
	"GET /healthz":       true,
	"GET /auth/login":    true,
	"GET /auth/callback": true,
	"POST /auth/logout":  true,
}

// TestNenhumaRotaNovaSemSessao percorre as rotas registradas no chi e exige que
// toda rota fora da allowlist responda 401 sem cookie de sessão.
func TestNenhumaRotaNovaSemSessao(t *testing.T) {
	srv := newTestServer(t, "http://127.0.0.1:1")
	router := srv.Routes(nil)

	var verificadas int
	err := chi.Walk(router, func(metodo, rota string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		chave := metodo + " " + rota
		if publicasPermitidas[chave] {
			return nil
		}

		// Substitui parâmetros de rota por um valor qualquer.
		caminho := strings.ReplaceAll(rota, "{name}", "qualquer")
		if strings.Contains(caminho, "{") {
			t.Errorf("%s: parâmetro de rota não previsto no teste (%s)", chave, caminho)
			return nil
		}

		req := httptest.NewRequest(metodo, caminho, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s respondeu %d sem sessão; esperava 401. "+
				"Rota nova precisa nascer atrás de RequireSession.", chave, rec.Code)
		}
		verificadas++
		return nil
	})
	if err != nil {
		t.Fatalf("chi.Walk: %v", err)
	}
	if verificadas == 0 {
		t.Fatal("nenhuma rota verificada — o teste não está enxergando o router")
	}
}

func TestHealthzEhPublicoENaoVazaNada(t *testing.T) {
	srv := newTestServer(t, "http://127.0.0.1:1")
	router := srv.Routes(nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("GET", "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("healthz respondeu %d", rec.Code)
	}
	corpo := rec.Body.String()
	for _, proibido := range []string{"nodeHost", "libraryCount", "version", "9320"} {
		if strings.Contains(corpo, proibido) {
			t.Errorf("healthz é público e não pode expor %q; corpo: %s", proibido, corpo)
		}
	}
}
```

`newTestServer` já existe em `internal/httpapi/api_test.go:155` com a assinatura `newTestServer(t *testing.T, nodeURL string) *Server`. Usar `newTestServer(t, "http://127.0.0.1:1")` — o Node não precisa responder, porque nenhuma das requisições deste teste passa do middleware de sessão.

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd backend-gateway && go test ./internal/httpapi/ -run 'Superficie|Healthz|SemSessao' -v`
Expected: FAIL — `/healthz` responde 404 (cai no NotFound), então `TestHealthzEhPublico` falha.

- [ ] **Step 3: Implementar o /healthz**

Em `internal/httpapi/api.go`, dentro de `Routes`, antes do `r.Route("/api", ...)`:

```go
	// Público de propósito, e deliberadamente vazio: é só o sinal de vida que o
	// `docker compose up -d --wait` e o rollback do deploy leem. Não expõe
	// versão, host do Node nem estado do acervo — isso é superfície de graça
	// para quem estiver sondando.
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"status":"ok"}`))
	})
```

- [ ] **Step 4: Rodar os testes**

Run: `cd backend-gateway && go test ./internal/httpapi/ -v`
Expected: PASS. Se `TestNenhumaRotaNovaSemSessao` acusar alguma rota, é achado real: proteger a rota, não relaxar a allowlist.

- [ ] **Step 5: Commitar**

```bash
git add backend-gateway/internal/httpapi/api.go backend-gateway/internal/httpapi/superficie_test.go
git commit -m "Healthz público e vazio, com teste que tranca a superfície sem sessão"
```

---

### Task 6 [B]: Log de auditoria das escritas

**Files:**
- Create: `backend-gateway/internal/httpapi/auditoria.go`
- Create: `backend-gateway/internal/httpapi/auditoria_test.go`
- Modify: `backend-gateway/internal/httpapi/api.go`

**Interfaces:**
- Consumes: `auth.EmailFromContext(ctx) string`.
- Produces: `auditoria(log func(formato string, args ...any)) func(http.Handler) http.Handler`.

- [ ] **Step 1: Escrever o teste que falha**

`backend-gateway/internal/httpapi/auditoria_test.go`:

```go
package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
)

func TestAuditoriaRegistraEscritaComEmail(t *testing.T) {
	var linhas []string
	h := auditoria(func(f string, a ...any) { linhas = append(linhas, fmt.Sprintf(f, a...)) })(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))

	req := httptest.NewRequest("POST", "/api/playlist/remove", nil)
	req = req.WithContext(auth.ContextComEmail(context.Background(), "aguergolet@gmail.com"))
	h.ServeHTTP(httptest.NewRecorder(), req)

	if len(linhas) != 1 {
		t.Fatalf("esperava 1 linha de auditoria, veio %d: %v", len(linhas), linhas)
	}
	linha := linhas[0]
	for _, esperado := range []string{"aguergolet@gmail.com", "POST", "/api/playlist/remove", "200"} {
		if !strings.Contains(linha, esperado) {
			t.Errorf("linha de auditoria sem %q: %s", esperado, linha)
		}
	}
}

func TestAuditoriaIgnoraLeitura(t *testing.T) {
	var linhas []string
	h := auditoria(func(f string, a ...any) { linhas = append(linhas, fmt.Sprintf(f, a...)) })(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	// GET de status roda a cada 2 segundos; auditar leitura só afogaria o log.
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/api/status", nil))

	if len(linhas) != 0 {
		t.Errorf("leitura não deve ser auditada, veio: %v", linhas)
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd backend-gateway && go test ./internal/httpapi/ -run Auditoria -v`
Expected: FAIL — `auditoria` e `auth.ContextComEmail` não existem.

- [ ] **Step 3: Expor o construtor de contexto em auth**

Em `internal/auth/middleware.go`, abaixo de `EmailFromContext`:

```go
// ContextComEmail injeta o e-mail autenticado. Existe para os testes montarem
// um request já autenticado sem forjar cookie.
func ContextComEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailKey, email)
}
```

E fazer `RequireSession` usá-la, trocando a última linha do handler por:

```go
			next.ServeHTTP(w, r.WithContext(ContextComEmail(r.Context(), email)))
```

- [ ] **Step 4: Implementar a auditoria**

`backend-gateway/internal/httpapi/auditoria.go`:

```go
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
```

- [ ] **Step 5: Ligar no router**

Em `internal/httpapi/api.go`, dentro de `r.Route("/api", ...)`, logo **depois** do `api.Use(auth.RequireSession(...))` (a ordem importa: a auditoria precisa do e-mail já no contexto):

```go
		api.Use(auditoria(log.Printf))
```

Acrescentar `"log"` aos imports de `api.go`.

- [ ] **Step 6: Rodar os testes**

Run: `cd backend-gateway && go test ./...`
Expected: PASS.

- [ ] **Step 7: Commitar**

```bash
git add backend-gateway/internal/httpapi/auditoria.go backend-gateway/internal/httpapi/auditoria_test.go \
        backend-gateway/internal/httpapi/api.go backend-gateway/internal/auth/middleware.go
git commit -m "Registra quem mudou a fila, o player e as playlists"
```

---

### Task 7 [B]: Cliente do catálogo de capa, com cache e degradação

**Files:**
- Create: `backend-gateway/internal/catalogo/client.go`
- Create: `backend-gateway/internal/catalogo/client_test.go`
- Modify: `backend-gateway/internal/config/config.go`
- Modify: `backend-gateway/internal/config/config_test.go`

**Interfaces:**
- Consumes: `POST {CatalogoURL}/api/servico/musicas/lookup` com header `X-Servico-Key` (Task 1).
- Produces: `catalogo.New(baseURL, chave string, ttl time.Duration) *catalogo.Client`; `(*Client).Lookup(ctx context.Context, arquivos []string) map[string]Info`; `catalogo.Info{Tipo, Artista, Titulo, ImagemURL, NomeExibicao string; Ano int}`. `Lookup` **nunca** devolve erro. Consumido pela Task 8.
- Produces: `config.Config.CatalogoURL string` e `config.Config.CatalogoChave string`.

- [ ] **Step 1: Escrever os testes que falham**

`backend-gateway/internal/catalogo/client_test.go`:

```go
package catalogo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func servidorFake(t *testing.T, chaveEsperada string, resposta any, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Servico-Key") != chaveEsperada {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resposta)
	}))
}

func TestLookupDevolveOQueOSiteConhece(t *testing.T) {
	srv := servidorFake(t, "chave", []map[string]any{
		{"arquivo": "a.mp3", "tipo": "musica", "artista": "AC/DC",
			"titulo": "Back in Black", "ano": 1980, "imagemUrl": "https://capa/512.jpg"},
	}, http.StatusOK)
	defer srv.Close()

	c := New(srv.URL, "chave", time.Hour)
	got := c.Lookup(context.Background(), []string{"a.mp3", "b.mp3"})

	if len(got) != 1 {
		t.Fatalf("esperava 1 resultado, veio %d", len(got))
	}
	if got["a.mp3"].ImagemURL != "https://capa/512.jpg" || got["a.mp3"].Ano != 1980 {
		t.Errorf("resultado errado: %+v", got["a.mp3"])
	}
}

func TestLookupDegradaEmFalha(t *testing.T) {
	for _, caso := range []struct {
		nome   string
		status int
	}{
		{"500 do site", http.StatusInternalServerError},
		{"401 chave errada", http.StatusUnauthorized},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			srv := servidorFake(t, "chave", nil, caso.status)
			defer srv.Close()

			c := New(srv.URL, "chave-errada", time.Hour)
			got := c.Lookup(context.Background(), []string{"a.mp3"})
			if len(got) != 0 {
				t.Errorf("falha do site tem de virar mapa vazio, veio %+v", got)
			}
		})
	}
}

func TestLookupDegradaComSiteFora(t *testing.T) {
	// Porta fechada: o painel não pode travar por causa de capa.
	c := New("http://127.0.0.1:1", "chave", time.Hour)
	feito := make(chan map[string]Info, 1)
	go func() { feito <- c.Lookup(context.Background(), []string{"a.mp3"}) }()

	select {
	case got := <-feito:
		if len(got) != 0 {
			t.Errorf("esperava mapa vazio, veio %+v", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Lookup travou; tem de ter timeout curto")
	}
}

func TestLookupDesligadoSemConfiguracao(t *testing.T) {
	c := New("", "", time.Hour)
	if got := c.Lookup(context.Background(), []string{"a.mp3"}); len(got) != 0 {
		t.Errorf("sem URL/chave o cliente fica desligado, veio %+v", got)
	}
}

func TestCacheNegativoNaoReconsulta(t *testing.T) {
	var chamadas int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chamadas++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New(srv.URL, "chave", time.Hour)
	c.Lookup(context.Background(), []string{"desconhecida.mp3"})
	c.Lookup(context.Background(), []string{"desconhecida.mp3"})

	if chamadas != 1 {
		t.Errorf("arquivo desconhecido não pode ser reperguntado a cada render; chamadas=%d", chamadas)
	}
}

func TestCacheExpiraNoTTL(t *testing.T) {
	var chamadas int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chamadas++
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"arquivo":"a.mp3","tipo":"musica","imagemUrl":"https://capa/1.jpg"}]`))
	}))
	defer srv.Close()

	agora := time.Unix(1_700_000_000, 0)
	c := New(srv.URL, "chave", time.Hour)
	c.now = func() time.Time { return agora }

	c.Lookup(context.Background(), []string{"a.mp3"})
	c.Lookup(context.Background(), []string{"a.mp3"})
	if chamadas != 1 {
		t.Fatalf("segunda consulta dentro do TTL deveria vir do cache; chamadas=%d", chamadas)
	}

	agora = agora.Add(time.Hour + time.Minute)
	c.Lookup(context.Background(), []string{"a.mp3"})
	if chamadas != 2 {
		t.Errorf("depois do TTL deveria reconsultar; chamadas=%d", chamadas)
	}
}

func TestLookupRespeitaTetoDe200PorChamada(t *testing.T) {
	var lotes []int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var corpo struct {
			Arquivos []string `json:"arquivos"`
		}
		json.NewDecoder(r.Body).Decode(&corpo)
		lotes = append(lotes, len(corpo.Arquivos))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	arquivos := make([]string, 250)
	for i := range arquivos {
		arquivos[i] = fmt.Sprintf("faixa-%d.mp3", i)
	}

	c := New(srv.URL, "chave", time.Hour)
	c.Lookup(context.Background(), arquivos)

	for _, n := range lotes {
		if n > 200 {
			t.Errorf("lote de %d arquivos estoura o teto do endpoint (200)", n)
		}
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd backend-gateway && go test ./internal/catalogo/ -v`
Expected: FAIL — pacote não existe.

- [ ] **Step 3: Implementar o cliente**

`backend-gateway/internal/catalogo/client.go`:

```go
// Package catalogo lê os dados de exibição de faixa (capa, artista, título,
// ano) do catálogo que a API do website mantém.
//
// Regra de ouro deste pacote: capa é enfeite, o painel é operação. Nenhuma
// falha aqui pode atrapalhar o controle do MPD — por isso Lookup não devolve
// erro, só o que conseguiu obter.
package catalogo

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// maxPorChamada espelha ServicoMusicasController.MaxArquivos no website.
const maxPorChamada = 200

// Info é o que o catálogo sabe de uma faixa.
type Info struct {
	Tipo         string `json:"tipo"`
	Artista      string `json:"artista"`
	Titulo       string `json:"titulo"`
	Ano          int    `json:"ano"`
	ImagemURL    string `json:"imagemUrl"`
	NomeExibicao string `json:"nomeExibicao"`
}

type respostaItem struct {
	Arquivo string `json:"arquivo"`
	Info
}

type entrada struct {
	info      Info
	encontrou bool
	em        time.Time
}

type Client struct {
	baseURL string
	chave   string
	ttl     time.Duration
	http    *http.Client

	mu    sync.Mutex
	cache map[string]entrada

	now func() time.Time
}

// New devolve um cliente. baseURL ou chave vazios deixam o cliente desligado —
// é assim que o ambiente de desenvolvimento roda sem o website por perto.
func New(baseURL, chave string, ttl time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		chave:   chave,
		ttl:     ttl,
		// Timeout curto: o painel espera por isso ao montar a tela.
		http:  &http.Client{Timeout: 3 * time.Second},
		cache: make(map[string]entrada),
		now:   time.Now,
	}
}

func (c *Client) ligado() bool { return c.baseURL != "" && c.chave != "" }

// Lookup devolve o que sabe sobre os arquivos pedidos. Arquivo sem linha no
// catálogo simplesmente não aparece no mapa. Nunca devolve erro.
func (c *Client) Lookup(ctx context.Context, arquivos []string) map[string]Info {
	resultado := make(map[string]Info)
	if !c.ligado() || len(arquivos) == 0 {
		return resultado
	}

	faltando := c.doCache(arquivos, resultado)
	if len(faltando) == 0 {
		return resultado
	}

	for inicio := 0; inicio < len(faltando); inicio += maxPorChamada {
		fim := min(inicio+maxPorChamada, len(faltando))
		c.buscarLote(ctx, faltando[inicio:fim], resultado)
	}
	return resultado
}

// doCache preenche o resultado com o que está em cache e devolve o que falta.
func (c *Client) doCache(arquivos []string, resultado map[string]Info) []string {
	agora := c.now()
	var faltando []string

	c.mu.Lock()
	defer c.mu.Unlock()

	vistos := make(map[string]bool, len(arquivos))
	for _, arquivo := range arquivos {
		if arquivo == "" || vistos[arquivo] {
			continue
		}
		vistos[arquivo] = true

		e, ok := c.cache[arquivo]
		if ok && agora.Sub(e.em) <= c.ttl {
			// Cache negativo: arquivo que o site não conhece fica registrado
			// como "não tem", senão ele seria reperguntado a cada render.
			if e.encontrou {
				resultado[arquivo] = e.info
			}
			continue
		}
		faltando = append(faltando, arquivo)
	}
	return faltando
}

func (c *Client) buscarLote(ctx context.Context, lote []string, resultado map[string]Info) {
	corpo, err := json.Marshal(map[string][]string{"arquivos": lote})
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/servico/musicas/lookup", bytes.NewReader(corpo))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Servico-Key", c.chave)

	resp, err := c.http.Do(req)
	if err != nil {
		// Sem a chave no log.
		log.Printf("catálogo indisponível: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("catálogo respondeu %d", resp.StatusCode)
		return
	}

	var itens []respostaItem
	if err := json.NewDecoder(resp.Body).Decode(&itens); err != nil {
		log.Printf("catálogo devolveu JSON inesperado: %v", err)
		return
	}

	agora := c.now()
	encontrados := make(map[string]Info, len(itens))
	for _, item := range itens {
		encontrados[item.Arquivo] = item.Info
		resultado[item.Arquivo] = item.Info
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, arquivo := range lote {
		info, ok := encontrados[arquivo]
		c.cache[arquivo] = entrada{info: info, encontrou: ok, em: agora}
	}
}
```

- [ ] **Step 4: Acrescentar a configuração**

Em `internal/config/config.go`, na struct:

```go
	// API do website, que serve o catálogo de faixas (capa, artista, ano).
	// Vazios deixam o enriquecimento desligado — o painel funciona sem ele.
	CatalogoURL   string
	CatalogoChave string
```

Em `Load`, junto aos outros `getenv`:

```go
		CatalogoURL:   strings.TrimRight(getenv("CATALOGO_URL"), "/"),
		CatalogoChave: getenv("CATALOGO_CHAVE"),
```

Teste em `config_test.go`:

```go
func TestCatalogoEhOpcional(t *testing.T) {
	env := map[string]string{
		"NODE_BACKEND_URL":     "http://node:9320",
		"NODE_API_KEY":         "k",
		"SESSION_SECRET":       strings.Repeat("s", 32),
		"ALLOWED_EMAILS":       "a@b.com",
		"GOOGLE_CLIENT_ID":     "id",
		"GOOGLE_CLIENT_SECRET": "secret",
		"OAUTH_REDIRECT_URL":   "https://batradio.morcegaofm.com.br/auth/callback",
	}
	c, err := Load(func(k string) string { return env[k] })
	if err != nil {
		t.Fatalf("catálogo não configurado não pode impedir o gateway de subir: %v", err)
	}
	if c.CatalogoURL != "" {
		t.Error("CatalogoURL deveria vir vazio")
	}
}
```

- [ ] **Step 5: Rodar os testes**

Run: `cd backend-gateway && go test ./internal/catalogo/ ./internal/config/ -v`
Expected: PASS.

- [ ] **Step 6: Commitar**

```bash
git add backend-gateway/internal/catalogo/ backend-gateway/internal/config/
git commit -m "Lê capa e ano do catálogo do site, com cache e sem travar o painel"
```

---

### Task 8 [B]: Capa no status, na fila e na tela

**Files:**
- Modify: `backend-gateway/internal/model/model.go`
- Modify: `backend-gateway/internal/httpapi/api.go`
- Modify: `backend-gateway/internal/httpapi/player.go`
- Modify: `backend-gateway/internal/httpapi/playlist.go`
- Create: `backend-gateway/internal/httpapi/enriquecer_test.go`
- Modify: `backend-gateway/main.go`
- Modify: `frontend-web/src/lib/types.ts`
- Modify: `frontend-web/src/components/PlayerBar.tsx`
- Modify: `frontend-web/src/components/PlayerBar.module.css`

**Interfaces:**
- Consumes: `catalogo.Client.Lookup` (Task 7).
- Produces: campos novos no JSON de `Song`: `imageUrl`, `displayName`, `year`, `kind` — todos `omitempty`. A SPA os lê pelo tipo `Song` de `lib/types.ts`.

- [ ] **Step 1: Escrever o teste que falha**

`backend-gateway/internal/httpapi/enriquecer_test.go`:

```go
package httpapi

import (
	"context"
	"testing"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/catalogo"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
)

type catalogoFake struct {
	dados map[string]catalogo.Info
	vezes int
}

func (f *catalogoFake) Lookup(_ context.Context, arquivos []string) map[string]catalogo.Info {
	f.vezes++
	out := make(map[string]catalogo.Info)
	for _, a := range arquivos {
		if info, ok := f.dados[a]; ok {
			out[a] = info
		}
	}
	return out
}

func TestEnriquecerPreencheCapaEMantemOResto(t *testing.T) {
	fake := &catalogoFake{dados: map[string]catalogo.Info{
		"a.mp3": {Tipo: "musica", ImagemURL: "https://capa/1.jpg", Ano: 1980, Artista: "AC/DC"},
	}}
	s := &Server{catalogo: fake}

	musicas := []model.Song{
		{File: "a.mp3", Artist: "AC/DC", Title: "Back in Black", Time: 255},
		{File: "b.mp3", Artist: "Desconhecido", Title: "Sem capa", Time: 120},
	}
	s.enriquecer(context.Background(), musicas)

	if musicas[0].ImageURL != "https://capa/1.jpg" || musicas[0].Year != 1980 {
		t.Errorf("faixa conhecida não foi enriquecida: %+v", musicas[0])
	}
	if musicas[0].Title != "Back in Black" || musicas[0].Time != 255 {
		t.Errorf("enriquecer não pode mexer no que veio do MPD: %+v", musicas[0])
	}
	if musicas[1].ImageURL != "" {
		t.Errorf("faixa desconhecida tem de ficar sem capa: %+v", musicas[1])
	}
}

func TestEnriquecerSemCatalogoNaoQuebra(t *testing.T) {
	s := &Server{catalogo: nil}
	musicas := []model.Song{{File: "a.mp3"}}
	s.enriquecer(context.Background(), musicas) // não pode entrar em pânico
	if musicas[0].ImageURL != "" {
		t.Error("sem catálogo não há capa")
	}
}

func TestEnriquecerNomeExibicaoVenceArtistaTitulo(t *testing.T) {
	// É o que resolve vinheta com nome de arquivo feio: a Batcaverna define o
	// nome de exibição e o BatRadio o respeita.
	fake := &catalogoFake{dados: map[string]catalogo.Info{
		"vinheta-01.mp3": {Tipo: "vinheta", NomeExibicao: "Vinheta de abertura"},
	}}
	s := &Server{catalogo: fake}

	musicas := []model.Song{{File: "vinheta-01.mp3", Artist: "", Title: "vinheta-01"}}
	s.enriquecer(context.Background(), musicas)

	if musicas[0].DisplayName != "Vinheta de abertura" || musicas[0].Kind != "vinheta" {
		t.Errorf("nome de exibição não chegou: %+v", musicas[0])
	}
}
```

- [ ] **Step 2: Rodar e ver falhar**

Run: `cd backend-gateway && go test ./internal/httpapi/ -run Enriquecer -v`
Expected: FAIL — `Server.catalogo`, `enriquecer` e os campos novos não existem.

- [ ] **Step 3: Acrescentar os campos ao model**

Em `internal/model/model.go`, na struct `Song`, depois de `ID`:

```go
	// Preenchidos pelo gateway a partir do catálogo do website; o Node/MPD
	// nunca manda estes campos. omitempty para não inchar a lista do acervo.
	ImageURL    string `json:"imageUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Year        int    `json:"year,omitempty"`
	Kind        string `json:"kind,omitempty"`
```

- [ ] **Step 4: Implementar o enriquecimento**

Em `internal/httpapi/api.go`, acrescentar a interface e o campo na `Server`:

```go
// catalogoLookup é o que o gateway usa do catálogo do website. Interface (e não
// o tipo concreto) para o teste injetar um fake.
type catalogoLookup interface {
	Lookup(ctx context.Context, arquivos []string) map[string]catalogo.Info
}
```

Na struct `Server`, depois de `poller`:

```go
	catalogo catalogoLookup
```

Em `NewServer`, acrescentar o parâmetro e repassar:

```go
func NewServer(cfg *config.Config, n *node.Client, cat catalogoLookup) *Server {
	return &Server{
		cfg:      cfg,
		node:     n,
		lib:      library.NewIndex(),
		poller:   NewPoller(n, 2*time.Second),
		catalogo: cat,
		now:      time.Now,
		randInt:  rand.Intn,
	}
}
```

E o método, no mesmo arquivo:

```go
// enriquecer sobrepõe os dados de exibição do catálogo nas faixas dadas,
// in-place. Silencioso por natureza: catálogo desligado ou fora do ar deixa as
// faixas exatamente como vieram do MPD.
func (s *Server) enriquecer(ctx context.Context, musicas []model.Song) {
	if s.catalogo == nil || len(musicas) == 0 {
		return
	}

	arquivos := make([]string, 0, len(musicas))
	for _, m := range musicas {
		if m.File != "" {
			arquivos = append(arquivos, m.File)
		}
	}

	info := s.catalogo.Lookup(ctx, arquivos)
	for i := range musicas {
		dados, ok := info[musicas[i].File]
		if !ok {
			continue
		}
		musicas[i].ImageURL = dados.ImagemURL
		musicas[i].DisplayName = dados.NomeExibicao
		musicas[i].Year = dados.Ano
		musicas[i].Kind = dados.Tipo
	}
}
```

Acrescentar o import de `catalogo` em `api.go`.

- [ ] **Step 5: Rodar o teste do enriquecimento**

Run: `cd backend-gateway && go test ./internal/httpapi/ -run Enriquecer -v`
Expected: PASS.

- [ ] **Step 6: Chamar em `handleStatus`**

`handleStatus` está em `internal/httpapi/api.go:187` (não em `player.go`). Substituir por:

```go
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	st, err := s.currentStatus(r.Context())
	if err != nil {
		writeNodeErr(w, err)
		return
	}
	s.enriquecerStatus(r.Context(), &st)
	writeJSON(w, http.StatusOK, st)
}

// enriquecerStatus completa só a faixa atual e a próxima. Nunca a fila inteira:
// o painel chama /api/status a cada 2 segundos.
func (s *Server) enriquecerStatus(ctx context.Context, st *model.Status) {
	var visiveis []model.Song
	if st.Current != nil {
		visiveis = append(visiveis, *st.Current)
	}
	if st.Next != nil {
		visiveis = append(visiveis, *st.Next)
	}
	if len(visiveis) == 0 {
		return
	}

	s.enriquecer(ctx, visiveis)

	i := 0
	if st.Current != nil {
		st.Current = &visiveis[i]
		i++
	}
	if st.Next != nil {
		st.Next = &visiveis[i]
	}
}
```

`model.Status` tem os campos `Current *Song` e `Next *Song` (`internal/model/model.go:66-67`).

- [ ] **Step 7: Chamar em `handleGetPlaylist`**

Em `internal/httpapi/playlist.go:58`, substituir o `writeJSON` final de `handleGetPlaylist` por:

```go
	// Só a página visível: a fila pode ter centenas de faixas e o teto do
	// endpoint do site é 200 por chamada. `filtered` vem de
	// schedule.ComputeTimes, que devolve slice nova a cada request — mexer nela
	// não toca a fila em cache.
	pagina := filtered[offset:end]
	musicas := make([]model.Song, len(pagina))
	for i := range pagina {
		musicas[i] = pagina[i].Song
	}
	s.enriquecer(r.Context(), musicas)
	for i := range pagina {
		pagina[i].Song = musicas[i]
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items":      pagina,
		"total":      total,
		"offset":     offset,
		"limit":      limit,
		"currentPos": st.Song,
		"version":    s.poller.Version(),
	})
```

`schedule.QueueItem` embute `model.Song` por valor (`internal/schedule/schedule.go:13-16`), então `pagina[i].Song = ...` é a atribuição certa.

Acrescentar aos imports de `playlist.go`:

```go
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
```

Não enriquecer `/api/library` (acervo inteiro, quase nada tem linha no catálogo) nem `/api/playlist/current` (devolve só `{"pos":N,"index":N}`, sem faixa).

- [ ] **Step 7: Ligar no main.go**

Em `backend-gateway/main.go`, antes de `httpapi.NewServer`:

```go
	cat := catalogo.New(cfg.CatalogoURL, cfg.CatalogoChave, time.Hour)
	if cfg.CatalogoURL == "" {
		log.Println("catálogo do site não configurado — o painel roda sem capa")
	}

	nodeClient := node.New(cfg.NodeBackendURL, cfg.NodeAPIKey)
	server := httpapi.NewServer(cfg, nodeClient, cat)
```

Acrescentar o import de `catalogo`.

- [ ] **Step 9: Rodar o Go inteiro**

Antes de rodar, corrigir o único call site existente: `internal/httpapi/api_test.go:164` passa `NewServer(cfg, node.New(nodeURL, "key"))`. Vira:

```go
	s := NewServer(cfg, node.New(nodeURL, "key"), nil)
```

`nil` é o caso "sem catálogo", que é o que os testes existentes querem.

Run: `cd backend-gateway && go test ./...`
Expected: PASS.

- [ ] **Step 10: Tipos no frontend**

Em `frontend-web/src/lib/types.ts`, na interface `Song`, depois de `id`:

```ts
  /** Preenchidos pelo gateway a partir do catálogo do site; ausentes se ele não souber. */
  imageUrl?: string
  displayName?: string
  year?: number
  kind?: string
```

- [ ] **Step 11: Mostrar a capa no PlayerBar**

Em `frontend-web/src/components/PlayerBar.tsx`, no bloco que mostra a faixa atual, trocar o logo fixo por capa com fallback. Componente novo no mesmo arquivo, acima de `PlayerBar`:

```tsx
// Capa da faixa. Sem capa (ou com URL que falha) cai no logo — o painel nunca
// fica com buraco na tela por causa do catálogo.
function Capa({ src, alt }: { src?: string; alt: string }) {
  const [falhou, setFalhou] = useState(false)
  const usar = src && !falhou ? src : logo
  return (
    <img
      className={styles.capa}
      src={usar}
      alt={alt}
      onError={() => setFalhou(true)}
    />
  )
}
```

E usar `<Capa src={status?.current?.imageUrl} alt="" />` no lugar do `<img src={logo} .../>` da faixa atual. Onde o nome da faixa é exibido, preferir `displayName`:

```tsx
const nome = status?.current?.displayName || `${status?.current?.artist} — ${status?.current?.title}`
```

Em `PlayerBar.module.css`, acrescentar:

```css
.capa {
  width: 48px;
  height: 48px;
  border-radius: 4px;
  object-fit: cover;
  flex-shrink: 0;
}
```

- [ ] **Step 12: Rodar o frontend**

Run: `cd frontend-web && npm test && npm run build`
Expected: PASS nos testes e build sem erro de tipo.

- [ ] **Step 13: Commitar**

```bash
git add backend-gateway/internal/model/model.go backend-gateway/internal/httpapi/ backend-gateway/main.go \
        frontend-web/src/lib/types.ts frontend-web/src/components/PlayerBar.tsx \
        frontend-web/src/components/PlayerBar.module.css
git commit -m "Mostra capa e nome de exibição do catálogo no tocando agora e na fila"
```

---

### Task 9 [B]: Compose de produção e documentação do deploy

**Files:**
- Create: `docker-compose.main.yml`
- Modify: `docker-compose.yml` (comentário de escopo)
- Modify: `backend-gateway/.env.example`
- Modify: `README.md`
- Create: `docs/DEPLOY.md`

**Interfaces:**
- Consumes: `GET /healthz` (Task 5); `CATALOGO_URL`/`CATALOGO_CHAVE` (Task 7).
- Produces: o arquivo que a Task 10 vai chamar por SSH.

Esta tarefa não tem teste automatizado; a verificação é o roteiro de fumaça do Step 6.

- [ ] **Step 1: Marcar o compose atual como de desenvolvimento**

No topo de `docker-compose.yml`, substituir o cabeçalho de comentário por:

```yaml
# BatRadio Web — ambiente de DESENVOLVIMENTO/LOCAL.
#
# NÃO use este arquivo em produção. Ele sobe um backend Node próprio, e no
# servidor da rádio já existe um rodando na porta 9320 — o que o cliente
# Windows do Renato usa. Dois processos falando com o mesmo MPD, ou a 9320
# trocada de dono, derrubam a operação dele.
#
# Produção: docker-compose.main.yml (só o gateway). Ver docs/DEPLOY.md.
```

- [ ] **Step 2: Criar o compose de produção**

`docker-compose.main.yml`:

```yaml
# BatRadio Web — produção (batradio.morcegaofm.com.br).
#
# Sobe SOMENTE o gateway. O backend Node da porta 9320 já roda neste host e é a
# interface do cliente Windows do Renato: este arquivo não o declara, não o
# constrói e não o reinicia.
#
# Uso no servidor:
#   docker compose -f docker-compose.main.yml up -d --wait
#
# Segredos vêm de um .env ao lado deste arquivo (nunca commitado).
services:
  gateway:
    image: ghcr.io/morcegao-fm/batradio-gateway:${TAG:-latest}
    restart: unless-stopped
    env_file: .env
    environment:
      PORT: "8080"
      # Node que JÁ roda no host. NODE_BACKEND_URL e NODE_API_KEY vêm do .env:
      # a chave tem de ser a mesma que o cliente do Renato usa, para nada no
      # Node precisar mudar.
      OAUTH_REDIRECT_URL: "https://batradio.morcegaofm.com.br/auth/callback"
      ALLOWED_EMAILS: "aguergolet@gmail.com,morcegaofm@gmail.com"
      COOKIE_SECURE: "true"
      DEV_MODE: "false"
      # API do website na mesma rede.
      CATALOGO_URL: "${CATALOGO_URL:-http://api:8080}"
    # Sem "ports": quem publica é o Traefik. Nada novo aberto no host.
    expose:
      - "8080"
    networks:
      - telogreug
    healthcheck:
      # O binário é distroless (sem shell, sem curl), então o healthcheck é o
      # próprio gateway: --health-check faz um GET em /healthz e sai 0 ou 1.
      test: ["CMD", "/batradio-gateway", "--health-check"]
      interval: 15s
      timeout: 5s
      retries: 5
      start_period: 20s
    labels:
      - "traefik.enable=true"
      - "traefik.docker.network=telogreug"
      - "traefik.http.routers.batradio.rule=Host(`batradio.morcegaofm.com.br`)"
      - "traefik.http.routers.batradio.entrypoints=websecure,web"
      - "traefik.http.routers.batradio.tls=true"
      - "traefik.http.routers.batradio.tls.certresolver=letsencrypt"
      - "traefik.http.services.batradio.loadbalancer.server.port=8080"

networks:
  telogreug:
    external: true
```

- [ ] **Step 3: Implementar o `--health-check`**

A imagem é distroless: não há `curl` nem shell para o healthcheck. O próprio binário resolve. No topo de `main()`, em `backend-gateway/main.go`, antes de qualquer carga de configuração:

```go
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
```

- [ ] **Step 4: Atualizar o .env.example**

Em `backend-gateway/.env.example`, trocar o valor de exemplo do redirect e acrescentar o bloco novo:

```
OAUTH_REDIRECT_URL=https://batradio.morcegaofm.com.br/auth/callback

# Cookie de sessão com Secure. true em produção; false só em http://localhost.
COOKIE_SECURE=true

# Catálogo de faixas da API do website (capa, artista, ano). Opcional: sem isso
# o painel roda, só que sem capa. A chave é a mesma de Servico__Chave na API.
CATALOGO_URL=http://api:8080
CATALOGO_CHAVE=
```

- [ ] **Step 5: Escrever o docs/DEPLOY.md**

`docs/DEPLOY.md`:

```markdown
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

# 3. A rede do Traefik existe e tem o certresolver letsencrypt:
docker network ls | grep telogreug

# 4. Nada já ocupa o nome do serviço novo:
docker ps -a --format '{{.Names}}' | grep -i batradio
```

`NODE_BACKEND_URL` sai daí:
- Node em container na `telogreug` → `http://<nome-do-container>:9320`.
- Node como processo no host → `http://host.docker.internal:9320` com
  `extra_hosts: ["host.docker.internal:host-gateway"]` no serviço `gateway`.

`NODE_API_KEY` tem de ser **a mesma chave que o cliente do Renato usa**
(`backend-server/config.json` no servidor, campo `webserver.apiKey`). Assim
nada no Node precisa mudar.

## Google OAuth

No Google Cloud Console → APIs & Services → Credentials, no OAuth client ID
do tipo Web application, acrescentar em *Authorized redirect URIs*:

```
https://batradio.morcegaofm.com.br/auth/callback
```

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
```

- [ ] **Step 6: Atualizar o README**

Na seção "Rodando em produção (Docker Compose)" do `README.md`, substituir o passo 4-5 por um ponteiro para `docs/DEPLOY.md` e acrescentar, na descrição de `backend-server/`, a frase:

```
- **`backend-server/`** — backend Node legado que conversa com o MPD. **Em
  produção ele já roda no servidor e é a interface do cliente Windows do
  Renato**: o compose de produção não o declara. Este diretório serve ao
  ambiente de desenvolvimento.
```

E corrigir a descrição do `frontend-windows-client/`, que hoje diz "descontinuado":

```
- **`frontend-windows-client/`** — cliente Windows Forms, **em produção na
  máquina do operador**. O frontend web não o substituiu ainda; o contrato do
  Node na porta 9320 existe para ele.
```

- [ ] **Step 7: Verificar que o build da imagem ainda funciona**

Run: `docker build -f Dockerfile.gateway -t batradio-gateway:teste .`
Expected: build completa sem erro.

Run: `docker run --rm batradio-gateway:teste --health-check; echo "saida=$?"`
Expected: `saida=1` (não há gateway ouvindo dentro desse container efêmero) — o que prova que o modo existe e falha fechado.

- [ ] **Step 8: Commitar**

```bash
git add docker-compose.main.yml docker-compose.yml backend-gateway/.env.example \
        backend-gateway/main.go README.md docs/DEPLOY.md
git commit -m "Compose de produção só com o gateway, atrás do Traefik"
```

---

### Task 10 [B]: CI e deploy automático

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `.github/workflows/deploy.yml`

- [ ] **Step 1: Criar o CI**

`.github/workflows/ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  gateway:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          cache-dependency-path: backend-gateway/go.sum
      - name: Testes do gateway
        working-directory: backend-gateway
        run: go test ./...
      - name: go vet
        working-directory: backend-gateway
        run: go vet ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
          cache: npm
          cache-dependency-path: frontend-web/package-lock.json
      - name: Instalar
        working-directory: frontend-web
        run: npm ci
      - name: Testes
        working-directory: frontend-web
        run: npm test
      - name: Build
        working-directory: frontend-web
        run: npm run build
```

- [ ] **Step 2: Criar o deploy**

`.github/workflows/deploy.yml`:

```yaml
name: Deploy

on:
  workflow_run:
    workflows: [CI]
    types: [completed]
    branches: [main]

permissions:
  contents: read
  packages: write

jobs:
  publicar:
    if: github.event.workflow_run.conclusion == 'success'
    runs-on: ubuntu-latest
    outputs:
      tag: ${{ github.event.workflow_run.head_sha }}
    steps:
      - uses: actions/checkout@v4
        with:
          ref: ${{ github.event.workflow_run.head_sha }}

      - uses: docker/setup-buildx-action@v3

      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Build local primeiro: o Trivy precisa aprovar antes de a imagem existir
      # no registry. Publicar e escanear depois seria publicar vulnerabilidade.
      - name: Build da imagem
        uses: docker/build-push-action@v6
        with:
          context: .
          file: Dockerfile.gateway
          load: true
          tags: batradio-gateway:${{ github.event.workflow_run.head_sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Trivy
        uses: aquasecurity/trivy-action@0.28.0
        with:
          image-ref: batradio-gateway:${{ github.event.workflow_run.head_sha }}
          severity: CRITICAL,HIGH
          exit-code: '1'
          ignore-unfixed: true

      - name: Publicar no GHCR
        run: |
          docker tag batradio-gateway:${{ github.event.workflow_run.head_sha }} \
            ghcr.io/morcegao-fm/batradio-gateway:${{ github.event.workflow_run.head_sha }}
          docker push ghcr.io/morcegao-fm/batradio-gateway:${{ github.event.workflow_run.head_sha }}

  implantar:
    needs: publicar
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Deploy por SSH, com rollback
        uses: appleboy/ssh-action@v1.2.0
        with:
          host: ${{ secrets.DEPLOY_SSH_HOST }}
          username: ${{ secrets.DEPLOY_SSH_USER }}
          key: ${{ secrets.DEPLOY_SSH_KEY }}
          port: ${{ secrets.DEPLOY_SSH_PORT || 22 }}
          script: |
            set -euo pipefail
            cd ~/apps/batradio

            # Guarda a tag atual para o rollback.
            ANTERIOR=$(grep '^TAG=' .env | cut -d= -f2 || echo "")
            NOVA=${{ needs.publicar.outputs.tag }}

            git pull --ff-only
            sed -i "s/^TAG=.*/TAG=$NOVA/" .env || echo "TAG=$NOVA" >> .env

            docker compose -f docker-compose.main.yml pull gateway

            if ! docker compose -f docker-compose.main.yml up -d --wait gateway; then
              echo "Healthcheck falhou; voltando para $ANTERIOR"
              sed -i "s/^TAG=.*/TAG=$ANTERIOR/" .env
              docker compose -f docker-compose.main.yml up -d --wait gateway
              exit 1
            fi

            # O caminho do Renato não pode ter sido tocado por este deploy.
            docker ps --format '{{.Names}}' | grep -qi 'node' && \
              echo "Node segue de pé" || echo "AVISO: nenhum container Node visível"
```

- [ ] **Step 3: Documentar os secrets no DEPLOY.md**

Acrescentar ao fim de `docs/DEPLOY.md`:

```markdown
## Secrets do GitHub

Em Settings → Environments → `production` (o mesmo padrão do website):

| Secret | Valor |
|---|---|
| `DEPLOY_SSH_HOST` | IP do servidor do playout |
| `DEPLOY_SSH_USER` | usuário SSH |
| `DEPLOY_SSH_KEY` | chave privada dedicada ao Actions (não a pessoal) |
| `DEPLOY_SSH_PORT` | porta SSH, se não for 22 (opcional) |

A imagem vai para o GHCR com o `GITHUB_TOKEN` do próprio workflow — nenhuma
credencial de registry precisa ser criada.
```

- [ ] **Step 4: Validar a sintaxe dos workflows**

Run: `python3 -c "import yaml,sys; [yaml.safe_load(open(f)) for f in ['.github/workflows/ci.yml','.github/workflows/deploy.yml']]; print('ok')"`
Expected: `ok`

- [ ] **Step 5: Commitar**

```bash
git add .github/workflows/ docs/DEPLOY.md
git commit -m "CI e deploy automático com Trivy, GHCR e rollback"
```

---

### Task 11 [B]: Skill `batradio-limites`

**Files:**
- Create: `.claude/skills/batradio-limites/SKILL.md`
- Modify: `CLAUDE.md` (criar, se não existir)

- [ ] **Step 1: Escrever a skill**

`.claude/skills/batradio-limites/SKILL.md`:

```markdown
---
name: batradio-limites
description: Use ao mexer em QUALQUER coisa do BatRadio — gateway Go, SPA, compose, deploy, backend Node ou cliente Windows. Define o que não pode ser tocado (a porta 9320 do operador), a regra de autenticação de toda rota nova e a fronteira com o website. Dispara em "batradio", "gateway", "9320", "node-backend", "MPD", "deploy", "compose", "Renato", "cliente Windows", "capa", "catálogo".
---

# Limites do BatRadio

Três regras que o código não conta sozinho, e cujo custo de violar é a rádio
fora do ar ou o painel de controle exposto.

## 1. A porta 9320 é de outra pessoa

O backend Node (`backend-server/`) roda em produção no host do playout, na
porta 9320, autenticado pelo header `batradio-apikey` em HTTP puro. **O cliente
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

Que a 9320 esteja na internet em HTTP é dívida conhecida, registrada em
`docs/superpowers/specs/2026-09-13-batradio-deploy-seguro-design.md`. Resolver
isso é um projeto com o operador na sala, não um refactor de passagem.

## 2. Toda rota nasce autenticada

O BatRadio controla o que está no ar. Quem entra nele muda a transmissão.

- Toda rota nova em `/api` vai atrás de `auth.RequireSession`. Sem exceção
  "temporária".
- A lista de rotas públicas é fechada: `/healthz`, `/auth/login`,
  `/auth/callback`, `/auth/logout` e os estáticos do SPA. Ela vive em
  `internal/httpapi/superficie_test.go`, e `TestNenhumaRotaNovaSemSessao`
  falha se alguém acrescentar rota desprotegida. **Se esse teste ficar
  vermelho, proteja a rota — não relaxe a allowlist.**
- `/healthz` é público e deliberadamente vazio. Não acrescente versão, host do
  Node nem contagem de acervo nele.
- `DEV_MODE=true` é login falso. O `config.Load` recusa subir se ele vier com
  `OAUTH_REDIRECT_URL` https. Não contorne isso.
- Acesso é a allowlist de `ALLOWED_EMAILS`, revalidada a cada request. Não
  existe papel, convite nem cadastro.

## 3. O website é dono dos dados de faixa; o BatRadio só lê

Capa, artista, título e ano vivem na tabela `musica_catalogo` da API .NET, e
são corrigidos à mão na Batcaverna.

- O BatRadio lê por `POST /api/servico/musicas/lookup` com header
  `X-Servico-Key`. **Somente leitura.** Não escreve no banco do website, não
  cria linha, não chama o iTunes.
- Capa é enfeite; o painel é operação. `internal/catalogo` nunca devolve erro:
  site fora do ar vira placeholder, e o controle do MPD segue.
- Não enriqueça `/api/library`: o acervo tem dezenas de milhares de faixas e a
  maioria não tem linha no catálogo.

## Confusão que já custou tempo

`batbelt` e o node do BatRadio são **serviços diferentes** no mesmo host:

| | batbelt | node do BatRadio |
|---|---|---|
| Porta | 8003 | 9320 |
| Auth | nenhuma | header `batradio-apikey` |
| Playlist | `GET /playlist` | `POST /playlist` |
| Consumidor | API do website | cliente Windows do operador |

## Segredos

`SESSION_SECRET`, `NODE_API_KEY`, `CATALOGO_CHAVE` e `Servico:Chave` vêm de
env. Nunca em código, em log, em teste ou em mensagem de erro devolvida ao
browser — inclusive na hora de logar falha de autenticação.
```

- [ ] **Step 2: Criar o CLAUDE.md apontando para a skill**

`CLAUDE.md` na raiz do repo:

```markdown
# Claude Code — BatRadio

Leia `.claude/skills/batradio-limites/SKILL.md` **antes** de mexer em qualquer
coisa deste repositório. Ele registra o que não pode ser tocado (a porta 9320,
que é a interface do cliente Windows do operador), a regra de que toda rota
nova nasce autenticada, e a fronteira com a API do website.

Desenho e plano desta rodada:
- `docs/superpowers/specs/2026-09-13-batradio-deploy-seguro-design.md`
- `docs/superpowers/plans/2026-09-13-batradio-deploy-seguro.md`
- `docs/DEPLOY.md`
```

- [ ] **Step 3: Conferir que o frontmatter é válido**

Run: `head -5 .claude/skills/batradio-limites/SKILL.md`
Expected: bloco `---` com `name:` e `description:` numa linha só cada.

- [ ] **Step 4: Commitar**

```bash
git add .claude/skills/batradio-limites/SKILL.md CLAUDE.md
git commit -m "Skill com os limites do BatRadio: 9320, autenticação e fronteira com o site"
```

---

## Ordem e dependências

```
Task 1 [W]  ─────────────────────────────┐
                                          ▼
Task 2 → Task 3 → Task 4 → Task 5 → Task 6 → Task 7 → Task 8 → Task 9 → Task 10 → Task 11
         (config)  (cookie) (borda) (superfície) (auditoria) (capa)  (UI)  (compose) (CI/CD) (skill)
```

Task 1 é independente e pode ir em paralelo. Tasks 2–6 são o endurecimento e
precisam estar prontas antes de o painel ficar acessível na internet (Task 9).
Task 7 depende da Task 1 estar **implantada** no website para funcionar de
verdade, mas os testes dela usam servidor fake e não bloqueiam.

## O que só você pode fazer

Passos fora do alcance do código, para agendar:

1. **Google Cloud Console** — acrescentar
   `https://batradio.morcegaofm.com.br/auth/callback` às Authorized redirect
   URIs do OAuth client.
2. **DNS** — apontar `batradio.morcegaofm.com.br` para o IP do servidor.
3. **Servidor** — rodar o levantamento de `docs/DEPLOY.md`, preencher o `.env`
   (incluindo a `batradio-apikey` que o operador já usa) e fazer a primeira
   subida manual.
4. **GitHub** — criar o Environment `production` com os secrets de SSH.
5. **Renato** — avisar antes da primeira subida e confirmar com ele que o
   cliente continuou funcionando depois.
