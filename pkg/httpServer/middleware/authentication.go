package middleware

import (
	"cine-circle-api/pkg/customError"
	"cine-circle-api/pkg/httpServer"
	"cine-circle-api/pkg/httpServer/authentication"
	"cine-circle-api/pkg/httpServer/httpError"
	"github.com/emicklei/go-restful"
)

// AuthenticateUser must be used on all secured endpoints. It will parse and check token provided in Authorization header.
// UserInfo (id, username, etc...) got from token claims will be added into request context in order to be accessible later without parsing token twice.
func AuthenticateUser() (filter func(*restful.Request, *restful.Response, *restful.FilterChain)) {
	filter = func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		// Parse and check token, then retrieve userInfo from it
		var userInfo authentication.UserInfo
		err := httpServer.ValidateTokenAndGetUserInfo(req, &userInfo)
		if err != nil {
			httpError.HandleHTTPError(req, res, customError.NewForbidden().WrapCode(invalidTokenErrorCode).WrapError(err))
			return
		}

		// Set userId in request context. in this way, we can retrieve it later without parsing token a second time.
		authentication.SetUserContextInRequest(req, userInfo)

		// If everything ok, process request
		chain.ProcessFilter(req, res)
	}
	return filter
}
