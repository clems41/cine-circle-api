package movieDom

import (
	"cine-circle/internal/domain"
	"time"
)

type Result struct {
	MovieID 		domain.IDType 		`json:"id"`
	ImdbID 			string 				`json:"imdbId"`
	Title 			string 				`json:"title"`
	Year 			string 				`json:"year"`
	Released 		time.Time 			`json:"released"`
	Runtime 		int 				`json:"runtime"`
	Genres 			[]string 			`json:"genres"`
	Directors 		[]string 			`json:"directors"`
	Actors	 		[]string 			`json:"actors"`
	Plot 			string 				`json:"plot"`
	Countries 		[]string 			`json:"countries"`
	Poster 			string 				`json:"poster"`
}

type OmdbRatingView struct {
	Source 			string 				`json:"Source"`
	Value  			string 				`json:"Value"`
}

type OmdbView struct {
	Title    		string 				`json:"Title"`
	Year     		string 				`json:"Year"`
	Rated    		string 				`json:"Rated"`
	Released 		string 				`json:"Released"`
	Runtime  		string 				`json:"Runtime"`
	Genre    		string 				`json:"Genre"`
	Director 		string 				`json:"Director"`
	Writer  		string 				`json:"Writer"`
	Actors   		string 				`json:"Actors"`
	Plot     		string 				`json:"Plot"`
	Language 		string 				`json:"Language"`
	Country  		string 				`json:"Country"`
	Awards   		string 				`json:"Awards"`
	Poster   		string 				`json:"Poster"`
	Ratings  		[]OmdbRatingView 	`json:"Ratings"`
	Metascore  		string 				`json:"Metascore"`
	Imdbrating 		string 				`json:"imdbRating"`
	Imdbvotes  		string 				`json:"imdbVotes"`
	Imdbid     		string 				`json:"imdbID"`
	Type       		string 				`json:"Type"`
	Dvd        		string 				`json:"DVD"`
	Boxoffice  		string 				`json:"BoxOffice"`
	Production 		string 				`json:"Production"`
	Website    		string 				`json:"Website"`
	Response   		string 				`json:"Response"`
	TotalSeasons 	string 				`json:"totalSeasons"`
}
