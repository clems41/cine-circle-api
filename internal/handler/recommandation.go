package handler

import (
	"cine-circle/internal/domain/recommendationDom"
	"github.com/emicklei/go-restful"
)

type recommendationHandler struct {
	service recommendationDom.Service
}

func NewRecommendationHandler(svc recommendationDom.Service) *recommendationHandler {
	return &recommendationHandler{
		service:    svc,
	}
}

func (api recommendationHandler) WebService() *restful.WebService {
	wsReco := &restful.WebService{}
	wsReco.Path("/v1/recommendations")

/*	wsReco.Route(wsReco.POST("/{movieId}").
		Param(wsReco.PathParameter("movieId", "ID of the movie to rate").DataType("int")).
		Doc("Add rating to movie for specific user").
		Writes(model.Rating{}).
		Returns(201, "Created", model.Rating{}).
		Returns(400, "Bad request, fields not validated", typedErrors.ErrApiBadRequest.CodeError()).
		Returns(422, "Not processable, impossible to serialize json to Rating",
			typedErrors.ErrApiUnprocessableEntity.CodeError()).
		Filter(authenticateUser(true)).
		To(AddRating))*/

	return wsReco
}

