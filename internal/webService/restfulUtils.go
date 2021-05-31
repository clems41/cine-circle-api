package webService

import (
	"cine-circle/internal/constant"
	"cine-circle/pkg/logger"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
)

type Handler interface {
	WebServices() []*restful.WebService
}

// AddHandlersToRestfulContainer : add webServices into restful container from handlers
func AddHandlersToRestfulContainer(container *restful.Container, handlers ...Handler) {
	for _, handler := range handlers {
		webServices := handler.WebServices()
		for _, webService := range webServices {
			container.Add(webService)
			for _, route := range webService.Routes() {
				logger.Sugar.Infof("%s \t %s", route.Method, route.Path)
			}
		}
	}
}

// WhoAmI : get user info from request
func WhoAmI(req *restful.Request) (user ActualUser, err error) {
	userInfoStr := req.HeaderParameter(constant.UserInfoRequestParameter)
	err = json.Unmarshal([]byte(userInfoStr), &user)
	if err != nil {
		return user, errors.WithStack(err)
	}
	return
}
