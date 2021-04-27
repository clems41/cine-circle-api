package test

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository"
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

// newSampler instantiates a new sampler object able to generate random resources for testing purpose
func newSampler(t *testing.T, DB *gorm.DB, populateDatabase bool) (sampler Sampler) {

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
	sampler.getUsersSample(NumberOfUsersToPopulateDatabase)
}

func (sampler *Sampler) getUserSample() (user *repository.User) {
	// HashAndSalt password for user
	password := getFakePassword()
	return sampler.getUserSampleWithSpecificPassword(password)
}

func (sampler *Sampler) getUsersSample(numberOfUsers int) (users []repository.User) {
	// HashAndSalt password for user
	for i := 0; i < numberOfUsers; i++ {
		users = append(users, *sampler.getUserSample())
	}
	return
}

func (sampler *Sampler) getUserSampleWithSpecificPassword(password string) (user *repository.User) {

	hashedPassword, err := utils.HashAndSaltPassword(password, constant.CostHashFunction)

	// Create new user
	if err != nil {
		sampler.t.Fatalf(err.Error())
	}
	user = &repository.User{
		Username:       strings.ToLower(fake.UserName()),
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
