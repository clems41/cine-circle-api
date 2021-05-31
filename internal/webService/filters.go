package webService

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/utils"
	"cine-circle/pkg/logger"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
)

// AuthenticateUser : Add filter for getting user infos (token, ID, etc...) in order to authenticate him. Add also all user info into request.
func AuthenticateUser() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		token, err := GetTokenFromAuthenticationHeader(req)
		if err != nil {
			HandleHTTPError(req, res, err)
			return
		}
		claims, err := CheckToken(token)
		if err != nil {
			HandleHTTPError(req, res, err)
			return
		}
		// Check if user exists in database
		user, err := ActualUserHandler.GetUserFromClaims(claims)
		if err != nil {
			HandleHTTPError(req, res, err)
			return
		}
		// Adding user info into request
		bytes, err := json.Marshal(user)
		if err != nil {
			HandleHTTPError(req, res, err)
			return
		}
		req.Request.Header.Set(constant.UserInfoRequestParameter, string(bytes))
		chain.ProcessFilter(req, res)
	}
	return filter
}

// LogRequest : Add filter for logging request
func LogRequest() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		logger.Sugar.Debugf("\t%s - %s", req.Request.Method, req.Request.URL.String())
		chain.ProcessFilter(req, res)
	}
	return filter
}

// CheckRoles : Add filter for checking if user has correct role for using this resource
func CheckRoles(resourceRoles []string) func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		user, err := WhoAmI(req)
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
