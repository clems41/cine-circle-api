package theMovieDatabase

/* Get Movie */

type MovieView struct {
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	BelongsToCollection BelongsToCollection `json:"belongs_to_collection"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	Id                  int                 `json:"id"`
	ImdbId              string              `json:"imdb_id"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
}

type BelongsToCollection struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ProductionCompany struct {
	Id            int     `json:"id"`
	LogoPath      *string `json:"logo_path"`
	Name          string  `json:"name"`
	OriginCountry string  `json:"origin_country"`
}

type ProductionCountry struct {
	Iso31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	Iso6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
}

/* Search */

type SearchView struct {
	Page         int               `json:"page"`
	Results      []SearchMovieView `json:"results"`
	TotalPages   int               `json:"total_pages"`
	TotalResults int               `json:"total_results"`
}

type SearchMovieView struct {
	BackdropPath     string   `json:"backdrop_path"`
	FirstAirDate     string   `json:"first_air_date,omitempty"`
	GenreIds         []int    `json:"genre_ids"`
	Id               int      `json:"id"`
	MediaType        string   `json:"media_type"`
	Name             string   `json:"name,omitempty"`
	OriginCountry    []string `json:"origin_country,omitempty"`
	OriginalLanguage string   `json:"original_language"`
	OriginalName     string   `json:"original_name,omitempty"`
	Overview         string   `json:"overview"`
	Popularity       float64  `json:"popularity"`
	PosterPath       string   `json:"poster_path"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	Adult            bool     `json:"adult,omitempty"`
	OriginalTitle    string   `json:"original_title,omitempty"`
	ReleaseDate      string   `json:"release_date,omitempty"`
	Title            string   `json:"title,omitempty"`
	Video            bool     `json:"video,omitempty"`
}
