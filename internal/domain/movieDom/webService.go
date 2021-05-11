package movieDom

type movieHandler struct {
	service Service
}

func NewMovieHandler(svc Service) *movieHandler {
	return &movieHandler{
		service:    svc,
	}
}

/*func (api movieHandler) WebService() *restful.WebService {
	wsMovie := &restful.WebService{}
	wsMovie.Path("/v1/movies")

	wsMovie.Route(wsMovie.GET("/").
		Doc("Get movie or series by search").
		Param(wsMovie.QueryParameter("title", "Get movie or series by title").DataType("string")).
		Param(wsMovie.QueryParameter("type", "Type of media to search (movie, series, episode)").DataType("string")).
		Writes(nil).
		Returns(200, "OK", SearchResult{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json",typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.SearchMovie))

	wsMovie.Route(wsMovie.GET("/{movieId}").
		Doc("Get movie by ID").
		Param(wsMovie.PathParameter("id", "Get movie by ID (based on IMDb ids)").DataType("string")).
		Writes(nil).
		Returns(200, "OK", Result{}).
		Returns(400, "Bad request, fields not validated", typedErrors.CustomError{}).
		Returns(401, "Unauthorized, user cannot access this route", typedErrors.CustomError{}).
		Returns(422, "Not processable, impossible to serialize json",typedErrors.CustomError{}).
		Filter(logRequest()).
		Filter(authenticateUser()).
		To(api.GetMovie))

	return wsMovie
}

func (api movieHandler) SearchMovie(req *restful.Request, res *restful.Response) {
	title := req.QueryParameter("title")
	mediaType := req.QueryParameter("type")
	search := Search{
		Title:     title,
		MediaType: mediaType,
	}
	result, err := api.service.SearchMovie(search)
	if err != nil {
		handleHTTPError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, result)
}

func (api movieHandler) GetMovie(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	movie, err := api.service.GetMovieByID(movieId)
	if err != nil {
		handleHTTPError(res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, movie)
}*/
