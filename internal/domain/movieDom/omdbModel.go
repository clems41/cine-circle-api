package movieDom

type OmdbRatingView struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

type OmdbView struct {
	Title        string           `json:"Title"`
	Year         string           `json:"Year"`
	Rated        string           `json:"Rated"`
	Released     string           `json:"Released"`
	Runtime      string           `json:"Runtime"`
	Genre        string           `json:"Genre"`
	Director     string           `json:"Director"`
	Writer       string           `json:"Writer"`
	Actors       string           `json:"Actors"`
	Plot         string           `json:"Plot"`
	Language     string           `json:"Language"`
	Country      string           `json:"Country"`
	Awards       string           `json:"Awards"`
	Poster       string           `json:"Poster"`
	Ratings      []OmdbRatingView `json:"Ratings"`
	Metascore    string           `json:"Metascore"`
	Imdbrating   string           `json:"imdbRating"`
	Imdbvotes    string           `json:"imdbVotes"`
	Imdbid       string           `json:"imdbID"`
	Type         string           `json:"Type"`
	Dvd          string           `json:"DVD"`
	Boxoffice    string           `json:"BoxOffice"`
	Production   string           `json:"Production"`
	Website      string           `json:"Website"`
	Response     string           `json:"Response"`
	TotalSeasons string           `json:"totalSeasons"`
}

type OmdbShortView struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	Poster string `json:"Poster"`
	Imdbid string `json:"imdbID"`
	Type   string `json:"Type"`
}

type OmdbSearchView struct {
	Search       []OmdbShortView `json:"Search"`
	TotalResults string          `json:"totalResults"`
	Response     string          `json:"Response"`
}
