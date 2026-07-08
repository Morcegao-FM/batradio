package model

import (
	"encoding/json"
	"testing"
)

func TestFlexIntDecodesStringAndNumber(t *testing.T) {
	var v struct {
		A FlexInt `json:"a"`
		B FlexInt `json:"b"`
	}
	if err := json.Unmarshal([]byte(`{"a":"3","b":7}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if int(v.A) != 3 || int(v.B) != 7 {
		t.Errorf("got a=%d b=%d", v.A, v.B)
	}
}

func TestFlexFloatDecodesStringAndNumber(t *testing.T) {
	var v struct {
		A FlexFloat `json:"a"`
		B FlexFloat `json:"b"`
	}
	if err := json.Unmarshal([]byte(`{"a":"10.5","b":2.25}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if float64(v.A) != 10.5 || float64(v.B) != 2.25 {
		t.Errorf("got a=%v b=%v", v.A, v.B)
	}
}

func TestFlexIntTruncatesFloatString(t *testing.T) {
	// MPD pode mandar Time como "352.133"
	var v struct {
		A FlexInt `json:"a"`
	}
	if err := json.Unmarshal([]byte(`{"a":"352.133"}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if int(v.A) != 352 {
		t.Errorf("got %d", v.A)
	}
}

func TestParseSongFromMPDPayload(t *testing.T) {
	raw := []byte(`{"file":"rock/acdc/back_in_black.mp3","Last-Modified":"2020-01-01T00:00:00Z",
		"Artist":"AC/DC","Title":"Back in Black","Album":"Back in Black","Genre":"Rock",
		"Time":"255","duration":"255.032","Pos":"4","Id":"17"}`)
	s, err := ParseSong(raw)
	if err != nil {
		t.Fatalf("ParseSong: %v", err)
	}
	if s.File != "rock/acdc/back_in_black.mp3" || s.Artist != "AC/DC" || s.Title != "Back in Black" {
		t.Errorf("basic fields wrong: %+v", s)
	}
	if s.Time != 255 || s.Pos != 4 || s.ID != 17 || s.Genre != "Rock" {
		t.Errorf("numeric/extra fields wrong: %+v", s)
	}
}

func TestParseSongsMissingFields(t *testing.T) {
	// Arquivos sem tags (ex.: vinhetas) só têm file/Time/Pos/Id.
	raw := []byte(`[{"file":"vinhetas/chacoalhe.mp3","Time":"12","Pos":"0","Id":"1"},
		{"file":"a.mp3","Artist":"X","Time":"60","Pos":"1","Id":"2"}]`)
	songs, err := ParseSongs(raw)
	if err != nil {
		t.Fatalf("ParseSongs: %v", err)
	}
	if len(songs) != 2 {
		t.Fatalf("len=%d", len(songs))
	}
	if songs[0].Artist != "" || songs[0].Time != 12 {
		t.Errorf("song0: %+v", songs[0])
	}
}
