package mediaProvider

type Service interface {
	Search(form SearchForm) (medias []MediaShort, total int64, err error)
	Get(mediaId string) (view Media, err error)
	GetProviderName() (name string) // useful to know which provider is used
}
