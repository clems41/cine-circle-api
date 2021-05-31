package movieDom

type MovieDBView struct {
	Adult               bool                       `json:"adult"`
	BackdropPath        string                     `json:"backdrop_path"`
	BelongsToCollection MovieDBCollection          `json:"belongs_to_collection"`
	Budget              int                        `json:"budget"`
	Genres              []MovieDBGenre             `json:"genres"`
	Homepage            string                     `json:"homepage"`
	Id                  int                        `json:"id"`
	ImdbId              string                     `json:"imdb_id"`
	OriginalLanguage    string                     `json:"original_language"`
	OriginalTitle       string                     `json:"original_title"`
	Overview            string                     `json:"overview"`
	Popularity          float64                    `json:"popularity"`
	PosterPath          string                     `json:"poster_path"`
	ProductionCompanies []MovieDBProductionCompany `json:"production_countries"`
	ReleaseDate         string                     `json:"release_date"`
	Revenue             int                        `json:"revenue"`
	Runtime             int                        `json:"runtime"`
	SpokenLanguages     []MovieDBSpokenLanguage    `json:"spoken_languages"`
	Status              string                     `json:"status"`
	Tagline             string                     `json:"tagline"`
	Title               string                     `json:"title"`
	Video               bool                       `json:"video"`
	VoteAverage         float64                    `json:"vote_average"`
	VoteCount           int                        `json:"vote_count"`
}

type MovieDBCollection struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

type MovieDBGenre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MovieDBProductionCompany struct {
	Id            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type MovieDBSpokenLanguage struct {
	EnglishName string `json:"english_name"`
	Iso6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
}

type MovieDBVideos struct {
	Id      int                    `json:"id"`
	Results []MovieDBVideosResult `json:"results"`
}

type MovieDBVideosResult struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
