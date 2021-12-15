package theMovieDatabase

const (
	providerName = "theMovieDB"
)

const (
	apiUrl          = "https://api.themoviedb.org/3/"
	tokenKey        = `eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIyNDAxNjYxOWFiMjNiMDYzNjMzYzgwZTY4MzFlN2NjYyIsInN1YiI6IjYwOGI3ZjZlOGM0MGY3MDA1N2U3ZDg4MCIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.TqHh6OC7IZ0s7err6njtR054Pi87kG6UaaER5WL04k0`
	movieSuffix     = "movie/"
	tvSuffix        = "suffix/"
	searchSuffix    = "search/multi"
	searchQueryName = "query"
	imageBaseUrl    = "https://image.tmdb.org/t/p/w500"
)

const (
	authorizationHeaderName           = "Authorization"
	languageQueryParameterName        = "language"
	languageFrenchQueryParameterValue = "fr-FR"
)

const (
	releaseDateLayout = "2006-01-02"
)
