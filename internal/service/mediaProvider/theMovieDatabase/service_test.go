package theMovieDatabase

import (
	"cine-circle-api/internal/constant/mediaConst"
	"cine-circle-api/internal/service/mediaProvider"
	"fmt"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Get(t *testing.T) {
	type Result struct {
		Title    string
		Overview string
	}

	expectedResultByLanguage := map[mediaProvider.Language]Result{
		mediaConst.FrenchLanguage: {
			Title:    "The Dark Knight : Le Chevalier noir",
			Overview: "La suite de Batman Begins, The Dark Knight, le réalisateur Christopher Nolan et l'acteur Christian Bale, qui endosse à nouveau le rôle de Batman/Bruce Wayne dans sa guerre permanente contre le crime. Avec l'aide du Lieutenant de Police Jim Gordon et du Procureur Harvey Dent, Batman entreprend de démanteler définitivement les organisations criminelles de Gotham. L'association s'avère efficace, mais le trio se heurte bientôt à un nouveau génie du crime, plus connu sous le nom du Joker, qui va plonger Gotham dans l'anarchie et pousser Batman à la limite entre héros et assassin. Heath Ledger interprète Le Joker : le méchant suprême et Aaron Eckhart joue le rôle de Dent. Maggie Gyllenhaal complète le casting en tant que Rachel Dawes. De retour après Batman Begins, Gary Oldman est à nouveau Gordon, Michael Caine interprète Alfred, et Morgan Freeman est Lucius Fox.",
		},
		mediaConst.EnglishLanguage: {
			Title:    "The Dark Knight",
			Overview: "Batman raises the stakes in his war on crime. With the help of Lt. Jim Gordon and District Attorney Harvey Dent, Batman sets out to dismantle the remaining criminal organizations that plague the streets. The partnership proves to be effective, but they soon find themselves prey to a reign of chaos unleashed by a rising criminal mastermind known to the terrified citizens of Gotham as the Joker.",
		},
	}

	movieId := 155
	svc := New()

	// Try to get tv show with movie id, should return error
	form := mediaProvider.MediaForm{
		Id:   fmt.Sprintf("%d", movieId),
		Type: mediaConst.TvType,
	}
	view, err := svc.Get(form)
	require.Error(t, err)

	// Everything should be OK now
	for language, expectedResult := range expectedResultByLanguage {
		form = mediaProvider.MediaForm{
			Id:       fmt.Sprintf("%d", movieId),
			Language: language,
			Type:     mediaConst.MovieType,
		}
		view, err = svc.Get(form)
		require.NoError(t, err)
		require.Equal(t, expectedResult.Title, view.Title)
		require.Equal(t, expectedResult.Overview, view.Overview)
	}

	// Try to get movie without specifying language, should be OK and use mediaProvider.DefaultLanguage
	form = mediaProvider.MediaForm{
		Id:   fmt.Sprintf("%d", movieId),
		Type: mediaConst.MovieType,
	}
	view, err = svc.Get(form)
	require.NoError(t, err)
	require.Equal(t, expectedResultByLanguage[mediaConst.DefaultLanguage].Title, view.Title)
	require.Equal(t, expectedResultByLanguage[mediaConst.DefaultLanguage].Overview, view.Overview)

	// Try to get movie with incorrect language, should be OK and use mediaProvider.DefaultLanguage
	form = mediaProvider.MediaForm{
		Id:       fmt.Sprintf("%d", movieId),
		Type:     mediaConst.MovieType,
		Language: mediaProvider.Language(fake.Words()),
	}
	view, err = svc.Get(form)
	require.NoError(t, err)
	require.Equal(t, expectedResultByLanguage[mediaConst.DefaultLanguage].Title, view.Title)
	require.Equal(t, expectedResultByLanguage[mediaConst.DefaultLanguage].Overview, view.Overview)

	// Try to get movie with missing type, should be OK and use mediaProvider.DefaultMediaType
	form = mediaProvider.MediaForm{
		Id:       fmt.Sprintf("%d", movieId),
		Language: mediaConst.DefaultLanguage,
	}
	view, err = svc.Get(form)
	require.NoError(t, err)
	require.Equal(t, expectedResultByLanguage[mediaConst.DefaultLanguage].Title, view.Title)
	require.Equal(t, expectedResultByLanguage[mediaConst.DefaultLanguage].Overview, view.Overview)

	// Try to get movie with incorrect type, should be OK and use mediaProvider.DefaultMediaType
	form = mediaProvider.MediaForm{
		Id:       fmt.Sprintf("%d", movieId),
		Type:     mediaProvider.MediaType(fake.Words()),
		Language: mediaConst.EnglishLanguage,
	}
	view, err = svc.Get(form)
	require.NoError(t, err)
	require.Equal(t, expectedResultByLanguage[mediaConst.EnglishLanguage].Title, view.Title)
	require.Equal(t, expectedResultByLanguage[mediaConst.EnglishLanguage].Overview, view.Overview)
}
