package omdb

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
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

func getDataFromOpenData(params []QueryParam) (model.CustomError, []byte) {
	timeStart := time.Now()
	client := &http.Client{}
	req, err := http.NewRequest("GET", omdbUrl, nil)
	if err != nil {
		return model.NewCustomError(err, http.StatusInternalServerError, model.ErrExternalSendRequestCode), nil
	}
	q := req.URL.Query()
	for _, param := range params  {
		q.Add(param.Key, param.Value)
	}
	q.Add("apikey", omdbAPIKey)
	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)
	if err != nil {
		return model.NewCustomError(err, http.StatusInternalServerError, model.ErrExternalReceiveResponseCode), nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	logger.Sugar.Debugf("Sending request to %s with query %s took %+v", omdbUrl, req.URL.RawQuery, time.Since(timeStart))
	return model.NewCustomError(err, http.StatusInternalServerError, model.ErrExternalReadBodyCode), body
}

func finMovieByQueryParams(params []QueryParam) (model.CustomError, model.Movie) {
	var omdbMovie model.OmdbMovie
	err, resp := getDataFromOpenData(params)
	if err.IsNotNil() {
		return err, model.Movie{}
	}
	err2 := json.Unmarshal(resp, &omdbMovie)
	if err2 != nil {
		return model.NewCustomError(err2, http.StatusInternalServerError, model.ErrExternalReadBodyCode), model.Movie{}
	}
	if omdbMovie.Response == "False" {
		return model.ErrInternalDatabaseResourceNotFound, model.Movie{}
	}
	return model.NoErr, omdbMovie.Movie()
}

func FindMovieByID(id string) (model.CustomError, model.Movie) {
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

func FindMovieBySearch(titleToSearch, mediaType string) (model.CustomError, model.MovieSearch) {
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
		return model.NewCustomError(err2, http.StatusInternalServerError, model.ErrExternalReadBodyCode), model.MovieSearch{}
	}
	if omdbMovieSearch.Response == "False" {
		return model.ErrInternalDatabaseResourceNotFound, model.MovieSearch{}
	}
	return model.NoErr, omdbMovieSearch.MovieSearch()
}

func MovieExists(id string) bool {
	err, movie := FindMovieByID(id)
	if err.IsNotNil() {
		return false
	}
	return movie.ID == id
}
