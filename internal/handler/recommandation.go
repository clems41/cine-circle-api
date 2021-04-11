package handler

import (
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func NewRecommendationHandler() *restful.WebService {
	wsReco := &restful.WebService{}
	wsReco.Path("/v1/recommendations")

	wsReco.Route(wsReco.POST("/{movieId}").
		Param(wsReco.PathParameter("movieId", "ID of the movie to rate").DataType("int")).
		Doc("Add rating to movie for specific user").
		Writes(model.Rating{}).
		Returns(201, "Created", model.Rating{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Rating",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(filterUser(true)).
		To(AddRating))

	return wsReco
}

func AddRating(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	_, username := service.CheckTokenAndGetUsername(req)
	var rating model.Rating
	err := req.ReadEntity(&rating)
	if err != nil {
		res.WriteHeaderAndEntity(typedErrors.ErrApiUnprocessableEntity.HttpCode(),
			typedErrors.ErrApiUnprocessableEntity.CodeError())
		return
	}
	err2, rating := service.AddRating(rating, movieId, username)
	if err2.IsNotNil() {
		res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, rating)
}

