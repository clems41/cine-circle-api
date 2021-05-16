package repository

import (
	"cine-circle/internal/domain/userDom"
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

	User						*userDom.Repository
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

	return
}

