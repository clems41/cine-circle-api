package postgresRepositories

import (
	"cine-circle-api/internal/model"
	"cine-circle-api/internal/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ repository.Media = (*mediaPgRepository)(nil)
var _ PgRepository = (*mediaPgRepository)(nil)

type mediaPgRepository struct {
	DB *gorm.DB
}

func NewMedia(DB *gorm.DB) *mediaPgRepository {
	return &mediaPgRepository{DB: DB}
}

func (repo *mediaPgRepository) Migrate() (err error) {
	err = repo.DB.
		AutoMigrate(&model.Media{})
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (repo *mediaPgRepository) Get(mediaID uint) (media model.Media, ok bool, err error) {
	err = repo.DB.
		Take(&media, mediaID).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return media, false, nil
		}
	} else {
		ok = true
	}
	return
}

func (repo *mediaPgRepository) Save(media *model.Media) (err error) {
	err = repo.DB.Save(media).Error
	return
}

func (repo *mediaPgRepository) GetFromProvider(mediaProviderName, mediaProviderId string) (media model.Media, ok bool, err error) {
	err = repo.DB.
		Take(&media, "media_provider_name = ? AND media_provider_id = ?", mediaProviderName, mediaProviderId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return media, false, nil
		} else {
			return media, false, errors.WithStack(err)
		}
	}
	return media, true, nil
}
