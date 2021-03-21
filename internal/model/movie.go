package model

import "time"

const (
	RatingSourceCineCircle = "Cine Circle"
)

type MovieRating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
	Comment string `json:"Comment"`
	PostedDate time.Time `json:"PostedDate"`
}

type Movie struct {
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
	Ratings  []MovieRating `json:"Ratings"`
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
}
