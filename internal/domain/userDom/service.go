package userDom

import (
	"cine-circle-api/internal/repository/instance/userRepository"
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/internal/service/mailService"
	"cine-circle-api/pkg/httpServer"
	"cine-circle-api/pkg/httpServer/authentication"
	"cine-circle-api/pkg/sql/gormUtils"
	"cine-circle-api/pkg/utils/securityUtils"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
)

var _ Service = (*service)(nil)

type Service interface {
	SignIn(form SignInForm) (view SignInView, err error)
	SignUp(form SignUpForm) (view SignUpView, err error)
	SendEmailConfirmation(form SendEmailConfirmationForm) (err error)
	ConfirmEmail(form ConfirmEmailForm) (err error)
	SendResetPasswordEmail(form SendResetPasswordEmailForm) (err error)
	ResetPassword(form ResetPasswordForm) (err error)
	GetOwnInfo(form GetOwnInfoForm) (view GetOwnInfoView, err error)
	Update(form UpdateForm) (view UpdateView, err error)
	UpdatePassword(form UpdatePasswordForm) (err error)
	Delete(form DeleteForm) (err error)
	Search(form SearchForm) (view SearchView, err error)
}

type service struct {
	repository  userRepository.Repository
	serviceMail mailService.Service
}

func NewService(serviceMail mailService.Service, repository userRepository.Repository) *service {
	return &service{
		repository:  repository,
		serviceMail: serviceMail,
	}
}

/* Login */

func (svc *service) SignIn(form SignInForm) (view SignInView, err error) {
	// Check if user exists
	user, ok, err := svc.repository.GetUserFromLogin(form.Login)
	if err != nil {
		return
	}
	if !ok {
		return view, errBadCredentials
	}

	// Check if password is correct
	err = securityUtils.CompareHashAndPassword(user.HashedPassword, form.Password)
	if err != nil {
		return view, errBadCredentials
	}

	// Generate token and view
	userInfo := authentication.UserInfo{
		Id: user.ID,
	}
	token, err := httpServer.GenerateTokenWithUserInfo(userInfo)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(user)
	view.Token = Token{
		ExpirationDate: token.ExpirationDate,
		TokenString:    token.TokenString,
	}
	return
}

func (svc *service) SignUp(form SignUpForm) (view SignUpView, err error) {
	hashedPassword, err := securityUtils.HashAndSaltPassword(form.Password)
	if err != nil {
		return view, errors.WithStack(err)
	}
	user := model.User{
		Username:       form.Username,
		HashedPassword: hashedPassword,
		LastName:       form.LastName,
		FirstName:      form.FirstName,
		Email:          form.Email,
		EmailConfirmed: false,
	}
	err = svc.repository.Create(&user)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(user)
	return
}

func (svc *service) SendEmailConfirmation(form SendEmailConfirmationForm) (err error) {
	// Get user info
	user, ok, err := svc.repository.GetUser(form.UserId)
	if err != nil {
		return
	}
	if !ok {
		return errUserUnauthorized
	}

	// Generate token
	token := fmt.Sprintf("%d", uuid.New().ID())

	// Generate and send email
	emailBody := strings.Replace(sendConfirmationEmailBody, sendConfirmationEmailFirstNameJoker, user.FirstName, 1)
	emailBody = strings.Replace(emailBody, sendConfirmationEmailTokenJoker, token, 1)
	emailForm := mailService.SendEmailForm{
		To:      []string{user.Email},
		Subject: sendConfirmationEmailSubject,
		Message: emailBody,
		Tags:    sendConfirmationEmailTags,
		Html:    false,
	}
	err = svc.serviceMail.Send(emailForm)
	if err != nil {
		return
	}

	// Save token into database
	user.EmailToken = token
	err = svc.repository.Save(&user)
	if err != nil {
		return
	}

	return
}

func (svc *service) ConfirmEmail(form ConfirmEmailForm) (err error) {
	// Get user info
	user, ok, err := svc.repository.GetUser(form.UserId)
	if err != nil {
		return
	}
	if !ok {
		return errUserUnauthorized
	}

	// Check if token is correct
	if user.EmailToken != form.EmailToken {
		return errUserUnauthorized
	}

	// Remove token in database and mark email as confirmed
	user.EmailToken = ""
	user.EmailConfirmed = true
	err = svc.repository.Save(&user)
	if err != nil {
		return
	}
	return
}

func (svc *service) SendResetPasswordEmail(form SendResetPasswordEmailForm) (err error) {
	// Get user info
	user, ok, err := svc.repository.GetUserFromLogin(form.Login)
	if err != nil {
		return
	}
	if !ok {
		return errUserUnauthorized
	}

	// Generate token
	token := fmt.Sprintf("%d", uuid.New().ID())

	// Generate and send email
	emailBody := strings.Replace(sendResetPasswordEmailBody, sendResetPasswordEmailFirstNameJoker, user.FirstName, 1)
	emailBody = strings.Replace(emailBody, sendResetPasswordTokenJoker, token, 1)
	emailForm := mailService.SendEmailForm{
		To:      []string{user.Email},
		Subject: sendResetPasswordEmailSubject,
		Message: emailBody,
		Tags:    sendResetPasswordEmailTags,
		Html:    false,
	}
	err = svc.serviceMail.Send(emailForm)
	if err != nil {
		return
	}

	// Save token into database
	user.PasswordToken = token
	err = svc.repository.Save(&user)
	if err != nil {
		return
	}
	return
}

func (svc *service) ResetPassword(form ResetPasswordForm) (err error) { // Get user info
	user, ok, err := svc.repository.GetUserFromLogin(form.Login)
	if err != nil {
		return
	}
	if !ok {
		return errUserUnauthorized
	}

	// Check if token is correct
	if user.PasswordToken != form.PasswordToken {
		return errUserUnauthorized
	}

	// Remove token in database and update password with new one
	newHashedPassword, err := securityUtils.HashAndSaltPassword(form.NewPassword)
	if err != nil {
		return
	}
	user.PasswordToken = ""
	user.HashedPassword = newHashedPassword
	err = svc.repository.Save(&user)
	if err != nil {
		return
	}
	return
}

func (svc *service) GetOwnInfo(form GetOwnInfoForm) (view GetOwnInfoView, err error) {
	user, ok, err := svc.repository.GetUser(form.UserId)
	if err != nil {
		return
	}
	if !ok {
		return view, errUserUnauthorized
	}
	view.CommonView = svc.fromModelToView(user)
	return
}

func (svc *service) Update(form UpdateForm) (view UpdateView, err error) {
	user, ok, err := svc.repository.GetUser(form.UserId)
	if err != nil {
		return
	}
	if !ok {
		return view, errUserUnauthorized
	}
	user.Username = form.Username
	user.LastName = form.LastName
	user.FirstName = form.FirstName
	user.Email = form.Email
	err = svc.repository.Save(&user)
	if err != nil {
		return
	}
	view.CommonView = svc.fromModelToView(user)
	return
}

func (svc *service) UpdatePassword(form UpdatePasswordForm) (err error) {
	// Check if old password is correct
	user, ok, err := svc.repository.GetUser(form.UserId)
	if err != nil {
		return
	}
	if !ok {
		return errUserUnauthorized
	}
	err = securityUtils.CompareHashAndPassword(user.HashedPassword, form.OldPassword)
	if err != nil {
		return errUserUnauthorized
	}

	// Update new one
	newHashedPassword, err := securityUtils.HashAndSaltPassword(form.NewPassword)
	if err != nil {
		return
	}
	user.HashedPassword = newHashedPassword
	err = svc.repository.Save(&user)
	if err != nil {
		return
	}
	return
}

func (svc *service) Delete(form DeleteForm) (err error) {
	// Check if password is correct
	user, ok, err := svc.repository.GetUser(form.UserId)
	if err != nil {
		return
	}
	if !ok {
		return errUserUnauthorized
	}
	err = securityUtils.CompareHashAndPassword(user.HashedPassword, form.Password)
	if err != nil {
		return errUserUnauthorized
	}

	ok, err = svc.repository.Delete(form.UserId)
	if err != nil {
		return
	}
	if !ok {
		return errUserCannotBeDeleted
	}
	return
}

func (svc *service) Search(form SearchForm) (view SearchView, err error) {
	repoForm := userRepository.SearchForm{
		PaginationQuery: gormUtils.PaginationQuery{
			Page:     form.Page,
			PageSize: form.PageSize,
		},
		Keyword: form.Keyword,
	}

	repoView, err := svc.repository.Search(repoForm)
	if err != nil {
		return
	}

	view.Page = form.BuildResult(repoView.Total)
	view.Users = make([]CommonView, 0)

	for _, user := range repoView.Users {
		view.Users = append(view.Users, svc.fromModelToView(user))
	}
	return
}

/* Private methods below */

func (svc *service) fromModelToView(user model.User) (view CommonView) {
	view = CommonView{
		Id:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Username:       user.Username,
		EmailConfirmed: user.EmailConfirmed,
	}
	return
}

func (svc *service) fromFormToModel(form CommonForm) (user model.User) {
	user = model.User{
		Username:  form.Username,
		LastName:  form.LastName,
		FirstName: form.FirstName,
		Email:     form.Email,
	}
	return
}
