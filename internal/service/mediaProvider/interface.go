package mediaProvider

type Service interface {
	Search(form SearchForm) (view SearchView, err error)
	Get(form MediaForm) (view MediaView, err error)
}
