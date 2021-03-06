package repository

import (
	"cine-circle/internal/domain/circleDom"
	"cine-circle/internal/domain/libraryDom"
	"cine-circle/internal/domain/movieDom"
	"cine-circle/internal/domain/recommendationDom"
	"cine-circle/internal/domain/userDom"
	"cine-circle/internal/domain/watchlistDom"
	logger "cine-circle/pkg/logger"
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

	User           *userDom.Repository
	Circle         *circleDom.Repository
	Movie          *movieDom.Repository
	Recommendation *recommendationDom.Repository
	Watchlist      *watchlistDom.Repository
	Library        *libraryDom.Repository
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

	r.User = userDom.NewRepository(DB)
	r.List = append(r.List, r.User)

	r.Circle = circleDom.NewRepository(DB)
	r.List = append(r.List, r.Circle)

	r.Movie = movieDom.NewRepository(DB)
	r.List = append(r.List, r.Movie)

	r.Recommendation = recommendationDom.NewRepository(DB)
	r.List = append(r.List, r.Recommendation)

	r.Watchlist = watchlistDom.NewRepository(DB)
	r.List = append(r.List, r.Watchlist)

	r.Library = libraryDom.NewRepository(DB)
	r.List = append(r.List, r.Library)

	return
}
