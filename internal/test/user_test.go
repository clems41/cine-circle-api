package test

import (
	"cine-circle/internal/domain"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/repository"
	"cine-circle/internal/utils"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestUser_Create(t *testing.T) {
	db, clean := OpenDatabase(t)
	defer clean()

	userRepository := repository.NewUserRepository(db)
	userService := userDom.NewService(userRepository)

	// Fields for creation
	username := fake.UserName()
	displayName := fake.FullName()
	password := getFakePassword()
	email := fake.EmailAddress()

	// Expected Result when creation is ok
	expectedResult := userDom.Result{
		Username:    strings.ToLower(username),
		DisplayName: displayName,
		Email:       email,
	}

	// test creation with missing mandatory field : Username
	creation := userDom.Creation{
		Username:    "",
		DisplayName: displayName,
		Password:    password,
		Email:       email,
	}
	result, err := userService.Create(creation)
	require.Error(t, err, "Should return error but got nil")

	// test creation with missing mandatory field : DisplayName
	creation = userDom.Creation{
		Username:    username,
		DisplayName: "",
		Password:    password,
		Email:       email,
	}
	result, err = userService.Create(creation)
	require.Error(t, err, "Should return error but got nil")

	// test creation with missing mandatory field : Password
	creation = userDom.Creation{
		Username:    username,
		DisplayName: displayName,
		Password:    "",
		Email:       email,
	}
	result, err = userService.Create(creation)
	require.Error(t, err, "Should return error but got nil")

	// test creation with missing mandatory field : Email
	creation = userDom.Creation{
		Username:    username,
		DisplayName: displayName,
		Password:    password,
		Email:       "",
	}
	result, err = userService.Create(creation)
	require.Error(t, err, "Should return error but got nil")

	// test creation with all correct fields
	creation = userDom.Creation{
		Username:    username,
		DisplayName: displayName,
		Password:    password,
		Email:       email,
	}

	// check if return no err
	result, err = userService.Create(creation)
	require.NoError(t, err, "Should not return error but got %v", err)

	// check if result is like expected
	expectedResult.UserID = result.UserID
	require.Equal(t, expectedResult, result)

	// check if password has been correctly salt and hash
	hashedPassword, err := userRepository.GetHashedPassword(userDom.Get{UserID: result.UserID})
	require.NoError(t, err, "Should not return error but got %v", err)
	assert.NotEqual(t, creation.Password, hashedPassword, "password should be hashed")
	err = utils.CompareHashAndPassword(hashedPassword, creation.Password)
	require.NoError(t, err, "passwords should be the same (using hash comparison) but got %v", err)
}

func TestUser_Update(t *testing.T) {
	db, clean := OpenDatabase(t)
	defer clean()

	userRepository := repository.NewUserRepository(db)
	userService := userDom.NewService(userRepository)

	sampler := newSampler(t, db, false)

	// Add existing user to database
	userSample := sampler.getUserSample()

	// New fields to update user with
	displayName := fake.FullName()
	email := fake.EmailAddress()

	// Expected Result when update is ok
	expectedResult := userDom.Result{
		UserID:      userSample.GetID(),
		Username:    userSample.Username,
		DisplayName: displayName,
		Email:       email,
	}

	// test update with missing mandatory field : UserID
	update := userDom.Update{
		UserID:      0,
		DisplayName: displayName,
		Email:       email,
	}
	result, err := userService.Update(update)
	require.Error(t, err, "Should return error but got nil")

	// test update with missing mandatory field : DisplayName
	update = userDom.Update{
		UserID:      userSample.GetID(),
		DisplayName: "",
		Email:       email,
	}
	result, err = userService.Update(update)
	require.Error(t, err, "Should return error but got nil")

	// test update with missing mandatory field : Email
	update = userDom.Update{
		UserID:      userSample.GetID(),
		DisplayName: displayName,
		Email:       "",
	}
	result, err = userService.Update(update)
	require.Error(t, err, "Should return error but got nil")

	// test update with all correct fields
	update = userDom.Update{
		UserID:      userSample.GetID(),
		DisplayName: displayName,
		Email:       email,
	}

	// check if return no err
	result, err = userService.Update(update)
	require.NoError(t, err, "Should not return error but got %v", err)

	// check if result is like expected
	require.Equal(t, expectedResult, result)
}

func TestUser_UpdatePassword(t *testing.T) {
	db, clean := OpenDatabase(t)
	defer clean()

	userRepository := repository.NewUserRepository(db)
	userService := userDom.NewService(userRepository)

	sampler := newSampler(t, db, false)

	// New fields to update user with
	oldPassword := getFakePassword()
	newPassword := getFakePassword()

	// Add existing user to database
	userSample := sampler.getUserSampleWithSpecificPassword(oldPassword)

	// Expected Result when updatePassword is ok
	expectedResult := userDom.Result{
		UserID:      userSample.GetID(),
		Username:    userSample.Username,
		DisplayName: userSample.DisplayName,
		Email:       userSample.Email,
	}

	// test updatePassword with missing mandatory field : UserID
	updatePassword := userDom.UpdatePassword{
		UserID:      0,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
	result, err := userService.UpdatePassword(updatePassword)
	require.Error(t, err, "Should return error but got nil")

	// test updatePassword with missing mandatory field : OldPassword
	updatePassword = userDom.UpdatePassword{
		UserID:      userSample.GetID(),
		OldPassword: "",
		NewPassword: newPassword,
	}
	result, err = userService.UpdatePassword(updatePassword)
	require.Error(t, err, "Should return error but got nil")

	// test updatePassword with wrong mandatory field : OldPassword
	updatePassword = userDom.UpdatePassword{
		UserID:      userSample.GetID(),
		OldPassword: "wrongPassword",
		NewPassword: newPassword,
	}
	result, err = userService.UpdatePassword(updatePassword)
	require.Error(t, err, "Should return error but got nil")

	// test updatePassword with missing mandatory field : NewPassword
	updatePassword = userDom.UpdatePassword{
		UserID:      userSample.GetID(),
		OldPassword: oldPassword,
		NewPassword: "",
	}
	result, err = userService.UpdatePassword(updatePassword)
	require.Error(t, err, "Should return error but got nil")

	// test updatePassword with all correct fields
	updatePassword = userDom.UpdatePassword{
		UserID:      userSample.GetID(),
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	// check if return no err
	result, err = userService.UpdatePassword(updatePassword)
	require.NoError(t, err, "Should not return error but got %v", err)

	// check if result is like expected
	require.Equal(t, expectedResult, result)

	// check if password has been correctly salt and hash
	hashedPassword, err := userRepository.GetHashedPassword(userDom.Get{UserID: result.UserID})
	require.NoError(t, err, "Should not return error but got %v", err)
	assert.NotEqual(t, updatePassword.OldPassword, hashedPassword, "password should be updated and hashed")
	assert.NotEqual(t, updatePassword.NewPassword, hashedPassword, "password should be hashed")
	err = utils.CompareHashAndPassword(hashedPassword, updatePassword.NewPassword)
	require.NoError(t, err, "passwords should be the same (using hash comparison) but got %v", err)
}

func TestUser_Delete(t *testing.T) {
	db, clean := OpenDatabase(t)
	defer clean()

	userRepository := repository.NewUserRepository(db)
	userService := userDom.NewService(userRepository)

	sampler := newSampler(t, db, false)

	// Add existing user to database
	userSample := sampler.getUserSample()

	// Check if user has been correctly created
	_, err := userService.Get(userDom.Get{UserID: userSample.GetID()})
	require.NoError(t, err, "Should not return error but got %v", err)

	// test update with missing mandatory field : UserID
	delete := userDom.Delete{
		UserID:      0,
	}
	err = userService.Delete(delete)
	require.Error(t, err, "Should return error but got nil")

	// test update with wrong mandatory field : UserID
	delete = userDom.Delete{
		UserID:      domain.IDType(99999999999999),
	}
	err = userService.Delete(delete)
	require.Error(t, err, "Should return error but got nil")

	// test update with all correct fields
	delete = userDom.Delete{
		UserID:      userSample.GetID(),
	}

	// check if return no err
	err = userService.Delete(delete)
	require.NoError(t, err, "Should not return error but got %v", err)
	_, err = userService.Get(userDom.Get{UserID: userSample.GetID()})
	require.Error(t, err, "Should return error because record must be deleted but got nil", err)
}

func TestUser_Get(t *testing.T) {
	db, clean := OpenDatabase(t)
	defer clean()

	userRepository := repository.NewUserRepository(db)
	userService := userDom.NewService(userRepository)

	sampler := newSampler(t, db, false)

	// Add existing user to database
	userSample := sampler.getUserSample()

	// Expected Result when updatePassword is ok
	expectedResult := userDom.Result{
		UserID:      userSample.GetID(),
		Username:    userSample.Username,
		DisplayName: userSample.DisplayName,
		Email:       userSample.Email,
	}

	// test update with all missing mandatory fields
	get := userDom.Get{
		UserID:   0,
		Username: "",
		Email:    "",
	}
	result, err := userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	// test update with wrong mandatory field : Username
	get = userDom.Get{
		UserID:   userSample.GetID(),
		Username: fake.UserName(),
		Email:    userSample.Email,
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	get = userDom.Get{
		Username: fake.UserName(),
		Email:    userSample.Email,
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	get = userDom.Get{
		Username: fake.UserName(),
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	// test update with wrong mandatory field : UserID
	get = userDom.Get{
		UserID:   domain.IDType(999999999999999999),
		Username: userSample.Username,
		Email:    userSample.Email,
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	get = userDom.Get{
		UserID:   domain.IDType(999999999999999999),
		Email:    userSample.Email,
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	get = userDom.Get{
		UserID:   domain.IDType(999999999999999999),
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	// test update with wrong mandatory field : Email
	get = userDom.Get{
		UserID:   userSample.GetID(),
		Username: userSample.Username,
		Email:    fake.EmailAddress(),
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	get = userDom.Get{
		UserID:   userSample.GetID(),
		Email:    fake.EmailAddress(),
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	get = userDom.Get{
		Email:    fake.EmailAddress(),
	}
	result, err = userService.Get(get)
	require.Error(t, err, "Should return error but got nil")

	// test update with all correct fields
	get = userDom.Get{
		UserID:   userSample.GetID(),
		Username: userSample.Username,
		Email:    userSample.Email,
	}

	// check if return no err
	result, err = userService.Get(get)
	require.NoError(t, err, "Should not return error but got %v", err)
	require.Equal(t, expectedResult, result)
}

func TestUser_Search(t *testing.T) {
	db, clean := OpenDatabase(t)
	defer clean()

	userRepository := repository.NewUserRepository(db)
	userService := userDom.NewService(userRepository)

	// Variables used for searching on username
	commonUsernamePart := fake.UserName()[:3]
	matchingUsername1 := fake.UserName() + commonUsernamePart + fake.UserName()
	matchingUsername2 := commonUsernamePart + fake.UserName()
	matchingUsername3 := fake.UserName() + commonUsernamePart

	// Variables used for searching on email
	commonEmailPart := fake.EmailAddress()[2:6]
	matchingEmail1 := fake.EmailAddress() + commonEmailPart + fake.EmailAddress()
	matchingEmail2 := commonEmailPart + fake.EmailAddress()
	matchingEmail3 := fake.EmailAddress() + commonEmailPart

	// Variables used for searching on displayName
	commonDisplayNamePart := fake.FullName()[3:9]
	matchingDisplayName1 := fake.FullName() + commonDisplayNamePart + fake.FullName()
	matchingDisplayName2 := commonDisplayNamePart + fake.FullName()
	matchingDisplayName3 := fake.FullName() + commonDisplayNamePart

	// Create some users using matching and not matching fields
	users := []repository.User {
		{
			Username:       matchingUsername2,
			DisplayName:    matchingDisplayName3,
			Email:          fake.EmailAddress(),
		},
		{
			Username:       matchingUsername1,
			DisplayName:    fake.FullName(),
			Email:          matchingEmail3,
		},
		{
			Username:       fake.UserName(),
			DisplayName:    matchingDisplayName1,
			Email:          matchingEmail1,
		},
		{
			Username:       matchingUsername3,
			DisplayName:    matchingDisplayName2,
			Email:          matchingEmail2,
		},
		{
			Username:       fake.UserName(),
			DisplayName:    fake.FullName(),
			Email:          fake.EmailAddress(),
		},
	}
	// Save users into database
	err := db.Create(&users).Error
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check with wrong length of keyword (<3) : should return error
	Filters := userDom.Filters{
		Keyword: "a",
	}
	result, err := userService.Search(Filters)
	require.Error(t, err, "Should return error because keyword is too small")
	require.Len(t, result, 0)

	// Check with keyword matching nothing : should return not error and length of 0
	Filters = userDom.Filters{
		Keyword: fake.UserName(),
	}
	result, err = userService.Search(Filters)
	require.NoError(t, err, "Should return not error")
	require.Len(t, result, 0)

	// Check with keyword matching username : should return not error and length of 3
	Filters = userDom.Filters{
		Keyword: commonUsernamePart,
	}
	result, err = userService.Search(Filters)
	require.NoError(t, err, "Should return not error")
	require.Len(t, result, 3)

	// Check with keyword matching email : should return not error and length of 3
	Filters = userDom.Filters{
		Keyword: commonEmailPart,
	}
	result, err = userService.Search(Filters)
	require.NoError(t, err, "Should return not error")
	require.Len(t, result, 3)

	// Check with keyword matching displayName : should return not error and length of 3
	Filters = userDom.Filters{
		Keyword: commonDisplayNamePart,
	}
	result, err = userService.Search(Filters)
	require.NoError(t, err, "Should return not error")
	require.Len(t, result, 3)
}