package theMovieDatabase

import (
	mediaProvider2 "cine-circle-api/external/mediaProvider"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/httpUtils"
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
	"strings"
	"time"
)

var _ mediaProvider2.Service = (*service)(nil)

type service struct {
}

func New() (svc *service) {
	return &service{}
}

func (svc *service) GetProviderName() (name string) {
	return providerName
}

func (svc *service) Search(form mediaProvider2.SearchForm) (view mediaProvider2.SearchView, err error) {
	url := apiUrl + searchSuffix

	queryParameters := map[string]string{
		searchQueryName: strings.ToLower(form.Keyword),
	}

	var search SearchView
	err = svc.askProvider(url, http.MethodGet, queryParameters, nil, &search)
	if err != nil {
		return
	}

	view.NumberOfItems = search.TotalResults
	view.CurrentPage = search.Page
	view.NumberOfPages = search.TotalPages
	for _, result := range search.Results {
		view.Result = append(view.Result, mediaProvider2.MovieShortView{
			Id:            fmt.Sprintf("%d", result.Id),
			Title:         result.Title,
			Language:      result.OriginalLanguage,
			OriginalTitle: result.OriginalTitle,
			PosterUrl:     svc.getImageUrl(result.PosterPath),
		})
	}

	return
}

func (svc *service) Get(form mediaProvider2.MovieForm) (view mediaProvider2.MovieView, err error) {
	url := apiUrl + movieSuffix + form.Id

	var movie MovieView
	err = svc.askProvider(url, http.MethodGet, nil, nil, &movie)
	if err != nil {
		return
	}

	releasedDate, err := time.Parse(releaseDateLayout, movie.ReleaseDate)
	if err != nil {
		return
	}

	var genres []string
	for _, genre := range movie.Genres {
		genres = append(genres, genre.Name)
	}

	view = mediaProvider2.MovieView{
		Id:            fmt.Sprintf("%d", movie.Id),
		Title:         movie.Title,
		BackdropUrl:   svc.getImageUrl(movie.BackdropPath),
		Genres:        genres,
		Language:      movie.OriginalLanguage,
		OriginalTitle: movie.OriginalTitle,
		Overview:      movie.Overview,
		PosterUrl:     svc.getImageUrl(movie.PosterPath),
		ReleaseDate:   releasedDate,
		Runtime:       movie.Runtime,
	}

	return
}

/* PRIVATE METHODS */

func (svc *service) askProvider(url, method string, queryParameters map[string]string, body, response interface{}) (err error) {
	startTime := time.Now()
	request := httpUtils.Request{
		Url:             url,
		Method:          method,
		ContentType:     restful.MIME_JSON,
		QueryParameters: queryParameters,
		Body:            body,
	}

	if request.QueryParameters == nil {
		request.QueryParameters = make(map[string]string)
	}
	request.QueryParameters[languageQueryParameterName] = languageFrenchQueryParameterValue

	request.HeadersParameters = map[string]string{
		authorizationHeaderName: fmt.Sprintf("Bearer %s", tokenKey),
	}

	httpStatus, err := httpUtils.SendRequest(request, &response, nil)
	logger.Infof("Request %s to %s returned %d in %v", request.Method, request.Url, httpStatus, time.Since(startTime))
	if err != nil {
		return
	}
	if httpStatus >= http.StatusBadRequest {
		return fmt.Errorf("mediaProvider return not OK status %d", httpStatus)
	}
	return
}

func (svc *service) getImageUrl(imagePath string) (url string) {
	return imageBaseUrl + imagePath
}
