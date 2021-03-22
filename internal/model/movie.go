package model

import (
	"github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

const (
	CineCircleSource = "Cine Circle"
	InternetMovieDatabaseSource = "Internet Movie Database"
	RottenTomatoesSource = "Rotten Tomatoes"
	MetacriticSource = "Metacritic"
	ImdbSource = "IMDb"
	TimeParserLayout = "2 Jan 2006"
	OmdbCurrency = "$"
)

type OmdbMovieRating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

type OmdbMovie struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []OmdbMovieRating `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	Imdbrating string `json:"imdbRating"`
	Imdbvotes  string `json:"imdbVotes"`
	Imdbid     string `json:"imdbID"`
	Type       string `json:"Type"`
	Dvd        string `json:"DVD"`
	Boxoffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
	TotalSeasons string `json:"totalSeasons"`
}

type Movie struct {
	ID       string `gorm:"primarykey"`
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Released time.Time `json:"Released"`
	Runtime  int `json:"Runtime"`
	Genres    pq.StringArray `json:"Genres"`
	Directors pq.StringArray `json:"Directors"`
	Writers   pq.StringArray `json:"Writers"`
	Actors   pq.StringArray `json:"Actors"`
	Plot     string `json:"Plot"`
	Languages pq.StringArray `json:"Languages"`
	Countries  pq.StringArray `json:"Countries"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	UserRatings  []Rating `json:"UserRatings"`
	PressRatings  []Rating `json:"PressRatings"`
	Metascore  string `json:"Metascore"`
	Imdbvotes  string `json:"imdbVotes"`
	Type       string `json:"Type"`
	Dvd        time.Time `json:"DVD"`
	BoxOffice  int `json:"BoxOffice"`
	BoxOfficeCurrency  string `json:"BoxOfficeCurrency"`
	Productions pq.StringArray `json:"Productions"`
	Website    string `json:"Website"`
	TotalSeasons int `json:"TotalSeasons"`
}

func (om OmdbMovie) Movie() Movie {
	movie := Movie{
		ID:                om.Imdbid,
		Title:             om.Title,
		Year:              om.Year,
		Plot:              om.Plot,
		Awards:            om.Awards,
		Poster:            om.Poster,
		Metascore:         om.Metascore,
		Imdbvotes:         om.Imdbvotes,
		Type:              om.Type,
		Website:           om.Website,
	}
	movie.Released, _ = time.Parse(TimeParserLayout, om.Released)
	movie.Runtime, _ = strconv.Atoi(strings.Replace(om.Runtime, " min", "", 1))
	movie.Genres = strings.Split(om.Genre, ", ")
	movie.Directors = strings.Split(om.Director, ", ")
	movie.Writers = strings.Split(om.Writer, ", ")
	movie.Actors = strings.Split(om.Actors, ", ")
	movie.Languages = strings.Split(om.Language, ", ")
	movie.Countries = strings.Split(om.Country, ", ")
	movie.Dvd, _ = time.Parse(TimeParserLayout, om.Dvd)
	movie.BoxOffice, _ = strconv.Atoi(strings.Replace(strings.Replace(om.Boxoffice, OmdbCurrency, "", 1), ",", "", -1))
	movie.BoxOfficeCurrency = OmdbCurrency
	movie.Productions = strings.Split(om.Production, ", ")
	movie.TotalSeasons, _ = strconv.Atoi(om.TotalSeasons)
	for _, pressRating := range om.Ratings {
		movie.PressRatings = append(movie.PressRatings, convertRatingDependingOnSource(pressRating))
	}
	imdbRating, _ := strconv.ParseFloat(om.Imdbrating, 64)
	movie.PressRatings = append(movie.PressRatings, Rating{
		Source:   ImdbSource,
		Value:    imdbRating,
	})
	return movie
}

func convertRatingDependingOnSource(pressRating OmdbMovieRating) Rating {
	rating := Rating{
		Source:   pressRating.Source,
	}
	var ratingValue float64
	switch pressRating.Source {
	case RottenTomatoesSource:
		valueStr := strings.Replace(pressRating.Value, "%", "", 1)
		ratingValue, _ = strconv.ParseFloat(valueStr, 64)
		ratingValue /= 10.0
	case MetacriticSource:
		valueStr := strings.Replace(pressRating.Value, "/100", "", 1)
		ratingValue, _ = strconv.ParseFloat(valueStr, 64)
		ratingValue /= 10.0
	default:
		valueStr := strings.Replace(pressRating.Value, "/10", "", 1)
		ratingValue, _ = strconv.ParseFloat(valueStr, 64)
	}
	rating.Value = ratingValue
	return rating
}