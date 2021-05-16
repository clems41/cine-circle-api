package circleDom

var _ Service = (*service)(nil)

type Service interface {
}

type service struct {
	r repository
}

func NewService(r repository) Service {
	return &service{
		r: r,
	}
}
