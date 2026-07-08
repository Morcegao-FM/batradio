// Package schedule calcula a grade de horários da fila (port fiel de
// BatRadioClient.TimeCalculation do cliente Windows) e o planejamento da
// inserção distribuída ("inserir música várias vezes").
package schedule

import (
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
)

type QueueItem struct {
	model.Song
	NextPresentation time.Time `json:"nextPresentation"`
}

// ComputeTimes calcula quando cada faixa da fila deve tocar. A âncora é a faixa
// atual (começou em now - elapsed). Para frente, acumula durações. Para trás:
// sem repeat, subtrai durações (faixas já tocadas); com repeat, as faixas
// anteriores tocarão de novo após o fim da fila (wrap), acumulando a partir da
// última. Se a faixa atual não está na fila, os horários ficam zerados — mesmo
// comportamento do cliente legado.
func ComputeTimes(songs []model.Song, st model.Status, now time.Time) []QueueItem {
	items := make([]QueueItem, len(songs))
	for i, s := range songs {
		items[i] = QueueItem{Song: s}
	}
	initial := -1
	for i, s := range songs {
		if s.Pos == st.Song {
			initial = i
			break
		}
	}
	if initial == -1 {
		return items
	}

	anchor := now.Add(-time.Duration(int(st.Elapsed)) * time.Second)
	items[initial].NextPresentation = anchor

	ref := anchor
	last := items[initial]
	for i := initial + 1; i < len(items); i++ {
		items[i].NextPresentation = ref.Add(time.Duration(last.Time) * time.Second)
		ref = items[i].NextPresentation
		last = items[i]
	}

	if !st.Repeat {
		ref = anchor
		for i := initial - 1; i >= 0; i-- {
			items[i].NextPresentation = ref.Add(-time.Duration(items[i].Time) * time.Second)
			ref = items[i].NextPresentation
		}
	} else {
		// Fiel ao C#: usa a duração da própria faixa ao acumular o wrap.
		for i := 0; i < initial; i++ {
			items[i].NextPresentation = ref.Add(time.Duration(items[i].Time) * time.Second)
			ref = items[i].NextPresentation
		}
	}
	return items
}

// PlanInsertions gera as posições onde inserir uma faixa para que toque `times`
// vezes espaçadas por `interval`, a partir de `from` (ou now) e, se `to` for
// informado, sem passar do fim do período. Retorna posições em ordem
// DECRESCENTE: inserindo de trás para frente, as posições anteriores continuam
// válidas.
func PlanInsertions(queue []QueueItem, times int, interval time.Duration, from, to *time.Time, now time.Time) []int {
	if times <= 0 || interval <= 0 {
		return nil
	}
	start := now
	if from != nil {
		start = *from
	}
	var positions []int
	for k := 0; k < times; k++ {
		target := start.Add(time.Duration(k) * interval)
		if to != nil && target.After(*to) {
			break
		}
		pos := len(queue)
		for i, item := range queue {
			if !item.NextPresentation.Before(target) {
				pos = i
				break
			}
		}
		positions = append(positions, pos)
	}
	// Ordena decrescente (targets crescem, então basta inverter).
	for i, j := 0, len(positions)-1; i < j; i, j = i+1, j-1 {
		positions[i], positions[j] = positions[j], positions[i]
	}
	return positions
}
