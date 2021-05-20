package recommendationDom

import (
	"cine-circle/internal/repository/repositoryModel"
	"cine-circle/internal/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ Service = (*service)(nil)

type Service interface {
	Create(creation Creation) (err error)
}

type service struct {
	r repository
}

func NewService(r repository) Service {
	return &service{
		r: r,
	}
}

func (service *service) Create(creation Creation) (err error) {
	err = creation.Valid()
	if err != nil {
		return
	}

	// Check if movie exists
	exists, err := service.r.CheckIfMovieExists(creation.MovieID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errMovieNotFound
		} else {
			return
		}
	}
	if !exists {
		return errMovieNotFound
	}

	// User need to be in each circle to send recommendation
	for _, circleId := range creation.CircleIDs {
		var userIDs []uint
		userIDs, err = service.r.GetUserIDsFromCircle(circleId)
		if err != nil {
			return
		}
		if !utils.ContainsID(userIDs, creation.SenderID) {
			return errUserUnauthorized
		}
	}

	// User need to be in touch with each user to send recommendation
	var userIDsCloseTo []uint
	userIDsCloseTo, err = service.r.GetUserIDsCloseToUser(creation.SenderID)
	if err != nil {
		return
	}
	for _, userID := range creation.UserIDs {
		if !utils.ContainsID(userIDsCloseTo, userID) {
			return errUserUnauthorized
		}
	}

	// Save recommendation
	var users []repositoryModel.User
	for _, userId := range creation.UserIDs {
		var user repositoryModel.User
		user.SetID(userId)
		users = append(users, user)
	}
	var circles []repositoryModel.Circle
	for _, circleID := range creation.CircleIDs {
		var circle repositoryModel.Circle
		circle.SetID(circleID)
		circles = append(circles, circle)
	}
	recommendation := repositoryModel.Recommendation{
		SenderID: creation.SenderID,
		MovieID:  creation.MovieID,
		Comment:  creation.Comment,
		Circles:  circles,
		Users:    users,
	}
	return service.r.Create(&recommendation)
}
