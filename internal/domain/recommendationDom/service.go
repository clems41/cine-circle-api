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
	List(filters Filters) (list ViewList, err error)
	ListUsers(usersFilters UsersFilters) (list UserList, err error)
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

func (service *service) List(filters Filters) (list ViewList, err error) {
	err = filters.Valid()
	if err != nil {
		return
	}

	// Update filters with right fields
	if filters.Field == "date" {
		filters.Field = "created_at"
	}

	recommendations, total, err := service.r.List(filters)
	if err != nil {
		return
	}
	list.Page = filters.PaginationRequest.BuildResult(total)
	for _, recommendation := range recommendations {
		list.Recommendations = append(list.Recommendations, service.toView(recommendation, filters.UserID))
	}
	return
}

func (service *service) ListUsers(usersFilters UsersFilters) (list UserList, err error) {
	users, total, err := service.r.ListUsers(usersFilters)
	if err != nil {
		return
	}
	list.Page = usersFilters.PaginationRequest.BuildResult(total)
	for _, user := range users {
		list.Users = append(list.Users, UserView{
			UserID:      user.GetID(),
			Username:    *user.Username,
			DisplayName: user.DisplayName,
		})
	}
	return
}

func (service *service) toView(recommendation repositoryModel.Recommendation, senderID uint) (view RecommendationView) {
	view = RecommendationView{
		ID:                 recommendation.GetID(),
		Date:               recommendation.CreatedAt,
		Sender:             UserView{
			UserID:      recommendation.Sender.GetID(),
			Username:    *recommendation.Sender.Username,
			DisplayName: recommendation.Sender.DisplayName,
		},
		Movie:              MovieView{
			ID:          recommendation.Movie.GetID(),
			Title:       recommendation.Movie.Title,
			PosterPath:  recommendation.Movie.PosterPath,
			ReleaseDate: recommendation.Movie.ReleaseDate,
		},
		Comment:            recommendation.Comment,
		Circles:            nil,
		Users:              nil,
	}
	if senderID == recommendation.SenderID {
		view.RecommendationType = recommendationSent
	} else {
		view.RecommendationType = recommendationReceived
	}
	for _, user := range recommendation.Users {
		view.Users = append(view.Users, UserView{
			UserID:      user.GetID(),
			Username:    *user.Username,
			DisplayName: user.DisplayName,
		})
	}
	for _, circle := range recommendation.Circles {
		circleView := CircleView{
			CircleID:    circle.GetID(),
			Name:        circle.Name,
			Description: circle.Description,
		}
		for _, user := range circle.Users {
			circleView.Users = append(circleView.Users, UserView{
				UserID:      user.GetID(),
				Username:    *user.Username,
				DisplayName: user.DisplayName,
			})
		}
		view.Circles = append(view.Circles, circleView)
	}
	return
}
