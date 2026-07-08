package library

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
)

func sampleSongs() []model.Song {
	return []model.Song{
		{File: "rock/led_zeppelin/whole_lotta_love.mp3", Artist: "Led Zeppelin", Title: "Whole Lotta Love", Album: "Led Zeppelin II", Genre: "Rock", Time: 333},
		{File: "rock/acdc/back_in_black.mp3", Artist: "AC/DC", Title: "Back in Black", Album: "Back in Black", Genre: "Rock", Time: 255},
		{File: "blues/bb_king/the_thrill_is_gone.mp3", Artist: "B.B. King", Title: "The Thrill Is Gone", Genre: "Blues", Time: 324},
		{File: "vinhetas/chacoalhe.mp3", Title: "Vinheta — Chacoalhe a Caverna", Artist: "Morcegão FM", Genre: "Vinheta", Time: 12},
	}
}

func TestSearchAllTermsMustMatch(t *testing.T) {
	idx := NewIndex()
	idx.Set(sampleSongs(), time.Now())

	items, total := idx.Search("zeppelin whole", 0, 50)
	if total != 1 || len(items) != 1 || items[0].Title != "Whole Lotta Love" {
		t.Fatalf("got total=%d items=%v", total, items)
	}
	// Termo presente + termo ausente → exclui
	_, total = idx.Search("zeppelin caverna", 0, 50)
	if total != 0 {
		t.Fatalf("expected 0, got %d", total)
	}
}

func TestSearchCaseInsensitiveAcrossFields(t *testing.T) {
	idx := NewIndex()
	idx.Set(sampleSongs(), time.Now())
	// bate no File
	_, total := idx.Search("BACK_IN_BLACK", 0, 50)
	if total != 1 {
		t.Fatalf("file match: got %d", total)
	}
	// bate no Album
	_, total = idx.Search("led zeppelin ii", 0, 50)
	if total != 1 {
		t.Fatalf("album match: got %d", total)
	}
	// bate no Genre
	_, total = idx.Search("vinheta", 0, 50)
	if total != 1 {
		t.Fatalf("genre match: got %d", total)
	}
}

func TestEmptyQueryReturnsAllSortedByFile(t *testing.T) {
	idx := NewIndex()
	idx.Set(sampleSongs(), time.Now())
	items, total := idx.Search("", 0, 50)
	if total != 4 || len(items) != 4 {
		t.Fatalf("total=%d len=%d", total, len(items))
	}
	for i := 1; i < len(items); i++ {
		if items[i-1].File > items[i].File {
			t.Fatalf("not sorted by file: %q > %q", items[i-1].File, items[i].File)
		}
	}
}

func TestPagination(t *testing.T) {
	songs := make([]model.Song, 250)
	for i := range songs {
		songs[i] = model.Song{File: fmt.Sprintf("f%04d.mp3", i), Title: fmt.Sprintf("T%d", i)}
	}
	idx := NewIndex()
	idx.Set(songs, time.Now())

	items, total := idx.Search("", 100, 50)
	if total != 250 || len(items) != 50 || items[0].File != "f0100.mp3" {
		t.Fatalf("total=%d len=%d first=%q", total, len(items), items[0].File)
	}
	// offset além do fim
	items, total = idx.Search("", 400, 50)
	if total != 250 || len(items) != 0 {
		t.Fatalf("past end: total=%d len=%d", total, len(items))
	}
	// limit clampado a 200
	items, _ = idx.Search("", 0, 1000)
	if len(items) != 200 {
		t.Fatalf("limit clamp: len=%d", len(items))
	}
	// limit default 50
	items, _ = idx.Search("", 0, 0)
	if len(items) != 50 {
		t.Fatalf("limit default: len=%d", len(items))
	}
}

func TestStats(t *testing.T) {
	idx := NewIndex()
	at := time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)
	idx.Set(sampleSongs(), at)
	count, refreshedAt := idx.Stats()
	if count != 4 || !refreshedAt.Equal(at) {
		t.Fatalf("count=%d at=%v", count, refreshedAt)
	}
}

func TestConcurrentSetAndSearch(t *testing.T) {
	idx := NewIndex()
	idx.Set(sampleSongs(), time.Now())
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				idx.Set(sampleSongs(), time.Now())
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				idx.Search("rock", 0, 10)
			}
		}()
	}
	wg.Wait()
}
