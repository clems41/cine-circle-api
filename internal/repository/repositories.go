package repository

import (
	"cine-circle/internal/logger"
	"gorm.io/gorm"
)

type DB interface {
	DB() *gorm.DB
	PrefixTables(prefix string)
}

type Repository interface {
	Migrate()
}

type Repositories struct {
	List []Repository

	Circle						*circleRepository
	Movie						*movieRepository
	Recommendation				*recommendationRepository
	User						*userRepository
	Watchlist					*watchlistRepository
}

func (rs Repositories) Migrate() {

	for _, r := range rs.List {
		if repo, ok := r.(Repository); ok {
			repo.Migrate()
		} else {
			logger.Sugar.Fatalf("not a repository %s", repo)
		}
	}
}

func NewAllRepositories(DB *gorm.DB) (r Repositories) {

	r.Circle = NewCircleRepository(DB)
	r.List = append(r.List, r.Circle)

	r.Movie = NewMovieRepository(DB)
	r.List = append(r.List, r.Movie)

	r.Recommendation = NewRecommendationRepository(DB)
	r.List = append(r.List, r.Recommendation)

	r.User = NewUserRepository(DB)
	r.List = append(r.List, r.User)

	r.Watchlist = NewWatchlistRepository(DB)
	r.List = append(r.List, r.Watchlist)

	return
}

