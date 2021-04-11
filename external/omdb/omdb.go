package omdb

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	omdbUrl = "http://www.omdbapi.com/"

	defaultAPIKey = "9d8fa748"
	defaultPlot = "full" //(full or short)

	envAPIKey = "OMDB_API_KEY"
	envPlot = "OMDB_PLOT"
)

const (
	movieMedia = "movie"
	seriesMedia = "series"
	episodeMedia = "episode"
)

type QueryParam struct {
	Key string
	Value string
}

var (
	omdbAPIKey string
	omdbPlot string
)

func init() {
	omdbAPIKey = utils.GetDefaultOrFromEnv(defaultAPIKey, envAPIKey)
	omdbPlot = utils.GetDefaultOrFromEnv(defaultPlot, envPlot)
}

func getDataFromOpenData(params []QueryParam) (typedErrors.CustomError, []byte) {
	timeStart := time.Now()
	client := &http.Client{}
	req, err := http.NewRequest("GET", omdbUrl, nil)
	if err != nil {
		return typedErrors.NewExternalSendRequestError(err), nil
	}
	q := req.URL.Query()
	for _, param := range params  {
		q.Add(param.Key, param.Value)
	}
	q.Add("apikey", omdbAPIKey)
	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)
	if err != nil {
		return typedErrors.NewExternalReadBodyError(err), nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	logger.Sugar.Debugf("Sending request to %s with query %s took %+v", omdbUrl, req.URL.RawQuery, time.Since(timeStart))
	return typedErrors.NewExternalReadBodyError(err), body
}

func finMovieByQueryParams(params []QueryParam) (typedErrors.CustomError, model.Movie) {
	var omdbMovie model.OmdbMovie
	err, resp := getDataFromOpenData(params)
	if err.IsNotNil() {
		return err, model.Movie{}
	}
	err2 := json.Unmarshal(resp, &omdbMovie)
	if err2 != nil {
		return typedErrors.NewExternalReadBodyError(err2), model.Movie{}
	}
	if omdbMovie.Response == "False" {
		return typedErrors.ErrRepositoryResourceNotFound, model.Movie{}
	}
	return typedErrors.NoErr, omdbMovie.Movie()
}

func FindMovieByID(id string) (typedErrors.CustomError, model.Movie) {
	params := []QueryParam{
		{
			Key:   "i",
			Value: id,
		},
		{
			Key:   "plot",
			Value: omdbPlot,
		},
	}
	return finMovieByQueryParams(params)
}

func FindMovieBySearch(titleToSearch, mediaType string) (typedErrors.CustomError, model.MovieSearch) {
	params := []QueryParam{
		{
			Key:   "s",
			Value: titleToSearch,
		},
		{
			Key:   "plot",
			Value: omdbPlot,
		},
	}
	if mediaType == movieMedia || mediaType == seriesMedia || mediaType == episodeMedia {
		params = append(params, QueryParam{
			Key:   "type",
			Value: mediaType,
		})
	}
	var omdbMovieSearch model.OmdbMovieSearch
	err, resp := getDataFromOpenData(params)
	if err.IsNotNil() {
		return err, model.MovieSearch{}
	}
	err2 := json.Unmarshal(resp, &omdbMovieSearch)
	if err2 != nil {
		return typedErrors.NewExternalReadBodyError(err2), model.MovieSearch{}
	}
	if omdbMovieSearch.Response == "False" {
		return typedErrors.ErrRepositoryResourceNotFound, model.MovieSearch{}
	}
	return typedErrors.NoErr, omdbMovieSearch.MovieSearch()
}

func MovieExists(id string) bool {
	err, movie := FindMovieByID(id)
	if err.IsNotNil() {
		return false
	}
	return movie.ID == id
}
