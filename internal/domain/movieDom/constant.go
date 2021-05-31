package movieDom

import "cine-circle/internal/utils"

// Used for Omdb
const (
	OmdbApiUrl = "http://www.omdbapi.com/"

	defaultOmdbAPIKey = "9d8fa748"
	envOmdbAPIKey = "EXTERNAL_API_KEY"

	defaultPlotValue = "full" //(full or short)
	envPlotValue = "EXTERNAL_PLOT_VALUE"

	MovieMedia = "movie"
	SeriesMedia = "series"

	ReleasedLayout = "02 Jan 2006"

	StringArraySeparator = ","
	RunTimeUnit = " min"
)

var (
	ExternalOmdbApiKey = utils.GetDefaultOrFromEnv(defaultOmdbAPIKey, envOmdbAPIKey)
	PlotValue = utils.GetDefaultOrFromEnv(defaultPlotValue, envPlotValue)
)

// Used for The Movie Database
const (
	MovieDBMovieID = "{movieId}"

	MovieDBApiUrl = "https://api.themoviedb.org/3"
	MovieDBApiMovieUrl = MovieDBApiUrl + "/movie/" + MovieDBMovieID
	MovieDBApiSearchUrl = MovieDBApiUrl + "/search/movie/"
	MovieDBApiVideoUrl = MovieDBApiMovieUrl+ "/videos"


	defaultMovieDBApiToken = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIyNDAxNjYxOWFiMjNiMDYzNjMzYzgwZTY4MzFlN2NjYyIsInN1YiI6IjYwOGI3ZjZlOGM0MGY3MDA1N2U3ZDg4MCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.TqHh6OC7IZ0s7err6njtR054Pi87kG6UaaER5WL04k0"
	envMovieDBApiKey = "EXTERNAL_API_TOKEN"

	releaseDateLayout = "2006-01-02"

	languageQueryParameterKey = "language"
	frenchValue = "fr-FR"
)

var (
	ExternalMovieDBApiToken = utils.GetDefaultOrFromEnv(defaultMovieDBApiToken, envMovieDBApiKey)
)
