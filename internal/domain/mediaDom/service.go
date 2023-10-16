package mediaDom

import (
	"cine-circle-api/external/mediaProvider"
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository"
)

var _ Service = (*service)(nil)

type Service interface {
	Get(form GetForm) (view GetView, err error)
	Search(form SearchForm) (view SearchView, err error)
}

type service struct {
	mediaProvider mediaProvider.Service
	repository    repository.Media
}

func NewService(mediaProvider mediaProvider.Service, repository repository.Media) Service {
	return &service{
		repository:    repository,
		mediaProvider: mediaProvider,
	}
}

func (svc *service) Get(form GetForm) (view GetView, err error) {
	// Media should be already in database (even marked as uncompleted)
	movie, ok, err := svc.repository.Get(form.MediaId)
	if err != nil {
		return
	}
	if !ok {
		return view, errMediaNotFound
	}

	// If movie is not completed, fill info from mediaProvider, then mark it as completed
	if !movie.Completed {
		var movieFromMediaProvider mediaProvider.Media
		movieFromMediaProvider, err = svc.mediaProvider.Get(movie.MediaProviderId)
		if err != nil {
			return
		}
		movie.Title = movieFromMediaProvider.Title
		movie.PosterUrl = movieFromMediaProvider.PosterUrl
		movie.Overview = movieFromMediaProvider.Overview
		movie.OriginalTitle = movieFromMediaProvider.OriginalTitle
		movie.Language = movieFromMediaProvider.Language
		movie.BackdropUrl = movieFromMediaProvider.BackdropUrl
		movie.Runtime = movieFromMediaProvider.Runtime
		movie.ReleaseDate = movieFromMediaProvider.ReleaseDate
		movie.Genres = movieFromMediaProvider.Genres
		movie.Completed = true
		err = svc.repository.Save(&movie)
		if err != nil {
			return
		}
	}

	// Fill view
	view = svc.fromModelToView(movie)
	return
}

func (svc *service) Search(form SearchForm) (view SearchView, err error) {
	// Call mediaProvider to get result
	formMediaProvider := mediaProvider.SearchForm{
		Page:    form.Page,
		Keyword: form.Keyword,
	}
	medias, total, err := svc.mediaProvider.Search(formMediaProvider)
	if err != nil {
		return
	}

	// Save all movie result into database and marked them as not completed and fill view
	for _, resultMedia := range medias {
		// Create only if not already exists, if already exists get previous ID to add it into result
		var alreadyExists bool
		var media model.Media
		media, alreadyExists, err = svc.repository.GetFromProvider(svc.mediaProvider.GetProviderName(), resultMedia.Id)
		if err != nil {
			return
		}
		if !alreadyExists {
			// Stored new movie from research if not already exists
			media = model.Media{
				MediaProviderName: svc.mediaProvider.GetProviderName(),
				MediaProviderId:   resultMedia.Id,
				Completed:         false,
			}
			err = svc.repository.Save(&media)
			if err != nil {
				return
			}
		}
		view.Result = append(view.Result, ResultView{
			Id:            media.ID,
			Title:         resultMedia.Title,
			Language:      resultMedia.Language,
			OriginalTitle: resultMedia.OriginalTitle,
			PosterUrl:     resultMedia.PosterUrl,
		})
	}

	// Fill other view fields
	view.Page = form.BuildResult(total)

	return
}

/* PRIVATE METHODS */

func (svc *service) fromModelToView(movie model.Media) (view GetView) {
	view = GetView{
		Id:            movie.ID,
		Title:         movie.Title,
		BackdropUrl:   movie.BackdropUrl,
		Genres:        movie.Genres,
		Language:      movie.Language,
		OriginalTitle: movie.OriginalTitle,
		Overview:      movie.Overview,
		PosterUrl:     movie.PosterUrl,
		ReleaseDate:   movie.ReleaseDate,
		Runtime:       movie.Runtime,
	}
	return
}
