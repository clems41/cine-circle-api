package watchlistDom

import (
	"cine-circle/internal/utils"
	webServicePkg "cine-circle/internal/webService"
	"github.com/emicklei/go-restful"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		service: svc,
	}
}

func (api handler) WebServices() (webServices []*restful.WebService) {
	wsWatchlist := &restful.WebService{}
	webServices = append(webServices, wsWatchlist)

	wsWatchlist.Path("/v1/watchlist")

	wsWatchlist.Route(wsWatchlist.POST("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to add into actualUser's watchlist").DataType("int")).
		Doc("Add movie into watchlist of actualUser").
		Returns(http.StatusOK, "Created", nil).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.AddMovie))

	wsWatchlist.Route(wsWatchlist.DELETE("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie to delete from actualUser's watchlist").DataType("int")).
		Doc("Remove movie from watchlist of actualUser").
		Returns(http.StatusNoContent, "Deleted", nil).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.DeleteMovie))

	wsWatchlist.Route(wsWatchlist.GET("/{movieId}").
		Param(wsWatchlist.PathParameter("movieId", "ID of the movie").DataType("int")).
		Doc("Check if movie is already in watchlist").
		Returns(http.StatusFound, "OK", true).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.AlreadyAdded))

	wsWatchlist.Route(wsWatchlist.GET("/").
		Doc("List all movies in watchlist").
		Param(wsWatchlist.QueryParameter("page", "num of page to get").DataType("int").DefaultValue("1")).
		Param(wsWatchlist.QueryParameter("pageSize", "number of element if one page").DataType("int").DefaultValue("10")).
		Returns(http.StatusOK, "OK", List{}).
		Returns(http.StatusUnauthorized, "Unauthorized, user cannot access this route", webServicePkg.FormattedJsonError{}).
		Returns(http.StatusNotFound, "Not found, impossible to find resource", webServicePkg.FormattedJsonError{}).
		Filter(webServicePkg.LogRequest()).
		Filter(webServicePkg.AuthenticateUser()).
		To(api.List))

	return
}

func (api handler) AddMovie(req *restful.Request, res *restful.Response) {
	// Get user info from token
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	var creation Creation
	creation.UserID = userFromRequest.ID
	movieIdStr := req.PathParameter("movieId")
	if movieIdStr != "" {
		creation.MovieID, err = utils.StrToID(movieIdStr)
		if err != nil {
			webServicePkg.HandleHTTPError(req, res, err)
			return
		}
	}

	err = api.service.AddMovie(creation)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (api handler) DeleteMovie(req *restful.Request, res *restful.Response) {
	// Get user info from token
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	var deletion Delete
	deletion.UserID = userFromRequest.ID
	movieIdStr := req.PathParameter("movieId")
	if movieIdStr != "" {
		deletion.MovieID, err = utils.StrToID(movieIdStr)
		if err != nil {
			webServicePkg.HandleHTTPError(req, res, err)
			return
		}
	}

	err = api.service.DeleteMovie(deletion)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusNoContent, "")
}

func (api handler) AlreadyAdded(req *restful.Request, res *restful.Response) {
	// Get user info from token
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	var check Check
	check.UserID = userFromRequest.ID
	movieIdStr := req.PathParameter("movieId")
	if movieIdStr != "" {
		check.MovieID, err = utils.StrToID(movieIdStr)
		if err != nil {
			webServicePkg.HandleHTTPError(req, res, err)
			return
		}
	}

	exists, err := api.service.AlreadyExists(check)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}
	if !exists {
		res.WriteHeaderAndEntity(http.StatusNotFound, exists)
		return
	}

	res.WriteHeaderAndEntity(http.StatusFound, exists)
}

func (api handler) List(req *restful.Request, res *restful.Response) {
	// Get user info from token
	userFromRequest, err := webServicePkg.WhoAmI(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	// get query parameters for filtering search
	var filters Filters
	filters.PaginationRequest, err = utils.ExtractPaginationRequest(req)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	filters.UserID = userFromRequest.ID

	list, err := api.service.List(filters)
	if err != nil {
		webServicePkg.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, list)
}
