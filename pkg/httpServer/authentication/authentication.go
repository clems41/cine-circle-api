package authentication

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
)

// WhoAmI return UserInfo contains in request context.
// To be working, middleware.AuthenticateUser filters should be added in endpoint filters
func WhoAmI(req *restful.Request) (userInfo UserInfo, err error) {
	// Get request context then userContext value
	requestContext := req.Request.Context()
	var ok bool
	userInfo, ok = requestContext.Value(userContextKey).(UserInfo)
	if !ok {
		return userInfo, fmt.Errorf("UserInfo cannot be found in request context")
	}
	return
}

// SetUserContextInRequest will add userInfo struct into request context.
// In this way, userInfo could be retrieved later without parsing token twice.
func SetUserContextInRequest(request *restful.Request, userInfo UserInfo) {
	// Create new user Context based on request context
	requestContext := request.Request.Context()
	userContext := context.WithValue(requestContext, userContextKey, userInfo)
	// Add userContext in request
	request.Request = request.Request.WithContext(userContext)
}
