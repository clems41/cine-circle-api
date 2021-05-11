package webService

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/domain"
	"cine-circle/internal/typedErrors"
	"cine-circle/internal/utils"
	"cine-circle/pkg/logger"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
	"strings"
)

type Handler interface {
	WebServices() []*restful.WebService
}

// AddHandler : add webServices into restful container from handlers
func AddHandler(container *restful.Container, handlers ...Handler) {
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

// WhoAmI : return username from token
func WhoAmI(req *restful.Request) (userID domain.IDType, err error) {
	token, err := GetTokenFromAuthenticationHeader(req)
	if err != nil {
		return
	}
	claims, err := CheckToken(token)
	if err != nil {
		return
	}
	subClaims := fmt.Sprintf("%v", claims["sub"])
	id, err := strconv.Atoi(subClaims)
	if err != nil {
		return userID, typedErrors.NewApiBadCredentialsErrorf("cannot find userID from token, got claims = %s", subClaims)
	}
	userID = domain.IDType(id)
	return
}

// HandleHTTPError : fill response with right error code in case of custom errors
func HandleHTTPError(req *restful.Request, res *restful.Response, err error) {
	e, ok := err.(typedErrors.CustomError)
	if ok {
		e.Print()
		res.WriteHeaderAndEntity(e.HttpCode(), e.CodeError())
	} else {
		res.WriteHeaderAndEntity(http.StatusInternalServerError, err.Error())
	}
	logger.Sugar.Errorf("Error occurs for request %s - %s : %s", req.Request.Method, req.Request.RequestURI, err.Error())
}

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
		logger.Sugar.Debugf("%s - %s - ", req.Request.Method, req.Request.URL.String())
		chain.ProcessFilter(req, res)
	}
	return filter
}

// GetTokenFromAuthenticationHeader : return token in header request
func GetTokenFromAuthenticationHeader(req *restful.Request) (token string, err error) {
	tokenHeader := req.HeaderParameter(constant.TokenHeader)
	if tokenHeader == "" {
		return token, typedErrors.ErrApiUserBadCredentials
	}
	res := strings.Split(tokenHeader, constant.BearerTokenDelimiterForHeader)
	if len(res) != 2 {
		return token, typedErrors.ErrApiUserBadCredentials
	}
	if res[0] != constant.TokenKind {
		return token, typedErrors.ErrApiUserBadCredentials
	}
	token = res[1]
	return
}

// CheckToken : check validity of token from header request and return claims
func CheckToken(token string) (claims jwt.MapClaims, err error) {
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tokenKey := utils.GetDefaultOrFromEnv(constant.SecretTokenDefault, constant.SecretTokenEnv)
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, typedErrors.NewApiBadCredentialsError(err)
		}
	}
	if tkn == nil || !tkn.Valid {
		logger.Sugar.Debugf("Error while getting token : Token not valid")
		return claims, typedErrors.ErrApiUserBadCredentials
	}
	return
}
