package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/schedule"
)

type queuePage struct {
	Items      []schedule.QueueItem `json:"items"`
	Total      int                  `json:"total"`
	CurrentPos int                  `json:"currentPos"`
	Version    int64                `json:"version"`
}

func TestGetPlaylistWithTimes(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	p := doJSON[queuePage](t, s, authedReq("GET", "/api/playlist", ""), 200)
	if p.Total != 3 || len(p.Items) != 3 || p.CurrentPos != 1 {
		t.Fatalf("page: total=%d len=%d current=%d", p.Total, len(p.Items), p.CurrentPos)
	}
	// fila: 3x60s, atual pos=1 elapsed=30 → pos1 = testNow-30s, pos0 = pos1-60s, pos2 = pos1+60s
	anchor := testNow.Add(-30 * 1e9)
	if !p.Items[1].NextPresentation.Equal(anchor) {
		t.Errorf("pos1 time: %v want %v", p.Items[1].NextPresentation, anchor)
	}
	if !p.Items[2].NextPresentation.Equal(anchor.Add(60 * 1e9)) {
		t.Errorf("pos2 time: %v", p.Items[2].NextPresentation)
	}
}

func TestGetPlaylistFiltered(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	p := doJSON[queuePage](t, s, authedReq("GET", "/api/playlist?q=faixa+2", ""), 200)
	if p.Total != 1 || p.Items[0].Title != "Faixa 2" {
		t.Fatalf("filtered: %+v", p)
	}
}

func TestPlaylistCurrent(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	cur := doJSON[map[string]int](t, s, authedReq("GET", "/api/playlist/current", ""), 200)
	if cur["pos"] != 1 || cur["index"] != 1 {
		t.Fatalf("current: %v", cur)
	}
}

func TestPlayerActions(t *testing.T) {
	f, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	for _, action := range []string{"toggle", "shuffle", "repeat", "crossfade"} {
		doJSON[map[string]any](t, s, authedReq("POST", "/api/player/"+action, ""), 200)
	}
	doJSON[map[string]any](t, s, authedReq("POST", "/api/player/play", `{"position":2}`), 200)
	joined := strings.Join(f.calls, "\n")
	for _, want := range []string{"POST /playorpause", "POST /shuffle", "POST /repeat", "POST /fadein", "POST /play position=2"} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing call %q in:\n%s", want, joined)
		}
	}
}

func TestPlayValidation(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, authedReq("POST", "/api/player/play", `{"position":-1}`))
	if rec.Code != 400 {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestAddBumpsVersion(t *testing.T) {
	f, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	v0 := s.poller.Version()
	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, authedReq("POST", "/api/playlist/add", `{"files":["x.mp3"],"position":2}`))
	if rec.Code != 204 {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
	if s.poller.Version() == v0 {
		t.Error("version not bumped")
	}
	joined := strings.Join(f.calls, "\n")
	if !strings.Contains(joined, `POST /addtoplaylist files=["x.mp3"] position=2`) {
		t.Errorf("calls:\n%s", joined)
	}
}

func TestRemoveDeletesDescending(t *testing.T) {
	f, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, authedReq("POST", "/api/playlist/remove", `{"positions":[0,2,1]}`))
	if rec.Code != 204 {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
	var deletes []string
	for _, c := range f.calls {
		if strings.Contains(c, "/delete") {
			deletes = append(deletes, c)
		}
	}
	want := []string{`POST /delete files=["2"]`, `POST /delete files=["1"]`, `POST /delete files=["0"]`}
	if len(deletes) != 3 {
		t.Fatalf("deletes: %v", deletes)
	}
	for i, w := range want {
		if deletes[i] != w {
			t.Errorf("delete[%d]=%q want %q", i, deletes[i], w)
		}
	}
}

func TestMove(t *testing.T) {
	f, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, authedReq("POST", "/api/playlist/move", `{"from":0,"to":2}`))
	if rec.Code != 204 {
		t.Fatalf("code=%d", rec.Code)
	}
	if !strings.Contains(strings.Join(f.calls, "\n"), "POST /move from-pos=0 to-pos=2") {
		t.Errorf("calls: %v", f.calls)
	}
}

func TestAddManyInsertsAtPlannedPositions(t *testing.T) {
	f, srv := newFakeNode(t)
	// fila longa: 24 faixas de 1h → alvos now, +4h, ... viram posições 1,5,9,...
	f.mu.Lock()
	f.queue = nil
	for i := 0; i < 24; i++ {
		f.queue = append(f.queue, map[string]string{
			"file": fmt.Sprintf("h%02d.mp3", i), "Title": fmt.Sprintf("Hora %d", i),
			"Time": "3600", "Pos": fmt.Sprint(i), "Id": fmt.Sprint(i + 100),
		})
	}
	f.status.Song = 0
	f.status.Elapsed = 0
	f.mu.Unlock()

	s := newTestServer(t, srv.URL)
	body := `{"file":"vinheta.mp3","timesPerDay":6,"intervalHours":4}`
	res := doJSON[map[string]any](t, s, authedReq("POST", "/api/playlist/add-many", body), 200)
	arr := res["insertedAt"].([]any)
	if len(arr) != 6 {
		t.Fatalf("insertedAt: %v", arr)
	}
	// alvo k=0 é exatamente o horário da faixa atual (pos 0) → posição 0; demais a cada 4 faixas de 1h
	want := []float64{20, 16, 12, 8, 4, 0}
	for i, w := range want {
		if arr[i].(float64) != w {
			t.Errorf("insertedAt[%d]=%v want %v (all=%v)", i, arr[i], w, arr)
		}
	}
	// e o Node recebeu 6 inserções na ordem decrescente
	var adds []string
	for _, c := range f.calls {
		if strings.Contains(c, "/addtoplaylist") {
			adds = append(adds, c)
		}
	}
	if len(adds) != 6 || !strings.Contains(adds[0], "position=20") || !strings.Contains(adds[5], "position=0") {
		t.Errorf("adds: %v", adds)
	}
}

func TestAddManyValidation(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	for _, body := range []string{
		`{"file":"","timesPerDay":6,"intervalHours":4}`,
		`{"file":"x.mp3","timesPerDay":0,"intervalHours":4}`,
		`{"file":"x.mp3","timesPerDay":100,"intervalHours":4}`,
		`{"file":"x.mp3","timesPerDay":6,"intervalHours":0}`,
	} {
		rec := httptest.NewRecorder()
		s.Routes(nil).ServeHTTP(rec, authedReq("POST", "/api/playlist/add-many", body))
		if rec.Code != 400 {
			t.Errorf("body %s: code=%d", body, rec.Code)
		}
	}
}

func TestSavedPlaylists(t *testing.T) {
	f, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)

	list := doJSON[map[string][]map[string]string](t, s, authedReq("GET", "/api/playlists", ""), 200)
	if len(list["items"]) != 1 || list["items"][0]["name"] != "Programação Padrão" {
		t.Fatalf("list: %v", list)
	}

	rec := httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, authedReq("POST", "/api/playlists", `{"name":"Nova"}`))
	if rec.Code != 204 {
		t.Fatalf("save: code=%d", rec.Code)
	}

	// load com playAtRandom → randInt fixo em 1 → play position=1... na fila de 3, 1+randInt(2)=2? ver implementação: pos = 1+randInt(len-2)? Contrato: 1 <= pos < len.
	st := doJSON[map[string]any](t, s, authedReq("POST", "/api/playlists/Cl%C3%A1ssicos/load", `{"playAtRandom":true}`), 200)
	if st["state"] != "play" {
		t.Fatalf("load: %v", st)
	}
	joined := strings.Join(f.calls, "\n")
	if !strings.Contains(joined, "POST /loadplaylist name=Clássicos") {
		t.Errorf("load call missing: %s", joined)
	}
	if !strings.Contains(joined, "POST /play position=") {
		t.Errorf("play call missing after random load: %s", joined)
	}

	rec = httptest.NewRecorder()
	s.Routes(nil).ServeHTTP(rec, authedReq("DELETE", "/api/playlists/Velha", ""))
	if rec.Code != 204 {
		t.Fatalf("delete: code=%d", rec.Code)
	}
	if !strings.Contains(strings.Join(f.calls, "\n"), "POST /removeplaylist name=Velha") {
		t.Errorf("remove call missing")
	}
}

func TestSavePlaylistValidation(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	for _, body := range []string{`{"name":""}`, `{"name":"a/b"}`} {
		rec := httptest.NewRecorder()
		s.Routes(nil).ServeHTTP(rec, authedReq("POST", "/api/playlists", body))
		if rec.Code != 400 {
			t.Errorf("body %s: code=%d", body, rec.Code)
		}
	}
}

func TestMe(t *testing.T) {
	_, srv := newFakeNode(t)
	s := newTestServer(t, srv.URL)
	me := doJSON[map[string]string](t, s, authedReq("GET", "/api/me", ""), 200)
	if me["email"] != testEmail {
		t.Fatalf("me: %v", me)
	}
}

var _ = json.Marshal // silencia import se não usado em edits futuros
