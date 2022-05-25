package recommendationDom

import (
	"cine-circle-api/internal/constant/recommendationConst"
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository"
	"cine-circle-api/pkg/utils/searchUtils"
	"time"
)

var _ Service = (*service)(nil)

type Service interface {
	Send(form SendForm) (view SendView, err error)
	Search(form SearchForm) (view SearchView, err error)
}

type service struct {
	repository       repository.Recommendation
	userRepository   repository.User
	mediaRepository  repository.Media
	circleRepository repository.Circle
}

func NewService(repository repository.Recommendation, userRepository repository.User,
	mediaRepository repository.Media, circleRepository repository.Circle) Service {
	return &service{
		repository:       repository,
		userRepository:   userRepository,
		mediaRepository:  mediaRepository,
		circleRepository: circleRepository,
	}
}

func (svc *service) Send(form SendForm) (view SendView, err error) {
	media, ok, err := svc.mediaRepository.Get(form.MediaId)
	if err != nil {
		return
	}
	if !ok {
		return view, errMediaNotFound
	}
	var circles []model.Circle
	for _, circleId := range form.CirclesIds {
		var circle model.Circle
		circle, ok, err = svc.circleRepository.Get(circleId)
		if err != nil {
			return
		}
		if !ok {
			return view, errCircleNotFound
		}
		circles = append(circles, circle)
	}
	sender, ok, err := svc.userRepository.Get(form.SenderId)
	if err != nil {
		return
	}
	if !ok {
		return view, errMediaNotFound
	}

	recommendation := model.Recommendation{
		SenderID: form.SenderId,
		Sender:   sender,
		Circles:  circles,
		Media:    media,
		MediaID:  media.ID,
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
	repoForm := repository.RecommendationSearchForm{
		PaginationRequest: searchUtils.PaginationRequest{
			Page:     form.Page,
			PageSize: form.PageSize,
		},
		UserId:  form.UserId,
		MediaId: uint(form.MediaId),
		Type:    form.Type,
	}

	recommendations, total, err := svc.repository.Search(repoForm)
	if err != nil {
		return
	}

	view.Page = form.BuildResult(total)
	view.Recommendations = make([]CommonView, 0)

	for _, recommendation := range recommendations {
		view.Recommendations = append(view.Recommendations, svc.fromModelToView(recommendation, form.UserId))
	}

	return
}

/* Private methods below */

func (svc *service) fromModelToView(recommendation model.Recommendation, userId uint) (view CommonView) {
	view = CommonView{
		Id: recommendation.ID,
		Sender: UserView{
			Id:        recommendation.Sender.ID,
			Firstname: recommendation.Sender.FirstName,
			Lastname:  recommendation.Sender.LastName,
			Username:  recommendation.Sender.Username,
		},
		Circles: nil,
		Media: MediaView{
			Id:            recommendation.Media.ID,
			Title:         recommendation.Media.Title,
			BackdropUrl:   recommendation.Media.BackdropUrl,
			Genres:        recommendation.Media.Genres,
			Language:      recommendation.Media.Language,
			OriginalTitle: recommendation.Media.OriginalTitle,
			Overview:      recommendation.Media.Overview,
			PosterUrl:     recommendation.Media.PosterUrl,
			ReleaseDate:   recommendation.Media.ReleaseDate,
			Runtime:       recommendation.Media.Runtime,
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
