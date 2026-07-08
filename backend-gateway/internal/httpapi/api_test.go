package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/auth"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/config"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/node"
)

var (
	testSecret = []byte(strings.Repeat("k", 32))
	testEmail  = "aguergolet@gmail.com"
	testNow    = time.Date(2026, 7, 8, 14, 0, 0, 0, time.UTC)
)

// fakeNode simula o backend Node com uma fila mutável em memória.
type fakeNode struct {
	mu     sync.Mutex
	t      *testing.T
	status model.Status
	queue  []map[string]string // payload cru estilo MPD
	files  string              // payload canned do acervo
	calls  []string            // "METHOD path key=val ..."
	down   bool
}

func (f *fakeNode) record(r *http.Request, keys ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	entry := r.Method + " " + r.URL.Path
	for _, k := range keys {
		entry += " " + k + "=" + r.Header.Get(k)
	}
	f.calls = append(f.calls, entry)
}

func (f *fakeNode) statusJSON() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	cur := ""
	for _, s := range f.queue {
		if s["Pos"] == fmt.Sprint(f.status.Song) {
			b, _ := json.Marshal(s)
			cur = string(b)
		}
	}
	repeat := "0"
	if f.status.Repeat {
		repeat = "1"
	}
	j := fmt.Sprintf(`{"state":%q,"song":"%d","elapsed":"%v","repeat":%q,"random":"0","xfade":"0","playlistlength":"%d"`,
		f.status.State, f.status.Song, f.status.Elapsed, repeat, len(f.queue))
	if cur != "" {
		j += `,"currentSong":` + cur
	}
	return j + "}"
}

func (f *fakeNode) queueJSON() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, _ := json.Marshal(f.queue)
	return string(b)
}

func newFakeNode(t *testing.T) (*fakeNode, *httptest.Server) {
	f := &fakeNode{t: t, status: model.Status{State: "play", Song: 1, Elapsed: 30}}
	// fila: 3 faixas de 60s
	for i := 0; i < 3; i++ {
		f.queue = append(f.queue, map[string]string{
			"file": fmt.Sprintf("s%d.mp3", i), "Artist": fmt.Sprintf("Artista %d", i),
			"Title": fmt.Sprintf("Faixa %d", i), "Time": "60",
			"Pos": fmt.Sprint(i), "Id": fmt.Sprint(i + 10),
		})
	}
	f.files = `[{"file":"rock/led/whole_lotta_love.mp3","Artist":"Led Zeppelin","Title":"Whole Lotta Love","Time":"333","Pos":"0","Id":"1"},
		{"file":"blues/bb/thrill.mp3","Artist":"B.B. King","Title":"The Thrill Is Gone","Time":"324","Pos":"1","Id":"2"}]`

	mux := http.NewServeMux()
	guard := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if f.down {
				w.WriteHeader(500)
				w.Write([]byte(`{"message":"down"}`))
				return
			}
			h(w, r)
		}
	}
	mux.HandleFunc("GET /status", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r)
		w.Write([]byte(f.statusJSON()))
	}))
	mux.HandleFunc("GET /list", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "type", "update")
		if r.Header.Get("type") == "playlist" {
			w.Write([]byte(`[{"playlist":"Programação Padrão","Last-Modified":"2026-07-04T10:00:00Z"}]`))
			return
		}
		w.Write([]byte(f.files))
	}))
	mux.HandleFunc("POST /playlist", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "update")
		w.Write([]byte(f.queueJSON()))
	}))
	mux.HandleFunc("POST /addtoplaylist", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "files", "position")
		w.Write([]byte(f.queueJSON()))
	}))
	mux.HandleFunc("POST /delete", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "files")
		w.Write([]byte(f.queueJSON()))
	}))
	mux.HandleFunc("POST /move", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "from-pos", "to-pos")
		w.Write([]byte(f.queueJSON()))
	}))
	for _, p := range []string{"/playorpause", "/repeat", "/shuffle", "/fadein"} {
		mux.HandleFunc("POST "+p, guard(func(w http.ResponseWriter, r *http.Request) {
			f.record(r)
			w.Write([]byte(f.statusJSON()))
		}))
	}
	mux.HandleFunc("POST /play", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "position")
		w.Write([]byte(f.statusJSON()))
	}))
	mux.HandleFunc("POST /loadplaylist", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "name")
		w.Write([]byte(f.queueJSON()))
	}))
	mux.HandleFunc("POST /saveplaylist", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "name")
		w.Write([]byte(f.queueJSON()))
	}))
	mux.HandleFunc("POST /removeplaylist", guard(func(w http.ResponseWriter, r *http.Request) {
		f.record(r, "name")
		w.Write([]byte(f.queueJSON()))
	}))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return f, srv
}

func newTestServer(t *testing.T, nodeURL string) *Server {
	cfg := &config.Config{
		Port:           "8080",
		StreamURL:      "https://stream.example/live",
		AllowedEmails:  []string{testEmail},
		SessionSecret:  testSecret,
		NodeBackendURL: nodeURL,
		NodeAPIKey:     "key",
	}
	s := NewServer(cfg, node.New(nodeURL, "key"))
	s.now = func() time.Time { return testNow }
	s.randInt = func(n int) int { return 1 }
	return s
}

func authedReq(method, target string, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	tok := auth.SignSession(testEmail, time.Now().Add(time.Hour), testSecret)
	req.AddCookie(&http.Cookie{Name: auth.SessionCookie, Value: tok})
	return req
}

func doJSON[T any](t *testing.T, s *Server, req *http.Request, wantCode int) T {
	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, req)
	if rec.Code != wantCode {
		t.Fatalf("%s %s: code=%d body=%s", req.Method, req.URL, rec.Code, rec.Body.String())
	}
	var out T
	if rec.Body.Len() > 0 {
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatalf("decode: %v body=%s", err, rec.Body.String())
		}
	}
	return out
}

func TestUnauthenticated(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, httptest.NewRequest("GET", "/api/status", nil))
	if rec.Code != 401 {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestStatus(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	st := doJSON[model.Status](t, s, authedReq("GET", "/api/status", ""), 200)
	if st.State != "play" || st.Song != 1 || st.PlaylistLength != 3 {
		t.Fatalf("status: %+v", st)
	}
}

func TestLibraryRefreshAndSearch(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	res := doJSON[map[string]any](t, s, authedReq("POST", "/api/library/refresh", ""), 200)
	if res["count"].(float64) != 2 {
		t.Fatalf("count=%v", res["count"])
	}
	type page struct {
		Items []model.Song `json:"items"`
		Total int          `json:"total"`
	}
	p := doJSON[page](t, s, authedReq("GET", "/api/library?q=zeppelin+whole&offset=0&limit=50", ""), 200)
	if p.Total != 1 || len(p.Items) != 1 || p.Items[0].Title != "Whole Lotta Love" {
		t.Fatalf("page: %+v", p)
	}
	p = doJSON[page](t, s, authedReq("GET", "/api/library", ""), 200)
	if p.Total != 2 {
		t.Fatalf("all: %+v", p)
	}
}

func TestServerInfo(t *testing.T) {
	f, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	// popula acervo e status
	doJSON[map[string]any](t, s, authedReq("POST", "/api/library/refresh", ""), 200)
	s.poller.Poll(context.Background())
	info := doJSON[map[string]any](t, s, authedReq("GET", "/api/server/info", ""), 200)
	if info["nodeOk"] != true || info["libraryCount"].(float64) != 2 || info["streamUrl"] != "https://stream.example/live" {
		t.Fatalf("info: %v", info)
	}
	f.down = true
	s.poller.Poll(context.Background())
	info = doJSON[map[string]any](t, s, authedReq("GET", "/api/server/info", ""), 200)
	if info["nodeOk"] != false {
		t.Fatalf("expected nodeOk=false: %v", info)
	}
}

func TestNodeDownIs502(t *testing.T) {
	f, srv := newFakeNode(t)
	f.down = true
	s := newTestServer(t, srv.URL)
	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, authedReq("GET", "/api/status", ""))
	if rec.Code != 502 {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
	var e struct {
		Error struct{ Code string } `json:"error"`
	}
	json.Unmarshal(rec.Body.Bytes(), &e)
	if e.Error.Code != "node_unavailable" {
		t.Fatalf("error code=%q", e.Error.Code)
	}
}

func TestSSEEmitsStatusAndPlaylist(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := authedReq("GET", "/api/events", "").WithContext(ctx)
	rec := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		s.Routes(nil).ServeHTTP(rec, req)
		close(done)
	}()
	time.Sleep(50 * time.Millisecond) // deixa o handler assinar
	s.poller.Poll(context.Background())
	s.poller.BumpPlaylist()
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	body := rec.Body.String()
	if !strings.Contains(body, "event: status\n") || !strings.Contains(body, `"state":"play"`) {
		t.Errorf("missing status event: %q", body)
	}
	if !strings.Contains(body, "event: playlist\n") || !strings.Contains(body, `"version":`) {
		t.Errorf("missing playlist event: %q", body)
	}
}

func TestPollerDetectsExternalQueueChange(t *testing.T) {
	f, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	s.poller.Poll(context.Background())
	v1 := s.poller.Version()
	// muda a fila "por fora" (outro cliente MPD)
	f.mu.Lock()
	f.queue = f.queue[:2]
	f.mu.Unlock()
	s.poller.Poll(context.Background())
	if s.poller.Version() == v1 {
		t.Fatal("version should bump when playlistlength changes externally")
	}
}
