package handler

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/logger"
	"cine-circle/internal/typedErrors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	"net/http"
	"strings"
)

var (
	CommonHandler *commonHandler
)

type handler interface {
	WebService() *restful.WebService
}

func AddWebService(container *restful.Container, handler handler) {
	container.Add(handler.WebService())

	rootPath := handler.WebService().RootPath()
	if rootPath == "" {
		rootPath = "/"
	}
	logger.Sugar.Infof("Routes for : %s", rootPath)
	for _, route := range handler.WebService().Routes() {
		logger.Sugar.Infof("%+v", route)
	}
}

type commonHandler struct {
	userService userDom.Service
}

func NewCommonHandler(userService userDom.Service) *commonHandler {
	return &commonHandler{
		userService: userService,
	}
}

func (handler *commonHandler) WhoAmI(req *restful.Request) (user userDom.Result, err error) {
	claims, err := checkToken(req)
	if err != nil {
		return
	}
	username := fmt.Sprintf("%v", claims["sub"])
	get := userDom.Get{Username: username}
	return handler.userService.Get(get)
}

func handleHTTPError(res *restful.Response, err error) {
	e, ok := err.(typedErrors.CustomError)
	if ok {
		e.Print()
		res.WriteHeaderAndEntity(e.HttpCode(), e.CodeError())
	} else {
		res.WriteHeaderAndEntity(http.StatusInternalServerError, "")
	}
}

// Add filter for getting user infos (token, ID, etc...) in order to authenticate him
func authenticateUser() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		_, err := checkToken(req)
		if err != nil {
			handleHTTPError(res, err)
			return
		} else {
			chain.ProcessFilter(req, res)
		}
	}
	return filter
}

// Add filter for logging request
func logRequest() func(*restful.Request, *restful.Response, *restful.FilterChain) {
	filter := func(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
		logger.Sugar.Debugf("%s - %s - ", req.Request.Method, req.Request.URL.String())
		chain.ProcessFilter(req, res)
	}
	return filter
}

func getTokenFromAuthenticationHeader(req *restful.Request) (token string, err error) {
	tokenHeader := req.HeaderParameter(constant.TokenHeader)
	if tokenHeader == "" {
		return token, typedErrors.ErrApiUserBadCredentials
	}
	res := strings.Split(tokenHeader, " ")
	if len(res) != 2 {
		return token, typedErrors.ErrApiUserBadCredentials
	}
	if res[0] != constant.TokenKind {
		return token, typedErrors.ErrApiUserBadCredentials
	}
	token = res[1]
	return
}

func checkToken(req *restful.Request) (claims jwt.MapClaims, err error) {
	token, err := getTokenFromAuthenticationHeader(req)
	if err != nil {
		return
	}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.TokenKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, typedErrors.NewApiBadCredentialsErrorf(err.Error())
		}
	}
	if tkn == nil || !tkn.Valid {
		logger.Sugar.Debugf("Error while getting token : Token not valid")
		return claims, typedErrors.ErrApiUserBadCredentials
	}
	return
}