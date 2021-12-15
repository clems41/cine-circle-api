package theMovieDatabase

import (
	"cine-circle-api/internal/constant/mediaConst"
	"cine-circle-api/internal/service/mediaProvider"
	"cine-circle-api/pkg/utils/httpUtils"
	"cine-circle-api/pkg/utils/sliceUtils"
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
	"time"
)

var _ mediaProvider.Service = (*service)(nil)

type service struct {
}

func New() (svc *service) {
	return &service{}
}

func (svc *service) Search(form mediaProvider.SearchForm) (view mediaProvider.SearchView, err error) {
	return
}

func (svc *service) Get(form mediaProvider.MediaForm) (view mediaProvider.MediaView, err error) {
	var urlSuffix string
	switch form.Type {
	case mediaConst.MovieType:
		urlSuffix = movieSuffix
	case mediaConst.TvType:
		urlSuffix = tvSuffix
	default:
		urlSuffix = movieSuffix
	}

	url := apiUrl + urlSuffix + form.Id

	var movie MovieView
	err = svc.askProvider(url, http.MethodGet, form.Language, nil, &movie)
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

	view = mediaProvider.MediaView{
		Id:            fmt.Sprintf("%d", movie.Id),
		Title:         movie.Title,
		BackdropUrl:   imageBaseUrl + movie.BackdropPath,
		Genres:        genres,
		Language:      movie.OriginalLanguage,
		OriginalTitle: movie.OriginalTitle,
		Overview:      movie.Overview,
		PosterUrl:     imageBaseUrl + movie.PosterPath,
		ReleaseDate:   releasedDate,
		Runtime:       time.Duration(movie.Runtime) * time.Minute,
	}

	return
}

/* PRIVATE METHODS */

func (svc *service) askProvider(url, method string, language mediaProvider.Language, body interface{}, response interface{}) (err error) {
	request := httpUtils.Request{
		Url:         url,
		Method:      method,
		ContentType: restful.MIME_JSON,
		Body:        body,
	}

	request.HeadersParameters = map[string]string{
		authorizationHeaderName: fmt.Sprintf("Bearer %s", tokenKey),
	}

	if !sliceUtils.SliceContainsStr(mediaConst.AllowedLanguages(), string(language)) {
		language = mediaConst.DefaultLanguage
	}
	request.QueryParameters = map[string]string{
		languageQueryParameterName: queryLanguageValue[string(language)],
	}

	httpStatus, err := httpUtils.SendRequest(request, &response, nil)
	if err != nil {
		return
	}
	if httpStatus >= http.StatusBadRequest {
		return fmt.Errorf("mediaProvider return not OK status %d", httpStatus)
	}
	return
}
