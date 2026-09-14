// Package model define os tipos de domínio do gateway e o parse tolerante do
// JSON vindo do backend Node (que repassa valores do MPD como strings).
package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// FlexInt aceita número JSON ou string numérica (inclusive "352.133", truncada).
type FlexInt int

func (f *FlexInt) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return fmt.Errorf("FlexInt: %w", err)
	}
	*f = FlexInt(int(v))
	return nil
}

// FlexFloat aceita número JSON ou string numérica.
type FlexFloat float64

func (f *FlexFloat) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		*f = 0
		return nil
	}
	v, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return fmt.Errorf("FlexFloat: %w", err)
	}
	*f = FlexFloat(v)
	return nil
}

type Song struct {
	File   string `json:"file"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Album  string `json:"album,omitempty"`
	Genre  string `json:"genre,omitempty"`
	Time   int    `json:"time"`
	Pos    int    `json:"pos"`
	ID     int    `json:"id"`

	// Preenchidos pelo gateway a partir do catálogo do website; o Node/MPD
	// nunca manda estes campos. omitempty para não inchar a lista do acervo.
	ImageURL    string `json:"imageUrl,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Year        int    `json:"year,omitempty"`
	Kind        string `json:"kind,omitempty"`
}

type Status struct {
	State          string  `json:"state"`
	Song           int     `json:"song"`
	Elapsed        float64 `json:"elapsed"`
	Repeat         bool    `json:"repeat"`
	Random         bool    `json:"random"`
	Crossfade      bool    `json:"crossfade"`
	PlaylistLength int     `json:"playlistLength"`
	Current        *Song   `json:"current,omitempty"`
	Next           *Song   `json:"next,omitempty"`
}

type PlaylistInfo struct {
	Name         string `json:"name"`
	LastModified string `json:"lastModified,omitempty"`
}

// mpdSong espelha as chaves cruas repassadas pelo Node.
type mpdSong struct {
	File   string  `json:"file"`
	Artist string  `json:"Artist"`
	Title  string  `json:"Title"`
	Album  string  `json:"Album"`
	Genre  string  `json:"Genre"`
	Time   FlexInt `json:"Time"`
	Pos    FlexInt `json:"Pos"`
	ID     FlexInt `json:"Id"`
}

func (m mpdSong) toSong() Song {
	title := m.Title
	if title == "" {
		// Sem tag: usa o nome do arquivo, como faz a UI legada ao exibir file.
		title = m.File
		if i := strings.LastIndex(title, "/"); i >= 0 {
			title = title[i+1:]
		}
	}
	return Song{
		File:   m.File,
		Artist: m.Artist,
		Title:  title,
		Album:  m.Album,
		Genre:  m.Genre,
		Time:   int(m.Time),
		Pos:    int(m.Pos),
		ID:     int(m.ID),
	}
}

func ParseSong(raw []byte) (Song, error) {
	var m mpdSong
	if err := json.Unmarshal(raw, &m); err != nil {
		return Song{}, err
	}
	return m.toSong(), nil
}

func ParseSongs(raw []byte) ([]Song, error) {
	var ms []mpdSong
	if err := json.Unmarshal(raw, &ms); err != nil {
		return nil, err
	}
	songs := make([]Song, 0, len(ms))
	for _, m := range ms {
		if m.File == "" {
			continue // entradas de diretório/playlist não têm file
		}
		songs = append(songs, m.toSong())
	}
	return songs, nil
}
