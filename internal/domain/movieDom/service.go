package movieDom

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/logger"
	"cine-circle/internal/typedErrors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var _ Service = (*service)(nil)

type Service interface {
	GetMovieByID(movieId string) (result Result, err error)
	SearchMovie(search Search) (searchResult SearchResult, err error)
	getMovieFromExternal(movieId string) (result Result, err error)
	searchMovieFromExternal(search Search) (searchResult SearchResult, err error)
}

type service struct {
	r Repository
}

type Repository interface {
	GetMovie(movieId string) (result Result, err error)
	SaveMovie(movie OmdbView) (result Result, err error)
}

func NewService(r Repository) Service {
	return &service{
		r:                              r,
	}
}

// GetMovieByID Try to get movie from database, if not exists, get it from External API
func (svc *service) GetMovieByID(movieId string) (result Result, err error) {
	result, err = svc.r.GetMovie(movieId)
	if err == typedErrors.ErrRepositoryResourceNotFound {
		return svc.getMovieFromExternal(movieId)
	}
	return
}

func (svc *service) SearchMovie(search Search) (searchResult SearchResult, err error) {
	return svc.searchMovieFromExternal(search)
}

// getMovieFromExternal get movie from external API (OMDb in this case)
func (svc *service) getMovieFromExternal(movieId string) (result Result, err error) {
	// queryParam i means movieID and plot is use for specify version of plot to get (short or full)
	params := []QueryParam{
		{
			Key:   "i",
			Value: movieId,
		},
		{
			Key:   "plot",
			Value: constant.PlotValue,
		},
	}

	// Send request to external API
	resp, err := svc.sendRequestToExternal(params)

	// Unmarshall response to get it as ExternalMovie
	var externalMovie OmdbView
	err = json.Unmarshal(resp, &externalMovie)
	if err != nil {
		return result, typedErrors.NewExternalReadBodyError(err)
	}

	// Check if Response from External API is correct
	if externalMovie.Response == "False" {
		return result, typedErrors.ErrRepositoryResourceNotFound
	}

	// Save movie into database if not already exists
	return svc.r.SaveMovie(externalMovie)
}

// searchMovieFromExternal search for movies or series from external API (OMDb in this case)
func (svc *service) searchMovieFromExternal(search Search) (searchResult SearchResult, err error) {
	// queryParam i means movieID, plot is use for specify version of plot to get (short or full) and type is used for searching among series or movies
	params := []QueryParam{
		{
			Key:   "s",
			Value: search.Title,
		},
		{
			Key:   "plot",
			Value: constant.PlotValue,
		},
	}
	if search.MediaType == constant.MovieMedia || search.MediaType == constant.SeriesMedia {
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
}

// sendRequestToExternal Send request to External API in order to get specific movie/serie or list of media depending on search
func (svc *service) sendRequestToExternal(params []QueryParam) (response []byte, err error) {
	// useful in debug mode to know how long it took for getting movie from external API
	timeStart := time.Now()
	// Prepare request to send with queryParams
	client := &http.Client{}
	req, err := http.NewRequest("GET", constant.ExternalApiUrl, nil)
	if err != nil {
		return response, typedErrors.NewExternalSendRequestError(err)
	}
	q := req.URL.Query()
	for _, param := range params  {
		q.Add(param.Key, param.Value)
	}

	// Adding API key to the request (mandatory)
	q.Add("apikey", constant.ExternalApiKey)
	req.URL.RawQuery = q.Encode()

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return response, typedErrors.NewExternalSendRequestError(err)
	}

	// Returning response for getting movie(s)
	response, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return response, typedErrors.NewExternalReadBodyError(err)
	}
	err = res.Body.Close()
	if err != nil {
		return response, typedErrors.NewExternalReadBodyError(err)
	}

	// Print how long request took time
	logger.Sugar.Debugf("Sending request to %s with queryParameters %s took %+v", constant.ExternalApiUrl, req.URL.RawQuery, time.Since(timeStart))
	return
}
