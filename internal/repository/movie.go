package repository

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/domain/movieDom"
	"cine-circle/internal/logger"
	"cine-circle/internal/typedErrors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var _ movieDom.Repository = (*movieRepository)(nil)

type Movie struct {
	Metadata
	ImdbID 			string 				`gorm:"uniqueIndex"`
	Title 			string				`gorm:"index"`
	Year 			string
	Released 		time.Time
	Runtime 		int
	Genres 			pq.StringArray 		`gorm:"type:text[]"`
	Directors 		pq.StringArray 		`gorm:"type:text[]"`
	Actors	 		pq.StringArray 		`gorm:"type:text[]"`
	Plot 			string
	Countries 		pq.StringArray 		`gorm:"type:text[]"`
	Poster 			string
	Type 			string
}

type movieRepository struct {
	DB *gorm.DB
}

func NewMovieRepository(DB *gorm.DB) *movieRepository {
	return &movieRepository{DB: DB}
}

func (r movieRepository) Migrate() {
	err := r.DB.AutoMigrate(&Movie{})
	if err != nil {
		logger.Sugar.Fatalf("Error occurs when migrating movieRepository : %s", err.Error())
	}
}

func (r movieRepository) GetMovie(movieId string) (result movieDom.Result, err error) {
	var movie Movie
	response := r.DB.
		Find(&movie, "imdb_id = ?", movieId)
	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}
	if response.RowsAffected == 0 {
		return result, typedErrors.ErrRepositoryResourceNotFound
	}

	result = r.movieToResult(movie)
	return
}

func (r movieRepository) SaveMovie(movieView movieDom.OmdbView) (result movieDom.Result, err error) {
	releasedTime, err := time.Parse(constant.ReleasedLayout, movieView.Released)
	if err != nil {
		return result, typedErrors.NewServiceGeneralError(err)
	}

	runTime, err := strconv.Atoi(strings.Replace(movieView.Runtime, constant.RunTimeUnit, "", -1))
	if err != nil {
		return result, typedErrors.NewServiceGeneralError(err)
	}

	genres := strings.Split(movieView.Genre, constant.StringArraySeparator)
	directors := strings.Split(movieView.Director, constant.StringArraySeparator)
	actors := strings.Split(movieView.Actors, constant.StringArraySeparator)
	countries := strings.Split(movieView.Country, constant.StringArraySeparator)

	movie := Movie{
		ImdbID:    movieView.Imdbid,
		Title:     movieView.Title,
		Year:      movieView.Year,
		Released:  releasedTime,
		Runtime:   runTime,
		Genres:    genres,
		Directors: directors,
		Actors:    actors,
		Plot:      movieView.Plot,
		Countries: countries,
		Poster:    movieView.Poster,
		Type:      movieView.Type,
	}

	err = r.DB.
		Create(&movie).
		Error

	if err != nil {
		return result, typedErrors.NewRepositoryQueryFailedError(err)
	}

	result = r.movieToResult(movie)
	return
}

func (r movieRepository) movieToResult(movie Movie) (result movieDom.Result) {
	result = movieDom.Result{
		ID:        movie.ImdbID,
		Title:     movie.Title,
		Year:      movie.Year,
		Released:  movie.Released,
		Runtime:   movie.Runtime,
		Genres:    movie.Genres,
		Directors: movie.Directors,
		Actors:    movie.Actors,
		Plot:      movie.Plot,
		Countries: movie.Countries,
		Poster:    movie.Poster,
		Type:      movie.Type,
	}
	return
}