package theMovieDatabase

import (
	"cine-circle-api/internal/service/mediaProvider"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Get(t *testing.T) {
	expectedOverview := "La suite de Batman Begins, The Dark Knight, le réalisateur Christopher Nolan et l'acteur Christian Bale, qui endosse à nouveau le rôle de Batman/Bruce Wayne dans sa guerre permanente contre le crime. Avec l'aide du Lieutenant de Police Jim Gordon et du Procureur Harvey Dent, Batman entreprend de démanteler définitivement les organisations criminelles de Gotham. L'association s'avère efficace, mais le trio se heurte bientôt à un nouveau génie du crime, plus connu sous le nom du Joker, qui va plonger Gotham dans l'anarchie et pousser Batman à la limite entre héros et assassin. Heath Ledger interprète Le Joker : le méchant suprême et Aaron Eckhart joue le rôle de Dent. Maggie Gyllenhaal complète le casting en tant que Rachel Dawes. De retour après Batman Begins, Gary Oldman est à nouveau Gordon, Michael Caine interprète Alfred, et Morgan Freeman est Lucius Fox."
	expectedTitle := "The Dark Knight : Le Chevalier noir"

	movieId := 155
	svc := New()

	form := mediaProvider.MovieForm{
		Id: fmt.Sprintf("%d", movieId),
	}
	view, err := svc.Get(form)
	require.NoError(t, err)
	require.Equal(t, expectedTitle, view.Title)
	require.Equal(t, expectedOverview, view.Overview)
}

func TestService_Search(t *testing.T) {
	svc := New()
	keyword := "The Dark Knight"

	form := mediaProvider.SearchForm{
		Page:    1,
		Keyword: keyword,
	}
	view, err := svc.Search(form)
	require.NoError(t, err)
	require.Equal(t, 1, view.CurrentPage)
	require.True(t, view.NumberOfPages > 0)
	require.True(t, view.NumberOfItems > 0)
	require.True(t, len(view.Result) > 0)
}
