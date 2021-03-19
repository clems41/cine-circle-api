package api

import (
	"cine-circle/internal/logger"
	"cine-circle/internal/model"
	"cine-circle/internal/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	omdbUrl = "http://www.omdbapi.com/"

	defaultAPIKey = "9d8fa748"
	defaultPlot = "full" //(full or short)

	envAPIKey = "OMDB_API_KEY"
	envPlot = "OMDB_PLOT"
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", omdbUrl, nil)
	if err != nil {
		return model.NewCustomError(err, http.StatusInternalServerError), nil
	}
	q := req.URL.Query()
	for _, param := range params  {
		q.Add(param.Key, param.Value)
	}
	q.Add("apikey", omdbAPIKey)
	req.URL.RawQuery = q.Encode()

	logger.Sugar.Debugf("Sending request to %s with query %s", omdbUrl, req.URL.RawQuery)
	res, err := client.Do(req)
	if err != nil {
		return model.NewCustomError(err, http.StatusInternalServerError), nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	return model.NewCustomError(err, http.StatusInternalServerError), body
}

func finMovieByQueryParams(params []QueryParam) (model.CustomError, model.Movie) {
	var movie model.Movie
	err, resp := getDataFromOpenData(params)
	if err.IsNotNil() {
		return err, movie
	}
	err2 := json.Unmarshal(resp, &movie)
	return model.NewCustomError(err2, http.StatusInternalServerError), movie
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

func FindMovieBySearch(search string) (model.CustomError, model.Movie) {
	params := []QueryParam{
		{
			Key:   "t",
			Value: search,
		},
		{
			Key:   "plot",
			Value: omdbPlot,
		},
	}
	return finMovieByQueryParams(params)
}
