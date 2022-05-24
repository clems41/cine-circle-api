package mediaProvider

type Service interface {
	Search(form SearchForm) (view SearchView, err error)
	Get(form MovieForm) (view MovieView, err error)
	GetProviderName() (name string) // useful to know which provider is used
}
