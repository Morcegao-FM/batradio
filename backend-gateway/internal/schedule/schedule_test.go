package schedule

import (
	"testing"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
)

var now = time.Date(2026, 7, 8, 14, 0, 0, 0, time.UTC)

// fila de 4 faixas de 60s, posições 0..3
func queue4() []model.Song {
	return []model.Song{
		{File: "s0.mp3", Pos: 0, Time: 60},
		{File: "s1.mp3", Pos: 1, Time: 60},
		{File: "s2.mp3", Pos: 2, Time: 60},
		{File: "s3.mp3", Pos: 3, Time: 60},
	}
}

func TestComputeTimesNoRepeat(t *testing.T) {
	st := model.Status{Song: 2, Elapsed: 30, Repeat: false}
	items := ComputeTimes(queue4(), st, now)
	if len(items) != 4 {
		t.Fatalf("len=%d", len(items))
	}
	// Âncora: faixa atual começou em now-30s
	anchor := now.Add(-30 * time.Second)
	want := []time.Time{
		anchor.Add(-120 * time.Second), // pos0 = pos1 - 60s
		anchor.Add(-60 * time.Second),  // pos1 = pos2 - 60s
		anchor,                         // pos2 (atual)
		anchor.Add(60 * time.Second),   // pos3
	}
	for i, w := range want {
		if !items[i].NextPresentation.Equal(w) {
			t.Errorf("pos%d: got %v want %v", i, items[i].NextPresentation, w)
		}
	}
}

func TestComputeTimesRepeatWrapsToEnd(t *testing.T) {
	st := model.Status{Song: 2, Elapsed: 30, Repeat: true}
	items := ComputeTimes(queue4(), st, now)
	anchor := now.Add(-30 * time.Second)
	end := anchor.Add(60 * time.Second) // horário da pos3 (última calculada pra frente)
	// Com repeat, pos0 e pos1 tocam DEPOIS do fim: pos0 = end+60, pos1 = pos0+60
	if !items[3].NextPresentation.Equal(end) {
		t.Fatalf("pos3: got %v want %v", items[3].NextPresentation, end)
	}
	if !items[0].NextPresentation.Equal(end.Add(60 * time.Second)) {
		t.Errorf("pos0: got %v want %v", items[0].NextPresentation, end.Add(60*time.Second))
	}
	if !items[1].NextPresentation.Equal(end.Add(120 * time.Second)) {
		t.Errorf("pos1: got %v want %v", items[1].NextPresentation, end.Add(120*time.Second))
	}
}

func TestComputeTimesCurrentSongMissing(t *testing.T) {
	st := model.Status{Song: 99, Elapsed: 30}
	items := ComputeTimes(queue4(), st, now)
	for i, it := range items {
		if !it.NextPresentation.IsZero() {
			t.Errorf("pos%d: expected zero time, got %v", i, it.NextPresentation)
		}
	}
}

func TestComputeTimesEmptyQueue(t *testing.T) {
	items := ComputeTimes(nil, model.Status{}, now)
	if len(items) != 0 {
		t.Fatalf("len=%d", len(items))
	}
}

// fila de 24h: 288 faixas de 5min
func queue24h() []QueueItem {
	items := make([]QueueItem, 288)
	for i := range items {
		items[i] = QueueItem{
			Song:             model.Song{File: "x.mp3", Pos: i, Time: 300},
			NextPresentation: now.Add(time.Duration(i) * 5 * time.Minute),
		}
	}
	return items
}

func TestPlanInsertionsEvery4Hours(t *testing.T) {
	positions := PlanInsertions(queue24h(), 6, 4*time.Hour, nil, nil, now)
	if len(positions) != 6 {
		t.Fatalf("len=%d %v", len(positions), positions)
	}
	// Decrescente
	for i := 1; i < len(positions); i++ {
		if positions[i-1] <= positions[i] {
			t.Fatalf("not descending: %v", positions)
		}
	}
	// Alvos: now, now+4h, ... → posições 0, 48, 96, 144, 192, 240 (5min por faixa)
	want := []int{240, 192, 144, 96, 48, 0}
	for i, w := range want {
		if positions[i] != w {
			t.Errorf("positions[%d]=%d want %d (all=%v)", i, positions[i], w, positions)
		}
	}
}

func TestPlanInsertionsWithPeriod(t *testing.T) {
	from := now.Add(1 * time.Hour)
	to := now.Add(8 * time.Hour) // alvos em +1h e +5h; +9h passa do fim (limite é inclusivo)
	positions := PlanInsertions(queue24h(), 6, 4*time.Hour, &from, &to, now)
	if len(positions) != 2 {
		t.Fatalf("len=%d %v", len(positions), positions)
	}
	if positions[0] != 60 || positions[1] != 12 {
		t.Errorf("got %v want [60 12]", positions)
	}
}

func TestPlanInsertionsTargetBeyondQueueEnd(t *testing.T) {
	short := queue24h()[:12] // 1h de fila
	positions := PlanInsertions(short, 3, 4*time.Hour, nil, nil, now)
	// Alvos: now (pos 0), now+4h (além do fim → 12), now+8h (→ 12)
	if len(positions) != 3 {
		t.Fatalf("len=%d %v", len(positions), positions)
	}
	if positions[0] != 12 || positions[1] != 12 || positions[2] != 0 {
		t.Errorf("got %v want [12 12 0]", positions)
	}
}

// Regressão: com repeat ligado, as faixas ANTES da atual têm horário de amanhã
// (wrap). O plano deve inserir perto do alvo cronológico, não na posição 0.
func TestPlanInsertionsWithRepeatWrap(t *testing.T) {
	st := model.Status{Song: 2, Elapsed: 30, Repeat: true}
	songs := []model.Song{
		{Pos: 0, Time: 3600}, {Pos: 1, Time: 3600},
		{Pos: 2, Time: 3600}, {Pos: 3, Time: 3600},
	}
	queue := ComputeTimes(songs, st, now)
	// Cronologia: pos2 (13:59:30), pos3 (14:59:30), pos0 (wrap 15:59:30), pos1 (16:59:30).
	positions := PlanInsertions(queue, 3, 1*time.Hour, nil, nil, now)
	// Alvos 14:00→pos3, 15:00→pos0 (wrap), 16:00→pos1; decrescente: [3 1 0].
	want := []int{3, 1, 0}
	if len(positions) != 3 {
		t.Fatalf("positions: %v", positions)
	}
	for i, w := range want {
		if positions[i] != w {
			t.Errorf("positions[%d]=%d want %d (all=%v)", i, positions[i], w, positions)
		}
	}
}

func TestPlanInsertionsZeroTimes(t *testing.T) {
	if got := PlanInsertions(queue24h(), 0, time.Hour, nil, nil, now); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}
