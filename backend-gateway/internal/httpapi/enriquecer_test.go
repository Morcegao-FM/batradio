package httpapi

import (
	"context"
	"testing"

	"github.com/Morcegao-FM/batradio/backend-gateway/internal/catalogo"
	"github.com/Morcegao-FM/batradio/backend-gateway/internal/model"
)

type catalogoFake struct {
	dados map[string]catalogo.Info
	vezes int
}

func (f *catalogoFake) Lookup(_ context.Context, arquivos []string) map[string]catalogo.Info {
	f.vezes++
	out := make(map[string]catalogo.Info)
	for _, a := range arquivos {
		if info, ok := f.dados[a]; ok {
			out[a] = info
		}
	}
	return out
}

func TestEnriquecerPreencheCapaEMantemOResto(t *testing.T) {
	fake := &catalogoFake{dados: map[string]catalogo.Info{
		"a.mp3": {Tipo: "musica", ImagemURL: "https://capa/1.jpg", Ano: 1980, Artista: "AC/DC"},
	}}
	s := &Server{catalogo: fake}

	musicas := []model.Song{
		{File: "a.mp3", Artist: "AC/DC", Title: "Back in Black", Time: 255},
		{File: "b.mp3", Artist: "Desconhecido", Title: "Sem capa", Time: 120},
	}
	s.enriquecer(context.Background(), musicas)

	if musicas[0].ImageURL != "https://capa/1.jpg" || musicas[0].Year != 1980 {
		t.Errorf("faixa conhecida não foi enriquecida: %+v", musicas[0])
	}
	if musicas[0].Title != "Back in Black" || musicas[0].Time != 255 {
		t.Errorf("enriquecer não pode mexer no que veio do MPD: %+v", musicas[0])
	}
	if musicas[1].ImageURL != "" {
		t.Errorf("faixa desconhecida tem de ficar sem capa: %+v", musicas[1])
	}
}

func TestEnriquecerSemCatalogoNaoQuebra(t *testing.T) {
	s := &Server{catalogo: nil}
	musicas := []model.Song{{File: "a.mp3"}}
	s.enriquecer(context.Background(), musicas) // não pode entrar em pânico
	if musicas[0].ImageURL != "" {
		t.Error("sem catálogo não há capa")
	}
}

func TestEnriquecerNomeExibicaoVenceArtistaTitulo(t *testing.T) {
	// É o que resolve vinheta com nome de arquivo feio: a Batcaverna define o
	// nome de exibição e o BatRadio o respeita.
	fake := &catalogoFake{dados: map[string]catalogo.Info{
		"vinheta-01.mp3": {Tipo: "vinheta", NomeExibicao: "Vinheta de abertura"},
	}}
	s := &Server{catalogo: fake}

	musicas := []model.Song{{File: "vinheta-01.mp3", Artist: "", Title: "vinheta-01"}}
	s.enriquecer(context.Background(), musicas)

	if musicas[0].DisplayName != "Vinheta de abertura" || musicas[0].Kind != "vinheta" {
		t.Errorf("nome de exibição não chegou: %+v", musicas[0])
	}
}

func TestEnriquecerStatusSoTocaAtualEProxima(t *testing.T) {
	fake := &catalogoFake{dados: map[string]catalogo.Info{
		"atual.mp3":   {ImagemURL: "https://capa/atual.jpg"},
		"proxima.mp3": {ImagemURL: "https://capa/prox.jpg"},
	}}
	s := &Server{catalogo: fake}

	st := model.Status{
		Current: &model.Song{File: "atual.mp3"},
		Next:    &model.Song{File: "proxima.mp3"},
	}
	s.enriquecerStatus(context.Background(), &st)

	if st.Current.ImageURL != "https://capa/atual.jpg" {
		t.Errorf("faixa atual sem capa: %+v", st.Current)
	}
	if st.Next.ImageURL != "https://capa/prox.jpg" {
		t.Errorf("próxima faixa sem capa: %+v", st.Next)
	}
	if fake.vezes != 1 {
		t.Errorf("status deveria consultar o catálogo uma vez só, veio %d", fake.vezes)
	}
}

func TestEnriquecerStatusSemFaixaNaoConsulta(t *testing.T) {
	fake := &catalogoFake{}
	s := &Server{catalogo: fake}
	st := model.Status{}
	s.enriquecerStatus(context.Background(), &st)
	if fake.vezes != 0 {
		t.Error("MPD parado não deve gerar consulta ao catálogo")
	}
}
