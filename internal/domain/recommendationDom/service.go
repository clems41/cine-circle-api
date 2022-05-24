package recommendationDom

import (
	"cine-circle-api/internal/constant/recommendationConst"
	"cine-circle-api/internal/repository"
	"cine-circle-api/internal/repository/instance/recommendationRepository"
	"cine-circle-api/internal/repository/postgres/pgModel"
	"cine-circle-api/pkg/sql/gormUtils"
	"time"
)

var _ Service = (*service)(nil)

type Service interface {
	Send(form SendForm) (view SendView, err error)
	Search(form SearchForm) (view SearchView, err error)
}

type service struct {
	repository       repository.Repository
	userRepository   repository.Repository
	mediaRepository  repository.Repository
	circleRepository repository.Repository
}

func NewService(repository repository.Repository, userRepository repository.Repository,
	mediaRepository repository.Repository, circleRepository repository.Repository) Service {
	return &service{
		repository:       repository,
		userRepository:   userRepository,
		mediaRepository:  mediaRepository,
		circleRepository: circleRepository,
	}
}

func (svc *service) Send(form SendForm) (view SendView, err error) {
	movie, ok, err := svc.mediaRepository.GetMovie(form.MediaId)
	if err != nil {
		return
	}
	if !ok {
		return view, errMediaNotFound
	}
	var circles []pgModel.Circle
	for _, circleId := range form.CirclesIds {
		var circle pgModel.Circle
		circle, ok, err = svc.circleRepository.Get(circleId)
		if err != nil {
			return
		}
		if !ok {
			return view, errCircleNotFound
		}
		circles = append(circles, circle)
	}
	recommendation := pgModel.Recommendation{
		SenderId: form.SenderId,
		Circles:  circles,
		MovieId:  movie.ID,
		Text:     form.Text,
		Date:     time.Now(),
	}
	err = svc.repository.Create(&recommendation)
	if err != nil {
		return
	}

	view.CommonView = svc.fromModelToView(recommendation, form.SenderId)
	return
}

func (svc *service) Search(form SearchForm) (view SearchView, err error) {
	repoForm := recommendationRepository.SearchForm{
		PaginationQuery: gormUtils.PaginationQuery{
			Page:     form.Page,
			PageSize: form.PageSize,
		},
		UserId:  form.UserId,
		MovieId: uint(form.MovieId),
		Type:    form.Type,
	}

	repoView, err := svc.repository.Search(repoForm)
	if err != nil {
		return
	}

	view.Page = form.BuildResult(repoView.Total)
	view.Recommendations = make([]CommonView, 0)

	for _, recommendation := range repoView.Recommendations {
		view.Recommendations = append(view.Recommendations, svc.fromModelToView(recommendation, form.UserId))
	}

	return
}

/* Private methods below */

func (svc *service) fromModelToView(recommendation pgModel.Recommendation, userId uint) (view CommonView) {
	view = CommonView{
		Id: recommendation.ID,
		Sender: UserView{
			Id:        recommendation.Sender.ID,
			Firstname: recommendation.Sender.FirstName,
			Lastname:  recommendation.Sender.LastName,
			Username:  recommendation.Sender.Username,
		},
		Circles: nil,
		Movie: MovieView{
			Id:            recommendation.Movie.ID,
			Title:         recommendation.Movie.Title,
			BackdropUrl:   recommendation.Movie.BackdropUrl,
			Genres:        recommendation.Movie.Genres,
			Language:      recommendation.Movie.Language,
			OriginalTitle: recommendation.Movie.OriginalTitle,
			Overview:      recommendation.Movie.Overview,
			PosterUrl:     recommendation.Movie.PosterUrl,
			ReleaseDate:   recommendation.Movie.ReleaseDate,
			Runtime:       recommendation.Movie.Runtime,
		},
		Text: recommendation.Text,
		Date: recommendation.Date,
	}

	for _, circle := range recommendation.Circles {
		circleView := CircleView{
			Id:          circle.ID,
			Name:        circle.Name,
			Description: circle.Description,
		}
		for _, user := range circle.Users {
			circleView.Users = append(circleView.Users, UserView{
				Id:        user.ID,
				Firstname: user.FirstName,
				Lastname:  user.LastName,
				Username:  user.Username,
			})
		}
		view.Circles = append(view.Circles, circleView)
	}

	if recommendation.Sender.ID == userId {
		view.Type = recommendationConst.SentType
	} else {
		view.Type = recommendationConst.ReceivedType
	}

	return
}
