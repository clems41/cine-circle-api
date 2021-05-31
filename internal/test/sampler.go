package test

import (
	"cine-circle/internal/constant"
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/utils"
	"github.com/google/uuid"
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
	users := sampler.GetUsers(NumberOfUsersToPopulateDatabase)

	// populate database with some circles (between 1 and 5)
	for range FakeRange(1, 5) {
		// Add user into circle (1 chance over 2)
		if FakeBool() {
			element := RandomElement(users)
			user, ok := element.(repositoryModel.User)
			if !ok {
				sampler.t.Fatalf("Element should be type User")
			}
			sampler.GetCircle(user)
		} else {
			sampler.GetCircle()
		}
	}

	// populate database with some movies (between 10 and 20)
	var movies []repositoryModel.Movie
	for range FakeRange(10, 20) {
		movies = append(movies, *sampler.GetMovie())
	}

	// populate databse with some recommendations (between 5 and 15)
	// Adding existing user into recommendation
	element := RandomElement(users)
	user, ok := element.(repositoryModel.User)
	if !ok {
		sampler.t.Fatalf("Element should be type User")
	}
	// Adding recommendation sent by user 1 time over 2
	if FakeBool() {
		sampler.GetRecommendationsSentByUser(&user)
	} else {
		sampler.GetRecommendationsReceivedByUser(&user)
	}
}

func (sampler *Sampler) GetUser() (user *repositoryModel.User) {
	// HashAndSalt password for user
	password := FakePassword()
	return sampler.GetUserWithSpecificPassword(password)
}

func (sampler *Sampler) GetUsers(numberOfUsers int) (users []repositoryModel.User) {
	// HashAndSalt password for user
	for i := 0; i < numberOfUsers; i++ {
		users = append(users, *sampler.GetUser())
	}
	return
}

func (sampler *Sampler) GetUserWithSpecificPassword(password string) (user *repositoryModel.User) {

	hashedPassword, err := utils.HashAndSaltPassword(password, constant.CostHashFunction)
	require.NoError(sampler.t, err)

	// Create new user
	username := strings.ToLower(fake.UserName() + fake.UserName())
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

func (sampler *Sampler) GetCircle(users ...repositoryModel.User) *repositoryModel.Circle {

	circle := repositoryModel.Circle{
		Name:        fake.Title(),
		Description: fake.Sentences(),
	}

	// Adding users
	for range FakeRange(4, 12) {
		circle.Users = append(circle.Users, *sampler.GetUser())
	}

	// Adding specific users
	circle.Users = append(circle.Users, users...)

	err := sampler.DB.
		Create(&circle).
		Error
	require.NoError(sampler.t, err)

	return &circle
}

func (sampler *Sampler) GetMovie() *repositoryModel.Movie {
	movie := repositoryModel.Movie{
		Title:            fake.Title(),
		ImdbId:           uuid.New().String(),
		BackdropPath:     fake.Street(),
		PosterPath:       fake.Street(),
		Genres:           fake.GetLangs(),
		OriginalLanguage: fake.Language(),
		OriginalTitle:    fake.Title(),
		Overview:         fake.Sentences(),
		ReleaseDate:      FakeTime(),
		Runtime:          FakeIntBetween(55, 236),
	}

	err := sampler.DB.
		Create(&movie).
		Error
	require.NoError(sampler.t, err)

	return &movie
}

func (sampler *Sampler) GetRecommendationsSentByUser(sender *repositoryModel.User) (list []repositoryModel.Recommendation) {
	for range FakeRange(4, 12) {
		movie := sampler.GetMovie()
		recommendation := repositoryModel.Recommendation{
			SenderID: sender.GetID(),
			MovieID:  movie.GetID(),
			Comment:  fake.Sentences(),
		}
		err := sampler.DB.Create(&recommendation).Error
		require.NoError(sampler.t, err)
		// Add users to sent recommendations
		for range FakeRange(1, 3) {
			user := sampler.GetUser()
			err = sampler.DB.Exec("INSERT INTO recommendation_user (recommendation_id,user_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), user.GetID()).Error
			require.NoError(sampler.t, err)
		}
		// Add circles to sent recommendations
		for range FakeRange(1, 2) {
			circle := sampler.GetCircle()
			err = sampler.DB.Exec("INSERT INTO recommendation_circle (recommendation_id,circle_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), circle.GetID()).Error
			require.NoError(sampler.t, err)
		}
		err = sampler.DB.
			Preload("Users").
			Preload("Movie").
			Preload("Sender").
			Preload("Circles").
			Preload("Circles.Users").
			Order("id").
			Take(&recommendation).
			Error
		require.NoError(sampler.t, err)
		list = append(list, recommendation)
	}
	return
}

func (sampler *Sampler) GetRecommendationsReceivedByUser(recipient *repositoryModel.User) (list []repositoryModel.Recommendation) {
	for range FakeRange(4, 12) {
		sender := sampler.GetUser()
		movie := sampler.GetMovie()
		recommendation := repositoryModel.Recommendation{
			SenderID: sender.GetID(),
			MovieID:  movie.GetID(),
			Comment:  fake.Sentences(),
		}
		err := sampler.DB.Create(&recommendation).Error
		require.NoError(sampler.t, err)
		// Create fake bool for choosing if user should be added into circles list or users list
		recipientShouldBeAddedIntoCircle := FakeBool()
		if !recipientShouldBeAddedIntoCircle {
			err = sampler.DB.Exec("INSERT INTO recommendation_user (recommendation_id,user_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), recipient.GetID()).Error
			require.NoError(sampler.t, err)
		} else {
			circle := sampler.GetCircle(*recipient)
			err = sampler.DB.Exec("INSERT INTO recommendation_circle (recommendation_id,circle_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), circle.GetID()).Error
			require.NoError(sampler.t, err)
		}
		// Add users to sent recommendations
		for range FakeRange(1, 3) {
			user := sampler.GetUser()
			err = sampler.DB.Exec("INSERT INTO recommendation_user (recommendation_id,user_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), user.GetID()).Error
			require.NoError(sampler.t, err)
		}
		// Add circles to sent recommendations
		for range FakeRange(1, 2) {
			circle := sampler.GetCircle()
			err = sampler.DB.Exec("INSERT INTO recommendation_circle (recommendation_id,circle_id) VALUES (?,?) ON CONFLICT DO NOTHING", recommendation.GetID(), circle.GetID()).Error
			require.NoError(sampler.t, err)
		}
		err = sampler.DB.
			Preload("Users").
			Preload("Movie").
			Preload("Sender").
			Preload("Circles").
			Preload("Circles.Users").
			Order("id").
			Take(&recommendation).
			Error
		require.NoError(sampler.t, err)
		list = append(list, recommendation)
	}
	return
}

func (sampler *Sampler) GetRecommendations() (list []repositoryModel.Recommendation) {
	for range FakeRange(1, 3) {
		sender := sampler.GetUser()
		list = append(list, sampler.GetRecommendationsSentByUser(sender)...)
	}
	return
}
