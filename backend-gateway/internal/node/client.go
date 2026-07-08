// Package node é o cliente HTTP do backend Node legado (server.js), que fala o
// protocolo antigo: parâmetros via headers e autenticação batradio-apikey.
package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
)

// ErrNodeUnavailable indica falha de rede ou resposta não-200 do Node.
var ErrNodeUnavailable = errors.New("backend da rádio indisponível")

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		// O Node usa timeout de 30s nas próprias operações MPD.
		http: &http.Client{Timeout: 35 * time.Second},
	}
}

// busyRetries e busyBackoff controlam as novas tentativas quando o Node
// responde "MPD is busy" (trava global dele durante scans do acervo/fila).
var (
	busyRetries = 4
	busyBackoff = 700 * time.Millisecond
)

func (c *Client) do(ctx context.Context, method, path string, headers map[string]string) ([]byte, error) {
	var lastErr error
	for attempt := 0; ; attempt++ {
		body, busy, err := c.doOnce(ctx, method, path, headers)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !busy || attempt >= busyRetries {
			return nil, lastErr
		}
		select {
		case <-ctx.Done():
			return nil, lastErr
		case <-time.After(busyBackoff * time.Duration(attempt+1)):
		}
	}
}

func (c *Client) doOnce(ctx context.Context, method, path string, headers map[string]string) (body []byte, busy bool, err error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("batradio-apikey", c.apiKey)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrNodeUnavailable, err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("%w: %v", ErrNodeUnavailable, err)
	}
	if resp.StatusCode != http.StatusOK {
		msg := string(body)
		var e struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(body, &e) == nil && e.Message != "" {
			msg = e.Message
		}
		busy = strings.Contains(strings.ToLower(msg), "busy")
		return nil, busy, fmt.Errorf("%w: %s (HTTP %d)", ErrNodeUnavailable, msg, resp.StatusCode)
	}
	return body, false, nil
}

// mpdStatus espelha o JSON cru do Node para /status.
type mpdStatus struct {
	State          string          `json:"state"`
	Song           model.FlexInt   `json:"song"`
	Elapsed        model.FlexFloat `json:"elapsed"`
	Repeat         model.FlexInt   `json:"repeat"`
	Random         model.FlexInt   `json:"random"`
	Xfade          model.FlexInt   `json:"xfade"`
	PlaylistLength model.FlexInt   `json:"playlistlength"`
	CurrentSong    json.RawMessage `json:"currentSong"`
	Next           json.RawMessage `json:"next"`
}

func parseStatus(raw []byte) (model.Status, error) {
	var ms mpdStatus
	if err := json.Unmarshal(raw, &ms); err != nil {
		return model.Status{}, fmt.Errorf("status inválido do Node: %w", err)
	}
	st := model.Status{
		State:          ms.State,
		Song:           int(ms.Song),
		Elapsed:        float64(ms.Elapsed),
		Repeat:         ms.Repeat == 1,
		Random:         ms.Random == 1,
		Crossfade:      ms.Xfade > 0,
		PlaylistLength: int(ms.PlaylistLength),
	}
	for _, pair := range []struct {
		raw  json.RawMessage
		dest **model.Song
	}{{ms.CurrentSong, &st.Current}, {ms.Next, &st.Next}} {
		if len(pair.raw) == 0 || string(pair.raw) == "null" {
			continue
		}
		if s, err := model.ParseSong(pair.raw); err == nil {
			song := s
			*pair.dest = &song
		}
	}
	return st, nil
}

func (c *Client) statusRequest(ctx context.Context, method, path string, headers map[string]string) (model.Status, error) {
	body, err := c.do(ctx, method, path, headers)
	if err != nil {
		return model.Status{}, err
	}
	return parseStatus(body)
}

func (c *Client) songsRequest(ctx context.Context, method, path string, headers map[string]string) ([]model.Song, error) {
	body, err := c.do(ctx, method, path, headers)
	if err != nil {
		return nil, err
	}
	songs, err := model.ParseSongs(body)
	if err != nil {
		return nil, fmt.Errorf("resposta inválida do Node em %s: %w", path, err)
	}
	return songs, nil
}

func (c *Client) GetStatus(ctx context.Context) (model.Status, error) {
	return c.statusRequest(ctx, http.MethodGet, "/status", nil)
}

func (c *Client) GetPlaylist(ctx context.Context, update bool) ([]model.Song, error) {
	return c.songsRequest(ctx, http.MethodPost, "/playlist", map[string]string{"update": strconv.FormatBool(update)})
}

func (c *Client) GetFiles(ctx context.Context, update bool) ([]model.Song, error) {
	return c.songsRequest(ctx, http.MethodGet, "/list", map[string]string{"type": "file", "update": strconv.FormatBool(update)})
}

func (c *Client) GetPlaylists(ctx context.Context) ([]model.PlaylistInfo, error) {
	body, err := c.do(ctx, http.MethodGet, "/list", map[string]string{"type": "playlist", "update": "false"})
	if err != nil {
		return nil, err
	}
	var raw []struct {
		Playlist     string `json:"playlist"`
		LastModified string `json:"Last-Modified"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("lista de playlists inválida: %w", err)
	}
	pls := make([]model.PlaylistInfo, 0, len(raw))
	for _, p := range raw {
		if p.Playlist == "" {
			continue
		}
		pls = append(pls, model.PlaylistInfo{Name: p.Playlist, LastModified: p.LastModified})
	}
	return pls, nil
}

func (c *Client) PlayOrPause(ctx context.Context) (model.Status, error) {
	return c.statusRequest(ctx, http.MethodPost, "/playorpause", nil)
}

func (c *Client) Play(ctx context.Context, pos int) (model.Status, error) {
	return c.statusRequest(ctx, http.MethodPost, "/play", map[string]string{"position": strconv.Itoa(pos)})
}

func (c *Client) Repeat(ctx context.Context) (model.Status, error) {
	return c.statusRequest(ctx, http.MethodPost, "/repeat", nil)
}

func (c *Client) Shuffle(ctx context.Context) (model.Status, error) {
	return c.statusRequest(ctx, http.MethodPost, "/shuffle", nil)
}

func (c *Client) Crossfade(ctx context.Context) (model.Status, error) {
	return c.statusRequest(ctx, http.MethodPost, "/fadein", nil)
}

func (c *Client) Add(ctx context.Context, files []string, pos int) ([]model.Song, error) {
	filesJSON, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}
	return c.songsRequest(ctx, http.MethodPost, "/addtoplaylist", map[string]string{
		"files":    string(filesJSON),
		"position": strconv.Itoa(pos),
	})
}

// Delete remove faixas pelas POSIÇÕES na fila (o Node repassa como argumento do
// comando MPD delete, e o cliente legado enviava posições).
func (c *Client) Delete(ctx context.Context, positions []int) ([]model.Song, error) {
	strs := make([]string, len(positions))
	for i, p := range positions {
		strs[i] = strconv.Itoa(p)
	}
	filesJSON, err := json.Marshal(strs)
	if err != nil {
		return nil, err
	}
	return c.songsRequest(ctx, http.MethodPost, "/delete", map[string]string{"files": string(filesJSON)})
}

func (c *Client) Move(ctx context.Context, from, to int) ([]model.Song, error) {
	return c.songsRequest(ctx, http.MethodPost, "/move", map[string]string{
		"from-pos": strconv.Itoa(from),
		"to-pos":   strconv.Itoa(to),
	})
}

func (c *Client) LoadPlaylist(ctx context.Context, name string) ([]model.Song, error) {
	return c.songsRequest(ctx, http.MethodPost, "/loadplaylist", map[string]string{"name": name})
}

func (c *Client) SavePlaylist(ctx context.Context, name string) error {
	_, err := c.do(ctx, http.MethodPost, "/saveplaylist", map[string]string{"name": name})
	return err
}

func (c *Client) RemovePlaylist(ctx context.Context, name string) error {
	_, err := c.do(ctx, http.MethodPost, "/removeplaylist", map[string]string{"name": name})
	return err
}
