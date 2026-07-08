// Package library mantém o índice em memória do acervo de músicas, com a mesma
// semântica de busca do cliente Windows legado (todos os termos devem ser
// substring, case-insensitive, da concatenação file+title+album+artist+genre).
package library

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
)

const (
	defaultLimit = 50
	maxLimit     = 200
)

type Index struct {
	mu          sync.RWMutex
	songs       []model.Song
	keys        []string // chave normalizada, mesma ordem de songs
	refreshedAt time.Time
}

func NewIndex() *Index { return &Index{} }

// Set substitui o conteúdo do índice; ordena por File e pré-computa as chaves.
func (i *Index) Set(songs []model.Song, at time.Time) {
	sorted := make([]model.Song, len(songs))
	copy(sorted, songs)
	sort.Slice(sorted, func(a, b int) bool { return sorted[a].File < sorted[b].File })

	keys := make([]string, len(sorted))
	for n, s := range sorted {
		keys[n] = searchKey(s)
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.songs = sorted
	i.keys = keys
	i.refreshedAt = at
}

func searchKey(s model.Song) string {
	return strings.ToLower(s.File + " " + s.Title + " " + s.Album + " " + s.Artist + " " + s.Genre)
}

// Search retorna a página [offset, offset+limit) das faixas que casam com q e o
// total de faixas filtradas.
func (i *Index) Search(q string, offset, limit int) ([]model.Song, int) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}
	terms := splitTerms(q)

	i.mu.RLock()
	defer i.mu.RUnlock()

	var matched []int
	if len(terms) == 0 {
		matched = make([]int, len(i.songs))
		for n := range i.songs {
			matched[n] = n
		}
	} else {
		for n, key := range i.keys {
			if matchesAll(key, terms) {
				matched = append(matched, n)
			}
		}
	}

	total := len(matched)
	if offset >= total {
		return []model.Song{}, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	items := make([]model.Song, 0, end-offset)
	for _, n := range matched[offset:end] {
		items = append(items, i.songs[n])
	}
	return items, total
}

func (i *Index) Stats() (int, time.Time) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return len(i.songs), i.refreshedAt
}

// MatchesQuery aplica a mesma semântica de busca do índice a uma faixa avulsa
// (usada para filtrar a fila).
func MatchesQuery(s model.Song, q string) bool {
	terms := splitTerms(q)
	if len(terms) == 0 {
		return true
	}
	return matchesAll(searchKey(s), terms)
}

func splitTerms(q string) []string {
	var terms []string
	for _, t := range strings.Fields(q) {
		terms = append(terms, strings.ToLower(t))
	}
	return terms
}

func matchesAll(key string, terms []string) bool {
	for _, t := range terms {
		if !strings.Contains(key, t) {
			return false
		}
	}
	return true
}
