package test

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/utils"
	"github.com/icrowley/fake"
	"gorm.io/gorm"
	"strings"
	"testing"
)

type Sampler struct {
	t  *testing.T
	DB *gorm.DB
}

// NewSampler instantiates a new sampler object able to generate random resources for testing purpose
func NewSampler(t *testing.T, DB *gorm.DB, populateDatabase bool) (sampler Sampler) {

	sampler.t = t
	sampler.DB = DB

	if populateDatabase {
		sampler.populateDatabase()
	}

	return
}

// populateDatabase inserts some random resources into database
func (sampler *Sampler) populateDatabase() {
	// populateDatabase with some users
	sampler.GetUsersSample(NumberOfUsersToPopulateDatabase)
}

func (sampler *Sampler) GetUserSample() (user *repositoryModel.User) {
	// HashAndSalt password for user
	password := FakePassword()
	return sampler.GetUserSampleWithSpecificPassword(password)
}

func (sampler *Sampler) GetUsersSample(numberOfUsers int) (users []repositoryModel.User) {
	// HashAndSalt password for user
	for i := 0; i < numberOfUsers; i++ {
		users = append(users, *sampler.GetUserSample())
	}
	return
}

func (sampler *Sampler) GetUserSampleWithSpecificPassword(password string) (user *repositoryModel.User) {

	hashedPassword, err := utils.HashAndSaltPassword(password, constant.CostHashFunction)

	// Create new user
	if err != nil {
		sampler.t.Fatalf(err.Error())
	}
	username := strings.ToLower(fake.UserName())
	user = &repositoryModel.User{
		Username:       &username,
		DisplayName:    fake.FullName(),
		Email:          fake.EmailAddress(),
		HashedPassword: hashedPassword,
	}

	// Save user into database
	err = sampler.DB.
		Create(user).
		Error
	if err != nil {
		sampler.t.Fatalf(err.Error())
	}
	return
}
