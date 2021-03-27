package api

import (
	"cine-circle/internal/service"
	"github.com/emicklei/go-restful"
	"net/http"
)

func GetToken(req *restful.Request, res *restful.Response) {
	auth := req.HeaderParameter("Authorization")
	err, token := service.GetTokenFromAuthentication(auth)
	if err.IsNotNil() {
		res.WriteHeaderAndEntity(err.HttpCode(), err.CodeError())
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, token)
}
