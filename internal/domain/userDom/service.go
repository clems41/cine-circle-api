package userDom

var _ Service = (*service)(nil)

type Service interface {
}

type service struct {
	r Repository
}

type Repository interface {
}

func NewService(r Repository) Service {
	return &service{
		r:                              r,
	}
}

