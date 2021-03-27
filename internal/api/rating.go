package api

import (
	"cine-circle/internal/model"
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func AddRating(req *restful.Request, res *restful.Response) {
	movieId := req.PathParameter("movieId")
	_, username := service.CheckTokenAndGetUsername(req)
	var rating model.Rating
	err := req.ReadEntity(&rating)
	if err != nil {
		res.WriteHeaderAndEntity(model.ErrInternalApiUnprocessableEntity.HttpCode(),
			model.ErrInternalApiUnprocessableEntity.CodeError())
		return
	}
	err2, rating := service.AddRating(rating, movieId, username)
	if err2.IsNotNil() {
		res.WriteHeaderAndEntity(err2.HttpCode(), err2.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, rating)
}

