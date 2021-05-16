package webService

import (
	"cine-circle/internal/utils"
	"cine-circle/pkg/logger"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
)

// AuthenticateUser : Add filter for getting user infos (token, ID, etc...) in order to authenticate him
func AuthenticateUser() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		token, err := GetTokenFromAuthenticationHeader(req)
		if err != nil {
			HandleHTTPError(req, res, err)
			return
		}
		_, err = CheckToken(token)
		if err != nil {
			HandleHTTPError(req, res, err)
			return
		} else {
			chain.ProcessFilter(req, res)
		}
	}
	return filter
}

// LogRequest : Add filter for logging request
func LogRequest() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		logger.Sugar.Debugf("%s - %s", req.Request.Method, req.Request.URL.String())
		chain.ProcessFilter(req, res)
	}
	return filter
}

// CheckRoles : Add filter for checking if user has correct role for using this resource
func CheckRoles(resourceRoles []string) func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		user, err := ActualUserHandler.WhoAmI(req)
		if err != nil {
			HandleHTTPError(req, res, err)
			return
		}
		if utils.StrSlicesHaveAtLeastOneMatch(resourceRoles, user.Roles) {
			chain.ProcessFilter(req, res)
		} else {
			// TODO update error with custom one
			err = errors.New("user cannot access this resource")
			HandleHTTPError(req, res, err)
			return
		}
	}
	return filter

}
