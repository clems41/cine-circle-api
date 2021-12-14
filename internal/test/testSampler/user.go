package testSampler

import (
	"cine-circle-api/internal/repository/model"
	"cine-circle-api/pkg/utils/securityUtils"
	"cine-circle-api/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
	"strings"
)

// GetUser retourne un user créé en DB
func (sampler *Sampler) GetUser() (user *model.User) {
	password := fakeData.Password()
	return sampler.GetUserWithPassword(password)
}

// GetUserWithEmailToken retourne un user créé en DB avec emailToken rempli comme s'il y a avait eu une demande de validation d'email
func (sampler *Sampler) GetUserWithEmailToken() (user *model.User) {
	password := fakeData.Password()
	hashedPassword, err := securityUtils.HashAndSaltPassword(password)
	require.NoError(sampler.t, err)

	user = &model.User{
		Username:       strings.ToLower(fakeData.UniqueUsername()),
		HashedPassword: hashedPassword,
		LastName:       fake.LastName(),
		FirstName:      fake.FirstName(),
		Email:          fakeData.UniqueEmail(),
		Role:           "",
		Active:         true,
		EmailToken:     "toto",
	}
	err = sampler.DB.Create(user).Error
	require.NoError(sampler.t, err)
	return
}

// GetUserWithPasswordToken retourne un user créé en DB avec PasswordToken rempli comme s'il y a avait eu une demande de réinitialisation de mot de passe
func (sampler *Sampler) GetUserWithPasswordToken() (user *model.User) {
	password := fakeData.Password()
	hashedPassword, err := securityUtils.HashAndSaltPassword(password)
	require.NoError(sampler.t, err)

	user = &model.User{
		Username:       strings.ToLower(fakeData.UniqueUsername()),
		HashedPassword: hashedPassword,
		LastName:       fake.LastName(),
		FirstName:      fake.FirstName(),
		Email:          fakeData.UniqueEmail(),
		Role:           "",
		Active:         true,
		PasswordToken:  "toto",
	}
	err = sampler.DB.Create(user).Error
	require.NoError(sampler.t, err)
	return
}

// GetUserWithPassword retourne un user créé en DB pour un mdp précis
func (sampler *Sampler) GetUserWithPassword(password string) (user *model.User) {
	hashedPassword, err := securityUtils.HashAndSaltPassword(password)
	require.NoError(sampler.t, err)

	user = &model.User{
		Username:       strings.ToLower(fakeData.UniqueUsername()),
		HashedPassword: hashedPassword,
		LastName:       fake.LastName(),
		FirstName:      fake.FirstName(),
		Email:          fakeData.UniqueEmail(),
		Active:         true,
	}
	err = sampler.DB.Create(user).Error
	require.NoError(sampler.t, err)
	return
}

// GetUserWithFirstName retourne un user créé en DB pour un firstname précis
func (sampler *Sampler) GetUserWithFirstName(firstName string) (user *model.User) {
	password := fakeData.Password()
	hashedPassword, err := securityUtils.HashAndSaltPassword(password)
	require.NoError(sampler.t, err)

	user = &model.User{
		Username:       strings.ToLower(fakeData.UniqueUsername()),
		HashedPassword: hashedPassword,
		LastName:       fake.LastName(),
		FirstName:      firstName,
		Email:          fakeData.UniqueEmail(),
		Active:         true,
	}
	err = sampler.DB.Create(user).Error
	require.NoError(sampler.t, err)
	return
}

// GetUserWithLastName retourne un user créé en DB pour un firstname précis
func (sampler *Sampler) GetUserWithLastName(lastName string) (user *model.User) {
	password := fakeData.Password()
	hashedPassword, err := securityUtils.HashAndSaltPassword(password)
	require.NoError(sampler.t, err)

	user = &model.User{
		Username:       strings.ToLower(fakeData.UniqueUsername()),
		HashedPassword: hashedPassword,
		LastName:       lastName,
		FirstName:      fake.FirstName(),
		Email:          fakeData.UniqueEmail(),
		Active:         true,
	}
	err = sampler.DB.Create(user).Error
	require.NoError(sampler.t, err)
	return
}

// GetUserWithUsername retourne un user créé en DB pour un firstname précis
func (sampler *Sampler) GetUserWithUsername(username string) (user *model.User) {
	password := fakeData.Password()
	hashedPassword, err := securityUtils.HashAndSaltPassword(password)
	require.NoError(sampler.t, err)

	user = &model.User{
		Username:       username,
		HashedPassword: hashedPassword,
		LastName:       fake.LastName(),
		FirstName:      fake.FirstName(),
		Email:          fakeData.UniqueEmail(),
		Active:         true,
	}
	err = sampler.DB.Create(user).Error
	require.NoError(sampler.t, err)
	return
}

// GetUserWithEmail retourne un user créé en DB pour un firstname précis
func (sampler *Sampler) GetUserWithEmail(email string) (user *model.User) {
	password := fakeData.Password()
	hashedPassword, err := securityUtils.HashAndSaltPassword(password)
	require.NoError(sampler.t, err)

	user = &model.User{
		Username:       strings.ToLower(fakeData.UniqueUsername()),
		HashedPassword: hashedPassword,
		LastName:       fake.LastName(),
		FirstName:      fake.FirstName(),
		Email:          email,
		Active:         true,
	}
	err = sampler.DB.Create(user).Error
	require.NoError(sampler.t, err)
	return
}
