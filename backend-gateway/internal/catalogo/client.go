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
			// como "não tem", senão seria reperguntado a cada render.
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
