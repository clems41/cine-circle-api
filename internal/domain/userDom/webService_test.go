package userDom

import (
	"cine-circle-api/external/mailService/mailServiceMock"
	"cine-circle-api/internal/model/testSampler"
	"cine-circle-api/internal/repository"
	"cine-circle-api/internal/repository/postgres/pgModel"
	"cine-circle-api/pkg/httpServer"
	"cine-circle-api/pkg/httpServer/authentication"
	"cine-circle-api/pkg/httpServer/httpServerMock"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/test/setupTestCase"
	"cine-circle-api/pkg/utils/securityUtils"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"cine-circle-api/pkg/utils/testUtils/testRuler"
	"encoding/base64"
	"fmt"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"testing"
	"time"
)

/* Sign In */

func TestHandler_SignIn(t *testing.T) {
	testPath := basePath + signInPath
	_, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	password := fakeData.Password()
	user := sampler.GetUserWithPassword(password)
	wrongUsername := strings.ToLower(fakeData.UniqueUsername())
	wrongEmail := fakeData.UniqueEmail()
	wrongPassword := fakeData.Password()

	// Try without authorization header --> NOK 400
	resp := httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, nil)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with wrong login (username) et wrong password --> NOK 401
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(wrongUsername, wrongPassword))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with wrong login (email) et wrong password --> NOK 401
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(wrongEmail, wrongPassword))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with wrong login (username) et correct password --> NOK 401
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(wrongUsername, password))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with wrong login (email) et correct password --> NOK 401
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(wrongEmail, password))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct login (username) et wrong password --> NOK 401
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(user.Username, wrongPassword))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct login (email) et wrong password --> NOK 401
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(user.Email, wrongPassword))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct login (username) et correct password --> OK 200
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(user.Username, password))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Try with correct login (email) et correct password --> OK 200
	resp = httpMock.SendRequestWithHeaderParameters(testPath, http.MethodPost, basicAuthHeader(user.Email, password))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view SignInView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":             user.ID,
			"FirstName":      user.FirstName,
			"LastName":       user.LastName,
			"Email":          user.Email,
			"Username":       user.Username,
			"EmailConfirmed": user.EmailConfirmed,
		},
		"Token": map[string]interface{}{
			"TokenString":    testRuler.NotEmptyField{},
			"ExpirationDate": testRuler.NotEmptyField{},
		},
	})

	// Check that expirationDate is not in the past and not more than 3 days
	require.True(t, view.Token.ExpirationDate.After(time.Now()) && view.Token.ExpirationDate.Before(time.Now().AddDate(0, 0, 3)),
		"%s should be not in the past, meaning before %s, but also before 3 days from now %s",
		view.Token.ExpirationDate.String(), time.Now().String(), time.Now().AddDate(0, 0, 3).String())

	// Check that token claims contains userInfo (id)
	var userInfo authentication.UserInfo
	claims, err := httpServer.ValidateToken(view.Token.TokenString)
	require.NoError(t, err)
	err = httpServer.GetUserInfoFromTokenClaims(claims, &userInfo)
	require.NoError(t, err)
	require.Equal(t, user.ID, userInfo.Id)
}

/* Sign Up */

func TestHandler_SignUp(t *testing.T) {
	testPath := basePath + signUpPath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	usernameAlreadyTaken := fakeData.UniqueUsername()
	sampler.GetUserWithUsername(usernameAlreadyTaken)
	emailAlreadyTaken := fakeData.UniqueEmail()
	sampler.GetUserWithEmail(emailAlreadyTaken)
	correctForm := SignUpForm{
		CommonForm: CommonForm{
			FirstName: fake.FirstName(),
			LastName:  fake.LastName(),
			Email:     fakeData.UniqueEmail(),
			Username:  strings.ToLower(fakeData.UniqueUsername()),
		},
		Password: fakeData.Password(),
	}

	// Try with missing firstname --> NOK 400
	wrongForm := correctForm
	wrongForm.FirstName = ""
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing lastname --> NOK 400
	wrongForm = correctForm
	wrongForm.LastName = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing email --> NOK 400
	wrongForm = correctForm
	wrongForm.Email = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing password --> NOK 400
	wrongForm = correctForm
	wrongForm.Password = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing username --> NOK 400
	wrongForm = correctForm
	wrongForm.Username = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid password (too short 6 < 8) --> NOK 400
	wrongForm = correctForm
	wrongForm.Password = fakeData.Password()[:6]
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid email (not email type) --> NOK 400
	wrongForm = correctForm
	wrongForm.Email = fake.Country()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid firstname (not only alpha) --> NOK 400
	wrongForm = correctForm
	wrongForm.FirstName = "toto*tata"
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid lastname (not only alpha) --> NOK 400
	wrongForm = correctForm
	wrongForm.LastName = "toto*tata"
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid username (not only lowercase) --> NOK 400
	wrongForm = correctForm
	wrongForm.FirstName = strings.ToLower(strings.ToLower(fakeData.UniqueUsername()))
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with username already taken --> NOK 400
	wrongForm = correctForm
	wrongForm.Username = usernameAlreadyTaken
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with email already taken --> NOK 400
	wrongForm = correctForm
	wrongForm.Email = emailAlreadyTaken
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with correct form --> OK 200
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Check view fields
	var view SignUpView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":             testRuler.NotEmptyField{},
			"FirstName":      correctForm.FirstName,
			"LastName":       correctForm.LastName,
			"Email":          correctForm.Email,
			"Username":       correctForm.Username,
			"EmailConfirmed": false,
		},
	})

	// Check that user has been successfully stored into database
	var user pgModel.User
	err := db.
		Take(&user, view.Id).
		Error
	require.NoError(t, err)
	require.Equal(t, view.Id, user.ID)
	require.Equal(t, view.Username, user.Username)
	require.Equal(t, view.LastName, user.LastName)
	require.Equal(t, view.FirstName, user.FirstName)
	require.Equal(t, view.Email, user.Email)

	// Check that password has been successfully hashed
	require.NoError(t, securityUtils.CompareHashAndPassword(user.HashedPassword, correctForm.Password))
}

/* Send confirmation email */

func TestHandler_SendConfirmationEmail(t *testing.T) {
	testPath := basePath + sendConfirmationEmailPath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(testPath, http.MethodGet)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with user authenticated --> OK 200
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(testPath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that emailToken has been created in database
	err := db.
		Take(&user, user.ID).
		Error
	require.NoError(t, err)
	require.NotEmpty(t, user.EmailToken, "EmailToken should not be empty because confirmation email has been sent")
}

/* Confirm email */

func TestHandler_ConfirmEmail(t *testing.T) {
	testPath := basePath + confirmEmailPath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUserWithEmailToken()
	correctForm := ConfirmEmailForm{
		EmailToken: user.EmailToken,
	}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with wrong emailToken --> NOK 401
	httpMock.AuthenticateUserPermanently(user)
	wrongForm := correctForm
	wrongForm.EmailToken = fake.Country()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct emailToken --> OK 200
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, correctForm)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that email has been confirmed in database
	err := db.
		Take(&user, user.ID).
		Error
	require.NoError(t, err)
	require.True(t, user.EmailConfirmed, "Email should be confirmed")
}

/* Send reset password email */

func TestHandler_SendResetPasswordEmail(t *testing.T) {
	testPath := basePath + sendResetPasswordEmailPath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	wrongLogin := fakeData.UniqueEmail()

	// Try with non-existing login --> NOK 401
	resp := httpMock.SendRequest(testPath+"/"+wrongLogin, http.MethodGet)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct login (username) --> OK 200
	resp = httpMock.SendRequest(testPath+"/"+user.Username, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Try with correct login (email) --> OK 200
	resp = httpMock.SendRequest(testPath+"/"+user.Email, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that passwordToken has been created in database
	err := db.
		Take(&user, user.ID).
		Error
	require.NoError(t, err)
	require.NotEmpty(t, user.PasswordToken, "PasswordToken should not be empty because reset password email has been sent")
}

/* Reset password */

func TestHandler_ResetPassword(t *testing.T) {
	testPath := basePath + resetPasswordPath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUserWithPasswordToken()
	wrongPasswordToken := fake.Country()
	newPassword := fakeData.Password()
	correctFormUsername := ResetPasswordForm{
		PasswordToken: user.PasswordToken,
		Login:         user.Username,
		NewPassword:   newPassword,
	}
	correctFormEmail := ResetPasswordForm{
		PasswordToken: user.PasswordToken,
		Login:         user.Email,
		NewPassword:   newPassword,
	}

	// Try with wrong passwordToken and wrong login (username) --> NOK 401
	wrongForm := correctFormUsername
	wrongForm.PasswordToken = wrongPasswordToken
	wrongForm.Login = strings.ToLower(fakeData.UniqueUsername())
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with wrong passwordToken and wrong login (email) --> NOK 401
	wrongForm = correctFormEmail
	wrongForm.PasswordToken = wrongPasswordToken
	wrongForm.Login = fakeData.UniqueEmail()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with wrong passwordToken and correct login (username) --> NOK 401
	wrongForm = correctFormUsername
	wrongForm.PasswordToken = wrongPasswordToken
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with wrong passwordToken and correct login (email) --> NOK 401
	wrongForm = correctFormEmail
	wrongForm.PasswordToken = wrongPasswordToken
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct passwordToken and wrong login (username) --> NOK 401
	wrongForm = correctFormUsername
	wrongForm.Login = strings.ToLower(fakeData.UniqueUsername())
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct passwordToken and wrong login (email) --> NOK 401
	wrongForm = correctFormEmail
	wrongForm.Login = fakeData.UniqueEmail()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct passwordToken and correct login but newPassword too short (6 < 8) --> NOK 400
	wrongForm = correctFormEmail
	wrongForm.PasswordToken = fakeData.Password()[:6]
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct passwordToken and correct login (username) --> OK 200
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, correctFormUsername)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that password has been updated in database
	err := db.
		Take(&user, user.ID).
		Error
	require.NoError(t, err)
	require.NoError(t, securityUtils.CompareHashAndPassword(user.HashedPassword, correctFormEmail.NewPassword),
		"Hashed password %s should match with new password %s", user.HashedPassword, correctFormEmail.NewPassword)

	// PasswordToken has been removed, so we need to add new one to check with email as login
	user = sampler.GetUserWithPasswordToken()
	correctFormEmail.NewPassword = fakeData.Password()
	correctFormEmail.Login = user.Email
	correctFormEmail.PasswordToken = user.PasswordToken

	// Try with correct passwordToken and correct login (email) --> OK 200
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPost, correctFormEmail)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that password has been updated in database
	err = db.
		Take(&user, user.ID).
		Error
	require.NoError(t, err)
	require.NoError(t, securityUtils.CompareHashAndPassword(user.HashedPassword, correctFormEmail.NewPassword),
		"Hashed password %s should match with new password %s", user.HashedPassword, correctFormEmail.NewPassword)
}

/* Update password */

func TestHandler_UpdatePassword(t *testing.T) {
	testPath := basePath + updatePasswordPath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	oldPassword := fakeData.Password()
	newPassword := fakeData.Password()
	user := sampler.GetUserWithPassword(oldPassword)
	correctForm := UpdatePasswordForm{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with missing oldPassword --> NOK 400
	httpMock.AuthenticateUserPermanently(user)
	wrongForm := correctForm
	wrongForm.OldPassword = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing newPassword --> NOK 400
	wrongForm = correctForm
	wrongForm.NewPassword = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with wrong newPassword (too short 6 < 8) --> NOK 400
	wrongForm = correctForm
	wrongForm.NewPassword = newPassword[:6]
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with wrong oldPassword --> NOK 401
	wrongForm = correctForm
	wrongForm.OldPassword = fakeData.Password()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct form --> OK 200
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that password has been successfully updated into database
	err := db.
		Take(&user, user.ID).
		Error
	require.NoError(t, err)
	require.NoError(t, securityUtils.CompareHashAndPassword(user.HashedPassword, correctForm.NewPassword),
		"Hashed password %s should match with new password %s", user.HashedPassword, correctForm.NewPassword)
}

/* Get own info */

func TestHandler_GetOwnInfo(t *testing.T) {
	testPath := basePath + ownUserPath
	_, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequest(testPath, http.MethodGet)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with authenticated user --> OK 200
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequest(testPath, http.MethodGet)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view GetOwnInfoView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":             user.ID,
			"FirstName":      user.FirstName,
			"LastName":       user.LastName,
			"Email":          user.Email,
			"Username":       user.Username,
			"EmailConfirmed": user.EmailConfirmed,
		},
	})
}

/* Update */

func TestHandler_Update(t *testing.T) {
	testPath := basePath + ownUserPath
	db, httpMock, sampler, ruler, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	user := sampler.GetUser()
	correctForm := UpdateForm{
		CommonForm: CommonForm{
			FirstName: fake.FirstName(),
			LastName:  fake.LastName(),
			Email:     fakeData.UniqueEmail(),
			Username:  strings.ToLower(fakeData.UniqueUsername()),
		},
	}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(testPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with missing firstname --> NOK 400
	httpMock.AuthenticateUserPermanently(user)
	wrongForm := correctForm
	wrongForm.FirstName = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing lastname --> NOK 400
	wrongForm = correctForm
	wrongForm.LastName = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing email --> NOK 400
	wrongForm = correctForm
	wrongForm.Email = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with missing username --> NOK 400
	wrongForm = correctForm
	wrongForm.Username = ""
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid email (not email type) --> NOK 400
	wrongForm = correctForm
	wrongForm.Email = fake.Country()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid firstname (not only alpha) --> NOK 400
	wrongForm = correctForm
	wrongForm.FirstName = fakeData.UniqueEmail()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid lastname (not only alpha) --> NOK 400
	wrongForm = correctForm
	wrongForm.LastName = fakeData.UniqueEmail()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with invalid username (not only lowercase) --> NOK 400
	wrongForm = correctForm
	wrongForm.Username = strings.ToUpper(strings.ToLower(fakeData.UniqueUsername()))
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, wrongForm)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with correct form --> OK 200
	resp = httpMock.SendRequestWithBody(testPath, http.MethodPut, correctForm)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view fields
	var view UpdateView
	httpMock.DecodeResponse(resp, &view)
	ruler.CheckStruct(view, map[string]interface{}{
		"CommonView": map[string]interface{}{
			"Id":             user.ID,
			"FirstName":      correctForm.FirstName,
			"LastName":       correctForm.LastName,
			"Email":          correctForm.Email,
			"Username":       correctForm.Username,
			"EmailConfirmed": user.EmailConfirmed,
		},
	})

	// Check that user info have been successfully updated in database
	err := db.
		Take(&user, user.ID).
		Error
	require.NoError(t, err)
	require.Equal(t, correctForm.FirstName, user.FirstName)
	require.Equal(t, correctForm.LastName, user.LastName)
	require.Equal(t, correctForm.Email, user.Email)
	require.Equal(t, correctForm.Username, user.Username)
}

/* Delete */

func TestHandler_Delete(t *testing.T) {
	testPath := basePath + ownUserPath
	db, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	password := fakeData.Password()
	user := sampler.GetUserWithPassword(password)
	correctForm := DeleteForm{
		Password: password,
	}

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithBody(testPath, http.MethodDelete, correctForm)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with wrong password --> NOK 401
	httpMock.AuthenticateUserPermanently(user)
	wrongForm := correctForm
	wrongForm.Password = fakeData.Password()
	resp = httpMock.SendRequestWithBody(testPath, http.MethodDelete, wrongForm)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Try with correct password --> OK 200
	resp = httpMock.SendRequestWithBody(testPath, http.MethodDelete, correctForm)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check that user has been deleted
	err := db.
		Take(&user, user.ID).
		Error
	require.Error(t, err)
}

func TestHandler_Search(t *testing.T) {
	testPath := basePath
	_, httpMock, sampler, _, tearDown := setupTestcase(t, true)
	defer tearDown()

	// Create testing data
	keyword := fake.Word() + fake.Word()
	user := sampler.GetUser()
	var users []*pgModel.User
	nbUsersFirstnameMatching := 6
	nbUsersLastnameMatching := 3
	nbUsersUsernameMatching := 4
	for idx := range fakeData.FakeRange(18, 30) { // create between 18 and 30 circles (random number) but only 6 + 3 + 4 with matching keyword
		if idx < nbUsersFirstnameMatching {
			users = append(users, sampler.GetUserWithFirstName(fakeData.UniqueUsername()+keyword+fakeData.UniqueUsername()))
		} else if idx < nbUsersFirstnameMatching+nbUsersLastnameMatching {
			users = append(users, sampler.GetUserWithLastName(keyword+fakeData.UniqueUsername()))
		} else if idx < nbUsersFirstnameMatching+nbUsersLastnameMatching+nbUsersUsernameMatching {
			users = append(users, sampler.GetUserWithUsername(fakeData.UniqueUsername()+keyword))
		}
	}
	page := 1
	pageSize := 20 // should fit on one page with keyword

	// Try without authentication --> NOK 403
	resp := httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, keyword))
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Try with no keyword --> NOK 400
	wrongKeyword := ""
	httpMock.AuthenticateUserPermanently(user)
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, wrongKeyword))
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with incorrect keyword (len < 3) --> NOK 400
	wrongKeyword = "ad"
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, wrongKeyword))
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Try with correct keyword --> OK 200
	resp = httpMock.SendRequestWithQueryParameters(testPath, http.MethodGet, searchQueryParameters(page, pageSize, keyword))
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Check view that contains all users matching keyword
	var view SearchView
	httpMock.DecodeResponse(resp, &view)
	require.Equal(t, page, view.CurrentPage)
	require.Equal(t, pageSize, view.PageSize)
	require.Equal(t, 1, view.NumberOfPages)
	require.Equal(t, nbUsersFirstnameMatching+nbUsersLastnameMatching+nbUsersUsernameMatching, view.NumberOfItems)
	require.Equal(t, nbUsersFirstnameMatching+nbUsersLastnameMatching+nbUsersUsernameMatching, len(view.Users))

	// Check that users contains keyword
	for _, userView := range view.Users {
		match := false
		if strings.Contains(strings.ToLower(userView.FirstName), strings.ToLower(keyword)) {
			match = true
		}
		if strings.Contains(strings.ToLower(userView.LastName), strings.ToLower(keyword)) {
			match = true
		}
		if strings.Contains(strings.ToLower(userView.Username), strings.ToLower(keyword)) {
			match = true
		}
		require.True(t, match, "userView %v should match with keyword %s", userView, keyword)
	}
}

// searchQueryParameters will return queryParameters based on search parameters
func searchQueryParameters(page, pageSize int, keyword string) (queryParameters map[string][]string) {
	queryParameters = map[string][]string{
		pageQueryParameter.Name:     {fmt.Sprintf("%d", page)},
		pageSizeQueryParameter.Name: {fmt.Sprintf("%d", pageSize)},
		keywordQueryParameter.Name:  {keyword},
	}
	return
}

// setupTestcase will instantiate project and return all objects that can be needed for testing
func setupTestcase(t *testing.T, populateDatabase bool) (db *gorm.DB, httpMock *httpServerMock.Server, sampler *testSampler.Sampler, ruler *testRuler.Ruler, tearDown func()) {
	db, tearDown = setupTestCase.OpenCleanDatabaseFromTemplate(t)
	repo := repository.New(db)
	mailMock := mailServiceMock.New()
	svc := NewService(mailMock, repo)
	ws := NewHandler(svc)
	httpMock = httpServerMock.New(t, logger.Logger(), ws)
	sampler = testSampler.New(t, db, populateDatabase)
	ruler = testRuler.New(t)
	return
}

// basicAuth return Authorization header with Basic Authentication header (ex: Basic bG9naW46cGFzc3dvcmQK)
func basicAuthHeader(username, password string) map[string]string {
	auth := username + ":" + password
	return map[string]string{"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))}
}
