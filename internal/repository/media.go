package repository

import (
	"cine-circle-api/internal/model"
)

type Media interface {
	Get(mediaID uint) (media model.Media, ok bool, err error)
	Save(media *model.Media) (err error)
	GetFromProvider(mediaProviderName, mediaProviderId string) (media model.Media, ok bool, err error)
}
