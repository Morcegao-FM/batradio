package node

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testKey = "test-api-key"

// fakeNode grava a última request de cada rota e devolve payloads canned.
type fakeNode struct {
	t        *testing.T
	mux      *http.ServeMux
	requests []recordedReq
}

type recordedReq struct {
	Path    string
	Headers http.Header
}

func newFakeNode(t *testing.T) (*fakeNode, *httptest.Server) {
	f := &fakeNode{t: t, mux: http.NewServeMux()}
	record := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f.requests = append(f.requests, recordedReq{Path: r.URL.Path, Headers: r.Header.Clone()})
			if r.Header.Get("batradio-apikey") != testKey {
				w.WriteHeader(403)
				json.NewEncoder(w).Encode(map[string]string{"message": "Invalid API Key"})
				return
			}
			h(w, r)
		}
	}
	statusPayload := `{"state":"play","song":"2","songid":"12","elapsed":"10.5","repeat":"1",
		"random":"0","xfade":"0","playlistlength":"8",
		"currentSong":{"file":"a.mp3","Artist":"Jimi Hendrix","Title":"Voodoo Child","Time":"312","Pos":"2","Id":"12"},
		"next":{"file":"b.mp3","Artist":"Led Zeppelin","Title":"Black Dog","Time":"294","Pos":"3","Id":"13"}}`
	songsPayload := `[{"file":"a.mp3","Artist":"A","Title":"TA","Time":"100","Pos":"0","Id":"1"},
		{"file":"b.mp3","Artist":"B","Title":"TB","Time":"200","Pos":"1","Id":"2"}]`
	playlistsPayload := `[{"playlist":"Programação Padrão","Last-Modified":"2026-07-04T10:00:00Z"},
		{"playlist":"Blues da Meia-Noite","Last-Modified":"2026-07-01T10:00:00Z"}]`

	f.mux.HandleFunc("GET /status", record(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(statusPayload)) }))
	f.mux.HandleFunc("POST /playlist", record(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(songsPayload)) }))
	f.mux.HandleFunc("GET /list", record(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("type") == "playlist" {
			w.Write([]byte(playlistsPayload))
			return
		}
		w.Write([]byte(songsPayload))
	}))
	for _, p := range []string{"/playorpause", "/play", "/repeat", "/shuffle", "/fadein"} {
		f.mux.HandleFunc("POST "+p, record(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(statusPayload)) }))
	}
	for _, p := range []string{"/addtoplaylist", "/delete", "/move", "/loadplaylist", "/saveplaylist", "/removeplaylist"} {
		f.mux.HandleFunc("POST "+p, record(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(songsPayload)) }))
	}
	srv := httptest.NewServer(f.mux)
	t.Cleanup(srv.Close)
	return f, srv
}

func (f *fakeNode) last() recordedReq {
	if len(f.requests) == 0 {
		f.t.Fatal("no requests recorded")
	}
	return f.requests[len(f.requests)-1]
}

func TestGetStatusParsesAndSendsKey(t *testing.T) {
	f, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	st, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus: %v", err)
	}
	if st.State != "play" || st.Song != 2 || st.Elapsed != 10.5 || !st.Repeat || st.Random || st.Crossfade {
		t.Errorf("status wrong: %+v", st)
	}
	if st.PlaylistLength != 8 || st.Current == nil || st.Current.Title != "Voodoo Child" || st.Next == nil {
		t.Errorf("current/next wrong: %+v", st)
	}
	if got := f.last().Headers.Get("batradio-apikey"); got != testKey {
		t.Errorf("apikey header: %q", got)
	}
}

func TestGetFilesSendsTypeAndUpdate(t *testing.T) {
	f, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	songs, err := c.GetFiles(context.Background(), true)
	if err != nil {
		t.Fatalf("GetFiles: %v", err)
	}
	if len(songs) != 2 {
		t.Fatalf("len=%d", len(songs))
	}
	h := f.last().Headers
	if h.Get("type") != "file" || h.Get("update") != "true" {
		t.Errorf("headers: type=%q update=%q", h.Get("type"), h.Get("update"))
	}
}

func TestGetPlaylistsParsesNames(t *testing.T) {
	_, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	pls, err := c.GetPlaylists(context.Background())
	if err != nil {
		t.Fatalf("GetPlaylists: %v", err)
	}
	if len(pls) != 2 || pls[0].Name != "Programação Padrão" || pls[0].LastModified == "" {
		t.Errorf("playlists wrong: %+v", pls)
	}
}

func TestAddSendsFilesJSONAndPosition(t *testing.T) {
	f, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	if _, err := c.Add(context.Background(), []string{"x.mp3", "y.mp3"}, 5); err != nil {
		t.Fatalf("Add: %v", err)
	}
	h := f.last().Headers
	if h.Get("files") != `["x.mp3","y.mp3"]` || h.Get("position") != "5" {
		t.Errorf("headers: files=%q position=%q", h.Get("files"), h.Get("position"))
	}
}

func TestDeleteSendsPositions(t *testing.T) {
	f, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	if _, err := c.Delete(context.Background(), []int{7, 3}); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got := f.last().Headers.Get("files"); got != `["7","3"]` {
		t.Errorf("files header: %q", got)
	}
}

func TestMoveSendsFromTo(t *testing.T) {
	f, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	if _, err := c.Move(context.Background(), 2, 9); err != nil {
		t.Fatalf("Move: %v", err)
	}
	h := f.last().Headers
	if h.Get("from-pos") != "2" || h.Get("to-pos") != "9" {
		t.Errorf("headers: from=%q to=%q", h.Get("from-pos"), h.Get("to-pos"))
	}
}

func TestPlaySendsPosition(t *testing.T) {
	f, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	if _, err := c.Play(context.Background(), 42); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if got := f.last().Headers.Get("position"); got != "42" {
		t.Errorf("position header: %q", got)
	}
}

func TestPlaylistNameOperations(t *testing.T) {
	f, srv := newFakeNode(t)
	c := New(srv.URL, testKey)
	if _, err := c.LoadPlaylist(context.Background(), "Blues da Meia-Noite"); err != nil {
		t.Fatalf("LoadPlaylist: %v", err)
	}
	if got := f.last().Headers.Get("name"); got != "Blues da Meia-Noite" {
		t.Errorf("name header: %q", got)
	}
	if err := c.SavePlaylist(context.Background(), "Nova"); err != nil {
		t.Fatalf("SavePlaylist: %v", err)
	}
	if err := c.RemovePlaylist(context.Background(), "Velha"); err != nil {
		t.Fatalf("RemovePlaylist: %v", err)
	}
}

func TestServerErrorIsNodeUnavailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"erro interno"}`))
	}))
	t.Cleanup(srv.Close)
	c := New(srv.URL, testKey)
	_, err := c.GetStatus(context.Background())
	if !errors.Is(err, ErrNodeUnavailable) {
		t.Fatalf("expected ErrNodeUnavailable, got %v", err)
	}
}

// O Node tem uma trava global "busy" (ex.: durante o scan do acervo). O
// cliente deve tentar de novo em vez de estourar 502 na cara do usuário.
func TestBusyIsRetried(t *testing.T) {
	oldBackoff := busyBackoff
	busyBackoff = time.Millisecond
	t.Cleanup(func() { busyBackoff = oldBackoff })

	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(500)
			w.Write([]byte(`{"message":"MPD is busy, try again later"}`))
			return
		}
		w.Write([]byte(`{"state":"play","song":"0","playlistlength":"1"}`))
	}))
	t.Cleanup(srv.Close)

	c := New(srv.URL, testKey)
	st, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
	if st.State != "play" || calls != 3 {
		t.Fatalf("state=%q calls=%d", st.State, calls)
	}
}

func TestBusyGivesUpAfterRetries(t *testing.T) {
	oldBackoff := busyBackoff
	busyBackoff = time.Millisecond
	t.Cleanup(func() { busyBackoff = oldBackoff })

	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"MPD is busy, try again later"}`))
	}))
	t.Cleanup(srv.Close)

	c := New(srv.URL, testKey)
	_, err := c.GetStatus(context.Background())
	if !errors.Is(err, ErrNodeUnavailable) {
		t.Fatalf("expected ErrNodeUnavailable, got %v", err)
	}
	if calls != busyRetries+1 {
		t.Fatalf("calls=%d want %d", calls, busyRetries+1)
	}
}

func TestNonBusyErrorIsNotRetried(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(500)
		w.Write([]byte(`{"message":"boom"}`))
	}))
	t.Cleanup(srv.Close)

	c := New(srv.URL, testKey)
	_, err := c.GetStatus(context.Background())
	if !errors.Is(err, ErrNodeUnavailable) {
		t.Fatalf("expected ErrNodeUnavailable, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("calls=%d want 1 (sem retry)", calls)
	}
}

func TestConnectionRefusedIsNodeUnavailable(t *testing.T) {
	c := New("http://127.0.0.1:1", testKey)
	_, err := c.GetStatus(context.Background())
	if !errors.Is(err, ErrNodeUnavailable) {
		t.Fatalf("expected ErrNodeUnavailable, got %v", err)
	}
}
