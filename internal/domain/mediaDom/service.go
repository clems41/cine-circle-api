package mediaDom

import (
	mediaProvider2 "cine-circle-api/external/mediaProvider"
	"cine-circle-api/internal/repository"
	"cine-circle-api/internal/repository/postgres/pgModel"
)

var _ Service = (*service)(nil)

type Service interface {
	Get(form GetForm) (view GetView, err error)
	Search(form SearchForm) (view SearchView, err error)
}

type service struct {
	mediaProvider mediaProvider2.Service
	repository    repository.Repository
}

func NewService(mediaProvider mediaProvider2.Service, repository repository.Repository) Service {
	return &service{
		repository:    repository,
		mediaProvider: mediaProvider,
	}
}

func (svc *service) Get(form GetForm) (view GetView, err error) {
	// Media should be already in database (even marked as uncompleted)
	movie, ok, err := svc.repository.GetMovie(form.MediaId)
	if err != nil {
		return
	}
	if !ok {
		return view, errMediaNotFound
	}

	// If movie is not completed, fill info from mediaProvider, then mark it as completed
	if !movie.Completed {
		var movieFromMediaProvider mediaProvider2.MovieView
		movieFromMediaProvider, err = svc.mediaProvider.Get(mediaProvider2.MovieForm{Id: movie.MediaProviderId})
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
	formMediaProvider := mediaProvider2.SearchForm{
		Page:    form.Page,
		Keyword: form.Keyword,
	}
	result, err := svc.mediaProvider.Search(formMediaProvider)
	if err != nil {
		return
	}

	// Save all movie result into database and marked them as not completed and fill view
	for _, media := range result.Result {
		// Create only if not already exists, if already exists get previous ID to add it into result
		var alreadyExists bool
		var movie pgModel.Movie
		movie, alreadyExists, err = svc.repository.GetMovieFromProvider(svc.mediaProvider.GetProviderName(), media.Id)
		if err != nil {
			return
		}
		if !alreadyExists {
			// Stored new movie from research if not already exists
			movie = pgModel.Movie{
				MediaProviderName: svc.mediaProvider.GetProviderName(),
				MediaProviderId:   media.Id,
				Completed:         false,
			}
			err = svc.repository.Create(&movie)
			if err != nil {
				return
			}
		}
		view.Result = append(view.Result, ResultView{
			Id:            movie.ID,
			Title:         media.Title,
			Language:      media.Language,
			OriginalTitle: media.OriginalTitle,
			PosterUrl:     media.PosterUrl,
		})
	}

	// Fill other view fields
	view.NumberOfPages = result.NumberOfPages
	view.CurrentPage = result.CurrentPage
	view.NumberOfItems = result.NumberOfItems
	view.PageSize = form.PageSize

	return
}

/* PRIVATE METHODS */

func (svc *service) fromModelToView(movie pgModel.Movie) (view GetView) {
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
