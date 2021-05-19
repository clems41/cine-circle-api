package movieDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/repositoryModel"
	logger "cine-circle/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

var _ Service = (*service)(nil)

type Service interface {
	Get(movieId uint) (view View, err error)
	//SearchMovie(search Search) (searchResult SearchResult, err error)
}

type service struct {
	r repository
}

func NewService(r repository) Service {
	return &service{
		r: r,
	}
}

// Get Try to get movie from database, if not exists, get it from External API
func (svc *service) Get(movieId uint) (view View, err error) {
	movie, err := svc.r.Get(movieId)
	if err != nil {
		// If movie has not been found in database, we get it from external api
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var movieDBView MovieDBView
			movieDBView, err = svc.getMovieFromMovieDB(movieId)
			if err != nil {
				return
			}
			if movieDBView.Id != int(movieId) {
				return view, errMovieNotFound
			}
			view = svc.fromMovieDBToView(movieDBView)
			// Save movie into database for next request
			movieToSave := svc.fromViewToRepo(view)
			err = svc.r.Save(&movieToSave)
			if err != nil {
				return
			}
		} else {
			return
		}
		// Case with movie from database
	} else {
		view = svc.fromRepoToView(movie)
	}
	return
}

// getMovieFromExternal get movie from external API (OMDb in this case)
func (svc *service) getMovieFromMovieDB(movieId uint) (view MovieDBView, err error) {
	params := []QueryParameter{
		{
			Key:   languageQueryParameterKey,
			Value: frenchValue,
		},
	}

	// Send request to external API
	url := MovieDBApiMovieUrl + fmt.Sprintf("%d", movieId)
	resp, err := svc.sendRequestToExternal(url, http.MethodGet, params)

	// Unmarshall response to get it as ExternalMovie
	err = json.Unmarshal(resp, &view)
	if err != nil {
		return view, errors.WithStack(err)
	}
	return
}

// sendRequestToExternal Send request to External API in order to get specific movie/serie or list of media depending on search
func (svc *service) sendRequestToExternal(url, method string, queryParameters []QueryParameter) (response []byte, err error) {
	// useful in debug mode to know how long it took for getting movie from external API
	timeStart := time.Now()
	// Prepare request to send with queryParams
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return response, errors.WithStack(err)
	}
	q := req.URL.Query()
	for _, param := range queryParameters {
		q.Add(param.Key, param.Value)
	}
	req.URL.RawQuery = q.Encode()

	// Adding API token to the request (mandatory)
	req.Header.Set(constant.TokenHeader, constant.TokenKind + constant.BearerTokenDelimiterForHeader + ExternalMovieDBApiToken)

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return response, errors.WithStack(err)
	}

	// Returning response for getting movie(s)
	response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return response, errors.WithStack(err)
	}
	err = res.Body.Close()
	if err != nil {
		return response, errors.WithStack(err)
	}

	// Print how long request took time
	logger.Sugar.Debugf("Sending request to %s with queryParameters %s took %+v", url, req.URL.RawQuery, time.Since(timeStart))
	return
}

func (svc *service) fromMovieDBToView(movieDBView MovieDBView) (view View) {
	view = View{
		ID:               uint(movieDBView.Id),
		Title:            movieDBView.Title,
		ImdbId:           movieDBView.ImdbId,
		BackdropPath:     movieDBView.BackdropPath,
		PosterPath:       movieDBView.PosterPath,
		OriginalLanguage: movieDBView.OriginalLanguage,
		OriginalTitle:    movieDBView.OriginalTitle,
		Overview:         movieDBView.Overview,
		Runtime:          movieDBView.Runtime,
	}
	for _, movieDBGenre := range movieDBView.Genres {
		view.Genres = append(view.Genres, movieDBGenre.Name)
	}
	releaseDate, err := time.Parse(releaseDateLayout, movieDBView.ReleaseDate)
	if err != nil {
		logger.Sugar.Errorf("Cannot get release date from %s", movieDBView.ReleaseDate)
	} else {
		view.ReleaseDate = releaseDate
	}
	return
}

func (svc *service) fromRepoToView(movie repositoryModel.Movie) (view View) {
	view = View{
		ID:               movie.GetID(),
		Title:            movie.Title,
		ImdbId:           movie.ImdbId,
		BackdropPath:     movie.BackdropPath,
		PosterPath:       movie.PosterPath,
		OriginalLanguage: movie.OriginalLanguage,
		OriginalTitle:    movie.OriginalTitle,
		Overview:         movie.Overview,
		Runtime:          movie.Runtime,
		Genres:           movie.Genres,
		ReleaseDate:      movie.ReleaseDate,
	}
	return
}

func (svc *service) fromViewToRepo(view View) (movie repositoryModel.Movie) {
	movie = repositoryModel.Movie{
		Title:            view.Title,
		ImdbId:           view.ImdbId,
		BackdropPath:     view.BackdropPath,
		PosterPath:       view.PosterPath,
		OriginalLanguage: view.OriginalLanguage,
		OriginalTitle:    view.OriginalTitle,
		Overview:         view.Overview,
		Runtime:          view.Runtime,
		Genres:           view.Genres,
		ReleaseDate:      view.ReleaseDate,
	}
	movie.SetID(view.ID)
	return
}

// searchMovieFromExternal search for movies or series from external API (OMDb in this case)
/*func (svc *service) searchMovieFromExternal(search Search) (searchResult SearchResult, err error) {
	// queryParam i means movieID, plot is use for specify version of plot to get (short or full) and type is used for searching among series or movies
	params := []QueryParam{
		{
			Key:   "s",
			Value: search.Title,
		},
		{
			Key:   "plot",
			Value: PlotValue,
		},
	}
	if search.MediaType == MovieMedia || search.MediaType == SeriesMedia {
		params = append(params, QueryParam{
			Key:   "type",
			Value: search.MediaType,
		})
	}

	// Sending request to external API
	resp, err := svc.sendRequestToExternal(params)

	// Unmarshall response to get it as ExternalMovie
	var externalSearch OmdbSearchView
	err = json.Unmarshal(resp, &externalSearch)
	if err != nil {
		return searchResult, typedErrors.NewExternalReadBodyError(err)
	}

	// Check if Response from External API is correct
	if externalSearch.Response == "False" {
		return searchResult, typedErrors.ErrRepositoryResourceNotFound
	}

	// Convert search result
	for _, externalShortMovie := range externalSearch.Search {
		resultShort := ResultShort{
			ID:  	externalShortMovie.Imdbid,
			Title:   externalShortMovie.Title,
			Year:    externalShortMovie.Year,
			Poster:  externalShortMovie.Poster,
			Type:    externalShortMovie.Type,
		}
		searchResult.Search = append(searchResult.Search, resultShort)
	}
	searchResult.TotalResults, err = strconv.Atoi(externalSearch.TotalResults)
	return
}*/
