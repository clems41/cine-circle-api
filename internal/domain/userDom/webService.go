package userDom

import (
	"cine-circle-api/internal/constant/swaggerConst"
	"cine-circle-api/internal/constant/webserviceConst"
	"cine-circle-api/pkg/customError"
	"cine-circle-api/pkg/httpServer/authentication"
	"cine-circle-api/pkg/httpServer/httpError"
	"cine-circle-api/pkg/httpServer/middleware"
	"cine-circle-api/pkg/utils/validationUtils"
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		service: svc,
	}
}

func (hd *handler) WebService() *restful.WebService {
	wsUser := new(restful.WebService)
	tags := []string{swaggerConst.UserTag}

	wsUser.Path(basePath)

	/* Login */

	wsUser.Route(wsUser.POST(signInPath).
		Produces(restful.MIME_JSON).
		Doc("Génération d'un token pour un user via une Basic Authentication (login:mdp)").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(SignInView{}).
		Returns(http.StatusOK, "Token généré", SignInView{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		To(hd.SignIn))

	wsUser.Route(wsUser.POST(signUpPath).
		Consumes(restful.MIME_JSON).
		Doc("Création de compte pour un nouvel user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(SignUpForm{}).
		Writes(SignUpView{}).
		Returns(http.StatusCreated, "Utilisateur créé", SignUpView{}).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		To(hd.SignUp))

	wsUser.Route(wsUser.GET(sendConfirmationEmailPath).
		Doc("Envoie d'un email pour confirmer son adresse à l'user authentifié").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Email envoyé", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.SendConfirmationEmail))

	wsUser.Route(wsUser.POST(confirmEmailPath).
		Doc("Confirmation de l'email en utilisant le token fourni dans l'email lors de l'envoie d'email de confirmation").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Email confirmé", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.ConfirmEmail))

	wsUser.Route(wsUser.GET(sendResetPasswordEmailPath+"/"+loginPathParameter.Joker()).
		Param(wsUser.PathParameter(loginPathParameter.String(), "Login (nom user ou email) de l'user qui souhaite réinitialiser son mot de passe").
			DataType("string").Required(true)).
		Doc("Envoie d'un email pour réinitialiser le mot de passe de l'user").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Email envoyé", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		To(hd.SendResetPasswordEmail))

	wsUser.Route(wsUser.POST(resetPasswordPath).
		Doc("Réinitialisation du mot de passe de l'user avec le token reçu dans le mail").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Mot de passe réinitialisé", nil).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		To(hd.ResetPassword))

	wsUser.Route(wsUser.PUT(updatePasswordPath).
		Consumes(restful.MIME_JSON).
		Doc("Modification du mot de passe de l'user authentifié par l'user authentifié").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(UpdatePasswordForm{}).
		Returns(http.StatusOK, "Mot de passe l'user authentifié mis à jour", nil).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.UpdatePassword))

	wsUser.Route(wsUser.PUT("/").
		Consumes(restful.MIME_JSON).
		Writes(restful.MIME_JSON).
		Doc("Modification des informations de l'user authentifié").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(UpdateForm{}).
		Returns(http.StatusOK, "Informations mises à jour", nil).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Update))

	wsUser.Route(wsUser.DELETE("/").
		Consumes(restful.MIME_JSON).
		Doc("Suppression du compte de l'user authentifié").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(UpdateForm{}).
		Returns(http.StatusOK, "Compte supprimé", nil).
		Returns(http.StatusBadRequest, webserviceConst.BadRequestMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Returns(http.StatusUnprocessableEntity, webserviceConst.UnprocessableEntityMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()).
		To(hd.Delete))

	wsUser.Route(wsUser.GET(ownInfoPath).
		Doc("Récupérer les informations de l'user actuel (token)").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(http.StatusOK, "Informations récupérées", GetOwnInfoView{}).
		Returns(http.StatusUnauthorized, webserviceConst.UnauthorizedMessage, httpError.FormattedJsonError{}).
		Filter(middleware.AuthenticateUser()). // Tous les users ont accès à cette route, pas besoin de préciser de rôle autorisé
		To(hd.GetOwnInfo))

	return wsUser
}

func (hd *handler) SignIn(req *restful.Request, res *restful.Response) {
	login, password, ok := req.Request.BasicAuth()
	if !ok {
		httpError.HandleHTTPError(req, res, errInvalidAuthorizationHeader)
		return
	}
	form := SignInForm{
		Password: password,
		Login:    login,
	}
	view, err := hd.service.SignIn(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) SignUp(req *restful.Request, res *restful.Response) {
	var form SignUpForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	view, err := hd.service.SignUp(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusCreated, view)
}

func (hd *handler) SendConfirmationEmail(req *restful.Request, res *restful.Response) {
	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	form := SendEmailConfirmationForm{
		UserId: user.Id,
	}
	err = hd.service.SendEmailConfirmation(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (hd *handler) ConfirmEmail(req *restful.Request, res *restful.Response) {
	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	var form ConfirmEmailForm
	err = validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}
	form.UserId = user.Id
	err = hd.service.ConfirmEmail(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (hd *handler) SendResetPasswordEmail(req *restful.Request, res *restful.Response) {
	login := req.PathParameter(loginPathParameter.String())
	form := SendResetPasswordEmailForm{Login: login}
	err := hd.service.SendResetPasswordEmail(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (hd *handler) ResetPassword(req *restful.Request, res *restful.Response) {
	var form ResetPasswordForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}

	err = hd.service.ResetPassword(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (hd *handler) UpdatePassword(req *restful.Request, res *restful.Response) {
	var form UpdatePasswordForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}

	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	form.UserId = user.Id

	err = hd.service.UpdatePassword(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, "")
}

func (hd *handler) Update(req *restful.Request, res *restful.Response) {
	var form UpdateForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}

	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	form.UserId = user.Id

	view, err := hd.service.Update(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}
	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) GetOwnInfo(req *restful.Request, res *restful.Response) {
	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	form := GetOwnInfoForm{UserId: user.Id}

	view, err := hd.service.GetOwnInfo(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, view)
}

func (hd *handler) Delete(req *restful.Request, res *restful.Response) {
	var form DeleteForm
	err := validationUtils.ReadEntityAndValidateStruct(req, &form)
	if err != nil {
		httpError.HandleHTTPError(req, res, customError.NewBadRequest().WrapError(err))
		return
	}

	user, err := authentication.WhoAmI(req)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	form.UserId = user.Id

	err = hd.service.Delete(form)
	if err != nil {
		httpError.HandleHTTPError(req, res, err)
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, "")
}
