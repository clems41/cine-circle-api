package test

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/utils"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
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
	require.NoError(sampler.t, err)

	// Create new user
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
	require.NoError(sampler.t, err)
	return
}

func (sampler *Sampler) GetCircle() *repositoryModel.Circle{

	circle := repositoryModel.Circle{
		Name:        fake.Title(),
		Description: fake.Sentences(),
	}

	nbUsers := FakeIntBetween(4, 12)
	for i := 0; i < nbUsers; i++ {
		circle.Users = append(circle.Users, *sampler.GetUserSample())
	}

	err := sampler.DB.
		Create(&circle).
		Error
	require.NoError(sampler.t, err)

	return &circle
}
